package certificate

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewCertificateCmd creates the certificate parent command with all certificate subcommands.
func NewCertificateCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "certificate",
		Aliases: []string{"cert", "certificates"},
		Short:   "Manage certificates",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))
	cmd.AddCommand(newCreateCmd(f, &domainID))
	cmd.AddCommand(newUpdateCmd(f, &domainID))
	cmd.AddCommand(newDeleteCmd(f, &domainID))
	cmd.AddCommand(newKeyCmd(f, &domainID))
	cmd.AddCommand(newKeysCmd(f, &domainID))
	cmd.AddCommand(newRotateCmd(f, &domainID))

	return cmd
}
