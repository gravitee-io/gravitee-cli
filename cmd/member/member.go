package member

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewMemberCmd creates the parent member command with all subcommands.
func NewMemberCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "Manage API members",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newAddCmd(f))
	cmd.AddCommand(newRemoveCmd(f))

	return cmd
}
