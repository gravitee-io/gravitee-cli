package passwordpolicy

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <policyID>",
		Short:   "Delete a password policy",
		Example: `  gio am password-policy delete my-policy-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDelete(f, *domainID, args[0])
		},
	}
}

func runDelete(f *factory.Factory, domainID, policyID string) error {
	if err := f.AM().DeletePasswordPolicy(domainID, policyID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Password policy '%s' deleted.", policyID)

	return nil
}
