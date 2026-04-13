package resource

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <resourceID>",
		Short:   "Delete a resource",
		Example: `  gio am resource delete my-resource-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDelete(f, *domainID, args[0])
		},
	}
}

func runDelete(f *factory.Factory, domainID, resourceID string) error {
	if err := f.AM().DeleteResource(domainID, resourceID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Resource '%s' deleted.", resourceID)

	return nil
}
