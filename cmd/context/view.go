package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newViewCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Display full details of the current context",
		Example: `  gio context view
  gio context view --context prod`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runView(f)
		},
	}
}

func runView(f *factory.Factory) error {
	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	cfg := f.Config

	contextName := cfg.Current
	if f.Overrides.Context != "" {
		contextName = f.Overrides.Context
	}

	if contextName == "" {
		return fmt.Errorf("no context configured\nHint: run 'gio login' to get started")
	}

	ctx, ok := cfg.Contexts[contextName]
	if !ok {
		return fmt.Errorf("context '%s' not found in config", contextName)
	}

	org := ctx.Org
	if org == "" {
		org = config.DefaultOrg
	}

	env := ctx.Env
	if env == "" {
		env = config.DefaultEnv
	}

	out := f.IOStreams.Out
	fmt.Fprintf(out, "Context:    %s\n", contextName)
	fmt.Fprintf(out, "Org:        %s\n", org)
	fmt.Fprintf(out, "Env:        %s\n", env)

	if ctx.APIM != nil {
		fmt.Fprintf(out, "APIM URL:   %s\n", ctx.APIM.URL)
		fmt.Fprintf(out, "APIM Token: %s\n", cmdutil.MaskToken(ctx.APIM.Token))
	}

	if ctx.AM != nil {
		fmt.Fprintf(out, "AM URL:     %s\n", ctx.AM.URL)
		fmt.Fprintf(out, "AM Token:   %s\n", cmdutil.MaskToken(ctx.AM.Token))
	}

	return nil
}
