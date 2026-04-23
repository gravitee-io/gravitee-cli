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
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newAppMemberCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var appID string

	cmd := &cobra.Command{
		Use:   "member",
		Short: "Manage application members",
	}

	cmd.PersistentFlags().StringVar(&appID, "app-id", "", "Application ID (required)")
	_ = cmd.MarkPersistentFlagRequired("app-id")

	cmd.AddCommand(newAppMemberListCmd(f, domainID, &appID))
	cmd.AddCommand(newAppMemberAddCmd(f, domainID, &appID))
	cmd.AddCommand(newAppMemberRemoveCmd(f, domainID, &appID))
	cmd.AddCommand(newAppMemberPermissionsCmd(f, domainID, &appID))

	return cmd
}

func newAppMemberListCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List application members",
		Example: `  gio am app member list --domain my-domain --app-id my-app`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListAppMembers(*domainID, *appID)
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

func newAppMemberAddCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	var (
		memberID   string
		memberType string
		role       string
	)

	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a member to an application",
		Example: `  gio am app member add --domain my-domain --app-id my-app --member-id user-1 --role owner`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, _ := json.Marshal(map[string]any{
				"memberId":   memberID,
				"memberType": memberType,
				"role":       role,
			})

			data, err := f.AM().AddAppMember(*domainID, *appID, json.RawMessage(body))
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

			p.PrintMessage("Member added successfully.")

			return nil
		},
	}

	cmd.Flags().StringVar(&memberID, "member-id", "", "Member ID (required)")
	cmd.Flags().StringVar(&memberType, "member-type", "USER", "Member type (USER, GROUP)")
	cmd.Flags().StringVar(&role, "role", "", "Role (required)")
	_ = cmd.MarkFlagRequired("member-id")
	_ = cmd.MarkFlagRequired("role")

	return cmd
}

func newAppMemberPermissionsCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "permissions",
		Short:   "Get application member permissions",
		Example: `  gio am app member permissions --domain my-domain --app-id my-app`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetAppMemberPermissions(*domainID, *appID)
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

func newAppMemberRemoveCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "remove <memberID>",
		Short:   "Remove a member from an application",
		Example: `  gio am app member remove member-1 --domain my-domain --app-id my-app`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RemoveAppMember(*domainID, *appID, args[0]); err != nil {
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
