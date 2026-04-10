//go:build e2e

package e2e

import (
	"encoding/json"
	"testing"
)

// TestAPIMSmoke validates the APIM e2e harness end-to-end: the docker-compose
// is up, fetchAPIMToken() obtained a valid PAT, GIO_APIM_* env vars are set,
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
