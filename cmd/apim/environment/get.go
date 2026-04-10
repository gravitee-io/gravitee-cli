package environment

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <envId>",
		Short:   "Get environment details",
		Example: `  gio apim environment get prod-1111-2222-3333-444455556666`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.RequireNonEmpty("envId", args[0]); err != nil {
				return err
			}

			return runGet(f, args[0])
		},
	}
}

func runGet(f *factory.Factory, envID string) error {
	data, err := f.APIM().GetEnvironment(envID)
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

	return printEnvDetail(p, data)
}

func printEnvDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fields := []struct {
		label string
		key   string
	}{
		{"Name", "name"},
		{"ID", "id"},
		{"Description", "description"},
	}

	for _, field := range fields {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
