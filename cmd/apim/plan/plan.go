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

package plan

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewPlanCmdRO creates the plan command with read-only subcommands.
func NewPlanCmdRO(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Manage API plans",
		Args:  cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)
	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))

	return cmd
}

// NewPlanCmd creates the parent plan command with all subcommands.
func NewPlanCmd(f *factory.Factory) *cobra.Command {
	cmd := NewPlanCmdRO(f)

	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(newDeleteCmd(f))
	cmd.AddCommand(newPublishCmd(f))
	cmd.AddCommand(newDeprecateCmd(f))
	cmd.AddCommand(newCloseCmd(f))

	return cmd
}
