//go:build e2e

package e2e

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestAPIMApplicationLifecycle covers application CRUD and quiet-mode output.
func TestAPIMApplicationLifecycle(t *testing.T) {
	appFixture := writeFixture(t, "app.json")
	appUpdatedFixture := writeFixture(t, "app-updated.json")

	out := runCLIExpectSuccess(t, "apim", "app", "create", "-f", appFixture, "-o", "json")

	appID := extractID(t, out)
	if appID == "" {
		t.Fatalf("app create returned no id: %s", out)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "app", "delete", appID)
	})

	runCLIExpectSuccess(t, "apim", "app", "get", appID, "-o", "json")
	runCLIExpectSuccess(t, "apim", "app", "list", "-o", "json")
	runCLIExpectSuccess(t, "apim", "app", "update", appID, "-f", appUpdatedFixture)

	got := runCLIExpectSuccess(t, "apim", "app", "get", appID, "-o", "json")

	var app map[string]any
	if err := json.Unmarshal([]byte(got), &app); err != nil {
		t.Fatalf("failed to parse updated app: %v", err)
	}

	desc, _ := app["description"].(string)
	if !strings.Contains(desc, "Updated description") {
		t.Errorf("expected updated description, got: %q", desc)
	}

	idOnly := runCLIExpectSuccess(t, "apim", "app", "get", appID, "-o", "id")
	if got := strings.TrimSpace(idOnly); got != appID {
		t.Errorf("expected id output %q, got %q", appID, got)
	}
}
