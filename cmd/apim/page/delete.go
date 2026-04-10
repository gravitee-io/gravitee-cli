package page

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "delete <pageId> --api <apiId>",
		Short:   "Delete a page",
		Example: `  gio apim page delete dddd1111-2222-3333-4444-555566667777 --api /my/api`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDelete(f, apiID, args[0])
		},
	}

	cmdutil.AddAPIFlag(cmd, &apiID)

	return cmd
}

func runDelete(f *factory.Factory, apiID, pageID string) error {
	if err := f.APIM().DeletePage(apiID, pageID); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return cmdutil.PrintActionResult(p, pageID, "deleted",
		fmt.Sprintf("Page '%s' deleted.", pageID))
}
