package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newCurrentCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "current",
		Short:   "Print the current context name",
		Example: `  gio context current`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runCurrent(f)
		},
	}
}

func runCurrent(f *factory.Factory) error {
	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	cfg := f.Config

	if cfg.Current == "" {
		return fmt.Errorf("no context configured\nHint: run 'gio login' to get started")
	}

	fmt.Fprintln(f.IOStreams.Out, cfg.Current)

	return nil
}
