package page

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewPageCmd creates the parent page command with all subcommands.
func NewPageCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "page",
		Short: "Manage API pages",
		Args:  cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(newDeleteCmd(f))
	cmd.AddCommand(newPublishCmd(f))
	cmd.AddCommand(newUnpublishCmd(f))

	return cmd
}
