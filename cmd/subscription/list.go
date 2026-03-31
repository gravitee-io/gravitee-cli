package subscription

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type listOptions struct {
	factory *factory.Factory
	apiID   string
	planID  string
	appID   string
	status  []string
	page    int
	perPage int
	all     bool
}

func newListCmd(f *factory.Factory) *cobra.Command {
	opts := &listOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "list --api <apiId>",
		Short: "List subscriptions for an API",
		Example: `  gio subscription list --api 8a7b3c4d-1234-5678-abcd-ef0123456789
  gio sub list --api 8a7b3c4d --status PENDING`,
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

	cmd.Flags().StringVar(&opts.apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")
	cmd.Flags().StringArrayVarP(&opts.status, "status", "s", []string{"ACCEPTED", "PENDING", "PAUSED"}, "Filter by status")
	cmd.Flags().StringVar(&opts.planID, "plan", "", "Filter by plan ID")
	cmd.Flags().StringVar(&opts.appID, "app", "", "Filter by application ID")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")

	return cmd
}

var validSubStatuses = []string{"PENDING", "ACCEPTED", "REJECTED", "PAUSED", "CLOSED"}

func (o *listOptions) validate() error {
	for _, s := range o.status {
		if err := cmdutil.ValidateEnum(s, "status", validSubStatuses); err != nil {
			return err
		}
	}

	return nil
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

	if len(o.status) > 0 {
		q.Set("statuses", strings.Join(o.status, ","))
	}

	if o.planID != "" {
		q.Set("planIds", o.planID)
	}

	if o.appID != "" {
		q.Set("applicationIds", o.appID)
	}

	return q.Encode()
}

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/subscriptions?%s", o.apiID, o.buildQuery(page)))

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

	if err := p.PrintList(resp.Data, subColumns()); err != nil {
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
		path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/subscriptions?%s", o.apiID, o.buildQuery(page)))

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

	if err := p.PrintList(allData, subColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintMessage("Showing %d results.", len(allData))
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

func nestedID(item interface{}, key string) string {
	m, ok := item.(map[string]interface{})
	if !ok {
		return ""
	}

	nested, ok := m[key].(map[string]interface{})
	if !ok {
		return cmdutil.StringField(item, key+"Id")
	}

	s, _ := nested["id"].(string)

	return s
}

func subColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
		{Name: "Plan", Value: func(i interface{}) string { return nestedID(i, "plan") }},
		{Name: "Application", Value: func(i interface{}) string { return nestedID(i, "application") }},
		{Name: "Status", Value: func(i interface{}) string { return cmdutil.StringField(i, "status") }},
		{Name: "Created", Value: func(i interface{}) string { return cmdutil.StringField(i, "createdAt") }},
	}
}
