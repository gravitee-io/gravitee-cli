// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
