package page

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newUpdateCmd(f *factory.Factory) *cobra.Command {
	var (
		apiID string
		file  string
	)

	cmd := &cobra.Command{
		Use:     "update <pageId> --api <apiId> -f <file>",
		Short:   "Update a page from a JSON file",
		Example: `  gio page update aaaa1111-2222-3333-4444-555566667777 --api 8a7b3c4d-1234-5678-abcd-ef0123456789 -f page-updated.json`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "page update"); err != nil {
				return err
			}

			return runUpdate(f, apiID, args[0], file)
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("api")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runUpdate(f *factory.Factory, apiID, pageID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/pages/%s", apiID, pageID))

	data, err := f.Client.Put(path, body)
	if err != nil {
		return fmt.Errorf("page update failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printPageDetail(p, data)
}
