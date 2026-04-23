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

package context

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newUseCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Switch to a different context",
		Example: `  gio context use prod
  gio context use staging`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runUse(f, args[0])
		},
	}
}

func runUse(f *factory.Factory, name string) error {
	name = config.NormalizeContextName(name)

	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	cfg := f.Config

	if _, ok := cfg.Contexts[name]; !ok {
		available := cfg.ContextNames()
		if len(available) == 0 {
			return fmt.Errorf("context '%s' not found\nHint: no contexts configured, run 'gio login' to get started", name)
		}

		return fmt.Errorf("context '%s' not found\nHint: available contexts: %s", name, strings.Join(available, ", "))
	}

	cfg.Current = name

	if err := cfg.SaveTo(f.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(f.IOStreams.Out, "Switched to context '%s'.\n", name)

	return nil
}
