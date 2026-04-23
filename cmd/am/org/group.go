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

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newOrgGroupCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "group",
		Aliases: []string{"groups"},
		Short:   "Manage organization groups",
	}

	cmd.AddCommand(newOrgGroupListCmd(f))
	cmd.AddCommand(newOrgGroupGetCmd(f))
	cmd.AddCommand(newOrgGroupCreateCmd(f))
	cmd.AddCommand(newOrgGroupUpdateCmd(f))
	cmd.AddCommand(newOrgGroupDeleteCmd(f))

	return cmd
}

// list

func newOrgGroupListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List organization groups",
		Example: `  gio am org group list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgGroupList(f)
		},
	}
}

func runOrgGroupList(f *factory.Factory) error {
	data, err := f.AM().ListOrgGroups()
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}

// get

func newOrgGroupGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <groupID>",
		Short:   "Get organization group details",
		Example: `  gio am org group get group-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgGroupGet(f, args[0])
		},
	}
}

func runOrgGroupGet(f *factory.Factory, groupID string) error {
	data, err := f.AM().GetOrgGroup(groupID)
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

	return printOrgGroupDetail(p, data)
}

func printOrgGroupDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"Description", "description"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}

// create

func newOrgGroupCreateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "create --file <group.json>",
		Short: "Create an organization group from a JSON file",
		Example: `  gio am org group create --file group.json
  gio am org group create -f group.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgGroupCreate(f, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runOrgGroupCreate(f *factory.Factory, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().CreateOrgGroup(body)
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

	return printOrgGroupDetail(p, data)
}

// update

func newOrgGroupUpdateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <groupID> --file <group.json>",
		Short: "Update an organization group from a JSON file",
		Example: `  gio am org group update group-id --file group.json
  gio am org group update group-id -f group.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgGroupUpdate(f, args[0], file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runOrgGroupUpdate(f *factory.Factory, groupID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateOrgGroup(groupID, body)
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

	return printOrgGroupDetail(p, data)
}

// delete

func newOrgGroupDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <groupID>",
		Short:   "Delete an organization group",
		Example: `  gio am org group delete group-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgGroupDelete(f, args[0])
		},
	}
}

func runOrgGroupDelete(f *factory.Factory, groupID string) error {
	if err := f.AM().DeleteOrgGroup(groupID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Organization group '%s' deleted.", groupID)

	return nil
}
