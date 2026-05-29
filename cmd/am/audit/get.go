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

package audit

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get <auditID>",
		Short:   "Get audit details",
		Example: `  gctl am audit get my-audit-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runGet(f, *domainID, args[0])
		},
	}
}

func runGet(f *factory.Factory, domainID, auditID string) error {
	data, err := f.AM().GetAudit(domainID, auditID)
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

	return printAuditDetail(p, data)
}

func printAuditDetail(p *printer.Printer, data []byte) error {
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
