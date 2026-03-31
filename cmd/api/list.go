package api

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
	sort    string
	order   string
	page    int
	perPage int
	all     bool
}

func newListCmd(f *factory.Factory) *cobra.Command {
	opts := &listOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List APIs",
		Example: `  gio api list
  gio api list --status STARTED --per-page 20`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.query, "query", "", "Search by name or description")
	cmd.Flags().StringVar(&opts.status, "status", "", "Filter by status: STARTED, STOPPED")
	cmd.Flags().StringVar(&opts.sort, "sort", "", "Sort field: name, updatedAt, createdAt")
	cmd.Flags().StringVar(&opts.order, "order", "asc", "Sort order: asc, desc")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")

	return cmd
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
	q.Set("perPage", strconv.Itoa(o.perPage))

	if o.query != "" {
		q.Set("q", o.query)
	}

	if o.status != "" {
		q.Set("status", o.status)
	}

	if o.sort != "" {
		q.Set("sortBy", o.sort)
	}

	if o.order != "" {
		q.Set("order", o.order)
	}

	return q.Encode()
}

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	path := cmdutil.V2EnvPath(f, "apis?"+o.buildQuery(page))

	data, err := f.Client.Get(path)
	if err != nil {
		return err
	}

	var resp paginatedResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	if err := p.PrintList(resp.Data, apiColumns()); err != nil {
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
		path := cmdutil.V2EnvPath(f, "apis?"+o.buildQuery(page))

		data, err := f.Client.Get(path)
		if err != nil {
			return err
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

	if err := p.PrintList(allData, apiColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintMessage("Showing %d results.", len(allData))
	}

	return nil
}

func apiColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
		{Name: "Status", Value: func(i interface{}) string { return cmdutil.StringField(i, "state") }},
		{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
		{Name: "Definition", Value: func(i interface{}) string { return cmdutil.StringField(i, "definitionVersion") }},
		{Name: "Updated", Value: func(i interface{}) string { return cmdutil.StringField(i, "updatedAt") }},
	}
}
