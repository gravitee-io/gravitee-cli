package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newListCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all contexts",
		Example: `  gio context list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runList(f)
		},
	}

	cmdutil.AddOutputFlags(cmd, f)

	return cmd
}

func runList(f *factory.Factory) error {
	if err := cmdutil.SetupConfig(f); err != nil {
		return err
	}

	cfg := f.Config
	names := cfg.ContextNames()

	if len(names) == 0 {
		fmt.Fprintln(f.IOStreams.Out, "No contexts configured.")

		return nil
	}

	// Build a list of maps for printer consumption.
	items := make([]map[string]any, 0, len(names))
	for _, name := range names {
		ctx := cfg.Contexts[name]

		current := ""
		if name == cfg.Current {
			current = "*"
		}

		hasAPIM := "no"
		if ctx.APIM != nil {
			hasAPIM = "yes"
		}

		hasAM := "no"
		if ctx.AM != nil {
			hasAM = "yes"
		}

		readOnly := "no"
		if ctx.ReadOnly {
			readOnly = "yes"
		}

		items = append(items, map[string]any{
			"current":  current,
			"name":     name,
			"org":      ctx.Org,
			"env":      ctx.Env,
			"apim":     hasAPIM,
			"am":       hasAM,
			"readOnly": readOnly,
		})
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintList(items, contextColumns())
}

func contextColumns() []printer.Column {
	return []printer.Column{
		{Name: "Current", Value: func(i any) string { return cmdutil.StringField(i, "current") }, Width: 1},
		{Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
		{Name: "Org", Value: func(i any) string { return cmdutil.StringField(i, "org") }},
		{Name: "Env", Value: func(i any) string { return cmdutil.StringField(i, "env") }},
		{Name: "APIM", Value: func(i any) string { return cmdutil.StringField(i, "apim") }},
		{Name: "AM", Value: func(i any) string { return cmdutil.StringField(i, "am") }},
		{Name: "Read-Only", Value: func(i any) string { return cmdutil.StringField(i, "readOnly") }},
	}
}
