package environment

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List environments",
		Example: `  gio apim environment list
  gio env list`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runList(f)
		},
	}
}

func runList(f *factory.Factory) error {
	data, err := f.APIM().ListEnvironments()
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if printer.IsStructured(f.OutputFormat) {
		return p.PrintDetail(data)
	}

	var items []any
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return p.PrintList(items, envColumns())
}

func envColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Description", Value: func(i any) string { return cmdutil.StringField(i, "description") }},
	}
}
