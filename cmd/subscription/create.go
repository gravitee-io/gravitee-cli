package subscription

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type createOptions struct {
	factory *factory.Factory
	apiID   string
	planID  string
	appID   string
	apiKey  string
}

func newCreateCmd(f *factory.Factory) *cobra.Command {
	opts := &createOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "create --api <apiId> --plan <planId> --app <appId>",
		Short:   "Create a subscription",
		Example: `  gio subscription create --api 8a7b3c4d --plan a1b2c3d4 --app e5f6a7b8`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "subscription create"); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")
	cmd.Flags().StringVar(&opts.planID, "plan", "", "Plan ID (required)")
	_ = cmd.MarkFlagRequired("plan")
	cmd.Flags().StringVar(&opts.appID, "app", "", "Application ID (required)")
	_ = cmd.MarkFlagRequired("app")
	cmd.Flags().StringVar(&opts.apiKey, "api-key", "", "Custom API key")

	return cmd
}

func (o *createOptions) run() error {
	f := o.factory
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/subscriptions", o.apiID))

	body := map[string]string{
		"planId":        o.planID,
		"applicationId": o.appID,
	}

	if o.apiKey != "" {
		body["customApiKey"] = o.apiKey
	}

	data, err := f.Client.Post(path, body)
	if err != nil {
		return fmt.Errorf("subscription creation failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printSubCreateDetail(p, data)
}
