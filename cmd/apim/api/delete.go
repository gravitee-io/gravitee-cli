package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
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
		Example: `  gio apim api delete /my/api
  gio apim api delete /my/api --force`,
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
