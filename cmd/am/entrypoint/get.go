package entrypoint

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Short:   "Get domain entrypoints",
		Example: `  gio am entrypoint get --domain my-domain`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runGet(f, *domainID)
		},
	}
}

func runGet(f *factory.Factory, domainID string) error {
	data, err := f.AM().GetEntrypoints(domainID)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}
