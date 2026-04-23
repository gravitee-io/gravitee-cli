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

package metadata

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
	page    int
	perPage int
	all     bool
}

func newListCmd(f *factory.Factory) *cobra.Command {
	opts := &listOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "list --api <apiId>",
		Short: "List metadata for an API",
		Example: `  gio apim metadata list --api /my/api
  gio apim metadata list --api 8a7b3c4d-1234-5678-abcd-ef0123456789 --all`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmdutil.AddAPIFlag(cmd, &opts.apiID)
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")

	return cmd
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

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	resp, err := f.APIM().ListMetadata(o.apiID, page, o.perPage)
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		return p.PrintDetail(resp)
	}

	if err := p.PrintList(resp.Data, metadataColumns()); err != nil {
		return err
	}

	pg := resp.Pagination
	cmdutil.PrintPaginationHint(p, pg.Page, pg.PerPage, pg.PageCount, pg.TotalCount, pg.PageItemsCount, o.all)

	return nil
}

func (o *listOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
	allData, err := apim.FetchAllPages(func(page int) (*apim.PaginatedResponse, error) {
		return f.APIM().ListMetadata(o.apiID, page, o.perPage)
	})
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		return p.PrintDetail(allData)
	}

	if err := p.PrintList(allData, metadataColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintHint("Showing %d results.", len(allData))
	}

	return nil
}

func metadataColumns() []printer.Column {
	return []printer.Column{
		{Name: "Key", Value: func(i any) string { return cmdutil.StringField(i, "key") }},
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "Value", Value: func(i any) string { return cmdutil.StringField(i, "value") }},
		{Name: "Format", Value: func(i any) string { return cmdutil.StringField(i, "format") }},
		{Name: "Updated", Value: func(i any) string { return cmdutil.StringField(i, "updatedAt") }},
	}
}
