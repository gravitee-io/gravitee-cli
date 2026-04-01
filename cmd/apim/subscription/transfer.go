package subscription

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type transferOptions struct {
	factory *factory.Factory
	apiID   string
	planID  string
}

func newTransferCmd(f *factory.Factory) *cobra.Command {
	opts := &transferOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "transfer <subId> --api <apiId> --plan <planId>",
		Short:   "Transfer a subscription to another plan",
		Example: `  gio apim subscription transfer 34f8c9e7 --api 8a7b3c4d --plan dd998877`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run(args[0])
		},
	}

	cmd.Flags().StringVar(&opts.apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")
	cmd.Flags().StringVar(&opts.planID, "plan", "", "Target plan ID (required)")
	_ = cmd.MarkFlagRequired("plan")

	return cmd
}

func (o *transferOptions) run(subID string) error {
	f := o.factory

	data, err := f.APIM().TransferSubscription(o.apiID, subID, o.planID)
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

	return printSubDetail(p, data)
}
