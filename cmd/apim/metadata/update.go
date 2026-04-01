package metadata

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newUpdateCmd(f *factory.Factory) *cobra.Command {
	var (
		apiID string
		file  string
	)

	cmd := &cobra.Command{
		Use:     "update <key> --api <apiId> -f <file>",
		Short:   "Update a metadata entry from a JSON file",
		Example: `  gio apim metadata update team-email --api 8a7b3c4d-1234-5678-abcd-ef0123456789 -f metadata-updated.json`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdate(f, apiID, args[0], file)
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("api")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runUpdate(f *factory.Factory, apiID, key, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.APIM().UpdateMetadata(apiID, key, body)
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

	return printMetadataDetail(p, data, apiID)
}
