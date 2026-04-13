package user

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newUserRoleCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "role",
		Short: "Manage user roles",
	}

	cmd.PersistentFlags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkPersistentFlagRequired("user-id")

	cmd.AddCommand(newUserRoleListCmd(f, domainID, &userID))
	cmd.AddCommand(newUserRoleAssignCmd(f, domainID, &userID))
	cmd.AddCommand(newUserRoleRevokeCmd(f, domainID, &userID))

	return cmd
}

func newUserRoleListCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List user roles",
		Example: `  gio am user role list --domain my-domain --user-id user-1
  gio am user role list --domain my-domain --user-id user-1 -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUserRoleList(f, *domainID, *userID)
		},
	}
}

func runUserRoleList(f *factory.Factory, domainID, userID string) error {
	data, err := f.AM().ListUserRoles(domainID, userID)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}

type userRoleAssignOptions struct {
	factory  *factory.Factory
	domainID *string
	userID   *string
	roles    string
}

func newUserRoleAssignCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	opts := &userRoleAssignOptions{factory: f, domainID: domainID, userID: userID}

	cmd := &cobra.Command{
		Use:     "assign",
		Short:   "Assign roles to a user",
		Example: `  gio am user role assign --domain my-domain --user-id user-1 --roles role1,role2`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.roles, "roles", "", "Comma-separated list of role IDs (required)")
	_ = cmd.MarkFlagRequired("roles")

	return cmd
}

func (o *userRoleAssignOptions) run() error {
	roleList := strings.Split(o.roles, ",")
	body, _ := json.Marshal(roleList)

	if err := o.factory.AM().AssignUserRoles(*o.domainID, *o.userID, json.RawMessage(body)); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(o.factory)
	if err != nil {
		return err
	}

	p.PrintMessage("Roles assigned to user '%s'.", *o.userID)

	return nil
}

func newUserRoleRevokeCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "revoke <roleID>",
		Short:   "Revoke a role from a user",
		Example: `  gio am user role revoke role-1 --domain my-domain --user-id user-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RevokeUserRole(*domainID, *userID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Role '%s' revoked from user '%s'.", args[0], *userID)

			return nil
		},
	}
}
