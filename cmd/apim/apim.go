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

	analyticscmd "gravitee.io/gctl/cmd/apim/analytics"
	apicmd "gravitee.io/gctl/cmd/apim/api"
	apikeycmd "gravitee.io/gctl/cmd/apim/apikey"
	appcmd "gravitee.io/gctl/cmd/apim/application"
	envcmd "gravitee.io/gctl/cmd/apim/environment"
	healthcmd "gravitee.io/gctl/cmd/apim/health"
	logcmd "gravitee.io/gctl/cmd/apim/log"
	membercmd "gravitee.io/gctl/cmd/apim/member"
	metadatacmd "gravitee.io/gctl/cmd/apim/metadata"
	pagecmd "gravitee.io/gctl/cmd/apim/page"
	plancmd "gravitee.io/gctl/cmd/apim/plan"
	plugincmd "gravitee.io/gctl/cmd/apim/plugin"
	subcmd "gravitee.io/gctl/cmd/apim/subscription"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newAPIMBaseCmd(f *factory.Factory) *cobra.Command {
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

	cmd.AddCommand(logcmd.NewLogCmd(f))
	cmd.AddCommand(analyticscmd.NewAnalyticsCmd(f))
	cmd.AddCommand(healthcmd.NewHealthCmd(f))
	cmd.AddCommand(metadatacmd.NewMetadataCmd(f))
	cmd.AddCommand(envcmd.NewEnvironmentCmd(f))
	cmd.AddCommand(plugincmd.NewPluginCmd(f))

	return cmd
}

// NewAPIMCmdRO creates the apim command with read-only subcommands only.
func NewAPIMCmdRO(f *factory.Factory) *cobra.Command {
	cmd := newAPIMBaseCmd(f)

	cmd.AddCommand(apicmd.NewAPICmdRO(f))
	cmd.AddCommand(plancmd.NewPlanCmdRO(f))
	cmd.AddCommand(subcmd.NewSubscriptionCmdRO(f))
	cmd.AddCommand(apikeycmd.NewAPIKeyCmdRO(f))
	cmd.AddCommand(membercmd.NewMemberCmdRO(f))
	cmd.AddCommand(pagecmd.NewPageCmdRO(f))
	cmd.AddCommand(appcmd.NewApplicationCmdRO(f))

	return cmd
}

// NewAPIMCmd creates the apim parent command with all APIM subcommands.
func NewAPIMCmd(f *factory.Factory) *cobra.Command {
	cmd := newAPIMBaseCmd(f)

	cmd.AddCommand(apicmd.NewAPICmd(f))
	cmd.AddCommand(plancmd.NewPlanCmd(f))
	cmd.AddCommand(subcmd.NewSubscriptionCmd(f))
	cmd.AddCommand(apikeycmd.NewAPIKeyCmd(f))
	cmd.AddCommand(membercmd.NewMemberCmd(f))
	cmd.AddCommand(pagecmd.NewPageCmd(f))
	cmd.AddCommand(appcmd.NewApplicationCmd(f))

	return cmd
}
