package domain

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domainID>",
		Short:   "Delete a security domain",
		Example: `  gio am domain delete my-domain-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDelete(f, args[0])
		},
	}
}

func runDelete(f *factory.Factory, domainID string) error {
	if err := f.AM().DeleteDomain(domainID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Domain '%s' deleted.", domainID)

	return nil
}
