package token

import (
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

// NewTokenCmd creates the parent "gio am token" command.
func NewTokenCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage user tokens",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newRevokeCmd(f))
	return cmd
}
