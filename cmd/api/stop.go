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
		Example: `  gio api stop 8a7b3c4d-1234-5678-abcd-ef0123456789`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "api stop"); err != nil {
				return err
			}

			path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/_stop", args[0]))

			if _, err := f.Client.Post(path, nil); err != nil {
				return fmt.Errorf("API stop failed: %w", err)
			}

			p := cmdutil.NewPrinter(f)
			p.PrintMessage("API '%s' stopped.", args[0])

			return nil
		},
	}
}
