package plan

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewPlanCmd creates the parent plan command with all subcommands.
func NewPlanCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Manage API plans",
		Args:  cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(newDeleteCmd(f))
	cmd.AddCommand(newPublishCmd(f))
	cmd.AddCommand(newDeprecateCmd(f))
	cmd.AddCommand(newCloseCmd(f))

	return cmd
}
