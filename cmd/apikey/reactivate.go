package apikey

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type reactivateOptions struct {
	factory      *factory.Factory
	apiID        string
	subscription string
}

func newReactivateCmd(f *factory.Factory) *cobra.Command {
	opts := &reactivateOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "reactivate <keyId> --api <apiId> --subscription <subId>",
		Short:   "Reactivate a previously revoked API key",
		Example: `  gio api-key reactivate 1a2b3c4d --api 8a7b3c4d --subscription aaaa1111`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "api-key reactivate"); err != nil {
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

func (o *reactivateOptions) run(keyID string) error {
	f := o.factory
	path := cmdutil.V2EnvPath(f, fmt.Sprintf(
		"apis/%s/subscriptions/%s/api-keys/%s/_reactivate",
		o.apiID, o.subscription, keyID,
	))

	data, err := f.Client.Post(path, nil)
	if err != nil {
		return fmt.Errorf("API key reactivate failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printKeyDetail(p, data)
}
