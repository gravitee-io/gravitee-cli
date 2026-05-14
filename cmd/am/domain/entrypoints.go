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

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newEntrypointsCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "entrypoints",
		Aliases: []string{"entrypoint"},
		Short:   "Manage a domain's entrypoint configuration (context-path or vhosts)",
	}

	cmd.AddCommand(newEntrypointsGetCmd(f))
	cmd.AddCommand(newEntrypointsSetPathCmd(f))
	cmd.AddCommand(newEntrypointsAddVhostCmd(f))
	cmd.AddCommand(newEntrypointsRemoveVhostCmd(f))
	cmd.AddCommand(newEntrypointsClearVhostsCmd(f))

	return cmd
}

// newEntrypointsCmdRO is the read-only entrypoints parent (get only).
func newEntrypointsCmdRO(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "entrypoints",
		Aliases: []string{"entrypoint"},
		Short:   "Inspect a domain's entrypoint configuration (context-path or vhosts)",
	}

	cmd.AddCommand(newEntrypointsGetCmd(f))

	return cmd
}

func newEntrypointsGetCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "get <domainID>",
		Short:   "Show the entrypoint configuration for a domain",
		Example: `  gio am domain entrypoints get my-domain-id -o json`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetDomain(args[0])
			if err != nil {
				return err
			}

			view, err := extractEntrypointView(data)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(view)
			}

			return printEntrypointsTable(p, view)
		},
	}
}

func newEntrypointsSetPathCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "set-path <domainID> <path>",
		Short:   "Switch the domain to context-path mode and set its path",
		Example: `  gio am domain entrypoints set-path my-domain-id /auth`,
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, _ := json.Marshal(map[string]any{
				"vhostMode": false,
				"path":      args[1],
			})

			return patchAndPrintEntrypoints(f, args[0], body)
		},
	}
}

func newEntrypointsAddVhostCmd(f *factory.Factory) *cobra.Command {
	var (
		path     string
		override bool
	)

	cmd := &cobra.Command{
		Use:   "add-vhost <domainID> <host>",
		Short: "Add a vhost to the domain (switches to vhost mode)",
		Example: `  gio am domain entrypoints add-vhost my-domain-id auth.example.com --path / --override
  gio am domain entrypoints add-vhost my-domain-id alt.example.com --path /auth`,
		Args: cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			domainID, host := args[0], args[1]

			current, err := f.AM().GetDomain(domainID)
			if err != nil {
				return err
			}

			view, err := extractEntrypointView(current)
			if err != nil {
				return err
			}

			for _, v := range view.Vhosts {
				if v.Host == host && v.Path == path {
					return fmt.Errorf("vhost %q with path %q already exists", host, path)
				}
			}

			view.Vhosts = append(view.Vhosts, vhost{Host: host, Path: path, OverrideEntrypoint: override})

			body, _ := json.Marshal(map[string]any{
				"vhostMode": true,
				"vhosts":    view.Vhosts,
			})

			return patchAndPrintEntrypoints(f, domainID, body)
		},
	}

	cmd.Flags().StringVar(&path, "path", "/", "Path served by this vhost")
	cmd.Flags().BoolVar(&override, "override", false, "Mark this vhost as the OIDC entrypoint override")

	return cmd
}

func newEntrypointsRemoveVhostCmd(f *factory.Factory) *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:     "remove-vhost <domainID> <host>",
		Short:   "Remove a vhost from the domain by host (and optionally path)",
		Example: `  gio am domain entrypoints remove-vhost my-domain-id auth.example.com`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			domainID, host := args[0], args[1]
			pathFilter := cmd.Flags().Changed("path")

			current, err := f.AM().GetDomain(domainID)
			if err != nil {
				return err
			}

			view, err := extractEntrypointView(current)
			if err != nil {
				return err
			}

			kept := make([]vhost, 0, len(view.Vhosts))
			removed := 0

			for _, v := range view.Vhosts {
				if v.Host == host && (!pathFilter || v.Path == path) {
					removed++
					continue
				}

				kept = append(kept, v)
			}

			if removed == 0 {
				return fmt.Errorf("no vhost matched host %q", host)
			}

			body, _ := json.Marshal(map[string]any{
				"vhostMode": len(kept) > 0,
				"vhosts":    kept,
			})

			return patchAndPrintEntrypoints(f, domainID, body)
		},
	}

	cmd.Flags().StringVar(&path, "path", "", "Only remove the vhost matching this path")

	return cmd
}

func newEntrypointsClearVhostsCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "clear-vhosts <domainID>",
		Short:   "Drop all vhosts and switch the domain back to context-path mode",
		Example: `  gio am domain entrypoints clear-vhosts my-domain-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, _ := json.Marshal(map[string]any{
				"vhostMode": false,
				"vhosts":    []vhost{},
			})

			return patchAndPrintEntrypoints(f, args[0], body)
		},
	}
}

type vhost struct {
	Host               string `json:"host"`
	Path               string `json:"path"`
	OverrideEntrypoint bool   `json:"overrideEntrypoint"`
}

type entrypointView struct {
	VhostMode bool    `json:"vhostMode"`
	Path      string  `json:"path,omitempty"`
	Vhosts    []vhost `json:"vhosts"`
}

func extractEntrypointView(data []byte) (entrypointView, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return entrypointView{}, fmt.Errorf("failed to parse domain: %w", err)
	}

	view := entrypointView{Vhosts: []vhost{}}

	if v, ok := m["vhostMode"].(bool); ok {
		view.VhostMode = v
	}

	if v, ok := m["path"].(string); ok {
		view.Path = v
	}

	raw, ok := m["vhosts"].([]any)
	if !ok {
		return view, nil
	}

	for _, item := range raw {
		entry, ok := item.(map[string]any)
		if !ok {
			continue
		}

		var v vhost

		if s, ok := entry["host"].(string); ok {
			v.Host = s
		}

		if s, ok := entry["path"].(string); ok {
			v.Path = s
		}

		if b, ok := entry["overrideEntrypoint"].(bool); ok {
			v.OverrideEntrypoint = b
		}

		view.Vhosts = append(view.Vhosts, v)
	}

	return view, nil
}

func patchAndPrintEntrypoints(f *factory.Factory, domainID string, body json.RawMessage) error {
	data, err := f.AM().PatchDomain(domainID, body)
	if err != nil {
		return err
	}

	view, err := extractEntrypointView(data)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(view)
	}

	return printEntrypointsTable(p, view)
}

func printEntrypointsTable(p *printer.Printer, view entrypointView) error {
	mode := "context-path"
	if view.VhostMode {
		mode = "vhost"
	}

	p.PrintMessage("%-16s%s", "Mode:", mode)

	if !view.VhostMode {
		path := view.Path
		if path == "" {
			path = "/"
		}

		p.PrintMessage("%-16s%s", "Path:", path)

		return nil
	}

	if len(view.Vhosts) == 0 {
		p.PrintMessage("%-16s(none)", "Vhosts:")

		return nil
	}

	p.PrintMessage("Vhosts:")

	for _, v := range view.Vhosts {
		override := ""
		if v.OverrideEntrypoint {
			override = "  (override)"
		}

		p.PrintMessage("  - %s%s%s", v.Host, v.Path, override)
	}

	return nil
}
