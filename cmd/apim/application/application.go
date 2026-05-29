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

package application

import (
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

// NewApplicationCmdRO creates the application command with read-only subcommands.
func NewApplicationCmdRO(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "application",
		Aliases: []string{"app"},
		Short:   "Manage applications",
		Args:    cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)
	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))

	return cmd
}

// NewApplicationCmd creates the parent application command with all subcommands.
func NewApplicationCmd(f *factory.Factory) *cobra.Command {
	cmd := NewApplicationCmdRO(f)

	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(newDeleteCmd(f))

	return cmd
}
