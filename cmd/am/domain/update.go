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

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

type updateOptions struct {
	factory         *factory.Factory
	domainID        string
	name            string
	description     string
	allowLocalhost  bool
	allowHTTPScheme bool
}

func newUpdateCmd(f *factory.Factory) *cobra.Command {
	opts := &updateOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "update <domainID>",
		Short: "Update a security domain",
		Example: `  gctl am domain update my-domain-id --name "New Name"
  gctl am domain update my-domain-id --description "Updated description"
  gctl am domain update my-domain-id --allow-localhost-redirect --allow-http-redirect`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.domainID = args[0]

			return opts.run(cmd)
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Domain name")
	cmd.Flags().StringVar(&opts.description, "description", "", "Domain description")
	cmd.Flags().BoolVar(&opts.allowLocalhost, "allow-localhost-redirect", false, "Allow loopback (localhost / 127.0.0.1) redirect URIs on registered clients")
	cmd.Flags().BoolVar(&opts.allowHTTPScheme, "allow-http-redirect", false, "Allow http:// scheme redirect URIs on registered clients")

	return cmd
}

func (o *updateOptions) run(cmd *cobra.Command) error {
	f := o.factory

	body := map[string]any{}
	if o.name != "" {
		body["name"] = o.name
	}

	if o.description != "" {
		body["description"] = o.description
	}

	allowLocalhostSet := cmd.Flags().Changed("allow-localhost-redirect")
	allowHTTPSet := cmd.Flags().Changed("allow-http-redirect")

	if allowLocalhostSet || allowHTTPSet {
		oidc, err := mergeClientRegistrationRedirectFlags(f, o.domainID, o.allowLocalhost, allowLocalhostSet, o.allowHTTPScheme, allowHTTPSet)
		if err != nil {
			return err
		}
		body["oidc"] = oidc
	}

	if len(body) == 0 {
		return fmt.Errorf("at least one flag (--name, --description, --allow-localhost-redirect, --allow-http-redirect) is required")
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().PatchDomain(o.domainID, json.RawMessage(raw))
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

// mergeClientRegistrationRedirectFlags performs a read-modify-write of the
// domain's `oidc` block so we can flip the two redirect-URI flags on
// `clientRegistrationSettings` without clobbering any other OIDC config
// (CIMD, etc).
func mergeClientRegistrationRedirectFlags(f *factory.Factory, domainID string, allowLocalhost, allowLocalhostSet, allowHTTP, allowHTTPSet bool) (map[string]any, error) {
	current, err := f.AM().GetDomain(domainID)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(current, &m); err != nil {
		return nil, fmt.Errorf("failed to parse domain: %w", err)
	}

	oidc, _ := m["oidc"].(map[string]any)
	if oidc == nil {
		oidc = map[string]any{}
	}

	settings, _ := oidc["clientRegistrationSettings"].(map[string]any)
	if settings == nil {
		settings = map[string]any{}
	}

	if allowLocalhostSet {
		settings["allowLocalhostRedirectUri"] = allowLocalhost
	}
	if allowHTTPSet {
		settings["allowHttpSchemeRedirectUri"] = allowHTTP
	}

	oidc["clientRegistrationSettings"] = settings

	return oidc, nil
}
