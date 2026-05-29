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

package domain

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newImportCmd(f *factory.Factory) *cobra.Command {
	var targetDomainID string
	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import domain configuration from a JSON export file",
		Example: `  gctl am domain import domain-export.json
  gctl am domain import domain-export.json --target existing-domain-id`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			return runImport(f, args[0], targetDomainID)
		},
	}
	cmd.Flags().StringVar(&targetDomainID, "target", "", "Target domain ID (creates new domain if not set)")
	return cmd
}

func runImport(f *factory.Factory, file, targetDomainID string) error {
	raw, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var exportData map[string]json.RawMessage
	if parseErr := json.Unmarshal(raw, &exportData); parseErr != nil {
		return fmt.Errorf("failed to parse export file: %w", parseErr)
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if targetDomainID == "" {
		var domainObj map[string]interface{}
		if err := json.Unmarshal(exportData["domain"], &domainObj); err != nil {
			return fmt.Errorf("failed to parse domain in export: %w", err)
		}
		body := map[string]interface{}{
			"name":        cmdutil.StringField(domainObj, "name"),
			"description": cmdutil.StringField(domainObj, "description"),
		}
		created, err := f.Client.Post(cmdutil.AMEnvPath(f, "domains"), body)
		if err != nil {
			return fmt.Errorf("failed to create domain: %w", err)
		}
		var newDomain map[string]interface{}
		if err := json.Unmarshal(created, &newDomain); err != nil {
			return fmt.Errorf("failed to parse CreateDomain response: %w", err)
		}
		targetDomainID = cmdutil.StringField(newDomain, "id")
		if targetDomainID == "" {
			return fmt.Errorf("CreateDomain response did not include an ID")
		}
		p.PrintMessage("Created domain '%s'.", targetDomainID)
	}

	totalImported, totalSkipped := 0, 0
	var allErrs []error
	for _, kind := range []string{"scopes", "roles", "groups", "applications"} {
		imported, skipped, errs := importItems(f, exportData, kind, targetDomainID, kind)
		totalImported += imported
		totalSkipped += skipped
		allErrs = append(allErrs, errs...)
	}

	for i, err := range allErrs {
		if i >= 5 {
			fmt.Fprintf(f.IOStreams.Err, "  ... and %d more errors\n", len(allErrs)-5)
			break
		}
		fmt.Fprintf(f.IOStreams.Err, "  - %v\n", err)
	}

	p.PrintMessage("Import complete: %d imported, %d skipped.", totalImported, totalSkipped)
	if totalSkipped > 0 {
		return fmt.Errorf("import partially failed: %d items skipped", totalSkipped)
	}
	return nil
}

// importItems creates resources from a JSON array in exportData.
// Returns (imported, skipped, errors). Each failed POST is counted as skipped
// and its error is included in the returned slice so the caller can report
// real failure counts to the user.
func importItems(f *factory.Factory, exportData map[string]json.RawMessage, key, domainID, resource string) (int, int, []error) {
	raw, ok := exportData[key]
	if !ok || len(raw) == 0 {
		return 0, 0, nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return 0, 1, []error{fmt.Errorf("%s: parse failed: %w", key, err)}
	}
	imported, skipped := 0, 0
	var errs []error
	for i, item := range items {
		path := cmdutil.AMDomainPathFor(f, domainID, resource)
		if _, err := f.Client.Post(path, item); err != nil {
			skipped++
			errs = append(errs, fmt.Errorf("%s[%d]: %w", key, i, err))
			continue
		}
		imported++
	}
	return imported, skipped, errs
}
