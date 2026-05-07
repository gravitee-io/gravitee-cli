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

package org

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newOrgReporterCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "reporter",
		Aliases: []string{"reporters"},
		Short:   "Manage organization reporters",
	}

	cmd.AddCommand(newOrgReporterListCmd(f))
	cmd.AddCommand(newOrgReporterGetCmd(f))
	cmd.AddCommand(newOrgReporterCreateCmd(f))
	cmd.AddCommand(newOrgReporterUpdateCmd(f))
	cmd.AddCommand(newOrgReporterDeleteCmd(f))

	return cmd
}

func newOrgReporterListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List organization reporters",
		Example: `  gio am org reporter list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			items, err := f.AM().ListOrgReporters()
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(items)
			}

			return p.PrintList(items, orgReporterColumns())
		},
	}
}

func newOrgReporterGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <reporterID>",
		Short:   "Get organization reporter details",
		Example: `  gio am org reporter get my-reporter-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetOrgReporter(args[0])
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(data)
			}

			return printOrgReporterDetail(p, data)
		},
	}
}

func newOrgReporterCreateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "create [-f <file>]",
		Short: "Create an organization reporter from a JSON file or stdin",
		Example: `  gio am org reporter create --file reporter.json
  gio am org reporter create -f reporter.json
  envsubst < reporter.json | gio am org reporter create`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONInput(file, f.IOStreams.In)
			if err != nil {
				return err
			}

			data, err := f.AM().CreateOrgReporter(body)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(data)
			}

			return printOrgReporterDetail(p, data)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a JSON file (optional - reads from stdin if omitted)")

	return cmd
}

func newOrgReporterUpdateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <reporterID> [-f <file>]",
		Short: "Update an organization reporter from a JSON file or stdin",
		Example: `  gio am org reporter update my-reporter-id --file reporter.json
  gio am org reporter update my-reporter-id -f reporter.json
  envsubst < reporter.json | gio am org reporter update my-reporter-id`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONInput(file, f.IOStreams.In)
			if err != nil {
				return err
			}

			data, err := f.AM().UpdateOrgReporter(args[0], body)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(data)
			}

			return printOrgReporterDetail(p, data)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a JSON file (optional - reads from stdin if omitted)")

	return cmd
}

func newOrgReporterDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <reporterID>",
		Short:   "Delete an organization reporter",
		Example: `  gio am org reporter delete my-reporter-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteOrgReporter(args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Organization reporter '%s' deleted.", args[0])

			return nil
		},
	}
}

func orgReporterColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Type", Value: func(i any) string { return cmdutil.StringField(i, "type") }},
		{Name: "Enabled", Value: func(i any) string {
			m, ok := i.(map[string]any)
			if !ok {
				return ""
			}

			if v, ok := m["enabled"].(bool); ok && v {
				return "true"
			}

			return "false"
		}},
	}
}

func printOrgReporterDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"Type", "type"},
		{"Enabled", "enabled"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
