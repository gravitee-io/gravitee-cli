package member

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newAddCmd(f *factory.Factory) *cobra.Command {
	var (
		apiID  string
		userID string
		role   string
	)

	cmd := &cobra.Command{
		Use:     "add --api <apiId> --user <userId> --role <role>",
		Short:   "Add a user as a member of an API with the specified role",
		Example: `  gio apim member add --api /my/api --user bbbb1111-2222-3333-4444-555566667777 --role OWNER`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runAdd(f, apiID, userID, role)
		},
	}

	cmdutil.AddAPIFlag(cmd, &apiID)
	cmd.Flags().StringVar(&userID, "user", "", "User ID to add as member (required)")
	cmd.Flags().StringVar(&role, "role", "", "Role to assign to the member (required)")
	_ = cmd.MarkFlagRequired("user")
	_ = cmd.MarkFlagRequired("role")

	return cmd
}

func runAdd(f *factory.Factory, apiID, userID, role string) error {
	data, err := f.APIM().AddMember(apiID, userID, role)
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

	return printMemberDetail(p, data)
}

func printMemberDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if v, ok := m["displayName"]; ok && v != nil {
		p.PrintMessage("%-16s%v", "Display Name:", v)
	}

	if v, ok := m["id"]; ok && v != nil {
		p.PrintMessage("%-16s%v", "ID:", v)
	}

	p.PrintMessage("%-16s%s", "Role:", roleFromMap(m))

	if v, ok := m["type"]; ok && v != nil {
		p.PrintMessage("%-16s%v", "Type:", v)
	}

	return nil
}

func roleFromMap(m map[string]any) string {
	roles, ok := m["roles"].([]any)
	if !ok || len(roles) == 0 {
		return ""
	}

	first, ok := roles[0].(map[string]any)
	if !ok {
		return ""
	}

	name, _ := first["name"].(string)

	return name
}
