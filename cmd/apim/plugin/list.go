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

package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

var validTypes = []string{"endpoints", "entrypoints", "policies"}

// typeLabels maps the plural API path segment to the singular label for the TYPE column.
var typeLabels = map[string]string{
	"endpoints":   "endpoint",
	"entrypoints": "entrypoint",
	"policies":    "policy",
}

type listOptions struct {
	factory    *factory.Factory
	pluginType string
}

func newListCmd(f *factory.Factory) *cobra.Command {
	opts := &listOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List plugins",
		Example: `  gctl apim plugin list
  gctl apim plugin list --type policies
  gctl apim plugin list --type endpoints -o json`,
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

	cmd.Flags().StringVarP(&opts.pluginType, "type", "t", "", "Filter by plugin type: endpoints, entrypoints, policies")

	return cmd
}

func (o *listOptions) validate() error {
	if o.pluginType == "" {
		return nil
	}

	for _, t := range validTypes {
		if o.pluginType == t {
			return nil
		}
	}

	return fmt.Errorf("invalid value '%s' for flag --type\nHint: allowed values are endpoints, entrypoints, policies", o.pluginType)
}

func (o *listOptions) run() error {
	f := o.factory
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if o.pluginType != "" {
		return o.fetchSingleType(f, p)
	}

	return o.fetchAllTypes(f, p)
}

func (o *listOptions) fetchSingleType(f *factory.Factory, p *printer.Printer) error {
	data, err := f.APIM().ListPlugins(o.pluginType)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(data)
	}

	return p.PrintList(data, pluginColumns())
}

func (o *listOptions) fetchAllTypes(f *factory.Factory, p *printer.Printer) error {
	var allItems []any

	for _, t := range validTypes {
		data, err := f.APIM().ListPlugins(t)
		if err != nil {
			return err
		}

		var items []any
		if err := json.Unmarshal(data, &items); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		label := typeLabels[t]
		for _, item := range items {
			if m, ok := item.(map[string]any); ok {
				m["type"] = label
			}

			allItems = append(allItems, item)
		}
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(allItems)
	}

	return p.PrintList(allItems, pluginColumnsWithType())
}

func pluginColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "Version", Value: func(i any) string { return cmdutil.StringField(i, "version") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Description", Value: func(i any) string { return cmdutil.StringField(i, "description") }},
	}
}

func pluginColumnsWithType() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "Type", Value: func(i any) string { return cmdutil.StringField(i, "type") }},
		{Name: "Version", Value: func(i any) string { return cmdutil.StringField(i, "version") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Description", Value: func(i any) string { return cmdutil.StringField(i, "description") }},
	}
}
