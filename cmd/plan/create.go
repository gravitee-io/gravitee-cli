package plan

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
	var (
		apiID string
		file  string
	)

	cmd := &cobra.Command{
		Use:     "create --api <apiId> -f <file>",
		Short:   "Create a plan from a JSON file",
		Example: `  gio plan create --api 8a7b3c4d-1234-5678-abcd-ef0123456789 -f plan.json`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			if err := cmdutil.CheckReadOnly(f, "plan create"); err != nil {
				return err
			}

			return runCreate(f, apiID, file)
		},
	}

	cmd.Flags().StringVar(&apiID, "api", "", "API ID (required)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("api")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runCreate(f *factory.Factory, apiID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	path := cmdutil.V2EnvPath(f, fmt.Sprintf("apis/%s/plans", apiID))

	data, err := f.Client.Post(path, body)
	if err != nil {
		return fmt.Errorf("plan creation failed: %w", err)
	}

	p := cmdutil.NewPrinter(f)

	if f.OutputFormat != printer.FormatTable {
		return p.PrintDetail(json.RawMessage(data))
	}

	return printPlanDetail(p, data)
}
