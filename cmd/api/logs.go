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

type logsOptions struct {
	factory        *factory.Factory
	apiID          string
	applicationIDs []string
	planIDs        []string
	methods        []string
	from           int64
	to             int64
	page           int
	perPage        int
	all            bool
}

func newLogsCmd(f *factory.Factory) *cobra.Command {
	opts := &logsOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "logs <apiId>",
		Short: "List API connection logs",
		Example: `  gio api logs 8a7b3c4d-... --from 1711497600000 --to 1711584000000
  gio api logs 8a7b3c4d-... --methods GET --per-page 20`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.apiID = args[0]

			return opts.run()
		},
	}

	cmd.Flags().Int64Var(&opts.from, "from", 0, "Start timestamp (epoch millis)")
	cmd.Flags().Int64Var(&opts.to, "to", 0, "End timestamp (epoch millis)")
	cmd.Flags().StringArrayVar(&opts.applicationIDs, "application-ids", nil, "Filter by application IDs")
	cmd.Flags().StringArrayVar(&opts.planIDs, "plan-ids", nil, "Filter by plan IDs")
	cmd.Flags().StringArrayVar(&opts.methods, "methods", nil, "Filter by HTTP methods")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")

	return cmd
}

func (o *logsOptions) run() error {
	f := o.factory
	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/logs?%s", o.apiID, o.buildQuery(o.page)))

		data, err := f.Client.Get(path)
		if err != nil {
			return err
		}

		return p.PrintDetail(json.RawMessage(data))
	}

	if o.all {
		return o.fetchAllLogs(f, p)
	}

	return o.fetchLogsPage(f, p, o.page)
}

func (o *logsOptions) buildQuery(page int) string {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("perPage", strconv.Itoa(o.perPage))

	if o.from != 0 {
		q.Set("from", strconv.FormatInt(o.from, 10))
	}

	if o.to != 0 {
		q.Set("to", strconv.FormatInt(o.to, 10))
	}

	for _, id := range o.applicationIDs {
		q.Add("applicationIds", id)
	}

	for _, id := range o.planIDs {
		q.Add("planIds", id)
	}

	for _, m := range o.methods {
		q.Add("methods", m)
	}

	return q.Encode()
}

func (o *logsOptions) fetchLogsPage(f *factory.Factory, p *printer.Printer, page int) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/logs?%s", o.apiID, o.buildQuery(page)))

	data, err := f.Client.Get(path)
	if err != nil {
		return err
	}

	var resp paginatedResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if err := p.PrintList(resp.Data, logColumns()); err != nil {
		return err
	}

	pg := resp.Pagination
	if pg.TotalCount > 0 {
		start := (pg.Page-1)*pg.PerPage + 1
		end := start + pg.PageItemsCount - 1

		if pg.PageCount > 1 && !o.all {
			p.PrintMessage("Showing %d-%d of %d (page %d/%d). Use --all to fetch all results.",
				start, end, pg.TotalCount, pg.Page, pg.PageCount)
		} else {
			p.PrintMessage("Showing %d-%d of %d (page %d/%d).",
				start, end, pg.TotalCount, pg.Page, pg.PageCount)
		}
	}

	return nil
}

func (o *logsOptions) fetchAllLogs(f *factory.Factory, p *printer.Printer) error {
	var allData []json.RawMessage

	for page := 1; ; page++ {
		path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/logs?%s", o.apiID, o.buildQuery(page)))

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

	if err := p.PrintList(allData, logColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintMessage("Showing %d results.", len(allData))
	}

	return nil
}

func logColumns() []printer.Column {
	return []printer.Column{
		{Name: "Timestamp", Value: func(i interface{}) string { return cmdutil.StringField(i, "timestamp") }},
		{Name: "Request ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "requestId") }},
		{Name: "Method", Value: func(i interface{}) string { return cmdutil.StringField(i, "method") }},
		{Name: "Status", Value: func(i interface{}) string { return cmdutil.StringField(i, "status") }},
		{Name: "Path", Value: func(i interface{}) string { return cmdutil.StringField(i, "path") }},
	}
}
