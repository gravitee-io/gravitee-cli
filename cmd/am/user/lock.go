package user

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newLockCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "lock <userID>",
		Short:   "Lock a user account",
		Example: `  gio am user lock user-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdateStatus(f, *domainID, args[0], false)
		},
	}
}

func newUnlockCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "unlock <userID>",
		Short:   "Unlock a user account",
		Example: `  gio am user unlock user-id --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdateStatus(f, *domainID, args[0], true)
		},
	}
}

func runUpdateStatus(f *factory.Factory, domainID, userID string, enabled bool) error {
	body := map[string]any{"enabled": enabled}
	raw, _ := json.Marshal(body)

	if _, err := f.AM().UpdateUserStatus(domainID, userID, json.RawMessage(raw)); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	action := "unlocked"
	if !enabled {
		action = "locked"
	}

	p.PrintMessage("User '%s' %s.", userID, action)

	return nil
}
