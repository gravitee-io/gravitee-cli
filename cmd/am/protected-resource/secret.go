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

package protectedresource

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newSecretCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var resourceID string

	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage protected resource secrets",
	}

	cmd.PersistentFlags().StringVar(&resourceID, "resource-id", "", "Protected resource ID (required)")
	_ = cmd.MarkPersistentFlagRequired("resource-id")

	cmd.AddCommand(newSecretListCmd(f, domainID, &resourceID))

	return cmd
}

func newSecretListCmd(f *factory.Factory, domainID, resourceID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List protected resource secrets",
		Example: `  gio am protected-resource secret list --domain my-domain --resource-id pr-1`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListProtectedResourceSecrets(*domainID, *resourceID)
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
