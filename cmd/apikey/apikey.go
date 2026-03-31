package apikey

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAPIKeyCmd creates the parent api-key command with all subcommands.
func NewAPIKeyCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-key",
		Short: "Manage API keys",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newRenewCmd(f))
	cmd.AddCommand(newRevokeCmd(f))
	cmd.AddCommand(newReactivateCmd(f))

	return cmd
}
