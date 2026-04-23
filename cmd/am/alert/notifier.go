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

package alert

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newNotifierCmd(f *factory.Factory, domainID *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "notifier",
		Aliases: []string{"notifiers"},
		Short:   "Manage alert notifiers",
	}

	cmd.AddCommand(newNotifierListCmd(f, domainID))
	cmd.AddCommand(newNotifierGetCmd(f, domainID))
	cmd.AddCommand(newNotifierCreateCmd(f, domainID))
	cmd.AddCommand(newNotifierUpdateCmd(f, domainID))
	cmd.AddCommand(newNotifierDeleteCmd(f, domainID))

	return cmd
}

func newNotifierListCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List alert notifiers",
		Example: `  gio am alert notifier list --domain my-domain`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runNotifierList(f, *domainID)
		},
	}
}

func runNotifierList(f *factory.Factory, domainID string) error {
	items, err := f.AM().ListAlertNotifiers(domainID)
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

	return p.PrintList(items, notifierColumns())
}

func notifierColumns() []printer.Column {
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

func newNotifierGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get <notifierID>",
		Short:   "Get alert notifier details",
		Example: `  gio am alert notifier get my-notifier-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runNotifierGet(f, *domainID, args[0])
		},
	}
}

func runNotifierGet(f *factory.Factory, domainID, notifierID string) error {
	data, err := f.AM().GetAlertNotifier(domainID, notifierID)
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

	return printNotifierDetail(p, data)
}

func printNotifierDetail(p *printer.Printer, data []byte) error {
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

func newNotifierCreateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "create --file <config.json>",
		Short: "Create an alert notifier from a JSON file",
		Example: `  gio am alert notifier create --domain my-domain --file notifier.json
  gio am alert notifier create --domain my-domain -f notifier.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runNotifierCreate(f, *domainID, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runNotifierCreate(f *factory.Factory, domainID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().CreateAlertNotifier(domainID, body)
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

	return printNotifierDetail(p, data)
}

func newNotifierUpdateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <notifierID> --file <config.json>",
		Short: "Update an alert notifier from a JSON file",
		Example: `  gio am alert notifier update my-notifier-id --domain my-domain --file notifier.json
  gio am alert notifier update my-notifier-id --domain my-domain -f notifier.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runNotifierUpdate(f, *domainID, args[0], file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runNotifierUpdate(f *factory.Factory, domainID, notifierID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateAlertNotifier(domainID, notifierID, body)
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

	return printNotifierDetail(p, data)
}

func newNotifierDeleteCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <notifierID>",
		Short:   "Delete an alert notifier",
		Example: `  gio am alert notifier delete my-notifier-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runNotifierDelete(f, *domainID, args[0])
		},
	}
}

func runNotifierDelete(f *factory.Factory, domainID, notifierID string) error {
	if err := f.AM().DeleteAlertNotifier(domainID, notifierID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Alert notifier '%s' deleted.", notifierID)

	return nil
}
