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

package api

import (
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newExportCmd(f *factory.Factory) *cobra.Command {
	var exclude []string

	cmd := &cobra.Command{
		Use:   "export <apiId>",
		Short: "Export an API definition",
		Example: `  gctl apim api export /my/api
  gctl apim api export 8a7b3c4d-1234-5678-abcd-ef0123456789 --exclude members --exclude pages`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			apiID, err := f.APIM().ResolveAPI(args[0])
			if err != nil {
				return err
			}

			return runExport(f, apiID, exclude)
		},
	}

	cmd.Flags().StringArrayVar(&exclude, "exclude", nil,
		"Exclude data from export: groups, members, metadata, pages, plans")

	return cmd
}

func runExport(f *factory.Factory, apiID string, exclude []string) error {
	data, err := f.APIM().ExportAPI(apiID, exclude)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}
