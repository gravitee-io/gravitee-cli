//go:build e2e

package e2e

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMemberCRUD(t *testing.T) {
	domainID := getDefaultDomainID(t)

	var userID string

	t.Run("create user for member tests", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "create",
			"--domain", domainID,
			"--username", "member-test-user",
			"--email", "member@test.com",
			"--password", "E2eTestPassword123!@#",
			"-o", "json")
		userID = extractID(t, out)

		if userID == "" {
			t.Fatal("user ID is empty")
		}
	})

	t.Run("list members", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "member", "list", "--domain", domainID)
	})

	t.Run("list members json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "member", "list", "--domain", domainID, "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON output: %v\nOutput: %s", err, out)
		}

		if _, ok := obj["memberships"]; !ok {
			t.Errorf("expected 'memberships' key in JSON output, got: %s", out)
		}
	})

	t.Run("cleanup user", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	})
}

func TestAuditOperations(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("list audits", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "audit", "list", "--domain", domainID)
	})

	t.Run("list audits with type filter", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "audit", "list", "--domain", domainID, "--type", "DOMAIN_CREATED")
	})

	t.Run("list audits json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "audit", "list", "--domain", domainID, "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON output: %v\nOutput: %s", err, out)
		}

		if _, ok := obj["data"]; !ok {
			t.Errorf("expected 'data' key in JSON output, got: %s", out)
		}
	})

	t.Run("list audits with pagination", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "audit", "list", "--domain", domainID, "--per-page", "2")
	})

	t.Run("list all audits", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "audit", "list", "--domain", domainID, "--all")
	})
}

func TestFlowOperations(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("list flows", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "flow", "list", "--domain", domainID)
	})

	t.Run("list flows json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "flow", "list", "--domain", domainID, "-o", "json")

		var arr []any
		if err := json.Unmarshal([]byte(out), &arr); err != nil {
			t.Fatalf("expected valid JSON array: %v\nOutput: %s", err, out)
		}
	})
}

func TestEntrypoint(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("get entrypoint", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "entrypoint", "get", "--domain", domainID)

		if !strings.Contains(out, "{") {
			t.Errorf("expected JSON output, got: %s", out)
		}

		// Entrypoint returns a JSON array, not an object.
		var arr []any
		if err := json.Unmarshal([]byte(out), &arr); err != nil {
			t.Fatalf("invalid JSON output: %v\nOutput: %s", err, out)
		}
	})
}

func TestIDPList(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("list idps", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "idp", "list", "--domain", domainID)
	})

	t.Run("list idps json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "idp", "list", "--domain", domainID, "-o", "json")

		var arr []any
		if err := json.Unmarshal([]byte(out), &arr); err != nil {
			t.Fatalf("expected valid JSON array: %v\nOutput: %s", err, out)
		}
	})
}
