package page

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newPublishCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "publish <pageId> --api <apiId>",
		Short:   "Publish a page",
		Example: `  gio page publish dddd1111-2222-3333-4444-555566667777 --api 8a7b3c4d-1234-5678-abcd-ef0123456789`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "page publish"); err != nil {
				return err
			}

			return runPublish(f, apiID, args[0])
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")

	return cmd
}

func runPublish(f *factory.Factory, apiID, pageID string) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/pages/%s/_publish", apiID, pageID))

	data, err := f.Client.Post(path, nil)
	if err != nil {
		return fmt.Errorf("page publish failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printPageDetail(p, data)
}
