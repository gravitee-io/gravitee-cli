package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type loginProductOptions struct {
	factory     *factory.Factory
	product     string
	url         string
	token       string
	contextName string
	org         string
	envID       string
	readOnly    bool
}

func newLoginCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [apim|am]",
		Short: "Configure credentials for a Gravitee product",
		Example: `  gio login apim --url https://apim.company.com --token gioat_abc123
  gio login apim --url https://apim.company.com --token gioat_abc123 --context prod --org ACME --env production
  gio login am --url https://am.company.com --token gioat_abc123 --context prod
  gio login`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runInteractiveLogin(f)
		},
	}

	cmd.AddCommand(newLoginProductCmd(f, "apim"))
	cmd.AddCommand(newLoginProductCmd(f, "am"))

	return cmd
}

func newLoginProductCmd(f *factory.Factory, product string) *cobra.Command {
	opts := &loginProductOptions{factory: f, product: product}

	cmd := &cobra.Command{
		Use:   product,
		Short: fmt.Sprintf("Configure credentials for a Gravitee %s instance", strings.ToUpper(product)),
		Example: fmt.Sprintf(`  gio login %s --url https://%s.company.com --token gioat_abc123
  gio login %s --url https://%s.company.com --token gioat_abc123 --context prod --org ACME --env production`,
			product, product, product, product),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return opts.run(cmd)
		},
	}

	cmd.Flags().StringVar(&opts.url, "url", "", "URL of the Gravitee control plane (required)")
	cmd.Flags().StringVar(&opts.token, "token", "", "Personal Access Token (required)")
	cmd.Flags().StringVar(&opts.contextName, "context", "default", "Context name")
	cmd.Flags().StringVar(&opts.org, "org", config.DefaultOrg, "Organization ID")
	cmd.Flags().StringVar(&opts.envID, "env-id", config.DefaultEnv, "Environment ID")
	cmd.Flags().BoolVar(&opts.readOnly, "read-only", false, "Enable read-only mode")
	_ = cmd.MarkFlagRequired("url")
	_ = cmd.MarkFlagRequired("token")

	return cmd
}

func (o *loginProductOptions) run(cmd *cobra.Command) error {
	if err := cmdutil.ValidateURL(o.url); err != nil {
		return err
	}

	if err := cmdutil.SetupConfig(o.factory); err != nil {
		return err
	}

	cfg := o.factory.Config
	ctx := cfg.EnsureContext(o.contextName)

	if cmd.Flags().Changed("org") {
		ctx.Org = o.org
	}

	if cmd.Flags().Changed("env-id") {
		ctx.Env = o.envID
	}

	if cmd.Flags().Changed("read-only") {
		ctx.ReadOnly = o.readOnly
	}

	ctx.SetProductConfig(o.product, &config.ProductConfig{
		URL:   o.url,
		Token: o.token,
	})

	cfg.Current = o.contextName

	if err := cfg.SaveTo(o.factory.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(o.factory.IOStreams.Out, "Context '%s' saved and set as current (%s configured).\n", o.contextName, strings.ToUpper(o.product))

	return nil
}

func runInteractiveLogin(f *factory.Factory) error {
	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	out := f.IOStreams.Out
	in := f.IOStreams.In

	// Prompt for product.
	fmt.Fprintln(out, "Which product do you want to configure?")
	fmt.Fprintln(out, "  [1] apim")
	fmt.Fprintln(out, "  [2] am")
	fmt.Fprint(out, "Choice: ")

	var choice string
	if _, err := fmt.Fscanln(in, &choice); err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	var product string

	switch strings.TrimSpace(choice) {
	case "1", "apim":
		product = "apim"
	case "2", "am":
		product = "am"
	default:
		return fmt.Errorf("invalid choice %q\nHint: enter 1 for apim or 2 for am", choice)
	}

	// Prompt for URL.
	fmt.Fprint(out, "URL: ")

	var url string
	if _, err := fmt.Fscanln(in, &url); err != nil {
		return fmt.Errorf("failed to read URL: %w", err)
	}

	if err := cmdutil.ValidateURL(url); err != nil {
		return err
	}

	// Prompt for token.
	fmt.Fprint(out, "Token: ")

	var token string
	if _, err := fmt.Fscanln(in, &token); err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}

	if token == "" {
		return fmt.Errorf("token is required")
	}

	// Prompt for context name.
	cfg := f.Config
	names := cfg.ContextNames()

	if len(names) > 0 {
		fmt.Fprintln(out, "Existing contexts:")

		for i, name := range names {
			marker := "  "
			if name == cfg.Current {
				marker = "* "
			}

			fmt.Fprintf(out, "  %s[%d] %s\n", marker, i+1, name)
		}

		fmt.Fprintf(out, "  [%d] Create new context\n", len(names)+1)
	}

	fmt.Fprint(out, "Context name (default): ")

	var contextName string
	if _, err := fmt.Fscanln(in, &contextName); err != nil {
		// Empty input means use default.
		contextName = "default"
	}

	if contextName == "" {
		contextName = "default"
	}

	ctx := cfg.EnsureContext(contextName)

	// Prompt for org (only if not already set on the context).
	if ctx.Org == "" {
		fmt.Fprintf(out, "Organization ID (%s): ", config.DefaultOrg)

		var org string
		if _, err := fmt.Fscanln(in, &org); err != nil {
			org = ""
		}

		if org != "" {
			ctx.Org = org
		}
	}

	// Prompt for env (only if not already set on the context).
	if ctx.Env == "" {
		fmt.Fprintf(out, "Environment ID (%s): ", config.DefaultEnv)

		var env string
		if _, err := fmt.Fscanln(in, &env); err != nil {
			env = ""
		}

		if env != "" {
			ctx.Env = env
		}
	}

	ctx.SetProductConfig(product, &config.ProductConfig{
		URL:   url,
		Token: token,
	})

	cfg.Current = contextName

	if err := cfg.SaveTo(f.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(out, "Context '%s' saved and set as current (%s configured).\n", contextName, strings.ToUpper(product))

	return nil
}
