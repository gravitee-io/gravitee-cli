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
		Example: `  gio apim member remove bbbb1111-2222-3333-4444-555566667777 --api /my/api`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runRemove(f, apiID, args[0])
		},
	}

	cmdutil.AddAPIFlag(cmd, &apiID)

	return cmd
}

func runRemove(f *factory.Factory, apiID, memberID string) error {
	if err := f.APIM().RemoveMember(apiID, memberID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return cmdutil.PrintActionResult(p, memberID, "removed",
		fmt.Sprintf("Member '%s' removed.", memberID))
}
