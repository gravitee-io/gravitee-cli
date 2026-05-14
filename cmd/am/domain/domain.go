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
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewDomainCmdRO creates the domain command with read-only subcommands.
func NewDomainCmdRO(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain",
		Aliases: []string{"dom"},
		Short:   "Manage security domains",
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newEntrypointsCmdRO(f))
	cmd.AddCommand(newCIMDCmdRO(f))

	return cmd
}

// NewDomainCmd creates the domain parent command with all domain subcommands.
func NewDomainCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain",
		Aliases: []string{"dom"},
		Short:   "Manage security domains",
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(newDeleteCmd(f))
	cmd.AddCommand(newEnableCmd(f))
	cmd.AddCommand(newDisableCmd(f))
	cmd.AddCommand(newUpdateCertSettingsCmd(f))
	cmd.AddCommand(newEntrypointsCmd(f))
	cmd.AddCommand(newCIMDCmd(f))

	return cmd
}
