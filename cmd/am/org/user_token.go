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

package org

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newOrgUserTokenCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user-token",
		Aliases: []string{"user-tokens"},
		Short:   "Manage organization user access tokens",
	}

	cmd.AddCommand(newOrgUserTokenListCmd(f))
	cmd.AddCommand(newOrgUserTokenCreateCmd(f))
	cmd.AddCommand(newOrgUserTokenRevokeCmd(f))

	return cmd
}

func newOrgUserTokenListCmd(f *factory.Factory) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:     "list --user-id <userID>",
		Short:   "List access tokens for an organization user",
		Example: `  gctl am org user-token list --user-id my-user-id`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListOrgUserTokens(userID)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			return p.PrintDetail(data)
		},
	}

	cmd.Flags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func newOrgUserTokenCreateCmd(f *factory.Factory) *cobra.Command {
	var (
		userID string
		name   string
	)

	cmd := &cobra.Command{
		Use:     "create --user-id <userID> --name <token-name>",
		Short:   "Create an access token for an organization user",
		Example: `  gctl am org user-token create --user-id my-user-id --name "my-token"`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			payload := map[string]string{"name": name}

			body, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("failed to build request body: %w", err)
			}

			data, err := f.AM().CreateOrgUserToken(userID, json.RawMessage(body))
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if f.OutputFormat != printer.FormatTable {
				return p.PrintDetail(data)
			}

			return printOrgUserTokenDetail(p, data)
		},
	}

	cmd.Flags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkFlagRequired("user-id")
	cmd.Flags().StringVar(&name, "name", "", "Token name (required)")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func newOrgUserTokenRevokeCmd(f *factory.Factory) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:     "revoke <tokenID> --user-id <userID>",
		Short:   "Revoke an access token for an organization user",
		Example: `  gctl am org user-token revoke my-token-id --user-id my-user-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RevokeOrgUserToken(userID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("User token '%s' revoked.", args[0])

			return nil
		},
	}

	cmd.Flags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func printOrgUserTokenDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Name", "name"},
		{"ID", "id"},
		{"Token", "token"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
