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

package dictionary

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newEntryCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var dictID string

	cmd := &cobra.Command{
		Use:   "entry",
		Short: "Manage dictionary entries",
	}

	cmd.PersistentFlags().StringVar(&dictID, "dict-id", "", "Dictionary ID (required)")
	_ = cmd.MarkPersistentFlagRequired("dict-id")

	cmd.AddCommand(newEntryListCmd(f, domainID, &dictID))
	cmd.AddCommand(newEntryUpdateCmd(f, domainID, &dictID))

	return cmd
}

func newEntryListCmd(f *factory.Factory, domainID, dictID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List dictionary entries",
		Example: `  gio am dictionary entry list --domain my-domain --dict-id my-dict`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListDictionaryEntries(*domainID, *dictID)
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

func newEntryUpdateCmd(f *factory.Factory, domainID, dictID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update [-f <file>]",
		Short: "Update dictionary entries from a JSON file or stdin",
		Example: `  gio am dictionary entry update --domain my-domain --dict-id my-dict --file entries.json
  envsubst < entries.json | gio am dictionary entry update --domain my-domain --dict-id my-dict`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONInput(file, f.IOStreams.In)
			if err != nil {
				return err
			}

			data, err := f.AM().UpdateDictionaryEntries(*domainID, *dictID, body)
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

			p.PrintMessage("Dictionary entries updated.")

			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a JSON file (optional - reads from stdin if omitted)")

	return cmd
}
