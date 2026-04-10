package api

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newRollbackCmd(f *factory.Factory) *cobra.Command {
	var eventID string

	cmd := &cobra.Command{
		Use:     "rollback <apiId> --event-id <eventId>",
		Short:   "Rollback an API to a previous version",
		Example: `  gio apim api rollback /my/api --event-id aaaa1111-bbbb-2222-cccc-333344445555`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			apiID, err := f.APIM().ResolveAPI(args[0])
			if err != nil {
				return err
			}

			return runRollback(f, apiID, eventID)
		},
	}

	cmd.Flags().StringVar(&eventID, "event-id", "", "Event ID to rollback to (required)")
	_ = cmd.MarkFlagRequired("event-id")

	return cmd
}

func runRollback(f *factory.Factory, apiID, eventID string) error {
	if err := f.APIM().RollbackAPI(apiID, eventID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if printer.IsStructured(p.Format) || p.Format == printer.FormatID {
		return p.PrintDetail(map[string]string{
			"id":      apiID,
			"eventId": eventID,
			"status":  "rolled-back",
		})
	}

	p.PrintMessage("API '%s' rolled back to event '%s'.", apiID, eventID)

	return nil
}
