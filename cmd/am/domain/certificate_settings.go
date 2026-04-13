package domain

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newUpdateCertSettingsCmd(f *factory.Factory) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:     "update-cert-settings <domainID> --file <settings.json>",
		Short:   "Update domain certificate settings",
		Example: `  gio am domain update-cert-settings my-domain-id --file settings.json`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdateCertSettings(f, args[0], file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runUpdateCertSettings(f *factory.Factory, domainID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateDomainCertificateSettings(domainID, body)
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

	p.PrintMessage("Certificate settings updated for domain '%s'.", domainID)

	return nil
}
