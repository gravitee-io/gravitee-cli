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

package app

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newAppAnalyticsCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var appID string

	cmd := &cobra.Command{
		Use:   "analytics",
		Short: "Application analytics",
	}

	cmd.PersistentFlags().StringVar(&appID, "app-id", "", "Application ID (required)")
	_ = cmd.MarkPersistentFlagRequired("app-id")

	cmd.AddCommand(newAppAnalyticsGetCmd(f, domainID, &appID))

	return cmd
}

func newAppAnalyticsGetCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	var (
		analyticsType string
		field         string
		from          string
		to            string
		interval      string
		size          int
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get application analytics",
		Example: `  gio am app analytics get --domain my-domain --app-id my-app --type count --field application_id
  gio am app analytics get --domain my-domain --app-id my-app --type count --from 1609459200000 --to 1612137600000`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			params := am.AnalyticsParams{
				Type:     analyticsType,
				Field:    field,
				From:     from,
				To:       to,
				Interval: interval,
				Size:     size,
			}

			data, err := f.AM().GetAppAnalytics(*domainID, *appID, params)
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

	cmd.Flags().StringVar(&analyticsType, "type", "", "Analytics type (e.g. count, date_histo, group_by)")
	cmd.Flags().StringVar(&field, "field", "", "Field to aggregate on")
	cmd.Flags().StringVar(&from, "from", "", "Start timestamp (epoch ms)")
	cmd.Flags().StringVar(&to, "to", "", "End timestamp (epoch ms)")
	cmd.Flags().StringVar(&interval, "interval", "", "Time interval")
	cmd.Flags().IntVar(&size, "size", 0, "Number of results")

	return cmd
}
