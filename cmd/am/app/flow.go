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

func newAppFlowCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var appID string

	cmd := &cobra.Command{
		Use:   "flow",
		Short: "Manage application flows",
	}

	cmd.PersistentFlags().StringVar(&appID, "app-id", "", "Application ID (required)")
	_ = cmd.MarkPersistentFlagRequired("app-id")

	cmd.AddCommand(newAppFlowListCmd(f, domainID, &appID))
	cmd.AddCommand(newAppFlowGetCmd(f, domainID, &appID))
	cmd.AddCommand(newAppFlowUpdateCmd(f, domainID, &appID))

	return cmd
}

func newAppFlowListCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List application flows",
		Example: `  gio am app flow list --domain my-domain --app-id my-app`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			items, err := f.AM().ListAppFlows(*domainID, *appID)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			raw, _ := json.Marshal(items)

			return p.PrintDetail(json.RawMessage(raw))
		},
	}
}

func newAppFlowGetCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get <flowID>",
		Short:   "Get an application flow",
		Example: `  gio am app flow get flow-1 --domain my-domain --app-id my-app`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetAppFlow(*domainID, *appID, args[0])
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

func newAppFlowUpdateCmd(f *factory.Factory, domainID, appID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update --file <flows.json>",
		Short: "Update application flows from a JSON file (bulk update)",
		Example: `  gio am app flow update --domain my-domain --app-id my-app --file flows.json
  gio am app flow update --domain my-domain --app-id my-app -f flows.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONFile(file)
			if err != nil {
				return err
			}

			data, err := f.AM().UpdateAppFlows(*domainID, *appID, body)
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

			p.PrintMessage("Application flows updated successfully.")

			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
