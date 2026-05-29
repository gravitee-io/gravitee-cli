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

package app

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newChangeTypeCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var (
		appID   string
		appType string
	)

	cmd := &cobra.Command{
		Use:   "change-type",
		Short: "Change the type of an application",
		Example: `  gctl am app change-type --domain my-domain --app-id my-app --type browser
  gctl am app change-type --domain my-domain --app-id my-app --type web`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.ValidateEnum(appType, "type", []string{"web", "native", "browser", "service", "resource_server"}); err != nil {
				return err
			}

			body, _ := json.Marshal(map[string]any{
				"type": strings.ToUpper(appType),
			})

			data, err := f.AM().ChangeAppType(*domainID, appID, json.RawMessage(body))
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

			p.PrintMessage(fmt.Sprintf("Application type changed to '%s'.", strings.ToUpper(appType)))

			return nil
		},
	}

	cmd.Flags().StringVar(&appID, "app-id", "", "Application ID (required)")
	cmd.Flags().StringVar(&appType, "type", "", "New application type: web, native, browser, service, resource_server (required)")
	_ = cmd.MarkFlagRequired("app-id")
	_ = cmd.MarkFlagRequired("type")

	return cmd
}
