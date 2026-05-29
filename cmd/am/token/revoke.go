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

package token

import (
	"fmt"

	"github.com/spf13/cobra"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newRevokeCmd(f *factory.Factory) *cobra.Command {
	var userID string
	cmd := &cobra.Command{
		Use:     "revoke <tokenId>",
		Short:   "Revoke a user token",
		Example: `  gctl am token revoke token-id --user user-uuid`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			return runRevoke(f, userID, args[0])
		},
	}
	cmd.Flags().StringVar(&userID, "user", "", "User ID (required)")
	_ = cmd.MarkFlagRequired("user")
	return cmd
}

func runRevoke(f *factory.Factory, userID, tokenID string) error {
	path := cmdutil.AMDomainPath(f, fmt.Sprintf("users/%s/tokens/%s", userID, tokenID))
	if err := f.Client.Delete(path); err != nil {
		return err
	}
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}
	p.PrintMessage("Token '%s' revoked.", tokenID)
	return nil
}
