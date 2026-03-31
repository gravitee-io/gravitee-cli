package member

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newRemoveCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "remove <memberId> --api <apiId>",
		Short:   "Remove a member from an API",
		Example: `  gio member remove bbbb1111-2222-3333-4444-555566667777 --api 8a7b3c4d-1234-5678-abcd-ef0123456789`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "member remove"); err != nil {
				return err
			}

			return runRemove(f, apiID, args[0])
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")

	return cmd
}

func runRemove(f *factory.Factory, apiID, memberID string) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/members/%s", apiID, memberID))

	if err := f.Client.Delete(path); err != nil {
		return fmt.Errorf("member removal failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)
	p.PrintMessage("Member '%s' removed.", memberID)

	return nil
}
