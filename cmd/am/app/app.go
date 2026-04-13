package app

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAppCmd creates the application parent command with all application subcommands.
func NewAppCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "app",
		Aliases: []string{"application"},
		Short:   "Manage applications",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))
	cmd.AddCommand(newCreateCmd(f, &domainID))
	cmd.AddCommand(newUpdateCmd(f, &domainID))
	cmd.AddCommand(newDeleteCmd(f, &domainID))
	cmd.AddCommand(newSecretCmd(f, &domainID))
	cmd.AddCommand(newAppMemberCmd(f, &domainID))
	cmd.AddCommand(newAppFlowCmd(f, &domainID))
	cmd.AddCommand(newAppEmailCmd(f, &domainID))
	cmd.AddCommand(newAppFormCmd(f, &domainID))
	cmd.AddCommand(newAppResourceCmd(f, &domainID))
	cmd.AddCommand(newAppAnalyticsCmd(f, &domainID))
	cmd.AddCommand(newChangeTypeCmd(f, &domainID))
	cmd.AddCommand(newAppResourcePolicyCmd(f, &domainID))

	return cmd
}
