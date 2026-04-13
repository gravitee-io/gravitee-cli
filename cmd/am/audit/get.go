package audit

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get <auditID>",
		Short:   "Get audit details",
		Example: `  gio am audit get my-audit-id --domain my-domain`,
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
