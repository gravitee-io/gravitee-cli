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

package certificate

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newKeyCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "key <certID>",
		Short:   "Get the public key of a certificate",
		Example: `  gio am certificate key cert-123 --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetCertificateKey(*domainID, args[0])
			if err != nil {
				return err
			}

			// /key returns raw text (SSH public key), not JSON.
			fmt.Fprintln(f.IOStreams.Out, string(data))

			return nil
		},
	}
}

func newKeysCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "keys <certID>",
		Short:   "Get all keys of a certificate",
		Example: `  gio am certificate keys cert-123 --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetCertificateKeys(*domainID, args[0])
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			return p.PrintDetail(data)
		},
	}
}
