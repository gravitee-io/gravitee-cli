package authorizationengine

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAuthorizationEngineCmd creates the authorization engine parent command with all subcommands.
func NewAuthorizationEngineCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "authorization-engine",
		Aliases: []string{"authz-engine", "ae"},
		Short:   "Manage authorization engines",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))
	cmd.AddCommand(newUpdateCmd(f, &domainID))

	return cmd
}
