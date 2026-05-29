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

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

type createOptions struct {
	factory            *factory.Factory
	domainID           *string
	username           string
	email              string
	password           string
	passwordStdin      bool
	firstName          string
	lastName           string
	preRegistration    bool
	preRegistrationSet bool
}

func newCreateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &createOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "create --username <username> --email <email>",
		Short: "Create a user",
		Example: `  gctl am user create --domain my-domain --username john --email john@example.com
  echo -n secret | gctl am user create --domain my-domain --username john --email john@example.com --password-stdin`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}
			if opts.passwordStdin || opts.password != "" {
				pw, err := cmdutil.ResolvePassword(opts.password, opts.passwordStdin, "Password: ", f.IOStreams.In, f.IOStreams.Err)
				if err != nil {
					return err
				}
				opts.password = pw
			}
			opts.preRegistrationSet = cmd.Flags().Changed("preRegistration")
			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.username, "username", "", "Username (required)")
	cmd.Flags().StringVar(&opts.email, "email", "", "Email address (required)")
	cmd.Flags().StringVar(&opts.password, "password", "", "Password (DEPRECATED: visible in process listings — prefer --password-stdin)")
	cmd.Flags().BoolVar(&opts.passwordStdin, "password-stdin", false, "Read password from stdin")
	cmd.Flags().StringVar(&opts.firstName, "firstName", "", "First name")
	cmd.Flags().StringVar(&opts.lastName, "lastName", "", "Last name")
	cmd.Flags().BoolVar(&opts.preRegistration, "preRegistration", false, "Pre-registration flag")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("email")

	return cmd
}

func (o *createOptions) run() error {
	f := o.factory

	body := map[string]any{
		"username": o.username,
		"email":    o.email,
	}

	if o.password != "" {
		body["password"] = o.password
	}

	if o.firstName != "" {
		body["firstName"] = o.firstName
	}

	if o.lastName != "" {
		body["lastName"] = o.lastName
	}

	if o.preRegistrationSet {
		body["preRegistration"] = o.preRegistration
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().CreateUser(*o.domainID, json.RawMessage(raw))
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
