// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/apim"
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
		Example: `  gio apim api logs 8a7b3c4d-... --from 1711497600000 --to 1711584000000
  gio apim api logs 8a7b3c4d-... --methods GET --per-page 20`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			apiID, err := f.APIM().ResolveAPI(args[0])
			if err != nil {
				return err
			}

			opts.apiID = apiID

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

func (o *logsOptions) params(page int) apim.ListAPILogsParams {
	return apim.ListAPILogsParams{
		ApplicationIDs: o.applicationIDs,
		PlanIDs:        o.planIDs,
		Methods:        o.methods,
		From:           o.from,
		To:             o.to,
		Page:           page,
		PerPage:        o.perPage,
	}
}

func (o *logsOptions) run() error {
	f := o.factory
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if o.all {
		return o.fetchAllLogs(f, p)
	}

	return o.fetchLogsPage(f, p, o.page)
}

func (o *logsOptions) fetchLogsPage(f *factory.Factory, p *printer.Printer, page int) error {
	resp, err := f.APIM().ListAPILogs(o.apiID, o.params(page))
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		return p.PrintDetail(resp)
	}

	if err := p.PrintList(resp.Data, logColumns()); err != nil {
		return err
	}

	pg := resp.Pagination
	cmdutil.PrintPaginationHint(p, pg.Page, pg.PerPage, pg.PageCount, pg.TotalCount, pg.PageItemsCount, o.all)

	return nil
}

func (o *logsOptions) fetchAllLogs(f *factory.Factory, p *printer.Printer) error {
	allData, err := apim.FetchAllPages(func(page int) (*apim.PaginatedResponse, error) {
		return f.APIM().ListAPILogs(o.apiID, o.params(page))
	})
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		return p.PrintDetail(allData)
	}

	if err := p.PrintList(allData, logColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintHint("Showing %d results.", len(allData))
	}

	return nil
}

func logColumns() []printer.Column {
	return []printer.Column{
		{Name: "Timestamp", Value: func(i any) string { return cmdutil.StringField(i, "timestamp") }},
		{Name: "Request ID", Value: func(i any) string { return cmdutil.StringField(i, "requestId") }},
		{Name: "Method", Value: func(i any) string { return cmdutil.StringField(i, "method") }},
		{Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "status") }},
		{Name: "Path", Value: func(i any) string { return cmdutil.StringField(i, "path") }},
	}
}
