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

package subscription

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/apim"
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
		Example: `  gio apim subscription create --api 8a7b3c4d --plan a1b2c3d4 --app e5f6a7b8`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmdutil.AddAPIFlag(cmd, &opts.apiID)
	cmd.Flags().StringVar(&opts.planID, "plan", "", "Plan ID (required)")
	_ = cmd.MarkFlagRequired("plan")
	cmd.Flags().StringVar(&opts.appID, "app", "", "Application ID (required)")
	_ = cmd.MarkFlagRequired("app")
	cmd.Flags().StringVar(&opts.apiKey, "api-key", "", "Custom API key")

	return cmd
}

func (o *createOptions) run() error {
	f := o.factory

	data, err := f.APIM().CreateSubscription(o.apiID, apim.CreateSubscriptionBody{
		PlanID:       o.planID,
		AppID:        o.appID,
		CustomAPIKey: o.apiKey,
	})
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

	return printSubCreateDetail(p, data)
}
