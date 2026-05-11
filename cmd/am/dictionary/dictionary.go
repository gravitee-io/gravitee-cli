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

package dictionary

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewDictionaryCmdRO creates the dictionary command with read-only subcommands.
func NewDictionaryCmdRO(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "dictionary",
		Aliases: []string{"dictionaries", "i18n"},
		Short:   "Manage i18n dictionaries",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))

	return cmd
}

// NewDictionaryCmd creates the dictionary parent command with all subcommands.
func NewDictionaryCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "dictionary",
		Aliases: []string{"dictionaries", "i18n"},
		Short:   "Manage i18n dictionaries",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))
	cmd.AddCommand(newCreateCmd(f, &domainID))
	cmd.AddCommand(newUpdateCmd(f, &domainID))
	cmd.AddCommand(newDeleteCmd(f, &domainID))
	cmd.AddCommand(newEntryCmd(f, &domainID))

	return cmd
}
