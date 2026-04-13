package dictionary

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewDictionaryCmd creates the dictionary parent command with all subcommands.
func NewDictionaryCmd(f *factory.Factory) *cobra.Command {
	var domainID string

	cmd := &cobra.Command{
		Use:     "dictionary",
		Aliases: []string{"dictionaries", "i18n"},
		Short:   "Manage i18n dictionaries",
	}

	cmd.PersistentFlags().StringVar(&domainID, "domain", "", "Domain ID (required)")
	_ = cmd.MarkPersistentFlagRequired("domain")

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f, &domainID))
	cmd.AddCommand(newGetCmd(f, &domainID))
	cmd.AddCommand(newCreateCmd(f, &domainID))
	cmd.AddCommand(newUpdateCmd(f, &domainID))
	cmd.AddCommand(newDeleteCmd(f, &domainID))
	cmd.AddCommand(newEntryCmd(f, &domainID))

	return cmd
}
