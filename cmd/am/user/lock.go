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

package user

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newLockCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "lock <userID>",
		Short:   "Lock a user account",
		Example: `  gio am user lock user-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdateStatus(f, *domainID, args[0], false)
		},
	}
}

func newUnlockCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "unlock <userID>",
		Short:   "Unlock a user account",
		Example: `  gio am user unlock user-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdateStatus(f, *domainID, args[0], true)
		},
	}
}

func runUpdateStatus(f *factory.Factory, domainID, userID string, enabled bool) error {
	body := map[string]any{"enabled": enabled}
	raw, _ := json.Marshal(body)

	if _, err := f.AM().UpdateUserStatus(domainID, userID, json.RawMessage(raw)); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	action := "unlocked"
	if !enabled {
		action = "locked"
	}

	p.PrintMessage("User '%s' %s.", userID, action)

	return nil
}
