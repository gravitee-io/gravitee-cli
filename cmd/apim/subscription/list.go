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

package subscription

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/apim"
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
		Example: `  gio apim subscription list --api /my/api
  gio apim sub list --api 8a7b3c4d --status PENDING`,
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

	cmdutil.AddAPIFlag(cmd, &opts.apiID)
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
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if o.all {
		return o.fetchAll(f, p)
	}

	return o.fetchPage(f, p, o.page)
}

func (o *listOptions) params(page int) apim.ListSubscriptionsParams {
	return apim.ListSubscriptionsParams{
		Statuses: o.status,
		PlanID:   o.planID,
		AppID:    o.appID,
		Page:     page,
		PerPage:  o.perPage,
	}
}

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	resp, err := f.APIM().ListSubscriptions(o.apiID, o.params(page))
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		return p.PrintDetail(resp)
	}

	if err := p.PrintList(resp.Data, subColumns()); err != nil {
		return err
	}

	pg := resp.Pagination
	cmdutil.PrintPaginationHint(p, pg.Page, pg.PerPage, pg.PageCount, pg.TotalCount, pg.PageItemsCount, o.all)

	return nil
}

func (o *listOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
	allData, err := apim.FetchAllPages(func(page int) (*apim.PaginatedResponse, error) {
		return f.APIM().ListSubscriptions(o.apiID, o.params(page))
	})
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		return p.PrintDetail(allData)
	}

	if err := p.PrintList(allData, subColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintHint("Showing %d results.", len(allData))
	}

	return nil
}

func nestedID(item any, key string) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	nested, ok := m[key].(map[string]any)
	if !ok {
		return cmdutil.StringField(item, key+"Id")
	}

	s, _ := nested["id"].(string)

	return s
}

func subColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Plan", Value: func(i any) string { return nestedID(i, "plan") }},
		{Name: "Application", Value: func(i any) string { return nestedID(i, "application") }},
		{Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "status") }},
		{Name: "Created", Value: func(i any) string { return cmdutil.StringField(i, "createdAt") }},
	}
}
