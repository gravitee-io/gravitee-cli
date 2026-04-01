package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a context",
		Example: `  gio context delete staging
  gio context delete old-prod`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runDelete(f, args[0])
		},
	}
}

func runDelete(f *factory.Factory, name string) error {
	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	cfg := f.Config

	if err := cfg.DeleteContext(name); err != nil {
		return err
	}

	if err := cfg.SaveTo(f.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(f.IOStreams.Out, "Context '%s' deleted.\n", name)

	return nil
}
