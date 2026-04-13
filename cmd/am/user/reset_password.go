package user

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type resetPasswordOptions struct {
	factory  *factory.Factory
	domainID *string
	password string
}

func newResetPasswordCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &resetPasswordOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:     "reset-password <userID> --password <newPassword>",
		Short:   "Reset a user's password",
		Example: `  gio am user reset-password user-id --domain my-domain --password newSecret123`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run(args[0])
		},
	}

	cmd.Flags().StringVar(&opts.password, "password", "", "New password (required)")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}

func (o *resetPasswordOptions) run(userID string) error {
	f := o.factory

	body := map[string]any{"password": o.password}
	raw, _ := json.Marshal(body)

	if err := f.AM().ResetPassword(*o.domainID, userID, json.RawMessage(raw)); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Password reset for user '%s'.", userID)

	return nil
}
