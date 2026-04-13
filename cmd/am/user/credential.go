package user

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newCredentialCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "credential",
		Short: "Manage user credentials",
	}

	cmd.PersistentFlags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkPersistentFlagRequired("user-id")

	cmd.AddCommand(newCredentialListCmd(f, domainID, &userID))
	cmd.AddCommand(newCredentialGetCmd(f, domainID, &userID))
	cmd.AddCommand(newCredentialRevokeCmd(f, domainID, &userID))

	return cmd
}

func newCredentialListCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List user credentials",
		Example: `  gio am user credential list --domain my-domain --user-id user-1
  gio am user credential list --domain my-domain --user-id user-1 -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runCredentialList(f, *domainID, *userID)
		},
	}
}

func runCredentialList(f *factory.Factory, domainID, userID string) error {
	items, err := f.AM().ListUserCredentials(domainID, userID)
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

	return p.PrintList(items, credentialColumns())
}

func newCredentialGetCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "get <credentialID>",
		Short: "Get user credential details",
		Example: `  gio am user credential get cred-1 --domain my-domain --user-id user-1
  gio am user credential get cred-1 --domain my-domain --user-id user-1 -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetUserCredential(*domainID, *userID, args[0])
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

func newCredentialRevokeCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "revoke <credentialID>",
		Short:   "Revoke a user credential",
		Example: `  gio am user credential revoke cred-1 --domain my-domain --user-id user-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().RevokeUserCredential(*domainID, *userID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Credential '%s' revoked.", args[0])

			return nil
		},
	}
}

func credentialColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Credential ID", Value: func(i any) string { return cmdutil.StringField(i, "credentialId") }},
		{Name: "Type", Value: func(i any) string { return cmdutil.StringField(i, "type") }},
	}
}
