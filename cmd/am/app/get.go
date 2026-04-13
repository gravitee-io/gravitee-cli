package app

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type getOptions struct {
	factory  *factory.Factory
	domainID *string
	appID    string
}

func newGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &getOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "get <appID>",
		Short: "Get application details",
		Example: `  gio am app get my-app-id --domain my-domain
  gio am app get my-app-id --domain my-domain -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.appID = args[0]

			return opts.run()
		},
	}

	return cmd
}

func (o *getOptions) run() error {
	f := o.factory

	data, err := f.AM().GetApplication(*o.domainID, o.appID)
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

	return printAppDetail(p, data)
}

func printAppDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"Type", "type"},
		{"Enabled", "enabled"},
		{"Description", "description"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
