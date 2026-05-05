package am

import (
	"encoding/json"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
	"github.com/spf13/cobra"
)

func newHealthCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "health",
		Aliases: []string{"ping"},
		Short:   "Check if the AM instance is reachable",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			return runHealth(f)
		},
	}
}

func runHealth(f *factory.Factory) error {
	data, err := f.Client.Get("/management/health")
	if err != nil {
		return err
	}
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	p.PrintMessage("AM instance is healthy.")

	return nil
}
