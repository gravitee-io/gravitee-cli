package subscription

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type acceptOptions struct {
	factory    *factory.Factory
	apiID      string
	reason     string
	startingAt string
	endingAt   string
	apiKey     string
}

func newAcceptCmd(f *factory.Factory) *cobra.Command {
	opts := &acceptOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "accept <subId> --api <apiId>",
		Short:   "Accept a pending subscription",
		Example: `  gio subscription accept cc556677 --api 8a7b3c4d --reason "Approved by ops team"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "subscription accept"); err != nil {
				return err
			}

			return opts.run(args[0])
		},
	}

	cmd.Flags().StringVar(&opts.apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")
	cmd.Flags().StringVar(&opts.reason, "reason", "", "Reason for accepting")
	cmd.Flags().StringVar(&opts.startingAt, "starting-at", "", "Start date (ISO 8601)")
	cmd.Flags().StringVar(&opts.endingAt, "ending-at", "", "End date (ISO 8601)")
	cmd.Flags().StringVar(&opts.apiKey, "api-key", "", "Custom API key")

	return cmd
}

func (o *acceptOptions) run(subID string) error {
	f := o.factory
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/subscriptions/%s/_accept", o.apiID, subID))

	body := make(map[string]string)

	if o.reason != "" {
		body["reason"] = o.reason
	}

	if o.startingAt != "" {
		body["startingAt"] = o.startingAt
	}

	if o.endingAt != "" {
		body["endingAt"] = o.endingAt
	}

	if o.apiKey != "" {
		body["customApiKey"] = o.apiKey
	}

	data, err := f.Client.Post(path, body)
	if err != nil {
		return fmt.Errorf("subscription accept failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printSubDetail(p, data)
}
