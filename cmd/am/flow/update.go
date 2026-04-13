package flow

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newUpdateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update --file <flows.json>",
		Short: "Update flows from a JSON file (bulk update)",
		Example: `  gio am flow update --domain my-domain --file flows.json
  gio am flow update --domain my-domain -f flows.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdate(f, *domainID, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runUpdate(f *factory.Factory, domainID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateFlows(domainID, body)
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

	p.PrintMessage("Flows updated successfully.")

	return nil
}
