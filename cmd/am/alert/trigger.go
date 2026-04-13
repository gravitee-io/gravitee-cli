package alert

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newTriggerCmd(f *factory.Factory, domainID *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trigger",
		Aliases: []string{"triggers"},
		Short:   "Manage alert triggers",
	}

	cmd.AddCommand(newTriggerGetCmd(f, domainID))
	cmd.AddCommand(newTriggerUpdateCmd(f, domainID))

	return cmd
}

func newTriggerGetCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Short:   "Get alert triggers",
		Example: `  gio am alert trigger get --domain my-domain`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runTriggerGet(f, *domainID)
		},
	}
}

func runTriggerGet(f *factory.Factory, domainID string) error {
	data, err := f.AM().GetAlertTriggers(domainID)
	if err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	return p.PrintDetail(data)
}

func newTriggerUpdateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update --file <triggers.json>",
		Short: "Update alert triggers from a JSON file",
		Example: `  gio am alert trigger update --domain my-domain --file triggers.json
  gio am alert trigger update --domain my-domain -f triggers.json`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runTriggerUpdate(f, *domainID, file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runTriggerUpdate(f *factory.Factory, domainID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateAlertTriggers(domainID, body)
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

	p.PrintMessage("Alert triggers updated successfully.")

	return nil
}
