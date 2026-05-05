package am

import (
	"encoding/json"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newWhoamiCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show information about the currently authenticated user",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			return runWhoami(f)
		},
	}
}

func runWhoami(f *factory.Factory) error {
	data, err := f.Client.Get("/management/user")
	if err != nil {
		return err
	}
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(json.RawMessage(data))
}
