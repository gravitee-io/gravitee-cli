package api

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/apim"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type analyticsOptions struct {
	factory      *factory.Factory
	aggregations string
	apiID        string
	field        string
	analytType   string
	ranges       string
	order        string
	query        string
	terms        []string
	from         int64
	to           int64
	interval     int64
	size         int
}

func newAnalyticsCmd(f *factory.Factory) *cobra.Command {
	opts := &analyticsOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "analytics <apiId>",
		Short:   "Get API analytics",
		Example: `  gio apim api analytics 8a7b3c4d-... --type COUNT --from 1711497600000 --to 1711584000000`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.apiID = args[0]

			return opts.run()
		},
	}

	cmd.Flags().Int64Var(&opts.from, "from", 0, "Start timestamp (epoch millis)")
	cmd.Flags().Int64Var(&opts.to, "to", 0, "End timestamp (epoch millis)")
	cmd.Flags().Int64Var(&opts.interval, "interval", 0, "Interval in milliseconds")
	cmd.Flags().StringVar(&opts.field, "field", "", "Aggregation field")
	cmd.Flags().StringVar(&opts.analytType, "type", "", "Analytics type: STATS, COUNT, HISTOGRAM, GROUP_BY")
	cmd.Flags().IntVar(&opts.size, "size", 0, "Result set size")
	cmd.Flags().StringVar(&opts.ranges, "ranges", "", "Aggregation ranges (e.g. 100:199;200:299)")
	cmd.Flags().StringVar(&opts.aggregations, "aggregations", "", "Aggregations (e.g. avg:response-time)")
	cmd.Flags().StringVar(&opts.order, "order", "", "Sort order")
	cmd.Flags().StringVar(&opts.query, "query", "", "Custom search query")
	cmd.Flags().StringArrayVar(&opts.terms, "terms", nil, "Filters (e.g. plan-id:xyz)")

	return cmd
}

func (o *analyticsOptions) run() error {
	data, err := o.factory.APIM().GetAPIAnalytics(o.apiID, apim.AnalyticsParams{
		Terms:        o.terms,
		Field:        o.field,
		Type:         o.analytType,
		Ranges:       o.ranges,
		Aggregations: o.aggregations,
		Order:        o.order,
		Query:        o.query,
		From:         o.from,
		To:           o.to,
		Interval:     o.interval,
		Size:         o.size,
	})
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(o.factory)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}
