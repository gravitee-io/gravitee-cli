package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type loginOptions struct {
	factory     *factory.Factory
	url         string
	token       string
	contextName string
	org         string
	envID       string
	readOnly    bool
}

func newLoginCmd(f *factory.Factory) *cobra.Command {
	opts := &loginOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Configure credentials for a Gravitee APIM instance",
		Example: `  gio login --url https://apim.company.com --token gioat_abc123
  gio login --url https://apim.company.com --token gioat_abc123 --context prod --env-id production --read-only`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return opts.run()
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

func (o *loginOptions) run() error {
	cfg := o.factory.Config
	if cfg == nil {
		cfg = &config.Config{Contexts: make(map[string]config.Context)}
	}

	cfg.Contexts[o.contextName] = config.Context{
		URL:      o.url,
		Token:    o.token,
		Org:      o.org,
		Env:      o.envID,
		ReadOnly: o.readOnly,
	}
	cfg.CurrentContext = o.contextName

	if err := cfg.SaveTo(o.factory.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(o.factory.IOStreams.Out, "Context '%s' saved and set as current.\n", o.contextName)

	return nil
}
