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

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newOrgUserCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Aliases: []string{"users"},
		Short:   "Manage organization users",
	}

	cmd.AddCommand(newOrgUserListCmd(f))
	cmd.AddCommand(newOrgUserGetCmd(f))
	cmd.AddCommand(newOrgUserCreateCmd(f))
	cmd.AddCommand(newOrgUserUpdateCmd(f))
	cmd.AddCommand(newOrgUserDeleteCmd(f))
	cmd.AddCommand(newOrgUserResetPasswordCmd(f))
	cmd.AddCommand(newOrgUserUpdateStatusCmd(f))
	cmd.AddCommand(newOrgUserUpdateUsernameCmd(f))
	cmd.AddCommand(newOrgUserBulkCmd(f))

	return cmd
}

// list

type orgUserListOptions struct {
	factory *factory.Factory
	page    int
	perPage int
	all     bool
}

func newOrgUserListCmd(f *factory.Factory) *cobra.Command {
	opts := &orgUserListOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List organization users",
		Example: `  gio am org user list --per-page 20`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.ValidatePagination(opts.page, opts.perPage); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")

	return cmd
}

func (o *orgUserListOptions) run() error {
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

func (o *orgUserListOptions) params(page int) am.ListOrgUsersParams {
	return am.ListOrgUsersParams{
		Page:    page,
		PerPage: o.perPage,
	}
}

func (o *orgUserListOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	resp, err := f.AM().ListOrgUsers(o.params(page - 1))
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		raw, _ := json.Marshal(resp)
		return p.PrintDetail(json.RawMessage(raw))
	}

	if err := p.PrintList(resp.Data, orgUserColumns()); err != nil {
		return err
	}

	if resp.TotalCount > len(resp.Data) {
		hint := " Use --all to fetch all results."
		if o.all {
			hint = ""
		}

		p.PrintHint("Showing %d of %d.%s", len(resp.Data), resp.TotalCount, hint)
	} else if resp.TotalCount > 0 {
		p.PrintHint("Showing %d results.", len(resp.Data))
	}

	return nil
}

func (o *orgUserListOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
	allData, err := am.FetchAllPages(func(page int) (*am.PaginatedResponse, error) {
		return f.AM().ListOrgUsers(o.params(page))
	}, o.perPage)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(allData)
	}

	if err := p.PrintList(allData, orgUserColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintHint("Showing %d results.", len(allData))
	}

	return nil
}

func orgUserColumns() []printer.Column {
	return []printer.Column{
		{Name: "Username", Value: func(i any) string { return cmdutil.StringField(i, "username") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Email", Value: func(i any) string { return cmdutil.StringField(i, "email") }},
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

// get

func newOrgUserGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <userID>",
		Short:   "Get organization user details",
		Example: `  gio am org user get user-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgUserGet(f, args[0])
		},
	}
}

func runOrgUserGet(f *factory.Factory, userID string) error {
	data, err := f.AM().GetOrgUser(userID)
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

	return printOrgUserDetail(p, data)
}

func printOrgUserDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Username", "username"},
		{"ID", "id"},
		{"Email", "email"},
		{"First Name", "firstName"},
		{"Last Name", "lastName"},
		{"Enabled", "enabled"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}

// create

func newOrgUserCreateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "create --file <user.json>",
		Short: "Create an organization user from a JSON file",
		Example: `  gio am org user create --file user.json
  gio am org user create -f user.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgUserCreate(f, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runOrgUserCreate(f *factory.Factory, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().CreateOrgUser(body)
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

	return printOrgUserDetail(p, data)
}

// update

func newOrgUserUpdateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <userID> --file <user.json>",
		Short: "Update an organization user from a JSON file",
		Example: `  gio am org user update user-id --file user.json
  gio am org user update user-id -f user.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgUserUpdate(f, args[0], file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runOrgUserUpdate(f *factory.Factory, userID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateOrgUser(userID, body)
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

	return printOrgUserDetail(p, data)
}

// delete

func newOrgUserDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <userID>",
		Short:   "Delete an organization user",
		Example: `  gio am org user delete user-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgUserDelete(f, args[0])
		},
	}
}

func runOrgUserDelete(f *factory.Factory, userID string) error {
	if err := f.AM().DeleteOrgUser(userID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Organization user '%s' deleted.", userID)

	return nil
}

// reset-password

func newOrgUserResetPasswordCmd(f *factory.Factory) *cobra.Command {
	var password string

	cmd := &cobra.Command{
		Use:     "reset-password <userID> --password <newPassword>",
		Short:   "Reset an organization user's password",
		Example: `  gio am org user reset-password user-id --password newSecret123`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body := map[string]any{"password": password}
			raw, _ := json.Marshal(body)

			if err := f.AM().ResetOrgUserPassword(args[0], json.RawMessage(raw)); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Password reset for organization user '%s'.", args[0])

			return nil
		},
	}

	cmd.Flags().StringVar(&password, "password", "", "New password (required)")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}

// update-status

func newOrgUserUpdateStatusCmd(f *factory.Factory) *cobra.Command {
	var enabled string

	cmd := &cobra.Command{
		Use:     "update-status <userID> --enabled <true|false>",
		Short:   "Update an organization user's status",
		Example: `  gio am org user update-status user-id --enabled true`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			var enabledVal bool
			switch enabled {
			case "true":
				enabledVal = true
			case "false":
				enabledVal = false
			default:
				return fmt.Errorf("invalid value '%s' for flag --enabled\nHint: allowed values are true, false", enabled)
			}

			body := map[string]any{"enabled": enabledVal}
			raw, _ := json.Marshal(body)

			data, err := f.AM().UpdateOrgUserStatus(args[0], json.RawMessage(raw))
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

			return printOrgUserDetail(p, data)
		},
	}

	cmd.Flags().StringVar(&enabled, "enabled", "", "Enable or disable the user (true/false) (required)")
	_ = cmd.MarkFlagRequired("enabled")

	return cmd
}

// update-username

func newOrgUserUpdateUsernameCmd(f *factory.Factory) *cobra.Command {
	var username string

	cmd := &cobra.Command{
		Use:     "update-username <userID>",
		Short:   "Update an organization user's username",
		Example: `  gio am org user update-username user-id --username newname`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body := map[string]any{"username": username}
			raw, _ := json.Marshal(body)

			if _, err := f.AM().UpdateOrgUsername(args[0], json.RawMessage(raw)); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Username updated for organization user '%s'.", args[0])

			return nil
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "New username (required)")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

// bulk

func newOrgUserBulkCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "bulk --file <operations.json>",
		Short: "Perform bulk operations on organization users",
		Example: `  gio am org user bulk --file operations.json
  gio am org user bulk -f operations.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgUserBulk(f, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON operations file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runOrgUserBulk(f *factory.Factory, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().BulkOrgUserOperation(body)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}
