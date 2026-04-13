package idp

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newUpdateCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "update <idpID> --file <config.json>",
		Short: "Update an identity provider from a JSON file",
		Example: `  gio am idp update my-idp-id --domain my-domain --file idp.json
  gio am idp update my-idp-id --domain my-domain -f idp.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runUpdate(f, *domainID, args[0], file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func runUpdate(f *factory.Factory, domainID, idpID, file string) error {
	body, err := cmdutil.ReadJSONFile(file)
	if err != nil {
		return err
	}

	data, err := f.AM().UpdateIdentityProvider(domainID, idpID, body)
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

	return printIDPDetail(p, data)
}
