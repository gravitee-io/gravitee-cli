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
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
	"github.com/spf13/cobra"
)

func newListCmd(f *factory.Factory) *cobra.Command {
	var userID string
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List user tokens",
		Example: `  gio am token list --user user-uuid`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			return runList(f, userID)
		},
	}
	cmd.Flags().StringVar(&userID, "user", "", "User ID (required)")
	_ = cmd.MarkFlagRequired("user")
	return cmd
}

func runList(f *factory.Factory, userID string) error {
	path := cmdutil.AMDomainPath(f, fmt.Sprintf("users/%s/tokens", userID))
	data, err := f.Client.Get(path)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}
	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}
	return p.PrintList(json.RawMessage(data), tokenColumns())
}

func tokenColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
		{Name: "Token", Value: func(i interface{}) string { return maskToken(cmdutil.StringField(i, "token")) }},
		{Name: "Expires At", Value: func(i interface{}) string { return cmdutil.StringField(i, "expiresAt") }},
		{Name: "Created At", Value: func(i interface{}) string { return cmdutil.StringField(i, "createdAt") }},
	}
}

// maskToken returns a redacted form of a bearer token: last 4 chars only.
// Tokens are bearer credentials — printing them in full to stdout puts them
// into terminal scrollback, shell history, and CI logs.
func maskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 4 {
		return "***"
	}
	return "***" + token[len(token)-4:]
}
