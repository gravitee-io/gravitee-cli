package api

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newStopCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "stop <apiId>",
		Short:   "Stop an API",
		Example: `  gio apim api stop /my/api`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			apiID, err := f.APIM().ResolveAPI(args[0])
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if err := f.APIM().StopAPI(apiID); err != nil {
				return err
			}

			return cmdutil.PrintActionResult(p, apiID, "stopped",
				fmt.Sprintf("API '%s' stopped.", apiID))
		},
	}
}
