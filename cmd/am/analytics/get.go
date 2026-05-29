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

	"gravitee.io/gctl/internal/am"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
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
		Example: `  gctl am analytics get --domain my-domain --type count
  gctl am analytics get --domain my-domain --type count --field status --from 2024-01-01 --to 2024-12-31`,
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
