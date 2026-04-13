package group

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
		Short: "Create a group",
		Example: `  gio am group create --domain my-domain --name "Admins"
  gio am group create --domain my-domain --name "Admins" --description "Administrator group"`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Group name (required)")
	cmd.Flags().StringVar(&opts.description, "description", "", "Group description")
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

	data, err := f.AM().CreateGroup(*o.domainID, json.RawMessage(raw))
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

	return printGroupDetail(p, data)
}
