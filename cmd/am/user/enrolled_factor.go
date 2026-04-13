package user

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newEnrolledFactorCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "factor",
		Short: "Manage user enrolled factors",
	}

	cmd.PersistentFlags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkPersistentFlagRequired("user-id")

	cmd.AddCommand(newEnrolledFactorListCmd(f, domainID, &userID))
	cmd.AddCommand(newEnrolledFactorDeleteCmd(f, domainID, &userID))

	return cmd
}

func newEnrolledFactorListCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List user enrolled factors",
		Example: `  gio am user factor list --domain my-domain --user-id user-1
  gio am user factor list --domain my-domain --user-id user-1 -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runEnrolledFactorList(f, *domainID, *userID)
		},
	}
}

func runEnrolledFactorList(f *factory.Factory, domainID, userID string) error {
	items, err := f.AM().ListUserFactors(domainID, userID)
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

	return p.PrintList(items, enrolledFactorColumns())
}

func newEnrolledFactorDeleteCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <factorID>",
		Short:   "Delete a user enrolled factor",
		Example: `  gio am user factor delete factor-1 --domain my-domain --user-id user-1`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := f.AM().DeleteUserFactor(*domainID, *userID, args[0]); err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			p.PrintMessage("Factor '%s' deleted.", args[0])

			return nil
		},
	}
}

func enrolledFactorColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Factor ID", Value: func(i any) string { return cmdutil.StringField(i, "factorId") }},
		{Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "status") }},
	}
}
