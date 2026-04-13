package user

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

type updateUsernameOptions struct {
	factory  *factory.Factory
	domainID *string
	username string
}

func newUpdateUsernameCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &updateUsernameOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:     "update-username <userID>",
		Short:   "Update a user's username",
		Example: `  gio am user update-username user-1 --domain my-domain --username newname`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run(args[0])
		},
	}

	cmd.Flags().StringVar(&opts.username, "username", "", "New username (required)")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

func (o *updateUsernameOptions) run(userID string) error {
	f := o.factory

	body := map[string]any{"username": o.username}
	raw, _ := json.Marshal(body)

	if _, err := f.AM().UpdateUsername(*o.domainID, userID, json.RawMessage(raw)); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	p.PrintMessage("Username updated for user '%s'.", userID)

	return nil
}
