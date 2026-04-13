package group

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newMemberCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var groupID string

	cmd := &cobra.Command{
		Use:   "member",
		Short: "Manage group members",
	}

	cmd.PersistentFlags().StringVar(&groupID, "group-id", "", "Group ID (required)")
	_ = cmd.MarkPersistentFlagRequired("group-id")

	cmd.AddCommand(newMemberListCmd(f, domainID, &groupID))
	cmd.AddCommand(newMemberAddCmd(f, domainID, &groupID))
	cmd.AddCommand(newMemberRemoveCmd(f, domainID, &groupID))

	return cmd
}

func newMemberListCmd(f *factory.Factory, domainID, groupID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List group members",
		Example: `  gio am group member list --domain my-domain --group-id my-group`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListGroupMembers(*domainID, *groupID)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			return p.PrintDetail(data)
		},
	}
}

func newMemberAddCmd(f *factory.Factory, domainID, groupID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "add <userID>",
		Short:   "Add a member to a group",
		Example: `  gio am group member add user-123 --domain my-domain --group-id my-group`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().AddGroupMember(*domainID, *groupID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Member '%s' added to group '%s'.", args[0], *groupID)

			return nil
		},
	}
}

func newMemberRemoveCmd(f *factory.Factory, domainID, groupID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "remove <memberID>",
		Short:   "Remove a member from a group",
		Example: `  gio am group member remove member-123 --domain my-domain --group-id my-group`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RemoveGroupMember(*domainID, *groupID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Member '%s' removed from group '%s'.", args[0], *groupID)

			return nil
		},
	}
}
