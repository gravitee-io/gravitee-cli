package context

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewContextCmd creates the context parent command with all subcommands.
func NewContextCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage CLI contexts",
		Long:  "Create, switch, and inspect CLI contexts that store connection details for Gravitee products.",
	}

	cmd.AddCommand(newUseCmd(f))
	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newCurrentCmd(f))
	cmd.AddCommand(newViewCmd(f))
	cmd.AddCommand(newSetCmd(f))
	cmd.AddCommand(newDeleteCmd(f))

	return cmd
}
