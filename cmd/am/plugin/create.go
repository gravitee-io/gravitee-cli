package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
	var name string
	var configFile string

	cmd := &cobra.Command{
		Use:   "create <type> <pluginId>",
		Short: "Create a resource instance from a plugin (use --config-file for non-interactive)",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			ptype, pluginID := args[0], args[1]
			apiPath, ok := pluginTypes[ptype]
			if !ok {
				return fmt.Errorf("unknown plugin type %q. Available: %s", ptype, sortedKeys(pluginTypes))
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if configFile == "" {
				return fmt.Errorf("--config-file is required (interactive mode not supported in gio CLI)")
			}
			raw, err := cmdutil.ReadJSONFile(configFile)
			if err != nil {
				return err
			}
			body := map[string]interface{}{
				"name":          name,
				"type":          pluginID,
				"configuration": string(raw),
			}
			path := cmdutil.AMDomainPath(f, apiPath)
			data, err := f.Client.Post(path, body)
			if err != nil {
				return err
			}
			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}
			return p.PrintDetail(json.RawMessage(data))
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Resource name (required)")
	cmd.Flags().StringVarP(&configFile, "config-file", "f", "", "JSON config file (required)")
	return cmd
}
