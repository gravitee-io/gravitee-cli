package user

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newBulkCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "bulk --file <operations.json>",
		Short: "Perform bulk user operations from a JSON file",
		Example: `  gio am user bulk --domain my-domain --file operations.json
  gio am user bulk --domain my-domain -f operations.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runBulk(f, *domainID, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON operations file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runBulk(f *factory.Factory, domainID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().BulkUserOperation(domainID, body)
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

	p.PrintMessage("Bulk operation completed successfully.")

	return nil
}
