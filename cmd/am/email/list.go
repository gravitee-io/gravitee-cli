package email

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newListCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List emails",
		Example: `  gio am email list --domain my-domain
  gio am emails list --domain my-domain`,
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
	items, err := f.AM().ListEmails(domainID)
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

	return p.PrintList(items, emailColumns())
}

func emailColumns() []printer.Column {
	return []printer.Column{
		{Name: "Template", Value: func(i any) string { return cmdutil.StringField(i, "template") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Enabled", Value: func(i any) string {
			m, ok := i.(map[string]any)
			if !ok {
				return ""
			}

			if v, ok := m["enabled"].(bool); ok && v {
				return "true"
			}

			return "false"
		}},
		{Name: "Subject", Value: func(i any) string { return cmdutil.StringField(i, "subject") }},
	}
}
