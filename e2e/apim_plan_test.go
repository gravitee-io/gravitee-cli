//go:build e2e

package e2e

import (
	"testing"
)

// TestAPIMPlanLifecycle covers plan CRUD + lifecycle transitions (publish,
// deprecate, close).
//
// Creates a dedicated API for the test; teardown removes it via --close-plans.
func TestAPIMPlanLifecycle(t *testing.T) {
	apiFixture := writeFixture(t, "api.json")
	planFixture := writeFixture(t, "plan.json")
	planUpdatedFixture := writeFixture(t, "plan-updated.json")

	apiOut := runCLIExpectSuccess(t, "apim", "api", "create", "-f", apiFixture, "-o", "json")

	apiID := extractID(t, apiOut)
	if apiID == "" {
		t.Fatalf("api create returned no id: %s", apiOut)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "api", "stop", apiID)
		_, _ = runCLI("apim", "api", "delete", apiID, "--close-plans")
	})

	t.Run("create update delete in STAGING", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "apim", "plan", "create",
			"--api", apiID,
			"-f", planFixture,
			"-o", "json")

		planID := extractID(t, out)
		if planID == "" {
			t.Fatalf("plan create returned no id: %s", out)
		}

		runCLIExpectSuccess(t, "apim", "plan", "get", planID, "--api", apiID, "-o", "json")
		runCLIExpectSuccess(t, "apim", "plan", "list", "--api", apiID, "-o", "json")
		runCLIExpectSuccess(t, "apim", "plan", "update", planID, "--api", apiID, "-f", planUpdatedFixture)
		runCLIExpectSuccess(t, "apim", "plan", "delete", planID, "--api", apiID)
	})

	t.Run("publish deprecate close", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "apim", "plan", "create",
			"--api", apiID,
			"-f", planFixture,
			"-o", "json")

		planID := extractID(t, out)
		if planID == "" {
			t.Fatalf("plan create returned no id: %s", out)
		}

		runCLIExpectSuccess(t, "apim", "plan", "publish", planID, "--api", apiID)
		runCLIExpectSuccess(t, "apim", "plan", "list", "--api", apiID, "--status", "PUBLISHED", "-o", "json")
		runCLIExpectSuccess(t, "apim", "plan", "deprecate", planID, "--api", apiID)
		runCLIExpectSuccess(t, "apim", "plan", "close", planID, "--api", apiID)
	})
}
