package org

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewOrgCmd creates the org parent command with all organization subcommands.
func NewOrgCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "org",
		Aliases: []string{"organization"},
		Short:   "Manage organization resources",
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newOrgUserCmd(f))
	cmd.AddCommand(newOrgGroupCmd(f))
	cmd.AddCommand(newOrgRoleCmd(f))
	cmd.AddCommand(newOrgSettingsCmd(f))
	cmd.AddCommand(newOrgMemberCmd(f))
	cmd.AddCommand(newOrgAuditCmd(f))
	cmd.AddCommand(newOrgReporterCmd(f))
	cmd.AddCommand(newOrgFormCmd(f))
	cmd.AddCommand(newOrgIDPCmd(f))
	cmd.AddCommand(newOrgEntrypointCmd(f))
	cmd.AddCommand(newOrgTagCmd(f))
	cmd.AddCommand(newOrgUserTokenCmd(f))

	return cmd
}
