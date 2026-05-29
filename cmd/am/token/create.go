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

	"github.com/spf13/cobra"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
	var userID, file string
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a user token",
		Example: `  gctl am token create --user user-uuid`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			return runCreate(f, userID, file)
		},
	}
	cmd.Flags().StringVar(&userID, "user", "", "User ID (required)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "JSON file with token definition")
	_ = cmd.MarkFlagRequired("user")
	return cmd
}

func runCreate(f *factory.Factory, userID, file string) error {
	var body json.RawMessage
	if file != "" {
		var err error
		body, err = cmdutil.ReadJSONFile(file)
		if err != nil {
			return err
		}
	} else {
		body = json.RawMessage(`{}`)
	}

	path := cmdutil.AMDomainPath(f, fmt.Sprintf("users/%s/tokens", userID))
	data, err := f.Client.Post(path, body)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	id := cmdutil.StringField(m, "id")
	tokenValue := cmdutil.StringField(m, "token")

	if tokenValue != "" {
		p.PrintMessage("Token created (ID: %s).", id)
		p.PrintMessage("")
		p.PrintMessage("Token value (store it now — it will not be shown again):")
		p.PrintMessage("  %s", tokenValue)
		return nil
	}
	p.PrintMessage("Token created (ID: %s).", id)
	return nil
}
