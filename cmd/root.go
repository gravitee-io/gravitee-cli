package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	amcmd "github.com/gravitee-io/gio-cli/cmd/am"
	apimcmd "github.com/gravitee-io/gio-cli/cmd/apim"
	contextcmd "github.com/gravitee-io/gio-cli/cmd/context"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewRootCmd creates the root gio command.
func NewRootCmd(version string) *cobra.Command {
	f := &factory.Factory{
		IOStreams: factory.DefaultIOStreams(),
	}

	cmd := &cobra.Command{
		Use:           "gio",
		Short:         "gio - Gravitee CLI",
		Long:          "gio is a command-line interface for the Gravitee platform.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().StringVar(&f.Overrides.Context, "context", "", "Override context")
	cmd.PersistentFlags().StringVar(&f.Overrides.Org, "org", "", "Override organization ID")
	cmd.PersistentFlags().StringVar(&f.Overrides.EnvID, "env", "", "Override environment ID")
	cmd.PersistentFlags().BoolVar(&f.Debug, "debug", false, "Show raw HTTP requests/responses")

	cmd.AddCommand(newLoginCmd(f))
	cmd.AddCommand(contextcmd.NewContextCmd(f))
	cmd.AddCommand(newCompletionCmd())
	cmd.AddCommand(newVersionCmd(f, version))

	cmd.AddCommand(apimcmd.NewAPIMCmd(f))
	cmd.AddCommand(amcmd.NewAMCmd(f))

	return cmd
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
