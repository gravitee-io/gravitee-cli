package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newCurrentContextCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "current-context",
		Short: "Display the current context name",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := f.Config

			if cfg.CurrentContext == "" {
				return fmt.Errorf("no context configured\nHint: run 'gio login' to get started")
			}

			fmt.Fprintln(f.IOStreams.Out, cfg.CurrentContext)

			return nil
		},
	}
}
