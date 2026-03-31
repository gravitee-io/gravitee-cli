package metadata

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "delete <key> --api <apiId>",
		Short:   "Delete a metadata entry",
		Example: `  gio metadata delete team-email --api 8a7b3c4d-1234-5678-abcd-ef0123456789`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "metadata delete"); err != nil {
				return err
			}

			return runDelete(f, apiID, args[0])
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")

	return cmd
}

func runDelete(f *factory.Factory, apiID, key string) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/metadata/%s", apiID, key))

	if err := f.Client.Delete(path); err != nil {
		return fmt.Errorf("metadata deletion failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)
	p.PrintMessage("Metadata '%s' deleted.", key)

	return nil
}
