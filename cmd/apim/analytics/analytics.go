// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analytics

import (
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/apim"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

// NewAnalyticsCmd creates the analytics command.
func NewAnalyticsCmd(f *factory.Factory) *cobra.Command {
	var (
		apiID        string
		aggregations string
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
	)

	cmd := &cobra.Command{
		Use:     "analytics --api <apiId>",
		Short:   "Get API analytics",
		Example: `  gctl apim analytics --api /my/api --type COUNT --from 1711497600000 --to 1711584000000`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			resolvedID, err := f.APIM().ResolveAPI(apiID)
			if err != nil {
				return err
			}

			data, err := f.APIM().GetAPIAnalytics(resolvedID, apim.AnalyticsParams{
				Terms:        terms,
				Field:        field,
				Type:         analytType,
				Ranges:       ranges,
				Aggregations: aggregations,
				Order:        order,
				Query:        query,
				From:         from,
				To:           to,
				Interval:     interval,
				Size:         size,
			})
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			return p.PrintDetail(data)
		},
	}

	cmdutil.AddOutputFlags(cmd, f)
	cmdutil.AddAPIFlag(cmd, &apiID)
	cmd.Flags().Int64Var(&from, "from", 0, "Start timestamp (epoch millis, required)")
	cmd.Flags().Int64Var(&to, "to", 0, "End timestamp (epoch millis, required)")
	cmd.Flags().Int64Var(&interval, "interval", 0, "Interval in milliseconds")
	cmd.Flags().StringVar(&field, "field", "", "Aggregation field")
	cmd.Flags().StringVar(&analytType, "type", "", "Analytics type: STATS, COUNT, HISTOGRAM, GROUP_BY")
	cmd.Flags().IntVar(&size, "size", 0, "Result set size")
	cmd.Flags().StringVar(&ranges, "ranges", "", "Aggregation ranges (e.g. 100:199;200:299)")
	cmd.Flags().StringVar(&aggregations, "aggregations", "", "Aggregations (e.g. avg:response-time)")
	cmd.Flags().StringVar(&order, "order", "", "Sort order")
	cmd.Flags().StringVar(&query, "query", "", "Custom search query")
	cmd.Flags().StringArrayVar(&terms, "terms", nil, "Filters (e.g. plan-id:xyz)")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")

	return cmd
}
