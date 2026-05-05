package plugin

import (
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

var pluginTypes = map[string]string{
	"idp":          "identities",
	"factor":       "factors",
	"certificate":  "certificates",
	"policy":       "policies",
	"resource":     "resources",
	"reporter":     "reporters",
	"botdetection": "bot-detections",
}

func NewPluginCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Explore and create resources from plugin schemas",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newSchemaCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	return cmd
}
