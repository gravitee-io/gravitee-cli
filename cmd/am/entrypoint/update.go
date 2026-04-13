package entrypoint

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newUpdateCmd(f *factory.Factory, _ *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <entrypointID> --file <entrypoint.json>",
		Short: "Update an entrypoint from a JSON file",
		Example: `  gio am entrypoint update my-entrypoint-id --domain my-domain --file entrypoint.json
  gio am entrypoint update my-entrypoint-id --domain my-domain -f entrypoint.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdate(f, args[0], file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runUpdate(f *factory.Factory, entrypointID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateEntrypoint(entrypointID, body)
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

	p.PrintMessage("Entrypoint '%s' updated.", entrypointID)

	return nil
}
