package api

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newLogCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "log <apiId> <requestId>",
		Short:   "Get details of a specific request log",
		Example: `  gio apim api log 8a7b3c4d-... req-aaaa-bbbb-cccc-dddd-eeeeeeee`,
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runLog(f, args[0], args[1])
		},
	}
}

func runLog(f *factory.Factory, apiID, requestID string) error {
	data, err := f.APIM().GetAPILog(apiID, requestID)
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

	return printLogDetail(p, data)
}

func printLogDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Request ID", "requestId"},
		{"Timestamp", "timestamp"},
		{"Method", "method"},
		{"Path", "path"},
		{"Status", "status"},
		{"Response Time", "responseTime"},
		{"Application", "application"},
		{"Plan", "plan"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-17s%v", field.label+":", v)
		}
	}

	return nil
}
