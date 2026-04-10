//go:build e2e

package e2e

import (
	"encoding/json"
	"testing"
)

// TestAPIMAPIKeyLifecycle covers API key operations on a subscription:
// list → renew → list → revoke → reactivate.
func TestAPIMAPIKeyLifecycle(t *testing.T) {
	apiID, planID := setupPublishedAPI(t, "plan.json")
	appID := setupApplication(t)

	subOut := runCLIExpectSuccess(t, "apim", "sub", "create",
		"--api", apiID,
		"--plan", planID,
		"--app", appID,
		"-o", "json")

	subID := extractID(t, subOut)
	if subID == "" {
		t.Fatalf("sub create returned no id: %s", subOut)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "sub", "close", subID, "--api", apiID)
	})

	runCLIExpectSuccess(t, "apim", "api-key", "list", "--api", apiID, "--subscription", subID, "-o", "json")
	runCLIExpectSuccess(t, "apim", "api-key", "renew", "--api", apiID, "--subscription", subID)

	// After renew, pick the first key from the list and exercise revoke/reactivate on it.
	listOut := runCLIExpectSuccess(t, "apim", "api-key", "list",
		"--api", apiID, "--subscription", subID, "-o", "json")

	var listResp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(listOut), &listResp); err != nil {
		t.Fatalf("failed to parse api-key list: %v", err)
	}

	if len(listResp.Data) == 0 {
		t.Fatal("expected at least one api-key after renew")
	}

	keyID := listResp.Data[0].ID

	runCLIExpectSuccess(t, "apim", "api-key", "revoke", keyID, "--api", apiID, "--subscription", subID)
	runCLIExpectSuccess(t, "apim", "api-key", "reactivate", keyID, "--api", apiID, "--subscription", subID)
}
