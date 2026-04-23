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
	"strings"
	"testing"
)

func TestUserUpdate(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-update-user",
		"--email", "update@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("update email", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "update", "--domain", domainID, userID, "--email", "updated@test.com")
		if !strings.Contains(out, "updated@test.com") {
			t.Errorf("expected updated email in output, got: %s", out)
		}
	})

	t.Run("update firstName", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "update", "--domain", domainID, userID, "--firstName", "John")
		if !strings.Contains(out, "John") {
			t.Errorf("expected firstName in output, got: %s", out)
		}
	})

	t.Run("update lastName", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "update", "--domain", domainID, userID, "--lastName", "Doe")
		if !strings.Contains(out, "Doe") {
			t.Errorf("expected lastName in output, got: %s", out)
		}
	})

	t.Run("update enabled false", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "update", "--domain", domainID, userID, "--enabled", "false")
		if !strings.Contains(out, "false") {
			t.Errorf("expected 'false' in output, got: %s", out)
		}
	})

	t.Run("update with no flags", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "user", "update", "--domain", domainID, userID)
		if !strings.Contains(out, "at least one flag") {
			t.Errorf("expected 'at least one flag' error, got: %s", out)
		}
	})

	t.Run("update with invalid enabled", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "user", "update", "--domain", domainID, userID, "--enabled", "maybe")
		if !strings.Contains(out, "invalid value") {
			t.Errorf("expected 'invalid value' error, got: %s", out)
		}
	})
}

func TestUserResetPassword(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "user", "create",
		"--domain", domainID,
		"--username", "e2e-reset-user",
		"--email", "reset@test.com",
		"--password", "E2eTestPassword123!@#",
		"-o", "json")
	userID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "user", "delete", "--domain", domainID, userID)
	}()

	t.Run("reset password", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "user", "reset-password",
			"--domain", domainID, "--password", "NewPassword456!@#", userID)
		if !strings.Contains(out, "reset") {
			t.Errorf("expected 'reset' in output, got: %s", out)
		}
	})

	t.Run("reset password missing password flag", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "user", "reset-password", "--domain", domainID, userID)
		if !strings.Contains(out, "required") {
			t.Errorf("expected 'required' error, got: %s", out)
		}
	})
}

func TestAppUpdate(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "app", "create",
		"--domain", domainID,
		"--name", "e2e-update-app",
		"--type", "service",
		"-o", "json")
	appID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "app", "delete", "--domain", domainID, appID)
	}()

	t.Run("update name", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "update", "--domain", domainID, appID, "--name", "e2e-updated-app")
		if !strings.Contains(out, "e2e-updated-app") {
			t.Errorf("expected updated name in output, got: %s", out)
		}
	})

	t.Run("update description", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "update", "--domain", domainID, appID, "--description", "Updated desc")
		if !strings.Contains(out, "Updated desc") {
			t.Errorf("expected updated description in output, got: %s", out)
		}
	})

	t.Run("update enabled false", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "update", "--domain", domainID, appID, "--enabled", "false")
	})

	t.Run("update enabled true", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "update", "--domain", domainID, appID, "--enabled", "true")
	})

	t.Run("update with no flags", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "update", "--domain", domainID, appID)
		if !strings.Contains(out, "at least one flag") {
			t.Errorf("expected 'at least one flag' error, got: %s", out)
		}
	})
}

func TestDomainListQuery(t *testing.T) {
	out := runCLIExpectSuccess(t, "am", "domain", "create",
		"--name", "e2e-query-test-domain",
		"-o", "json")
	domainID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "domain", "delete", domainID)
	}()

	t.Run("list with query flag accepted", func(t *testing.T) {
		// AM API may not support query search - just verify the command succeeds.
		runCLIExpectSuccess(t, "am", "domain", "list", "--query", "e2e-query-test")
	})

	t.Run("list with non-matching query", func(t *testing.T) {
		// With unsupported query, AM may return all or none - just verify no error.
		runCLIExpectSuccess(t, "am", "domain", "list", "--query", "nonexistent-xyz-12345")
	})
}

func TestScopeUpdateDescription(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "scope", "create",
		"--domain", domainID,
		"--key", "e2e_desc_scope",
		"--name", "Desc Test",
		"--description", "original",
		"-o", "json")
	scopeID := extractID(t, out)

	defer func() {
		runCLIExpectSuccess(t, "am", "scope", "delete", "--domain", domainID, scopeID)
	}()

	t.Run("update name", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "scope", "update", "--domain", domainID, scopeID, "--name", "Updated Name")
		if !strings.Contains(out, "Updated Name") {
			t.Errorf("expected updated name in output, got: %s", out)
		}
	})

	t.Run("update description", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "scope", "update", "--domain", domainID, scopeID, "--description", "updated desc")
		if !strings.Contains(out, "updated desc") {
			t.Errorf("expected updated description in output, got: %s", out)
		}
	})

	t.Run("update with no flags", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "scope", "update", "--domain", domainID, scopeID)
		if !strings.Contains(out, "at least one flag") {
			t.Errorf("expected 'at least one flag' error, got: %s", out)
		}
	})
}
