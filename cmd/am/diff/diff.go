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

package diff

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/config"
	"gravitee.io/gctl/internal/factory"
)

type resourceSpec struct {
	name          string
	path          string
	keyField      string
	compareFields []string
	paginated     bool
}

var resourceSpecs = []resourceSpec{
	{"scopes", "scopes", "key", []string{"name", "description"}, true},
	{"roles", "roles", "name", []string{"description", "assignableType"}, true},
	{"groups", "groups", "name", []string{"description"}, true},
	{"applications", "applications", "name", []string{"description", "type"}, true},
	{"identities", "identities", "name", []string{"type"}, false},
	{"certificates", "certificates", "name", []string{"type"}, false},
	{"factors", "factors", "name", []string{"factorType"}, false},
	{"flows", "flows", "type", []string{"enabled"}, false},
}

func NewDiffCmd(f *factory.Factory) *cobra.Command {
	var fromCtx, toCtx string
	var fromDomain, toDomain string

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Compare domain configuration between two contexts",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runDiff(f, fromCtx, toCtx, fromDomain, toDomain)
		},
	}
	cmd.Flags().StringVar(&fromCtx, "from", "", "Source context name (required)")
	cmd.Flags().StringVar(&toCtx, "to", "", "Target context name (required)")
	cmd.Flags().StringVar(&fromDomain, "from-domain", "", "Override domain ID for source context")
	cmd.Flags().StringVar(&toDomain, "to-domain", "", "Override domain ID for target context")
	return cmd
}

func runDiff(f *factory.Factory, fromCtx, toCtx, fromDomain, toDomain string) error {
	if fromCtx == "" || toCtx == "" {
		return fmt.Errorf("--from and --to are required")
	}

	fromResolved, toResolved, err := resolveContexts(f.Config, fromCtx, toCtx, fromDomain, toDomain)
	if err != nil {
		return err
	}

	fromClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: fromResolved.URL, Token: fromResolved.Token})
	toClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: toResolved.URL, Token: toResolved.Token})

	fmt.Fprintf(f.IOStreams.Out, "Comparing %s/%s -> %s/%s\n\n",
		fromCtx, fromResolved.Domain, toCtx, toResolved.Domain)

	for _, spec := range resourceSpecs {
		fromItems, err := fetchItems(fromClient, fromResolved, spec.path, spec.paginated)
		if err != nil {
			fmt.Fprintf(f.IOStreams.Out, "  [%s] error fetching from: %v\n", spec.name, err)
			continue
		}
		toItems, err := fetchItems(toClient, toResolved, spec.path, spec.paginated)
		if err != nil {
			fmt.Fprintf(f.IOStreams.Out, "  [%s] error fetching to: %v\n", spec.name, err)
			continue
		}
		result := compareResources(fromItems, toItems, spec.keyField, spec.compareFields)
		if result.Added+result.Removed+result.Changed == 0 {
			fmt.Fprintf(f.IOStreams.Out, "  [%s] no differences\n", spec.name)
			continue
		}
		fmt.Fprintf(f.IOStreams.Out, "  [%s] +%d -%d ~%d\n", spec.name, result.Added, result.Removed, result.Changed)
		for _, line := range result.Lines {
			fmt.Fprintf(f.IOStreams.Out, "    %s\n", line)
		}
	}
	return nil
}

func resolveContexts(cfg *config.Config, fromCtx, toCtx, fromDomain, toDomain string) (*config.ResolvedContext, *config.ResolvedContext, error) {
	fromResolved, err := cfg.Resolve(config.Overrides{Context: fromCtx, Domain: fromDomain}, "am")
	if err != nil {
		return nil, nil, fmt.Errorf("--from: %w", err)
	}
	toResolved, err := cfg.Resolve(config.Overrides{Context: toCtx, Domain: toDomain}, "am")
	if err != nil {
		return nil, nil, fmt.Errorf("--to: %w", err)
	}
	if fromResolved.Domain == "" || toResolved.Domain == "" {
		return nil, nil, fmt.Errorf("both contexts must have a domain set (use --from-domain / --to-domain to override)")
	}
	return fromResolved, toResolved, nil
}

func fetchItems(c client.GraviteeClient, r *config.ResolvedContext, path string, paginated bool) ([]map[string]interface{}, error) {
	fullPath := fmt.Sprintf("/management/organizations/%s/environments/%s/domains/%s/%s",
		r.Org, r.Env, r.Domain, path)
	if paginated {
		fullPath += "?" + url.Values{"page": {"0"}, "size": {"1000"}}.Encode()
	}
	data, err := c.Get(fullPath)
	if err != nil {
		return nil, err
	}
	if paginated {
		var resp struct {
			Data []map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, err
		}
		return resp.Data, nil
	}
	var items []map[string]interface{}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}
