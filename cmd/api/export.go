package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newExportCmd(f *factory.Factory) *cobra.Command {
	var exclude []string

	cmd := &cobra.Command{
		Use:   "export <apiId>",
		Short: "Export an API definition",
		Example: `  gio api export 8a7b3c4d-1234-5678-abcd-ef0123456789
  gio api export 8a7b3c4d-1234-5678-abcd-ef0123456789 --exclude members --exclude pages`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runExport(f, args[0], exclude)
		},
	}

	cmd.Flags().StringArrayVar(&exclude, "exclude", nil,
		"Exclude data from export: groups, members, metadata, pages, plans")

	return cmd
}

func runExport(f *factory.Factory, apiID string, exclude []string) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/_export/definition", apiID))

	if len(exclude) > 0 {
		escaped := make([]string, len(exclude))
		for i, e := range exclude {
			escaped[i] = url.QueryEscape(e)
		}

		path += "?excludeAdditionalData=" + strings.Join(escaped, ",")
	}

	data, err := f.Client.Get(path)
	if err != nil {
		return err
	}

	p := cmdutil.NewPrinter(f)

	return p.PrintDetail(json.RawMessage(data))
}
