package subscription

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type rejectOptions struct {
	factory *factory.Factory
	apiID   string
	reason  string
}

func newRejectCmd(f *factory.Factory) *cobra.Command {
	opts := &rejectOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "reject <subId> --api <apiId>",
		Short:   "Reject a pending subscription",
		Example: `  gio subscription reject cc556677 --api 8a7b3c4d --reason "Insufficient justification"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "subscription reject"); err != nil {
				return err
			}

			return opts.run(args[0])
		},
	}

	cmd.Flags().StringVar(&opts.apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")
	cmd.Flags().StringVar(&opts.reason, "reason", "", "Reason for rejecting")

	return cmd
}

func (o *rejectOptions) run(subID string) error {
	f := o.factory
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/subscriptions/%s/_reject", o.apiID, subID))

	body := make(map[string]string)

	if o.reason != "" {
		body["reason"] = o.reason
	}

	data, err := f.Client.Post(path, body)
	if err != nil {
		return fmt.Errorf("subscription reject failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printSubDetail(p, data)
}
