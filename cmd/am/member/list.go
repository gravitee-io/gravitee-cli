package member

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newListCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List domain members",
		Example: `  gio am member list --domain my-domain`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runList(f, *domainID)
		},
	}
}

func runList(f *factory.Factory, domainID string) error {
	data, err := f.AM().ListMembers(domainID)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(data)
	}

	// The response has a "memberships" array - extract it for table display.
	memberships, err := extractMemberships(data)
	if err != nil {
		return err
	}

	return p.PrintList(memberships, memberColumns())
}

func extractMemberships(data json.RawMessage) ([]json.RawMessage, error) {
	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse member response: %w", err)
	}

	raw, ok := wrapper["memberships"]
	if !ok {
		return nil, nil
	}

	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("failed to parse memberships array: %w", err)
	}

	return items, nil
}

func memberColumns() []printer.Column {
	return []printer.Column{
		{Name: "MemberID", Value: func(i any) string { return cmdutil.StringField(i, "memberId") }},
		{Name: "RoleID", Value: func(i any) string { return cmdutil.StringField(i, "roleId") }},
		{Name: "MemberType", Value: func(i any) string { return cmdutil.StringField(i, "memberType") }},
	}
}
