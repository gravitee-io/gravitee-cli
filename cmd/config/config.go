package config

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewConfigCmd creates the parent config command with all subcommands.
func NewConfigCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(newSetContextCmd(f))
	cmd.AddCommand(newUseContextCmd(f))
	cmd.AddCommand(newCurrentContextCmd(f))
	cmd.AddCommand(newGetContextsCmd(f))
	cmd.AddCommand(newViewCmd(f))

	return cmd
}
