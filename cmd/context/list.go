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

package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

func newListCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all contexts",
		Example: `  gctl context list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runList(f)
		},
	}

	cmdutil.AddOutputFlags(cmd, f)

	return cmd
}

func runList(f *factory.Factory) error {
	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	cfg := f.Config
	names := cfg.ContextNames()

	if len(names) == 0 {
		fmt.Fprintln(f.IOStreams.Out, "No contexts configured.")

		return nil
	}

	items := make([]map[string]any, 0, len(names))
	for _, name := range names {
		ctx := cfg.Contexts[name]

		current := ""
		if name == cfg.Current {
			current = "*"
		}

		hasAPIM := "no"
		if ctx.APIM != nil {
			hasAPIM = "yes"
		}

		hasAM := "no"
		if ctx.AM != nil {
			hasAM = "yes"
		}

		items = append(items, map[string]any{
			"current": current,
			"name":    name,
			"org":     ctx.Org,
			"env":     ctx.Env,
			"apim":    hasAPIM,
			"am":      hasAM,
		})
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintList(items, contextColumns())
}

func contextColumns() []printer.Column {
	return []printer.Column{
		{Name: "Current", Value: func(i any) string { return cmdutil.StringField(i, "current") }, Width: 1},
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "Org", Value: func(i any) string { return cmdutil.StringField(i, "org") }},
		{Name: "Env", Value: func(i any) string { return cmdutil.StringField(i, "env") }},
		{Name: "APIM", Value: func(i any) string { return cmdutil.StringField(i, "apim") }},
		{Name: "AM", Value: func(i any) string { return cmdutil.StringField(i, "am") }},
	}
}
