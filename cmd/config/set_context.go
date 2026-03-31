package config

import (
	"fmt"

	"github.com/spf13/cobra"

	iconfig "github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type setContextOptions struct {
	factory     *factory.Factory
	contextName string
	url         string
	token       string
	org         string
	envID       string
	readOnly    bool
}

func newSetContextCmd(f *factory.Factory) *cobra.Command {
	opts := &setContextOptions{factory: f}

	cmd := &cobra.Command{
		Use:     "set-context <name>",
		Short:   "Create or update a context",
		Example: `  gio config set-context prod --url https://apim-prod.company.com --token gioat_xyz --read-only`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			opts.contextName = args[0]

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.url, "url", "", "URL of the Gravitee control plane (required)")
	cmd.Flags().StringVar(&opts.token, "token", "", "Personal Access Token (required)")
	cmd.Flags().StringVar(&opts.org, "org", iconfig.DefaultOrg, "Organization ID")
	cmd.Flags().StringVar(&opts.envID, "env-id", iconfig.DefaultEnv, "Environment ID")
	cmd.Flags().BoolVar(&opts.readOnly, "read-only", false, "Enable read-only mode")
	_ = cmd.MarkFlagRequired("url")
	_ = cmd.MarkFlagRequired("token")

	return cmd
}

func (o *setContextOptions) run() error {
	cfg := o.factory.Config

	cfg.Contexts[o.contextName] = iconfig.Context{
		URL:      o.url,
		Token:    o.token,
		Org:      o.org,
		Env:      o.envID,
		ReadOnly: o.readOnly,
	}

	if err := cfg.SaveTo(o.factory.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(o.factory.IOStreams.Out, "Context '%s' saved.\n", o.contextName)

	return nil
}
