package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

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
		Example: `  gio api analytics 8a7b3c4d-... --type COUNT --from 1711497600000 --to 1711584000000`,
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
	q := o.buildQuery()
	path := cmdutil.V2EnvPath(o.factory, fmt.Sprintf("apis/%s/analytics?%s", o.apiID, q))

	data, err := o.factory.Client.Get(path)
	if err != nil {
		return err
	}

	p := cmdutil.NewPrinter(o.factory)

	return p.PrintDetail(json.RawMessage(data))
}

func (o *analyticsOptions) buildQuery() string {
	q := url.Values{}

	if o.from != 0 {
		q.Set("from", strconv.FormatInt(o.from, 10))
	}

	if o.to != 0 {
		q.Set("to", strconv.FormatInt(o.to, 10))
	}

	if o.interval != 0 {
		q.Set("interval", strconv.FormatInt(o.interval, 10))
	}

	if o.field != "" {
		q.Set("field", o.field)
	}

	if o.analytType != "" {
		q.Set("type", o.analytType)
	}

	if o.size != 0 {
		q.Set("size", strconv.Itoa(o.size))
	}

	if o.ranges != "" {
		q.Set("ranges", o.ranges)
	}

	if o.aggregations != "" {
		q.Set("aggregations", o.aggregations)
	}

	if o.order != "" {
		q.Set("order", o.order)
	}

	if o.query != "" {
		q.Set("query", o.query)
	}

	for _, term := range o.terms {
		q.Add("terms", term)
	}

	return q.Encode()
}
