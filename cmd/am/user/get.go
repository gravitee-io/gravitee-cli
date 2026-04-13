package user

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type getOptions struct {
	factory  *factory.Factory
	domainID *string
	userID   string
}

func newGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &getOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "get <userID>",
		Short: "Get user details",
		Example: `  gio am user get user-id --domain my-domain
  gio am user get user-id --domain my-domain -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.userID = args[0]

			return opts.run()
		},
	}

	return cmd
}

func (o *getOptions) run() error {
	f := o.factory

	data, err := f.AM().GetUser(*o.domainID, o.userID)
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

func printUserDetail(p *printer.Printer, data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, field := range []struct{ label, key string }{
		{"Username", "username"},
		{"ID", "id"},
		{"Email", "email"},
		{"First Name", "firstName"},
		{"Last Name", "lastName"},
		{"Enabled", "enabled"},
	} {
		if v, ok := m[field.key]; ok && v != nil {
			p.PrintMessage("%-16s%v", field.label+":", v)
		}
	}

	return nil
}
