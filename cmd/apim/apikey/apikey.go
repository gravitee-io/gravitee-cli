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

package apikey

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAPIKeyCmd creates the parent api-key command with all subcommands.
func NewAPIKeyCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-key",
		Short: "Manage API keys",
		Args:  cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newRenewCmd(f))
	cmd.AddCommand(newRevokeCmd(f))
	cmd.AddCommand(newReactivateCmd(f))

	return cmd
}
