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
	"strings"
	"testing"
)

// TestAPIMAPILifecycle covers the full API CRUD + lifecycle flow.
//
// Sub-tests run sequentially and share an API created in "create". The teardown
// (stop + delete --close-plans) runs via t.Cleanup so every exit path leaves
// APIM in a clean state.
func TestAPIMAPILifecycle(t *testing.T) {
	apiFixture := writeFixture(t, "api.json")
	apiUpdatedFixture := writeFixture(t, "api-updated.json")
	planFixture := writeFixture(t, "plan.json")

	out := runCLIExpectSuccess(t, "apim", "api", "create", "-f", apiFixture, "-o", "json")

	apiID := extractID(t, out)
	if apiID == "" {
		t.Fatalf("api create returned no id: %s", out)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "api", "stop", apiID)
		_, _ = runCLI("apim", "api", "delete", apiID, "--close-plans")
	})

	t.Run("get", func(t *testing.T) {
		runCLIExpectSuccess(t, "apim", "api", "get", apiID, "-o", "json")
	})

	t.Run("list", func(t *testing.T) {
		runCLIExpectSuccess(t, "apim", "api", "list", "-o", "json")
	})

	t.Run("update", func(t *testing.T) {
		runCLIExpectSuccess(t, "apim", "api", "update", apiID, "-f", apiUpdatedFixture)

		got := runCLIExpectSuccess(t, "apim", "api", "get", apiID, "-o", "json")

		var api map[string]any
		if err := json.Unmarshal([]byte(got), &api); err != nil {
			t.Fatalf("failed to parse updated api: %v", err)
		}

		desc, _ := api["description"].(string)
		if !strings.Contains(desc, "Updated description") {
			t.Errorf("expected updated description, got: %q", desc)
		}
	})

	t.Run("get -o id returns id", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "apim", "api", "get", apiID, "-o", "id")

		got := strings.TrimSpace(out)
		if got != apiID {
			t.Errorf("expected id output %q, got %q", apiID, got)
		}
	})

	t.Run("export", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "apim", "api", "export", apiID, "-o", "json")

		if !strings.HasPrefix(strings.TrimSpace(out), "{") {
			t.Errorf("expected JSON object from export, got: %s", out)
		}
	})

	t.Run("deploy and lifecycle", func(t *testing.T) {
		planOut := runCLIExpectSuccess(t, "apim", "plan", "create",
			"--api", apiID,
			"-f", planFixture,
			"-o", "json")

		planID := extractID(t, planOut)
		if planID == "" {
			t.Fatalf("plan create returned no id: %s", planOut)
		}

		// Publish must succeed for deploy to work. Tolerate any error here -
		// the lifecycle sub-asserts below will fail if deploy breaks.
		_, _ = runCLI("apim", "plan", "publish", planID, "--api", apiID)

		runCLIExpectSuccess(t, "apim", "api", "deploy", apiID)
		runCLIExpectSuccess(t, "apim", "api", "start", apiID)
		runCLIExpectSuccess(t, "apim", "api", "stop", apiID)
	})

	t.Run("observability wiring", func(t *testing.T) {
		// These commands must not error even if the underlying data is empty
		// (no traffic means logs/analytics return empty paginated responses).
		runCLIExpectSuccess(t, "apim", "health", "--api", apiID)
		runCLIExpectSuccess(t, "apim", "log", "list", "--api", apiID, "-o", "json")
		runCLIExpectSuccess(t, "apim", "analytics",
			"--api", apiID,
			"--type", "COUNT",
			"--from", "1",
			"--to", "9999999999999",
			"-o", "json")
	})

	t.Run("rollback rejects unknown event", func(t *testing.T) {
		// Exercising a real rollback would require fetching a prior event id,
		// which is not exposed by the CLI. Sending a synthetic UUID proves the
		// command reaches the server and propagates its error untouched.
		bogusEvent := "00000000-0000-0000-0000-000000000000"

		errOut := runCLIExpectError(t, "apim", "api", "rollback", apiID, "--event-id", bogusEvent)

		if !strings.Contains(strings.ToLower(errOut), "not found") && !strings.Contains(strings.ToLower(errOut), "event") {
			t.Errorf("expected event-not-found error, got: %s", errOut)
		}
	})
}
