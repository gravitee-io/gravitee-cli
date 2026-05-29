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
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newCertCredentialCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "cert-credential",
		Short: "Manage user certificate credentials",
	}

	cmd.PersistentFlags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkPersistentFlagRequired("user-id")

	cmd.AddCommand(newCertCredentialListCmd(f, domainID, &userID))
	cmd.AddCommand(newCertCredentialGetCmd(f, domainID, &userID))
	cmd.AddCommand(newCertCredentialEnrollCmd(f, domainID, &userID))
	cmd.AddCommand(newCertCredentialRevokeCmd(f, domainID, &userID))

	return cmd
}

func newCertCredentialListCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List user certificate credentials",
		Example: `  gctl am user cert-credential list --domain my-domain --user-id user-1
  gctl am user cert-credential list --domain my-domain --user-id user-1 -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runCertCredentialList(f, *domainID, *userID)
		},
	}
}

func runCertCredentialList(f *factory.Factory, domainID, userID string) error {
	items, err := f.AM().ListUserCertCredentials(domainID, userID)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(items)
	}

	return p.PrintList(items, certCredentialColumns())
}

func newCertCredentialGetCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get <credentialID>",
		Short:   "Get a user certificate credential",
		Example: `  gctl am user cert-credential get cred-1 --domain my-domain --user-id user-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetUserCertCredential(*domainID, *userID, args[0])
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
}

func newCertCredentialEnrollCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "enroll [-f <file>]",
		Short: "Enroll a user certificate credential from a JSON file or stdin",
		Example: `  gctl am user cert-credential enroll --domain my-domain --user-id user-1 --file cert.json
  gctl am user cert-credential enroll --domain my-domain --user-id user-1 -f cert.json
  envsubst < cert.json | gctl am user cert-credential enroll --domain my-domain --user-id user-1`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runCertCredentialEnroll(f, *domainID, *userID, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to a JSON file (optional - reads from stdin if omitted)")

	return cmd
}

func runCertCredentialEnroll(f *factory.Factory, domainID, userID, file string) error {
	body, err := cmdutil.ReadJSONInput(file, f.IOStreams.In)
	if err != nil {
		return err
	}

	data, err := f.AM().EnrollUserCertCredential(domainID, userID, body)
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

	p.PrintMessage("Certificate credential enrolled successfully.")

	return nil
}

func newCertCredentialRevokeCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "revoke <credentialID>",
		Short:   "Revoke a user certificate credential",
		Example: `  gctl am user cert-credential revoke cred-1 --domain my-domain --user-id user-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RevokeUserCertCredential(*domainID, *userID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Certificate credential '%s' revoked.", args[0])

			return nil
		},
	}
}

func certCredentialColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "status") }},
		{Name: "Created At", Value: func(i any) string { return cmdutil.StringField(i, "createdAt") }},
	}
}
