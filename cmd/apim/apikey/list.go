package apikey

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/apim"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type listOptions struct {
	factory      *factory.Factory
	apiID        string
	subscription string
	page         int
	perPage      int
	all          bool
}

func newListCmd(f *factory.Factory) *cobra.Command {
	opts := &listOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "list --api <apiId> --subscription <subId>",
		Short: "List API keys for a subscription",
		Example: `  gio apim api-key list --api 8a7b3c4d --subscription aaaa1111
  gio apim api-key list --api 8a7b3c4d --subscription aaaa1111 --all`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.apiID, "api", "", "API ID (required)")
	cmd.Flags().StringVar(&opts.subscription, "subscription", "", "Subscription ID (required)")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
	cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")

	_ = cmd.MarkFlagRequired("api")
	_ = cmd.MarkFlagRequired("subscription")

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

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	resp, err := f.APIM().ListAPIKeys(o.apiID, o.subscription, page, o.perPage)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(resp)
	}

	if err := p.PrintList(resp.Data, apiKeyColumns()); err != nil {
		return err
	}

	pg := resp.Pagination
	cmdutil.PrintPaginationHint(p, pg.Page, pg.PerPage, pg.PageCount, pg.TotalCount, pg.PageItemsCount, o.all)

	return nil
}

func (o *listOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
	allData, err := apim.FetchAllPages(func(page int) (*apim.PaginatedResponse, error) {
		return f.APIM().ListAPIKeys(o.apiID, o.subscription, page, o.perPage)
	})
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(allData)
	}

	if err := p.PrintList(allData, apiKeyColumns()); err != nil {
		return err
	}

	if len(allData) > 0 {
		p.PrintMessage("Showing %d results.", len(allData))
	}

	return nil
}

func apiKeyColumns() []printer.Column {
	return []printer.Column{
		{Name: "Key", Value: func(i any) string { return cmdutil.StringField(i, "key") }},
		{Name: "Revoked", Value: func(i any) string { return boolField(i, "revoked") }},
		{Name: "Expired", Value: func(i any) string { return boolField(i, "expired") }},
		{Name: "Created", Value: func(i any) string { return cmdutil.StringField(i, "createdAt") }},
		{Name: "Expire At", Value: func(i any) string {
			s := cmdutil.StringField(i, "expireAt")
			if s == "" {
				return "-"
			}

			return s
		}},
	}
}
