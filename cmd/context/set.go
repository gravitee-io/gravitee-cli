package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type setOptions struct {
	factory  *factory.Factory
	org      string
	envID    string
	readOnly bool
}

func newSetCmd(f *factory.Factory) *cobra.Command {
	opts := &setOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "set <name>",
		Short: "Create or update shared fields of a context",
		Long:  "Create or update shared fields (org, env, read-only) of a context. Use 'gio login' to set product URLs and tokens.",
		Example: `  gio context set prod --org ACME --env production
  gio context set staging --read-only
  gio context set dev --org DEV --env development --read-only=false`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.run(cmd, args[0])
		},
	}

	cmd.Flags().StringVar(&opts.org, "org", "", "Organization ID")
	cmd.Flags().StringVar(&opts.envID, "env-id", "", "Environment ID")
	cmd.Flags().BoolVar(&opts.readOnly, "read-only", false, "Enable read-only mode")

	return cmd
}

func (o *setOptions) run(cmd *cobra.Command, name string) error {
	if err := cmdutil.SetupConfig(o.factory); err != nil {
		return err
	}

	cfg := o.factory.Config
	ctx := cfg.EnsureContext(name)

	if cmd.Flags().Changed("org") {
		ctx.Org = o.org
	}

	if cmd.Flags().Changed("env-id") {
		ctx.Env = o.envID
	}

	if cmd.Flags().Changed("read-only") {
		ctx.ReadOnly = o.readOnly
	}

	if err := cfg.SaveTo(o.factory.ConfigPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(o.factory.IOStreams.Out, "Context '%s' updated.\n", name)

	return nil
}
