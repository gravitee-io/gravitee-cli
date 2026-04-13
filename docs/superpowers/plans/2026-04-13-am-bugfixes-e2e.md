# AM CLI Bugfixes + E2E Tests Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix all 11 bugs found during manual testing of the AM CLI and add E2E tests running against a real Gravitee AM instance in GitHub Actions.

**Architecture:** Bug fixes are isolated per-file changes to command handlers and service layer. E2E tests use `//go:build e2e` tag and run in a separate GH Actions workflow with Docker Compose (MongoDB + AM Management API). The E2E test binary is the compiled CLI itself, invoked via `exec.Command`.

**Tech Stack:** Go 1.26, Cobra CLI, Docker Compose, GitHub Actions, Gravitee AM 4.x

---

### Task 1: Add pagination validation helper to cmdutil

**Files:**
- Modify: `internal/cmdutil/cmdutil.go`
- Modify: `internal/cmdutil/cmdutil_test.go`

Fixes bugs #3 and #4: `--page 0` and `--per-page 0` cause server errors.

- [ ] **Step 1: Write failing tests for ValidatePagination**

Add to `internal/cmdutil/cmdutil_test.go`:

```go
func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name    string
		page    int
		perPage int
		wantErr string
	}{
		{"valid", 1, 10, ""},
		{"page zero", 0, 10, "--page must be >= 1"},
		{"page negative", -1, 10, "--page must be >= 1"},
		{"per-page zero", 1, 0, "--per-page must be >= 1"},
		{"per-page negative", 1, -5, "--per-page must be >= 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePagination(tt.page, tt.perPage)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/cmdutil/ -run TestValidatePagination -v`
Expected: FAIL — `ValidatePagination` undefined.

- [ ] **Step 3: Implement ValidatePagination**

Add to `internal/cmdutil/cmdutil.go`:

```go
// ValidatePagination checks that page and perPage are positive.
func ValidatePagination(page, perPage int) error {
	if page < 1 {
		return fmt.Errorf("--page must be >= 1, got %d", page)
	}
	if perPage < 1 {
		return fmt.Errorf("--per-page must be >= 1, got %d", perPage)
	}
	return nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/cmdutil/ -run TestValidatePagination -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/cmdutil/cmdutil.go internal/cmdutil/cmdutil_test.go
git commit -m "feat: add ValidatePagination helper to cmdutil"
```

---

### Task 2: Add pagination validation to all paginated list commands

**Files:**
- Modify: `cmd/am/domain/list.go`
- Modify: `cmd/am/app/list.go`
- Modify: `cmd/am/user/list.go`
- Modify: `cmd/am/role/list.go`
- Modify: `cmd/am/scope/list.go`
- Modify: `cmd/am/group/list.go`
- Modify: `cmd/am/audit/list.go`

Add `cmdutil.ValidatePagination(opts.page, opts.perPage)` call in each list command's `RunE`, before `opts.run()`. All follow the same pattern.

- [ ] **Step 1: Add validation to domain list**

In `cmd/am/domain/list.go`, in the `RunE` function, after `RequireContext` and before `return opts.run()`:

```go
if err := cmdutil.ValidatePagination(opts.page, opts.perPage); err != nil {
    return err
}
```

- [ ] **Step 2: Add validation to all other paginated list commands**

Apply the same pattern to `cmd/am/app/list.go`, `cmd/am/user/list.go`, `cmd/am/role/list.go`, `cmd/am/scope/list.go`, `cmd/am/group/list.go`, `cmd/am/audit/list.go`. Each has the same `RunE` structure — add the validation between `RequireContext` and `opts.run()`.

- [ ] **Step 3: Write a unit test for domain list with invalid page**

Add to `cmd/am/domain/list_test.go`:

```go
t.Run("rejects page zero", func(t *testing.T) {
    tc := testutil.NewFactory(&testutil.NoOpClient, false)
    err := testutil.Execute(newListCmd(tc.Factory), "--page", "0")
    testutil.AssertErrorContains(t, err, "--page must be >= 1")
})

t.Run("rejects per-page zero", func(t *testing.T) {
    tc := testutil.NewFactory(&testutil.NoOpClient, false)
    err := testutil.Execute(newListCmd(tc.Factory), "--per-page", "0")
    testutil.AssertErrorContains(t, err, "--per-page must be >= 1")
})
```

- [ ] **Step 4: Run all tests**

Run: `go test ./cmd/am/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/am/domain/list.go cmd/am/domain/list_test.go cmd/am/app/list.go cmd/am/user/list.go cmd/am/role/list.go cmd/am/scope/list.go cmd/am/group/list.go cmd/am/audit/list.go
git commit -m "fix: validate --page and --per-page are >= 1 in all list commands"
```

---

### Task 3: Fix domain create — add dataPlaneId

**Files:**
- Modify: `cmd/am/domain/create.go`
- Modify: `cmd/am/domain/create_test.go`

Fixes bug #1: `domain create` fails with 400 because `dataPlaneId` is not sent.

- [ ] **Step 1: Write failing test for dataPlaneId in body**

Add to `cmd/am/domain/create_test.go`:

```go
t.Run("includes dataPlaneId in request body", func(t *testing.T) {
    var sentBody map[string]any
    fake := &client.FakeClient{
        PostFunc: func(_ string, body any) ([]byte, error) {
            raw, _ := body.(json.RawMessage)
            json.Unmarshal(raw, &sentBody)
            return json.Marshal(map[string]any{"id": "new", "name": "Test"})
        },
    }
    tc := testutil.NewFactory(fake, false)

    err := testutil.Execute(newCreateCmd(tc.Factory), "--name", "Test")

    testutil.AssertNoError(t, err)
    if sentBody["dataPlaneId"] != "default" {
        t.Errorf("expected dataPlaneId 'default', got %v", sentBody["dataPlaneId"])
    }
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/am/domain/ -run TestCreateDomain/includes_dataPlaneId -v`
Expected: FAIL

- [ ] **Step 3: Update create.go to include dataPlaneId**

In `cmd/am/domain/create.go`:

Add a new flag to `createOptions`:

```go
type createOptions struct {
	factory     *factory.Factory
	name        string
	description string
	dataPlaneID string
}
```

Register the flag in `newCreateCmd`:

```go
cmd.Flags().StringVar(&opts.dataPlaneID, "data-plane-id", "default", "Data plane ID")
```

Update `run()` to include it in the body:

```go
body := map[string]any{
    "name":        o.name,
    "dataPlaneId": o.dataPlaneID,
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./cmd/am/domain/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/am/domain/create.go cmd/am/domain/create_test.go
git commit -m "fix: include dataPlaneId in domain create request body"
```

---

### Task 4: Add app type validation

**Files:**
- Modify: `cmd/am/app/create.go`
- Modify: `cmd/am/app/create_test.go`

Fixes bug #6: invalid `--type` is not validated by CLI.

- [ ] **Step 1: Write failing test**

Add to `cmd/am/app/create_test.go`:

```go
t.Run("rejects invalid app type", func(t *testing.T) {
    tc := testutil.NewFactory(&testutil.NoOpClient, false)
    err := testutil.Execute(newCreateCmd(tc.Factory, &domainID), "--name", "x", "--type", "invalid")
    testutil.AssertErrorContains(t, err, "invalid value 'invalid' for flag --type")
})
```

Note: `domainID` should be defined as `domainID := "test-domain"` at the top of the test function or use the pattern already established in the test file.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/am/app/ -run TestCreateApplication/rejects_invalid -v`
Expected: FAIL

- [ ] **Step 3: Add validation to create.go RunE**

In `cmd/am/app/create.go`, in `RunE`, after `RequireContext` and before `opts.run()`:

```go
if err := cmdutil.ValidateEnum(opts.appType, "type", []string{"web", "native", "browser", "service", "resource_server"}); err != nil {
    return err
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./cmd/am/app/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/am/app/create.go cmd/am/app/create_test.go
git commit -m "fix: validate --type flag in app create command"
```

---

### Task 5: Add --query to group list

**Files:**
- Modify: `internal/am/group.go`
- Modify: `internal/am/mock.go`
- Modify: `cmd/am/group/list.go`
- Modify: `cmd/am/group/list_test.go`

Fixes bug #7: `group list` lacks `--query` unlike other paginated list commands.

- [ ] **Step 1: Add Query field to ListGroupsParams**

In `internal/am/group.go`, update the struct:

```go
type ListGroupsParams struct {
	Query   string
	Page    int
	PerPage int
}
```

And update `ListGroups` to include query:

```go
func (s *service) ListGroups(domainID string, params ListGroupsParams) (*PaginatedResponse, error) {
	q := client.BuildQuery(map[string]string{
		"page": client.Itoa(params.Page),
		"size": client.Itoa(params.PerPage),
		"q":    params.Query,
	})
	// ... rest unchanged
}
```

- [ ] **Step 2: Add --query flag to group list command**

In `cmd/am/group/list.go`, add `query` to `listOptions`:

```go
type listOptions struct {
	factory  *factory.Factory
	domainID *string
	query    string
	page     int
	perPage  int
	all      bool
}
```

Register the flag:

```go
cmd.Flags().StringVar(&opts.query, "query", "", "Search by name")
```

Update `params()`:

```go
func (o *listOptions) params(page int) am.ListGroupsParams {
	return am.ListGroupsParams{
		Query:   o.query,
		Page:    page,
		PerPage: o.perPage,
	}
}
```

- [ ] **Step 3: Update mock**

In `internal/am/mock.go`, the `ListGroupsFunc` signature already matches `func(string, ListGroupsParams)` — no change needed since `ListGroupsParams` is the same type. The mock will pick up the new field automatically.

- [ ] **Step 4: Run tests**

Run: `go test ./cmd/am/group/ -v && go test ./internal/am/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/am/group.go cmd/am/group/list.go
git commit -m "feat: add --query flag to group list command"
```

---

### Task 6: Fix member list — show roleId in table

**Files:**
- Modify: `cmd/am/member/list.go`

Fixes bug #8: Role column is empty because API returns `roleId` not `role`.

- [ ] **Step 1: Update memberColumns**

In `cmd/am/member/list.go`, change `memberColumns()`:

```go
func memberColumns() []printer.Column {
	return []printer.Column{
		{Name: "MemberID", Value: func(i any) string { return cmdutil.StringField(i, "memberId") }},
		{Name: "RoleID", Value: func(i any) string { return cmdutil.StringField(i, "roleId") }},
		{Name: "MemberType", Value: func(i any) string { return cmdutil.StringField(i, "memberType") }},
	}
}
```

- [ ] **Step 2: Run tests**

Run: `go test ./cmd/am/member/ -v`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add cmd/am/member/list.go
git commit -m "fix: show roleId instead of empty role column in member list"
```

---

### Task 7: Remove certificate-settings command

**Files:**
- Delete: `cmd/am/certificate-settings/` (entire directory)
- Modify: `cmd/am/am.go`
- Modify: `internal/am/certificate_settings.go`
- Modify: `internal/am/service.go`
- Modify: `internal/am/mock.go`

Fixes bug #9: endpoint returns 405 — it doesn't exist in AM 4.x.

- [ ] **Step 1: Remove the import and AddCommand from am.go**

In `cmd/am/am.go`, remove the import:

```go
certsettingscmd "github.com/gravitee-io/gio-cli/cmd/am/certificate-settings"
```

And remove the line:

```go
cmd.AddCommand(certsettingscmd.NewCertificateSettingsCmd(f))
```

- [ ] **Step 2: Remove CertificateSettingsService from service.go**

In `internal/am/service.go`, remove `CertificateSettingsService` from the `Service` interface embed list.

- [ ] **Step 3: Remove the mock functions**

In `internal/am/mock.go`, remove the `GetCertificateSettingsFunc` field and the `GetCertificateSettings` method.

- [ ] **Step 4: Delete the files**

```bash
rm -rf cmd/am/certificate-settings/
rm internal/am/certificate_settings.go
```

- [ ] **Step 5: Build and test**

Run: `go build ./... && go test ./...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "fix: remove certificate-settings command (endpoint does not exist in AM 4.x)"
```

---

### Task 8: Create E2E docker-compose and test infrastructure

**Files:**
- Create: `e2e/docker-compose.yml`
- Create: `e2e/am_test.go`
- Modify: `hack/make/test.mk`

- [ ] **Step 1: Create docker-compose.yml**

Create `e2e/docker-compose.yml`:

```yaml
services:
  mongodb:
    image: mongo:6.0
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s

  management:
    image: graviteeio/am-management-api:4
    ports:
      - "8093:8093"
      - "18093:18093"
    depends_on:
      mongodb:
        condition: service_healthy
    environment:
      - gravitee_repositories_management_mongodb_uri=mongodb://mongodb:27017/graviteeam?serverSelectionTimeoutMS=5000&connectTimeoutMS=5000&socketTimeoutMS=5000
      - gravitee_repositories_oauth2_mongodb_uri=mongodb://mongodb:27017/graviteeam?serverSelectionTimeoutMS=5000&connectTimeoutMS=5000&socketTimeoutMS=5000
      - gravitee_dataPlanes_0_id=default
      - gravitee_dataPlanes_0_type=mongodb
      - gravitee_dataPlanes_0_mongodb_uri=mongodb://mongodb:27017/graviteeam?serverSelectionTimeoutMS=5000&connectTimeoutMS=5000&socketTimeoutMS=5000
      - gravitee_services_core_http_host=0.0.0.0
    healthcheck:
      test: ["CMD-SHELL", "curl -sf -u admin:adminadmin http://localhost:18093/_node/health || exit 1"]
      interval: 15s
      timeout: 5s
      retries: 15
      start_period: 45s
```

- [ ] **Step 2: Create E2E test file scaffold**

Create `e2e/am_test.go`:

```go
//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var (
	cliBinary string
	amURL     = "http://localhost:8093"
)

func TestMain(m *testing.M) {
	// Build the CLI binary
	binary, err := buildCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build CLI: %v\n", err)
		os.Exit(1)
	}
	cliBinary = binary
	defer os.Remove(binary)

	// Wait for AM to be ready
	if err := waitForAM(amURL, 3*time.Minute); err != nil {
		fmt.Fprintf(os.Stderr, "AM not ready: %v\n", err)
		os.Exit(1)
	}

	// Login
	if err := loginToAM(); err != nil {
		fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func buildCLI() (string, error) {
	binary := "../dist/gio-e2e-test"
	cmd := exec.Command("go", "build", "-o", binary, "..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return binary, cmd.Run()
}

func waitForAM(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(baseURL + "/management/organizations/DEFAULT/environments/DEFAULT/domains")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 401 {
				// AM is up — it returns 401 for unauthenticated requests
				return nil
			}
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("AM did not become ready within %s", timeout)
}

func loginToAM() error {
	// Get token via basic auth
	req, _ := http.NewRequest("POST", amURL+"/management/auth/token", nil)
	req.SetBasicAuth("admin", "adminadmin")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token: %w", err)
	}

	// Login via CLI
	out, err := runCLI("login", "am", "--url", amURL, "--token", tokenResp.AccessToken)
	if err != nil {
		return fmt.Errorf("CLI login failed: %s: %w", out, err)
	}
	return nil
}

func runCLI(args ...string) (string, error) {
	cmd := exec.Command(cliBinary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	output := stdout.String() + stderr.String()
	return strings.TrimSpace(output), err
}

func runCLIExpectSuccess(t *testing.T, args ...string) string {
	t.Helper()
	out, err := runCLI(args...)
	if err != nil {
		t.Fatalf("CLI command failed: %v\nOutput: %s\nArgs: %v", err, out, args)
	}
	return out
}

func runCLIExpectError(t *testing.T, args ...string) string {
	t.Helper()
	out, err := runCLI(args...)
	if err == nil {
		t.Fatalf("expected CLI command to fail, but it succeeded\nOutput: %s\nArgs: %v", out, args)
	}
	return out
}

func extractID(t *testing.T, jsonOutput string) string {
	t.Helper()
	var obj map[string]any
	if err := json.Unmarshal([]byte(jsonOutput), &obj); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, jsonOutput)
	}
	id, ok := obj["id"].(string)
	if !ok {
		t.Fatalf("no 'id' field in JSON output: %s", jsonOutput)
	}
	return id
}
```

- [ ] **Step 3: Add E2E test target to Makefile**

Add to `hack/make/test.mk`:

```makefile
test-e2e:
	cd $(ROOT_DIR) && go test -tags e2e -v -timeout 10m ./e2e/
```

- [ ] **Step 4: Commit**

```bash
git add e2e/docker-compose.yml e2e/am_test.go hack/make/test.mk
git commit -m "feat: add E2E test infrastructure with Docker Compose for AM"
```

---

### Task 9: Write E2E tests for domain CRUD

**Files:**
- Modify: `e2e/am_test.go`

- [ ] **Step 1: Add domain CRUD tests**

Append to `e2e/am_test.go`:

```go
func TestDomainCRUD(t *testing.T) {
	// List domains (should have at least one default)
	t.Run("list domains", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "list")
		if !strings.Contains(out, "Showing") {
			t.Errorf("expected 'Showing' in output, got: %s", out)
		}
	})

	// Create a domain
	var domainID string
	t.Run("create domain", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "create", "--name", "e2e-test-domain", "--description", "E2E test", "-o", "json")
		domainID = extractID(t, out)
		if domainID == "" {
			t.Fatal("domain ID is empty")
		}
	})

	// Get domain
	t.Run("get domain", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "get", domainID)
		if !strings.Contains(out, "e2e-test-domain") {
			t.Errorf("expected domain name in output, got: %s", out)
		}
	})

	// Get domain JSON
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

	// Update domain
	t.Run("update domain", func(t *testing.T) {
		out := runCLIExpectSuccess(t, "am", "domain", "update", domainID, "--name", "e2e-renamed")
		if !strings.Contains(out, "e2e-renamed") {
			t.Errorf("expected updated name, got: %s", out)
		}
	})

	// Enable/disable
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

	// Delete domain
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
```

- [ ] **Step 2: Commit**

```bash
git add e2e/am_test.go
git commit -m "feat: add E2E tests for domain CRUD"
```

---

### Task 10: Write E2E tests for app, user, role, scope, group

**Files:**
- Modify: `e2e/am_test.go`

- [ ] **Step 1: Add app CRUD tests**

Append to `e2e/am_test.go`:

```go
func TestAppCRUD(t *testing.T) {
	// Use the default domain
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
			"--domain", domainID, "--key", "e2e_scope", "--name", "E2E Scope", "-o", "json")
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

// getDefaultDomainID lists domains and returns the first domain ID.
func getDefaultDomainID(t *testing.T) string {
	t.Helper()
	out := runCLIExpectSuccess(t, "am", "domain", "list", "-o", "json")
	var resp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse domain list: %v", err)
	}
	if len(resp.Data) == 0 {
		t.Fatal("no domains found")
	}
	return resp.Data[0].ID
}
```

- [ ] **Step 2: Commit**

```bash
git add e2e/am_test.go
git commit -m "feat: add E2E tests for app, user, role, scope, group CRUD and output formats"
```

---

### Task 11: Create GitHub Actions E2E workflow

**Files:**
- Create: `.github/workflows/e2e.yml`

- [ ] **Step 1: Create the workflow file**

Create `.github/workflows/e2e.yml`:

```yaml
name: E2E

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  e2e-am:
    runs-on: ubuntu-latest
    timeout-minutes: 15
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Start AM services
        working-directory: e2e
        run: docker compose up -d --wait --wait-timeout 180

      - name: Show AM logs on failure
        if: failure()
        working-directory: e2e
        run: docker compose logs management

      - name: Run E2E tests
        run: make test-e2e

      - name: Stop AM services
        if: always()
        working-directory: e2e
        run: docker compose down -v
```

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/e2e.yml
git commit -m "feat: add GitHub Actions E2E workflow with Gravitee AM"
```

---

### Task 12: Run all tests and verify build

**Files:** None (verification only)

- [ ] **Step 1: Run unit tests**

```bash
go test ./...
```

Expected: All PASS.

- [ ] **Step 2: Run linter**

```bash
go vet ./...
```

Expected: No issues.

- [ ] **Step 3: Build**

```bash
go build -o gio .
```

Expected: Success.

- [ ] **Step 4: Quick smoke test the fixed bugs**

```bash
./gio am domain create --name "final-test" -o json   # Should include dataPlaneId
./gio am domain list --page 0                         # Should error, not 500
./gio am app create --domain X --name x --type bad    # Should error before hitting API
./gio am certificate-settings --help                  # Should not exist
```

- [ ] **Step 5: Final commit if any fixes needed**
