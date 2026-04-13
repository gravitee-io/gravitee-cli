package org

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newOrgFormCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "form",
		Aliases: []string{"forms"},
		Short:   "Manage organization forms",
	}

	cmd.AddCommand(newOrgFormGetCmd(f))
	cmd.AddCommand(newOrgFormCreateCmd(f))
	cmd.AddCommand(newOrgFormUpdateCmd(f))
	cmd.AddCommand(newOrgFormDeleteCmd(f))

	return cmd
}

func newOrgFormGetCmd(f *factory.Factory) *cobra.Command {
	var template string

	cmd := &cobra.Command{
		Use:     "get --template <TEMPLATE>",
		Short:   "Get organization form by template",
		Example: `  gio am org form get --template LOGIN`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetOrgForm(template)
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

			return printOrgFormDetail(p, data)
		},
	}

	cmd.Flags().StringVar(&template, "template", "", "Form template name (required)")
	_ = cmd.MarkFlagRequired("template")

	return cmd
}

func newOrgFormCreateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "create --file <config.json>",
		Short: "Create an organization form from a JSON file",
		Example: `  gio am org form create --file form.json
  gio am org form create -f form.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONFile(file)
			if err != nil {
				return err
			}

			data, err := f.AM().CreateOrgForm(body)
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

			return printOrgFormDetail(p, data)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newOrgFormUpdateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <formID> --file <config.json>",
		Short: "Update an organization form from a JSON file",
		Example: `  gio am org form update my-form-id --file form.json
  gio am org form update my-form-id -f form.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONFile(file)
			if err != nil {
				return err
			}

			data, err := f.AM().UpdateOrgForm(args[0], body)
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

			return printOrgFormDetail(p, data)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newOrgFormDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <formID>",
		Short:   "Delete an organization form",
		Example: `  gio am org form delete my-form-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteOrgForm(args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Organization form '%s' deleted.", args[0])

			return nil
		},
	}
}

func printOrgFormDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"ID", "id"},
		{"Template", "template"},
		{"Enabled", "enabled"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
