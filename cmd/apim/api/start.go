package api

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newStartCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "start <apiId>",
		Short:   "Start an API",
		Example: `  gio apim api start 8a7b3c4d-1234-5678-abcd-ef0123456789`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.APIM().StartAPI(args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}
			p.PrintMessage("API '%s' started.", args[0])

			return nil
		},
	}
}
