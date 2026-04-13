package passwordpolicy

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewPasswordPolicyCmd creates the password-policy parent command with all password policy subcommands.
func NewPasswordPolicyCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "password-policy",
		Aliases: []string{"pp", "password-policies"},
		Short:   "Manage password policies",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))
	cmd.AddCommand(newCreateCmd(f, &domainID))
	cmd.AddCommand(newUpdateCmd(f, &domainID))
	cmd.AddCommand(newDeleteCmd(f, &domainID))
	cmd.AddCommand(newActiveCmd(f, &domainID))
	cmd.AddCommand(newSetDefaultCmd(f, &domainID))
	cmd.AddCommand(newEvaluateCmd(f, &domainID))

	return cmd
}
