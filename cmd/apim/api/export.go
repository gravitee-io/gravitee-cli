package api

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newExportCmd(f *factory.Factory) *cobra.Command {
	var exclude []string

	cmd := &cobra.Command{
		Use:   "export <apiId>",
		Short: "Export an API definition",
		Example: `  gio apim api export 8a7b3c4d-1234-5678-abcd-ef0123456789
  gio apim api export 8a7b3c4d-1234-5678-abcd-ef0123456789 --exclude members --exclude pages`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runExport(f, args[0], exclude)
		},
	}

	cmd.Flags().StringArrayVar(&exclude, "exclude", nil,
		"Exclude data from export: groups, members, metadata, pages, plans")

	return cmd
}

func runExport(f *factory.Factory, apiID string, exclude []string) error {
	data, err := f.APIM().ExportAPI(apiID, exclude)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}
