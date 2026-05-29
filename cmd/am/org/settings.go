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
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newOrgSettingsCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Manage organization settings",
	}

	cmd.AddCommand(newOrgSettingsGetCmd(f))
	cmd.AddCommand(newOrgSettingsUpdateCmd(f))

	return cmd
}

// get

func newOrgSettingsGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Short:   "Get organization settings",
		Example: `  gctl am org settings get`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgSettingsGet(f)
		},
	}
}

func runOrgSettingsGet(f *factory.Factory) error {
	data, err := f.AM().GetOrgSettings()
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}

// update

func newOrgSettingsUpdateCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update [-f <file>]",
		Short: "Update organization settings from a JSON file or stdin",
		Example: `  gctl am org settings update --file settings.json
  gctl am org settings update -f settings.json
  envsubst < settings.json | gctl am org settings update`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runOrgSettingsUpdate(f, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a JSON file (optional - reads from stdin if omitted)")

	return cmd
}

func runOrgSettingsUpdate(f *factory.Factory, file string) error {
	body, err := cmdutil.ReadJSONInput(file, f.IOStreams.In)
	if err != nil {
		return err
	}

	data, err := f.AM().PatchOrgSettings(body)
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

	p.PrintMessage("Organization settings updated.")

	return nil
}
