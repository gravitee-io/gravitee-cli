package analytics

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type getOptions struct {
	factory  *factory.Factory
	domainID *string
	params   am.AnalyticsParams
}

func newGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &getOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get domain analytics",
		Example: `  gio am analytics get --domain my-domain --type count
  gio am analytics get --domain my-domain --type count --field status --from 2024-01-01 --to 2024-12-31`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.params.Type, "type", "", "Analytics type (e.g. count)")
	cmd.Flags().StringVar(&opts.params.Field, "field", "", "Analytics field (e.g. status)")
	cmd.Flags().StringVar(&opts.params.From, "from", "", "Start date")
	cmd.Flags().StringVar(&opts.params.To, "to", "", "End date")
	cmd.Flags().StringVar(&opts.params.Interval, "interval", "", "Interval in milliseconds")
	cmd.Flags().IntVar(&opts.params.Size, "size", 0, "Number of results")

	return cmd
}

func (o *getOptions) run() error {
	f := o.factory

	data, err := f.AM().GetAnalytics(*o.domainID, o.params)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}
