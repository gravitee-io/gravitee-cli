package alert

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAlertCmd creates the alert parent command with notifier and trigger subcommands.
func NewAlertCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "alert",
		Aliases: []string{"alerts"},
		Short:   "Manage alerts (notifiers and triggers)",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newNotifierCmd(f, &domainID))
	cmd.AddCommand(newTriggerCmd(f, &domainID))

	return cmd
}
