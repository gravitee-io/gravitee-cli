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

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCopyCmd(f *factory.Factory) *cobra.Command {
	var targetName string
	cmd := &cobra.Command{
		Use:     "copy <sourceDomainId>",
		Short:   "Copy a domain to a new domain in the same workspace",
		Example: `  gio am domain copy abc-123 --name my-copy`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			return runCopy(f, args[0], targetName)
		},
	}
	cmd.Flags().StringVar(&targetName, "name", "", "Name for the new domain (required)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func runCopy(f *factory.Factory, sourceDomainID, targetName string) error {
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	body, err := json.Marshal(map[string]interface{}{"name": targetName})
	if err != nil {
		return err
	}
	created, err := f.Client.Post(cmdutil.AMEnvPath(f, "domains"), body)
	if err != nil {
		return err
	}
	var newDomain map[string]interface{}
	if parseErr := json.Unmarshal(created, &newDomain); parseErr != nil {
		return fmt.Errorf("failed to parse CreateDomain response: %w", parseErr)
	}
	targetDomainID := cmdutil.StringField(newDomain, "id")
	if targetDomainID == "" {
		return fmt.Errorf("CreateDomain response did not include an ID")
	}

	p.PrintMessage("Created domain '%s' (%s). Copying resources...", targetName, targetDomainID)

	exported, err := exportToMemory(f, sourceDomainID)
	if err != nil {
		return fmt.Errorf("failed to export source domain (new domain '%s' was created but is empty — delete it manually if not needed): %w", targetDomainID, err)
	}

	totalImported, totalSkipped := 0, 0
	var allErrs []error
	for _, kind := range []string{"scopes", "roles", "groups", "applications"} {
		imported, skipped, errs := importItems(f, exported, kind, targetDomainID, kind)
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

	p.PrintMessage("Copy complete: %d imported, %d skipped.", totalImported, totalSkipped)
	if totalSkipped > 0 {
		return fmt.Errorf("copy partially failed: %d items skipped", totalSkipped)
	}
	return nil
}
