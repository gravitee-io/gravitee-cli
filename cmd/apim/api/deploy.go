package api

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeployCmd(f *factory.Factory) *cobra.Command {
	var label string

	cmd := &cobra.Command{
		Use:   "deploy <apiId>",
		Short: "Deploy an API",
		Example: `  gio apim api deploy 8a7b3c4d-1234-5678-abcd-ef0123456789
  gio apim api deploy 8a7b3c4d-1234-5678-abcd-ef0123456789 --label "v2.1.0 hotfix"`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runDeploy(f, args[0], label)
		},
	}

	cmd.Flags().StringVar(&label, "label", "", "Deployment label (32 characters max)")

	return cmd
}

func runDeploy(f *factory.Factory, apiID, label string) error {
	if len(label) > 32 {
		return fmt.Errorf("deployment label exceeds 32 characters")
	}

	if err := f.APIM().DeployAPI(apiID, label); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}
	p.PrintMessage("API '%s' deployment requested.", apiID)

	return nil
}
