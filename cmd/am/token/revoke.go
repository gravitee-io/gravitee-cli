package token

import (
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newRevokeCmd(f *factory.Factory) *cobra.Command {
	var userID string
	cmd := &cobra.Command{
		Use:     "revoke <tokenId>",
		Short:   "Revoke a user token",
		Example: `  gio am token revoke token-id --user user-uuid`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			return runRevoke(f, userID, args[0])
		},
	}
	cmd.Flags().StringVar(&userID, "user", "", "User ID (required)")
	_ = cmd.MarkFlagRequired("user")
	return cmd
}

func runRevoke(f *factory.Factory, userID, tokenID string) error {
	path := cmdutil.AMDomainPath(f, fmt.Sprintf("users/%s/tokens/%s", userID, tokenID))
	if err := f.Client.Delete(path); err != nil {
		return err
	}
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}
	p.PrintMessage("Token '%s' revoked.", tokenID)
	return nil
}
