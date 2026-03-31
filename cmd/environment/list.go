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
		Example: `  gio environment list
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
	path := fmt.Sprintf("/management/organizations/%s/environments", f.Resolved.Org)

	data, err := f.Client.Get(path)
	if err != nil {
		return err
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		var raw json.RawMessage = data

		return p.PrintDetail(raw)
	}

	var items []interface{}
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return p.PrintList(items, envColumns())
}

func envColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
		{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
		{Name: "Description", Value: func(i interface{}) string { return cmdutil.StringField(i, "description") }},
	}
}
