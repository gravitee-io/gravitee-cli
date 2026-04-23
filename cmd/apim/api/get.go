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

package api

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <apiId>",
		Short:   "Get API details",
		Example: `  gio apim api get /my/api`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.RequireNonEmpty("apiId", args[0]); err != nil {
				return err
			}

			apiID, err := f.APIM().ResolveAPI(args[0])
			if err != nil {
				return err
			}

			return runGet(f, apiID)
		},
	}
}

func runGet(f *factory.Factory, apiID string) error {
	data, err := f.APIM().GetAPI(apiID)
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

	return printAPIDetail(p, data)
}

func printAPIDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fields := []struct {
		label string
		key   string
	}{
		{"Name", "name"},
		{"ID", "id"},
		{"Status", "state"},
		{"Definition", "definitionVersion"},
		{"Type", "type"},
		{"Visibility", "visibility"},
		{"Lifecycle", "lifecycleState"},
		{"Updated", "updatedAt"},
	}

	for _, field := range fields {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
