# AM Remaining Commands — Implementation Plan (Part 1: Tasks 1–7)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port the remaining 13 TypeScript commands from `am-tooling/am-cli` to Go `gio-cli` on the `am` branch.

**Architecture:** Direct HTTP via `f.Client.Get/Post/Delete`. URL helpers: `cmdutil.AMDomainPath(f, path)`, `cmdutil.AMEnvPath(f, path)`. Guards: `cmdutil.RequireAMContext(f)` / `cmdutil.RequireAMDomain(f)`. Tests: `newTestFactory(fake, readOnly)` in `helpers_test.go`.

**Tech Stack:** Go 1.22+, cobra, `internal/client`, `internal/cmdutil`, `internal/printer`, `internal/config`

**Source reference:** `/Users/rpo/Documents/Projects/Gravitee/AccessManagement/am-tooling/am-cli/src/commands/`

---

## File Structure

| File | Responsibility |
|------|---------------|
| `cmd/am/logout.go` | Clear token from current or all contexts |
| `cmd/am/status.go` | Show context/session info (no network) |
| `cmd/am/doctor.go` | Diagnostic checks with connectivity test |
| `cmd/am/logs/logs.go` | Parent `logs` command |
| `cmd/am/logs/list.go` | Poll audit logs, `--follow` mode |
| `cmd/am/plugin/plugin.go` | Parent `plugin` command |
| `cmd/am/plugin/list.go` | List platform plugins by type |
| `cmd/am/plugin/schema.go` | Show plugin configuration schema |
| `cmd/am/plugin/create.go` | Create resource from plugin (non-interactive via `--config-file`) |
| `cmd/am/diff/diff.go` | Compare two domain configs across contexts |
| `cmd/am/diff/compare.go` | Pure comparison logic (testable) |
| `cmd/am/lint/lint.go` | Parent + run all 14 security rules |
| `cmd/am/lint/rules.go` | 14 rule implementations (pure functions) |
| `cmd/am/watch/watch.go` | Live dashboard polling audits |
| `cmd/am/watch/render.go` | Pure render logic (testable) |
| `cmd/am/shell/shell.go` | Interactive REPL |
| `cmd/am/oidctest/oidctest.go` | Parent `test` command |
| `cmd/am/oidctest/discover.go` | OIDC discovery |
| `cmd/am/oidctest/login.go` | ROPC flow test |
| `cmd/am/oidctest/clientcreds.go` | client_credentials flow test |
| `cmd/am/trace/trace.go` | Auth path trace command |
| `cmd/am/trace/checks.go` | 7 pure check functions |
| `cmd/am/supportdump/supportdump.go` | Diagnostic dump command |
| `cmd/am/supportdump/redact.go` | Secret redaction (pure) |
| `cmd/am/completion.go` | Shell completion wrapper |
| `cmd/am/am.go` | Wire all new commands (modify existing) |

---

## Task 1: `gio am logout`

**Files:**
- Create: `cmd/am/logout.go`
- Modify: `cmd/am/am.go`

### TypeScript behaviour (logout.ts)
- `--all` flag: delete `token` from every context in `config.Contexts`
- Default: delete token from `config.Contexts[config.CurrentContext]`
- Save config after mutation
- Print success or warning if no token stored

- [ ] **Step 1: Write the failing test** — add to `cmd/am/health_test.go` (same `am` package):

```go
func TestLogout(t *testing.T) {
    cfg := &config.Config{
        Contexts:       map[string]config.Context{"ctx1": {Token: "tok"}},
        CurrentContext: "ctx1",
    }
    f, out := newAMTestFactory(nil, cfg)
    cmd := newLogoutCmd(f)
    cmd.SetArgs([]string{})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if f.Config.Contexts["ctx1"].Token != "" {
        t.Error("expected token to be cleared")
    }
    if !strings.Contains(out.String(), "Logged out") {
        t.Errorf("expected success message, got: %s", out.String())
    }
}

func TestLogoutAll(t *testing.T) {
    cfg := &config.Config{
        Contexts: map[string]config.Context{
            "ctx1": {Token: "tok1"},
            "ctx2": {Token: "tok2"},
        },
        CurrentContext: "ctx1",
    }
    f, out := newAMTestFactory(nil, cfg)
    cmd := newLogoutCmd(f)
    cmd.SetArgs([]string{"--all"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    for name, ctx := range f.Config.Contexts {
        if ctx.Token != "" {
            t.Errorf("expected token cleared for %s", name)
        }
    }
    if !strings.Contains(out.String(), "2") {
        t.Errorf("expected count in message, got: %s", out.String())
    }
}
```

Note: `newAMTestFactory` already exists in `cmd/am/health_test.go` but only takes a client. We need to extend it or add a variant. Check `cmd/am/health_test.go` for the current signature. If it does not accept a `*config.Config`, add `newAMTestFactoryWithConfig(client, cfg)`.

- [ ] **Step 2: Run test — verify it fails**

```
cd /Users/rpo/Documents/Projects/Gravitee/gio-cli && go test ./cmd/am/ -run TestLogout -v
```
Expected: `FAIL — newLogoutCmd undefined`

- [ ] **Step 3: Implement `cmd/am/logout.go`**

```go
package am

import (
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newLogoutCmd(f *factory.Factory) *cobra.Command {
	var all bool
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Clear stored authentication token",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := f.Config
			if all {
				count := 0
				for name, ctx := range cfg.Contexts {
					if ctx.Token != "" {
						ctx.Token = ""
						cfg.Contexts[name] = ctx
						count++
					}
				}
				if count == 0 {
					fmt.Fprintln(f.IOStreams.Out, "No stored tokens to clear.")
					return nil
				}
				if err := cfg.SaveTo(f.ConfigPath); err != nil {
					return fmt.Errorf("failed to save config: %w", err)
				}
				fmt.Fprintf(f.IOStreams.Out, "Cleared tokens for %d workspace(s).\n", count)
				return nil
			}
			if cfg.CurrentContext == "" {
				fmt.Fprintln(f.IOStreams.Out, "No workspace selected.")
				return nil
			}
			ctx, ok := cfg.Contexts[cfg.CurrentContext]
			if !ok || ctx.Token == "" {
				fmt.Fprintf(f.IOStreams.Out, "No token stored for workspace %q.\n", cfg.CurrentContext)
				return nil
			}
			ctx.Token = ""
			cfg.Contexts[cfg.CurrentContext] = ctx
			if err := cfg.SaveTo(f.ConfigPath); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			fmt.Fprintf(f.IOStreams.Out, "Logged out from workspace %q.\n", cfg.CurrentContext)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&all, "all", "a", false, "Clear tokens for all workspaces")
	return cmd
}
```

- [ ] **Step 4: Wire in `cmd/am/am.go`** — add `cmd.AddCommand(newLogoutCmd(f))` after the existing `cmd.AddCommand(newLoginCmd(f))` line.

- [ ] **Step 5: Run test — verify pass**

```
go test ./cmd/am/ -run TestLogout -v
```
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add cmd/am/logout.go cmd/am/am.go cmd/am/health_test.go
git commit -m "feat: add gio am logout command"
```

---

## Task 2: `gio am status`

**Files:**
- Create: `cmd/am/status.go`
- Modify: `cmd/am/am.go`, `cmd/am/health_test.go`

### Behaviour (status.ts)
- No network call
- Show: workspace, org, env, domain, authenticated (yes/no/expired), CLI version
- Read from `f.Config` + `f.Resolved` (Resolved may be nil)

- [ ] **Step 1: Write the failing test** — add to `cmd/am/health_test.go`:

```go
func TestStatus(t *testing.T) {
    cfg := &config.Config{
        Contexts: map[string]config.Context{
            "myws": {URL: "https://am.example.com", Token: "tok", Org: "DEFAULT", Env: "DEFAULT"},
        },
        CurrentContext: "myws",
    }
    f, out := newAMTestFactoryWithConfig(nil, cfg)
    cmd := newStatusCmd(f)
    cmd.SetArgs([]string{})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "myws") {
        t.Errorf("expected workspace name, got: %s", out.String())
    }
    if !strings.Contains(out.String(), "https://am.example.com") {
        t.Errorf("expected URL, got: %s", out.String())
    }
}
```

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/ -run TestStatus -v
```

- [ ] **Step 3: Implement `cmd/am/status.go`**

```go
package am

import (
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newStatusCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current CLI context and session status",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := f.Config
			out := f.IOStreams.Out

			if cfg == nil || cfg.CurrentContext == "" {
				fmt.Fprintln(out, "workspace:     (not set)")
				fmt.Fprintln(out, "authenticated: no")
				return nil
			}

			ctx, ok := cfg.Contexts[cfg.CurrentContext]
			domain := ""
			if f.Resolved != nil {
				domain = f.Resolved.Domain
			}

			fmt.Fprintf(out, "workspace:     %s @ %s\n", cfg.CurrentContext, ctx.URL)
			fmt.Fprintf(out, "organization:  %s\n", ctx.Org)
			fmt.Fprintf(out, "environment:   %s\n", ctx.Env)
			if domain != "" {
				fmt.Fprintf(out, "domain:        %s\n", domain)
			} else {
				fmt.Fprintln(out, "domain:        (not set)")
			}
			if !ok || ctx.Token == "" {
				fmt.Fprintln(out, "authenticated: no")
			} else {
				fmt.Fprintln(out, "authenticated: yes")
			}
			return nil
		},
	}
}
```

- [ ] **Step 4: Wire in `am.go`** — `cmd.AddCommand(newStatusCmd(f))`

- [ ] **Step 5: Run — verify pass**

```
go test ./cmd/am/ -run TestStatus -v
```

- [ ] **Step 6: Commit**

```bash
git add cmd/am/status.go cmd/am/am.go cmd/am/health_test.go
git commit -m "feat: add gio am status command"
```

---

## Task 3: `gio am doctor`

**Files:**
- Create: `cmd/am/doctor.go`
- Modify: `cmd/am/am.go`, `cmd/am/health_test.go`

### Behaviour (doctor.ts)
- Check 1: config file exists and has contexts
- Check 2: current context is set
- Check 3: token is present
- Check 4: domain is set
- Check 5: connectivity — GET `/management/user` (requires AM context + token)
- Output: table of checks with OK/WARN/FAIL

- [ ] **Step 1: Write the failing test** — add to `cmd/am/health_test.go`:

```go
func TestDoctor(t *testing.T) {
    fake := &client.FakeClient{
        GetFunc: func(path string) ([]byte, error) {
            if strings.Contains(path, "/management/user") {
                return []byte(`{"id":"u1","username":"admin"}`), nil
            }
            return nil, fmt.Errorf("unexpected path: %s", path)
        },
    }
    f, out := newAMTestFactory(fake, false)
    cmd := newDoctorCmd(f)
    cmd.SetArgs([]string{})
    _ = cmd.Execute() // doctor always exits 0 in success
    if !strings.Contains(out.String(), "OK") {
        t.Errorf("expected OK checks, got: %s", out.String())
    }
}
```

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/ -run TestDoctor -v
```

- [ ] **Step 3: Implement `cmd/am/doctor.go`**

```go
package am

import (
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

type doctorCheck struct {
	label  string
	status string // "OK", "WARN", "FAIL"
	detail string
}

func newDoctorCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run diagnostic checks on the CLI configuration",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			checks := runDoctorChecks(f)
			out := f.IOStreams.Out
			for _, c := range checks {
				fmt.Fprintf(out, "  [%-4s] %-20s %s\n", c.status, c.label, c.detail)
			}
			return nil
		},
	}
}

func runDoctorChecks(f *factory.Factory) []doctorCheck {
	var checks []doctorCheck

	// 1. Config
	if f.Config == nil || len(f.Config.Contexts) == 0 {
		checks = append(checks, doctorCheck{"config", "FAIL", "No contexts configured — run 'gio am login'"})
		return checks
	}
	checks = append(checks, doctorCheck{"config", "OK", fmt.Sprintf("%d context(s) found", len(f.Config.Contexts))})

	// 2. Current context
	if f.Config.CurrentContext == "" {
		checks = append(checks, doctorCheck{"context", "WARN", "No current context set"})
		return checks
	}
	ctx, ok := f.Config.Contexts[f.Config.CurrentContext]
	if !ok {
		checks = append(checks, doctorCheck{"context", "FAIL", fmt.Sprintf("Context %q not found", f.Config.CurrentContext)})
		return checks
	}
	checks = append(checks, doctorCheck{"context", "OK", fmt.Sprintf("%s @ %s", f.Config.CurrentContext, ctx.URL)})

	// 3. Token
	if ctx.Token == "" {
		checks = append(checks, doctorCheck{"auth", "FAIL", "No token stored — run 'gio am login'"})
		return checks
	}
	checks = append(checks, doctorCheck{"auth", "OK", "Token present"})

	// 4. Domain
	domain := ""
	if f.Resolved != nil {
		domain = f.Resolved.Domain
	}
	if domain == "" {
		checks = append(checks, doctorCheck{"domain", "WARN", "No domain set — run 'gio am set domain <id>'"})
	} else {
		checks = append(checks, doctorCheck{"domain", "OK", domain})
	}

	// 5. Connectivity
	if err := cmdutil.RequireAMContext(f); err == nil && f.Client != nil {
		_, err := f.Client.Get("/management/user")
		if err != nil {
			checks = append(checks, doctorCheck{"connect", "FAIL", fmt.Sprintf("Cannot reach AM: %v", err)})
		} else {
			checks = append(checks, doctorCheck{"connect", "OK", "AM management API reachable"})
		}
	}

	return checks
}
```

- [ ] **Step 4: Wire in `am.go`** — `cmd.AddCommand(newDoctorCmd(f))`

- [ ] **Step 5: Run — verify pass**

```
go test ./cmd/am/ -run TestDoctor -v
```

- [ ] **Step 6: Commit**

```bash
git add cmd/am/doctor.go cmd/am/am.go cmd/am/health_test.go
git commit -m "feat: add gio am doctor command"
```

---

## Task 4: `gio am logs`

**Files:**
- Create: `cmd/am/logs/logs.go`, `cmd/am/logs/list.go`, `cmd/am/logs/helpers_test.go`, `cmd/am/logs/logs_test.go`
- Modify: `cmd/am/am.go`

### Behaviour (logs.ts)
- `--lines N` (default 20): show last N audit events
- `--follow`: poll every `--interval` seconds, track seen IDs, print new events
- `--interval N` (default 5): polling interval seconds
- `--type`, `--status`: filter (same as `gio am audit list`)
- Output format per event: `timestamp  STATUS  TYPE  actor → target`

- [ ] **Step 1: Write the failing test**

Create `cmd/am/logs/helpers_test.go`:

```go
package logs

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts:       map[string]config.Context{"test": {URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT"}},
		CurrentContext: "test",
	}
	f := &factory.Factory{
		Config:   cfg,
		Resolved: &config.ResolvedContext{Name: "test", URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am", ReadOnly: readOnly},
		Client:   c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
```

Create `cmd/am/logs/logs_test.go`:

```go
package logs

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestLogsListBasic(t *testing.T) {
	resp := map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"id": "e1", "type": "USER_LOGIN", "timestamp": 1700000000000,
				"outcome": map[string]interface{}{"status": "SUCCESS"},
				"actor":   map[string]interface{}{"displayName": "admin"},
				"target":  map[string]interface{}{"displayName": "alice"},
			},
		},
		"currentPage": 0, "totalCount": 1,
	}
	data, _ := json.Marshal(resp)
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) { return data, nil },
	}
	f, out := newTestFactory(fake, false)
	cmd := NewLogsCmd(f)
	cmd.SetArgs([]string{"--lines", "5"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "USER_LOGIN") {
		t.Errorf("expected USER_LOGIN in output, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "admin") {
		t.Errorf("expected actor in output, got: %s", out.String())
	}
}

func TestFormatEvent(t *testing.T) {
	e := auditEvent{
		ID:        "e1",
		EventType: "USER_LOGIN",
		Status:    "SUCCESS",
		Actor:     "admin",
		Target:    "alice",
		Timestamp: "2023-11-14 22:13:20",
	}
	line := formatEvent(e)
	if !strings.Contains(line, "USER_LOGIN") {
		t.Errorf("expected type, got: %s", line)
	}
	if !strings.Contains(line, "admin") {
		t.Errorf("expected actor, got: %s", line)
	}
	if !strings.Contains(line, "alice") {
		t.Errorf("expected target, got: %s", line)
	}
}
```

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/logs/ -v
```
Expected: `FAIL — NewLogsCmd undefined`

- [ ] **Step 3: Implement `cmd/am/logs/logs.go`**

```go
package logs

import (
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewLogsCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Tail audit logs (alias: gio am audit list --follow)",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newListCmd(f))
	// make list the default action
	cmd.RunE = newListCmd(f).RunE
	return cmd
}
```

- [ ] **Step 4: Implement `cmd/am/logs/list.go`**

```go
package logs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

type auditEvent struct {
	ID        string
	EventType string
	Status    string
	Actor     string
	Target    string
	Timestamp string
}

type paginatedResp struct {
	Data       []json.RawMessage `json:"data"`
	TotalCount int               `json:"totalCount"`
}

func newListCmd(f *factory.Factory) *cobra.Command {
	var lines int
	var follow bool
	var interval int
	var filterType string
	var filterStatus string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show recent audit log events",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			if follow {
				return runFollow(f, interval, filterType, filterStatus)
			}
			return runOnce(f, lines, filterType, filterStatus)
		},
	}
	cmd.Flags().IntVarP(&lines, "lines", "n", 20, "Number of recent events to show")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Poll for new events continuously")
	cmd.Flags().IntVar(&interval, "interval", 5, "Polling interval in seconds (used with --follow)")
	cmd.Flags().StringVar(&filterType, "type", "", "Filter by audit type")
	cmd.Flags().StringVar(&filterStatus, "status", "", "Filter by status (SUCCESS, FAILURE)")
	return cmd
}

func buildQuery(page, size int, filterType, filterStatus string) string {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("size", strconv.Itoa(size))
	if filterType != "" {
		q.Set("type", filterType)
	}
	if filterStatus != "" {
		q.Set("status", filterStatus)
	}
	return q.Encode()
}

func fetchEvents(f *factory.Factory, size int, filterType, filterStatus string) ([]auditEvent, error) {
	path := cmdutil.AMDomainPath(f, "audits?"+buildQuery(0, size, filterType, filterStatus))
	data, err := f.Client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp paginatedResp
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	events := make([]auditEvent, 0, len(resp.Data))
	for _, raw := range resp.Data {
		events = append(events, parseEvent(raw))
	}
	return events, nil
}

func parseEvent(raw json.RawMessage) auditEvent {
	var m map[string]interface{}
	_ = json.Unmarshal(raw, &m)
	e := auditEvent{
		ID:        stringField(m, "id"),
		EventType: stringField(m, "type"),
	}
	if outcome, ok := m["outcome"].(map[string]interface{}); ok {
		e.Status = stringField(outcome, "status")
	}
	if actor, ok := m["actor"].(map[string]interface{}); ok {
		e.Actor = stringField(actor, "displayName")
		if e.Actor == "" {
			e.Actor = stringField(actor, "id")
		}
	}
	if target, ok := m["target"].(map[string]interface{}); ok {
		e.Target = stringField(target, "displayName")
		if e.Target == "" {
			e.Target = stringField(target, "id")
		}
	}
	if ts, ok := m["timestamp"].(float64); ok {
		e.Timestamp = time.UnixMilli(int64(ts)).UTC().Format("2006-01-02 15:04:05")
	}
	return e
}

func formatEvent(e auditEvent) string {
	target := ""
	if e.Target != "" {
		target = " → " + e.Target
	}
	return fmt.Sprintf("%s  %-8s  %-30s  %s%s", e.Timestamp, e.Status, e.EventType, e.Actor, target)
}

func runOnce(f *factory.Factory, lines int, filterType, filterStatus string) error {
	events, err := fetchEvents(f, lines, filterType, filterStatus)
	if err != nil {
		return err
	}
	for _, e := range events {
		fmt.Fprintln(f.IOStreams.Out, formatEvent(e))
	}
	return nil
}

func runFollow(f *factory.Factory, intervalSec int, filterType, filterStatus string) error {
	seen := make(map[string]bool)
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	// Initial fetch
	events, err := fetchEvents(f, 50, filterType, filterStatus)
	if err != nil {
		return err
	}
	for _, e := range events {
		seen[e.ID] = true
		fmt.Fprintln(f.IOStreams.Out, formatEvent(e))
	}

	for range ticker.C {
		events, err := fetchEvents(f, 50, filterType, filterStatus)
		if err != nil {
			fmt.Fprintf(f.IOStreams.Out, "error: %v\n", err)
			continue
		}
		for _, e := range events {
			if !seen[e.ID] {
				seen[e.ID] = true
				fmt.Fprintln(f.IOStreams.Out, formatEvent(e))
			}
		}
	}
	return nil
}

func stringField(m map[string]interface{}, key string) string {
	s, _ := m[key].(string)
	return strings.TrimSpace(s)
}
```

- [ ] **Step 5: Wire in `am.go`** — add import and `cmd.AddCommand(logscmd.NewLogsCmd(f))`

- [ ] **Step 6: Run — verify pass**

```
go test ./cmd/am/logs/ -v
```

- [ ] **Step 7: Commit**

```bash
git add cmd/am/logs/ cmd/am/am.go
git commit -m "feat: add gio am logs command with --follow support"
```

---

## Task 5: `gio am plugin list/schema/create`

**Files:**
- Create: `cmd/am/plugin/plugin.go`, `cmd/am/plugin/list.go`, `cmd/am/plugin/schema.go`, `cmd/am/plugin/create.go`, `cmd/am/plugin/helpers_test.go`, `cmd/am/plugin/plugin_test.go`
- Modify: `cmd/am/am.go`

### Behaviour (plugin.ts)
- `plugin list <type>` — GET `/management/platform/plugins/{apiType}` (no domain needed, uses AMEnvPath base)
- `plugin schema <type> <pluginId>` — GET `/management/platform/plugins/{apiType}/{pluginId}/schema`
- `plugin create <type> <pluginId>` — non-interactive via `--config-file <json>`, `--name <name>`, POST to domain
- Valid types: `idp`, `factor`, `certificate`, `policy`, `resource`, `reporter`, `botdetection`

Platform plugins URL: `{baseURL}/management/platform/plugins/{apiType}` — this is NOT org/env scoped. Use direct path building: `"/management/platform/plugins/" + apiType`.

- [ ] **Step 1: Write the failing test**

Create `cmd/am/plugin/helpers_test.go`:

```go
package plugin

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts:       map[string]config.Context{"test": {URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT"}},
		CurrentContext: "test",
	}
	f := &factory.Factory{
		Config:   cfg,
		Resolved: &config.ResolvedContext{Name: "test", URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am", ReadOnly: readOnly},
		Client:   c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
```

Create `cmd/am/plugin/plugin_test.go`:

```go
package plugin

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gravitee-io/gio-cli/internal/client"
)

func TestPluginList(t *testing.T) {
	plugins := []map[string]interface{}{
		{"id": "github-am-idp", "name": "GitHub Identity Provider", "version": "2.4.0"},
	}
	data, _ := json.Marshal(plugins)
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/platform/plugins/identities") {
				t.Errorf("unexpected path: %s", path)
			}
			return data, nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := NewPluginCmd(f)
	cmd.SetArgs([]string{"list", "idp"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "GitHub Identity Provider") {
		t.Errorf("expected plugin name, got: %s", out.String())
	}
}

func TestPluginSchema(t *testing.T) {
	schema := map[string]interface{}{
		"properties": map[string]interface{}{
			"clientId":     map[string]interface{}{"type": "string", "title": "Client ID"},
			"clientSecret": map[string]interface{}{"type": "string", "title": "Client Secret"},
		},
	}
	data, _ := json.Marshal(schema)
	fake := &client.FakeClient{
		GetFunc: func(path string) ([]byte, error) {
			if !strings.Contains(path, "/platform/plugins/identities/github-am-idp/schema") {
				t.Errorf("unexpected path: %s", path)
			}
			return data, nil
		},
	}
	f, out := newTestFactory(fake, false)
	cmd := NewPluginCmd(f)
	cmd.SetArgs([]string{"schema", "idp", "github-am-idp"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "clientId") {
		t.Errorf("expected field in output, got: %s", out.String())
	}
}

func TestPluginCreateWithFile(t *testing.T) {
	var posted map[string]interface{}
	fake := &client.FakeClient{
		PostFunc: func(path string, body interface{}) ([]byte, error) {
			if b, ok := body.([]byte); ok {
				_ = json.Unmarshal(b, &posted)
			}
			return []byte(`{"id":"idp-new","name":"My GitHub"}`), nil
		},
	}
	f, out := newTestFactory(fake, false)

	// write temp config file
	import os
	tmp, _ := os.CreateTemp("", "*.json")
	tmp.WriteString(`{"clientId":"abc","clientSecret":"xyz"}`)
	tmp.Close()
	defer os.Remove(tmp.Name())

	cmd := NewPluginCmd(f)
	cmd.SetArgs([]string{"create", "idp", "github-am-idp", "--name", "My GitHub", "--config-file", tmp.Name()})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "My GitHub") {
		t.Errorf("expected name in output, got: %s", out.String())
	}
}
```

Note: the `import os` inside TestPluginCreateWithFile is incorrect Go syntax — move it to the package-level imports block.

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/plugin/ -v
```

- [ ] **Step 3: Implement `cmd/am/plugin/plugin.go`**

```go
package plugin

import (
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

var pluginTypes = map[string]string{
	"idp":         "identities",
	"factor":      "factors",
	"certificate": "certificates",
	"policy":      "policies",
	"resource":    "resources",
	"reporter":    "reporters",
	"botdetection": "bot-detections",
}

func NewPluginCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Explore and create resources from plugin schemas",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newSchemaCmd(f))
	cmd.AddCommand(newCreateCmd(f))
	return cmd
}
```

- [ ] **Step 4: Implement `cmd/am/plugin/list.go`**

```go
package plugin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list <type>",
		Short: "List available plugins of a given type",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			ptype := args[0]
			apiType, ok := pluginTypes[ptype]
			if !ok {
				return fmt.Errorf("unknown plugin type %q. Available: %s", ptype, strings.Join(keys(pluginTypes), ", "))
			}
			data, err := f.Client.Get("/management/platform/plugins/" + apiType)
			if err != nil {
				return err
			}
			p := cmdutil.NewPrinter(f)
			var items []json.RawMessage
			if err := json.Unmarshal(data, &items); err != nil {
				return p.PrintDetail(json.RawMessage(data))
			}
			return p.PrintList(items, pluginColumns())
		},
	}
}

func pluginColumns() []cmdutil_printer_columns {
	// inline to avoid import cycle — use printer.Column directly
	return nil // defined below
}
```

Actually, use the same pattern as other commands:

```go
package plugin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/gravitee-io/gio-cli/internal/printer"
	"github.com/spf13/cobra"
)

func newListCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list <type>",
		Short: "List available plugins of a given type",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			ptype := args[0]
			apiType, ok := pluginTypes[ptype]
			if !ok {
				return fmt.Errorf("unknown plugin type %q. Available: %s", ptype, strings.Join(mapKeys(pluginTypes), ", "))
			}
			data, err := f.Client.Get("/management/platform/plugins/" + apiType)
			if err != nil {
				return err
			}
			p := cmdutil.NewPrinter(f)
			var items []json.RawMessage
			if err := json.Unmarshal(data, &items); err != nil {
				return p.PrintDetail(json.RawMessage(data))
			}
			return p.PrintList(items, []printer.Column{
				{Name: "ID", Value: func(i interface{}) string { return cmdutil.StringField(i, "id") }},
				{Name: "Name", Value: func(i interface{}) string { return cmdutil.StringField(i, "name") }},
				{Name: "Version", Value: func(i interface{}) string { return cmdutil.StringField(i, "version") }},
			})
		},
	}
}

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
```

- [ ] **Step 5: Implement `cmd/am/plugin/schema.go`**

```go
package plugin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newSchemaCmd(f *factory.Factory) *cobra.Command {
	var raw bool
	cmd := &cobra.Command{
		Use:   "schema <type> <pluginId>",
		Short: "Show configuration schema for a plugin",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			ptype, pluginID := args[0], args[1]
			apiType, ok := pluginTypes[ptype]
			if !ok {
				return fmt.Errorf("unknown plugin type %q", ptype)
			}
			path := fmt.Sprintf("/management/platform/plugins/%s/%s/schema", apiType, pluginID)
			data, err := f.Client.Get(path)
			if err != nil {
				return err
			}
			if raw {
				fmt.Fprintln(f.IOStreams.Out, string(data))
				return nil
			}
			var schema map[string]interface{}
			if err := json.Unmarshal(data, &schema); err != nil {
				fmt.Fprintln(f.IOStreams.Out, string(data))
				return nil
			}
			printSchema(f, schema, pluginID)
			return nil
		},
	}
	cmd.Flags().BoolVar(&raw, "raw", false, "Show raw JSON schema")
	return cmd
}

func printSchema(f *factory.Factory, schema map[string]interface{}, pluginID string) {
	fmt.Fprintf(f.IOStreams.Out, "Schema: %s\n\n", pluginID)
	props, _ := schema["properties"].(map[string]interface{})
	for key, val := range props {
		prop, _ := val.(map[string]interface{})
		title, _ := prop["title"].(string)
		propType, _ := prop["type"].(string)
		desc, _ := prop["description"].(string)
		if title == "" {
			title = key
		}
		line := fmt.Sprintf("  %-30s %-10s", key, propType)
		if desc != "" {
			line += "  " + desc
		} else if title != key {
			line += "  " + title
		}
		fmt.Fprintln(f.IOStreams.Out, line)
	}
	// Show enum values
	for key, val := range props {
		prop, _ := val.(map[string]interface{})
		if enums, ok := prop["enum"].([]interface{}); ok {
			strs := make([]string, 0, len(enums))
			for _, e := range enums {
				strs = append(strs, fmt.Sprintf("%v", e))
			}
			fmt.Fprintf(f.IOStreams.Out, "  %s: {%s}\n", key, strings.Join(strs, ", "))
		}
	}
}
```

- [ ] **Step 6: Implement `cmd/am/plugin/create.go`**

```go
package plugin

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
	var name string
	var configFile string

	cmd := &cobra.Command{
		Use:   "create <type> <pluginId>",
		Short: "Create a resource instance from a plugin (use --config-file for non-interactive)",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			if err := cmdutil.CheckReadOnly(f, "plugin create"); err != nil {
				return err
			}
			ptype, pluginID := args[0], args[1]
			apiPath, ok := pluginTypes[ptype]
			if !ok {
				return fmt.Errorf("unknown plugin type %q", ptype)
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if configFile == "" {
				return fmt.Errorf("--config-file is required (interactive mode not supported in gio CLI)")
			}
			raw, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}
			body := map[string]interface{}{
				"name":          name,
				"type":          pluginID,
				"configuration": string(raw),
			}
			bodyJSON, _ := json.Marshal(body)
			path := cmdutil.AMDomainPath(f, apiPath)
			data, err := f.Client.Post(path, bodyJSON)
			if err != nil {
				return err
			}
			p := cmdutil.NewPrinter(f)
			return p.PrintDetail(json.RawMessage(data))
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Resource name (required)")
	cmd.Flags().StringVarP(&configFile, "config-file", "f", "", "JSON config file (required)")
	return cmd
}
```

- [ ] **Step 7: Wire in `am.go`** — add `cmd.AddCommand(plugincmd.NewPluginCmd(f))`

- [ ] **Step 8: Run — verify pass**

```
go test ./cmd/am/plugin/ -v
```

- [ ] **Step 9: Commit**

```bash
git add cmd/am/plugin/ cmd/am/am.go
git commit -m "feat: add gio am plugin list/schema/create commands"
```

---

## Task 6: `gio am diff`

**Files:**
- Create: `cmd/am/diff/diff.go`, `cmd/am/diff/compare.go`, `cmd/am/diff/helpers_test.go`, `cmd/am/diff/diff_test.go`
- Modify: `cmd/am/am.go`

### Behaviour (diff.ts)
- `--from <context>` and `--to <context>` — named contexts in `~/.gio/config.json`
- Fetch from each: scopes, roles, groups, apps, idps, certificates, factors, flows
- For each resource type: compare by name/clientId (key field), show +added, -removed, ~changed
- Changed: compare specific fields per type (see resourceSpecs below)
- Output: table with `+/-/~` prefix

Resource key fields and compare fields:
```
scopes:       key=key,         compare=[name, description]
roles:        key=name,        compare=[description, assignableType]
groups:       key=name,        compare=[description]
apps:         key=name,        compare=[description, type]
idps:         key=name,        compare=[type]
certificates: key=name,        compare=[type]
factors:      key=name,        compare=[factorType]
flows:        key=type,        compare=[enabled, pre, post]
```

- [ ] **Step 1: Write the failing test**

Create `cmd/am/diff/helpers_test.go`:

```go
package diff

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts: map[string]config.Context{
			"ctx-a": {URL: "http://am-a", Token: "tok-a", Org: "DEFAULT", Env: "DEFAULT"},
			"ctx-b": {URL: "http://am-b", Token: "tok-b", Org: "DEFAULT", Env: "DEFAULT"},
		},
		CurrentContext: "ctx-a",
	}
	f := &factory.Factory{
		Config:   cfg,
		Resolved: &config.ResolvedContext{Name: "ctx-a", URL: "http://am-a", Token: "tok-a", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am"},
		Client:   c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
```

Create `cmd/am/diff/diff_test.go`:

```go
package diff

import (
	"strings"
	"testing"
)

func TestDiffObjects(t *testing.T) {
	from := map[string]interface{}{"name": "foo", "description": "old desc"}
	to := map[string]interface{}{"name": "foo", "description": "new desc"}
	changes := diffObjects(from, to, []string{"description"})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Field != "description" {
		t.Errorf("expected field 'description', got %q", changes[0].Field)
	}
}

func TestCompareResourcesAdded(t *testing.T) {
	fromItems := []map[string]interface{}{
		{"name": "scope-a"},
	}
	toItems := []map[string]interface{}{
		{"name": "scope-a"},
		{"name": "scope-b"},
	}
	result := compareResources(fromItems, toItems, "name", []string{"name"})
	if result.Added != 1 {
		t.Errorf("expected 1 added, got %d", result.Added)
	}
}

func TestCompareResourcesRemoved(t *testing.T) {
	fromItems := []map[string]interface{}{
		{"name": "scope-a"},
		{"name": "scope-b"},
	}
	toItems := []map[string]interface{}{
		{"name": "scope-a"},
	}
	result := compareResources(fromItems, toItems, "name", []string{"name"})
	if result.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", result.Removed)
	}
}

func TestFormatDiffLine(t *testing.T) {
	line := formatDiffLine("+", "scope", "scope-b", nil)
	if !strings.Contains(line, "+") {
		t.Error("expected + prefix")
	}
	if !strings.Contains(line, "scope-b") {
		t.Error("expected resource name")
	}
}
```

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/diff/ -v
```

- [ ] **Step 3: Implement `cmd/am/diff/compare.go`**

```go
package diff

import "fmt"

type FieldChange struct {
	Field string
	From  interface{}
	To    interface{}
}

type DiffResult struct {
	Added   int
	Removed int
	Changed int
	Lines   []string
}

func diffObjects(from, to map[string]interface{}, compareFields []string) []FieldChange {
	var changes []FieldChange
	for _, field := range compareFields {
		fromVal := fmt.Sprintf("%v", from[field])
		toVal := fmt.Sprintf("%v", to[field])
		if fromVal != toVal {
			changes = append(changes, FieldChange{Field: field, From: from[field], To: to[field]})
		}
	}
	return changes
}

func compareResources(fromItems, toItems []map[string]interface{}, keyField string, compareFields []string) DiffResult {
	result := DiffResult{}

	fromMap := make(map[string]map[string]interface{}, len(fromItems))
	for _, item := range fromItems {
		key := fmt.Sprintf("%v", item[keyField])
		fromMap[key] = item
	}

	toMap := make(map[string]map[string]interface{}, len(toItems))
	for _, item := range toItems {
		key := fmt.Sprintf("%v", item[keyField])
		toMap[key] = item
	}

	// Added
	for key, item := range toMap {
		if _, exists := fromMap[key]; !exists {
			result.Added++
			result.Lines = append(result.Lines, formatDiffLine("+", keyField, key, item))
		}
	}

	// Removed
	for key, item := range fromMap {
		if _, exists := toMap[key]; !exists {
			result.Removed++
			result.Lines = append(result.Lines, formatDiffLine("-", keyField, key, item))
		}
	}

	// Changed
	for key, toItem := range toMap {
		fromItem, exists := fromMap[key]
		if !exists {
			continue
		}
		changes := diffObjects(fromItem, toItem, compareFields)
		if len(changes) > 0 {
			result.Changed++
			for _, c := range changes {
				result.Lines = append(result.Lines, fmt.Sprintf("~ %-20s %-30s %v → %v", key, c.Field, c.From, c.To))
			}
		}
	}

	return result
}

func formatDiffLine(prefix, keyField, key string, item map[string]interface{}) string {
	return fmt.Sprintf("%s %-30s", prefix, key)
}
```

- [ ] **Step 4: Implement `cmd/am/diff/diff.go`**

```go
package diff

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

type resourceSpec struct {
	path          string
	keyField      string
	compareFields []string
	paginated     bool
}

var resourceSpecs = []struct {
	name string
	resourceSpec
}{
	{"scopes", resourceSpec{"scopes", "key", []string{"name", "description"}, true}},
	{"roles", resourceSpec{"roles", "name", []string{"description", "assignableType"}, true}},
	{"groups", resourceSpec{"groups", "name", []string{"description"}, true}},
	{"applications", resourceSpec{"applications", "name", []string{"description", "type"}, true}},
	{"identities", resourceSpec{"identities", "name", []string{"type"}, false}},
	{"certificates", resourceSpec{"certificates", "name", []string{"type"}, false}},
	{"factors", resourceSpec{"factors", "name", []string{"factorType"}, false}},
	{"flows", resourceSpec{"flows", "type", []string{"enabled"}, false}},
}

func NewDiffCmd(f *factory.Factory) *cobra.Command {
	var fromCtx, toCtx string
	var fromDomain, toDomain string

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Compare domain configuration between two contexts",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if fromCtx == "" || toCtx == "" {
				return fmt.Errorf("--from and --to are required")
			}
			cfg := f.Config
			fromResolved, err := resolveContext(cfg, fromCtx, fromDomain)
			if err != nil {
				return fmt.Errorf("--from: %w", err)
			}
			toResolved, err := resolveContext(cfg, toCtx, toDomain)
			if err != nil {
				return fmt.Errorf("--to: %w", err)
			}
			if fromResolved.Domain == "" || toResolved.Domain == "" {
				return fmt.Errorf("both contexts must have a domain set (use --from-domain / --to-domain to override)")
			}

			fromClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: fromResolved.URL, Token: fromResolved.Token})
			toClient := client.NewHTTPClient(client.HTTPClientConfig{BaseURL: toResolved.URL, Token: toResolved.Token})

			fmt.Fprintf(f.IOStreams.Out, "Comparing %s/%s → %s/%s\n\n",
				fromCtx, fromResolved.Domain, toCtx, toResolved.Domain)

			for _, spec := range resourceSpecs {
				fromItems, err := fetchItems(fromClient, fromResolved, spec.path, spec.paginated)
				if err != nil {
					fmt.Fprintf(f.IOStreams.Out, "  [%s] error fetching from: %v\n", spec.name, err)
					continue
				}
				toItems, err := fetchItems(toClient, toResolved, spec.path, spec.paginated)
				if err != nil {
					fmt.Fprintf(f.IOStreams.Out, "  [%s] error fetching to: %v\n", spec.name, err)
					continue
				}
				result := compareResources(fromItems, toItems, spec.keyField, spec.compareFields)
				if result.Added+result.Removed+result.Changed == 0 {
					fmt.Fprintf(f.IOStreams.Out, "  [%s] no differences\n", spec.name)
					continue
				}
				fmt.Fprintf(f.IOStreams.Out, "  [%s] +%d -%d ~%d\n", spec.name, result.Added, result.Removed, result.Changed)
				for _, line := range result.Lines {
					fmt.Fprintf(f.IOStreams.Out, "    %s\n", line)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&fromCtx, "from", "", "Source context name (required)")
	cmd.Flags().StringVar(&toCtx, "to", "", "Target context name (required)")
	cmd.Flags().StringVar(&fromDomain, "from-domain", "", "Override domain ID for source context")
	cmd.Flags().StringVar(&toDomain, "to-domain", "", "Override domain ID for target context")
	return cmd
}

func resolveContext(cfg *config.Config, contextName, domainOverride string) (*config.ResolvedContext, error) {
	return cfg.Resolve(config.Overrides{Context: contextName, Domain: domainOverride})
}

func fetchItems(c client.GraviteeClient, r *config.ResolvedContext, path string, paginated bool) ([]map[string]interface{}, error) {
	fullPath := fmt.Sprintf("/management/organizations/%s/environments/%s/domains/%s/%s",
		r.Org, r.Env, r.Domain, path)
	if paginated {
		fullPath += "?" + url.Values{"page": {"0"}, "size": {"1000"}}.Encode()
	}
	data, err := c.Get(fullPath)
	if err != nil {
		return nil, err
	}
	if paginated {
		var resp struct {
			Data []map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, err
		}
		return resp.Data, nil
	}
	var items []map[string]interface{}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}
```

- [ ] **Step 5: Wire in `am.go`** — `cmd.AddCommand(diffcmd.NewDiffCmd(f))`

- [ ] **Step 6: Run — verify pass**

```
go test ./cmd/am/diff/ -v
```

- [ ] **Step 7: Commit**

```bash
git add cmd/am/diff/ cmd/am/am.go
git commit -m "feat: add gio am diff command"
```

---

## Task 7: `gio am lint`

**Files:**
- Create: `cmd/am/lint/lint.go`, `cmd/am/lint/rules.go`, `cmd/am/lint/helpers_test.go`, `cmd/am/lint/lint_test.go`
- Modify: `cmd/am/am.go`

### Behaviour (lint.ts)
- Fetch domain, apps, certs, idps, scopes, factors
- Run 14 rules, each returns `[]LintFinding{rule, severity, resource, message}`
- Severity: `critical` or `warning`
- Score: `10 - 2*criticals - 1*warnings` (min 0)
- `--ci` flag: exit code 1 if any `critical` findings

Rules:
1. `implicit-grant` (critical) — app has implicit grant type
2. `no-pkce` (warning) — auth_code app has no PKCE
3. `long-token-lifetime` (warning) — access token > 1h (3600s)
4. `long-refresh-lifetime` (warning) — refresh token > 30d (2592000s)
5. `no-idp` (critical) — app has no identity providers
6. `localhost-redirect` (warning) — redirect URI contains localhost
7. `http-redirect` (warning) — redirect URI starts with http://
8. `wildcard-redirect` (critical) — redirect URI contains *
9. `app-disabled` (warning) — app.enabled == false
10. `no-factors` (warning) — domain has factors configured but app has no MFA requirement
11. `cert-expiry` (critical) — certificate `expiresAt` is within 30 days or already expired
12. `unused-scope` (warning) — scope not referenced by any app
13. `password-grant-no-mfa` (warning) — app has password grant but no MFA
14. `empty-domain` (warning) — domain has 0 apps

- [ ] **Step 1: Write the failing test**

Create `cmd/am/lint/helpers_test.go`:

```go
package lint

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient, readOnly bool) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts:       map[string]config.Context{"test": {URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT"}},
		CurrentContext: "test",
	}
	f := &factory.Factory{
		Config:   cfg,
		Resolved: &config.ResolvedContext{Name: "test", URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am", ReadOnly: readOnly},
		Client:   c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
```

Create `cmd/am/lint/lint_test.go`:

```go
package lint

import (
	"strings"
	"testing"
	"time"
)

func TestRuleImplicitGrant(t *testing.T) {
	app := map[string]interface{}{
		"name": "My App",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"grantTypes": []interface{}{"implicit"},
			},
		},
	}
	findings := ruleImplicitGrant([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != "critical" {
		t.Errorf("expected critical, got %q", findings[0].Severity)
	}
}

func TestRuleWildcardRedirect(t *testing.T) {
	app := map[string]interface{}{
		"name": "My App",
		"settings": map[string]interface{}{
			"oauth": map[string]interface{}{
				"redirectUris": []interface{}{"https://*.example.com/callback"},
			},
		},
	}
	findings := ruleWildcardRedirect([]map[string]interface{}{app})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestRuleCertExpiry(t *testing.T) {
	soon := time.Now().Add(10 * 24 * time.Hour).UnixMilli()
	cert := map[string]interface{}{
		"name":      "My Cert",
		"expiresAt": float64(soon),
	}
	findings := ruleCertExpiry([]map[string]interface{}{cert})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != "critical" {
		t.Errorf("expected critical, got %q", findings[0].Severity)
	}
}

func TestRuleEmptyDomain(t *testing.T) {
	findings := ruleEmptyDomain([]map[string]interface{}{})
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for empty domain, got %d", len(findings))
	}
}

func TestCalculateScore(t *testing.T) {
	findings := []LintFinding{
		{Severity: "critical"},
		{Severity: "critical"},
		{Severity: "warning"},
	}
	score := calculateScore(findings)
	// 10 - 2*2 - 1*1 = 5
	if score != 5 {
		t.Errorf("expected score 5, got %d", score)
	}
}

func TestScoreFloorZero(t *testing.T) {
	var findings []LintFinding
	for i := 0; i < 10; i++ {
		findings = append(findings, LintFinding{Severity: "critical"})
	}
	score := calculateScore(findings)
	if score != 0 {
		t.Errorf("expected score 0, got %d", score)
	}
}
```

Note: `ruleNoPkce` and others follow the same pattern — you may add tests for them too.

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/lint/ -v
```

- [ ] **Step 3: Implement `cmd/am/lint/rules.go`**

```go
package lint

import (
	"fmt"
	"strings"
	"time"
)

type LintFinding struct {
	Rule     string
	Severity string // "critical" or "warning"
	Resource string
	Message  string
}

func ruleImplicitGrant(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		grants := oauthGrantTypes(app)
		for _, g := range grants {
			if g == "implicit" {
				out = append(out, LintFinding{
					Rule: "implicit-grant", Severity: "critical",
					Resource: appName(app),
					Message:  "Uses implicit grant type (deprecated, insecure)",
				})
				break
			}
		}
	}
	return out
}

func ruleNoPkce(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		grants := oauthGrantTypes(app)
		hasAuthCode := false
		for _, g := range grants {
			if g == "authorization_code" {
				hasAuthCode = true
			}
		}
		if !hasAuthCode {
			continue
		}
		oauth := oauthSettings(app)
		forcePkce, _ := oauth["forcePKCE"].(bool)
		requirePkce, _ := oauth["requirePKCE"].(bool)
		if !forcePkce && !requirePkce {
			out = append(out, LintFinding{
				Rule: "no-pkce", Severity: "warning",
				Resource: appName(app),
				Message:  "Authorization code flow without PKCE enforcement",
			})
		}
	}
	return out
}

func ruleLongTokenLifetime(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		oauth := oauthSettings(app)
		if v, ok := oauth["accessTokenValiditySeconds"].(float64); ok && v > 3600 {
			out = append(out, LintFinding{
				Rule: "long-token-lifetime", Severity: "warning",
				Resource: appName(app),
				Message:  fmt.Sprintf("Access token lifetime %gs exceeds 1 hour", v),
			})
		}
	}
	return out
}

func ruleLongRefreshLifetime(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		oauth := oauthSettings(app)
		if v, ok := oauth["refreshTokenValiditySeconds"].(float64); ok && v > 2592000 {
			out = append(out, LintFinding{
				Rule: "long-refresh-lifetime", Severity: "warning",
				Resource: appName(app),
				Message:  fmt.Sprintf("Refresh token lifetime %gs exceeds 30 days", v),
			})
		}
	}
	return out
}

func ruleNoIdp(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		idps, _ := app["identityProviders"].([]interface{})
		if len(idps) == 0 {
			idps2, _ := app["identities"].([]interface{})
			if len(idps2) == 0 {
				out = append(out, LintFinding{
					Rule: "no-idp", Severity: "critical",
					Resource: appName(app),
					Message:  "No identity providers assigned",
				})
			}
		}
	}
	return out
}

func ruleLocalhostRedirect(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		for _, uri := range redirectUris(app) {
			if strings.Contains(uri, "localhost") || strings.Contains(uri, "127.0.0.1") {
				out = append(out, LintFinding{
					Rule: "localhost-redirect", Severity: "warning",
					Resource: appName(app),
					Message:  fmt.Sprintf("Redirect URI contains localhost: %s", uri),
				})
				break
			}
		}
	}
	return out
}

func ruleHttpRedirect(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		for _, uri := range redirectUris(app) {
			if strings.HasPrefix(uri, "http://") && !strings.Contains(uri, "localhost") {
				out = append(out, LintFinding{
					Rule: "http-redirect", Severity: "warning",
					Resource: appName(app),
					Message:  fmt.Sprintf("Non-secure redirect URI: %s", uri),
				})
				break
			}
		}
	}
	return out
}

func ruleWildcardRedirect(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		for _, uri := range redirectUris(app) {
			if strings.Contains(uri, "*") {
				out = append(out, LintFinding{
					Rule: "wildcard-redirect", Severity: "critical",
					Resource: appName(app),
					Message:  fmt.Sprintf("Wildcard redirect URI: %s", uri),
				})
				break
			}
		}
	}
	return out
}

func ruleAppDisabled(apps []map[string]interface{}) []LintFinding {
	var out []LintFinding
	for _, app := range apps {
		enabled, _ := app["enabled"].(bool)
		if !enabled {
			out = append(out, LintFinding{
				Rule: "app-disabled", Severity: "warning",
				Resource: appName(app),
				Message:  "Application is disabled",
			})
		}
	}
	return out
}

func ruleNoFactors(apps []map[string]interface{}, factors []map[string]interface{}) []LintFinding {
	if len(factors) == 0 {
		return nil
	}
	var out []LintFinding
	for _, app := range apps {
		appFactors, _ := app["factors"].([]interface{})
		if len(appFactors) == 0 {
			out = append(out, LintFinding{
				Rule: "no-factors", Severity: "warning",
				Resource: appName(app),
				Message:  "Domain has MFA factors but app has none assigned",
			})
		}
	}
	return out
}

func ruleCertExpiry(certs []map[string]interface{}) []LintFinding {
	var out []LintFinding
	threshold := time.Now().Add(30 * 24 * time.Hour).UnixMilli()
	for _, cert := range certs {
		expiresAt, ok := cert["expiresAt"].(float64)
		if !ok {
			continue
		}
		if int64(expiresAt) <= threshold {
			name, _ := cert["name"].(string)
			out = append(out, LintFinding{
				Rule: "cert-expiry", Severity: "critical",
				Resource: name,
				Message:  fmt.Sprintf("Certificate expires within 30 days (expiresAt: %d)", int64(expiresAt)),
			})
		}
	}
	return out
}

func ruleUnusedScope(apps []map[string]interface{}, scopes []map[string]interface{}) []LintFinding {
	used := make(map[string]bool)
	for _, app := range apps {
		oauth := oauthSettings(app)
		if scopeSettings, ok := oauth["scopeSettings"].([]interface{}); ok {
			for _, s := range scopeSettings {
				if sm, ok := s.(map[string]interface{}); ok {
					if key, ok := sm["scope"].(string); ok {
						used[key] = true
					}
				}
			}
		}
	}
	var out []LintFinding
	for _, scope := range scopes {
		key, _ := scope["key"].(string)
		if key != "" && !used[key] {
			out = append(out, LintFinding{
				Rule: "unused-scope", Severity: "warning",
				Resource: key,
				Message:  fmt.Sprintf("Scope %q is not used by any application", key),
			})
		}
	}
	return out
}

func rulePasswordGrantNoMfa(apps []map[string]interface{}, factors []map[string]interface{}) []LintFinding {
	if len(factors) == 0 {
		return nil
	}
	var out []LintFinding
	for _, app := range apps {
		grants := oauthGrantTypes(app)
		hasPassword := false
		for _, g := range grants {
			if g == "password" {
				hasPassword = true
			}
		}
		if !hasPassword {
			continue
		}
		appFactors, _ := app["factors"].([]interface{})
		if len(appFactors) == 0 {
			out = append(out, LintFinding{
				Rule: "password-grant-no-mfa", Severity: "warning",
				Resource: appName(app),
				Message:  "Password grant without MFA factor",
			})
		}
	}
	return out
}

func ruleEmptyDomain(apps []map[string]interface{}) []LintFinding {
	if len(apps) == 0 {
		return []LintFinding{{
			Rule: "empty-domain", Severity: "warning",
			Resource: "domain",
			Message:  "Domain has no applications",
		}}
	}
	return nil
}

func calculateScore(findings []LintFinding) int {
	criticals := 0
	warnings := 0
	for _, f := range findings {
		switch f.Severity {
		case "critical":
			criticals++
		case "warning":
			warnings++
		}
	}
	score := 10 - 2*criticals - warnings
	if score < 0 {
		return 0
	}
	return score
}

// Helpers

func appName(app map[string]interface{}) string {
	if name, ok := app["name"].(string); ok {
		return name
	}
	return "unknown"
}

func oauthSettings(app map[string]interface{}) map[string]interface{} {
	settings, _ := app["settings"].(map[string]interface{})
	oauth, _ := settings["oauth"].(map[string]interface{})
	return oauth
}

func oauthGrantTypes(app map[string]interface{}) []string {
	oauth := oauthSettings(app)
	raw, _ := oauth["grantTypes"].([]interface{})
	result := make([]string, 0, len(raw))
	for _, g := range raw {
		if s, ok := g.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func redirectUris(app map[string]interface{}) []string {
	oauth := oauthSettings(app)
	raw, _ := oauth["redirectUris"].([]interface{})
	result := make([]string, 0, len(raw))
	for _, u := range raw {
		if s, ok := u.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
```

- [ ] **Step 4: Implement `cmd/am/lint/lint.go`**

```go
package lint

import (
	"encoding/json"
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

type lintContext struct {
	apps    []map[string]interface{}
	certs   []map[string]interface{}
	factors []map[string]interface{}
	scopes  []map[string]interface{}
}

func NewLintCmd(f *factory.Factory) *cobra.Command {
	var ci bool

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Run security audit rules against the current domain (14 rules, scored 0-10)",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			ctx, err := collectLintData(f)
			if err != nil {
				return err
			}
			findings := runAllRules(ctx)
			score := calculateScore(findings)

			out := f.IOStreams.Out
			if len(findings) == 0 {
				fmt.Fprintf(out, "No findings. Score: 10/10\n")
				return nil
			}
			for _, finding := range findings {
				fmt.Fprintf(out, "  [%-8s] %-30s %-30s %s\n",
					finding.Severity, finding.Rule, finding.Resource, finding.Message)
			}
			fmt.Fprintf(out, "\nScore: %d/10 (%d critical, %d warning)\n",
				score, countBySeverity(findings, "critical"), countBySeverity(findings, "warning"))

			if ci && countBySeverity(findings, "critical") > 0 {
				return fmt.Errorf("lint failed: critical findings present")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&ci, "ci", false, "Exit with code 1 if any critical findings are found")
	return cmd
}

func collectLintData(f *factory.Factory) (lintContext, error) {
	var ctx lintContext
	// Apps
	data, err := f.Client.Get(cmdutil.AMDomainPath(f, "applications?page=0&size=1000"))
	if err != nil {
		return ctx, err
	}
	var appsResp struct{ Data []map[string]interface{} `json:"data"` }
	if err := json.Unmarshal(data, &appsResp); err != nil {
		return ctx, err
	}
	ctx.apps = appsResp.Data

	// Certs
	if data, err = f.Client.Get(cmdutil.AMDomainPath(f, "certificates")); err == nil {
		_ = json.Unmarshal(data, &ctx.certs)
	}
	// Factors
	if data, err = f.Client.Get(cmdutil.AMDomainPath(f, "factors")); err == nil {
		_ = json.Unmarshal(data, &ctx.factors)
	}
	// Scopes
	data, err = f.Client.Get(cmdutil.AMDomainPath(f, "scopes?page=0&size=1000"))
	if err == nil {
		var scopesResp struct{ Data []map[string]interface{} `json:"data"` }
		if json.Unmarshal(data, &scopesResp) == nil {
			ctx.scopes = scopesResp.Data
		}
	}
	return ctx, nil
}

func runAllRules(ctx lintContext) []LintFinding {
	var all []LintFinding
	all = append(all, ruleImplicitGrant(ctx.apps)...)
	all = append(all, ruleNoPkce(ctx.apps)...)
	all = append(all, ruleLongTokenLifetime(ctx.apps)...)
	all = append(all, ruleLongRefreshLifetime(ctx.apps)...)
	all = append(all, ruleNoIdp(ctx.apps)...)
	all = append(all, ruleLocalhostRedirect(ctx.apps)...)
	all = append(all, ruleHttpRedirect(ctx.apps)...)
	all = append(all, ruleWildcardRedirect(ctx.apps)...)
	all = append(all, ruleAppDisabled(ctx.apps)...)
	all = append(all, ruleNoFactors(ctx.apps, ctx.factors)...)
	all = append(all, ruleCertExpiry(ctx.certs)...)
	all = append(all, ruleUnusedScope(ctx.apps, ctx.scopes)...)
	all = append(all, rulePasswordGrantNoMfa(ctx.apps, ctx.factors)...)
	all = append(all, ruleEmptyDomain(ctx.apps)...)
	return all
}

func countBySeverity(findings []LintFinding, severity string) int {
	count := 0
	for _, f := range findings {
		if f.Severity == severity {
			count++
		}
	}
	return count
}
```

- [ ] **Step 5: Wire in `am.go`** — `cmd.AddCommand(lintcmd.NewLintCmd(f))`

- [ ] **Step 6: Run — verify pass**

```
go test ./cmd/am/lint/ -v
```

- [ ] **Step 7: Build check**

```
go build ./...
```

- [ ] **Step 8: Commit**

```bash
git add cmd/am/lint/ cmd/am/am.go
git commit -m "feat: add gio am lint command with 14 security rules"
```

---

*Continue in Part 2: Tasks 8–13 (watch, shell, test-oidc, trace, support-dump, completion)*

See: `docs/superpowers/plans/2026-05-02-am-remaining-commands-part2.md`
