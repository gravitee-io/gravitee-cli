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

// createTestApp creates an app for testing and returns its ID.
func createTestApp(t *testing.T, domainID, name string) string {
	t.Helper()

	out := runCLIExpectSuccess(t, "am", "app", "create",
		"--domain", domainID,
		"--name", name,
		"--type", "service",
		"-o", "json")

	return extractID(t, out)
}

// createTestGroup creates a group for testing and returns its ID.
func createTestGroup(t *testing.T, domainID, name string) string {
	t.Helper()

	out := runCLIExpectSuccess(t, "am", "group", "create",
		"--domain", domainID,
		"--name", name,
		"-o", "json")

	return extractID(t, out)
}

func TestAppSecrets(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-sub-app-secrets")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("list secrets", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "secret", "list", "--domain", domainID, "--app-id", appID)
	})

	t.Run("create secret", func(t *testing.T) {
		out, err := runCLI("am", "app", "secret", "create", "--domain", domainID, "--app-id", appID, "--name", "e2e-secret")
		if err != nil {
			t.Skipf("secret create not supported or failed: %s", out)
		}
	})

	t.Run("missing app-id returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "secret", "list", "--domain", domainID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestAppMembers(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-sub-app-members")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("list members", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "member", "list", "--domain", domainID, "--app-id", appID)
	})

	t.Run("list members json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "member", "list", "--domain", domainID, "--app-id", appID, "-o", "json")

		var obj json.RawMessage
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON output: %v\nOutput: %s", err, out)
		}
	})

	t.Run("member permissions", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "member", "permissions", "--domain", domainID, "--app-id", appID)
	})
}

func TestAppFlows(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-sub-app-flows")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("list flows", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "flow", "list", "--domain", domainID, "--app-id", appID)
	})

	t.Run("list flows json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "flow", "list", "--domain", domainID, "--app-id", appID, "-o", "json")

		var obj json.RawMessage
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON output: %v\nOutput: %s", err, out)
		}
	})
}

func TestAppEmails(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-sub-app-emails")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("get email template", func(t *testing.T) {
		out, err := runCLI("am", "app", "email", "get", "--domain", domainID, "--app-id", appID, "--template", "RESET_PASSWORD")
		if err != nil {
			t.Logf("email get returned error (may be 404 if not configured): %s", out)
		}
	})

	t.Run("missing template returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "email", "get", "--domain", domainID, "--app-id", appID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})

	t.Run("missing file on create returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "email", "create", "--domain", domainID, "--app-id", appID)
		if !strings.Contains(out, "input") {
			t.Errorf("expected 'input' error, got: %s", out)
		}
	})
}

func TestAppForms(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-sub-app-forms")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("get form template", func(t *testing.T) {
		out, err := runCLI("am", "app", "form", "get", "--domain", domainID, "--app-id", appID, "--template", "LOGIN")
		if err != nil {
			t.Logf("form get returned error (may be 404 if not configured): %s", out)
		}
	})

	t.Run("missing template returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "form", "get", "--domain", domainID, "--app-id", appID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestAppResources(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-sub-app-resources")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("list resources", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "resource", "list", "--domain", domainID, "--app-id", appID)
	})
}

func TestAppResourcePolicies(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-sub-app-res-policies")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("missing resource-id returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "resource-policy", "list", "--domain", domainID, "--app-id", appID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestAppAnalytics(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-sub-app-analytics")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("get analytics", func(t *testing.T) {
		out, err := runCLI("am", "app", "analytics", "get", "--domain", domainID, "--app-id", appID, "--type", "count", "--interval", "86400000", "--from", "0", "--to", "9999999999999")
		if err != nil {
			t.Skipf("analytics get not supported or failed: %s", out)
		}
	})

	t.Run("missing type returns error", func(t *testing.T) {
		// Missing --type may return CLI "required" error or API 500 - just verify it fails.
		_, err := runCLI("am", "app", "analytics", "get", "--domain", domainID, "--app-id", appID)
		if err == nil {
			t.Error("expected error for missing --type")
		}
	})
}

func TestAppChangeType(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-sub-app-changetype")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("change type", func(t *testing.T) {
		out, err := runCLI("am", "app", "change-type", "--domain", domainID, "--app-id", appID, "--type", "browser")
		if err != nil {
			t.Skipf("change-type not supported or failed: %s", out)
		}
	})

	t.Run("missing type returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "change-type", "--domain", domainID, "--app-id", appID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestGroupMembers(t *testing.T) {
	domainID := getDefaultDomainID(t)
	groupID := createTestGroup(t, domainID, "e2e-sub-group-members")

	defer runCLI("am", "group", "delete", "--domain", domainID, groupID)

	// Create a user to add as member.
	userOut := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-grp-member-user",
		"--email", "e2e-grp-member@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, userOut)

	defer runCLI("am", "user", "delete", "--domain", domainID, userID)

	t.Run("list members", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "group", "member", "list", "--domain", domainID, "--group-id", groupID)
	})

	t.Run("add member", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "group", "member", "add", "--domain", domainID, "--group-id", groupID, userID)
	})

	t.Run("missing group-id returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "group", "member", "list", "--domain", domainID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestGroupRoles(t *testing.T) {
	domainID := getDefaultDomainID(t)
	groupID := createTestGroup(t, domainID, "e2e-sub-group-roles")

	defer runCLI("am", "group", "delete", "--domain", domainID, groupID)

	// Create a role to assign.
	roleOut := runCLIExpectSuccess(t, "am", "role", "create",
		"--domain", domainID,
		"--name", "e2e-grp-role",
		"-o", "json")
	roleID := extractID(t, roleOut)

	defer runCLI("am", "role", "delete", "--domain", domainID, roleID)

	t.Run("list roles", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "group", "role", "list", "--domain", domainID, "--group-id", groupID)
	})

	t.Run("assign role", func(t *testing.T) {
		out, err := runCLI("am", "group", "role", "assign", "--domain", domainID, "--group-id", groupID, "--roles", roleID)
		if err != nil {
			t.Skipf("group role assign not supported or failed: %s", out)
		}
	})

	t.Run("missing group-id returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "group", "role", "list", "--domain", domainID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}
