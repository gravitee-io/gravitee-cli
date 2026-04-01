package api

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
	var closePlans bool

	cmd := &cobra.Command{
		Use:   "delete <apiId>",
		Short: "Delete an API",
		Example: `  gio apim api delete 8a7b3c4d-1234-5678-abcd-ef0123456789
  gio apim api delete 8a7b3c4d-1234-5678-abcd-ef0123456789 --close-plans`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDelete(f, args[0], closePlans)
		},
	}

	cmd.Flags().BoolVar(&closePlans, "close-plans", false, "Force deletion by closing API plans")

	return cmd
}

func runDelete(f *factory.Factory, apiID string, closePlans bool) error {
	if err := f.APIM().DeleteAPI(apiID, closePlans); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}
	p.PrintMessage("API '%s' deleted.", apiID)

	return nil
}
