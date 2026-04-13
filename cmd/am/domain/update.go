package domain

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type updateOptions struct {
	factory     *factory.Factory
	domainID    string
	name        string
	description string
}

func newUpdateCmd(f *factory.Factory) *cobra.Command {
	opts := &updateOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "update <domainID>",
		Short: "Update a security domain",
		Example: `  gio am domain update my-domain-id --name "New Name"
  gio am domain update my-domain-id --description "Updated description"`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			opts.domainID = args[0]

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Domain name")
	cmd.Flags().StringVar(&opts.description, "description", "", "Domain description")

	return cmd
}

func (o *updateOptions) run() error {
	f := o.factory

	body := map[string]any{}
	if o.name != "" {
		body["name"] = o.name
	}

	if o.description != "" {
		body["description"] = o.description
	}

	if len(body) == 0 {
		return fmt.Errorf("at least one flag (--name, --description) is required")
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().PatchDomain(o.domainID, json.RawMessage(raw))
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

	return printDomainDetail(p, data)
}
