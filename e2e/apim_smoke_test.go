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

// TestAPIMSmoke validates the APIM e2e harness end-to-end: the docker-compose
// is up, fetchAPIMToken() obtained a valid PAT, GCTL_APIM_* env vars are set,
// and the CLI can hit the management API and parse a JSON response.
//
// If this test fails, no other apim_*_test.go can pass - fix this first.
func TestAPIMSmoke(t *testing.T) {
	out := runCLIExpectSuccess(t, "apim", "env", "list", "-o", "json")

	var envs []map[string]any
	if err := json.Unmarshal([]byte(out), &envs); err != nil {
		t.Fatalf("expected JSON array from 'apim env list', got: %s\nparse error: %v", out, err)
	}

	if len(envs) == 0 {
		t.Fatal("expected at least one environment, got empty list")
	}
}
