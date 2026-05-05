package plugin

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
	"github.com/spf13/cobra"
)

func newListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list <type>",
		Short: "List available plugins of a given type",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			ptype := args[0]
			apiType, ok := pluginTypes[ptype]
			if !ok {
				return fmt.Errorf("unknown plugin type %q. Available: %s", ptype, sortedKeys(pluginTypes))
			}
			data, err := f.Client.Get("/management/platform/plugins/" + apiType)
			if err != nil {
				return err
			}
			p, err := cmdutil.NewPrinter(f)
			if err != nil {
				return err
			}
			var items []json.RawMessage
			if err := json.Unmarshal(data, &items); err != nil {
				return p.PrintDetail(json.RawMessage(data))
			}
			return p.PrintList(items, []printer.Column{
				{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
				{Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
				{Name: "Version", Value: func(i interface{}) string { return cmdutil.StringField(i, "version") }},
			})
		},
	}
}

func sortedKeys(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
