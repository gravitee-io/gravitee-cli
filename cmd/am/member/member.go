package member

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewMemberCmd creates the member parent command with all member subcommands.
func NewMemberCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "member",
		Aliases: []string{"members"},
		Short:   "Manage domain members",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newAddCmd(f, &domainID))
	cmd.AddCommand(newRemoveCmd(f, &domainID))
	cmd.AddCommand(newPermissionsCmd(f, &domainID))

	return cmd
}
