package apikey

import (
	"fmt"

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
		Example: `  gio api-key revoke 1a2b3c4d --api 8a7b3c4d --subscription aaaa1111`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "api-key revoke"); err != nil {
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
	path := cmdutil.V2EnvPath(f, fmt.Sprintf(
		"apis/%s/subscriptions/%s/api-keys/%s/_revoke",
		o.apiID, o.subscription, keyID,
	))

	if _, err := f.Client.Post(path, nil); err != nil {
		return fmt.Errorf("API key revoke failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)
	p.PrintMessage("API key '%s' revoked.", keyID)

	return nil
}
