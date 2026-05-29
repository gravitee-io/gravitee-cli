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

package health

import (
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

// NewHealthCmd creates the health command.
func NewHealthCmd(f *factory.Factory) *cobra.Command {
	var (
		apiID string
		field string
	)

	cmd := &cobra.Command{
		Use:     "health --api <apiId>",
		Short:   "Get API health check availability",
		Example: `  gctl apim health --api /my/api`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			resolvedID, err := f.APIM().ResolveAPI(apiID)
			if err != nil {
				return err
			}

			data, err := f.APIM().GetAPIHealth(resolvedID, field)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			if len(data) == 0 {
				p.PrintMessage("No health check data available for this API.")

				return nil
			}

			return p.PrintDetail(data)
		},
	}

	cmdutil.AddOutputFlags(cmd, f)
	cmdutil.AddAPIFlag(cmd, &apiID)
	cmd.Flags().StringVar(&field, "field", "endpoint", "Grouping field")

	return cmd
}
