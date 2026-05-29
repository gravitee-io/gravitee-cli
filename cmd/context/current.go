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
)

func newCurrentCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "current",
		Short:   "Print the current context name",
		Example: `  gctl context current`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runCurrent(f)
		},
	}
}

func runCurrent(f *factory.Factory) error {
	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	cfg := f.Config

	if cfg.Current == "" {
		return fmt.Errorf("no context configured\nHint: run 'gctl login' to get started")
	}

	fmt.Fprintln(f.IOStreams.Out, cfg.Current)

	return nil
}
