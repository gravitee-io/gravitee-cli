package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

var validTypes = []string{"endpoints", "entrypoints", "policies"}

// typeLabels maps the plural API path segment to the singular label for the TYPE column.
var typeLabels = map[string]string{
	"endpoints":   "endpoint",
	"entrypoints": "entrypoint",
	"policies":    "policy",
}

type listOptions struct {
	factory    *factory.Factory
	pluginType string
}

func newListCmd(f *factory.Factory) *cobra.Command {
	opts := &listOptions{factory: f}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List plugins",
		Example: `  gio plugin list
  gio plugin list --type policies
  gio plugin list --type endpoints -o json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := opts.validate(); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().StringVarP(&opts.pluginType, "type", "t", "", "Filter by plugin type: endpoints, entrypoints, policies")

	return cmd
}

func (o *listOptions) validate() error {
	if o.pluginType == "" {
		return nil
	}

	for _, t := range validTypes {
		if o.pluginType == t {
			return nil
		}
	}

	return fmt.Errorf("invalid value '%s' for flag --type\nHint: allowed values are endpoints, entrypoints, policies", o.pluginType)
}

func (o *listOptions) run() error {
	f := o.factory
	p := cmdutil.NewPrinter(f)

	if o.pluginType != "" {
		return o.fetchSingleType(f, p)
	}

	return o.fetchAllTypes(f, p)
}

func (o *listOptions) fetchSingleType(f *factory.Factory, p *printer.Printer) error {
	path := fmt.Sprintf("/management/v2/organizations/%s/plugins/%s", f.Resolved.Org, o.pluginType)

	data, err := f.Client.Get(path)
	if err != nil {
		return fmt.Errorf("plugin list failed: %w", err)
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return p.PrintList(json.RawMessage(data), pluginColumns())
}

func (o *listOptions) fetchAllTypes(f *factory.Factory, p *printer.Printer) error {
	var allItems []interface{}

	for _, t := range validTypes {
		path := fmt.Sprintf("/management/v2/organizations/%s/plugins/%s", f.Resolved.Org, t)

		data, err := f.Client.Get(path)
		if err != nil {
			return fmt.Errorf("plugin list failed: %w", err)
		}

		var items []interface{}
		if err := json.Unmarshal(data, &items); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		label := typeLabels[t]
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				m["type"] = label
			}

			allItems = append(allItems, item)
		}
	}

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(allItems)
	}

	return p.PrintList(allItems, pluginColumnsWithType())
}

func pluginColumns() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
		{Name: "Version", Value: func(i interface{}) string { return cmdutil.StringField(i, "version") }},
		{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
		{Name: "Description", Value: func(i interface{}) string { return cmdutil.StringField(i, "description") }},
	}
}

func pluginColumnsWithType() []printer.Column {
	return []printer.Column{
		{Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
		{Name: "Type", Value: func(i interface{}) string { return cmdutil.StringField(i, "type") }},
		{Name: "Version", Value: func(i interface{}) string { return cmdutil.StringField(i, "version") }},
		{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
		{Name: "Description", Value: func(i interface{}) string { return cmdutil.StringField(i, "description") }},
	}
}
