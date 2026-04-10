package api

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newHealthCmd(f *factory.Factory) *cobra.Command {
	var field string

	cmd := &cobra.Command{
		Use:     "health <apiId>",
		Short:   "Get API health check availability",
		Example: `  gio apim api health /my/api`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			apiID, err := f.APIM().ResolveAPI(args[0])
			if err != nil {
				return err
			}

			return runHealth(f, apiID, field)
		},
	}

	cmd.Flags().StringVar(&field, "field", "endpoint", "Grouping field")

	return cmd
}

func runHealth(f *factory.Factory, apiID, field string) error {
	data, err := f.APIM().GetAPIHealth(apiID, field)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		p.PrintMessage("No health check data available for this API.")

		return nil
	}

	return p.PrintDetail(data)
}
