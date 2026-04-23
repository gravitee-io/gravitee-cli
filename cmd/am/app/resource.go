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

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newAppResourceCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var appID string

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Manage application resources",
	}

	cmd.PersistentFlags().StringVar(&appID, "app-id", "", "Application ID (required)")
	_ = cmd.MarkPersistentFlagRequired("app-id")

	cmd.AddCommand(newAppResourceListCmd(f, domainID, &appID))
	cmd.AddCommand(newAppResourceGetCmd(f, domainID, &appID))

	return cmd
}

func newAppResourceListCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List application resources",
		Example: `  gio am app resource list --domain my-domain --app-id my-app`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListAppResources(*domainID, *appID)
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

func newAppResourceGetCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get <resourceID>",
		Short:   "Get an application resource",
		Example: `  gio am app resource get res-1 --domain my-domain --app-id my-app`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetAppResource(*domainID, *appID, args[0])
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
