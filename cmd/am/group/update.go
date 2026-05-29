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

package group

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
	"gravitee.io/gctl/internal/printer"
)

type updateOptions struct {
	factory     *factory.Factory
	domainID    *string
	groupID     string
	name        string
	description string
}

func newUpdateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &updateOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "update <groupID>",
		Short: "Update a group",
		Example: `  gctl am group update my-group-id --domain my-domain --name "New Name"
  gctl am group update my-group-id --domain my-domain --description "Updated description"`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.groupID = args[0]

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Group name")
	cmd.Flags().StringVar(&opts.description, "description", "", "Group description")

	return cmd
}

func (o *updateOptions) run() error {
	f := o.factory

	body := map[string]any{}
	if o.name != "" {
		body["name"] = o.name
	}

	if o.description != "" {
		body["description"] = o.description
	}

	if len(body) == 0 {
		return fmt.Errorf("at least one flag (--name, --description) is required")
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().UpdateGroup(*o.domainID, o.groupID, json.RawMessage(raw))
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

	return printGroupDetail(p, data)
}
