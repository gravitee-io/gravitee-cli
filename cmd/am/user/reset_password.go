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
)

type resetPasswordOptions struct {
	factory       *factory.Factory
	domainID      *string
	password      string
	passwordStdin bool
}

func newResetPasswordCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &resetPasswordOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "reset-password <userID>",
		Short: "Reset a user's password",
		Example: `  echo -n 'newSecret123' | gctl am user reset-password user-id --password-stdin
  gctl am user reset-password user-id  # interactive prompt`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}
			pw, err := cmdutil.ResolvePassword(opts.password, opts.passwordStdin, "New password: ", f.IOStreams.In, f.IOStreams.Err)
			if err != nil {
				return err
			}
			return opts.run(args[0], pw)
		},
	}

	cmd.Flags().StringVar(&opts.password, "password", "", "New password (DEPRECATED: visible in process listings — prefer --password-stdin)")
	cmd.Flags().BoolVar(&opts.passwordStdin, "password-stdin", false, "Read new password from stdin")

	return cmd
}

func (o *resetPasswordOptions) run(userID, password string) error {
	f := o.factory

	body := map[string]any{"password": password}
	raw, _ := json.Marshal(body)

	if err := f.AM().ResetPassword(*o.domainID, userID, json.RawMessage(raw)); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Password reset for user '%s'.", userID)

	return nil
}
