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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type updateOptions struct {
	factory   *factory.Factory
	domainID  *string
	userID    string
	email     string
	firstName string
	lastName  string
	enabled   string
}

func newUpdateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &updateOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "update <userID>",
		Short: "Update a user",
		Example: `  gio am user update user-id --domain my-domain --email newemail@example.com
  gio am user update user-id --domain my-domain --firstName John --lastName Doe
  gio am user update user-id --domain my-domain --enabled false`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.userID = args[0]

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.email, "email", "", "Email address")
	cmd.Flags().StringVar(&opts.firstName, "firstName", "", "First name")
	cmd.Flags().StringVar(&opts.lastName, "lastName", "", "Last name")
	cmd.Flags().StringVar(&opts.enabled, "enabled", "", "Enable or disable the user (true/false)")

	return cmd
}

func (o *updateOptions) run() error {
	f := o.factory

	body := map[string]any{}

	if o.email != "" {
		body["email"] = o.email
	}

	if o.firstName != "" {
		body["firstName"] = o.firstName
	}

	if o.lastName != "" {
		body["lastName"] = o.lastName
	}

	if o.enabled != "" {
		switch o.enabled {
		case "true":
			body["enabled"] = true
		case "false":
			body["enabled"] = false
		default:
			return fmt.Errorf("invalid value '%s' for flag --enabled\nHint: allowed values are true, false", o.enabled)
		}
	}

	if len(body) == 0 {
		return fmt.Errorf("at least one flag (--email, --firstName, --lastName, --enabled) is required")
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().UpdateUser(*o.domainID, o.userID, json.RawMessage(raw))
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

	return printUserDetail(p, data)
}
