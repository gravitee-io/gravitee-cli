package subscription

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newCloseCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "close <subId> --api <apiId>",
		Short:   "Close a subscription permanently",
		Example: `  gio subscription close 34f8c9e7 --api 8a7b3c4d`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "subscription close"); err != nil {
				return err
			}

			return runClose(f, apiID, args[0])
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")

	return cmd
}

func runClose(f *factory.Factory, apiID, subID string) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/subscriptions/%s/_close", apiID, subID))

	data, err := f.Client.Post(path, nil)
	if err != nil {
		return fmt.Errorf("subscription close failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printSubDetail(p, data)
}
