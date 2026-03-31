package page

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "delete <pageId> --api <apiId>",
		Short:   "Delete a page",
		Example: `  gio page delete dddd1111-2222-3333-4444-555566667777 --api 8a7b3c4d-1234-5678-abcd-ef0123456789`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "page delete"); err != nil {
				return err
			}

			return runDelete(f, apiID, args[0])
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")

	return cmd
}

func runDelete(f *factory.Factory, apiID, pageID string) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/pages/%s", apiID, pageID))

	if err := f.Client.Delete(path); err != nil {
		return fmt.Errorf("page deletion failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)
	p.PrintMessage("Page '%s' deleted.", pageID)

	return nil
}
