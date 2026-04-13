package entrypoint

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory, _ *string) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <entrypointID>",
		Short:   "Delete an entrypoint",
		Example: `  gio am entrypoint delete my-entrypoint-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteEntrypoint(args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Entrypoint '%s' deleted.", args[0])

			return nil
		},
	}
}
