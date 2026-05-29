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
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

type setDomainOptions struct {
	factory  *factory.Factory
	domainID string
	clear    bool
}

func newSetCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set AM context values",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newSetDomainCmd(f))
	return cmd
}

func newSetDomainCmd(f *factory.Factory) *cobra.Command {
	opts := &setDomainOptions{factory: f}
	cmd := &cobra.Command{
		Use:     "domain <id>",
		Short:   "Set active AM domain",
		Example: "  gctl am set domain my-domain-id\n  gctl am set domain --clear",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			if len(args) == 1 {
				opts.domainID = args[0]
			}
			if !opts.clear && opts.domainID == "" {
				return fmt.Errorf("provide a domain ID or use --clear")
			}
			return opts.run()
		},
	}
	cmd.Flags().BoolVar(&opts.clear, "clear", false, "Unset current domain")
	return cmd
}

func (o *setDomainOptions) run() error {
	cfg := o.factory.Config
	contextName := cfg.Current
	ctx := cfg.EnsureContext(contextName)

	if o.clear {
		ctx.Domain = ""
	} else {
		ctx.Domain = o.domainID
	}

	if err := cfg.SaveTo(o.factory.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	if o.clear {
		fmt.Fprintf(o.factory.IOStreams.Out, "Domain cleared for context '%s'.\n", contextName)
	} else {
		fmt.Fprintf(o.factory.IOStreams.Out, "Domain set to '%s' for context '%s'.\n", o.domainID, contextName)
	}
	return nil
}
