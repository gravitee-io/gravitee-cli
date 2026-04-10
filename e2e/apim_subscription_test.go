//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

// TestAPIMSubscriptionAuto covers a subscription lifecycle with an AUTO-validated plan:
// create → get → list → pause → resume → close.
func TestAPIMSubscriptionAuto(t *testing.T) {
	apiID, planID := setupPublishedAPI(t, "plan.json")
	appID := setupApplication(t)

	out := runCLIExpectSuccess(t, "apim", "sub", "create",
		"--api", apiID,
		"--plan", planID,
		"--app", appID,
		"-o", "json")

	subID := extractID(t, out)
	if subID == "" {
		t.Fatalf("sub create returned no id: %s", out)
	}

	runCLIExpectSuccess(t, "apim", "sub", "get", subID, "--api", apiID, "-o", "json")
	runCLIExpectSuccess(t, "apim", "sub", "list", "--api", apiID, "-o", "json")
	runCLIExpectSuccess(t, "apim", "sub", "pause", subID, "--api", apiID)
	runCLIExpectSuccess(t, "apim", "sub", "resume", subID, "--api", apiID)
	runCLIExpectSuccess(t, "apim", "sub", "close", subID, "--api", apiID)
}

// TestAPIMSubscriptionManual verifies the CLI serialization of accept/reject
// with a MANUAL-validation plan. APIM auto-accepts subscriptions created by the
// API owner (single-user setup), so accept/reject may fail at the server layer
// with a state-transition error - but they must NOT fail with a nil-body
// serialization bug from the client.
func TestAPIMSubscriptionManual(t *testing.T) {
	apiID, planID := setupPublishedAPI(t, "plan-apikey-manual.json")
	appID := setupApplication(t)

	out := runCLIExpectSuccess(t, "apim", "sub", "create",
		"--api", apiID,
		"--plan", planID,
		"--app", appID,
		"-o", "json")

	subID := extractID(t, out)
	if subID == "" {
		t.Fatalf("sub create returned no id: %s", out)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "sub", "close", subID, "--api", apiID)
	})

	// Accept / reject: either succeed (if still PENDING) or fail with a state
	// transition error from the server. Both are fine. The client-side nil-body
	// bug would surface as a "must not be null" / "required body" error - we
	// assert it doesn't happen.
	assertNoBodySerializationBug(t, "sub accept",
		"apim", "sub", "accept", subID, "--api", apiID, "--reason", "smoke test")

	assertNoBodySerializationBug(t, "sub reject",
		"apim", "sub", "reject", subID, "--api", apiID, "--reason", "smoke test")
}

// TestAPIMSubscriptionTransfer covers transferring a subscription from one
// published plan to another.
func TestAPIMSubscriptionTransfer(t *testing.T) {
	apiID, plan1ID := setupPublishedAPI(t, "plan.json")

	// Add a second published plan on the same API.
	planFixture := writeFixture(t, "plan.json")
	out := runCLIExpectSuccess(t, "apim", "plan", "create",
		"--api", apiID,
		"-f", planFixture,
		"-o", "json")

	plan2ID := extractID(t, out)
	if plan2ID == "" {
		t.Fatalf("plan2 create returned no id: %s", out)
	}

	runCLIExpectSuccess(t, "apim", "plan", "publish", plan2ID, "--api", apiID)

	appID := setupApplication(t)

	subOut := runCLIExpectSuccess(t, "apim", "sub", "create",
		"--api", apiID,
		"--plan", plan1ID,
		"--app", appID,
		"-o", "json")

	subID := extractID(t, subOut)
	if subID == "" {
		t.Fatalf("sub create returned no id: %s", subOut)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "sub", "close", subID, "--api", apiID)
	})

	runCLIExpectSuccess(t, "apim", "sub", "transfer", subID, "--api", apiID, "--plan", plan2ID)
}

// setupPublishedAPI creates an API from api.json, a plan from the given fixture,
// publishes the plan, deploys and starts the API. Cleanup (stop + delete with
// --close-plans) is registered via t.Cleanup.
//
// Returns the API ID and the plan ID.
func setupPublishedAPI(t *testing.T, planFixtureName string) (apiID, planID string) {
	t.Helper()

	apiFixture := writeFixture(t, "api.json")
	planFixture := writeFixture(t, planFixtureName)

	apiOut := runCLIExpectSuccess(t, "apim", "api", "create", "-f", apiFixture, "-o", "json")

	apiID = extractID(t, apiOut)
	if apiID == "" {
		t.Fatalf("api create returned no id: %s", apiOut)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "api", "stop", apiID)
		_, _ = runCLI("apim", "api", "delete", apiID, "--close-plans")
	})

	planOut := runCLIExpectSuccess(t, "apim", "plan", "create",
		"--api", apiID,
		"-f", planFixture,
		"-o", "json")

	planID = extractID(t, planOut)
	if planID == "" {
		t.Fatalf("plan create returned no id: %s", planOut)
	}

	runCLIExpectSuccess(t, "apim", "plan", "publish", planID, "--api", apiID)
	runCLIExpectSuccess(t, "apim", "api", "deploy", apiID)
	runCLIExpectSuccess(t, "apim", "api", "start", apiID)

	return apiID, planID
}

// setupApplication creates an application from app.json and registers its
// deletion via t.Cleanup.
func setupApplication(t *testing.T) string {
	t.Helper()

	appFixture := writeFixture(t, "app.json")

	out := runCLIExpectSuccess(t, "apim", "app", "create", "-f", appFixture, "-o", "json")

	appID := extractID(t, out)
	if appID == "" {
		t.Fatalf("app create returned no id: %s", out)
	}

	t.Cleanup(func() {
		_, _ = runCLI("apim", "app", "delete", appID)
	})

	return appID
}

// assertNoBodySerializationBug runs a CLI command, tolerates any failure (the
// operation may succeed or fail for legitimate reasons), but fails the test if
// the error contains a "must not be null" or "required body" signature that
// would indicate a client-side nil-body serialization bug.
func assertNoBodySerializationBug(t *testing.T, name string, args ...string) {
	t.Helper()

	out, _ := runCLI(args...)

	lower := strings.ToLower(out)
	if strings.Contains(lower, "must not be null") || strings.Contains(lower, "required body") {
		t.Errorf("%s: client nil-body serialization bug detected\noutput: %s", name, out)
	}
}
