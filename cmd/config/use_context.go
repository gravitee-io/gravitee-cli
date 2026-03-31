package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newUseContextCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "use-context <name>",
		Short:   "Switch the current context",
		Example: `  gio config use-context prod`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			cfg := f.Config

			if _, ok := cfg.Contexts[name]; !ok {
				available := make([]string, 0, len(cfg.Contexts))
				for n := range cfg.Contexts {
					available = append(available, n)
				}

				sort.Strings(available)

				return fmt.Errorf("context '%s' not found\nHint: available contexts: %s. See 'gio config get-contexts'", name, strings.Join(available, ", "))
			}

			cfg.CurrentContext = name

			if err := cfg.SaveTo(f.ConfigPath); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Switched to context '%s'.\n", name)

			return nil
		},
	}
}
