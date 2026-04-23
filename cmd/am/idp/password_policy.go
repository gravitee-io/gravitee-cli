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

package idp

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newPasswordPolicyCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var idpID string

	cmd := &cobra.Command{
		Use:     "password-policy",
		Aliases: []string{"pp"},
		Short:   "Manage identity provider password policy",
	}

	cmd.PersistentFlags().StringVar(&idpID, "idp-id", "", "Identity provider ID (required)")
	_ = cmd.MarkPersistentFlagRequired("idp-id")

	cmd.AddCommand(newPPAssignCmd(f, domainID, &idpID))

	return cmd
}

func newPPAssignCmd(f *factory.Factory, domainID, idpID *string) *cobra.Command {
	var policyID string

	cmd := &cobra.Command{
		Use:   "assign --policy-id <policyID>",
		Short: "Assign a password policy to an identity provider",
		Example: `  gio am idp password-policy assign --domain my-domain --idp-id my-idp --policy-id pp-123
  gio am idp password-policy assign --domain my-domain --idp-id my-idp --policy-id ""  # unassign`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, _ := json.Marshal(map[string]string{"passwordPolicy": policyID})

			data, err := f.AM().UpdateIDPPasswordPolicy(*domainID, *idpID, json.RawMessage(body))
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

			if policyID == "" {
				p.PrintMessage("Password policy unassigned from identity provider '%s'.", *idpID)
			} else {
				p.PrintMessage("Password policy '%s' assigned to identity provider '%s'.", policyID, *idpID)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&policyID, "policy-id", "", "Password policy ID (required, empty string to unassign)")
	_ = cmd.MarkFlagRequired("policy-id")

	return cmd
}
