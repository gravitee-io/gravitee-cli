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

func TestDomainCRUD(t *testing.T) {
	t.Run("list domains (may be empty on fresh AM)", func(t *testing.T) {
		// On a fresh AM instance there are no domains - just verify the command succeeds.
		runCLIExpectSuccess(t, "am", "domain", "list")
	})

	var domainID string

	t.Run("create domain", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "create",
			"--name", "e2e-test-domain",
			"--description", "E2E test",
			"-o", "json")
		domainID = extractID(t, out)

		if domainID == "" {
			t.Fatal("domain ID is empty")
		}
	})

	t.Run("get domain", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "get", domainID)
		if !strings.Contains(out, "e2e-test-domain") {
			t.Errorf("expected domain name in output, got: %s", out)
		}
	})

	t.Run("get domain json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "get", domainID, "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		if obj["name"] != "e2e-test-domain" {
			t.Errorf("expected name 'e2e-test-domain', got %v", obj["name"])
		}
	})

	t.Run("update domain", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "update", domainID, "--name", "e2e-renamed")
		if !strings.Contains(out, "e2e-renamed") {
			t.Errorf("expected updated name, got: %s", out)
		}
	})

	t.Run("disable domain", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "disable", domainID)
		if !strings.Contains(out, "disabled") {
			t.Errorf("expected 'disabled' message, got: %s", out)
		}
	})

	t.Run("enable domain", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "enable", domainID)
		if !strings.Contains(out, "enabled") {
			t.Errorf("expected 'enabled' message, got: %s", out)
		}
	})

	t.Run("delete domain", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "delete", domainID)
		if !strings.Contains(out, "deleted") {
			t.Errorf("expected 'deleted' message, got: %s", out)
		}
	})
}

func TestDomainEdgeCases(t *testing.T) {
	t.Run("page 0 returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "domain", "list", "--page", "0")
		if !strings.Contains(out, "--page must be >= 1") {
			t.Errorf("expected page validation error, got: %s", out)
		}
	})

	t.Run("per-page 0 returns error", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "domain", "list", "--per-page", "0")
		if !strings.Contains(out, "--per-page must be >= 1") {
			t.Errorf("expected per-page validation error, got: %s", out)
		}
	})

	t.Run("create without name fails", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "domain", "create")
		if !strings.Contains(out, "required") {
			t.Errorf("expected required flag error, got: %s", out)
		}
	})

	t.Run("update without flags fails", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "domain", "update", "fake-id")
		if !strings.Contains(out, "at least one flag") {
			t.Errorf("expected at-least-one-flag error, got: %s", out)
		}
	})

	t.Run("get nonexistent domain fails", func(t *testing.T) {
		_, err := runCLI("am", "domain", "get", "nonexistent-domain-id")
		if err == nil {
			t.Error("expected error for nonexistent domain")
		}
	})
}

func TestAppCRUD(t *testing.T) {
	domainID := getDefaultDomainID(t)

	var appID string

	t.Run("create app", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "create",
			"--domain", domainID,
			"--name", "e2e-app",
			"--type", "service",
			"-o", "json")
		appID = extractID(t, out)
	})

	t.Run("list apps", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "list", "--domain", domainID)
		if !strings.Contains(out, "e2e-app") {
			t.Errorf("expected app in list, got: %s", out)
		}
	})

	t.Run("get app", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "get", "--domain", domainID, appID)
		if !strings.Contains(out, "e2e-app") {
			t.Errorf("expected app name, got: %s", out)
		}
	})

	t.Run("update app", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "update", "--domain", domainID, appID, "--name", "e2e-app-renamed")
		if !strings.Contains(out, "e2e-app-renamed") {
			t.Errorf("expected updated name, got: %s", out)
		}
	})

	t.Run("delete app", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "delete", "--domain", domainID, appID)
	})

	t.Run("invalid type rejected", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "create",
			"--domain", domainID,
			"--name", "x",
			"--type", "invalid")
		if !strings.Contains(out, "invalid value") {
			t.Errorf("expected type validation error, got: %s", out)
		}
	})
}

func TestUserCRUD(t *testing.T) {
	domainID := getDefaultDomainID(t)

	var userID string

	t.Run("create user", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "create",
			"--domain", domainID,
			"--username", "e2e-user",
			"--email", "e2e@test.com",
			"--password", "E2eTestPassword123!@#",
			"-o", "json")
		userID = extractID(t, out)
	})

	t.Run("get user", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "get", "--domain", domainID, userID)
		if !strings.Contains(out, "e2e-user") {
			t.Errorf("expected username, got: %s", out)
		}
	})

	t.Run("lock user", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "lock", "--domain", domainID, userID)
		if !strings.Contains(out, "locked") {
			t.Errorf("expected locked message, got: %s", out)
		}
	})

	t.Run("unlock user", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "unlock", "--domain", domainID, userID)
		if !strings.Contains(out, "unlocked") {
			t.Errorf("expected unlocked message, got: %s", out)
		}
	})

	t.Run("delete user", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	})
}

func TestRoleScopeGroupCRUD(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("role CRUD", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "role", "create",
			"--domain", domainID, "--name", "e2e-role", "-o", "json")
		roleID := extractID(t, out)

		runCLIExpectSuccess(t, "am", "role", "get", "--domain", domainID, roleID)
		runCLIExpectSuccess(t, "am", "role", "update", "--domain", domainID, roleID, "--name", "e2e-role-updated")
		runCLIExpectSuccess(t, "am", "role", "delete", "--domain", domainID, roleID)
	})

	t.Run("scope CRUD", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "scope", "create",
			"--domain", domainID, "--key", "e2e_scope", "--name", "E2E Scope", "--description", "E2E test scope", "-o", "json")
		scopeID := extractID(t, out)

		runCLIExpectSuccess(t, "am", "scope", "get", "--domain", domainID, scopeID)
		runCLIExpectSuccess(t, "am", "scope", "update", "--domain", domainID, scopeID, "--name", "E2E Updated")
		runCLIExpectSuccess(t, "am", "scope", "delete", "--domain", domainID, scopeID)
	})

	t.Run("group CRUD", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "group", "create",
			"--domain", domainID, "--name", "e2e-group", "-o", "json")
		groupID := extractID(t, out)

		runCLIExpectSuccess(t, "am", "group", "get", "--domain", domainID, groupID)
		runCLIExpectSuccess(t, "am", "group", "update", "--domain", domainID, groupID, "--name", "e2e-group-updated")
		runCLIExpectSuccess(t, "am", "group", "delete", "--domain", domainID, groupID)
	})
}

func TestOutputFormats(t *testing.T) {
	t.Run("domain list json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "list", "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON output: %v", err)
		}

		if _, ok := obj["data"]; !ok {
			t.Error("expected 'data' key in JSON output")
		}
	})

	t.Run("domain list yaml", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "list", "-o", "yaml")
		if !strings.Contains(out, "data:") {
			t.Errorf("expected YAML with 'data:', got: %s", out)
		}
	})

	t.Run("domain list quiet", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "list", "-q")
		if out != "" {
			t.Errorf("expected empty output in quiet mode, got: %s", out)
		}
	})

	t.Run("domain list no-headers", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "list", "--no-headers")
		if strings.Contains(out, "NAME") {
			t.Error("expected no headers but found NAME")
		}
	})

	t.Run("invalid format rejected", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "domain", "list", "-o", "xml")
		if !strings.Contains(out, "invalid output format") {
			t.Errorf("expected format error, got: %s", out)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("nonexistent context", func(t *testing.T) {
		// Isolate HOME and clear the AM env vars so the CLI falls back to the
		// config file (and reports the missing context) instead of bypassing
		// resolution via GCTL_AM_URL/GCTL_AM_TOKEN set in TestMain.
		t.Setenv("HOME", t.TempDir())
		t.Setenv("GCTL_AM_URL", "")
		t.Setenv("GCTL_AM_TOKEN", "")

		out := runCLIExpectError(t, "am", "domain", "list", "--context", "nonexistent")
		if !strings.Contains(out, "not found") {
			t.Errorf("expected context not found error, got: %s", out)
		}
	})

	t.Run("missing domain flag", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "list")
		if !strings.Contains(out, "required") {
			t.Errorf("expected required flag error, got: %s", out)
		}
	})
}
