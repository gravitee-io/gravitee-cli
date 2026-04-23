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

// TestAPIMEnvironmentGet retrieves the default environment and verifies the server returns a JSON object.
func TestAPIMEnvironmentGet(t *testing.T) {
	out := runCLIExpectSuccess(t, "apim", "env", "get", "DEFAULT", "-o", "json")

	var env map[string]any
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("expected JSON object from 'apim env get', got: %s\nparse error: %v", out, err)
	}

	if env["id"] == nil && env["hrids"] == nil && env["name"] == nil {
		t.Errorf("expected at least one of id/hrids/name in environment response, got: %v", env)
	}
}
