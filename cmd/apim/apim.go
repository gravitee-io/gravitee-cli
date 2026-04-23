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

package apim

import (
	"fmt"

	"github.com/spf13/cobra"

	apicmd "github.com/gravitee-io/gio-cli/cmd/apim/api"
	apikeycmd "github.com/gravitee-io/gio-cli/cmd/apim/apikey"
	appcmd "github.com/gravitee-io/gio-cli/cmd/apim/application"
	envcmd "github.com/gravitee-io/gio-cli/cmd/apim/environment"
	membercmd "github.com/gravitee-io/gio-cli/cmd/apim/member"
	metadatacmd "github.com/gravitee-io/gio-cli/cmd/apim/metadata"
	pagecmd "github.com/gravitee-io/gio-cli/cmd/apim/page"
	plancmd "github.com/gravitee-io/gio-cli/cmd/apim/plan"
	plugincmd "github.com/gravitee-io/gio-cli/cmd/apim/plugin"
	subcmd "github.com/gravitee-io/gio-cli/cmd/apim/subscription"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAPIMCmd creates the apim parent command with all APIM subcommands.
func NewAPIMCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apim",
		Short: "Gravitee API Management",
		Long:  "Manage Gravitee APIM resources: APIs, plans, subscriptions, applications, and more.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmdutil.SetupConfig(f); err != nil {
				return err
			}

			if err := cmdutil.ResolveProductContext(f, "apim"); err != nil {
				return err
			}

			return cmdutil.ResolveAPIMFlags(f, cmd)
		},
	}

	defaultHelp := cmd.HelpFunc()
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		// Silence errors: help must render even without a configured context.
		_ = cmdutil.SetupConfig(f)
		_ = cmdutil.ResolveProductContext(f, "apim")
		if header := cmdutil.ContextHeader(f, "apim"); header != "" {
			fmt.Fprint(c.OutOrStdout(), header+"\n")
		}

		defaultHelp(c, args)
	})

	cmd.AddCommand(apicmd.NewAPICmd(f))
	cmd.AddCommand(plancmd.NewPlanCmd(f))
	cmd.AddCommand(subcmd.NewSubscriptionCmd(f))
	cmd.AddCommand(apikeycmd.NewAPIKeyCmd(f))
	cmd.AddCommand(membercmd.NewMemberCmd(f))
	cmd.AddCommand(pagecmd.NewPageCmd(f))
	cmd.AddCommand(metadatacmd.NewMetadataCmd(f))
	cmd.AddCommand(appcmd.NewApplicationCmd(f))
	cmd.AddCommand(envcmd.NewEnvironmentCmd(f))
	cmd.AddCommand(plugincmd.NewPluginCmd(f))

	return cmd
}
