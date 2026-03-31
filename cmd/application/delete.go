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
		Example: `  gio app delete aaaa1111-2222-3333-4444-555566667777`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "app delete"); err != nil {
				return err
			}

			return runDelete(f, args[0])
		},
	}

	return cmd
}

func runDelete(f *factory.Factory, appID string) error {
	path := cmdutil.V1EnvPath(f, fmt.Sprintf("applications/%s", appID))

	if err := f.Client.Delete(path); err != nil {
		return fmt.Errorf("application deletion failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)
	p.PrintMessage("Application '%s' deleted.", appID)

	return nil
}
