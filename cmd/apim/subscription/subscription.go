package subscription

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewSubscriptionCmd creates the parent subscription command with all subcommands.
func NewSubscriptionCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subscription",
		Aliases: []string{"sub"},
		Short:   "Manage subscriptions",
		Args:    cobra.NoArgs,
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newGetCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newAcceptCmd(f))
	cmd.AddCommand(newRejectCmd(f))
	cmd.AddCommand(newPauseCmd(f))
	cmd.AddCommand(newResumeCmd(f))
	cmd.AddCommand(newCloseCmd(f))
	cmd.AddCommand(newTransferCmd(f))

	return cmd
}
