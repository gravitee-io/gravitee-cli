package user

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type listOptions struct {
	factory  *factory.Factory
	domainID *string
	query    string
	filter   string
	page     int
	perPage  int
	all      bool
}

func newListCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &listOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Example: `  gio am user list --domain my-domain
  gio am user list --domain my-domain --query admin --per-page 20`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.ValidatePagination(opts.page, opts.perPage); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.query, "query", "", "Search query")
	cmd.Flags().StringVar(&opts.filter, "filter", "", "SCIM filter expression")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")

	return cmd
}

func (o *listOptions) run() error {
	f := o.factory
	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if o.all {
		return o.fetchAll(f, p)
	}

	return o.fetchPage(f, p, o.page)
}

func (o *listOptions) params(page int) am.ListUsersParams {
	return am.ListUsersParams{
		Query:   o.query,
		Filter:  o.filter,
		Page:    page,
		PerPage: o.perPage,
	}
}

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	resp, err := f.AM().ListUsers(*o.domainID, o.params(page-1)) // Convert 1-based CLI page to 0-based API page
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		raw, _ := json.Marshal(resp)
		return p.PrintDetail(json.RawMessage(raw))
	}

	if err := p.PrintList(resp.Data, userColumns()); err != nil {
		return err
	}

	if resp.TotalCount > len(resp.Data) {
		hint := " Use --all to fetch all results."
		if o.all {
			hint = ""
		}

		p.PrintMessage("Showing %d of %d.%s", len(resp.Data), resp.TotalCount, hint)
	} else if resp.TotalCount > 0 {
		p.PrintMessage("Showing %d results.", len(resp.Data))
	}

	return nil
}

func (o *listOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
	allData, err := am.FetchAllPages(func(page int) (*am.PaginatedResponse, error) {
		return f.AM().ListUsers(*o.domainID, o.params(page))
	}, o.perPage)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(allData)
	}

	if err := p.PrintList(allData, userColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintMessage("Showing %d results.", len(allData))
	}

	return nil
}

func userColumns() []printer.Column {
	return []printer.Column{
		{Name: "Username", Value: func(i any) string { return cmdutil.StringField(i, "username") }},
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Email", Value: func(i any) string { return cmdutil.StringField(i, "email") }},
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
	}
}
