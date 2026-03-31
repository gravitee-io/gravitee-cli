package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	iconfig "github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newViewCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Display the current context configuration",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runView(f)
		},
	}
}

func runView(f *factory.Factory) error {
	cfg := f.Config

	if cfg.CurrentContext == "" {
		return fmt.Errorf("no context configured\nHint: run 'gio login' to get started")
	}

	ctx, ok := cfg.Contexts[cfg.CurrentContext]
	if !ok {
		return fmt.Errorf("current context '%s' not found in config", cfg.CurrentContext)
	}

	env := ctx.Env
	if env == "" {
		env = iconfig.DefaultEnv
	}

	org := ctx.Org
	if org == "" {
		org = iconfig.DefaultOrg
	}

	ro := "no"
	if ctx.ReadOnly {
		ro = "yes"
	}

	out := f.IOStreams.Out
	fmt.Fprintf(out, "Context:    %s\n", cfg.CurrentContext)
	fmt.Fprintf(out, "URL:        %s\n", ctx.URL)
	fmt.Fprintf(out, "Org:        %s\n", org)
	fmt.Fprintf(out, "Env:        %s\n", env)
	fmt.Fprintf(out, "Read-only:  %s\n", ro)
	fmt.Fprintf(out, "Token:      %s\n", maskToken(ctx.Token))

	return nil
}

func maskToken(token string) string {
	if len(token) <= 3 {
		return "***"
	}

	return strings.Repeat("*", len(token)-3) + token[len(token)-3:]
}
