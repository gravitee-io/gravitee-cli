package user

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type createOptions struct {
	factory         *factory.Factory
	domainID        *string
	username        string
	email           string
	password        string
	firstName       string
	lastName        string
	preRegistration bool
}

func newCreateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &createOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "create --username <username> --email <email>",
		Short: "Create a user",
		Example: `  gio am user create --domain my-domain --username john --email john@example.com
  gio am user create --domain my-domain --username john --email john@example.com --password secret --firstName John --lastName Doe`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.username, "username", "", "Username (required)")
	cmd.Flags().StringVar(&opts.email, "email", "", "Email address (required)")
	cmd.Flags().StringVar(&opts.password, "password", "", "Password")
	cmd.Flags().StringVar(&opts.firstName, "firstName", "", "First name")
	cmd.Flags().StringVar(&opts.lastName, "lastName", "", "Last name")
	cmd.Flags().BoolVar(&opts.preRegistration, "preRegistration", false, "Pre-registration flag")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("email")

	return cmd
}

func (o *createOptions) run() error {
	f := o.factory

	body := map[string]any{
		"username": o.username,
		"email":    o.email,
	}

	if o.password != "" {
		body["password"] = o.password
	}

	if o.firstName != "" {
		body["firstName"] = o.firstName
	}

	if o.lastName != "" {
		body["lastName"] = o.lastName
	}

	if o.preRegistration {
		body["preRegistration"] = true
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().CreateUser(*o.domainID, json.RawMessage(raw))
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

	return printUserDetail(p, data)
}
