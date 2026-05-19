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

package auth

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAuthCmd creates the `am auth` parent command. Bootstrap is a write
// operation (it creates a PAT and may persist it to config) so the read-only
// binary deliberately does not expose this subtree.
func NewAuthCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "auth",
		Short:  "Authentication helpers for first-mile CLI setup",
		Long:   "Authentication helpers intended for local development, CI lab setups, and demos.\nProduction users should mint PATs via the AM console and configure them with `gio context`.",
		Hidden: true,
	}

	cmd.AddCommand(newBootstrapCmd(f))

	return cmd
}
