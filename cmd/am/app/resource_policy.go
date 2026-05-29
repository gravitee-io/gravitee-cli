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

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newAppResourcePolicyCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var (
		appID      string
		resourceID string
	)

	cmd := &cobra.Command{
		Use:   "resource-policy",
		Short: "Manage application resource policies",
	}

	cmd.PersistentFlags().StringVar(&appID, "app-id", "", "Application ID (required)")
	cmd.PersistentFlags().StringVar(&resourceID, "resource-id", "", "Resource ID (required)")
	_ = cmd.MarkPersistentFlagRequired("app-id")
	_ = cmd.MarkPersistentFlagRequired("resource-id")

	cmd.AddCommand(newAppResourcePolicyListCmd(f, domainID, &appID, &resourceID))
	cmd.AddCommand(newAppResourcePolicyGetCmd(f, domainID, &appID, &resourceID))

	return cmd
}

func newAppResourcePolicyListCmd(f *factory.Factory, domainID, appID, resourceID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List application resource policies",
		Example: `  gctl am app resource-policy list --domain my-domain --app-id my-app --resource-id res-1`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListAppResourcePolicies(*domainID, *appID, *resourceID)
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
}

func newAppResourcePolicyGetCmd(f *factory.Factory, domainID, appID, resourceID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get <policyID>",
		Short:   "Get an application resource policy",
		Example: `  gctl am app resource-policy get policy-1 --domain my-domain --app-id my-app --resource-id res-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetAppResourcePolicy(*domainID, *appID, *resourceID, args[0])
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
}
