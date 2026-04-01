package am

import (
	"fmt"

	"github.com/spf13/cobra"

	domaincmd "github.com/gravitee-io/gio-cli/cmd/am/domain"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewAMCmd creates the am parent command with all AM subcommands.
func NewAMCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "am",
		Short: "Gravitee Access Management",
		Long:  "Manage Gravitee AM resources: domains, applications, users, identity providers, and more.",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.SetupConfig(f); err != nil {
				return err
			}
			return cmdutil.ResolveProductContext(f, "am")
		},
	}

	// Override help to show context info.
	defaultHelp := cmd.HelpFunc()
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		_ = cmdutil.SetupConfig(f)
		_ = cmdutil.ResolveProductContext(f, "am")
		if header := cmdutil.ContextHeader(f, "am"); header != "" {
			fmt.Fprint(c.OutOrStdout(), header+"\n")
		}

		defaultHelp(c, args)
	})

	cmd.AddCommand(domaincmd.NewDomainCmd(f))

	return cmd
}
