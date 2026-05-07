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

package alert

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newTriggerCmd(f *factory.Factory, domainID *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trigger",
		Aliases: []string{"triggers"},
		Short:   "Manage alert triggers",
	}

	cmd.AddCommand(newTriggerGetCmd(f, domainID))
	cmd.AddCommand(newTriggerUpdateCmd(f, domainID))

	return cmd
}

func newTriggerGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Short:   "Get alert triggers",
		Example: `  gio am alert trigger get --domain my-domain`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runTriggerGet(f, *domainID)
		},
	}
}

func runTriggerGet(f *factory.Factory, domainID string) error {
	data, err := f.AM().GetAlertTriggers(domainID)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}

func newTriggerUpdateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update [-f <file>]",
		Short: "Update alert triggers from a JSON file or stdin",
		Example: `  gio am alert trigger update --domain my-domain --file triggers.json
  gio am alert trigger update --domain my-domain -f triggers.json
  envsubst < triggers.json | gio am alert trigger update --domain my-domain`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runTriggerUpdate(f, *domainID, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a JSON file (optional - reads from stdin if omitted)")

	return cmd
}

func runTriggerUpdate(f *factory.Factory, domainID, file string) error {
	body, err := cmdutil.ReadJSONInput(file, f.IOStreams.In)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateAlertTriggers(domainID, body)
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

	p.PrintMessage("Alert triggers updated successfully.")

	return nil
}
