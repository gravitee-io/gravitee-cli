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

package am

import (
	"encoding/json"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
	"github.com/spf13/cobra"
)

func newHealthCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "health",
		Aliases: []string{"ping"},
		Short:   "Check if the AM instance is reachable",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			return runHealth(f)
		},
	}
}

func runHealth(f *factory.Factory) error {
	data, err := f.Client.Get("/management/health")
	if err != nil {
		return err
	}
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	p.PrintMessage("AM instance is healthy.")

	return nil
}
