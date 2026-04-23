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

package idp

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type listOptions struct {
	factory      *factory.Factory
	domainID     *string
	userProvider bool
}

func newListCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &listOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List identity providers",
		Example: `  gio am idp list --domain my-domain
  gio am idp list --domain my-domain --user-provider`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().BoolVar(&opts.userProvider, "user-provider", false, "Filter user providers only")

	return cmd
}

func (o *listOptions) run() error {
	f := o.factory

	items, err := f.AM().ListIdentityProviders(*o.domainID, o.userProvider)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(items)
	}

	return p.PrintList(items, idpColumns())
}

func idpColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Type", Value: func(i any) string { return cmdutil.StringField(i, "type") }},
		{Name: "Enabled", Value: func(i any) string {
			m, ok := i.(map[string]any)
			if !ok {
				return ""
			}

			if v, ok := m["enabled"].(bool); ok && v {
				return "true"
			}

			return "false"
		}},
	}
}
