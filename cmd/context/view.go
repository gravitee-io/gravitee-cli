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

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/config"
	"gravitee.io/gctl/internal/factory"
)

func newViewCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Display full details of the current context",
		Example: `  gctl context view
  gctl context view --context prod`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runView(f)
		},
	}
}

func runView(f *factory.Factory) error {
	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	cfg := f.Config

	contextName := cfg.Current
	if f.Overrides.Context != "" {
		contextName = f.Overrides.Context
	}

	if contextName == "" {
		return fmt.Errorf("no context configured\nHint: run 'gctl login' to get started")
	}

	ctx, ok := cfg.Contexts[contextName]
	if !ok {
		return fmt.Errorf("context '%s' not found in config", contextName)
	}

	org := ctx.Org
	if org == "" {
		org = config.DefaultOrg
	}

	env := ctx.Env
	if env == "" {
		env = config.DefaultEnv
	}

	out := f.IOStreams.Out
	fmt.Fprintf(out, "Context:    %s\n", contextName)
	fmt.Fprintf(out, "Org:        %s\n", org)
	fmt.Fprintf(out, "Env:        %s\n", env)

	if ctx.APIM != nil {
		fmt.Fprintf(out, "APIM URL:   %s\n", ctx.APIM.URL)
		fmt.Fprintf(out, "APIM Token: %s\n", cmdutil.MaskToken(ctx.APIM.Token))
	}

	if ctx.AM != nil {
		fmt.Fprintf(out, "AM URL:     %s\n", ctx.AM.URL)
		fmt.Fprintf(out, "AM Token:   %s\n", cmdutil.MaskToken(ctx.AM.Token))
	}

	return nil
}
