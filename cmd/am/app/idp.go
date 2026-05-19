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

func newAppIdpCmd(f *factory.Factory, domainID *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "idp",
		Short: "Manage identity providers bound to an application",
	}

	cmd.AddCommand(newAppIdpListCmd(f, domainID))
	cmd.AddCommand(newAppIdpAddCmd(f, domainID))
	cmd.AddCommand(newAppIdpRemoveCmd(f, domainID))

	return cmd
}

func newAppIdpListCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list <appID>",
		Short:   "List identity providers bound to an application",
		Example: `  gio am app idp list my-app --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			bindings, err := readAppIdpBindings(f, *domainID, args[0])
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(bindings)
			}

			return printIdpBindings(p, bindings)
		},
	}
}

func newAppIdpAddCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var (
		priority      int
		selectionRule string
	)

	cmd := &cobra.Command{
		Use:   "add <appID> <idpID>",
		Short: "Bind an identity provider to an application",
		Example: `  gio am app idp add my-app my-idp --domain my-domain
  gio am app idp add my-app my-idp --domain my-domain --priority 10`,
		Args: cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}
			return runAppIdpAdd(f, *domainID, args[0], args[1], priority, selectionRule)
		},
	}

	cmd.Flags().IntVar(&priority, "priority", 0, "Priority of this IdP binding (lower runs first)")
	cmd.Flags().StringVar(&selectionRule, "selection-rule", "", "EL expression used to select this IdP at runtime")

	return cmd
}

func newAppIdpRemoveCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "remove <appID> <idpID>",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove an identity provider binding from an application",
		Example: `  gio am app idp remove my-app my-idp --domain my-domain`,
		Args:    cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}
			return runAppIdpRemove(f, *domainID, args[0], args[1])
		},
	}
}

// readAppIdpBindings returns the application's `identityProviders` array as
// a slice of maps. An app with no bindings returns an empty slice (not nil).
func readAppIdpBindings(f *factory.Factory, domainID, appID string) ([]map[string]any, error) {
	data, err := f.AM().GetApplication(domainID, appID)
	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse application: %w", err)
	}

	raw, _ := m["identityProviders"].([]any)
	out := make([]map[string]any, 0, len(raw))
	for i, item := range raw {
		b, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("identityProviders entry %d is not an object", i)
		}
		out = append(out, b)
	}
	return out, nil
}

func runAppIdpAdd(f *factory.Factory, domainID, appID, idpID string, priority int, selectionRule string) error {
	bindings, err := readAppIdpBindings(f, domainID, appID)
	if err != nil {
		return err
	}

	for _, b := range bindings {
		if id, _ := b["identity"].(string); id == idpID {
			// Update existing binding in place.
			b["priority"] = priority
			if selectionRule == "" {
				b["selectionRule"] = nil
			} else {
				b["selectionRule"] = selectionRule
			}
			return patchAppIdpBindings(f, domainID, appID, bindings)
		}
	}

	binding := map[string]any{
		"identity":      idpID,
		"priority":      priority,
		"selectionRule": nil,
	}
	if selectionRule != "" {
		binding["selectionRule"] = selectionRule
	}

	bindings = append(bindings, binding)

	return patchAppIdpBindings(f, domainID, appID, bindings)
}

func runAppIdpRemove(f *factory.Factory, domainID, appID, idpID string) error {
	bindings, err := readAppIdpBindings(f, domainID, appID)
	if err != nil {
		return err
	}

	filtered := bindings[:0]
	removed := false
	for _, b := range bindings {
		if id, _ := b["identity"].(string); id == idpID {
			removed = true
			continue
		}
		filtered = append(filtered, b)
	}

	if !removed {
		return fmt.Errorf("identity provider %q is not bound to application %q", idpID, appID)
	}

	return patchAppIdpBindings(f, domainID, appID, filtered)
}

func patchAppIdpBindings(f *factory.Factory, domainID, appID string, bindings []map[string]any) error {
	body, _ := json.Marshal(map[string]any{"identityProviders": bindings})

	data, err := f.AM().PatchApplication(domainID, appID, json.RawMessage(body))
	if err != nil {
		return err
	}

	updated, err := readBindingsFromApp(data)
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

	return printIdpBindings(p, updated)
}

func readBindingsFromApp(data []byte) ([]map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse application: %w", err)
	}
	raw, _ := m["identityProviders"].([]any)
	out := make([]map[string]any, 0, len(raw))
	for i, item := range raw {
		b, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("identityProviders entry %d is not an object", i)
		}
		out = append(out, b)
	}
	return out, nil
}

func printIdpBindings(p *printer.Printer, bindings []map[string]any) error {
	if len(bindings) == 0 {
		p.PrintMessage("No identity providers bound.")
		return nil
	}

	p.PrintMessage("%-40s%-10s%s", "IDENTITY", "PRIORITY", "SELECTION RULE")
	for _, b := range bindings {
		id, _ := b["identity"].(string)
		priority := 0
		if v, ok := b["priority"].(float64); ok {
			priority = int(v)
		}
		rule, _ := b["selectionRule"].(string)
		p.PrintMessage("%-40s%-10d%s", id, priority, rule)
	}
	return nil
}
