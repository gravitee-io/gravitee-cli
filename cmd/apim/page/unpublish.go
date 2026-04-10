package page

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newUnpublishCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "unpublish <pageId> --api <apiId>",
		Short:   "Unpublish a page",
		Example: `  gio apim page unpublish dddd1111-2222-3333-4444-555566667777 --api /my/api`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUnpublish(f, apiID, args[0])
		},
	}

	cmdutil.AddAPIFlag(cmd, &apiID)

	return cmd
}

func runUnpublish(f *factory.Factory, apiID, pageID string) error {
	data, err := f.APIM().UnpublishPage(apiID, pageID)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(data)
	}

	return printPageDetail(p, data)
}
