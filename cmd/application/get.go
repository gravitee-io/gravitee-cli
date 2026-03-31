package application

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newGetCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <appId>",
		Short:   "Get application details",
		Example: `  gio app get aaaa1111-2222-3333-4444-555566667777`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runGet(f, args[0])
		},
	}

	return cmd
}

func runGet(f *factory.Factory, appID string) error {
	path := cmdutil.V1EnvPath(f, fmt.Sprintf("applications/%s", appID))

	data, err := f.Client.Get(path)
	if err != nil {
		return err
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printAppDetail(p, data)
}
