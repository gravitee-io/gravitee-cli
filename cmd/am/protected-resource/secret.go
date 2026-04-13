package protectedresource

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newSecretCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var resourceID string

	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage protected resource secrets",
	}

	cmd.PersistentFlags().StringVar(&resourceID, "resource-id", "", "Protected resource ID (required)")
	_ = cmd.MarkPersistentFlagRequired("resource-id")

	cmd.AddCommand(newSecretListCmd(f, domainID, &resourceID))

	return cmd
}

func newSecretListCmd(f *factory.Factory, domainID, resourceID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List protected resource secrets",
		Example: `  gio am protected-resource secret list --domain my-domain --resource-id pr-1`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListProtectedResourceSecrets(*domainID, *resourceID)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			return p.PrintDetail(data)
		},
	}
}
