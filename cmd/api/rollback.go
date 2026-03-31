package api

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newRollbackCmd(f *factory.Factory) *cobra.Command {
	var eventID string

	cmd := &cobra.Command{
		Use:     "rollback <apiId> --event-id <eventId>",
		Short:   "Rollback an API to a previous version",
		Example: `  gio api rollback 8a7b3c4d-1234-5678-abcd-ef0123456789 --event-id aaaa1111-bbbb-2222-cccc-333344445555`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "api rollback"); err != nil {
				return err
			}

			return runRollback(f, args[0], eventID)
		},
	}

	cmd.Flags().StringVar(&eventID, "event-id", "", "Event ID to rollback to (required)")
	_ = cmd.MarkFlagRequired("event-id")

	return cmd
}

func runRollback(f *factory.Factory, apiID, eventID string) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/_rollback", apiID))
	body := map[string]string{"eventId": eventID}

	if _, err := f.Client.Post(path, body); err != nil {
		return fmt.Errorf("API rollback failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)
	p.PrintMessage("API '%s' rolled back to event '%s'.", apiID, eventID)

	return nil
}
