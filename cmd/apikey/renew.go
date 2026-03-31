package apikey

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type renewOptions struct {
	factory      *factory.Factory
	apiID        string
	subscription string
}

func newRenewCmd(f *factory.Factory) *cobra.Command {
	opts := &renewOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "renew --api <apiId> --subscription <subId>",
		Short: "Generate a new API key for a subscription",
		Example: `  gio api-key renew --api 8a7b3c4d --subscription aaaa1111
  gio api-key renew --api 8a7b3c4d --subscription aaaa1111 -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "api-key renew"); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.apiID, "api", "", "API ID (required)")
	cmd.Flags().StringVar(&opts.subscription, "subscription", "", "Subscription ID (required)")

	_ = cmd.MarkFlagRequired("api")
	_ = cmd.MarkFlagRequired("subscription")

	return cmd
}

func (o *renewOptions) run() error {
	f := o.factory
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/subscriptions/%s/api-keys/_renew", o.apiID, o.subscription))

	data, err := f.Client.Post(path, nil)
	if err != nil {
		return fmt.Errorf("API key renew failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printKeyDetail(p, data)
}

func printKeyDetail(p *printer.Printer, data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Key", "key"},
		{"Subscription", "subscription"},
		{"API", "api"},
		{"Revoked", "revoked"},
		{"Expired", "expired"},
		{"Created", "createdAt"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
