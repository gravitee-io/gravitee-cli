package metadata

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
	var (
		apiID string
		file  string
	)

	cmd := &cobra.Command{
		Use:     "create --api <apiId> -f <file>",
		Short:   "Create a metadata entry from a JSON file",
		Example: `  gio apim metadata create --api 8a7b3c4d-1234-5678-abcd-ef0123456789 -f metadata.json`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runCreate(f, apiID, file)
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("api")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runCreate(f *factory.Factory, apiID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.APIM().CreateMetadata(apiID, body)
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
