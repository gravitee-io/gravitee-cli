package metadata

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewMetadataCmd creates the parent metadata command with all subcommands.
func NewMetadataCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Manage API metadata",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(newDeleteCmd(f))

	return cmd
}
