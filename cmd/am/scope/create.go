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

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

type createOptions struct {
	factory     *factory.Factory
	domainID    *string
	key         string
	name        string
	description string
}

func newCreateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &createOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "create --key <key> --name <name>",
		Short: "Create a scope",
		Example: `  gctl am scope create --domain my-domain --key openid --name "OpenID"
  gctl am scope create --domain my-domain --key profile --name "Profile" --description "Access to profile"`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.key, "key", "", "Scope key (required)")
	cmd.Flags().StringVar(&opts.name, "name", "", "Scope name (required)")
	cmd.Flags().StringVar(&opts.description, "description", "", "Scope description")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func (o *createOptions) run() error {
	f := o.factory

	body := map[string]any{
		"key":  o.key,
		"name": o.name,
	}

	if o.description != "" {
		body["description"] = o.description
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().CreateScope(*o.domainID, json.RawMessage(raw))
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
