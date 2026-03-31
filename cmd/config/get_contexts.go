package config

import (
	"sort"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	iconfig "github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newGetContextsCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get-contexts",
		Short: "List all configured contexts",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runGetContexts(f)
		},
	}
}

type contextRow struct {
	Current  string `json:"current"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Env      string `json:"env"`
	ReadOnly string `json:"readOnly"`
}

func runGetContexts(f *factory.Factory) error {
	cfg := f.Config
	p := printer.New(f.OutputFormat, f.IOStreams.Out, f.Quiet)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(cfg)
	}

	rows := buildContextRows(cfg)
	columns := contextColumns()

	return p.PrintList(rows, columns)
}

func buildContextRows(cfg *iconfig.Config) []contextRow {
	names := make([]string, 0, len(cfg.Contexts))
	for name := range cfg.Contexts {
		names = append(names, name)
	}

	sort.Strings(names)

	rows := make([]contextRow, 0, len(names))
	for _, name := range names {
		ctx := cfg.Contexts[name]

		marker := " "
		if name == cfg.CurrentContext {
			marker = "*"
		}

		env := ctx.Env
		if env == "" {
			env = iconfig.DefaultEnv
		}

		ro := "no"
		if ctx.ReadOnly {
			ro = "yes"
		}

		rows = append(rows, contextRow{
			Current:  marker,
			Name:     name,
			URL:      ctx.URL,
			Env:      env,
			ReadOnly: ro,
		})
	}

	return rows
}

func contextColumns() []printer.Column {
	return []printer.Column{
		{Name: " ", Width: 1, Value: func(item interface{}) string { return cmdutil.StringField(item, "current") }},
		{Name: "Name", Value: func(item interface{}) string { return cmdutil.StringField(item, "name") }},
		{Name: "URL", Value: func(item interface{}) string { return cmdutil.StringField(item, "url") }},
		{Name: "Env", Value: func(item interface{}) string { return cmdutil.StringField(item, "env") }},
		{Name: "Read-Only", Value: func(item interface{}) string { return cmdutil.StringField(item, "readOnly") }},
	}
}
