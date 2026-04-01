package domain

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewDomainCmd creates the domain parent command with all domain subcommands.
func NewDomainCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain",
		Aliases: []string{"dom"},
		Short:   "Manage security domains",
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))

	return cmd
}
