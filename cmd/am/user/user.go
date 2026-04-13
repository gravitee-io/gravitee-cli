package user

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewUserCmd creates the user parent command with all user subcommands.
func NewUserCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "user",
		Aliases: []string{"users"},
		Short:   "Manage users",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))
	cmd.AddCommand(newCreateCmd(f, &domainID))
	cmd.AddCommand(newUpdateCmd(f, &domainID))
	cmd.AddCommand(newDeleteCmd(f, &domainID))
	cmd.AddCommand(newLockCmd(f, &domainID))
	cmd.AddCommand(newUnlockCmd(f, &domainID))
	cmd.AddCommand(newResetPasswordCmd(f, &domainID))
	cmd.AddCommand(newConsentCmd(f, &domainID))
	cmd.AddCommand(newUserRoleCmd(f, &domainID))
	cmd.AddCommand(newDeviceCmd(f, &domainID))
	cmd.AddCommand(newCredentialCmd(f, &domainID))
	cmd.AddCommand(newEnrolledFactorCmd(f, &domainID))
	cmd.AddCommand(newUserAuditCmd(f, &domainID))
	cmd.AddCommand(newSendRegistrationCmd(f, &domainID))
	cmd.AddCommand(newUpdateUsernameCmd(f, &domainID))
	cmd.AddCommand(newIdentityCmd(f, &domainID))
	cmd.AddCommand(newCertCredentialCmd(f, &domainID))
	cmd.AddCommand(newBulkCmd(f, &domainID))

	return cmd
}
