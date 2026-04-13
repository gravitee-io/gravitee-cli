package dataplane

import (
	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

// NewDataPlaneCmd creates the data-plane parent command with all subcommands.
func NewDataPlaneCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "data-plane",
		Aliases: []string{"dp", "data-planes"},
		Short:   "Manage data planes",
	}

	cmdutil.AddOutputFlags(cmd, f)

	cmd.AddCommand(newListCmd(f))

	return cmd
}

func newListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List data planes",
		Example: `  gio am data-plane list`,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			data, err := f.AM().ListDataPlanes()
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
