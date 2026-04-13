package org

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newOrgRoleCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "role",
		Aliases: []string{"roles"},
		Short:   "Manage organization roles",
	}

	cmd.AddCommand(newOrgRoleListCmd(f))
	cmd.AddCommand(newOrgRoleGetCmd(f))
	cmd.AddCommand(newOrgRoleCreateCmd(f))
	cmd.AddCommand(newOrgRoleUpdateCmd(f))
	cmd.AddCommand(newOrgRoleDeleteCmd(f))

	return cmd
}

// list

func newOrgRoleListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List organization roles",
		Example: `  gio am org role list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgRoleList(f)
		},
	}
}

func runOrgRoleList(f *factory.Factory) error {
	items, err := f.AM().ListOrgRoles()
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

	return p.PrintList(items, orgRoleColumns())
}

func orgRoleColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Description", Value: func(i any) string { return cmdutil.StringField(i, "description") }},
		{Name: "System", Value: func(i any) string {
			m, ok := i.(map[string]any)
			if !ok {
				return ""
			}

			if v, ok := m["system"].(bool); ok && v {
				return "true"
			}

			return "false"
		}},
	}
}

// get

func newOrgRoleGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <roleID>",
		Short:   "Get organization role details",
		Example: `  gio am org role get role-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgRoleGet(f, args[0])
		},
	}
}

func runOrgRoleGet(f *factory.Factory, roleID string) error {
	data, err := f.AM().GetOrgRole(roleID)
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

	return printOrgRoleDetail(p, data)
}

func printOrgRoleDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"Description", "description"},
		{"System", "system"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}

// create

func newOrgRoleCreateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "create --file <role.json>",
		Short: "Create an organization role from a JSON file",
		Example: `  gio am org role create --file role.json
  gio am org role create -f role.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgRoleCreate(f, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runOrgRoleCreate(f *factory.Factory, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().CreateOrgRole(body)
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

	return printOrgRoleDetail(p, data)
}

// update

func newOrgRoleUpdateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <roleID> --file <role.json>",
		Short: "Update an organization role from a JSON file",
		Example: `  gio am org role update role-id --file role.json
  gio am org role update role-id -f role.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgRoleUpdate(f, args[0], file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runOrgRoleUpdate(f *factory.Factory, roleID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateOrgRole(roleID, body)
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

	return printOrgRoleDetail(p, data)
}

// delete

func newOrgRoleDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <roleID>",
		Short:   "Delete an organization role",
		Example: `  gio am org role delete role-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgRoleDelete(f, args[0])
		},
	}
}

func runOrgRoleDelete(f *factory.Factory, roleID string) error {
	if err := f.AM().DeleteOrgRole(roleID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Organization role '%s' deleted.", roleID)

	return nil
}
