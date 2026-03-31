package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newHealthCmd(f *factory.Factory) *cobra.Command {
	var field string

	cmd := &cobra.Command{
		Use:     "health <apiId>",
		Short:   "Get API health check availability",
		Example: `  gio api health 8a7b3c4d-1234-5678-abcd-ef0123456789`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runHealth(f, args[0], field)
		},
	}

	cmd.Flags().StringVar(&field, "field", "endpoint", "Grouping field")

	return cmd
}

func runHealth(f *factory.Factory, apiID, field string) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/health/availability?field=%s", apiID, url.QueryEscape(field)))

	data, err := f.Client.Get(path)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		p := cmdutil.NewPrinter(f)
		p.PrintMessage("No health check data available for this API.")

		return nil
	}

	p := cmdutil.NewPrinter(f)

	return p.PrintDetail(json.RawMessage(data))
}
