package application

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <appId>",
		Short:   "Delete an application",
		Example: `  gio apim app delete aaaa1111-2222-3333-4444-555566667777`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDelete(f, args[0])
		},
	}

	return cmd
}

func runDelete(f *factory.Factory, appID string) error {
	if err := f.APIM().DeleteApplication(appID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return cmdutil.PrintActionResult(p, appID, "deleted",
		fmt.Sprintf("Application '%s' deleted.", appID))
}
