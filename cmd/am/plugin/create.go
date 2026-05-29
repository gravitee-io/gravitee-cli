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

package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
	var name string
	var configFile string

	cmd := &cobra.Command{
		Use:   "create <type> <pluginId>",
		Short: "Create a resource instance from a plugin (use --config-file for non-interactive)",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			ptype, pluginID := args[0], args[1]
			apiPath, ok := pluginTypes[ptype]
			if !ok {
				return fmt.Errorf("unknown plugin type %q. Available: %s", ptype, sortedKeys(pluginTypes))
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if configFile == "" {
				return fmt.Errorf("--config-file is required (interactive mode not supported in gctl CLI)")
			}
			raw, err := cmdutil.ReadJSONFile(configFile)
			if err != nil {
				return err
			}
			body := map[string]interface{}{
				"name":          name,
				"type":          pluginID,
				"configuration": string(raw),
			}
			path := cmdutil.AMDomainPath(f, apiPath)
			data, err := f.Client.Post(path, body)
			if err != nil {
				return err
			}
			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}
			return p.PrintDetail(json.RawMessage(data))
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Resource name (required)")
	cmd.Flags().StringVarP(&configFile, "config-file", "f", "", "JSON config file (required)")
	return cmd
}
