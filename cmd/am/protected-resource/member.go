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

func newMemberCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var resourceID string

	cmd := &cobra.Command{
		Use:   "member",
		Short: "Manage protected resource members",
	}

	cmd.PersistentFlags().StringVar(&resourceID, "resource-id", "", "Protected resource ID (required)")
	_ = cmd.MarkPersistentFlagRequired("resource-id")

	cmd.AddCommand(newMemberListCmd(f, domainID, &resourceID))
	cmd.AddCommand(newMemberRemoveCmd(f, domainID, &resourceID))

	return cmd
}

func newMemberListCmd(f *factory.Factory, domainID, resourceID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List protected resource members",
		Example: `  gio am protected-resource member list --domain my-domain --resource-id pr-1`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListProtectedResourceMembers(*domainID, *resourceID)
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

func newMemberRemoveCmd(f *factory.Factory, domainID, resourceID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "remove <memberID>",
		Short:   "Remove a member from a protected resource",
		Example: `  gio am protected-resource member remove member-1 --domain my-domain --resource-id pr-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RemoveProtectedResourceMember(*domainID, *resourceID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Member '%s' removed.", args[0])

			return nil
		},
	}
}
