# AM Missing Features Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement all features missing from the `am` branch of gio-cli: health, whoami, factor, flow, group, audit, token, and domain export/import/copy.

**Architecture:** All commands follow the `am` branch pattern — direct HTTP calls via `f.Client`, path construction via `cmdutil.AMDomainPath`/`AMEnvPath`, guards via `cmdutil.RequireAMContext`/`RequireAMDomain`. No service layer. Tests use `client.FakeClient` + package-level `newTestFactory` helper.

**Tech Stack:** Go, cobra, `internal/client.GraviteeClient`, `internal/cmdutil`, `internal/printer`

---

## File Structure

### New files to create

```
cmd/am/health.go                  — health command (GET /management/health)
cmd/am/health_test.go             — unit tests for health
cmd/am/whoami.go                  — whoami command (GET /management/user)
cmd/am/whoami_test.go             — unit tests for whoami
cmd/am/factor/factor.go           — factor parent command wiring
cmd/am/factor/list.go             — list MFA factors (GET domain/factors)
cmd/am/factor/get.go              — get single factor (GET domain/factors/{id})
cmd/am/factor/factor_test.go      — unit tests for factor commands
cmd/am/factor/helpers_test.go     — newTestFactory for factor package
cmd/am/flow/flow.go               — flow parent command wiring
cmd/am/flow/list.go               — list flows (GET domain/flows)
cmd/am/flow/get.go                — get single flow (GET domain/flows/{id})
cmd/am/flow/flow_test.go          — unit tests for flow commands
cmd/am/flow/helpers_test.go       — newTestFactory for flow package
cmd/am/group/group.go             — group parent command wiring
cmd/am/group/list.go              — list groups paginated (GET domain/groups)
cmd/am/group/get.go               — get single group (GET domain/groups/{id})
cmd/am/group/create.go            — create group from --file or flags (POST domain/groups)
cmd/am/group/delete.go            — delete group (DELETE domain/groups/{id})
cmd/am/group/group_test.go        — unit tests for group commands
cmd/am/group/helpers_test.go      — newTestFactory for group package
cmd/am/audit/audit.go             — audit parent command wiring
cmd/am/audit/list.go              — list audit logs paginated (GET domain/audits)
cmd/am/audit/get.go               — get single audit log (GET domain/audits/{id})
cmd/am/audit/audit_test.go        — unit tests for audit commands
cmd/am/audit/helpers_test.go      — newTestFactory for audit package
cmd/am/token/token.go             — token parent command wiring
cmd/am/token/list.go              — list user tokens (GET domain/users/{userId}/tokens)
cmd/am/token/create.go            — create user token (POST domain/users/{userId}/tokens)
cmd/am/token/revoke.go            — revoke user token (DELETE domain/users/{userId}/tokens/{id})
cmd/am/token/token_test.go        — unit tests for token commands
cmd/am/token/helpers_test.go      — newTestFactory for token package
cmd/am/domain/export.go           — export domain config to JSON
cmd/am/domain/import.go           — import domain config from JSON file
cmd/am/domain/copy.go             — copy domain (export + import into new domain)
cmd/am/domain/export_test.go      — unit tests for export/import/copy
```

### Files to modify

```
cmd/am/am.go                      — wire health, whoami, factor, flow, group, audit, token
cmd/am/domain/domain.go           — wire export, import, copy
```

---

## Shared Patterns Reference

Every package needs `helpers_test.go`:
```go
package <pkg>

import (
    "bytes"
    "github.com/gravitee-io/gio-cli/internal/client"
    "github.com/gravitee-io/gio-cli/internal/config"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(fc *client.FakeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
    out := &bytes.Buffer{}
    return &factory.Factory{
        Config: &config.Config{
            CurrentContext: "am-test",
            Contexts: map[string]config.Context{
                "am-test": {
                    URL: "https://am-test.com", Token: "tok",
                    Org: "DEFAULT", Env: "DEFAULT",
                    Type: "am", Domain: "test-domain",
                    ReadOnly: readOnly,
                },
            },
        },
        Resolved: &config.ResolvedContext{
            Name: "am-test", URL: "https://am-test.com", Token: "tok",
            Org: "DEFAULT", Env: "DEFAULT",
            Type: "am", Domain: "test-domain",
            ReadOnly: readOnly,
        },
        Client:       fc,
        IOStreams:     factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
        OutputFormat: "table",
    }, out
}
```

Paginated response pattern (copy into each package that lists paginated resources):
```go
type amPaginatedResponse struct {
    Data        []json.RawMessage `json:"data"`
    CurrentPage int               `json:"currentPage"`
    TotalCount  int               `json:"totalCount"`
}
```

---

## Task 1: `gio am health` + `gio am whoami`

**Files:**
- Create: `cmd/am/health.go`
- Create: `cmd/am/health_test.go`
- Create: `cmd/am/whoami.go`
- Create: `cmd/am/whoami_test.go`
- Modify: `cmd/am/am.go`

- [ ] **Step 1: Write failing tests**

Create `cmd/am/health_test.go`:
```go
package am

import (
    "testing"
    "strings"
    "bytes"
    "github.com/gravitee-io/gio-cli/internal/client"
    "github.com/gravitee-io/gio-cli/internal/config"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newAMTestFactory(fc *client.FakeClient) (*factory.Factory, *bytes.Buffer) {
    out := &bytes.Buffer{}
    return &factory.Factory{
        Config: &config.Config{
            CurrentContext: "am-test",
            Contexts: map[string]config.Context{
                "am-test": {URL: "https://am.example.com", Token: "tok", Type: "am", Org: "DEFAULT", Env: "DEFAULT"},
            },
        },
        Resolved: &config.ResolvedContext{
            Name: "am-test", URL: "https://am.example.com", Token: "tok",
            Org: "DEFAULT", Env: "DEFAULT", Type: "am",
        },
        Client:       fc,
        IOStreams:     factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
        OutputFormat: "table",
    }, out
}

func TestHealth(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if path != "/management/health" {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`{"status":"UP"}`), nil
        },
    }
    f, out := newAMTestFactory(fake)
    cmd := newHealthCmd(f)
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "healthy") {
        t.Errorf("expected 'healthy' in output, got: %s", out.String())
    }
}

func TestWhoami(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if path != "/management/user" {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`{"username":"admin","email":"admin@example.com"}`), nil
        },
    }
    f, out := newAMTestFactory(fake)
    cmd := newWhoamiCmd(f)
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "admin") {
        t.Errorf("expected 'admin' in output, got: %s", out.String())
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/rpo/Documents/Projects/Gravitee/gio-cli
go test ./cmd/am/... -run "TestHealth|TestWhoami" -v 2>&1 | tail -10
```
Expected: FAIL with "undefined: newHealthCmd"

- [ ] **Step 3: Implement `cmd/am/health.go`**

```go
package am

import (
    "encoding/json"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

func newHealthCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "health",
        Aliases: []string{"ping"},
        Short:   "Check if the AM instance is reachable",
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireAMContext(f); err != nil {
                return err
            }
            return runHealth(f)
        },
    }
}

func runHealth(f *factory.Factory) error {
    data, err := f.Client.Get("/management/health")
    if err != nil {
        return err
    }

    p := cmdutil.NewPrinter(f)

    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }

    p.PrintMessage("AM instance is healthy.")
    return nil
}
```

- [ ] **Step 4: Implement `cmd/am/whoami.go`**

```go
package am

import (
    "encoding/json"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

func newWhoamiCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:   "whoami",
        Short: "Show information about the currently authenticated user",
        Args:  cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireAMContext(f); err != nil {
                return err
            }
            return runWhoami(f)
        },
    }
}

func runWhoami(f *factory.Factory) error {
    data, err := f.Client.Get("/management/user")
    if err != nil {
        return err
    }

    p := cmdutil.NewPrinter(f)
    return p.PrintDetail(json.RawMessage(data))
}
```

- [ ] **Step 5: Wire into `cmd/am/am.go`**

Add to imports and `NewAMCmd`:
```go
// In imports, no new packages needed — health and whoami are in the same package.

// In NewAMCmd body, add after existing cmd.AddCommand calls:
cmd.AddCommand(newHealthCmd(f))
cmd.AddCommand(newWhoamiCmd(f))
```

- [ ] **Step 6: Run tests to verify they pass**

```bash
go test ./cmd/am/... -run "TestHealth|TestWhoami" -v 2>&1 | tail -10
```
Expected: PASS

- [ ] **Step 7: Run full test suite**

```bash
go test ./... 2>&1 | tail -5
```
Expected: ok for all packages

- [ ] **Step 8: Commit**

```bash
git add cmd/am/health.go cmd/am/health_test.go cmd/am/whoami.go cmd/am/whoami_test.go cmd/am/am.go
git commit -m "feat: add gio am health and whoami commands"
```

---

## Task 2: `gio am factor list/get`

**Files:**
- Create: `cmd/am/factor/factor.go`
- Create: `cmd/am/factor/list.go`
- Create: `cmd/am/factor/get.go`
- Create: `cmd/am/factor/factor_test.go`
- Create: `cmd/am/factor/helpers_test.go`
- Modify: `cmd/am/am.go`

The factor API returns a JSON array (not paginated). `list` returns `[]factor` and `get` returns a single factor object.

- [ ] **Step 1: Write failing tests**

Create `cmd/am/factor/helpers_test.go`:
```go
package factor

import (
    "bytes"
    "github.com/gravitee-io/gio-cli/internal/client"
    "github.com/gravitee-io/gio-cli/internal/config"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(fc *client.FakeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
    out := &bytes.Buffer{}
    return &factory.Factory{
        Config: &config.Config{
            CurrentContext: "am-test",
            Contexts: map[string]config.Context{
                "am-test": {
                    URL: "https://am-test.com", Token: "tok",
                    Org: "DEFAULT", Env: "DEFAULT",
                    Type: "am", Domain: "test-domain",
                    ReadOnly: readOnly,
                },
            },
        },
        Resolved: &config.ResolvedContext{
            Name: "am-test", URL: "https://am-test.com", Token: "tok",
            Org: "DEFAULT", Env: "DEFAULT",
            Type: "am", Domain: "test-domain",
            ReadOnly: readOnly,
        },
        Client:       fc,
        IOStreams:     factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
        OutputFormat: "table",
    }, out
}
```

Create `cmd/am/factor/factor_test.go`:
```go
package factor

import (
    "strings"
    "testing"
    "github.com/gravitee-io/gio-cli/internal/client"
)

func TestFactorList(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if !strings.Contains(path, "/factors") {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`[{"id":"factor-1","name":"SMS Factor","factorType":"SMS"}]`), nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newListCmd(f)
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "SMS Factor") {
        t.Errorf("expected 'SMS Factor' in output, got: %s", out.String())
    }
}

func TestFactorGet(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if !strings.Contains(path, "/factors/factor-1") {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`{"id":"factor-1","name":"SMS Factor","factorType":"SMS"}`), nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newGetCmd(f)
    cmd.SetArgs([]string{"factor-1"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "factor-1") {
        t.Errorf("expected 'factor-1' in output, got: %s", out.String())
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./cmd/am/factor/... -v 2>&1 | tail -10
```
Expected: FAIL with "no Go files"

- [ ] **Step 3: Implement `cmd/am/factor/factor.go`**

```go
package factor

import (
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

// NewFactorCmd creates the parent "gio am factor" command.
func NewFactorCmd(f *factory.Factory) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "factor",
        Short: "Manage MFA factors",
        Args:  cobra.NoArgs,
    }
    cmd.AddCommand(newListCmd(f))
    cmd.AddCommand(newGetCmd(f))
    return cmd
}
```

- [ ] **Step 4: Implement `cmd/am/factor/list.go`**

```go
package factor

import (
    "encoding/json"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

func newListCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "list",
        Short:   "List MFA factors",
        Example: `  gio am factor list`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            return runList(f)
        },
    }
}

func runList(f *factory.Factory) error {
    path := cmdutil.AMDomainPath(f, "factors")
    data, err := f.Client.Get(path)
    if err != nil {
        return err
    }
    p := cmdutil.NewPrinter(f)
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }
    return p.PrintList(json.RawMessage(data), factorColumns())
}

func factorColumns() []printer.Column {
    return []printer.Column{
        {Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
        {Name: "Type", Value: func(i interface{}) string { return cmdutil.StringField(i, "factorType") }},
        {Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
    }
}
```

- [ ] **Step 5: Implement `cmd/am/factor/get.go`**

```go
package factor

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

func newGetCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "get <factorId>",
        Short:   "Get MFA factor details",
        Example: `  gio am factor get my-factor-id`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            return runGet(f, args[0])
        },
    }
}

func runGet(f *factory.Factory, factorID string) error {
    path := cmdutil.AMDomainPath(f, fmt.Sprintf("factors/%s", factorID))
    data, err := f.Client.Get(path)
    if err != nil {
        return err
    }
    p := cmdutil.NewPrinter(f)
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }
    return printFactorDetail(p, data)
}

func printFactorDetail(p *printer.Printer, data []byte) error {
    var m map[string]interface{}
    if err := json.Unmarshal(data, &m); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    for _, field := range []struct{ label, key string }{
        {"Name", "name"}, {"ID", "id"}, {"Type", "factorType"},
    } {
        if v, ok := m[field.key]; ok && v != nil {
            p.PrintMessage("%-16s%v", field.label+":", v)
        }
    }
    return nil
}
```

- [ ] **Step 6: Wire into `cmd/am/am.go`**

Add to imports:
```go
factorcmd "github.com/gravitee-io/gio-cli/cmd/am/factor"
```

Add to `NewAMCmd`:
```go
cmd.AddCommand(factorcmd.NewFactorCmd(f))
```

- [ ] **Step 7: Run tests to verify they pass**

```bash
go test ./cmd/am/factor/... -v 2>&1 | tail -10
go test ./... 2>&1 | tail -5
```
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add cmd/am/factor/ cmd/am/am.go
git commit -m "feat: add gio am factor list/get commands"
```

---

## Task 3: `gio am flow list/get`

**Files:**
- Create: `cmd/am/flow/flow.go`
- Create: `cmd/am/flow/list.go`
- Create: `cmd/am/flow/get.go`
- Create: `cmd/am/flow/flow_test.go`
- Create: `cmd/am/flow/helpers_test.go`
- Modify: `cmd/am/am.go`

The flow API returns a JSON array (not paginated).

- [ ] **Step 1: Write failing tests**

Create `cmd/am/flow/helpers_test.go`:
```go
package flow

import (
    "bytes"
    "github.com/gravitee-io/gio-cli/internal/client"
    "github.com/gravitee-io/gio-cli/internal/config"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(fc *client.FakeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
    out := &bytes.Buffer{}
    return &factory.Factory{
        Config: &config.Config{
            CurrentContext: "am-test",
            Contexts: map[string]config.Context{
                "am-test": {
                    URL: "https://am-test.com", Token: "tok",
                    Org: "DEFAULT", Env: "DEFAULT",
                    Type: "am", Domain: "test-domain",
                    ReadOnly: readOnly,
                },
            },
        },
        Resolved: &config.ResolvedContext{
            Name: "am-test", URL: "https://am-test.com", Token: "tok",
            Org: "DEFAULT", Env: "DEFAULT",
            Type: "am", Domain: "test-domain",
            ReadOnly: readOnly,
        },
        Client:       fc,
        IOStreams:     factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
        OutputFormat: "table",
    }, out
}
```

Create `cmd/am/flow/flow_test.go`:
```go
package flow

import (
    "strings"
    "testing"
    "github.com/gravitee-io/gio-cli/internal/client"
)

func TestFlowList(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if !strings.Contains(path, "/flows") {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`[{"id":"flow-1","name":"Login Flow","type":"ROOT"}]`), nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newListCmd(f)
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "Login Flow") {
        t.Errorf("expected 'Login Flow' in output, got: %s", out.String())
    }
}

func TestFlowGet(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if !strings.Contains(path, "/flows/flow-1") {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`{"id":"flow-1","name":"Login Flow","type":"ROOT"}`), nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newGetCmd(f)
    cmd.SetArgs([]string{"flow-1"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "flow-1") {
        t.Errorf("expected 'flow-1' in output, got: %s", out.String())
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./cmd/am/flow/... -v 2>&1 | tail -10
```
Expected: FAIL with "no Go files"

- [ ] **Step 3: Implement `cmd/am/flow/flow.go`**

```go
package flow

import (
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

// NewFlowCmd creates the parent "gio am flow" command.
func NewFlowCmd(f *factory.Factory) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "flow",
        Short: "Manage authentication flows",
        Args:  cobra.NoArgs,
    }
    cmd.AddCommand(newListCmd(f))
    cmd.AddCommand(newGetCmd(f))
    return cmd
}
```

- [ ] **Step 4: Implement `cmd/am/flow/list.go`**

```go
package flow

import (
    "encoding/json"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

func newListCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "list",
        Short:   "List authentication flows",
        Example: `  gio am flow list`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            return runList(f)
        },
    }
}

func runList(f *factory.Factory) error {
    path := cmdutil.AMDomainPath(f, "flows")
    data, err := f.Client.Get(path)
    if err != nil {
        return err
    }
    p := cmdutil.NewPrinter(f)
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }
    return p.PrintList(json.RawMessage(data), flowColumns())
}

func flowColumns() []printer.Column {
    return []printer.Column{
        {Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
        {Name: "Type", Value: func(i interface{}) string { return cmdutil.StringField(i, "type") }},
        {Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
    }
}
```

- [ ] **Step 5: Implement `cmd/am/flow/get.go`**

```go
package flow

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

func newGetCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "get <flowId>",
        Short:   "Get authentication flow details",
        Example: `  gio am flow get my-flow-id`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            return runGet(f, args[0])
        },
    }
}

func runGet(f *factory.Factory, flowID string) error {
    path := cmdutil.AMDomainPath(f, fmt.Sprintf("flows/%s", flowID))
    data, err := f.Client.Get(path)
    if err != nil {
        return err
    }
    p := cmdutil.NewPrinter(f)
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }
    return printFlowDetail(p, data)
}

func printFlowDetail(p *printer.Printer, data []byte) error {
    var m map[string]interface{}
    if err := json.Unmarshal(data, &m); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    for _, field := range []struct{ label, key string }{
        {"Name", "name"}, {"ID", "id"}, {"Type", "type"},
    } {
        if v, ok := m[field.key]; ok && v != nil {
            p.PrintMessage("%-16s%v", field.label+":", v)
        }
    }
    return nil
}
```

- [ ] **Step 6: Wire into `cmd/am/am.go`**

Add to imports:
```go
flowcmd "github.com/gravitee-io/gio-cli/cmd/am/flow"
```

Add to `NewAMCmd`:
```go
cmd.AddCommand(flowcmd.NewFlowCmd(f))
```

- [ ] **Step 7: Run tests and full suite**

```bash
go test ./cmd/am/flow/... -v 2>&1 | tail -10
go test ./... 2>&1 | tail -5
```
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add cmd/am/flow/ cmd/am/am.go
git commit -m "feat: add gio am flow list/get commands"
```

---

## Task 4: `gio am group list/get/create/delete`

**Files:**
- Create: `cmd/am/group/group.go`
- Create: `cmd/am/group/list.go`
- Create: `cmd/am/group/get.go`
- Create: `cmd/am/group/create.go`
- Create: `cmd/am/group/delete.go`
- Create: `cmd/am/group/group_test.go`
- Create: `cmd/am/group/helpers_test.go`
- Modify: `cmd/am/am.go`

Groups use paginated list (same `amPaginatedResponse` pattern as roles/users).

- [ ] **Step 1: Write failing tests**

Create `cmd/am/group/helpers_test.go`:
```go
package group

import (
    "bytes"
    "github.com/gravitee-io/gio-cli/internal/client"
    "github.com/gravitee-io/gio-cli/internal/config"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(fc *client.FakeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
    out := &bytes.Buffer{}
    return &factory.Factory{
        Config: &config.Config{
            CurrentContext: "am-test",
            Contexts: map[string]config.Context{
                "am-test": {
                    URL: "https://am-test.com", Token: "tok",
                    Org: "DEFAULT", Env: "DEFAULT",
                    Type: "am", Domain: "test-domain",
                    ReadOnly: readOnly,
                },
            },
        },
        Resolved: &config.ResolvedContext{
            Name: "am-test", URL: "https://am-test.com", Token: "tok",
            Org: "DEFAULT", Env: "DEFAULT",
            Type: "am", Domain: "test-domain",
            ReadOnly: readOnly,
        },
        Client:       fc,
        IOStreams:     factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
        OutputFormat: "table",
    }, out
}
```

Create `cmd/am/group/group_test.go`:
```go
package group

import (
    "encoding/json"
    "strings"
    "testing"

    "github.com/gravitee-io/gio-cli/internal/client"
)

func TestGroupList(t *testing.T) {
    resp := map[string]interface{}{
        "data":        []map[string]interface{}{{"id": "group-1", "name": "Admins", "description": "Admin group"}},
        "currentPage": 0,
        "totalCount":  1,
    }
    data, _ := json.Marshal(resp)
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if !strings.Contains(path, "/groups?") {
                t.Errorf("unexpected path: %s", path)
            }
            return data, nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newListCmd(f)
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "Admins") {
        t.Errorf("expected 'Admins' in output, got: %s", out.String())
    }
}

func TestGroupGet(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if !strings.Contains(path, "/groups/group-1") {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`{"id":"group-1","name":"Admins","description":"Admin group"}`), nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newGetCmd(f)
    cmd.SetArgs([]string{"group-1"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "group-1") {
        t.Errorf("expected 'group-1' in output, got: %s", out.String())
    }
}

func TestGroupCreate(t *testing.T) {
    fake := &client.FakeClient{
        PostFunc: func(path string, body interface{}) ([]byte, error) {
            if !strings.Contains(path, "/groups") {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`{"id":"group-new","name":"DevTeam"}`), nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newCreateCmd(f)
    cmd.SetArgs([]string{"--name", "DevTeam", "--description", "Developers"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "DevTeam") {
        t.Errorf("expected 'DevTeam' in output, got: %s", out.String())
    }
}

func TestGroupDelete(t *testing.T) {
    deleted := false
    fake := &client.FakeClient{
        DeleteFunc: func(path string) error {
            if !strings.Contains(path, "/groups/group-1") {
                t.Errorf("unexpected path: %s", path)
            }
            deleted = true
            return nil
        },
    }
    f, _ := newTestFactory(fake, false)
    cmd := newDeleteCmd(f)
    cmd.SetArgs([]string{"group-1"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !deleted {
        t.Error("expected Delete to be called")
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./cmd/am/group/... -v 2>&1 | tail -10
```
Expected: FAIL with "no Go files"

- [ ] **Step 3: Implement `cmd/am/group/group.go`**

```go
package group

import (
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

// NewGroupCmd creates the parent "gio am group" command.
func NewGroupCmd(f *factory.Factory) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "group",
        Short: "Manage groups",
        Args:  cobra.NoArgs,
    }
    cmd.AddCommand(newListCmd(f))
    cmd.AddCommand(newGetCmd(f))
    cmd.AddCommand(newCreateCmd(f))
    cmd.AddCommand(newDeleteCmd(f))
    return cmd
}
```

- [ ] **Step 4: Implement `cmd/am/group/list.go`**

```go
package group

import (
    "encoding/json"
    "fmt"
    "net/url"
    "strconv"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

type listOptions struct {
    factory *factory.Factory
    query   string
    page    int
    size    int
    all     bool
}

type amPaginatedResponse struct {
    Data        []json.RawMessage `json:"data"`
    CurrentPage int               `json:"currentPage"`
    TotalCount  int               `json:"totalCount"`
}

func newListCmd(f *factory.Factory) *cobra.Command {
    opts := &listOptions{factory: f}
    cmd := &cobra.Command{
        Use:     "list",
        Short:   "List groups",
        Example: `  gio am group list`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            return opts.run()
        },
    }
    cmd.Flags().IntVarP(&opts.page, "page", "p", 0, "Page number (0-based)")
    cmd.Flags().IntVarP(&opts.size, "size", "s", 20, "Results per page")
    cmd.Flags().StringVarP(&opts.query, "query", "q", "", "Search by name")
    cmd.Flags().BoolVarP(&opts.all, "all", "a", false, "Fetch all pages")
    return cmd
}

func (o *listOptions) buildQuery(page int) string {
    q := url.Values{}
    q.Set("page", strconv.Itoa(page))
    q.Set("size", strconv.Itoa(o.size))
    if o.query != "" {
        q.Set("q", o.query)
    }
    return q.Encode()
}

func (o *listOptions) run() error {
    f := o.factory
    p := cmdutil.NewPrinter(f)
    if o.all {
        return o.fetchAll(f, p)
    }
    return o.fetchPage(f, p, o.page)
}

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
    path := cmdutil.AMDomainPath(f, "groups?"+o.buildQuery(page))
    data, err := f.Client.Get(path)
    if err != nil {
        return err
    }
    var resp amPaginatedResponse
    if err := json.Unmarshal(data, &resp); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }
    if err := p.PrintList(resp.Data, groupColumns()); err != nil {
        return err
    }
    if resp.TotalCount > 0 {
        p.PrintMessage("Showing %d of %d total.", len(resp.Data), resp.TotalCount)
    }
    return nil
}

func (o *listOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
    var allData []json.RawMessage
    size := 100
    for page := 0; page <= 1000; page++ {
        q := url.Values{}
        q.Set("page", strconv.Itoa(page))
        q.Set("size", strconv.Itoa(size))
        if o.query != "" {
            q.Set("q", o.query)
        }
        path := cmdutil.AMDomainPath(f, "groups?"+q.Encode())
        data, err := f.Client.Get(path)
        if err != nil {
            return err
        }
        var resp amPaginatedResponse
        if err := json.Unmarshal(data, &resp); err != nil {
            return fmt.Errorf("failed to parse response: %w", err)
        }
        allData = append(allData, resp.Data...)
        if len(allData) >= resp.TotalCount || len(resp.Data) < size {
            break
        }
    }
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(allData)
    }
    if err := p.PrintList(allData, groupColumns()); err != nil {
        return err
    }
    if len(allData) > 0 {
        p.PrintMessage("Showing %d results.", len(allData))
    }
    return nil
}

func groupColumns() []printer.Column {
    return []printer.Column{
        {Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
        {Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
        {Name: "Description", Value: func(i interface{}) string { return cmdutil.StringField(i, "description") }},
    }
}
```

- [ ] **Step 5: Implement `cmd/am/group/get.go`**

```go
package group

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

func newGetCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "get <groupId>",
        Short:   "Get group details",
        Example: `  gio am group get my-group-id`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            return runGet(f, args[0])
        },
    }
}

func runGet(f *factory.Factory, groupID string) error {
    path := cmdutil.AMDomainPath(f, fmt.Sprintf("groups/%s", groupID))
    data, err := f.Client.Get(path)
    if err != nil {
        return err
    }
    p := cmdutil.NewPrinter(f)
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }
    return printGroupDetail(p, data)
}

func printGroupDetail(p *printer.Printer, data []byte) error {
    var m map[string]interface{}
    if err := json.Unmarshal(data, &m); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    for _, field := range []struct{ label, key string }{
        {"Name", "name"}, {"ID", "id"}, {"Description", "description"},
    } {
        if v, ok := m[field.key]; ok && v != nil {
            p.PrintMessage("%-16s%v", field.label+":", v)
        }
    }
    return nil
}
```

- [ ] **Step 6: Implement `cmd/am/group/create.go`**

```go
package group

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
    var name, description, file string
    cmd := &cobra.Command{
        Use:     "create",
        Short:   "Create a group",
        Example: `  gio am group create --name Admins --description "Admin group"`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            if err := cmdutil.CheckReadOnly(f, "group create"); err != nil {
                return err
            }
            return runCreate(f, name, description, file)
        },
    }
    cmd.Flags().StringVar(&name, "name", "", "Group name")
    cmd.Flags().StringVar(&description, "description", "", "Group description")
    cmd.Flags().StringVarP(&file, "file", "f", "", "JSON file with group definition (overrides flags)")
    return cmd
}

func runCreate(f *factory.Factory, name, description, file string) error {
    var body json.RawMessage
    if file != "" {
        var err error
        body, err = cmdutil.ReadJSONFile(file)
        if err != nil {
            return err
        }
    } else {
        if name == "" {
            return fmt.Errorf("--name is required when --file is not provided")
        }
        payload := map[string]interface{}{"name": name}
        if description != "" {
            payload["description"] = description
        }
        b, err := json.Marshal(payload)
        if err != nil {
            return fmt.Errorf("failed to build request body: %w", err)
        }
        body = b
    }

    path := cmdutil.AMDomainPath(f, "groups")
    data, err := f.Client.Post(path, body)
    if err != nil {
        return err
    }

    p := cmdutil.NewPrinter(f)
    var m map[string]interface{}
    if err := json.Unmarshal(data, &m); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    p.PrintMessage("Group '%s' created (ID: %s).", cmdutil.StringField(m, "name"), cmdutil.StringField(m, "id"))
    return nil
}
```

- [ ] **Step 7: Implement `cmd/am/group/delete.go`**

```go
package group

import (
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "delete <groupId>",
        Short:   "Delete a group",
        Example: `  gio am group delete my-group-id`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            if err := cmdutil.CheckReadOnly(f, "group delete"); err != nil {
                return err
            }
            return runDelete(f, args[0])
        },
    }
}

func runDelete(f *factory.Factory, groupID string) error {
    path := cmdutil.AMDomainPath(f, fmt.Sprintf("groups/%s", groupID))
    if err := f.Client.Delete(path); err != nil {
        return err
    }
    p := cmdutil.NewPrinter(f)
    p.PrintMessage("Group '%s' deleted.", groupID)
    return nil
}
```

- [ ] **Step 8: Wire into `cmd/am/am.go`**

Add to imports:
```go
groupcmd "github.com/gravitee-io/gio-cli/cmd/am/group"
```

Add to `NewAMCmd`:
```go
cmd.AddCommand(groupcmd.NewGroupCmd(f))
```

- [ ] **Step 9: Run tests and full suite**

```bash
go test ./cmd/am/group/... -v 2>&1 | tail -15
go test ./... 2>&1 | tail -5
```
Expected: PASS

- [ ] **Step 10: Commit**

```bash
git add cmd/am/group/ cmd/am/am.go
git commit -m "feat: add gio am group list/get/create/delete commands"
```

---

## Task 5: `gio am audit list/get`

**Files:**
- Create: `cmd/am/audit/audit.go`
- Create: `cmd/am/audit/list.go`
- Create: `cmd/am/audit/get.go`
- Create: `cmd/am/audit/audit_test.go`
- Create: `cmd/am/audit/helpers_test.go`
- Modify: `cmd/am/am.go`

Audits use paginated list.

- [ ] **Step 1: Write failing tests**

Create `cmd/am/audit/helpers_test.go`:
```go
package audit

import (
    "bytes"
    "github.com/gravitee-io/gio-cli/internal/client"
    "github.com/gravitee-io/gio-cli/internal/config"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(fc *client.FakeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
    out := &bytes.Buffer{}
    return &factory.Factory{
        Config: &config.Config{
            CurrentContext: "am-test",
            Contexts: map[string]config.Context{
                "am-test": {
                    URL: "https://am-test.com", Token: "tok",
                    Org: "DEFAULT", Env: "DEFAULT",
                    Type: "am", Domain: "test-domain",
                    ReadOnly: readOnly,
                },
            },
        },
        Resolved: &config.ResolvedContext{
            Name: "am-test", URL: "https://am-test.com", Token: "tok",
            Org: "DEFAULT", Env: "DEFAULT",
            Type: "am", Domain: "test-domain",
            ReadOnly: readOnly,
        },
        Client:       fc,
        IOStreams:     factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
        OutputFormat: "table",
    }, out
}
```

Create `cmd/am/audit/audit_test.go`:
```go
package audit

import (
    "encoding/json"
    "strings"
    "testing"

    "github.com/gravitee-io/gio-cli/internal/client"
)

func TestAuditList(t *testing.T) {
    resp := map[string]interface{}{
        "data": []map[string]interface{}{
            {"id": "audit-1", "type": "USER_LOGIN", "outcome": map[string]interface{}{"status": "SUCCESS"}},
        },
        "currentPage": 0,
        "totalCount":  1,
    }
    data, _ := json.Marshal(resp)
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if !strings.Contains(path, "/audits?") {
                t.Errorf("unexpected path: %s", path)
            }
            return data, nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newListCmd(f)
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "USER_LOGIN") {
        t.Errorf("expected 'USER_LOGIN' in output, got: %s", out.String())
    }
}

func TestAuditGet(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if !strings.Contains(path, "/audits/audit-1") {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`{"id":"audit-1","type":"USER_LOGIN"}`), nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newGetCmd(f)
    cmd.SetArgs([]string{"audit-1"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "audit-1") {
        t.Errorf("expected 'audit-1' in output, got: %s", out.String())
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./cmd/am/audit/... -v 2>&1 | tail -10
```
Expected: FAIL with "no Go files"

- [ ] **Step 3: Implement `cmd/am/audit/audit.go`**

```go
package audit

import (
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

// NewAuditCmd creates the parent "gio am audit" command.
func NewAuditCmd(f *factory.Factory) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "audit",
        Short: "View audit logs",
        Args:  cobra.NoArgs,
    }
    cmd.AddCommand(newListCmd(f))
    cmd.AddCommand(newGetCmd(f))
    return cmd
}
```

- [ ] **Step 4: Implement `cmd/am/audit/list.go`**

```go
package audit

import (
    "encoding/json"
    "fmt"
    "net/url"
    "strconv"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

type listOptions struct {
    factory    *factory.Factory
    auditType  string
    status     string
    page       int
    size       int
    all        bool
}

type amPaginatedResponse struct {
    Data        []json.RawMessage `json:"data"`
    CurrentPage int               `json:"currentPage"`
    TotalCount  int               `json:"totalCount"`
}

func newListCmd(f *factory.Factory) *cobra.Command {
    opts := &listOptions{factory: f}
    cmd := &cobra.Command{
        Use:     "list",
        Short:   "List audit logs",
        Example: `  gio am audit list`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            return opts.run()
        },
    }
    cmd.Flags().IntVarP(&opts.page, "page", "p", 0, "Page number (0-based)")
    cmd.Flags().IntVarP(&opts.size, "size", "s", 20, "Results per page")
    cmd.Flags().StringVar(&opts.auditType, "type", "", "Filter by audit type")
    cmd.Flags().StringVar(&opts.status, "status", "", "Filter by status (SUCCESS, FAILURE)")
    cmd.Flags().BoolVarP(&opts.all, "all", "a", false, "Fetch all pages")
    return cmd
}

func (o *listOptions) buildQuery(page int) string {
    q := url.Values{}
    q.Set("page", strconv.Itoa(page))
    q.Set("size", strconv.Itoa(o.size))
    if o.auditType != "" {
        q.Set("type", o.auditType)
    }
    if o.status != "" {
        q.Set("status", o.status)
    }
    return q.Encode()
}

func (o *listOptions) run() error {
    f := o.factory
    p := cmdutil.NewPrinter(f)
    if o.all {
        return o.fetchAll(f, p)
    }
    return o.fetchPage(f, p, o.page)
}

func (o *listOptions) fetchPage(f *factory.Factory, p *printer.Printer, page int) error {
    path := cmdutil.AMDomainPath(f, "audits?"+o.buildQuery(page))
    data, err := f.Client.Get(path)
    if err != nil {
        return err
    }
    var resp amPaginatedResponse
    if err := json.Unmarshal(data, &resp); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }
    if err := p.PrintList(resp.Data, auditColumns()); err != nil {
        return err
    }
    if resp.TotalCount > 0 {
        p.PrintMessage("Showing %d of %d total.", len(resp.Data), resp.TotalCount)
    }
    return nil
}

func (o *listOptions) fetchAll(f *factory.Factory, p *printer.Printer) error {
    var allData []json.RawMessage
    size := 100
    for page := 0; page <= 1000; page++ {
        q := url.Values{}
        q.Set("page", strconv.Itoa(page))
        q.Set("size", strconv.Itoa(size))
        if o.auditType != "" {
            q.Set("type", o.auditType)
        }
        if o.status != "" {
            q.Set("status", o.status)
        }
        path := cmdutil.AMDomainPath(f, "audits?"+q.Encode())
        data, err := f.Client.Get(path)
        if err != nil {
            return err
        }
        var resp amPaginatedResponse
        if err := json.Unmarshal(data, &resp); err != nil {
            return fmt.Errorf("failed to parse response: %w", err)
        }
        allData = append(allData, resp.Data...)
        if len(allData) >= resp.TotalCount || len(resp.Data) < size {
            break
        }
    }
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(allData)
    }
    if err := p.PrintList(allData, auditColumns()); err != nil {
        return err
    }
    if len(allData) > 0 {
        p.PrintMessage("Showing %d results.", len(allData))
    }
    return nil
}

func auditColumns() []printer.Column {
    return []printer.Column{
        {Name: "Type", Value: func(i interface{}) string { return cmdutil.StringField(i, "type") }},
        {Name: "Status", Value: func(i interface{}) string {
            m, ok := i.(map[string]interface{})
            if !ok {
                return ""
            }
            if outcome, ok := m["outcome"].(map[string]interface{}); ok {
                return cmdutil.StringField(outcome, "status")
            }
            return ""
        }},
        {Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
    }
}
```

- [ ] **Step 5: Implement `cmd/am/audit/get.go`**

```go
package audit

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

func newGetCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "get <auditId>",
        Short:   "Get audit log details",
        Example: `  gio am audit get my-audit-id`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            return runGet(f, args[0])
        },
    }
}

func runGet(f *factory.Factory, auditID string) error {
    path := cmdutil.AMDomainPath(f, fmt.Sprintf("audits/%s", auditID))
    data, err := f.Client.Get(path)
    if err != nil {
        return err
    }
    p := cmdutil.NewPrinter(f)
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }
    return printAuditDetail(p, data)
}

func printAuditDetail(p *printer.Printer, data []byte) error {
    var m map[string]interface{}
    if err := json.Unmarshal(data, &m); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    p.PrintMessage("%-16s%v", "ID:", cmdutil.StringField(m, "id"))
    p.PrintMessage("%-16s%v", "Type:", cmdutil.StringField(m, "type"))
    if outcome, ok := m["outcome"].(map[string]interface{}); ok {
        p.PrintMessage("%-16s%v", "Status:", cmdutil.StringField(outcome, "status"))
        if msg := cmdutil.StringField(outcome, "message"); msg != "" {
            p.PrintMessage("%-16s%v", "Message:", msg)
        }
    }
    return nil
}
```

- [ ] **Step 6: Wire into `cmd/am/am.go`**

Add to imports:
```go
auditcmd "github.com/gravitee-io/gio-cli/cmd/am/audit"
```

Add to `NewAMCmd`:
```go
cmd.AddCommand(auditcmd.NewAuditCmd(f))
```

- [ ] **Step 7: Run tests and full suite**

```bash
go test ./cmd/am/audit/... -v 2>&1 | tail -10
go test ./... 2>&1 | tail -5
```
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add cmd/am/audit/ cmd/am/am.go
git commit -m "feat: add gio am audit list/get commands"
```

---

## Task 6: `gio am token list/create/revoke`

**Files:**
- Create: `cmd/am/token/token.go`
- Create: `cmd/am/token/list.go`
- Create: `cmd/am/token/create.go`
- Create: `cmd/am/token/revoke.go`
- Create: `cmd/am/token/token_test.go`
- Create: `cmd/am/token/helpers_test.go`
- Modify: `cmd/am/am.go`

Token API: `GET/POST/DELETE /domains/{domainId}/users/{userId}/tokens`. Response is a JSON array (not paginated).

- [ ] **Step 1: Write failing tests**

Create `cmd/am/token/helpers_test.go`:
```go
package token

import (
    "bytes"
    "github.com/gravitee-io/gio-cli/internal/client"
    "github.com/gravitee-io/gio-cli/internal/config"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(fc *client.FakeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
    out := &bytes.Buffer{}
    return &factory.Factory{
        Config: &config.Config{
            CurrentContext: "am-test",
            Contexts: map[string]config.Context{
                "am-test": {
                    URL: "https://am-test.com", Token: "tok",
                    Org: "DEFAULT", Env: "DEFAULT",
                    Type: "am", Domain: "test-domain",
                    ReadOnly: readOnly,
                },
            },
        },
        Resolved: &config.ResolvedContext{
            Name: "am-test", URL: "https://am-test.com", Token: "tok",
            Org: "DEFAULT", Env: "DEFAULT",
            Type: "am", Domain: "test-domain",
            ReadOnly: readOnly,
        },
        Client:       fc,
        IOStreams:     factory.IOStreams{Out: out, Err: &bytes.Buffer{}},
        OutputFormat: "table",
    }, out
}
```

Create `cmd/am/token/token_test.go`:
```go
package token

import (
    "strings"
    "testing"

    "github.com/gravitee-io/gio-cli/internal/client"
)

func TestTokenList(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if !strings.Contains(path, "/users/user-1/tokens") {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`[{"id":"token-1","token":"abc","createdAt":1000}]`), nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newListCmd(f)
    cmd.SetArgs([]string{"--user", "user-1"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "token-1") {
        t.Errorf("expected 'token-1' in output, got: %s", out.String())
    }
}

func TestTokenCreate(t *testing.T) {
    fake := &client.FakeClient{
        PostFunc: func(path string, body interface{}) ([]byte, error) {
            if !strings.Contains(path, "/users/user-1/tokens") {
                t.Errorf("unexpected path: %s", path)
            }
            return []byte(`{"id":"token-new","token":"xyz"}`), nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newCreateCmd(f)
    cmd.SetArgs([]string{"--user", "user-1"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "token-new") {
        t.Errorf("expected 'token-new' in output, got: %s", out.String())
    }
}

func TestTokenRevoke(t *testing.T) {
    revoked := false
    fake := &client.FakeClient{
        DeleteFunc: func(path string) error {
            if !strings.Contains(path, "/users/user-1/tokens/token-1") {
                t.Errorf("unexpected path: %s", path)
            }
            revoked = true
            return nil
        },
    }
    f, _ := newTestFactory(fake, false)
    cmd := newRevokeCmd(f)
    cmd.SetArgs([]string{"token-1", "--user", "user-1"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !revoked {
        t.Error("expected Delete to be called")
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./cmd/am/token/... -v 2>&1 | tail -10
```
Expected: FAIL with "no Go files"

- [ ] **Step 3: Implement `cmd/am/token/token.go`**

```go
package token

import (
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

// NewTokenCmd creates the parent "gio am token" command.
func NewTokenCmd(f *factory.Factory) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "token",
        Short: "Manage user tokens",
        Args:  cobra.NoArgs,
    }
    cmd.AddCommand(newListCmd(f))
    cmd.AddCommand(newCreateCmd(f))
    cmd.AddCommand(newRevokeCmd(f))
    return cmd
}
```

- [ ] **Step 4: Implement `cmd/am/token/list.go`**

```go
package token

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
    "github.com/spf13/cobra"
)

func newListCmd(f *factory.Factory) *cobra.Command {
    var userID string
    cmd := &cobra.Command{
        Use:     "list",
        Short:   "List user tokens",
        Example: `  gio am token list --user user-uuid`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            return runList(f, userID)
        },
    }
    cmd.Flags().StringVar(&userID, "user", "", "User ID (required)")
    _ = cmd.MarkFlagRequired("user")
    return cmd
}

func runList(f *factory.Factory, userID string) error {
    path := cmdutil.AMDomainPath(f, fmt.Sprintf("users/%s/tokens", userID))
    data, err := f.Client.Get(path)
    if err != nil {
        return err
    }
    p := cmdutil.NewPrinter(f)
    if f.OutputFormat != printer.FormatTable {
        return p.PrintDetail(json.RawMessage(data))
    }
    return p.PrintList(json.RawMessage(data), tokenColumns())
}

func tokenColumns() []printer.Column {
    return []printer.Column{
        {Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
        {Name: "Token", Value: func(i interface{}) string { return cmdutil.StringField(i, "token") }},
    }
}
```

- [ ] **Step 5: Implement `cmd/am/token/create.go`**

```go
package token

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
    var userID, file string
    cmd := &cobra.Command{
        Use:     "create",
        Short:   "Create a user token",
        Example: `  gio am token create --user user-uuid`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            if err := cmdutil.CheckReadOnly(f, "token create"); err != nil {
                return err
            }
            return runCreate(f, userID, file)
        },
    }
    cmd.Flags().StringVar(&userID, "user", "", "User ID (required)")
    cmd.Flags().StringVarP(&file, "file", "f", "", "JSON file with token definition")
    _ = cmd.MarkFlagRequired("user")
    return cmd
}

func runCreate(f *factory.Factory, userID, file string) error {
    var body json.RawMessage
    if file != "" {
        var err error
        body, err = cmdutil.ReadJSONFile(file)
        if err != nil {
            return err
        }
    } else {
        body = json.RawMessage(`{}`)
    }

    path := cmdutil.AMDomainPath(f, fmt.Sprintf("users/%s/tokens", userID))
    data, err := f.Client.Post(path, body)
    if err != nil {
        return err
    }

    p := cmdutil.NewPrinter(f)
    var m map[string]interface{}
    if err := json.Unmarshal(data, &m); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    p.PrintMessage("Token created (ID: %s).", cmdutil.StringField(m, "id"))
    return nil
}
```

- [ ] **Step 6: Implement `cmd/am/token/revoke.go`**

```go
package token

import (
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

func newRevokeCmd(f *factory.Factory) *cobra.Command {
    var userID string
    cmd := &cobra.Command{
        Use:     "revoke <tokenId>",
        Short:   "Revoke a user token",
        Example: `  gio am token revoke token-id --user user-uuid`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireAMDomain(f); err != nil {
                return err
            }
            if err := cmdutil.CheckReadOnly(f, "token revoke"); err != nil {
                return err
            }
            return runRevoke(f, userID, args[0])
        },
    }
    cmd.Flags().StringVar(&userID, "user", "", "User ID (required)")
    _ = cmd.MarkFlagRequired("user")
    return cmd
}

func runRevoke(f *factory.Factory, userID, tokenID string) error {
    path := cmdutil.AMDomainPath(f, fmt.Sprintf("users/%s/tokens/%s", userID, tokenID))
    if err := f.Client.Delete(path); err != nil {
        return err
    }
    p := cmdutil.NewPrinter(f)
    p.PrintMessage("Token '%s' revoked.", tokenID)
    return nil
}
```

- [ ] **Step 7: Wire into `cmd/am/am.go`**

Add to imports:
```go
tokencmd "github.com/gravitee-io/gio-cli/cmd/am/token"
```

Add to `NewAMCmd`:
```go
cmd.AddCommand(tokencmd.NewTokenCmd(f))
```

- [ ] **Step 8: Run tests and full suite**

```bash
go test ./cmd/am/token/... -v 2>&1 | tail -10
go test ./... 2>&1 | tail -5
```
Expected: PASS

- [ ] **Step 9: Commit**

```bash
git add cmd/am/token/ cmd/am/am.go
git commit -m "feat: add gio am token list/create/revoke commands"
```

---

## Task 7: `gio am domain export/import/copy`

**Files:**
- Create: `cmd/am/domain/export.go`
- Create: `cmd/am/domain/import.go`
- Create: `cmd/am/domain/copy.go`
- Create: `cmd/am/domain/export_test.go`
- Modify: `cmd/am/domain/domain.go`

Export fetches domain + all child resources concurrently. Import reads a JSON export file and recreates resources. Copy = export + import into a newly created domain.

All three commands use `RequireAMContext(f)` (not `RequireAMDomain`) because the domain ID is passed as a CLI argument, not from config.

- [ ] **Step 1: Write failing tests**

Create `cmd/am/domain/export_test.go`:
```go
package domain

import (
    "encoding/json"
    "strings"
    "testing"

    "github.com/gravitee-io/gio-cli/internal/client"
)

func TestDomainExport(t *testing.T) {
    callCount := 0
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            callCount++
            switch {
            case strings.HasSuffix(path, "/domains/dom-1"):
                return []byte(`{"id":"dom-1","name":"Test Domain"}`), nil
            case strings.Contains(path, "/applications"):
                return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
            case strings.Contains(path, "/identityProviders"):
                return []byte(`[]`), nil
            case strings.Contains(path, "/roles"):
                return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
            case strings.Contains(path, "/scopes"):
                return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
            case strings.Contains(path, "/factors"):
                return []byte(`[]`), nil
            case strings.Contains(path, "/groups"):
                return []byte(`{"data":[],"currentPage":0,"totalCount":0}`), nil
            case strings.Contains(path, "/flows"):
                return []byte(`[]`), nil
            case strings.Contains(path, "/certificates"):
                return []byte(`[]`), nil
            }
            return nil, nil
        },
    }
    f, out := newTestFactory(fake, false)
    cmd := newExportCmd(f)
    cmd.SetArgs([]string{"dom-1"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    var result map[string]json.RawMessage
    if err := json.Unmarshal(out.Bytes(), &result); err != nil {
        t.Fatalf("export output is not valid JSON: %v", err)
    }
    if _, ok := result["domain"]; !ok {
        t.Error("export JSON missing 'domain' key")
    }
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./cmd/am/domain/... -run TestDomainExport -v 2>&1 | tail -10
```
Expected: FAIL with "undefined: newExportCmd"

- [ ] **Step 3: Implement `cmd/am/domain/export.go`**

```go
package domain

import (
    "encoding/json"
    "fmt"
    "net/url"
    "os"
    "strconv"
    "sync"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

func newExportCmd(f *factory.Factory) *cobra.Command {
    var file string
    cmd := &cobra.Command{
        Use:     "export <domainId>",
        Short:   "Export domain configuration to JSON",
        Example: `  gio am domain export abc-123 -f domain-export.json`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireAMContext(f); err != nil {
                return err
            }
            return runExport(f, args[0], file)
        },
    }
    cmd.Flags().StringVarP(&file, "file", "f", "", "Output file path (default: stdout)")
    return cmd
}

func runExport(f *factory.Factory, domainID, file string) error {
    export, err := exportToMemory(f, domainID)
    if err != nil {
        return err
    }
    out, err := json.MarshalIndent(export, "", "  ")
    if err != nil {
        return err
    }
    if file != "" {
        return os.WriteFile(file, out, 0600)
    }
    fmt.Fprintln(f.IOStreams.Out, string(out))
    return nil
}

// exportToMemory fetches domain and all child resources concurrently.
func exportToMemory(f *factory.Factory, domainID string) (map[string]json.RawMessage, error) {
    domainPath := cmdutil.AMEnvPath(f, fmt.Sprintf("domains/%s", domainID))
    domainData, err := f.Client.Get(domainPath)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch domain: %w", err)
    }

    jobs := []struct {
        key string
        fn  func() (json.RawMessage, error)
    }{
        {"applications", func() (json.RawMessage, error) {
            return fetchAllPaginated(f, domainID, "applications")
        }},
        {"identityProviders", func() (json.RawMessage, error) {
            data, err := f.Client.Get(cmdutil.AMDomainPathFor(f, domainID, "identityProviders"))
            if err != nil {
                return nil, err
            }
            return json.RawMessage(data), nil
        }},
        {"roles", func() (json.RawMessage, error) {
            return fetchAllPaginated(f, domainID, "roles")
        }},
        {"scopes", func() (json.RawMessage, error) {
            return fetchAllPaginated(f, domainID, "scopes")
        }},
        {"factors", func() (json.RawMessage, error) {
            data, err := f.Client.Get(cmdutil.AMDomainPathFor(f, domainID, "factors"))
            if err != nil {
                return nil, err
            }
            return json.RawMessage(data), nil
        }},
        {"groups", func() (json.RawMessage, error) {
            return fetchAllPaginated(f, domainID, "groups")
        }},
        {"flows", func() (json.RawMessage, error) {
            data, err := f.Client.Get(cmdutil.AMDomainPathFor(f, domainID, "flows"))
            if err != nil {
                return nil, err
            }
            return json.RawMessage(data), nil
        }},
        {"certificates", func() (json.RawMessage, error) {
            data, err := f.Client.Get(cmdutil.AMDomainPathFor(f, domainID, "certificates"))
            if err != nil {
                return nil, err
            }
            return json.RawMessage(data), nil
        }},
    }

    results := make(map[string]json.RawMessage)
    var mu sync.Mutex
    var wg sync.WaitGroup
    var firstErr error

    for _, job := range jobs {
        wg.Add(1)
        go func(j struct {
            key string
            fn  func() (json.RawMessage, error)
        }) {
            defer wg.Done()
            data, err := j.fn()
            mu.Lock()
            defer mu.Unlock()
            if err != nil && firstErr == nil {
                firstErr = fmt.Errorf("fetch %s: %w", j.key, err)
                return
            }
            results[j.key] = data
        }(job)
    }
    wg.Wait()

    if firstErr != nil {
        return nil, firstErr
    }

    return map[string]json.RawMessage{
        "domain":            domainData,
        "applications":      results["applications"],
        "identityProviders": results["identityProviders"],
        "roles":             results["roles"],
        "scopes":            results["scopes"],
        "factors":           results["factors"],
        "groups":            results["groups"],
        "flows":             results["flows"],
        "certificates":      results["certificates"],
    }, nil
}

// fetchAllPaginated fetches all pages of a paginated resource and returns the combined data as a JSON array.
func fetchAllPaginated(f *factory.Factory, domainID, resource string) (json.RawMessage, error) {
    var all []json.RawMessage
    size := 100
    for page := 0; page <= 1000; page++ {
        q := url.Values{}
        q.Set("page", strconv.Itoa(page))
        q.Set("size", strconv.Itoa(size))
        path := cmdutil.AMDomainPathFor(f, domainID, resource+"?"+q.Encode())
        data, err := f.Client.Get(path)
        if err != nil {
            return nil, err
        }
        var resp amPaginatedResponse
        if err := json.Unmarshal(data, &resp); err != nil {
            return nil, fmt.Errorf("failed to parse %s response: %w", resource, err)
        }
        all = append(all, resp.Data...)
        if len(all) >= resp.TotalCount || len(resp.Data) < size {
            break
        }
    }
    b, err := json.Marshal(all)
    if err != nil {
        return nil, err
    }
    return json.RawMessage(b), nil
}
```

Note: `cmdutil.AMDomainPathFor` is a new helper needed for export (builds path with an explicit domainID rather than from config). We need to add it to `internal/cmdutil/cmdutil.go`:

```go
// AMDomainPathFor builds a domain-scoped path using an explicit domainID (not from config).
func AMDomainPathFor(f *factory.Factory, domainID, path string) string {
    return client.AMDomainPath(f.Resolved.Org, f.Resolved.Env, domainID, path)
}
```

- [ ] **Step 4: Add `AMDomainPathFor` to `internal/cmdutil/cmdutil.go`**

Add after the existing `AMDomainPath` function in `internal/cmdutil/cmdutil.go`:
```go
// AMDomainPathFor builds a domain-scoped path using an explicit domainID.
func AMDomainPathFor(f *factory.Factory, domainID, path string) string {
    return client.AMDomainPath(f.Resolved.Org, f.Resolved.Env, domainID, path)
}
```

- [ ] **Step 5: Implement `cmd/am/domain/import.go`**

```go
package domain

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

func newImportCmd(f *factory.Factory) *cobra.Command {
    var targetDomainID string
    cmd := &cobra.Command{
        Use:   "import <file>",
        Short: "Import domain configuration from a JSON export file",
        Example: `  gio am domain import domain-export.json
  gio am domain import domain-export.json --target existing-domain-id`,
        Args: cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireAMContext(f); err != nil {
                return err
            }
            if err := cmdutil.CheckReadOnly(f, "domain import"); err != nil {
                return err
            }
            return runImport(f, args[0], targetDomainID)
        },
    }
    cmd.Flags().StringVar(&targetDomainID, "target", "", "Target domain ID (creates new domain if not set)")
    return cmd
}

func runImport(f *factory.Factory, file, targetDomainID string) error {
    raw, err := os.ReadFile(file)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }

    var exportData map[string]json.RawMessage
    if err := json.Unmarshal(raw, &exportData); err != nil {
        return fmt.Errorf("failed to parse export file: %w", err)
    }

    p := cmdutil.NewPrinter(f)

    if targetDomainID == "" {
        var domainObj map[string]interface{}
        if err := json.Unmarshal(exportData["domain"], &domainObj); err != nil {
            return fmt.Errorf("failed to parse domain in export: %w", err)
        }
        body, err := json.Marshal(map[string]interface{}{
            "name":        cmdutil.StringField(domainObj, "name"),
            "description": cmdutil.StringField(domainObj, "description"),
        })
        if err != nil {
            return err
        }
        created, err := f.Client.Post(cmdutil.AMEnvPath(f, "domains"), body)
        if err != nil {
            return fmt.Errorf("failed to create domain: %w", err)
        }
        var newDomain map[string]interface{}
        if err := json.Unmarshal(created, &newDomain); err != nil {
            return fmt.Errorf("failed to parse CreateDomain response: %w", err)
        }
        targetDomainID = cmdutil.StringField(newDomain, "id")
        if targetDomainID == "" {
            return fmt.Errorf("CreateDomain response did not include an ID")
        }
        p.PrintMessage("Created domain '%s'.", targetDomainID)
    }

    imported, skipped := 0, 0
    add := func(i, s int) { imported += i; skipped += s }

    add(importItems(f, exportData, "scopes", targetDomainID, "scopes"))
    add(importItems(f, exportData, "roles", targetDomainID, "roles"))
    add(importItems(f, exportData, "groups", targetDomainID, "groups"))
    add(importItems(f, exportData, "applications", targetDomainID, "applications"))

    p.PrintMessage("Import complete: %d imported, %d skipped.", imported, skipped)
    return nil
}

// importItems creates resources from a JSON array in exportData under the given key.
// Returns (imported, skipped) counts.
func importItems(f *factory.Factory, exportData map[string]json.RawMessage, key, domainID, resource string) (int, int) {
    raw, ok := exportData[key]
    if !ok || len(raw) == 0 {
        return 0, 0
    }

    var items []json.RawMessage
    if err := json.Unmarshal(raw, &items); err != nil {
        return 0, 1
    }

    imported := 0
    for _, item := range items {
        path := cmdutil.AMDomainPathFor(f, domainID, resource)
        if _, err := f.Client.Post(path, item); err == nil {
            imported++
        }
    }
    return imported, 0
}
```

- [ ] **Step 6: Implement `cmd/am/domain/copy.go`**

```go
package domain

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/spf13/cobra"
)

func newCopyCmd(f *factory.Factory) *cobra.Command {
    var targetName string
    cmd := &cobra.Command{
        Use:     "copy <sourceDomainId>",
        Short:   "Copy a domain to a new domain in the same workspace",
        Example: `  gio am domain copy abc-123 --name my-copy`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireAMContext(f); err != nil {
                return err
            }
            if err := cmdutil.CheckReadOnly(f, "domain copy"); err != nil {
                return err
            }
            return runCopy(f, args[0], targetName)
        },
    }
    cmd.Flags().StringVar(&targetName, "name", "", "Name for the new domain (required)")
    _ = cmd.MarkFlagRequired("name")
    return cmd
}

func runCopy(f *factory.Factory, sourceDomainID, targetName string) error {
    p := cmdutil.NewPrinter(f)

    body, err := json.Marshal(map[string]interface{}{"name": targetName})
    if err != nil {
        return err
    }
    created, err := f.Client.Post(cmdutil.AMEnvPath(f, "domains"), body)
    if err != nil {
        return err
    }
    var newDomain map[string]interface{}
    if err := json.Unmarshal(created, &newDomain); err != nil {
        return fmt.Errorf("failed to parse CreateDomain response: %w", err)
    }
    targetDomainID := cmdutil.StringField(newDomain, "id")
    if targetDomainID == "" {
        return fmt.Errorf("CreateDomain response did not include an ID")
    }

    p.PrintMessage("Created domain '%s' (%s). Copying resources...", targetName, targetDomainID)

    exported, err := exportToMemory(f, sourceDomainID)
    if err != nil {
        return fmt.Errorf("failed to export source domain (new domain '%s' was created but is empty — delete it manually if not needed): %w", targetDomainID, err)
    }

    imported, skipped := 0, 0
    add := func(i, s int) { imported += i; skipped += s }

    add(importItems(f, exported, "scopes", targetDomainID, "scopes"))
    add(importItems(f, exported, "roles", targetDomainID, "roles"))
    add(importItems(f, exported, "groups", targetDomainID, "groups"))
    add(importItems(f, exported, "applications", targetDomainID, "applications"))

    p.PrintMessage("Copy complete: %d imported, %d skipped.", imported, skipped)
    return nil
}
```

- [ ] **Step 7: Wire export/import/copy into `cmd/am/domain/domain.go`**

The current `domain.go` has:
```go
cmd.AddCommand(newListCmd(f))
cmd.AddCommand(newGetCmd(f))
cmd.AddCommand(newCreateCmd(f))
cmd.AddCommand(newUpdateCmd(f))
cmd.AddCommand(newDeleteCmd(f))
cmd.AddCommand(newEnableCmd(f))
cmd.AddCommand(newDisableCmd(f))
```

Add after the last `AddCommand`:
```go
cmd.AddCommand(newExportCmd(f))
cmd.AddCommand(newImportCmd(f))
cmd.AddCommand(newCopyCmd(f))
```

- [ ] **Step 8: Run tests and full suite**

```bash
go test ./cmd/am/domain/... -run TestDomainExport -v 2>&1 | tail -10
go test ./... 2>&1 | tail -5
```
Expected: PASS

- [ ] **Step 9: Commit**

```bash
git add cmd/am/domain/export.go cmd/am/domain/import.go cmd/am/domain/copy.go \
        cmd/am/domain/export_test.go cmd/am/domain/domain.go \
        internal/cmdutil/cmdutil.go
git commit -m "feat: add gio am domain export/import/copy commands"
```

---

## Final Verification

- [ ] **Verify full build and tests pass**

```bash
go build ./... 2>&1 | head -20
go test ./... 2>&1
```
Expected: No errors, all tests pass.

- [ ] **Verify all new commands are reachable**

```bash
./gio am --help 2>&1 | grep -E "health|whoami|factor|flow|group|audit|token"
./gio am domain --help 2>&1 | grep -E "export|import|copy"
```
Expected: All 9 command groups appear in help output.

- [ ] **Final commit if any cleanup needed**

```bash
git add -A
git commit -m "chore: final cleanup and verification of AM missing features"
```

---

## Self-Review

### Spec Coverage
- ✅ `gio am health` — Task 1
- ✅ `gio am whoami` — Task 1
- ✅ `gio am factor list/get` — Task 2
- ✅ `gio am flow list/get` — Task 3
- ✅ `gio am group list/get/create/delete` — Task 4
- ✅ `gio am audit list/get` — Task 5
- ✅ `gio am token list/create/revoke` — Task 6
- ✅ `gio am domain export/import/copy` — Task 7

### Architecture Consistency
- All commands use `f.Client` directly (no service layer) — ✅ matches `am` branch pattern
- All domain-scoped commands use `RequireAMDomain(f)` — ✅
- `health`/`whoami`/`export`/`import`/`copy` use `RequireAMContext(f)` (domain from arg or not needed) — ✅
- `cmdutil.NewPrinter(f)` returns `*printer.Printer` (no error) — ✅
- `amPaginatedResponse` struct repeated per package (Go idiom, avoids cross-package coupling) — ✅

### Placeholder Scan
No TBD, TODO, or incomplete steps found. Every step has exact code.

### Type Consistency
- `cmdutil.AMDomainPathFor` added in Task 7 and used in `export.go`, `import.go`, `copy.go` — consistent
- `amPaginatedResponse` defined in each package that uses pagination (group, audit) — no cross-package conflicts
- `importItems` defined in `import.go`, referenced in `copy.go` — both in same package `domain` — ✅
