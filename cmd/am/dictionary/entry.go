package dictionary

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
)

func newEntryCmd(f *factory.Factory, domainID *string) *cobra.Command {
	var dictID string

	cmd := &cobra.Command{
		Use:   "entry",
		Short: "Manage dictionary entries",
	}

	cmd.PersistentFlags().StringVar(&dictID, "dict-id", "", "Dictionary ID (required)")
	_ = cmd.MarkPersistentFlagRequired("dict-id")

	cmd.AddCommand(newEntryListCmd(f, domainID, &dictID))
	cmd.AddCommand(newEntryUpdateCmd(f, domainID, &dictID))

	return cmd
}

func newEntryListCmd(f *factory.Factory, domainID, dictID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List dictionary entries",
		Example: `  gio am dictionary entry list --domain my-domain --dict-id my-dict`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListDictionaryEntries(*domainID, *dictID)
			if err != nil {
				return err
			}

			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}

			return p.PrintDetail(data)
		},
	}
}

func newEntryUpdateCmd(f *factory.Factory, domainID, dictID *string) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:     "update --file <entries.json>",
		Short:   "Update dictionary entries from a JSON file",
		Example: `  gio am dictionary entry update --domain my-domain --dict-id my-dict --file entries.json`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			body, err := cmdutil.ReadJSONFile(file)
			if err != nil {
				return err
			}

			data, err := f.AM().UpdateDictionaryEntries(*domainID, *dictID, body)
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

			p.PrintMessage("Dictionary entries updated.")

			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON entries file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
