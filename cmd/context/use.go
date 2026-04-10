package context

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newUseCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Switch to a different context",
		Example: `  gio context use prod
  gio context use staging`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runUse(f, args[0])
		},
	}
}

func runUse(f *factory.Factory, name string) error {
	name = config.NormalizeContextName(name)

	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	cfg := f.Config

	if _, ok := cfg.Contexts[name]; !ok {
		available := cfg.ContextNames()
		if len(available) == 0 {
			return fmt.Errorf("context '%s' not found\nHint: no contexts configured, run 'gio login' to get started", name)
		}

		return fmt.Errorf("context '%s' not found\nHint: available contexts: %s", name, strings.Join(available, ", "))
	}

	cfg.Current = name

	if err := cfg.SaveTo(f.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(f.IOStreams.Out, "Switched to context '%s'.\n", name)

	return nil
}
