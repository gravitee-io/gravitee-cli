package am

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newStatusCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current CLI context and session status",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := f.Config
			out := f.IOStreams.Out

			if cfg == nil || cfg.Current == "" {
				fmt.Fprintln(out, "context:       (not set)")
				fmt.Fprintln(out, "authenticated: no")

				return nil
			}

			ctx, ok := cfg.Contexts[cfg.Current]
			domain := ""

			if f.Resolved != nil {
				domain = f.Resolved.Domain
			}

			amURL := ""
			if ok && ctx != nil && ctx.AM != nil {
				amURL = ctx.AM.URL
			}

			fmt.Fprintf(out, "context:       %s @ %s\n", cfg.Current, amURL)

			if ok && ctx != nil {
				fmt.Fprintf(out, "organization:  %s\n", ctx.Org)
				fmt.Fprintf(out, "environment:   %s\n", ctx.Env)
			}

			if domain != "" {
				fmt.Fprintf(out, "domain:        %s\n", domain)
			} else {
				fmt.Fprintln(out, "domain:        (not set)")
			}

			if !ok || ctx == nil || ctx.AM == nil || ctx.AM.Token == "" {
				fmt.Fprintln(out, "authenticated: no")
			} else {
				fmt.Fprintln(out, "authenticated: yes")
			}

			return nil
		},
	}
}
