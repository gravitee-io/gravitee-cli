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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type updateOptions struct {
	factory     *factory.Factory
	domainID    *string
	appID       string
	name        string
	description string
	enabled     string
	template    string
}

func newUpdateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &updateOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "update <appID>",
		Short: "Update an application",
		Example: `  gio am app update my-app-id --domain my-domain --name "New Name"
  gio am app update my-app-id --domain my-domain --enabled false`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.appID = args[0]

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Application name")
	cmd.Flags().StringVar(&opts.description, "description", "", "Application description")
	cmd.Flags().StringVar(&opts.enabled, "enabled", "", "Enable or disable application (true/false)")
	cmd.Flags().StringVar(&opts.template, "template", "", "Mark application as a template (true/false)")

	return cmd
}

func (o *updateOptions) run() error {
	f := o.factory

	body := map[string]any{}
	if o.name != "" {
		body["name"] = o.name
	}

	if o.description != "" {
		body["description"] = o.description
	}

	if o.enabled != "" {
		switch o.enabled {
		case "true":
			body["enabled"] = true
		case "false":
			body["enabled"] = false
		default:
			return fmt.Errorf("--enabled must be 'true' or 'false', got %q", o.enabled)
		}
	}

	if o.template != "" {
		switch o.template {
		case "true":
			body["template"] = true
		case "false":
			body["template"] = false
		default:
			return fmt.Errorf("--template must be 'true' or 'false', got %q", o.template)
		}
	}

	if len(body) == 0 {
		return fmt.Errorf("at least one flag (--name, --description, --enabled, --template) is required")
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().PatchApplication(*o.domainID, o.appID, json.RawMessage(raw))
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

	return printAppDetail(p, data)
}
