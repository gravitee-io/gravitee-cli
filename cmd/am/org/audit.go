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

func newOrgAuditCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "audit",
		Aliases: []string{"audits"},
		Short:   "Manage organization audits",
	}

	cmd.AddCommand(newOrgAuditListCmd(f))
	cmd.AddCommand(newOrgAuditGetCmd(f))

	return cmd
}

// list

type orgAuditListOptions struct {
	factory  *factory.Factory
	typeFlag string
	status   string
	from     string
	to       string
	page     int
	perPage  int
	all      bool
}

func newOrgAuditListCmd(f *factory.Factory) *cobra.Command {
	opts := &orgAuditListOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organization audits",
		Example: `  gio am org audit list
  gio am org audit list --type USER_LOGIN --status SUCCESS
  gio am org audit list --all`,
		Args: cobra.NoArgs,
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

	cmd.Flags().StringVar(&opts.typeFlag, "type", "", "Filter by audit type")
	cmd.Flags().StringVar(&opts.status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&opts.from, "from", "", "Filter from date")
	cmd.Flags().StringVar(&opts.to, "to", "", "Filter to date")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")

	return cmd
}

func (o *orgAuditListOptions) run() error {
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

func (o *orgAuditListOptions) params(page int) am.ListOrgAuditsParams {
	return am.ListOrgAuditsParams{
		Type:    o.typeFlag,
		Status:  o.status,
		From:    o.from,
		To:      o.to,
		Page:    page,
		PerPage: o.perPage,
	}
}

func (o *orgAuditListOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	resp, err := f.AM().ListOrgAudits(o.params(page - 1))
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		raw, _ := json.Marshal(resp)
		return p.PrintDetail(json.RawMessage(raw))
	}

	if err := p.PrintList(resp.Data, orgAuditColumns()); err != nil {
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

func (o *orgAuditListOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
	allData, err := am.FetchAllPages(func(page int) (*am.PaginatedResponse, error) {
		return f.AM().ListOrgAudits(o.params(page))
	}, o.perPage)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(allData)
	}

	if err := p.PrintList(allData, orgAuditColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintHint("Showing %d results.", len(allData))
	}

	return nil
}

func orgAuditColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Type", Value: func(i any) string { return cmdutil.StringField(i, "type") }},
		{Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "status") }},
		{Name: "Actor", Value: func(i any) string {
			m, ok := i.(map[string]any)
			if !ok {
				return ""
			}

			actor, ok := m["actor"].(map[string]any)
			if !ok {
				return ""
			}

			if v, ok := actor["displayName"].(string); ok {
				return v
			}

			return ""
		}},
		{Name: "Target", Value: func(i any) string {
			m, ok := i.(map[string]any)
			if !ok {
				return ""
			}

			target, ok := m["target"].(map[string]any)
			if !ok {
				return ""
			}

			if v, ok := target["displayName"].(string); ok {
				return v
			}

			return ""
		}},
		{Name: "Timestamp", Value: func(i any) string { return cmdutil.StringField(i, "timestamp") }},
	}
}

// get

func newOrgAuditGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <auditID>",
		Short:   "Get organization audit details",
		Example: `  gio am org audit get my-audit-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgAuditGet(f, args[0])
		},
	}
}

func runOrgAuditGet(f *factory.Factory, auditID string) error {
	data, err := f.AM().GetOrgAudit(auditID)
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

	return printOrgAuditDetail(p, data)
}

func printOrgAuditDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"ID", "id"},
		{"Type", "type"},
		{"Status", "status"},
		{"Timestamp", "timestamp"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	if actor, ok := m["actor"].(map[string]any); ok {
		if v, ok := actor["displayName"].(string); ok {
			p.PrintMessage("%-16s%v", "Actor:", v)
		}
	}

	if target, ok := m["target"].(map[string]any); ok {
		if v, ok := target["displayName"].(string); ok {
			p.PrintMessage("%-16s%v", "Target:", v)
		}
	}

	return nil
}
