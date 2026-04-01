package api

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAPICmd creates the parent api command with all subcommands.
func NewAPICmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Manage APIs",
		Args:  cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newUpdateCmd(f))
	cmd.AddCommand(newDeleteCmd(f))
	cmd.AddCommand(newStartCmd(f))
	cmd.AddCommand(newStopCmd(f))
	cmd.AddCommand(newDeployCmd(f))
	cmd.AddCommand(newImportCmd(f))
	cmd.AddCommand(newExportCmd(f))
	cmd.AddCommand(newRollbackCmd(f))
	cmd.AddCommand(newAnalyticsCmd(f))
	cmd.AddCommand(newHealthCmd(f))
	cmd.AddCommand(newLogsCmd(f))
	cmd.AddCommand(newLogCmd(f))

	return cmd
}
