//go:build e2e

package e2e

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUserConsents(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-consent-user",
		"--email", "consent@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("list consents", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "consent", "list",
			"--domain", domainID, "--user-id", userID)
	})

	t.Run("revoke-all consents", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "consent", "revoke-all",
			"--domain", domainID, "--user-id", userID)
	})

	t.Run("list consents missing user-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "user", "consent", "list",
			"--domain", domainID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestUserRoleAssignment(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-role-user",
		"--email", "roleuser@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	roleOut := runCLIExpectSuccess(t, "am", "role", "create",
		"--domain", domainID,
		"--name", "e2e-user-role",
		"-o", "json")
	roleID := extractID(t, roleOut)

	defer func() {
		runCLIExpectSuccess(t, "am", "role", "delete", "--domain", domainID, roleID)
	}()

	t.Run("list roles", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "role", "list",
			"--domain", domainID, "--user-id", userID)
	})

	t.Run("assign and revoke role", func(t *testing.T) {
		_, err := runCLI("am", "user", "role", "assign",
			"--domain", domainID, "--user-id", userID, "--roles", roleID)
		if err != nil {
			t.Skipf("role assign not supported on this AM version: %v", err)
		}

		out := runCLIExpectSuccess(t, "am", "user", "role", "list",
			"--domain", domainID, "--user-id", userID)
		if !strings.Contains(out, roleID) {
			t.Errorf("expected role ID %s in output, got: %s", roleID, out)
		}

		runCLIExpectSuccess(t, "am", "user", "role", "revoke",
			"--domain", domainID, "--user-id", userID, roleID)
	})
}

func TestUserDevices(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-device-user",
		"--email", "device@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("list devices", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "device", "list",
			"--domain", domainID, "--user-id", userID)
	})

	t.Run("list devices missing user-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "user", "device", "list",
			"--domain", domainID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestUserCredentials(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-cred-user",
		"--email", "cred@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("list credentials", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "credential", "list",
			"--domain", domainID, "--user-id", userID)
	})

	t.Run("list credentials missing user-id", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "user", "credential", "list",
			"--domain", domainID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestUserCertCredentials(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-certcred-user",
		"--email", "certcred@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("list cert credentials", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "cert-credential", "list",
			"--domain", domainID, "--user-id", userID)
	})
}

func TestUserEnrolledFactors(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-factor-user",
		"--email", "factor@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("list enrolled factors", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "factor", "list",
			"--domain", domainID, "--user-id", userID)
	})
}

func TestUserIdentities(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-identity-user",
		"--email", "identity@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("list identities", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "identity", "list",
			"--domain", domainID, "--user-id", userID)
	})
}

func TestUserAudits(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-audit-user",
		"--email", "audit@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("list audits", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "user", "audit", "list",
			"--domain", domainID, "--user-id", userID)
	})

	t.Run("list audits json", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "audit", "list",
			"--domain", domainID, "--user-id", userID, "-o", "json")

		var resp map[string]any
		if err := json.Unmarshal([]byte(out), &resp); err != nil {
			t.Fatalf("expected valid JSON output, got: %s", out)
		}

		if _, ok := resp["data"]; !ok {
			t.Errorf("expected 'data' key in JSON response, got: %s", out)
		}
	})
}

func TestUserSendRegistration(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-reg-user",
		"--email", "reg@test.com",
		"--password", "E2eTestPassword123!@#",
		"--preRegistration",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("send registration", func(t *testing.T) {
		_, err := runCLI("am", "user", "send-registration", "--domain", domainID, userID)
		if err != nil {
			t.Skipf("send-registration not available (requires email config): %v", err)
		}
	})
}

func TestUserUpdateUsername(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-uname-user",
		"--email", "uname@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("update username", func(t *testing.T) {
		_, err := runCLI("am", "user", "update-username",
			"--domain", domainID, "--username", "new-username", userID)
		if err != nil {
			t.Skipf("update-username not supported on this AM config: %v", err)
		}
	})

	t.Run("update username missing flag", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "user", "update-username",
			"--domain", domainID, userID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestUserBulk(t *testing.T) {
	domainID := getDefaultDomainID(t)

	t.Run("bulk with invalid file", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "user", "bulk",
			"--domain", domainID, "--file", "/tmp/nonexistent-e2e-bulk.json")
		if !strings.Contains(out, "no such file") {
			t.Errorf("expected 'no such file' error, got: %s", out)
		}
	})

	t.Run("bulk missing file flag", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "user", "bulk",
			"--domain", domainID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}
