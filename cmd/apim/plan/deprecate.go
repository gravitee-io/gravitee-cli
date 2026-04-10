package plan

import (
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
		Example: `  gio apim plan deprecate aaaa1111-2222-3333-4444-555566667777 --api /my/api`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDeprecate(f, apiID, args[0])
		},
	}

	cmdutil.AddAPIFlag(cmd, &apiID)

	return cmd
}

func runDeprecate(f *factory.Factory, apiID, planID string) error {
	data, err := f.APIM().DeprecatePlan(apiID, planID)
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

	return printPlanDetail(p, data)
}
