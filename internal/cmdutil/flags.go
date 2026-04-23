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

package cmdutil

import (
	"github.com/spf13/cobra"
)

// AddAPIFlag registers the standard --api flag on cmd, bound to target,
// and marks it required. All APIM subcommands that target a specific API
// should use this so the flag's behavior (description, shortcut, validation)
// lives in one place.
func AddAPIFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "api", "",
		"API id or context path (e.g. /my/api) (required)")
	_ = cmd.MarkFlagRequired("api")

	// MarkFlagRequired checks presence, not content: --api "" would pass.
	// Reject it explicitly before RunE so callers never hit the server with
	// an empty id (which silently resolves to wrong endpoints - see B3).
	prev := cmd.PreRunE
	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		if prev != nil {
			if err := prev(c, args); err != nil {
				return err
			}
		}

		return RequireNonEmpty("--api", *target)
	}
}
