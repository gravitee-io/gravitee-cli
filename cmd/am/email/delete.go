package email

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <emailID>",
		Short:   "Delete an email",
		Example: `  gio am email delete my-email-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDelete(f, *domainID, args[0])
		},
	}
}

func runDelete(f *factory.Factory, domainID, emailID string) error {
	if err := f.AM().DeleteEmail(domainID, emailID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Email '%s' deleted.", emailID)

	return nil
}
