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

package am

import (
	"fmt"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/factory"
)

func newStatusCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current CLI context and session status",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := f.Config
			out := f.IOStreams.Out

			if cfg == nil || cfg.Current == "" {
				fmt.Fprintln(out, "context:       (not set)")
				fmt.Fprintln(out, "authenticated: no")

				return nil
			}

			ctx, ok := cfg.Contexts[cfg.Current]
			domain := ""

			if f.Resolved != nil {
				domain = f.Resolved.Domain
			}

			amURL := ""
			if ok && ctx != nil && ctx.AM != nil {
				amURL = ctx.AM.URL
			}

			fmt.Fprintf(out, "context:       %s @ %s\n", cfg.Current, amURL)

			if ok && ctx != nil {
				fmt.Fprintf(out, "organization:  %s\n", ctx.Org)
				fmt.Fprintf(out, "environment:   %s\n", ctx.Env)
			}

			if domain != "" {
				fmt.Fprintf(out, "domain:        %s\n", domain)
			} else {
				fmt.Fprintln(out, "domain:        (not set)")
			}

			if !ok || ctx == nil || ctx.AM == nil || ctx.AM.Token == "" {
				fmt.Fprintln(out, "authenticated: no")
			} else {
				fmt.Fprintln(out, "authenticated: yes")
			}

			return nil
		},
	}
}
