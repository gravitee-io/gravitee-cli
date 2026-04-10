package subscription

import (
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
		Example: `  gio apim subscription reject cc556677 --api 8a7b3c4d --reason "Insufficient justification"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run(args[0])
		},
	}

	cmdutil.AddAPIFlag(cmd, &opts.apiID)
	cmd.Flags().StringVar(&opts.reason, "reason", "", "Reason for rejecting")

	return cmd
}

func (o *rejectOptions) run(subID string) error {
	f := o.factory

	data, err := f.APIM().RejectSubscription(o.apiID, subID, o.reason)
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
