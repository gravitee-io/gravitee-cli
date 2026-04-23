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
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/apim"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type listOptions struct {
	factory *factory.Factory
	query   string
	status  string
	page    int
	perPage int
	all     bool
	wide    bool
}

func newListCmd(f *factory.Factory) *cobra.Command {
	opts := &listOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List APIs",
		Example: `  gio apim api list
  gio apim api list --status STARTED --per-page 20`,
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
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")
	cmd.Flags().BoolVarP(&opts.wide, "wide", "w", false, "Show additional columns (tags, categories, owner, portal, visibility)")

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

func (o *listOptions) params(page int) apim.ListAPIsParams {
	return apim.ListAPIsParams{
		Query:   o.query,
		Status:  o.status,
		Page:    page,
		PerPage: o.perPage,
	}
}

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	resp, err := f.APIM().ListAPIs(o.params(page))
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		raw, _ := json.Marshal(resp)
		return p.PrintDetail(json.RawMessage(raw))
	}

	if err := p.PrintList(resp.Data, apiColumns(o.wide)); err != nil {
		return err
	}

	pg := resp.Pagination
	cmdutil.PrintPaginationHint(p, pg.Page, pg.PerPage, pg.PageCount, pg.TotalCount, pg.PageItemsCount, o.all)

	return nil
}

func (o *listOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
	allData, err := apim.FetchAllPages(func(page int) (*apim.PaginatedResponse, error) {
		return f.APIM().ListAPIs(o.params(page))
	})
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		return p.PrintDetail(allData)
	}

	if err := p.PrintList(allData, apiColumns(o.wide)); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintHint("Showing %d results.", len(allData))
	}

	return nil
}

func apiColumns(wide bool) []printer.Column {
	cols := []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "API Type", Value: apiTypeLabel},
		{Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "state") }},
		{Name: "Access", Value: apiAccessPath},
	}

	if wide {
		cols = append(cols,
			printer.Column{Name: "Tags", Value: func(i any) string { return joinStringArray(i, "tags") }},
			printer.Column{Name: "Categories", Value: func(i any) string { return joinStringArray(i, "categories") }},
			printer.Column{Name: "Owner", Value: apiOwnerName},
			printer.Column{Name: "Portal", Value: func(i any) string { return cmdutil.StringField(i, "lifecycleState") }},
			printer.Column{Name: "Visibility", Value: func(i any) string { return cmdutil.StringField(i, "visibility") }},
		)
	}

	cols = append(cols, printer.Column{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }})

	return cols
}

func apiTypeLabel(item any) string {
	defVersion := cmdutil.StringField(item, "definitionVersion")

	switch defVersion {
	case "V4":
		apiType := cmdutil.StringField(item, "type")
		if apiType != "" {
			return "V4 " + apiType
		}

		return "V4"
	case "V2":
		return "V2 HTTP Proxy"
	case "V1":
		return "V1 HTTP Proxy"
	}

	return defVersion
}

func apiAccessPath(item any) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	// V1/V2 APIs have a direct contextPath field
	if cp, cpOK := m["contextPath"].(string); cpOK && cp != "" {
		return cp
	}

	// V4 APIs: extract from listeners[].paths[].path
	listeners, ok := m["listeners"].([]any)
	if !ok {
		return ""
	}

	for _, l := range listeners {
		lm, ok := l.(map[string]any)
		if !ok {
			continue
		}

		paths, ok := lm["paths"].([]any)
		if !ok {
			continue
		}

		for _, p := range paths {
			pm, ok := p.(map[string]any)
			if !ok {
				continue
			}

			if path, ok := pm["path"].(string); ok {
				return path
			}
		}
	}

	return ""
}

func apiOwnerName(item any) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	owner, ok := m["primaryOwner"].(map[string]any)
	if !ok {
		return ""
	}

	s, _ := owner["displayName"].(string)

	return s
}

func joinStringArray(item any, key string) string {
	m, ok := item.(map[string]any)
	if !ok {
		return ""
	}

	arr, ok := m[key].([]any)
	if !ok || len(arr) == 0 {
		return ""
	}

	parts := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			parts = append(parts, s)
		}
	}

	return strings.Join(parts, ", ")
}
