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

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newPauseCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "pause <subId> --api <apiId>",
		Short:   "Pause an active subscription",
		Example: `  gio apim subscription pause 34f8c9e7 --api 8a7b3c4d`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runPause(f, apiID, args[0])
		},
	}

	cmdutil.AddAPIFlag(cmd, &apiID)

	return cmd
}

func runPause(f *factory.Factory, apiID, subID string) error {
	data, err := f.APIM().PauseSubscription(apiID, subID)
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
