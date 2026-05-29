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

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newEnableCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "enable <domainID>",
		Short:   "Enable a security domain",
		Example: `  gctl am domain enable my-domain-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runSetEnabled(f, args[0], true)
		},
	}
}

func newDisableCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "disable <domainID>",
		Short:   "Disable a security domain",
		Example: `  gctl am domain disable my-domain-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runSetEnabled(f, args[0], false)
		},
	}
}

func runSetEnabled(f *factory.Factory, domainID string, enabled bool) error {
	body := map[string]any{"enabled": enabled}
	raw, _ := json.Marshal(body)

	if _, err := f.AM().PatchDomain(domainID, json.RawMessage(raw)); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	action := "enabled"
	if !enabled {
		action = "disabled"
	}

	p.PrintMessage("Domain '%s' %s.", domainID, action)

	return nil
}
