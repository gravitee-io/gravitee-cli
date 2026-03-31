package application

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type listOptions struct {
	factory *factory.Factory
	query   string
	status  string
	order   string
	page    int
	perPage int
	all     bool
}

func newListCmd(f *factory.Factory) *cobra.Command {
	opts := &listOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List applications",
		Example: `  gio app list
  gio app list --query "Mobile" --order -updated_at
  gio app list --status ARCHIVED`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := opts.validate(); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.query, "query", "", "Search applications by name")
	cmd.Flags().StringVar(&opts.status, "status", "ACTIVE", "Filter by status: ACTIVE, ARCHIVED")
	cmd.Flags().StringVar(&opts.order, "order", "name", "Sort field: name, updated_at, -name, -updated_at")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")

	return cmd
}

var (
	validAppStatuses = []string{"ACTIVE", "ARCHIVED"}
	validAppOrders   = []string{"name", "updated_at", "-name", "-updated_at"}
)

func (o *listOptions) validate() error {
	if err := cmdutil.ValidateEnum(o.status, "status", validAppStatuses); err != nil {
		return err
	}

	if err := cmdutil.ValidateEnum(o.order, "order", validAppOrders); err != nil {
		return err
	}

	return nil
}

type paginatedResponse struct {
	Data       []json.RawMessage `json:"data"`
	Pagination struct {
		Page           int `json:"page"`
		PerPage        int `json:"perPage"`
		PageCount      int `json:"pageCount"`
		TotalCount     int `json:"totalCount"`
		PageItemsCount int `json:"pageItemsCount"`
	} `json:"pagination"`
}

func (o *listOptions) run() error {
	f := o.factory
	p := cmdutil.NewPrinter(f)

	if o.all {
		return o.fetchAll(f, p)
	}

	return o.fetchPage(f, p, o.page)
}

func (o *listOptions) buildQuery(page int) string {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("size", strconv.Itoa(o.perPage))

	if o.status != "" {
		q.Set("status", o.status)
	}

	if o.order != "" {
		q.Set("order", o.order)
	}

	if o.query != "" {
		q.Set("query", o.query)
	}

	return q.Encode()
}

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	path := cmdutil.V1EnvPath(f, fmt.Sprintf("applications/_paged?%s", o.buildQuery(page)))

	data, err := f.Client.Get(path)
	if err != nil {
		return fmt.Errorf("application list failed: %w", err)
	}

	var resp paginatedResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	if err := p.PrintList(resp.Data, appColumns()); err != nil {
		return err
	}

	pg := resp.Pagination
	start := (pg.Page-1)*pg.PerPage + 1
	end := start + pg.PageItemsCount - 1

	if pg.PageCount > 1 {
		hint := " Use --all to fetch all results."
		if o.all || pg.Page == pg.PageCount {
			hint = ""
		}

		p.PrintMessage("Showing %d-%d of %d (page %d/%d).%s",
			start, end, pg.TotalCount, pg.Page, pg.PageCount, hint)
	} else if pg.TotalCount > 0 {
		p.PrintMessage("Showing %d-%d of %d (page %d/%d).",
			start, end, pg.TotalCount, pg.Page, pg.PageCount)
	}

	return nil
}

func (o *listOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
	var allData []json.RawMessage

	for page := 1; ; page++ {
		path := cmdutil.V1EnvPath(f, fmt.Sprintf("applications/_paged?%s", o.buildQuery(page)))

		data, err := f.Client.Get(path)
		if err != nil {
			return fmt.Errorf("application list failed: %w", err)
		}

		var resp paginatedResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		allData = append(allData, resp.Data...)

		if resp.Pagination.PageCount <= 0 || page >= resp.Pagination.PageCount || page > 1000 {
			break
		}
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(allData)
	}

	if err := p.PrintList(allData, appColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintMessage("Showing %d results.", len(allData))
	}

	return nil
}

func appColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
		{Name: "Type", Value: func(i interface{}) string { return cmdutil.StringField(i, "type") }},
		{Name: "Status", Value: func(i interface{}) string { return cmdutil.StringField(i, "status") }},
		{Name: "Owner", Value: ownerDisplayName},
		{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
		{Name: "Updated", Value: func(i interface{}) string { return cmdutil.StringField(i, "updated_at") }},
	}
}
