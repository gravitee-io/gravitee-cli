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

	"github.com/spf13/cobra"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newWhoamiCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show information about the currently authenticated user",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			return runWhoami(f)
		},
	}
}

func runWhoami(f *factory.Factory) error {
	data, err := f.Client.Get("/management/user")
	if err != nil {
		return err
	}
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(json.RawMessage(data))
}
