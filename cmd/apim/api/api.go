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

package api

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAPICmd creates the parent api command with all subcommands.
func NewAPICmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Manage APIs",
		Args:  cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(newDeleteCmd(f))
	cmd.AddCommand(newStartCmd(f))
	cmd.AddCommand(newStopCmd(f))
	cmd.AddCommand(newDeployCmd(f))
	cmd.AddCommand(newImportCmd(f))
	cmd.AddCommand(newExportCmd(f))
	cmd.AddCommand(newRollbackCmd(f))
	cmd.AddCommand(newAnalyticsCmd(f))
	cmd.AddCommand(newHealthCmd(f))
	cmd.AddCommand(newLogsCmd(f))
	cmd.AddCommand(newLogCmd(f))

	return cmd
}
