package protectedresource

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewProtectedResourceCmd creates the protected resource parent command with all subcommands.
func NewProtectedResourceCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "protected-resource",
		Aliases: []string{"protected-resources", "pr"},
		Short:   "Manage protected resources",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))
	cmd.AddCommand(newCreateCmd(f, &domainID))
	cmd.AddCommand(newUpdateCmd(f, &domainID))
	cmd.AddCommand(newDeleteCmd(f, &domainID))
	cmd.AddCommand(newMemberCmd(f, &domainID))
	cmd.AddCommand(newSecretCmd(f, &domainID))

	return cmd
}
