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
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

type reactivateOptions struct {
	factory      *factory.Factory
	apiID        string
	subscription string
}

func newReactivateCmd(f *factory.Factory) *cobra.Command {
	opts := &reactivateOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "reactivate <keyId> --api <apiId> --subscription <subId>",
		Short:   "Reactivate a previously revoked API key",
		Example: `  gctl apim api-key reactivate 1a2b3c4d --api 8a7b3c4d --subscription aaaa1111`,
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

func (o *reactivateOptions) run(keyID string) error {
	f := o.factory

	data, err := f.APIM().ReactivateAPIKey(o.apiID, o.subscription, keyID)
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

	return printKeyDetail(p, data)
}
