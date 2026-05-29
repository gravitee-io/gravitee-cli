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

package org

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newOrgMemberCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "member",
		Aliases: []string{"members"},
		Short:   "Manage organization members",
	}

	cmd.AddCommand(newOrgMemberListCmd(f))
	cmd.AddCommand(newOrgMemberAddCmd(f))
	cmd.AddCommand(newOrgMemberRemoveCmd(f))

	return cmd
}

// list

func newOrgMemberListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List organization members",
		Example: `  gctl am org member list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgMemberList(f)
		},
	}
}

func runOrgMemberList(f *factory.Factory) error {
	data, err := f.AM().ListOrgMembers()
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

	memberships, err := extractOrgMemberships(data)
	if err != nil {
		return err
	}

	return p.PrintList(memberships, orgMemberColumns())
}

func extractOrgMemberships(data json.RawMessage) ([]json.RawMessage, error) {
	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse member response: %w", err)
	}

	raw, ok := wrapper["memberships"]
	if !ok {
		return nil, nil
	}

	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("failed to parse memberships array: %w", err)
	}

	return items, nil
}

func orgMemberColumns() []printer.Column {
	return []printer.Column{
		{Name: "MemberID", Value: func(i any) string { return cmdutil.StringField(i, "memberId") }},
		{Name: "RoleID", Value: func(i any) string { return cmdutil.StringField(i, "roleId") }},
		{Name: "MemberType", Value: func(i any) string { return cmdutil.StringField(i, "memberType") }},
	}
}

// add

func newOrgMemberAddCmd(f *factory.Factory) *cobra.Command {
	var (
		memberID   string
		role       string
		memberType string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a member to the organization",
		Example: `  gctl am org member add --member-id user-123 --role role-456
  gctl am org member add --member-id user-123 --role role-456 --member-type USER`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgMemberAdd(f, memberID, role, memberType)
		},
	}

	cmd.Flags().StringVar(&memberID, "member-id", "", "Member user ID (required)")
	cmd.Flags().StringVar(&role, "role", "", "Role ID (required)")
	cmd.Flags().StringVar(&memberType, "member-type", "USER", "Member type")
	_ = cmd.MarkFlagRequired("member-id")
	_ = cmd.MarkFlagRequired("role")

	return cmd
}

func runOrgMemberAdd(f *factory.Factory, memberID, role, memberType string) error {
	body, _ := json.Marshal(map[string]string{
		"memberId":   memberID,
		"memberType": memberType,
		"role":       role,
	})

	data, err := f.AM().AddOrgMember(body)
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

	p.PrintMessage("Member '%s' added to organization.", memberID)

	return nil
}

// remove

func newOrgMemberRemoveCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "remove <memberID>",
		Short:   "Remove a member from the organization",
		Example: `  gctl am org member remove member-123`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgMemberRemove(f, args[0])
		},
	}
}

func runOrgMemberRemove(f *factory.Factory, memberID string) error {
	if err := f.AM().RemoveOrgMember(memberID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Member '%s' removed from organization.", memberID)

	return nil
}
