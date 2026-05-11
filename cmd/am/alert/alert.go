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

package alert

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAlertCmdRO creates the alert command with read-only notifier and trigger subcommands.
func NewAlertCmdRO(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "alert",
		Aliases: []string{"alerts"},
		Short:   "Manage alerts (notifiers and triggers)",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	notifierCmd := &cobra.Command{
		Use:     "notifier",
		Aliases: []string{"notifiers"},
		Short:   "Manage alert notifiers",
	}
	notifierCmd.AddCommand(newNotifierListCmd(f, &domainID))
	notifierCmd.AddCommand(newNotifierGetCmd(f, &domainID))
	cmd.AddCommand(notifierCmd)

	triggerCmd := &cobra.Command{
		Use:     "trigger",
		Aliases: []string{"triggers"},
		Short:   "Manage alert triggers",
	}
	triggerCmd.AddCommand(newTriggerGetCmd(f, &domainID))
	cmd.AddCommand(triggerCmd)

	return cmd
}

// NewAlertCmd creates the alert parent command with notifier and trigger subcommands.
func NewAlertCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "alert",
		Aliases: []string{"alerts"},
		Short:   "Manage alerts (notifiers and triggers)",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newNotifierCmd(f, &domainID))
	cmd.AddCommand(newTriggerCmd(f, &domainID))

	return cmd
}
