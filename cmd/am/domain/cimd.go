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
	"strings"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newCIMDCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cimd",
		Short: "Manage Client-Initiated Identity Management (CIMD) settings for a domain",
	}

	cmd.AddCommand(newCIMDGetCmd(f))
	cmd.AddCommand(newCIMDEnableCmd(f))
	cmd.AddCommand(newCIMDDisableCmd(f))

	return cmd
}

// newCIMDCmdRO is the read-only CIMD parent (get only).
func newCIMDCmdRO(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cimd",
		Short: "Inspect Client-Initiated Identity Management (CIMD) settings for a domain",
	}

	cmd.AddCommand(newCIMDGetCmd(f))

	return cmd
}

func newCIMDGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <domainID>",
		Short:   "Show CIMD settings for a domain",
		Example: `  gctl am domain cimd get my-domain-id -o json`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetDomain(args[0])
			if err != nil {
				return err
			}

			settings, err := extractCIMDSettings(data)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(settings)
			}

			return printCIMDTable(p, settings)
		},
	}
}

type cimdOptions struct {
	allowPrivate      bool
	allowHTTP         bool
	templateID        string
	allowedDomainsCSV string
	fetchTimeoutMs    int
	maxResponseSizeKB int
	cacheTTLSeconds   int
}

func newCIMDEnableCmd(f *factory.Factory) *cobra.Command {
	opts := &cimdOptions{}

	cmd := &cobra.Command{
		Use:   "enable <domainID>",
		Short: "Enable CIMD on a domain and configure its settings",
		Example: `  gctl am domain cimd enable my-domain-id --template-id my-template-app
  gctl am domain cimd enable my-domain-id --allow-private --allow-http --allowed-domains "a.com,b.com"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runCIMDEnable(f, cmd, args[0], opts)
		},
	}

	cmd.Flags().BoolVar(&opts.allowPrivate, "allow-private", false, "Allow private IP addresses for client metadata URIs")
	cmd.Flags().BoolVar(&opts.allowHTTP, "allow-http", false, "Allow unsecured HTTP client metadata URIs")
	cmd.Flags().StringVar(&opts.templateID, "template-id", "", "ID of the template application to bind")
	cmd.Flags().StringVar(&opts.allowedDomainsCSV, "allowed-domains", "", "Comma-separated list of allowed client metadata domains")
	cmd.Flags().IntVar(&opts.fetchTimeoutMs, "fetch-timeout-ms", 0, "Fetch timeout for client metadata in milliseconds")
	cmd.Flags().IntVar(&opts.maxResponseSizeKB, "max-response-size-kb", 0, "Maximum client metadata response size in KB")
	cmd.Flags().IntVar(&opts.cacheTTLSeconds, "cache-ttl-seconds", 0, "Cache TTL for client metadata in seconds")

	return cmd
}

func newCIMDDisableCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "disable <domainID>",
		Short:   "Disable CIMD on a domain",
		Example: `  gctl am domain cimd disable my-domain-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runCIMDPatch(f, args[0], map[string]any{"enabled": false})
		},
	}
}

func runCIMDEnable(f *factory.Factory, cmd *cobra.Command, domainID string, opts *cimdOptions) error {
	settings := map[string]any{"enabled": true}

	if cmd.Flags().Changed("allow-private") {
		settings["allowPrivateIpAddress"] = opts.allowPrivate
	}

	if cmd.Flags().Changed("allow-http") {
		settings["allowUnsecuredHttpUri"] = opts.allowHTTP
	}

	if opts.templateID != "" {
		settings["templateId"] = opts.templateID
	}

	if opts.allowedDomainsCSV != "" {
		settings["allowedDomains"] = splitCSV(opts.allowedDomainsCSV)
	}

	if cmd.Flags().Changed("fetch-timeout-ms") {
		settings["fetchTimeoutMs"] = opts.fetchTimeoutMs
	}

	if cmd.Flags().Changed("max-response-size-kb") {
		settings["maxResponseSizeKb"] = opts.maxResponseSizeKB
	}

	if cmd.Flags().Changed("cache-ttl-seconds") {
		settings["cacheTtlSeconds"] = opts.cacheTTLSeconds
	}

	return runCIMDPatch(f, domainID, settings)
}

// runCIMDPatch merges the supplied cimdSettings into the existing `oidc` block
// (read-modify-write) so the PATCH does not clobber unrelated OIDC fields.
func runCIMDPatch(f *factory.Factory, domainID string, settings map[string]any) error {
	current, err := f.AM().GetDomain(domainID)
	if err != nil {
		return err
	}

	var m map[string]any
	if err = json.Unmarshal(current, &m); err != nil {
		return fmt.Errorf("failed to parse domain: %w", err)
	}

	oidc, _ := m["oidc"].(map[string]any)
	if oidc == nil {
		oidc = map[string]any{}
	}

	existing, _ := oidc["cimdSettings"].(map[string]any)
	if existing == nil {
		existing = map[string]any{}
	}

	for k, v := range settings {
		existing[k] = v
	}

	oidc["cimdSettings"] = existing

	body, _ := json.Marshal(map[string]any{"oidc": oidc})

	data, err := f.AM().PatchDomain(domainID, body)
	if err != nil {
		return err
	}

	updated, err := extractCIMDSettings(data)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(updated)
	}

	return printCIMDTable(p, updated)
}

func extractCIMDSettings(data []byte) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse domain: %w", err)
	}

	oidc, ok := m["oidc"].(map[string]any)
	if !ok {
		return map[string]any{"enabled": false}, nil
	}

	settings, ok := oidc["cimdSettings"].(map[string]any)
	if !ok {
		return map[string]any{"enabled": false}, nil
	}

	return settings, nil
}

func printCIMDTable(p *printer.Printer, settings map[string]any) error {
	for _, field := range []struct{ label, key string }{
		{"Enabled", "enabled"},
		{"Template ID", "templateId"},
		{"Allowed domains", "allowedDomains"},
		{"Allow private IP", "allowPrivateIpAddress"},
		{"Allow HTTP", "allowUnsecuredHttpUri"},
		{"Fetch timeout (ms)", "fetchTimeoutMs"},
		{"Max response (KB)", "maxResponseSizeKb"},
		{"Cache TTL (s)", "cacheTtlSeconds"},
	} {
		v, ok := settings[field.key]
		if !ok || v == nil {
			continue
		}

		p.PrintMessage("%-22s%v", field.label+":", v)
	}

	return nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}
