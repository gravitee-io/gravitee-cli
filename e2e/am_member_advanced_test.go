//go:build e2e

package e2e

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMemberAddRemove(t *testing.T) {
	domainID := getDefaultDomainID(t)

	// Member add requires an org-level user and an org-level role with assignableType=domain.
	// 1. Get the org admin user ID.
	orgUserOut := runCLIExpectSuccess(t, "am", "org", "user", "list", "-o", "json")

	var orgUsers struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(orgUserOut), &orgUsers); err != nil || len(orgUsers.Data) == 0 {
		t.Fatal("no org users available for member test")
	}

	orgUserID := orgUsers.Data[0].ID

	// 2. Get DOMAIN_USER role (org-level role with assignableType=domain).
	orgRoleOut := runCLIExpectSuccess(t, "am", "org", "role", "list", "-o", "json")

	var domainUserRoleID string

	var orgRoles []map[string]any
	if err := json.Unmarshal([]byte(orgRoleOut), &orgRoles); err != nil {
		// Try paginated format.
		var paged struct {
			Data []map[string]any `json:"data"`
		}
		if err2 := json.Unmarshal([]byte(orgRoleOut), &paged); err2 == nil {
			orgRoles = paged.Data
		}
	}

	for _, r := range orgRoles {
		name, _ := r["name"].(string)
		assignable, _ := r["assignableType"].(string)
		if name == "DOMAIN_USER" && assignable == "domain" {
			domainUserRoleID, _ = r["id"].(string)
			break
		}
	}

	if domainUserRoleID == "" {
		t.Fatal("DOMAIN_USER role not found in org roles")
	}

	// 3. Add member.
	t.Run("add member", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "member", "add",
			"--domain", domainID,
			"--member-id", orgUserID,
			"--role", domainUserRoleID)

		lower := strings.ToLower(out)
		if !strings.Contains(lower, "added") && !strings.Contains(lower, "member") {
			t.Errorf("expected success message, got: %s", out)
		}
	})

	// 4. List members and verify user is present.
	t.Run("list members contains user", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "member", "list", "--domain", domainID)

		if !strings.Contains(out, orgUserID) {
			t.Errorf("expected member list to contain user ID %s, got: %s", orgUserID, out)
		}
	})

	// 5. List members JSON and find membership ID.
	var membershipID string

	t.Run("list members json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "member", "list", "--domain", domainID, "-o", "json")

		var resp struct {
			Memberships []struct {
				ID       string `json:"id"`
				MemberID string `json:"memberId"`
			} `json:"memberships"`
		}

		if err := json.Unmarshal([]byte(out), &resp); err != nil {
			t.Fatalf("failed to parse memberships JSON: %v\nOutput: %s", err, out)
		}

		for _, m := range resp.Memberships {
			if m.MemberID == orgUserID {
				membershipID = m.ID
				break
			}
		}

		if membershipID == "" {
			t.Fatalf("could not find membership for user %s", orgUserID)
		}
	})

	// 6. Remove member.
	t.Run("remove member", func(t *testing.T) {
		if membershipID == "" {
			t.Fatal("no membership ID")
		}

		runCLIExpectSuccess(t, "am", "member", "remove", "--domain", domainID, membershipID)
	})

	// 7. Validation tests.
	t.Run("add member missing member-id", func(t *testing.T) {
		runCLIExpectError(t, "am", "member", "add",
			"--domain", domainID, "--role", domainUserRoleID)
	})

	t.Run("add member missing role", func(t *testing.T) {
		runCLIExpectError(t, "am", "member", "add",
			"--domain", domainID, "--member-id", orgUserID)
	})

	t.Run("remove member missing id", func(t *testing.T) {
		runCLIExpectError(t, "am", "member", "remove", "--domain", domainID)
	})
}

func TestAuditGet(t *testing.T) {
	domainID := getDefaultDomainID(t)

	var auditID string

	// Step 1-2: List audits and extract first audit ID.
	t.Run("list audits and extract id", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "audit", "list", "--domain", domainID, "-o", "json")

		var resp struct {
			Data []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
			TotalCount  int `json:"totalCount"`
			CurrentPage int `json:"currentPage"`
		}

		if err := json.Unmarshal([]byte(out), &resp); err != nil {
			t.Fatalf("failed to parse audit list JSON: %v\nOutput: %s", err, out)
		}

		if len(resp.Data) == 0 {
			t.Skip("no audit entries available — skipping audit get tests")
		}

		auditID = resp.Data[0].ID

		if auditID == "" {
			t.Fatal("first audit ID is empty")
		}
	})

	// Step 3: Get single audit by ID.
	t.Run("get audit", func(t *testing.T) {
		if auditID == "" {
			t.Skip("no audit ID from previous step")
		}

		out := runCLIExpectSuccess(t, "am", "audit", "get", "--domain", domainID, auditID)

		if !strings.Contains(out, auditID) {
			t.Errorf("expected output to contain audit ID %s, got: %s", auditID, out)
		}
	})

	// Step 4: Get audit as JSON.
	t.Run("get audit json", func(t *testing.T) {
		if auditID == "" {
			t.Skip("no audit ID from previous step")
		}

		out := runCLIExpectSuccess(t, "am", "audit", "get", "--domain", domainID, auditID, "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON output: %v\nOutput: %s", err, out)
		}

		if _, ok := obj["id"]; !ok {
			t.Errorf("expected 'id' field in audit JSON, got: %s", out)
		}
	})

	// Step 5: Get nonexistent audit should fail.
	t.Run("get nonexistent audit", func(t *testing.T) {
		_, err := runCLI("am", "audit", "get", "--domain", domainID, "nonexistent-audit-id")
		if err == nil {
			t.Error("expected error for nonexistent audit ID")
		}
	})
}
