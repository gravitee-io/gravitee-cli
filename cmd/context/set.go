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

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

type setOptions struct {
	factory *factory.Factory
	org     string
	envID   string
}

func newSetCmd(f *factory.Factory) *cobra.Command {
	opts := &setOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "set <name>",
		Short: "Create or update shared fields of a context",
		Long:  "Create or update shared fields (org, env) of a context. Use 'gio login' to set product URLs and tokens.",
		Example: `  gio context set prod --org ACME --env production
  gio context set dev --org DEV --env development`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.run(cmd, args[0])
		},
	}

	cmd.Flags().StringVar(&opts.org, "org", "", "Organization ID")
	cmd.Flags().StringVar(&opts.envID, "env", "", "Environment ID")

	return cmd
}

func (o *setOptions) run(cmd *cobra.Command, name string) error {
	name = config.NormalizeContextName(name)

	if err := cmdutil.SetupConfig(o.factory); err != nil {
		return err
	}

	cfg := o.factory.Config
	ctx := cfg.EnsureContext(name)

	if cmd.Flags().Changed("org") {
		ctx.Org = o.org
	}

	if cmd.Flags().Changed("env") {
		ctx.Env = o.envID
	}

	if err := cfg.SaveTo(o.factory.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(o.factory.IOStreams.Out, "Context '%s' updated.\n", name)

	return nil
}
