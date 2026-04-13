package role

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type createOptions struct {
	factory     *factory.Factory
	domainID    *string
	name        string
	description string
}

func newCreateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &createOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "create --name <name>",
		Short: "Create a role",
		Example: `  gio am role create --domain my-domain --name "Admin"
  gio am role create --domain my-domain --name "Admin" --description "Administrator role"`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Role name (required)")
	cmd.Flags().StringVar(&opts.description, "description", "", "Role description")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func (o *createOptions) run() error {
	f := o.factory

	body := map[string]any{"name": o.name}
	if o.description != "" {
		body["description"] = o.description
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().CreateRole(*o.domainID, json.RawMessage(raw))
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

	return printRoleDetail(p, data)
}
