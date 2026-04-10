//go:build e2e

package e2e

import (
	"encoding/json"
	"testing"
)

// TestAPIMPluginList exercises plugin listing for each supported type.
func TestAPIMPluginList(t *testing.T) {
	for _, pluginType := range []string{"endpoints", "entrypoints", "policies"} {
		t.Run(pluginType, func(t *testing.T) {
			out := runCLIExpectSuccess(t, "apim", "plugin", "list", "--type", pluginType, "-o", "json")

			var plugins []map[string]any
			if err := json.Unmarshal([]byte(out), &plugins); err != nil {
				t.Fatalf("expected JSON array from 'plugin list --type %s', got: %s\nparse error: %v", pluginType, out, err)
			}
		})
	}
}
