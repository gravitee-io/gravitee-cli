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

package cmdutil

import (
	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/factory"
)

// ResolveAPIMFlags rewrites the --api flag (if present and set) with the id
// resolved by the APIM service. Meant to run in PersistentPreRunE so every
// APIM subcommand accepting --api gets path → id resolution for free.
func ResolveAPIMFlags(f *factory.Factory, cmd *cobra.Command) error {
	flag := cmd.Flags().Lookup("api")
	if flag == nil {
		return nil
	}

	val := flag.Value.String()
	if val == "" {
		return nil
	}

	svc := f.APIM()
	if svc == nil {
		return nil
	}

	resolved, err := svc.ResolveAPI(val)
	if err != nil {
		return err
	}

	return flag.Value.Set(resolved)
}
