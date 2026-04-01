package apikey

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type revokeOptions struct {
	factory      *factory.Factory
	apiID        string
	subscription string
}

func newRevokeCmd(f *factory.Factory) *cobra.Command {
	opts := &revokeOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "revoke <keyId> --api <apiId> --subscription <subId>",
		Short:   "Revoke an API key",
		Example: `  gio apim api-key revoke 1a2b3c4d --api 8a7b3c4d --subscription aaaa1111`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run(args[0])
		},
	}

	cmd.Flags().StringVar(&opts.apiID, "api", "", "API ID (required)")
	cmd.Flags().StringVar(&opts.subscription, "subscription", "", "Subscription ID (required)")

	_ = cmd.MarkFlagRequired("api")
	_ = cmd.MarkFlagRequired("subscription")

	return cmd
}

func (o *revokeOptions) run(keyID string) error {
	f := o.factory

	if err := f.APIM().RevokeAPIKey(o.apiID, o.subscription, keyID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}
	p.PrintMessage("API key '%s' revoked.", keyID)

	return nil
}
