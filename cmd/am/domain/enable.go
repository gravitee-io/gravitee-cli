package domain

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newEnableCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "enable <domainID>",
		Short:   "Enable a security domain",
		Example: `  gio am domain enable my-domain-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runSetEnabled(f, args[0], true)
		},
	}
}

func newDisableCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "disable <domainID>",
		Short:   "Disable a security domain",
		Example: `  gio am domain disable my-domain-id`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireContext(f); err != nil {
				return err
			}

			return runSetEnabled(f, args[0], false)
		},
	}
}

func runSetEnabled(f *factory.Factory, domainID string, enabled bool) error {
	body := map[string]any{"enabled": enabled}
	raw, _ := json.Marshal(body)

	if _, err := f.AM().PatchDomain(domainID, json.RawMessage(raw)); err != nil {
		return err
	}

	p, err := cmdutil.NewPrinter(f)
	if err != nil {
		return err
	}

	action := "enabled"
	if !enabled {
		action = "disabled"
	}

	p.PrintMessage("Domain '%s' %s.", domainID, action)

	return nil
}
