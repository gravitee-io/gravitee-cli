package passwordpolicy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newListCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List password policies",
		Example: `  gio am password-policy list --domain my-domain
  gio am pp list --domain my-domain`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runList(f, *domainID)
		},
	}
}

func runList(f *factory.Factory, domainID string) error {
	items, err := f.AM().ListPasswordPolicies(domainID)
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

	return p.PrintList(items, policyColumns())
}

func policyColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "MinLength", Value: func(i any) string {
			m, ok := i.(map[string]any)
			if !ok {
				return ""
			}

			if v, ok := m["minLength"]; ok && v != nil {
				return fmt.Sprintf("%v", v)
			}

			return ""
		}},
		{Name: "MaxLength", Value: func(i any) string {
			m, ok := i.(map[string]any)
			if !ok {
				return ""
			}

			if v, ok := m["maxLength"]; ok && v != nil {
				return fmt.Sprintf("%v", v)
			}

			return ""
		}},
	}
}
