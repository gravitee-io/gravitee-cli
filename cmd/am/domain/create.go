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

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type createOptions struct {
	factory     *factory.Factory
	name        string
	description string
	dataPlaneID string
}

func newCreateCmd(f *factory.Factory) *cobra.Command {
	opts := &createOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "create --name <name>",
		Short: "Create a security domain",
		Example: `  gio am domain create --name "My Domain"
  gio am domain create --name "My Domain" --description "Production domain"`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Domain name (required)")
	cmd.Flags().StringVar(&opts.description, "description", "", "Domain description")
	cmd.Flags().StringVar(&opts.dataPlaneID, "data-plane-id", "default", "Data plane ID")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func (o *createOptions) run() error {
	f := o.factory

	body := map[string]any{
		"name":        o.name,
		"dataPlaneId": o.dataPlaneID,
	}
	if o.description != "" {
		body["description"] = o.description
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().CreateDomain(json.RawMessage(raw))
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

	return printDomainDetail(p, data)
}
