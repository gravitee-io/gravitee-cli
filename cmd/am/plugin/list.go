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
	"sort"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
	"github.com/spf13/cobra"
)

func newListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list <type>",
		Short: "List available plugins of a given type",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			ptype := args[0]
			apiType, ok := pluginTypes[ptype]
			if !ok {
				return fmt.Errorf("unknown plugin type %q. Available: %s", ptype, sortedKeys(pluginTypes))
			}
			data, err := f.Client.Get("/management/platform/plugins/" + apiType)
			if err != nil {
				return err
			}
			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}
			var items []json.RawMessage
			if err := json.Unmarshal(data, &items); err != nil {
				return p.PrintDetail(json.RawMessage(data))
			}
			return p.PrintList(items, []printer.Column{
				{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
				{Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
				{Name: "Version", Value: func(i interface{}) string { return cmdutil.StringField(i, "version") }},
			})
		},
	}
}

func sortedKeys(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
