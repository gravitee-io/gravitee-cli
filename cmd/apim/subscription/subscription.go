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

package subscription

import (
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

// NewSubscriptionCmdRO creates the subscription command with read-only subcommands.
func NewSubscriptionCmdRO(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subscription",
		Aliases: []string{"sub"},
		Short:   "Manage subscriptions",
		Args:    cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)
	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))

	return cmd
}

// NewSubscriptionCmd creates the parent subscription command with all subcommands.
func NewSubscriptionCmd(f *factory.Factory) *cobra.Command {
	cmd := NewSubscriptionCmdRO(f)

	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newAcceptCmd(f))
	cmd.AddCommand(newRejectCmd(f))
	cmd.AddCommand(newPauseCmd(f))
	cmd.AddCommand(newResumeCmd(f))
	cmd.AddCommand(newCloseCmd(f))
	cmd.AddCommand(newTransferCmd(f))

	return cmd
}
