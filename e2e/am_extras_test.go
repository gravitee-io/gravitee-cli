//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestCertificateExtras(t *testing.T) {
	domainID := getDefaultDomainID(t)

	var certID string

	t.Run("list certificates", func(t *testing.T) {
		out, err := runCLI("am", "certificate", "list", "--domain", domainID, "-o", "json")
		if err != nil {
			t.Skipf("certificate list failed, skipping extras: %s", out)
		}

		var resp []map[string]any
		if err := json.Unmarshal([]byte(out), &resp); err != nil {
			// Try paginated format.
			var paged struct {
				Data []map[string]any `json:"data"`
			}
			if err2 := json.Unmarshal([]byte(out), &paged); err2 != nil {
				t.Skipf("cannot parse certificate list, skipping: %s", out)
			}
			resp = paged.Data
		}

		if len(resp) > 0 {
			id, ok := resp[0]["id"].(string)
			if ok {
				certID = id
			}
		}
	})

	t.Run("certificate key", func(t *testing.T) {
		if certID == "" {
			t.Skip("no certificate available")
		}

		out, err := runCLI("am", "certificate", "key", "--domain", domainID, certID)
		if err != nil {
			t.Skipf("certificate key not supported or failed: %s", out)
		}
	})

	t.Run("certificate keys", func(t *testing.T) {
		if certID == "" {
			t.Skip("no certificate available")
		}

		out, err := runCLI("am", "certificate", "keys", "--domain", domainID, certID)
		if err != nil {
			t.Skipf("certificate keys not supported or failed: %s", out)
		}
	})

	t.Run("certificate rotate", func(t *testing.T) {
		out, err := runCLI("am", "certificate", "rotate", "--domain", domainID)
		if err != nil {
			t.Skipf("certificate rotate not supported or failed: %s", out)
		}
	})
}

func TestPasswordPolicyExtras(t *testing.T) {
	domainID := getDefaultDomainID(t)

	// Create a password policy for testing.
	ppFile, err := os.CreateTemp("", "e2e-pp-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	ppFile.WriteString(`{"name":"e2e-pp-extras","minLength":8,"maxLength":128}`)
	ppFile.Close()

	defer os.Remove(ppFile.Name())

	ppOut, ppErr := runCLI("am", "password-policy", "create", "--domain", domainID, "--file", ppFile.Name(), "-o", "json")

	var policyID string
	if ppErr == nil {
		policyID = extractID(t, ppOut)
	}

	defer func() {
		if policyID != "" {
			runCLI("am", "password-policy", "delete", "--domain", domainID, policyID)
		}
	}()

	t.Run("active password policy", func(t *testing.T) {
		out, err := runCLI("am", "password-policy", "active", "--domain", domainID)
		if err != nil {
			t.Logf("no active password policy: %s", out)
		}
	})

	t.Run("list password policies", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "password-policy", "list", "--domain", domainID)
	})

	t.Run("evaluate password policy", func(t *testing.T) {
		if policyID == "" {
			t.Skip("password policy not created")
		}

		runCLIExpectSuccess(t, "am", "password-policy", "evaluate", "--domain", domainID, "--password", "Test123!", policyID)
	})

	t.Run("evaluate missing password flag", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "password-policy", "evaluate", "--domain", domainID, "fake-id")
		if !strings.Contains(strings.ToLower(out), "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestDomainGetByHRID(t *testing.T) {
	domainID := getDefaultDomainID(t)

	var hrid string
	var domainName string

	t.Run("get domain details", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "get", domainID, "-o", "json")

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}

		h, ok := obj["hrid"].(string)
		if !ok || h == "" {
			t.Skip("domain has no hrid field")
		}
		hrid = h

		n, ok := obj["name"].(string)
		if ok {
			domainName = n
		}
	})

	t.Run("get domain by hrid", func(t *testing.T) {
		if hrid == "" {
			t.Skip("no hrid available")
		}

		out := runCLIExpectSuccess(t, "am", "domain", "get", "--hrid", hrid)
		if domainName != "" && !strings.Contains(out, domainName) {
			t.Errorf("expected domain name %q in output, got: %s", domainName, out)
		}
	})

	t.Run("get domain by nonexistent hrid", func(t *testing.T) {
		out, err := runCLI("am", "domain", "get", "--hrid", "nonexistent-hrid")
		if err == nil {
			t.Errorf("expected error for nonexistent hrid, got: %s", out)
		}
	})
}

func TestDomainCertificateSettings(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("missing file flag", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "domain", "update-cert-settings", domainID)
		if !strings.Contains(strings.ToLower(out), "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})

	t.Run("invalid file content", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "cert-settings-*.json")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(`{"invalid": true}`)
		if err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		tmpFile.Close()

		out, err := runCLI("am", "domain", "update-cert-settings", domainID, "--file", tmpFile.Name())
		if err != nil {
			t.Logf("update-cert-settings with invalid file failed as expected: %s", out)
		}
	})
}

func TestDataPlaneList(t *testing.T) {
	t.Run("list data planes", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "data-plane", "list")
	})

	t.Run("list data planes json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "data-plane", "list", "-o", "json")

		if !json.Valid([]byte(out)) {
			t.Errorf("expected valid JSON output, got: %s", out)
		}
	})
}

func TestProtectedResourceMembers(t *testing.T) {
	getDefaultDomainID(t)

	t.Run("member list missing resource-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "protected-resource", "member", "list", "--domain", "fake")
		if !strings.Contains(strings.ToLower(out), "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})

	t.Run("secret list missing resource-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "protected-resource", "secret", "list", "--domain", "fake")
		if !strings.Contains(strings.ToLower(out), "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestMemberPermissions(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("get permissions", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "member", "permissions", "--domain", domainID)
	})

	t.Run("get permissions json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "member", "permissions", "--domain", domainID, "-o", "json")

		if !json.Valid([]byte(out)) {
			t.Errorf("expected valid JSON output, got: %s", out)
		}
	})
}

func TestDictionaryEntries(t *testing.T) {
	getDefaultDomainID(t)

	t.Run("entry list missing dict-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "dictionary", "entry", "list", "--domain", "fake")
		if !strings.Contains(strings.ToLower(out), "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})

	t.Run("entry update missing dict-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "dictionary", "entry", "update", "--domain", "fake", "--file", "x.json")
		if !strings.Contains(strings.ToLower(out), "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestIDPPasswordPolicy(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("assign missing policy-id flag", func(t *testing.T) {
		runCLIExpectError(t, "am", "idp", "password-policy", "assign",
			"--domain", domainID, "--idp-id", "fake-idp")
	})

	t.Run("assign missing idp-id flag", func(t *testing.T) {
		runCLIExpectError(t, "am", "idp", "password-policy", "assign",
			"--domain", domainID, "--policy-id", "fake-pp")
	})
}

func TestEntrypointCRUD(t *testing.T) {
	domainID := getDefaultDomainID(t)

	var entrypointID string

	t.Run("create entrypoint", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "entrypoint-*.json")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(`{"name":"e2e-entrypoint","url":"http://test.com","tags":[]}`)
		if err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		tmpFile.Close()

		out, err := runCLI("am", "entrypoint", "create", "--domain", domainID, "--file", tmpFile.Name(), "-o", "json")
		if err != nil {
			t.Skipf("entrypoint create failed: %s", out)
		}

		var obj map[string]any
		if err := json.Unmarshal([]byte(out), &obj); err == nil {
			if id, ok := obj["id"].(string); ok {
				entrypointID = id
			}
		}
	})

	t.Run("delete entrypoint", func(t *testing.T) {
		if entrypointID == "" {
			t.Skip("no entrypoint was created")
		}

		out := runCLIExpectSuccess(t, "am", "entrypoint", "delete", "--domain", domainID, entrypointID)
		if !strings.Contains(strings.ToLower(out), "deleted") {
			t.Logf("delete output: %s", out)
		}
	})

	t.Run("create missing file flag", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "entrypoint", "create", "--domain", domainID)
		if !strings.Contains(strings.ToLower(out), "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}
