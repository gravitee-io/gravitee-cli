package plan

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newDeprecateCmd(f *factory.Factory) *cobra.Command {
	var apiID string

	cmd := &cobra.Command{
		Use:     "deprecate <planId> --api <apiId>",
		Short:   "Deprecate a published plan",
		Example: `  gio plan deprecate aaaa1111-2222-3333-4444-555566667777 --api 8a7b3c4d-1234-5678-abcd-ef0123456789`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "plan deprecate"); err != nil {
				return err
			}

			return runDeprecate(f, apiID, args[0])
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	_ = cmd.MarkFlagRequired("api")

	return cmd
}

func runDeprecate(f *factory.Factory, apiID, planID string) error {
	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/plans/%s/_deprecate", apiID, planID))

	data, err := f.Client.Post(path, nil)
	if err != nil {
		return fmt.Errorf("plan deprecate failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printPlanDetail(p, data)
}
