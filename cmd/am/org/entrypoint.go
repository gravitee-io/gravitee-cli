package org

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newOrgEntrypointCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "entrypoint",
		Aliases: []string{"entrypoints"},
		Short:   "Manage organization entrypoints",
	}

	cmd.AddCommand(newOrgEntrypointListCmd(f))
	cmd.AddCommand(newOrgEntrypointGetCmd(f))
	cmd.AddCommand(newOrgEntrypointCreateCmd(f))
	cmd.AddCommand(newOrgEntrypointUpdateCmd(f))
	cmd.AddCommand(newOrgEntrypointDeleteCmd(f))

	return cmd
}

func newOrgEntrypointListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List organization entrypoints",
		Example: `  gio am org entrypoint list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			items, err := f.AM().ListOrgEntrypoints()
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

			return p.PrintList(items, orgEntrypointColumns())
		},
	}
}

func newOrgEntrypointGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <entrypointID>",
		Short:   "Get organization entrypoint details",
		Example: `  gio am org entrypoint get my-entrypoint-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetOrgEntrypoint(args[0])
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

			return printOrgEntrypointDetail(p, data)
		},
	}
}

func newOrgEntrypointCreateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "create --file <config.json>",
		Short: "Create an organization entrypoint from a JSON file",
		Example: `  gio am org entrypoint create --file entrypoint.json
  gio am org entrypoint create -f entrypoint.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONFile(file)
			if err != nil {
				return err
			}

			data, err := f.AM().CreateOrgEntrypoint(body)
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

			return printOrgEntrypointDetail(p, data)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newOrgEntrypointUpdateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <entrypointID> --file <config.json>",
		Short: "Update an organization entrypoint from a JSON file",
		Example: `  gio am org entrypoint update my-entrypoint-id --file entrypoint.json
  gio am org entrypoint update my-entrypoint-id -f entrypoint.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONFile(file)
			if err != nil {
				return err
			}

			data, err := f.AM().UpdateOrgEntrypoint(args[0], body)
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

			return printOrgEntrypointDetail(p, data)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newOrgEntrypointDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <entrypointID>",
		Short:   "Delete an organization entrypoint",
		Example: `  gio am org entrypoint delete my-entrypoint-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteOrgEntrypoint(args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Organization entrypoint '%s' deleted.", args[0])

			return nil
		},
	}
}

func orgEntrypointColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "URL", Value: func(i any) string { return cmdutil.StringField(i, "url") }},
	}
}

func printOrgEntrypointDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"URL", "url"},
		{"Default", "defaultEntrypoint"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
