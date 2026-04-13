//go:build e2e

package e2e

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPaginationEdgeCases(t *testing.T) {
	domainID := getDefaultDomainID(t)

	// Create 3 scopes for pagination testing.
	var scopeIDs [3]string

	for i, sc := range []struct {
		key, name string
	}{
		{"e2e_page_1", "Page1"},
		{"e2e_page_2", "Page2"},
		{"e2e_page_3", "Page3"},
	} {
		out := runCLIExpectSuccess(t, "am", "scope", "create",
			"--domain", domainID,
			"--key", sc.key,
			"--name", sc.name,
			"--description", "test",
			"-o", "json")
		scopeIDs[i] = extractID(t, out)
	}

	// Ensure cleanup runs even if a subtest fails.
	defer func() {
		for _, id := range scopeIDs {
			if id != "" {
				runCLI("am", "scope", "delete", "--domain", domainID, id)
			}
		}
	}()

	t.Run("partial page with per-page 2", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "scope", "list", "--domain", domainID, "--per-page", "2")
		if !strings.Contains(out, "Showing 2 of") {
			t.Errorf("expected 'Showing 2 of' in output, got: %s", out)
		}
	})

	t.Run("second page", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "scope", "list", "--domain", domainID, "--per-page", "2", "--page", "2")
	})

	t.Run("all flag shows all scopes", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "scope", "list", "--domain", domainID, "--all")
		for _, name := range []string{"Page1", "Page2", "Page3"} {
			if !strings.Contains(out, name) {
				t.Errorf("expected %q in output, got: %s", name, out)
			}
		}
	})

	t.Run("large per-page shows all scopes", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "scope", "list", "--domain", domainID, "--per-page", "1000")
		for _, name := range []string{"Page1", "Page2", "Page3"} {
			if !strings.Contains(out, name) {
				t.Errorf("expected %q in output, got: %s", name, out)
			}
		}
	})

	t.Run("negative page returns error", func(t *testing.T) {
		runCLIExpectError(t, "am", "scope", "list", "--domain", domainID, "--page", "-1")
	})

	t.Run("negative per-page returns error", func(t *testing.T) {
		runCLIExpectError(t, "am", "scope", "list", "--domain", domainID, "--per-page", "-1")
	})
}

func TestUserSCIMFilter(t *testing.T) {
	domainID := getDefaultDomainID(t)

	// Create 2 users for filter testing.
	out1 := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-filter-user1",
		"--email", "filter1@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID1 := extractID(t, out1)

	out2 := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-filter-user2",
		"--email", "filter2@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID2 := extractID(t, out2)

	defer func() {
		if userID1 != "" {
			runCLI("am", "user", "delete", "--domain", domainID, userID1)
		}
		if userID2 != "" {
			runCLI("am", "user", "delete", "--domain", domainID, userID2)
		}
	}()

	t.Run("SCIM filter returns only matching user", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "list",
			"--domain", domainID,
			"--filter", `userName eq "e2e-filter-user1"`)
		if !strings.Contains(out, "e2e-filter-user1") {
			t.Errorf("expected e2e-filter-user1 in output, got: %s", out)
		}
		if strings.Contains(out, "e2e-filter-user2") {
			t.Errorf("did not expect e2e-filter-user2 in output, got: %s", out)
		}
	})

	t.Run("query flag is accepted", func(t *testing.T) {
		// AM API may not support --query for users; just verify the command doesn't error.
		runCLIExpectSuccess(t, "am", "user", "list",
			"--domain", domainID,
			"--query", "e2e-filter")
	})
}

func TestDebugOutput(t *testing.T) {
	out := runCLIExpectSuccess(t, "am", "domain", "list", "--debug")

	if !strings.Contains(out, "> GET") {
		t.Errorf("expected debug request log '> GET' in output, got: %s", out)
	}
	if !strings.Contains(out, "< HTTP") {
		t.Errorf("expected debug response log '< HTTP' in output, got: %s", out)
	}
}

func TestQueryFlags(t *testing.T) {
	domainID := getDefaultDomainID(t)

	// Create a role for query testing.
	out := runCLIExpectSuccess(t, "am", "role", "create",
		"--domain", domainID,
		"--name", "e2e-query-role",
		"-o", "json")
	roleID := extractID(t, out)

	defer func() {
		if roleID != "" {
			runCLI("am", "role", "delete", "--domain", domainID, roleID)
		}
	}()

	t.Run("role list with query succeeds", func(t *testing.T) {
		// AM API search behavior varies, just verify the command succeeds.
		runCLIExpectSuccess(t, "am", "role", "list", "--domain", domainID, "--query", "e2e-query")
	})

	t.Run("scope list with nonexistent query returns empty", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "scope", "list", "--domain", domainID, "--query", "nonexistent")

		// Verify either "No results found" or an empty JSON data array.
		hasNoResults := strings.Contains(out, "No results found")

		var resp struct {
			Data []json.RawMessage `json:"data"`
		}

		emptyJSON := false
		if err := json.Unmarshal([]byte(out), &resp); err == nil {
			emptyJSON = len(resp.Data) == 0
		}

		if !hasNoResults && !emptyJSON {
			// If it's a table output with no data rows, that's also fine.
			// Just check it doesn't contain unexpected scope names.
			if strings.Contains(out, "e2e_page") || strings.Contains(out, "e2e_scope") {
				t.Errorf("expected empty result for nonexistent query, got: %s", out)
			}
		}
	})
}

func TestListAllFlag(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("role list all", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "role", "list", "--domain", domainID, "--all")
	})

	t.Run("user list all", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "list", "--domain", domainID, "--all")
	})

	t.Run("app list all", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "list", "--domain", domainID, "--all")
	})
}
