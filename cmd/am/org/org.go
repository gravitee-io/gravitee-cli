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

package org

import (
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

// NewOrgCmd creates the org parent command with all organization subcommands.
func NewOrgCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "org",
		Aliases: []string{"organization"},
		Short:   "Manage organization resources",
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newOrgUserCmd(f))
	cmd.AddCommand(newOrgGroupCmd(f))
	cmd.AddCommand(newOrgRoleCmd(f))
	cmd.AddCommand(newOrgSettingsCmd(f))
	cmd.AddCommand(newOrgMemberCmd(f))
	cmd.AddCommand(newOrgAuditCmd(f))
	cmd.AddCommand(newOrgReporterCmd(f))
	cmd.AddCommand(newOrgFormCmd(f))
	cmd.AddCommand(newOrgIDPCmd(f))
	cmd.AddCommand(newOrgEntrypointCmd(f))
	cmd.AddCommand(newOrgTagCmd(f))
	cmd.AddCommand(newOrgUserTokenCmd(f))

	return cmd
}
