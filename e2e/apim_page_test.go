//go:build e2e

package e2e

import (
	"testing"
)

// TestAPIMPageLifecycle covers documentation page CRUD + publish/unpublish.
func TestAPIMPageLifecycle(t *testing.T) {
	apiFixture := writeFixture(t, "api.json")
	pageFixture := writeFixture(t, "page.json")
	pageUpdatedFixture := writeFixture(t, "page-updated.json")

	apiOut := runCLIExpectSuccess(t, "apim", "api", "create", "-f", apiFixture, "-o", "json")

	apiID := extractID(t, apiOut)
	if apiID == "" {
		t.Fatalf("api create returned no id: %s", apiOut)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "api", "stop", apiID)
		_, _ = runCLI("apim", "api", "delete", apiID, "--close-plans")
	})

	pageOut := runCLIExpectSuccess(t, "apim", "page", "create",
		"--api", apiID,
		"-f", pageFixture,
		"-o", "json")

	pageID := extractID(t, pageOut)
	if pageID == "" {
		t.Fatalf("page create returned no id: %s", pageOut)
	}

	runCLIExpectSuccess(t, "apim", "page", "get", pageID, "--api", apiID, "-o", "json")
	runCLIExpectSuccess(t, "apim", "page", "list", "--api", apiID, "-o", "json")
	runCLIExpectSuccess(t, "apim", "page", "update", pageID, "--api", apiID, "-f", pageUpdatedFixture)
	runCLIExpectSuccess(t, "apim", "page", "publish", pageID, "--api", apiID)
	runCLIExpectSuccess(t, "apim", "page", "unpublish", pageID, "--api", apiID)
	runCLIExpectSuccess(t, "apim", "page", "delete", pageID, "--api", apiID)
}
