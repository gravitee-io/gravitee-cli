package token

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
	"github.com/spf13/cobra"
)

func newListCmd(f *factory.Factory) *cobra.Command {
	var userID string
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List user tokens",
		Example: `  gio am token list --user user-uuid`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			return runList(f, userID)
		},
	}
	cmd.Flags().StringVar(&userID, "user", "", "User ID (required)")
	_ = cmd.MarkFlagRequired("user")
	return cmd
}

func runList(f *factory.Factory, userID string) error {
	path := cmdutil.AMDomainPath(f, fmt.Sprintf("users/%s/tokens", userID))
	data, err := f.Client.Get(path)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}
	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}
	return p.PrintList(json.RawMessage(data), tokenColumns())
}

func tokenColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
		{Name: "Token", Value: func(i interface{}) string { return cmdutil.StringField(i, "token") }},
	}
}
