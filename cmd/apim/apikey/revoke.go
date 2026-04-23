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

package apikey

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type revokeOptions struct {
	factory      *factory.Factory
	apiID        string
	subscription string
}

func newRevokeCmd(f *factory.Factory) *cobra.Command {
	opts := &revokeOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "revoke <keyId> --api <apiId> --subscription <subId>",
		Short:   "Revoke an API key",
		Example: `  gio apim api-key revoke 1a2b3c4d --api 8a7b3c4d --subscription aaaa1111`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run(args[0])
		},
	}

	cmdutil.AddAPIFlag(cmd, &opts.apiID)
	cmd.Flags().StringVar(&opts.subscription, "subscription", "", "Subscription ID (required)")

	_ = cmd.MarkFlagRequired("subscription")

	return cmd
}

func (o *revokeOptions) run(keyID string) error {
	f := o.factory

	if err := f.APIM().RevokeAPIKey(o.apiID, o.subscription, keyID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return cmdutil.PrintActionResult(p, keyID, "revoked",
		fmt.Sprintf("API key '%s' revoked.", keyID))
}
