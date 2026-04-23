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

package scope

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type getOptions struct {
	factory  *factory.Factory
	domainID *string
	scopeID  string
}

func newGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &getOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "get <scopeID>",
		Short: "Get scope details",
		Example: `  gio am scope get my-scope-id --domain my-domain
  gio am scope get my-scope-id --domain my-domain -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.scopeID = args[0]

			return opts.run()
		},
	}

	return cmd
}

func (o *getOptions) run() error {
	f := o.factory

	data, err := f.AM().GetScope(*o.domainID, o.scopeID)
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

	return printScopeDetail(p, data)
}

func printScopeDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Key", "key"},
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
