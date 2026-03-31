package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	apicmd "github.com/gravitee-io/gio-cli/cmd/api"
	apikeycmd "github.com/gravitee-io/gio-cli/cmd/apikey"
	appcmd "github.com/gravitee-io/gio-cli/cmd/application"
	configcmd "github.com/gravitee-io/gio-cli/cmd/config"
	envcmd "github.com/gravitee-io/gio-cli/cmd/environment"
	membercmd "github.com/gravitee-io/gio-cli/cmd/member"
	metadatacmd "github.com/gravitee-io/gio-cli/cmd/metadata"
	pagecmd "github.com/gravitee-io/gio-cli/cmd/page"
	plancmd "github.com/gravitee-io/gio-cli/cmd/plan"
	plugincmd "github.com/gravitee-io/gio-cli/cmd/plugin"
	subcmd "github.com/gravitee-io/gio-cli/cmd/subscription"
	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewRootCmd creates the root gio command with all subcommands.
func NewRootCmd(version string) *cobra.Command {
	var overrides config.Overrides
	var debug bool

	f := &factory.Factory{
		IOStreams: factory.DefaultIOStreams(),
	}

	cmd := &cobra.Command{
		Use:     "gio",
		Short:   "gio - Gravitee APIM CLI",
		Long:    "gio is a command-line interface for Gravitee API Management.",
		Version: version,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if f.ConfigPath == "" {
				p, err := config.Path()
				if err != nil {
					return err
				}

				f.ConfigPath = p
			}

			cfg, err := config.LoadFrom(f.ConfigPath)
			if err != nil {
				return err
			}

			f.Config = cfg

			resolved, err := cfg.Resolve(overrides)
			if err != nil {
				// Allow commands that don't need a context (login, config, version).
				return nil //nolint:nilerr // Context resolution failure is not fatal for all commands.
			}

			f.Resolved = resolved
			f.Client = client.NewHTTPClient(client.HTTPClientConfig{
				BaseURL:  resolved.URL,
				Token:    resolved.Token,
				Debug:    debug,
				DebugOut: f.IOStreams.Err,
			})

			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags.
	cmd.PersistentFlags().StringVar(&overrides.Context, "context", "", "Override context")
	cmd.PersistentFlags().StringVar(&overrides.Org, "org", "", "Override organization ID")
	cmd.PersistentFlags().StringVar(&overrides.EnvID, "env-id", "", "Override environment ID")
	cmd.PersistentFlags().StringVarP(&f.OutputFormat, "output", "o", "table", "Output format: table, json, yaml")
	cmd.PersistentFlags().BoolVarP(&f.Quiet, "quiet", "q", false, "Suppress output except errors")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Show raw HTTP requests/responses")

	// Subcommands.
	cmd.AddCommand(newLoginCmd(f))
	cmd.AddCommand(configcmd.NewConfigCmd(f))
	cmd.AddCommand(apicmd.NewAPICmd(f))
	cmd.AddCommand(plancmd.NewPlanCmd(f))
	cmd.AddCommand(subcmd.NewSubscriptionCmd(f))
	cmd.AddCommand(apikeycmd.NewAPIKeyCmd(f))
	cmd.AddCommand(membercmd.NewMemberCmd(f))
	cmd.AddCommand(pagecmd.NewPageCmd(f))
	cmd.AddCommand(metadatacmd.NewMetadataCmd(f))
	cmd.AddCommand(appcmd.NewApplicationCmd(f))
	cmd.AddCommand(envcmd.NewEnvironmentCmd(f))
	cmd.AddCommand(plugincmd.NewPluginCmd(f))

	// Customize help template to show context info.
	cmd.SetHelpFunc(contextualHelp(cmd.HelpFunc()))

	return cmd
}

func contextualHelp(defaultHelp func(*cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		// Only show context info on the root command.
		if cmd.Parent() == nil || cmd.Name() == "gio" {
			cfg, err := config.Load()
			if err == nil {
				printContextHeader(cmd, cfg)
			}
		}

		defaultHelp(cmd, args)
	}
}

func printContextHeader(cmd *cobra.Command, cfg *config.Config) {
	if cfg.CurrentContext == "" || len(cfg.Contexts) == 0 {
		fmt.Fprint(cmd.OutOrStderr(), "\n  No context configured. Run 'gio login' to get started.\n\n")

		return
	}

	ctx, ok := cfg.Contexts[cfg.CurrentContext]
	if !ok {
		return
	}

	env := ctx.Env
	if env == "" {
		env = config.DefaultEnv
	}

	mode := ""
	if ctx.ReadOnly {
		mode = "\n  Mode:      read-only"
	}

	fmt.Fprintf(cmd.OutOrStderr(), "\n  Context:   %s\n  URL:       %s\n  Env:       %s%s\n\n",
		cfg.CurrentContext, ctx.URL, env, mode)
}

// Execute runs the root command.
func Execute(version string) int {
	cmd := NewRootCmd(version)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "Error:", err)

		return 1
	}

	return 0
}
