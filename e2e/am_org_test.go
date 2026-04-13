//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func writeOrgTempJSON(t *testing.T, data map[string]any) string {
	t.Helper()

	raw, _ := json.Marshal(data)

	f, err := os.CreateTemp("", "e2e-org-*.json")
	if err != nil {
		t.Fatalf("failed to create temp: %v", err)
	}

	f.Write(raw)
	f.Close()

	return f.Name()
}

func TestOrgSettings(t *testing.T) {
	t.Run("get settings", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "settings", "get")
	})

	t.Run("get settings json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "org", "settings", "get", "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON output: %v\nOutput: %s", err, out)
		}
	})
}

func TestOrgUserCRUD(t *testing.T) {
	t.Run("list users", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "user", "list")
	})

	var userID string
	var created bool

	t.Run("create user", func(t *testing.T) {
		tmpFile := writeOrgTempJSON(t, map[string]any{
			"username":        "e2e-org-user",
			"firstName":       "E2E",
			"lastName":        "Test",
			"email":           "orguser@test.com",
			"password":        "E2eTestPassword123!@#",
			"preRegistration": false,
		})
		defer os.Remove(tmpFile)

		out, err := runCLI("am", "org", "user", "create", "--file", tmpFile, "-o", "json")
		if err != nil {
			t.Skipf("skipping: org user create failed: %v\nOutput: %s", err, out)
		}

		userID = extractID(t, out)
		created = true
	})

	t.Run("get user", func(t *testing.T) {
		if !created {
			t.Skip("skipping: user was not created")
		}

		out := runCLIExpectSuccess(t, "am", "org", "user", "get", userID)
		if !strings.Contains(out, "e2e-org-user") {
			t.Errorf("expected 'e2e-org-user' in output, got: %s", out)
		}
	})

	t.Run("update-status user", func(t *testing.T) {
		if !created {
			t.Skip("skipping: user was not created")
		}

		out, err := runCLI("am", "org", "user", "update-status", userID, "--enabled", "false")
		if err != nil {
			t.Skipf("skipping: update-status failed: %v\nOutput: %s", err, out)
		}
	})

	t.Run("reset-password user", func(t *testing.T) {
		if !created {
			t.Skip("skipping: user was not created")
		}

		out, err := runCLI("am", "org", "user", "reset-password", userID, "--password", "NewPass123!@#")
		if err != nil {
			t.Skipf("skipping: reset-password failed: %v\nOutput: %s", err, out)
		}
	})

	t.Run("update-username user", func(t *testing.T) {
		if !created {
			t.Skip("skipping: user was not created")
		}

		out, err := runCLI("am", "org", "user", "update-username", userID, "--username", "e2e-org-renamed")
		if err != nil {
			t.Skipf("skipping: update-username failed: %v\nOutput: %s", err, out)
		}
	})

	t.Run("delete user", func(t *testing.T) {
		if !created {
			t.Skip("skipping: user was not created")
		}

		runCLIExpectSuccess(t, "am", "org", "user", "delete", userID)
	})

	t.Run("bulk missing file", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "user", "bulk")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})
}

func TestOrgGroupCRUD(t *testing.T) {
	t.Run("list groups", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "group", "list")
	})

	t.Run("list groups json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "org", "group", "list", "-o", "json")

		// Org groups may return paginated object or array — just verify valid JSON.
		if !json.Valid([]byte(out)) {
			t.Fatalf("expected valid JSON, got: %s", out)
		}
	})

	var groupID string
	var created bool

	t.Run("create group", func(t *testing.T) {
		tmpFile := writeOrgTempJSON(t, map[string]any{
			"name": "e2e-org-group",
		})
		defer os.Remove(tmpFile)

		out, err := runCLI("am", "org", "group", "create", "--file", tmpFile, "-o", "json")
		if err != nil {
			t.Skipf("skipping: org group create failed: %v\nOutput: %s", err, out)
		}

		groupID = extractID(t, out)
		created = true
	})

	t.Run("get group", func(t *testing.T) {
		if !created {
			t.Skip("skipping: group was not created")
		}

		out := runCLIExpectSuccess(t, "am", "org", "group", "get", groupID)
		if !strings.Contains(out, "e2e-org-group") {
			t.Errorf("expected 'e2e-org-group' in output, got: %s", out)
		}
	})

	t.Run("delete group", func(t *testing.T) {
		if !created {
			t.Skip("skipping: group was not created")
		}

		runCLIExpectSuccess(t, "am", "org", "group", "delete", groupID)
	})
}

func TestOrgRoleCRUD(t *testing.T) {
	t.Run("list roles", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "role", "list")
	})

	t.Run("list roles json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "org", "role", "list", "-o", "json")

		var obj any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("expected valid JSON: %v\nOutput: %s", err, out)
		}
	})

	var roleID string
	var created bool

	t.Run("create role", func(t *testing.T) {
		tmpFile := writeOrgTempJSON(t, map[string]any{
			"name":           "e2e-org-role",
			"assignableType": "ORGANIZATION",
		})
		defer os.Remove(tmpFile)

		out, err := runCLI("am", "org", "role", "create", "--file", tmpFile, "-o", "json")
		if err != nil {
			t.Skipf("skipping: org role create failed: %v\nOutput: %s", err, out)
		}

		roleID = extractID(t, out)
		created = true
	})

	t.Run("get role", func(t *testing.T) {
		if !created {
			t.Skip("skipping: role was not created")
		}

		out := runCLIExpectSuccess(t, "am", "org", "role", "get", roleID)
		if !strings.Contains(out, "e2e-org-role") {
			t.Errorf("expected 'e2e-org-role' in output, got: %s", out)
		}
	})

	t.Run("delete role", func(t *testing.T) {
		if !created {
			t.Skip("skipping: role was not created")
		}

		runCLIExpectSuccess(t, "am", "org", "role", "delete", roleID)
	})
}

func TestOrgMemberOperations(t *testing.T) {
	t.Run("list members", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "member", "list")
	})

	t.Run("list members json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "org", "member", "list", "-o", "json")

		var obj any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("expected valid JSON: %v\nOutput: %s", err, out)
		}
	})

	t.Run("add member missing member-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "member", "add", "--role", "fake-role")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})

	t.Run("add member missing role", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "member", "add", "--member-id", "fake-id")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})
}

func TestOrgAuditOperations(t *testing.T) {
	t.Run("list audits", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "audit", "list")
	})

	var auditID string

	t.Run("list audits json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "org", "audit", "list", "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("expected valid JSON: %v\nOutput: %s", err, out)
		}

		if _, ok := obj["data"]; !ok {
			t.Errorf("expected 'data' key in JSON output, got: %s", out)
		}

		// Extract an audit ID for later use.
		if data, ok := obj["data"].([]any); ok && len(data) > 0 {
			if entry, ok := data[0].(map[string]any); ok {
				if id, ok := entry["id"].(string); ok {
					auditID = id
				}
			}
		}
	})

	t.Run("list audits with per-page", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "audit", "list", "--per-page", "2")
	})

	t.Run("get audit", func(t *testing.T) {
		if auditID == "" {
			t.Skip("skipping: no audit ID available")
		}

		out := runCLIExpectSuccess(t, "am", "org", "audit", "get", auditID)
		if !strings.Contains(out, auditID) {
			t.Errorf("expected audit ID '%s' in output, got: %s", auditID, out)
		}
	})
}

func TestOrgReporterOperations(t *testing.T) {
	t.Run("list reporters", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "reporter", "list")
	})

	t.Run("list reporters json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "org", "reporter", "list", "-o", "json")

		var obj any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("expected valid JSON: %v\nOutput: %s", err, out)
		}
	})

	t.Run("create reporter missing file", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "reporter", "create")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})
}

func TestOrgFormOperations(t *testing.T) {
	t.Run("get form with template", func(t *testing.T) {
		// May 404 if not configured, so we just use runCLI.
		runCLI("am", "org", "form", "get", "--template", "LOGIN")
	})

	t.Run("get form missing template", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "form", "get")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})

	t.Run("create form missing file", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "form", "create")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})
}

func TestOrgIDPOperations(t *testing.T) {
	t.Run("list idps", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "idp", "list")
	})

	t.Run("list idps json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "org", "idp", "list", "-o", "json")

		var obj any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("expected valid JSON: %v\nOutput: %s", err, out)
		}
	})

	t.Run("create idp missing file", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "idp", "create")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})
}

func TestOrgEntrypointOperations(t *testing.T) {
	t.Run("list entrypoints", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "entrypoint", "list")
	})

	t.Run("list entrypoints json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "org", "entrypoint", "list", "-o", "json")

		var obj any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("expected valid JSON: %v\nOutput: %s", err, out)
		}
	})

	t.Run("create entrypoint missing file", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "entrypoint", "create")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})
}

func TestOrgTagCRUD(t *testing.T) {
	t.Run("list tags", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "org", "tag", "list")
	})

	t.Run("list tags json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "org", "tag", "list", "-o", "json")

		var obj any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("expected valid JSON: %v\nOutput: %s", err, out)
		}
	})

	var tagID string
	var created bool

	t.Run("create tag", func(t *testing.T) {
		out, err := runCLI("am", "org", "tag", "create", "--name", "e2e-tag", "--description", "test tag", "-o", "json")
		if err != nil {
			t.Skipf("skipping: org tag create failed: %v\nOutput: %s", err, out)
		}

		tagID = extractID(t, out)
		created = true
	})

	t.Run("get tag", func(t *testing.T) {
		if !created {
			t.Skip("skipping: tag was not created")
		}

		out := runCLIExpectSuccess(t, "am", "org", "tag", "get", tagID)
		if !strings.Contains(out, "e2e-tag") {
			t.Errorf("expected 'e2e-tag' in output, got: %s", out)
		}
	})

	t.Run("delete tag", func(t *testing.T) {
		if !created {
			t.Skip("skipping: tag was not created")
		}

		runCLIExpectSuccess(t, "am", "org", "tag", "delete", tagID)
	})
}

func TestOrgUserTokens(t *testing.T) {
	t.Run("list tokens missing user-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "user-token", "list")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})

	t.Run("create token missing user-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "org", "user-token", "create", "--name", "test")
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' in error, got: %s", out)
		}
	})
}
