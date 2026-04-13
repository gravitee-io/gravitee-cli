package org

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newOrgTagCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tag",
		Aliases: []string{"tags"},
		Short:   "Manage organization sharding tags",
	}

	cmd.AddCommand(newOrgTagListCmd(f))
	cmd.AddCommand(newOrgTagGetCmd(f))
	cmd.AddCommand(newOrgTagCreateCmd(f))
	cmd.AddCommand(newOrgTagUpdateCmd(f))
	cmd.AddCommand(newOrgTagDeleteCmd(f))

	return cmd
}

func newOrgTagListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List organization sharding tags",
		Example: `  gio am org tag list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			items, err := f.AM().ListOrgTags()
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

			return p.PrintList(items, orgTagColumns())
		},
	}
}

func newOrgTagGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <tagID>",
		Short:   "Get organization sharding tag details",
		Example: `  gio am org tag get my-tag-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetOrgTag(args[0])
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

			return printOrgTagDetail(p, data)
		},
	}
}

func newOrgTagCreateCmd(f *factory.Factory) *cobra.Command {
	var (
		name        string
		description string
	)

	cmd := &cobra.Command{
		Use:   "create --name <tag-name> [--description <description>]",
		Short: "Create an organization sharding tag",
		Example: `  gio am org tag create --name "my-tag"
  gio am org tag create --name "my-tag" --description "Tag description"`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			payload := map[string]string{"name": name}
			if description != "" {
				payload["description"] = description
			}

			body, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("failed to build request body: %w", err)
			}

			data, err := f.AM().CreateOrgTag(json.RawMessage(body))
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

			return printOrgTagDetail(p, data)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Tag name (required)")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&description, "description", "", "Tag description")

	return cmd
}

func newOrgTagUpdateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <tagID> --file <config.json>",
		Short: "Update an organization sharding tag from a JSON file",
		Example: `  gio am org tag update my-tag-id --file tag.json
  gio am org tag update my-tag-id -f tag.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONFile(file)
			if err != nil {
				return err
			}

			data, err := f.AM().UpdateOrgTag(args[0], body)
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

			return printOrgTagDetail(p, data)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newOrgTagDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <tagID>",
		Short:   "Delete an organization sharding tag",
		Example: `  gio am org tag delete my-tag-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteOrgTag(args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Organization tag '%s' deleted.", args[0])

			return nil
		},
	}
}

func orgTagColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Description", Value: func(i any) string { return cmdutil.StringField(i, "description") }},
	}
}

func printOrgTagDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"Description", "description"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
