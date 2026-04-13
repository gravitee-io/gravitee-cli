package member

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newRemoveCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "remove <memberID>",
		Short:   "Remove a member from a domain",
		Example: `  gio am member remove member-123 --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runRemove(f, *domainID, args[0])
		},
	}
}

func runRemove(f *factory.Factory, domainID, memberID string) error {
	if err := f.AM().RemoveMember(domainID, memberID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Member '%s' removed.", memberID)

	return nil
}
