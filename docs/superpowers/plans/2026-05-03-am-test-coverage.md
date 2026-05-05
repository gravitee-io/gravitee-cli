# AM CLI Test Coverage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add missing command-level tests for `audit`, `trace`, `supportdump`, `lint`, and `diff` packages so every non-trivial command is covered by `go test ./...` which runs on every CI push.

**Architecture:** All tests use the existing in-process `client.FakeClient` pattern (inject via `factory.Factory`) except `diff` which creates its own HTTP clients internally — that one uses `net/http/httptest.NewServer`. No new production code changes; only test files.

**Tech Stack:** Go `testing`, `net/http/httptest`, `encoding/json`, `strings`, existing `client.FakeClient` / `factory.Factory` test helpers.

---

## File Structure

| File | Action | Purpose |
|------|--------|---------|
| `cmd/am/audit/audit_test.go` | Modify | Add `TestParseEvent`, `TestFormatEvent`, `TestAuditListColumns`, `TestAuditListAll` |
| `cmd/am/trace/helpers_test.go` | **Create** | `newTestFactory` helper (package `trace`) |
| `cmd/am/trace/trace_test.go` | Modify | Add 5 missing check-function tests + `TestRunTrace` e2e |
| `cmd/am/supportdump/helpers_test.go` | **Create** | `newTestFactory` helper (package `supportdump`) |
| `cmd/am/supportdump/supportdump_test.go` | Modify | Add `TestSupportDumpSingleDomain` |
| `cmd/am/lint/lint_test.go` | Modify | Add `TestLintCmd` e2e |
| `cmd/am/diff/diff_test.go` | Modify | Add `TestDiffCmd` e2e using httptest |

---

### Task 1: `audit` — parseEvent, formatEvent, columns, --all

**Files:**
- Modify: `cmd/am/audit/audit_test.go`

- [ ] **Step 1: Add four tests to audit_test.go**

Append to `cmd/am/audit/audit_test.go` (inside the existing `package audit` file, after the last function):

```go
func TestParseEvent(t *testing.T) {
	raw := json.RawMessage(`{
		"id": "audit-1",
		"type": "USER_LOGIN",
		"outcome": {"status": "SUCCESS"},
		"actor": {"displayName": "john"},
		"target": {"displayName": "MyApp"},
		"timestamp": 1700000000000
	}`)
	e := parseEvent(raw)
	if e.ID != "audit-1" {
		t.Errorf("expected id 'audit-1', got %q", e.ID)
	}
	if e.EventType != "USER_LOGIN" {
		t.Errorf("expected type 'USER_LOGIN', got %q", e.EventType)
	}
	if e.Status != "SUCCESS" {
		t.Errorf("expected status 'SUCCESS', got %q", e.Status)
	}
	if e.Actor != "john" {
		t.Errorf("expected actor 'john', got %q", e.Actor)
	}
	if e.Target != "MyApp" {
		t.Errorf("expected target 'MyApp', got %q", e.Target)
	}
	if e.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestFormatEvent(t *testing.T) {
	e := auditEvent{
		Timestamp: "2024-01-01 00:00:00",
		Status:    "SUCCESS",
		EventType: "USER_LOGIN",
		Actor:     "john",
		Target:    "MyApp",
	}
	out := formatEvent(e)
	for _, want := range []string{"SUCCESS", "USER_LOGIN", "john", "→ MyApp"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in formatEvent output, got: %s", want, out)
		}
	}
}

func TestAuditListColumns(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id":      "a-1",
				"type":    "USER_LOGIN",
				"outcome": map[string]interface{}{"status": "SUCCESS"},
				"actor":   map[string]interface{}{"displayName": "john"},
				"target":  map[string]interface{}{"displayName": "MyApp"},
				"timestamp": float64(1700000000000),
			},
		},
		"totalCount": 1,
	}
	data, _ := json.Marshal(resp)
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) { return data, nil },
	}
	f, out := newTestFactory(fake, false)
	cmd := newListCmd(f)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"john", "MyApp", "2023"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("expected %q in output, got: %s", want, out.String())
		}
	}
}

func TestAuditListAll(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{"id": "a-1", "type": "LOGIN"},
			{"id": "a-2", "type": "LOGOUT"},
		},
		"totalCount": 2,
	}
	data, _ := json.Marshal(resp)
	calls := 0
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			calls++
			return data, nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := newListCmd(f)
	cmd.SetArgs([]string{"--all"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 page fetch for 2-item result, got %d", calls)
	}
	if !strings.Contains(out.String(), "LOGIN") {
		t.Errorf("expected 'LOGIN' in output, got: %s", out.String())
	}
}
```

Also add `"strings"` to the import block at the top of `audit_test.go` if not already present. The file already imports `"encoding/json"` and `"strings"` — check before adding.

- [ ] **Step 2: Run tests**

```
go test ./cmd/am/audit/... -v -run "TestParseEvent|TestFormatEvent|TestAuditListColumns|TestAuditListAll"
```

Expected: all 4 PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/am/audit/audit_test.go
git commit -m "test: add parseEvent, formatEvent, columns, and --all tests for audit list"
```

---

### Task 2: `trace` — helpers + 5 missing check functions + runTrace e2e

**Files:**
- Create: `cmd/am/trace/helpers_test.go`
- Modify: `cmd/am/trace/trace_test.go`

- [ ] **Step 1: Create cmd/am/trace/helpers_test.go**

```go
package trace

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(fc *client.FakeClient) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	return &factory.Factory{
		Config: &config.Config{
			CurrentContext: "am-test",
			Contexts: map[string]config.Context{
				"am-test": {
					URL: "https://am-test.com", Token: "tok",
					Org: "DEFAULT", Env: "DEFAULT",
					Type: "am", Domain: "test-domain",
				},
			},
		},
		Resolved: &config.ResolvedContext{
			Name: "am-test", URL: "https://am-test.com", Token: "tok",
			Org: "DEFAULT", Env: "DEFAULT",
			Type: "am", Domain: "test-domain",
		},
		Client:       fc,
		IOStreams:    factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
		OutputFormat: "table",
	}, out
}
```

- [ ] **Step 2: Append 5 check-function tests + 1 e2e to trace_test.go**

Add these at the bottom of `cmd/am/trace/trace_test.go`. The file already has `package trace` and `import "testing"` — add `"encoding/json"`, `"fmt"`, `"strings"` to the import block.

```go
func TestCheckGrantTypes_UserFacing(t *testing.T) {
	app := map[string]interface{}{
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"grantTypes": []interface{}{"authorization_code", "client_credentials"},
			},
		},
	}
	step := checkGrantTypes(app)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q: %s", step.Status, step.Detail)
	}
	if !strings.Contains(step.Detail, "authorization_code") {
		t.Errorf("expected 'authorization_code' in detail, got %q", step.Detail)
	}
}

func TestCheckGrantTypes_NoUserFacing(t *testing.T) {
	app := map[string]interface{}{
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"grantTypes": []interface{}{"client_credentials"},
			},
		},
	}
	step := checkGrantTypes(app)
	if step.Status != "warn" {
		t.Errorf("expected warn, got %q", step.Status)
	}
}

func TestCheckMfa_NoFactors(t *testing.T) {
	step := checkMfa(map[string]interface{}{}, nil)
	if step.Status != "ok" {
		t.Errorf("expected ok when no domain factors, got %q", step.Status)
	}
}

func TestCheckMfa_Enrolled(t *testing.T) {
	user := map[string]interface{}{
		"factors": []interface{}{"factor-1"},
	}
	domainFactors := []map[string]interface{}{{"id": "factor-1", "name": "TOTP"}}
	step := checkMfa(user, domainFactors)
	if step.Status != "ok" {
		t.Errorf("expected ok when factor enrolled, got %q: %s", step.Status, step.Detail)
	}
}

func TestCheckMfa_RequiredNotEnrolled(t *testing.T) {
	user := map[string]interface{}{}
	domainFactors := []map[string]interface{}{{"id": "factor-1", "name": "TOTP"}}
	step := checkMfa(user, domainFactors)
	if step.Status != "warn" {
		t.Errorf("expected warn when factors required but not enrolled, got %q", step.Status)
	}
	if !strings.Contains(step.Detail, "TOTP") {
		t.Errorf("expected factor name 'TOTP' in detail, got %q", step.Detail)
	}
}

func TestCheckFlows_WithPolicies(t *testing.T) {
	flows := []map[string]interface{}{
		{
			"type": "LOGIN",
			"pre": []interface{}{
				map[string]interface{}{"name": "IP-Filter"},
			},
		},
	}
	step := checkFlows(flows)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q", step.Status)
	}
	if !strings.Contains(step.Detail, "IP-Filter") {
		t.Errorf("expected policy name in detail, got %q", step.Detail)
	}
}

func TestCheckFlows_NoPolicies(t *testing.T) {
	step := checkFlows(nil)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q", step.Status)
	}
	if !strings.Contains(step.Detail, "No pre-login") {
		t.Errorf("expected 'No pre-login' in detail, got %q", step.Detail)
	}
}

func TestCheckConsent_Skipped(t *testing.T) {
	app := map[string]interface{}{
		"settings": map[string]interface{}{
			"advanced": map[string]interface{}{"skipConsent": true},
		},
	}
	step := checkConsent(app)
	if !strings.Contains(step.Detail, "skipped") {
		t.Errorf("expected 'skipped' in detail, got %q", step.Detail)
	}
}

func TestCheckTokenConfig(t *testing.T) {
	app := map[string]interface{}{
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"accessTokenValiditySeconds":  float64(3600),
				"refreshTokenValiditySeconds": float64(86400),
				"idTokenValiditySeconds":      nil,
			},
		},
	}
	step := checkTokenConfig(app)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q", step.Status)
	}
	if !strings.Contains(step.Detail, "3600") {
		t.Errorf("expected access token value in detail, got %q", step.Detail)
	}
	if !strings.Contains(step.Detail, "default") {
		t.Errorf("expected 'default' for nil id token, got %q", step.Detail)
	}
}

func TestRunTrace(t *testing.T) {
	user := map[string]interface{}{
		"id": "user-1", "username": "john", "email": "john@example.com",
		"enabled": true, "accountNonLocked": true,
		"source": "idp-1",
	}
	app := map[string]interface{}{
		"id": "app-1", "name": "MyApp",
		"identityProviders": []interface{}{
			map[string]interface{}{"identity": "idp-1"},
		},
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"grantTypes": []interface{}{"authorization_code"},
			},
			"advanced": map[string]interface{}{"skipConsent": false},
		},
	}
	userBytes, _ := json.Marshal(user)
	appBytes, _ := json.Marshal(app)
	empty, _ := json.Marshal([]interface{}{})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			switch {
			case strings.Contains(path, "/users/user-1"):
				return userBytes, nil
			case strings.Contains(path, "/applications/app-1"):
				return appBytes, nil
			case strings.Contains(path, "/identities"):
				return empty, nil
			case strings.Contains(path, "/factors"):
				return empty, nil
			case strings.Contains(path, "/flows"):
				return empty, nil
			}
			return nil, fmt.Errorf("unexpected path: %s", path)
		},
	}
	f, out := newTestFactory(fake)
	cmd := NewTraceCmd(f)
	cmd.SetArgs([]string{"--user", "user-1", "--app", "app-1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "john@example.com") {
		t.Errorf("expected user email in output, got: %s", output)
	}
	if !strings.Contains(output, "All checks passed") {
		t.Errorf("expected 'All checks passed' verdict, got: %s", output)
	}
}
```

- [ ] **Step 3: Run tests**

```
go test ./cmd/am/trace/... -v
```

Expected: all 14 tests PASS (7 existing + 7 new check tests + TestRunTrace = 15 total).

- [ ] **Step 4: Commit**

```bash
git add cmd/am/trace/helpers_test.go cmd/am/trace/trace_test.go
git commit -m "test: add missing check function tests and runTrace e2e for trace package"
```

---

### Task 3: `supportdump` — helpers + e2e command test

**Files:**
- Create: `cmd/am/supportdump/helpers_test.go`
- Modify: `cmd/am/supportdump/supportdump_test.go`

- [ ] **Step 1: Create cmd/am/supportdump/helpers_test.go**

```go
package supportdump

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(fc *client.FakeClient) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	return &factory.Factory{
		Config: &config.Config{
			CurrentContext: "am-test",
			Contexts: map[string]config.Context{
				"am-test": {
					URL: "https://am-test.com", Token: "tok",
					Org: "DEFAULT", Env: "DEFAULT",
					Type: "am", Domain: "test-domain",
				},
			},
		},
		Resolved: &config.ResolvedContext{
			Name: "am-test", URL: "https://am-test.com", Token: "tok",
			Org: "DEFAULT", Env: "DEFAULT",
			Type: "am", Domain: "test-domain",
		},
		Client:       fc,
		IOStreams:    factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
		OutputFormat: "table",
	}, out
}
```

- [ ] **Step 2: Append e2e test to supportdump_test.go**

Add at the bottom of `cmd/am/supportdump/supportdump_test.go`. The file's current import block only uses `"testing"` — add `"encoding/json"` and `"strings"`.

```go
func TestSupportDumpSingleDomain(t *testing.T) {
	domain := map[string]interface{}{"id": "test-domain", "name": "Test Domain", "enabled": true}
	emptyPaginated := map[string]interface{}{"data": []interface{}{}, "totalCount": 0}
	emptyArr := []interface{}{}

	domainBytes, _ := json.Marshal(domain)
	emptyPaginatedBytes, _ := json.Marshal(emptyPaginated)
	emptyArrBytes, _ := json.Marshal(emptyArr)

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			switch {
			case strings.HasSuffix(strings.Split(path, "?")[0], "/domains/test-domain"):
				return domainBytes, nil
			case strings.Contains(path, "/applications"),
				strings.Contains(path, "/roles"),
				strings.Contains(path, "/scopes"),
				strings.Contains(path, "/groups"):
				return emptyPaginatedBytes, nil
			case strings.Contains(path, "/audits"):
				return emptyPaginatedBytes, nil
			default:
				return emptyArrBytes, nil
			}
		},
	}
	f, out := newTestFactory(fake)
	cmd := NewSupportDumpCmd(f)
	cmd.SetArgs([]string{"--no-redact"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "_metadata") {
		t.Errorf("expected '_metadata' in JSON output, got: %s", output)
	}
	if !strings.Contains(output, "test-domain") {
		t.Errorf("expected domain ID in JSON output, got: %s", output)
	}
	if !strings.Contains(output, "Test Domain") {
		t.Errorf("expected domain name in JSON output, got: %s", output)
	}
}
```

- [ ] **Step 3: Run tests**

```
go test ./cmd/am/supportdump/... -v
```

Expected: all 4 PASS (3 existing + 1 new).

- [ ] **Step 4: Commit**

```bash
git add cmd/am/supportdump/helpers_test.go cmd/am/supportdump/supportdump_test.go
git commit -m "test: add e2e command test for support-dump single-domain mode"
```

---

### Task 4: `lint` — e2e command test

**Files:**
- Modify: `cmd/am/lint/lint_test.go`

- [ ] **Step 1: Append e2e test to lint_test.go**

The file already imports `"strings"` and `"testing"`. Add `"encoding/json"` to the import block, then append:

```go
func TestLintCmd(t *testing.T) {
	// App with implicit grant triggers a critical finding
	apps := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"name": "bad-app",
				"settings": map[string]interface{}{
					"oauth": map[string]interface{}{
						"grantTypes": []interface{}{"implicit"},
					},
				},
			},
		},
		"totalCount": 1,
	}
	appsBytes, _ := json.Marshal(apps)
	emptyArr, _ := json.Marshal([]interface{}{})
	emptyList, _ := json.Marshal(map[string]interface{}{"data": []interface{}{}, "totalCount": 0})

	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			switch {
			case strings.Contains(path, "applications"):
				return appsBytes, nil
			case strings.Contains(path, "scopes"):
				return emptyList, nil
			default:
				return emptyArr, nil
			}
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := NewLintCmd(f)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "critical") {
		t.Errorf("expected 'critical' finding in output, got: %s", output)
	}
	if !strings.Contains(output, "Score:") {
		t.Errorf("expected 'Score:' in output, got: %s", output)
	}
}
```

Note: `newTestFactory` already exists in `cmd/am/lint/helpers_test.go` — do not redefine it.
Also: `client` package is already imported in `helpers_test.go` within the same `package lint` — `lint_test.go` can use `&client.FakeClient{}` because they share the package. Add `"github.com/gravitee-io/gio-cli/internal/client"` to `lint_test.go`'s import block.

- [ ] **Step 2: Run tests**

```
go test ./cmd/am/lint/... -v -run TestLintCmd
```

Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add cmd/am/lint/lint_test.go
git commit -m "test: add e2e command test for lint (critical finding + score output)"
```

---

### Task 5: `diff` — e2e command test with httptest

**Files:**
- Modify: `cmd/am/diff/diff_test.go`

The `diff` command creates its own HTTP clients via `client.NewHTTPClient` using the resolved context URLs. The existing `newTestFactory` in `helpers_test.go` sets up `ctx-a` (http://am-a) and `ctx-b` (http://am-b). The test overrides these URLs to point at real `httptest.NewServer` instances.

- [ ] **Step 1: Append e2e test to diff_test.go**

The file already imports `"strings"` and `"testing"`. Add `"encoding/json"`, `"net/http"`, `"net/http/httptest"` to the import block. The `config` package is used in the test to update context entries — add `"github.com/gravitee-io/gio-cli/internal/config"` if not already imported.

```go
func TestDiffCmd(t *testing.T) {
	// from: has scope-a only
	fromHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "scopes") {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{{"key": "scope-a", "name": "Scope A"}},
			})
			return
		}
		if strings.Contains(r.URL.Path, "roles") ||
			strings.Contains(r.URL.Path, "groups") ||
			strings.Contains(r.URL.Path, "applications") {
			json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
			return
		}
		json.NewEncoder(w).Encode([]interface{}{})
	})
	// to: has scope-a + scope-b (one extra)
	toHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "scopes") {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"key": "scope-a", "name": "Scope A"},
					{"key": "scope-b", "name": "Scope B"},
				},
			})
			return
		}
		if strings.Contains(r.URL.Path, "roles") ||
			strings.Contains(r.URL.Path, "groups") ||
			strings.Contains(r.URL.Path, "applications") {
			json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{}})
			return
		}
		json.NewEncoder(w).Encode([]interface{}{})
	})

	fromServer := httptest.NewServer(fromHandler)
	defer fromServer.Close()
	toServer := httptest.NewServer(toHandler)
	defer toServer.Close()

	f, out := newTestFactory(nil)
	// Point ctx-a and ctx-b at the httptest servers
	f.Config.Contexts["ctx-a"] = config.Context{
		URL: fromServer.URL, Token: "tok-a",
		Org: "DEFAULT", Env: "DEFAULT", Domain: "dom-a", Type: "am",
	}
	f.Config.Contexts["ctx-b"] = config.Context{
		URL: toServer.URL, Token: "tok-b",
		Org: "DEFAULT", Env: "DEFAULT", Domain: "dom-b", Type: "am",
	}

	cmd := NewDiffCmd(f)
	cmd.SetArgs([]string{"--from", "ctx-a", "--to", "ctx-b"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "scopes") {
		t.Errorf("expected 'scopes' in output, got: %s", output)
	}
	if !strings.Contains(output, "+1") {
		t.Errorf("expected '+1' added scope in diff output, got: %s", output)
	}
}
```

Note: `newTestFactory` in `helpers_test.go` takes `c client.GraviteeClient`. Passing `nil` is valid — `diff`'s RunE never calls `f.Client`, so the nil interface is never dereferenced.

- [ ] **Step 2: Run tests**

```
go test ./cmd/am/diff/... -v -run TestDiffCmd
```

Expected: PASS. The httptest servers respond to the exact paths `fetchItems` constructs (`/management/organizations/DEFAULT/environments/DEFAULT/domains/dom-a/scopes?...`).

- [ ] **Step 3: Run full suite**

```
go test ./...
```

Expected: all tests PASS. This is the same command run by `make test-cover` in CI.

- [ ] **Step 4: Commit**

```bash
git add cmd/am/diff/diff_test.go
git commit -m "test: add e2e command test for diff using httptest servers"
```

---

## Self-Review

**Spec coverage:**
- ✓ `audit`: `parseEvent`, `formatEvent`, column assertions, `--all`
- ✓ `trace`: all 7 check functions covered, `runTrace` e2e covered
- ✓ `supportdump`: e2e single-domain execution covered
- ✓ `lint`: e2e command execution with critical finding + score output
- ✓ `diff`: e2e with httptest, scopes diff detected

**Placeholder scan:** No TBD, no "add appropriate…", all code blocks are complete.

**Type consistency:**
- `auditEvent` struct fields match `parseEvent` return and `formatEvent` parameter throughout Task 1
- `TraceStep.Status` string values ("ok", "warn", "block") match existing `checkUserStatus` pattern
- `client.FakeClient` field `GetFunc func(path string) ([]byte, error)` matches everywhere
- `newTestFactory` signatures are consistent within each package

**CI impact:** `go test ./...` (run by `make test-cover` in `.github/workflows/ci.yml`) automatically picks up all new test files. No CI config changes needed.
