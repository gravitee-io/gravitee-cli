package botdetection

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newUpdateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <botDetectionID> --file <config.json>",
		Short: "Update a bot detection from a JSON file",
		Example: `  gio am bot-detection update my-bot-detection-id --domain my-domain --file bot-detection.json
  gio am bot-detection update my-bot-detection-id --domain my-domain -f bot-detection.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdate(f, *domainID, args[0], file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runUpdate(f *factory.Factory, domainID, botDetectionID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateBotDetection(domainID, botDetectionID, body)
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

	return printBotDetectionDetail(p, data)
}
