package certificate

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newKeyCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "key <certID>",
		Short:   "Get the public key of a certificate",
		Example: `  gio am certificate key cert-123 --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetCertificateKey(*domainID, args[0])
			if err != nil {
				return err
			}

			// /key returns raw text (SSH public key), not JSON.
			fmt.Fprintln(f.IOStreams.Out, string(data))

			return nil
		},
	}
}

func newKeysCmd(f *factory.Factory, domainID *string) *cobra.Command {
	return &cobra.Command{
		Use:     "keys <certID>",
		Short:   "Get all keys of a certificate",
		Example: `  gio am certificate keys cert-123 --domain my-domain`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().GetCertificateKeys(*domainID, args[0])
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
