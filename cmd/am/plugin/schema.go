package plugin

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newSchemaCmd(f *factory.Factory) *cobra.Command {
	var raw bool
	cmd := &cobra.Command{
		Use:   "schema <type> <pluginId>",
		Short: "Show configuration schema for a plugin",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			ptype, pluginID := args[0], args[1]
			apiType, ok := pluginTypes[ptype]
			if !ok {
				return fmt.Errorf("unknown plugin type %q. Available: %s", ptype, sortedKeys(pluginTypes))
			}
			path := fmt.Sprintf("/management/platform/plugins/%s/%s/schema", apiType, pluginID)
			data, err := f.Client.Get(path)
			if err != nil {
				return err
			}
			if raw {
				fmt.Fprintln(f.IOStreams.Out, string(data))
				return nil
			}
			var schema map[string]interface{}
			if err := json.Unmarshal(data, &schema); err != nil {
				fmt.Fprintln(f.IOStreams.Out, string(data))
				return nil
			}
			printSchema(f, schema, pluginID)
			return nil
		},
	}
	cmd.Flags().BoolVar(&raw, "raw", false, "Show raw JSON schema")
	return cmd
}

func printSchema(f *factory.Factory, schema map[string]interface{}, pluginID string) {
	fmt.Fprintf(f.IOStreams.Out, "Schema: %s\n\n", pluginID)
	props, _ := schema["properties"].(map[string]interface{})

	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := props[key]
		prop, _ := val.(map[string]interface{})
		title, _ := prop["title"].(string)
		propType, _ := prop["type"].(string)
		desc, _ := prop["description"].(string)
		if title == "" {
			title = key
		}
		line := fmt.Sprintf("  %-30s %-10s", key, propType)
		if desc != "" {
			line += "  " + desc
		} else if title != key {
			line += "  " + title
		}
		fmt.Fprintln(f.IOStreams.Out, line)
	}

	for _, key := range keys {
		val := props[key]
		prop, _ := val.(map[string]interface{})
		if enums, ok := prop["enum"].([]interface{}); ok {
			strs := make([]string, 0, len(enums))
			for _, e := range enums {
				strs = append(strs, fmt.Sprintf("%v", e))
			}
			fmt.Fprintf(f.IOStreams.Out, "  %s: {%s}\n", key, strings.Join(strs, ", "))
		}
	}
}
