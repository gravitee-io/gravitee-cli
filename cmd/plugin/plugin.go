package plugin

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewPluginCmd creates the parent plugin command with all subcommands.
func NewPluginCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(newListCmd(f))

	return cmd
}
