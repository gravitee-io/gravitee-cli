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

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newOrgIDPCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "idp",
		Aliases: []string{"identity-provider", "identity-providers"},
		Short:   "Manage organization identity providers",
	}

	cmd.AddCommand(newOrgIDPListCmd(f))
	cmd.AddCommand(newOrgIDPGetCmd(f))
	cmd.AddCommand(newOrgIDPCreateCmd(f))
	cmd.AddCommand(newOrgIDPUpdateCmd(f))
	cmd.AddCommand(newOrgIDPDeleteCmd(f))

	return cmd
}

func newOrgIDPListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List organization identity providers",
		Example: `  gctl am org idp list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			items, err := f.AM().ListOrgIdentityProviders()
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

			return p.PrintList(items, orgIDPColumns())
		},
	}
}

func newOrgIDPGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <idpID>",
		Short:   "Get organization identity provider details",
		Example: `  gctl am org idp get my-idp-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetOrgIdentityProvider(args[0])
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

			return printOrgIDPDetail(p, data)
		},
	}
}

func newOrgIDPCreateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "create [-f <file>]",
		Short: "Create an organization identity provider from a JSON file or stdin",
		Example: `  gctl am org idp create --file idp.json
  gctl am org idp create -f idp.json
  envsubst < idp.json | gctl am org idp create`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONInput(file, f.IOStreams.In)
			if err != nil {
				return err
			}

			data, err := f.AM().CreateOrgIdentityProvider(body)
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

			return printOrgIDPDetail(p, data)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a JSON file (optional - reads from stdin if omitted)")

	return cmd
}

func newOrgIDPUpdateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <idpID> [-f <file>]",
		Short: "Update an organization identity provider from a JSON file or stdin",
		Example: `  gctl am org idp update my-idp-id --file idp.json
  gctl am org idp update my-idp-id -f idp.json
  envsubst < idp.json | gctl am org idp update my-idp-id`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONInput(file, f.IOStreams.In)
			if err != nil {
				return err
			}

			data, err := f.AM().UpdateOrgIdentityProvider(args[0], body)
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

			return printOrgIDPDetail(p, data)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a JSON file (optional - reads from stdin if omitted)")

	return cmd
}

func newOrgIDPDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <idpID>",
		Short:   "Delete an organization identity provider",
		Example: `  gctl am org idp delete my-idp-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteOrgIdentityProvider(args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Organization identity provider '%s' deleted.", args[0])

			return nil
		},
	}
}

func orgIDPColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Type", Value: func(i any) string { return cmdutil.StringField(i, "type") }},
	}
}

func printOrgIDPDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"Type", "type"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
