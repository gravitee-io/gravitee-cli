package apikey

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

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
		Example: `  gio api-key list --api 8a7b3c4d --subscription aaaa1111
  gio api-key list --api 8a7b3c4d --subscription aaaa1111 --all`,
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

type paginatedResponse struct {
	Data       []json.RawMessage `json:"data"`
	Pagination struct {
		Page           int `json:"page"`
		PerPage        int `json:"perPage"`
		PageCount      int `json:"pageCount"`
		TotalCount     int `json:"totalCount"`
		PageItemsCount int `json:"pageItemsCount"`
	} `json:"pagination"`
}

func (o *listOptions) run() error {
	f := o.factory
	p := cmdutil.NewPrinter(f)

	if o.all {
		return o.fetchAll(f, p)
	}

	return o.fetchPage(f, p, o.page)
}

func (o *listOptions) buildQuery(page int) string {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("perPage", strconv.Itoa(o.perPage))

	return q.Encode()
}

func (o *listOptions) basePath() string {
	return fmt.Sprintf("apis/%s/subscriptions/%s/api-keys", o.apiID, o.subscription)
}

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
	path := cmdutil.V2EnvPath(f, o.basePath()+"?"+o.buildQuery(page))

	data, err := f.Client.Get(path)
	if err != nil {
		return fmt.Errorf("API key list failed: %w", err)
	}

	var resp paginatedResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	if err := p.PrintList(resp.Data, apiKeyColumns()); err != nil {
		return err
	}

	pg := resp.Pagination
	start := (pg.Page-1)*pg.PerPage + 1
	end := start + pg.PageItemsCount - 1

	if pg.PageCount > 1 {
		hint := " Use --all to fetch all results."
		if o.all || pg.Page == pg.PageCount {
			hint = ""
		}

		p.PrintMessage("Showing %d-%d of %d (page %d/%d).%s",
			start, end, pg.TotalCount, pg.Page, pg.PageCount, hint)
	} else if pg.TotalCount > 0 {
		p.PrintMessage("Showing %d-%d of %d (page %d/%d).",
			start, end, pg.TotalCount, pg.Page, pg.PageCount)
	}

	return nil
}

func (o *listOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
	var allData []json.RawMessage

	for page := 1; ; page++ {
		path := cmdutil.V2EnvPath(f, o.basePath()+"?"+o.buildQuery(page))

		data, err := f.Client.Get(path)
		if err != nil {
			return fmt.Errorf("API key list failed: %w", err)
		}

		var resp paginatedResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		allData = append(allData, resp.Data...)

		if resp.Pagination.PageCount <= 0 || page >= resp.Pagination.PageCount || page > 1000 {
			break
		}
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
		{Name: "Key", Value: func(i interface{}) string { return cmdutil.StringField(i, "key") }},
		{Name: "Revoked", Value: func(i interface{}) string { return boolField(i, "revoked") }},
		{Name: "Expired", Value: func(i interface{}) string { return boolField(i, "expired") }},
		{Name: "Created", Value: func(i interface{}) string { return cmdutil.StringField(i, "createdAt") }},
		{Name: "Expire At", Value: func(i interface{}) string {
			s := cmdutil.StringField(i, "expireAt")
			if s == "" {
				return "-"
			}

			return s
		}},
	}
}
