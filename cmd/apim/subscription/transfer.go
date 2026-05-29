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

package subscription

import (
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

type transferOptions struct {
	factory *factory.Factory
	apiID   string
	planID  string
}

func newTransferCmd(f *factory.Factory) *cobra.Command {
	opts := &transferOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "transfer <subId> --api <apiId> --plan <planId>",
		Short:   "Transfer a subscription to another plan",
		Example: `  gctl apim subscription transfer 34f8c9e7 --api 8a7b3c4d --plan dd998877`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run(args[0])
		},
	}

	cmdutil.AddAPIFlag(cmd, &opts.apiID)
	cmd.Flags().StringVar(&opts.planID, "plan", "", "Target plan ID (required)")
	_ = cmd.MarkFlagRequired("plan")

	return cmd
}

func (o *transferOptions) run(subID string) error {
	f := o.factory

	data, err := f.APIM().TransferSubscription(o.apiID, subID, o.planID)
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

	return printSubDetail(p, data)
}
