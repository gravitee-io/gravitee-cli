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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"gravitee.io/gctl/internal/client"
	"gravitee.io/gctl/internal/cmdutil"
	"gravitee.io/gctl/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <apiId>",
		Short: "Delete an API",
		Long: `Delete an API.

Without --force, the server rejects deletion if the API is running or has open
plans. With --force, the CLI stops the API (if running) and closes all plans
before deletion. Closing plans cascades to subscriptions server-side, so all
consumers lose access. This is irreversible.`,
		Example: `  gctl apim api delete /my/api
  gctl apim api delete /my/api --force`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			apiID, err := f.APIM().ResolveAPI(args[0])
			if err != nil {
				return err
			}

			return runDelete(f, apiID, force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false,
		"Stop the API if running and close all plans before deletion (cascades to subscriptions)")

	return cmd
}

func runDelete(f *factory.Factory, apiID string, force bool) error {
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	var stopped bool

	if force {
		stopped, err = stopIfRunning(f, apiID)
		if err != nil {
			return err
		}
	}

	if err := f.APIM().DeleteAPI(apiID, force); err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.Status == 400 {
			return fmt.Errorf("%w\nHint: the API is running or has open plans. Retry with --force to stop it and close plans in one step", err)
		}

		return err
	}

	var msg string
	switch {
	case stopped:
		msg = fmt.Sprintf("API '%s' stopped and deleted (plans closed).", apiID)
	case force:
		msg = fmt.Sprintf("API '%s' deleted (plans closed).", apiID)
	default:
		msg = fmt.Sprintf("API '%s' deleted.", apiID)
	}

	return cmdutil.PrintActionResult(p, apiID, "deleted", msg)
}

// stopIfRunning stops the API iff its runtime state is STARTED. Returns true if a stop was issued.
func stopIfRunning(f *factory.Factory, apiID string) (bool, error) {
	data, err := f.APIM().GetAPI(apiID)
	if err != nil {
		return false, err
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return false, fmt.Errorf("parse API state: %w", err)
	}

	if state, _ := m["state"].(string); state != "STARTED" {
		return false, nil
	}

	if err := f.APIM().StopAPI(apiID); err != nil {
		return false, err
	}

	return true, nil
}
