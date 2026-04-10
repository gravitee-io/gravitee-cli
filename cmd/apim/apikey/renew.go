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
		Example: `  gio apim api-key renew --api 8a7b3c4d --subscription aaaa1111
  gio apim api-key renew --api 8a7b3c4d --subscription aaaa1111 -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmdutil.AddAPIFlag(cmd, &opts.apiID)
	cmd.Flags().StringVar(&opts.subscription, "subscription", "", "Subscription ID (required)")

	_ = cmd.MarkFlagRequired("subscription")

	return cmd
}

func (o *renewOptions) run() error {
	f := o.factory

	data, err := f.APIM().RenewAPIKey(o.apiID, o.subscription)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(data)
	}

	return printKeyDetail(p, data)
}

func printKeyDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
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
