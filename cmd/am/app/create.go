package app

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

type createOptions struct {
	factory      *factory.Factory
	domainID     *string
	name         string
	appType      string
	description  string
	redirectURIs string
}

func newCreateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	opts := &createOptions{factory: f, domainID: domainID}

	cmd := &cobra.Command{
		Use:   "create --name <name> --type <type>",
		Short: "Create an application",
		Example: `  gio am app create --domain my-domain --name "My App" --type web
  gio am app create --domain my-domain --name "My App" --type browser --redirect-uris "http://localhost:4200/callback"`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.ValidateEnum(opts.appType, "type", []string{"web", "native", "browser", "service", "resource_server"}); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Application name (required)")
	cmd.Flags().StringVar(&opts.appType, "type", "", "Application type: web, native, browser, service, resource_server (required)")
	cmd.Flags().StringVar(&opts.description, "description", "", "Application description")
	cmd.Flags().StringVar(&opts.redirectURIs, "redirect-uris", "", "Comma-separated redirect URIs")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("type")

	return cmd
}

func (o *createOptions) run() error {
	f := o.factory

	body := map[string]any{
		"name": o.name,
		"type": o.appType,
	}

	if o.description != "" {
		body["description"] = o.description
	}

	if o.redirectURIs != "" {
		uris := strings.Split(o.redirectURIs, ",")
		for i := range uris {
			uris[i] = strings.TrimSpace(uris[i])
		}

		body["redirectUris"] = uris
	}

	raw, _ := json.Marshal(body)

	data, err := f.AM().CreateApplication(*o.domainID, json.RawMessage(raw))
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

	return printAppDetail(p, data)
}
