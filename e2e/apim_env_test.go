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
