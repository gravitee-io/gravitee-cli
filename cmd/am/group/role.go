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

package group

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newRoleCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage group roles",
	}

	cmd.PersistentFlags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	_ = cmd.MarkPersistentFlagRequired("group-id")

	cmd.AddCommand(newRoleListCmd(f, domainID, &groupID))
	cmd.AddCommand(newRoleAssignCmd(f, domainID, &groupID))
	cmd.AddCommand(newRoleRevokeCmd(f, domainID, &groupID))

	return cmd
}

func newRoleListCmd(f *factory.Factory, domainID, groupID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List group roles",
		Example: `  gio am group role list --domain my-domain --group-id my-group`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListGroupRoles(*domainID, *groupID)
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

func newRoleAssignCmd(f *factory.Factory, domainID, groupID *string) *cobra.Command {
	var roles string

	cmd := &cobra.Command{
		Use:     "assign --roles role1,role2",
		Short:   "Assign roles to a group",
		Example: `  gio am group role assign --domain my-domain --group-id my-group --roles role1,role2`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			roleList := strings.Split(roles, ",")
			raw, _ := json.Marshal(roleList)

			data, err := f.AM().AssignGroupRoles(*domainID, *groupID, json.RawMessage(raw))
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

			p.PrintMessage("Roles assigned to group '%s'.", *groupID)

			return nil
		},
	}

	cmd.Flags().StringVar(&roles, "roles", "", "Comma-separated list of role IDs (required)")
	_ = cmd.MarkFlagRequired("roles")

	return cmd
}

func newRoleRevokeCmd(f *factory.Factory, domainID, groupID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "revoke <roleID>",
		Short:   "Revoke a role from a group",
		Example: `  gio am group role revoke role-123 --domain my-domain --group-id my-group`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RevokeGroupRole(*domainID, *groupID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Role '%s' revoked from group '%s'.", args[0], *groupID)

			return nil
		},
	}
}
