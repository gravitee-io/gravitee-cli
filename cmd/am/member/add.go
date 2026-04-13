package member

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newAddCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var (
		memberID   string
		role       string
		memberType string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a member to a domain",
		Example: `  gio am member add --domain my-domain --member-id user-123 --role role-456
  gio am member add --domain my-domain --member-id user-123 --role role-456 --member-type USER`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runAdd(f, *domainID, memberID, role, memberType)
		},
	}

	cmd.Flags().StringVar(&memberID, "member-id", "", "Member user ID (required)")
	cmd.Flags().StringVar(&role, "role", "", "Role ID (required)")
	cmd.Flags().StringVar(&memberType, "member-type", "USER", "Member type")
	_ = cmd.MarkFlagRequired("member-id")
	_ = cmd.MarkFlagRequired("role")

	return cmd
}

func runAdd(f *factory.Factory, domainID, memberID, role, memberType string) error {
	body, _ := json.Marshal(map[string]string{
		"memberId":   memberID,
		"memberType": memberType,
		"role":       role,
	})

	data, err := f.AM().AddMember(domainID, body)
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

	p.PrintMessage("Member '%s' added.", memberID)

	return nil
}
