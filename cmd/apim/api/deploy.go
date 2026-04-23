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

package api

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeployCmd(f *factory.Factory) *cobra.Command {
	var label string

	cmd := &cobra.Command{
		Use:   "deploy <apiId>",
		Short: "Deploy an API",
		Example: `  gio apim api deploy /my/api
  gio apim api deploy 8a7b3c4d-1234-5678-abcd-ef0123456789 --label "v2.1.0 hotfix"`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			apiID, err := f.APIM().ResolveAPI(args[0])
			if err != nil {
				return err
			}

			return runDeploy(f, apiID, label)
		},
	}

	cmd.Flags().StringVar(&label, "label", "", "Deployment label (32 characters max)")

	return cmd
}

func runDeploy(f *factory.Factory, apiID, label string) error {
	if len(label) > 32 {
		return fmt.Errorf("deployment label exceeds 32 characters")
	}

	if err := f.APIM().DeployAPI(apiID, label); err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.Status == 400 {
			return fmt.Errorf("%w\nHint: ensure the API has at least one published plan before deploying ('gio apim plan publish')", err)
		}

		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return cmdutil.PrintActionResult(p, apiID, "deployed",
		fmt.Sprintf("API '%s' deployment requested.", apiID))
}
