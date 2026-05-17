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
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type createOptions struct {
	factory      *factory.Factory
	domainID     *string
	name         string
	appType      string
	description  string
	redirectURIs string
}

func newCreateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &createOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "create --name <name> --type <type>",
		Short: "Create an application",
		Example: `  gio am app create --domain my-domain --name "My App" --type web
  gio am app create --domain my-domain --name "My App" --type browser --redirect-uris "http://localhost:4200/callback"`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.ValidateEnum(opts.appType, "type", []string{"web", "native", "browser", "service", "resource_server"}); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Application name (required)")
	cmd.Flags().StringVar(&opts.appType, "type", "", "Application type: web, native, browser, service, resource_server (required)")
	cmd.Flags().StringVar(&opts.description, "description", "", "Application description")
	cmd.Flags().StringVar(&opts.redirectURIs, "redirect-uris", "", "Comma-separated redirect URIs")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("type")

	return cmd
}

func (o *createOptions) run() error {
	f := o.factory

	body := map[string]any{
		"name": o.name,
		"type": o.appType,
	}

	if o.description != "" {
		body["description"] = o.description
	}

	if o.redirectURIs != "" {
		uris := strings.Split(o.redirectURIs, ",")
		for i := range uris {
			uris[i] = strings.TrimSpace(uris[i])
		}

		body["redirectUris"] = uris
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().CreateApplication(*o.domainID, json.RawMessage(raw))
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

	if err := printAppDetail(p, data); err != nil {
		return err
	}

	printInitialClientSecret(p, data)

	return nil
}

// printInitialClientSecret surfaces the OAuth2 client secret that AM
// auto-generates for service apps on create. The secret is only present in
// the create response — subsequent reads omit it, so we make sure the user
// sees it the one time they can.
func printInitialClientSecret(p *printer.Printer, data []byte) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}

	settings, _ := m["settings"].(map[string]any)
	if settings == nil {
		return
	}
	oauth, _ := settings["oauth"].(map[string]any)
	if oauth == nil {
		return
	}

	clientID, _ := oauth["clientId"].(string)
	secret, _ := oauth["clientSecret"].(string)
	if secret == "" {
		return
	}

	p.PrintMessage("")
	if clientID != "" {
		p.PrintMessage("Client ID:      %s", clientID)
	}
	p.PrintMessage("Client secret (store it now — it will not be shown again):")
	p.PrintMessage("  %s", secret)
}
