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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	amcmd "github.com/gravitee-io/gio-cli/cmd/am"
	apimcmd "github.com/gravitee-io/gio-cli/cmd/apim"
	contextcmd "github.com/gravitee-io/gio-cli/cmd/context"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewRootCmd creates the root gio command.
func NewRootCmd(version string) *cobra.Command {
	f := &factory.Factory{
		IOStreams: factory.DefaultIOStreams(),
	}

	cmd := &cobra.Command{
		Use:           "gio",
		Short:         "gio - Gravitee CLI",
		Long:          "gio is a command-line interface for the Gravitee platform.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&f.Overrides.Context, "context", "", "Override context")
	cmd.PersistentFlags().StringVar(&f.Overrides.Org, "org", "", "Override organization ID")
	cmd.PersistentFlags().StringVar(&f.Overrides.EnvID, "env", "", "Override environment ID")
	cmd.PersistentFlags().BoolVar(&f.Debug, "debug", false, "Show raw HTTP requests/responses")

	cmd.AddCommand(newLoginCmd(f))
	cmd.AddCommand(contextcmd.NewContextCmd(f))
	cmd.AddCommand(newCompletionCmd())
	cmd.AddCommand(newVersionCmd(f, version))

	cmd.AddCommand(apimcmd.NewAPIMCmd(f))
	cmd.AddCommand(amcmd.NewAMCmd(f))

	return cmd
}

// NewRootCmdRO creates the root gio-ro command with read-only subcommands only.
func NewRootCmdRO(version string) *cobra.Command {
	f := &factory.Factory{
		IOStreams: factory.DefaultIOStreams(),
	}

	cmd := &cobra.Command{
		Use:           "gio-ro",
		Short:         "gio-ro - Gravitee CLI (read-only)",
		Long:          "gio-ro is a read-only command-line interface for the Gravitee platform. Only commands that do not modify state are available.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&f.Overrides.Context, "context", "", "Override context")
	cmd.PersistentFlags().StringVar(&f.Overrides.Org, "org", "", "Override organization ID")
	cmd.PersistentFlags().StringVar(&f.Overrides.EnvID, "env", "", "Override environment ID")
	cmd.PersistentFlags().BoolVar(&f.Debug, "debug", false, "Show raw HTTP requests/responses")

	cmd.AddCommand(newLoginCmd(f))
	cmd.AddCommand(contextcmd.NewContextCmdRO(f))
	cmd.AddCommand(newCompletionCmd())
	cmd.AddCommand(newVersionCmd(f, version))

	cmd.AddCommand(apimcmd.NewAPIMCmdRO(f))

	return cmd
}

// Execute runs the root command.
func Execute(version string) int {
	cmd := NewRootCmd(version)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "Error:", err)

		return 1
	}

	return 0
}

// ExecuteRO runs the read-only root command.
func ExecuteRO(version string) int {
	cmd := NewRootCmdRO(version)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "Error:", err)

		return 1
	}

	return 0
}
