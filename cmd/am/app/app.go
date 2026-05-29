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
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

// NewAppCmdRO creates the application command with read-only subcommands.
func NewAppCmdRO(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "app",
		Aliases: []string{"application"},
		Short:   "Manage applications",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))

	return cmd
}

// NewAppCmd creates the application parent command with all application subcommands.
func NewAppCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "app",
		Aliases: []string{"application"},
		Short:   "Manage applications",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))
	cmd.AddCommand(newCreateCmd(f, &domainID))
	cmd.AddCommand(newUpdateCmd(f, &domainID))
	cmd.AddCommand(newDeleteCmd(f, &domainID))
	cmd.AddCommand(newSecretCmd(f, &domainID))
	cmd.AddCommand(newAppMemberCmd(f, &domainID))
	cmd.AddCommand(newAppFlowCmd(f, &domainID))
	cmd.AddCommand(newAppEmailCmd(f, &domainID))
	cmd.AddCommand(newAppFormCmd(f, &domainID))
	cmd.AddCommand(newAppResourceCmd(f, &domainID))
	cmd.AddCommand(newAppAnalyticsCmd(f, &domainID))
	cmd.AddCommand(newChangeTypeCmd(f, &domainID))
	cmd.AddCommand(newAppResourcePolicyCmd(f, &domainID))
	cmd.AddCommand(newAppIdpCmd(f, &domainID))

	return cmd
}
