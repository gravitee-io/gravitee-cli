package user

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newIdentityCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "identity",
		Short: "Manage user identities",
	}

	cmd.PersistentFlags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkPersistentFlagRequired("user-id")

	cmd.AddCommand(newIdentityListCmd(f, domainID, &userID))
	cmd.AddCommand(newIdentityUnlinkCmd(f, domainID, &userID))

	return cmd
}

func newIdentityListCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List user identities",
		Example: `  gio am user identity list --domain my-domain --user-id user-1
  gio am user identity list --domain my-domain --user-id user-1 -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runIdentityList(f, *domainID, *userID)
		},
	}
}

func runIdentityList(f *factory.Factory, domainID, userID string) error {
	items, err := f.AM().ListUserIdentities(domainID, userID)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(items)
	}

	return p.PrintList(items, identityColumns())
}

func newIdentityUnlinkCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "unlink <identityID>",
		Short:   "Unlink a user identity",
		Example: `  gio am user identity unlink identity-1 --domain my-domain --user-id user-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().UnlinkUserIdentity(*domainID, *userID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Identity '%s' unlinked.", args[0])

			return nil
		},
	}
}

func identityColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Provider ID", Value: func(i any) string { return cmdutil.StringField(i, "providerId") }},
		{Name: "User ID", Value: func(i any) string { return cmdutil.StringField(i, "userId") }},
	}
}
