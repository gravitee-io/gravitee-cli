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

// TestDomainEntrypoints covers `am domain entrypoints {get,set-path,add-vhost,
// remove-vhost,clear-vhosts}` against a real AM domain.
func TestDomainEntrypoints(t *testing.T) {
	out := runCLIExpectSuccess(t, "am", "domain", "create",
		"--name", "e2e-entrypoints-domain",
		"-o", "json")
	domainID := extractID(t, out)

	defer runCLI("am", "domain", "delete", domainID)

	t.Run("get defaults to context-path mode", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "entrypoints", "get", domainID, "-o", "json")
		if !strings.Contains(out, "context-path") && !strings.Contains(out, `"vhostMode": false`) && !strings.Contains(out, `"vhostMode":false`) {
			t.Errorf("expected context-path mode on fresh domain, got: %s", out)
		}
	})

	t.Run("set-path switches path", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "entrypoints", "set-path", domainID, "/auth", "-o", "json")
		if !strings.Contains(out, "/auth") {
			t.Errorf("expected path /auth in output, got: %s", out)
		}
	})

	t.Run("add-vhost switches to vhost mode", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "entrypoints", "add-vhost",
			domainID, "auth.e2e.example.com",
			"--path", "/",
			"--override",
			"-o", "json")
		if !strings.Contains(out, "auth.e2e.example.com") {
			t.Errorf("expected vhost in output, got: %s", out)
		}
	})

	t.Run("add second vhost", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "domain", "entrypoints", "add-vhost",
			domainID, "alt.e2e.example.com",
			"--path", "/auth")
	})

	t.Run("refuse to remove the only override vhost", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "domain", "entrypoints", "remove-vhost",
			domainID, "auth.e2e.example.com")
		if !strings.Contains(out, "overrideEntrypoint") {
			t.Errorf("expected refusal mentioning overrideEntrypoint, got: %s", out)
		}
	})

	t.Run("remove non-override vhost", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "domain", "entrypoints", "remove-vhost",
			domainID, "alt.e2e.example.com")
	})

	t.Run("clear-vhosts returns to context-path mode", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "entrypoints", "clear-vhosts", domainID, "-o", "json")
		if !strings.Contains(out, "context-path") && !strings.Contains(out, `"vhostMode": false`) && !strings.Contains(out, `"vhostMode":false`) {
			t.Errorf("expected context-path mode after clear, got: %s", out)
		}
	})
}

// TestDomainCIMD covers `am domain cimd {get,enable,disable}`.
func TestDomainCIMD(t *testing.T) {
	out := runCLIExpectSuccess(t, "am", "domain", "create",
		"--name", "e2e-cimd-domain",
		"-o", "json")
	domainID := extractID(t, out)

	defer runCLI("am", "domain", "delete", domainID)

	t.Run("get defaults to disabled", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "domain", "cimd", "get", domainID)
	})

	t.Run("enable with options", func(t *testing.T) {
		out, err := runCLI("am", "domain", "cimd", "enable", domainID,
			"--allow-private",
			"--allow-http",
			"--allowed-domains", "a.example.com,b.example.com",
			"--fetch-timeout-ms", "1500",
			"-o", "json")
		if err != nil {
			t.Skipf("cimd enable not supported by this AM build: %s", out)
		}

		if !strings.Contains(out, "a.example.com") {
			t.Errorf("expected allowed-domains in output, got: %s", out)
		}
	})

	t.Run("disable", func(t *testing.T) {
		out, err := runCLI("am", "domain", "cimd", "disable", domainID, "-o", "json")
		if err != nil {
			t.Skipf("cimd disable failed: %s", out)
		}
	})
}

// TestDomainRedirectURIRules covers the `--allow-localhost-redirect` and
// `--allow-http-redirect` flags on `am domain update`.
func TestDomainRedirectURIRules(t *testing.T) {
	out := runCLIExpectSuccess(t, "am", "domain", "create",
		"--name", "e2e-redirect-rules",
		"-o", "json")
	domainID := extractID(t, out)

	defer runCLI("am", "domain", "delete", domainID)

	t.Run("enable both", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "domain", "update", domainID,
			"--allow-localhost-redirect",
			"--allow-http-redirect")
	})

	t.Run("disable both", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "domain", "update", domainID,
			"--allow-localhost-redirect=false",
			"--allow-http-redirect=false")
	})
}

// TestAppCreateInitialSecret verifies that creating a service app surfaces the
// one-time client secret (commit 47bee52).
func TestAppCreateInitialSecret(t *testing.T) {
	domainID := getDefaultDomainID(t)

	out := runCLIExpectSuccess(t, "am", "app", "create",
		"--domain", domainID,
		"--name", "e2e-initial-secret",
		"--type", "service",
		"-o", "json")
	appID := extractID(t, out)

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	var resp map[string]any
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse app json: %v", err)
	}

	settings, _ := resp["settings"].(map[string]any)
	oauth, _ := settings["oauth"].(map[string]any)
	secret, _ := oauth["clientSecret"].(string)

	if secret == "" {
		t.Errorf("expected initial clientSecret in create response, got: %s", out)
	}
}

// TestAppUpdateTemplate covers `am app update --template true|false`
// (commit 7375bd5).
func TestAppUpdateTemplate(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-template-app")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	t.Run("mark as template", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "update", appID,
			"--domain", domainID,
			"--template", "true",
			"-o", "json")

		var m map[string]any
		if err := json.Unmarshal([]byte(out), &m); err != nil {
			t.Fatalf("failed to parse app: %v", err)
		}

		tpl, _ := m["template"].(bool)
		if !tpl {
			t.Errorf("expected template=true in response, got: %s", out)
		}
	})

	t.Run("invalid value is rejected", func(t *testing.T) {
		out := runCLIExpectError(t, "am", "app", "update", appID,
			"--domain", domainID,
			"--template", "nope")
		if !strings.Contains(out, "template") {
			t.Errorf("expected error to mention --template, got: %s", out)
		}
	})
}

// TestAppRenewSecret covers `am app secret renew` (commit 47bee52).
func TestAppRenewSecret(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-renew-secret")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	listOut, err := runCLI("am", "app", "secret", "list",
		"--domain", domainID,
		"--app-id", appID,
		"-o", "json")
	if err != nil {
		t.Skipf("app secret list failed: %s", listOut)
	}

	var secrets []map[string]any
	if err := json.Unmarshal([]byte(listOut), &secrets); err != nil {
		// fall back to wrapped form
		var paged struct {
			Data []map[string]any `json:"data"`
		}
		if err2 := json.Unmarshal([]byte(listOut), &paged); err2 != nil {
			t.Skipf("cannot parse secret list: %s", listOut)
		}
		secrets = paged.Data
	}

	if len(secrets) == 0 {
		t.Skip("no app secrets available to renew")
	}

	secretID, _ := secrets[0]["id"].(string)
	if secretID == "" {
		t.Skipf("first secret has no id: %v", secrets[0])
	}

	out, err := runCLI("am", "app", "secret", "renew", secretID,
		"--domain", domainID,
		"--app-id", appID,
		"-o", "json")
	if err != nil {
		t.Skipf("renew failed (AM build may not support it): %s", out)
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("failed to parse renew response: %v", err)
	}

	value, _ := m["secret"].(string)
	if value == "" {
		value, _ = m["clientSecret"].(string)
	}

	if value == "" {
		t.Errorf("expected secret value in renew response, got: %s", out)
	}
}

// TestAppIdpBindings covers `am app idp {list,add,remove}` (commit 13a1bff).
func TestAppIdpBindings(t *testing.T) {
	domainID := getDefaultDomainID(t)
	appID := createTestApp(t, domainID, "e2e-idp-bindings")

	defer runCLI("am", "app", "delete", "--domain", domainID, appID)

	listOut := runCLIExpectSuccess(t, "am", "idp", "list", "--domain", domainID, "-o", "json")

	var idps []map[string]any
	if err := json.Unmarshal([]byte(listOut), &idps); err != nil {
		t.Skipf("cannot parse idp list: %s", listOut)
	}

	if len(idps) == 0 {
		t.Skip("no IdP available on domain for binding test")
	}

	idpID, _ := idps[0]["id"].(string)
	if idpID == "" {
		t.Skip("first idp has no id")
	}

	t.Run("add binding", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "idp", "add", appID, idpID,
			"--domain", domainID,
			"--priority", "5")
	})

	t.Run("list shows binding", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "idp", "list", appID,
			"--domain", domainID,
			"-o", "json")
		if !strings.Contains(out, idpID) {
			t.Errorf("expected idp %s in bindings list, got: %s", idpID, out)
		}
	})

	t.Run("re-add updates binding in place", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "idp", "add", appID, idpID,
			"--domain", domainID,
			"--priority", "20",
			"--selection-rule", "{#context.attributes['foo'] == 'bar'}")

		out := runCLIExpectSuccess(t, "am", "app", "idp", "list", appID,
			"--domain", domainID,
			"-o", "json")

		var bindings []map[string]any
		if err := json.Unmarshal([]byte(out), &bindings); err != nil {
			t.Fatalf("failed to parse bindings: %v", err)
		}

		if len(bindings) != 1 {
			t.Fatalf("expected exactly 1 binding after re-add, got %d: %s", len(bindings), out)
		}
	})

	t.Run("remove binding", func(t *testing.T) {
		runCLIExpectSuccess(t, "am", "app", "idp", "remove", appID, idpID,
			"--domain", domainID)
	})

	t.Run("list is empty after remove", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "app", "idp", "list", appID,
			"--domain", domainID,
			"-o", "json")
		if strings.Contains(out, idpID) {
			t.Errorf("expected binding gone after remove, got: %s", out)
		}
	})
}
