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

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newConsentCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "consent",
		Short: "Manage user consents",
	}

	cmd.PersistentFlags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkPersistentFlagRequired("user-id")

	cmd.AddCommand(newConsentListCmd(f, domainID, &userID))
	cmd.AddCommand(newConsentRevokeCmd(f, domainID, &userID))
	cmd.AddCommand(newConsentRevokeAllCmd(f, domainID, &userID))

	return cmd
}

func newConsentListCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List user consents",
		Example: `  gio am user consent list --domain my-domain --user-id user-1
  gio am user consent list --domain my-domain --user-id user-1 -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runConsentList(f, *domainID, *userID)
		},
	}
}

func runConsentList(f *factory.Factory, domainID, userID string) error {
	items, err := f.AM().ListUserConsents(domainID, userID)
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

	return p.PrintList(items, consentColumns())
}

func newConsentRevokeCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "revoke <consentID>",
		Short:   "Revoke a user consent",
		Example: `  gio am user consent revoke consent-1 --domain my-domain --user-id user-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RevokeUserConsent(*domainID, *userID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Consent '%s' revoked.", args[0])

			return nil
		},
	}
}

func newConsentRevokeAllCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "revoke-all",
		Short:   "Revoke all user consents",
		Example: `  gio am user consent revoke-all --domain my-domain --user-id user-1`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RevokeAllUserConsents(*domainID, *userID); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("All consents revoked for user '%s'.", *userID)

			return nil
		},
	}
}

func consentColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Client ID", Value: func(i any) string { return cmdutil.StringField(i, "clientId") }},
		{Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "status") }},
	}
}
