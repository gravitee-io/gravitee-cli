package user

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/am"
	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newUserAuditCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var userID string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Manage user audits",
	}

	cmd.PersistentFlags().StringVar(&userID, "user-id", "", "User ID (required)")
	_ = cmd.MarkPersistentFlagRequired("user-id")

	cmd.AddCommand(newUserAuditListCmd(f, domainID, &userID))
	cmd.AddCommand(newUserAuditGetCmd(f, domainID, &userID))

	return cmd
}

type userAuditListOptions struct {
	factory   *factory.Factory
	domainID  *string
	userID    *string
	auditType string
	status    string
	from      string
	to        string
	page      int
	perPage   int
}

func newUserAuditListCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	opts := &userAuditListOptions{factory: f, domainID: domainID, userID: userID}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List user audits",
		Example: `  gio am user audit list --domain my-domain --user-id user-1
  gio am user audit list --domain my-domain --user-id user-1 --type LOGIN --status SUCCESS
  gio am user audit list --domain my-domain --user-id user-1 --from 1609459200000 --to 1612137600000`,
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

	cmd.Flags().StringVar(&opts.auditType, "type", "", "Audit event type")
	cmd.Flags().StringVar(&opts.status, "status", "", "Audit event status")
	cmd.Flags().StringVar(&opts.from, "from", "", "From timestamp (epoch ms)")
	cmd.Flags().StringVar(&opts.to, "to", "", "To timestamp (epoch ms)")
	cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")

	return cmd
}

func (o *userAuditListOptions) run() error {
	f := o.factory

	params := am.ListUserAuditsParams{
		Type:    o.auditType,
		Status:  o.status,
		From:    o.from,
		To:      o.to,
		Page:    o.page - 1, // Convert 1-based CLI page to 0-based API page
		PerPage: o.perPage,
	}

	resp, err := f.AM().ListUserAudits(*o.domainID, *o.userID, params)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	if f.OutputFormat != printer.FormatTable {
		raw, _ := json.Marshal(resp)
		return p.PrintDetail(json.RawMessage(raw))
	}

	if err := p.PrintList(resp.Data, userAuditColumns()); err != nil {
		return err
	}

	if resp.TotalCount > len(resp.Data) {
		p.PrintMessage("Showing %d of %d.", len(resp.Data), resp.TotalCount)
	} else if resp.TotalCount > 0 {
		p.PrintMessage("Showing %d results.", len(resp.Data))
	}

	return nil
}

func newUserAuditGetCmd(f *factory.Factory, domainID, userID *string) *cobra.Command {
	return &cobra.Command{
		Use:   "get <auditID>",
		Short: "Get user audit details",
		Example: `  gio am user audit get audit-1 --domain my-domain --user-id user-1
  gio am user audit get audit-1 --domain my-domain --user-id user-1 -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetUserAudit(*domainID, *userID, args[0])
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			return p.PrintDetail(data)
		},
	}
}

func userAuditColumns() []printer.Column {
	return []printer.Column{
		{Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
		{Name: "Type", Value: func(i any) string { return cmdutil.StringField(i, "type") }},
		{Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "status") }},
	}
}
