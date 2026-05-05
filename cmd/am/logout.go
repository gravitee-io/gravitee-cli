package am

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newLogoutCmd(f *factory.Factory) *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Clear stored authentication token",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := f.Config
			if all {
				count := 0

				for _, ctx := range cfg.Contexts {
					if ctx.AM != nil && ctx.AM.Token != "" {
						ctx.AM.Token = ""
						count++
					}
				}

				if count == 0 {
					fmt.Fprintln(f.IOStreams.Out, "No stored AM tokens to clear.")
					return nil
				}

				if f.ConfigPath != "" {
					if err := cfg.SaveTo(f.ConfigPath); err != nil {
						return fmt.Errorf("failed to save config: %w", err)
					}
				}

				fmt.Fprintf(f.IOStreams.Out, "Cleared AM tokens for %d context(s).\n", count)

				return nil
			}

			if cfg.Current == "" {
				fmt.Fprintln(f.IOStreams.Out, "No context selected.")
				return nil
			}

			ctx, ok := cfg.Contexts[cfg.Current]
			if !ok || ctx.AM == nil || ctx.AM.Token == "" {
				fmt.Fprintf(f.IOStreams.Out, "No AM token stored for context %q.\n", cfg.Current)
				return nil
			}

			ctx.AM.Token = ""

			if f.ConfigPath != "" {
				if err := cfg.SaveTo(f.ConfigPath); err != nil {
					return fmt.Errorf("failed to save config: %w", err)
				}
			}

			fmt.Fprintf(f.IOStreams.Out, "Logged out from context %q.\n", cfg.Current)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Clear AM tokens for all contexts")

	return cmd
}
