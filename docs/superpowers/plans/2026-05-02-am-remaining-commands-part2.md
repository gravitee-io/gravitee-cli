# AM Remaining Commands — Implementation Plan (Part 2: Tasks 8–13)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development. Steps use checkbox (`- [ ]`) syntax for tracking.

**Continuation of:** `docs/superpowers/plans/2026-05-02-am-remaining-commands.md`

---

## Task 8: `gio am watch`

**Files:**
- Create: `cmd/am/watch/watch.go`, `cmd/am/watch/render.go`, `cmd/am/watch/helpers_test.go`, `cmd/am/watch/watch_test.go`
- Modify: `cmd/am/am.go`

### Behaviour (watch.ts)
- Poll `GET /domains/{id}/audits?page=0&size=50` every `--interval` seconds (default 5)
- Clear terminal, render dashboard with: header, stats (total/success/failure/rate), top 5 event types bar chart, top 5 errors, 15 most recent events
- `Ctrl+C` to stop (handle os.Interrupt signal)

- [ ] **Step 1: Write the failing test**

Create `cmd/am/watch/helpers_test.go`:

```go
package watch

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts:       map[string]config.Context{"test": {URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT"}},
		CurrentContext: "test",
	}
	f := &factory.Factory{
		Config:   cfg,
		Resolved: &config.ResolvedContext{Name: "test", URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am"},
		Client:   c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
```

Create `cmd/am/watch/watch_test.go`:

```go
package watch

import (
	"strings"
	"testing"
)

func TestBuildDashboardData(t *testing.T) {
	rawEvents := []map[string]interface{}{
		{
			"id": "e1", "type": "USER_LOGIN", "timestamp": float64(1700000000000),
			"outcome": map[string]interface{}{"status": "SUCCESS"},
			"actor":   map[string]interface{}{"displayName": "admin"},
		},
		{
			"id": "e2", "type": "USER_LOGIN", "timestamp": float64(1700000001000),
			"outcome": map[string]interface{}{"status": "FAILURE"},
			"actor":   map[string]interface{}{"displayName": "bob"},
		},
	}
	data := buildDashboardData(rawEvents, "my-domain", "test-ws")
	if data.Stats.Total != 2 {
		t.Errorf("expected 2 total, got %d", data.Stats.Total)
	}
	if data.Stats.Successes != 1 {
		t.Errorf("expected 1 success, got %d", data.Stats.Successes)
	}
	if data.Stats.Failures != 1 {
		t.Errorf("expected 1 failure, got %d", data.Stats.Failures)
	}
	if len(data.Stats.TopTypes) == 0 {
		t.Error("expected top types")
	}
}

func TestRender(t *testing.T) {
	data := DashboardData{
		DomainName:  "my-domain",
		Workspace:   "test-ws",
		RefreshedAt: "2023-11-14 22:13:20",
		Events: []AuditEvent{
			{ID: "e1", EventType: "USER_LOGIN", Status: "SUCCESS", Actor: "admin", Timestamp: "2023-11-14 22:13:20"},
		},
		Stats: DashboardStats{Total: 1, Successes: 1, Failures: 0},
	}
	out := render(data, 5)
	if !strings.Contains(out, "my-domain") {
		t.Error("expected domain name in render")
	}
	if !strings.Contains(out, "USER_LOGIN") {
		t.Error("expected event type in render")
	}
}
```

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/watch/ -v
```

- [ ] **Step 3: Implement `cmd/am/watch/render.go`**

```go
package watch

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type AuditEvent struct {
	ID        string
	EventType string
	Status    string
	Actor     string
	Target    string
	Timestamp string
	RawTs     int64
}

type TypeCount struct {
	Type  string
	Count int
}

type DashboardStats struct {
	Total    int
	Successes int
	Failures  int
	TopTypes  []TypeCount
	TopErrors []TypeCount
}

type DashboardData struct {
	DomainName  string
	Workspace   string
	RefreshedAt string
	Events      []AuditEvent
	Stats       DashboardStats
}

func buildDashboardData(rawEvents []map[string]interface{}, domainName, workspace string) DashboardData {
	events := make([]AuditEvent, 0, len(rawEvents))
	for _, e := range rawEvents {
		ev := AuditEvent{
			ID:        stringField(e, "id"),
			EventType: stringField(e, "type"),
		}
		if outcome, ok := e["outcome"].(map[string]interface{}); ok {
			ev.Status = stringField(outcome, "status")
		}
		if actor, ok := e["actor"].(map[string]interface{}); ok {
			ev.Actor = stringField(actor, "displayName")
		}
		if target, ok := e["target"].(map[string]interface{}); ok {
			ev.Target = stringField(target, "displayName")
		}
		if ts, ok := e["timestamp"].(float64); ok {
			ev.RawTs = int64(ts)
			ev.Timestamp = time.UnixMilli(int64(ts)).UTC().Format("2006-01-02 15:04:05")
		}
		events = append(events, ev)
	}
	// sort descending by timestamp
	sortDesc(events)

	total := len(events)
	successes := 0
	failures := 0
	typeCounts := make(map[string]int)
	errorCounts := make(map[string]int)

	for _, ev := range events {
		if ev.Status == "SUCCESS" {
			successes++
		} else {
			failures++
			errorCounts[ev.EventType]++
		}
		typeCounts[ev.EventType]++
	}

	return DashboardData{
		DomainName:  domainName,
		Workspace:   workspace,
		RefreshedAt: time.Now().UTC().Format("2006-01-02 15:04:05"),
		Events:      events,
		Stats: DashboardStats{
			Total:     total,
			Successes: successes,
			Failures:  failures,
			TopTypes:  topN(typeCounts, 5),
			TopErrors: topN(errorCounts, 5),
		},
	}
}

func render(data DashboardData, intervalSec int) string {
	var sb strings.Builder
	hr := strings.Repeat("─", 80)

	sb.WriteString(fmt.Sprintf("  Gravitee AM — %s (%s)\n", data.DomainName, data.Workspace))
	sb.WriteString(fmt.Sprintf("  %s\n\n", hr))

	successRate := 0
	if data.Stats.Total > 0 {
		successRate = int(math.Round(float64(data.Stats.Successes) / float64(data.Stats.Total) * 100))
	}
	sb.WriteString(fmt.Sprintf("  Events: %d    ✓ %d success    ✗ %d failure    Success rate: %d%%\n\n",
		data.Stats.Total, data.Stats.Successes, data.Stats.Failures, successRate))

	if len(data.Stats.TopTypes) > 0 {
		sb.WriteString("  Event types:\n")
		for _, t := range data.Stats.TopTypes {
			barLen := 0
			if data.Stats.Total > 0 {
				barLen = int(math.Round(float64(t.Count) / float64(data.Stats.Total) * 30))
			}
			bar := strings.Repeat("█", barLen)
			sb.WriteString(fmt.Sprintf("    %-25s %s %d\n", t.Type, bar, t.Count))
		}
		sb.WriteString("\n")
	}

	if len(data.Stats.TopErrors) > 0 {
		sb.WriteString("  Top errors:\n")
		for _, e := range data.Stats.TopErrors {
			sb.WriteString(fmt.Sprintf("    %-30s %d\n", e.Type, e.Count))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("  Recent events:\n")
	maxEvents := 15
	if len(data.Events) < maxEvents {
		maxEvents = len(data.Events)
	}
	for _, ev := range data.Events[:maxEvents] {
		icon := "✓"
		if ev.Status != "SUCCESS" {
			icon = "✗"
		}
		target := ""
		if ev.Target != "" {
			target = " → " + ev.Target
		}
		sb.WriteString(fmt.Sprintf("    %s %s %-20s %s%s\n",
			ev.Timestamp, icon, ev.EventType, ev.Actor, target))
	}

	sb.WriteString(fmt.Sprintf("\n  %s\n", hr))
	sb.WriteString(fmt.Sprintf("  Last refresh: %s    Interval: %ds    Ctrl+C to stop\n", data.RefreshedAt, intervalSec))
	return sb.String()
}

func stringField(m map[string]interface{}, key string) string {
	s, _ := m[key].(string)
	return s
}

func sortDesc(events []AuditEvent) {
	for i := 1; i < len(events); i++ {
		for j := i; j > 0 && events[j].RawTs > events[j-1].RawTs; j-- {
			events[j], events[j-1] = events[j-1], events[j]
		}
	}
}

func topN(counts map[string]int, n int) []TypeCount {
	result := make([]TypeCount, 0, len(counts))
	for k, v := range counts {
		result = append(result, TypeCount{Type: k, Count: v})
	}
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && result[j].Count > result[j-1].Count; j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	if len(result) > n {
		result = result[:n]
	}
	return result
}
```

- [ ] **Step 4: Implement `cmd/am/watch/watch.go`**

```go
package watch

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewWatchCmd(f *factory.Factory) *cobra.Command {
	var intervalSec int

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Live dashboard — monitor logins, errors, and audit events in real-time",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			if intervalSec < 1 {
				return fmt.Errorf("--interval must be >= 1")
			}

			domainName := f.Resolved.Domain

			refresh := func() {
				data, err := f.Client.Get(cmdutil.AMDomainPath(f, "audits?page=0&size=50"))
				if err != nil {
					return
				}
				var resp struct {
					Data []map[string]interface{} `json:"data"`
				}
				if err := json.Unmarshal(data, &resp); err != nil {
					return
				}
				dashboard := buildDashboardData(resp.Data, domainName, f.Config.CurrentContext)
				// clear screen
				fmt.Fprint(f.IOStreams.Out, "\033[2J\033[H")
				fmt.Fprint(f.IOStreams.Out, render(dashboard, intervalSec))
			}

			// Initial render
			refresh()

			ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
			defer ticker.Stop()

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

			for {
				select {
				case <-ticker.C:
					refresh()
				case <-sig:
					fmt.Fprintln(f.IOStreams.Out, "\nStopped.")
					return nil
				}
			}
		},
	}
	cmd.Flags().IntVar(&intervalSec, "interval", 5, "Refresh interval in seconds")
	return cmd
}
```

- [ ] **Step 5: Wire in `am.go`** — `cmd.AddCommand(watchcmd.NewWatchCmd(f))`

- [ ] **Step 6: Run — verify pass**

```
go test ./cmd/am/watch/ -v
```

- [ ] **Step 7: Commit**

```bash
git add cmd/am/watch/ cmd/am/am.go
git commit -m "feat: add gio am watch live dashboard command"
```

---

## Task 9: `gio am shell`

**Files:**
- Create: `cmd/am/shell/shell.go`, `cmd/am/shell/shell_test.go`
- Modify: `cmd/am/am.go`

### Behaviour (shell.ts)
- Interactive REPL: prompt `[workspace:domain] am> `
- Tab completion over known subcommands
- Built-in: `exit`/`quit`, `clear`, `help`
- Each line parsed with `splitArgs` (handles quoted strings)
- Run the parsed command via the parent cobra command

- [ ] **Step 1: Write the failing test**

Create `cmd/am/shell/shell_test.go`:

```go
package shell

import (
	"testing"
)

func TestSplitArgs(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"domain list", []string{"domain", "list"}},
		{`role get "my role"`, []string{"role", "get", "my role"}},
		{"user list --page 1", []string{"user", "list", "--page", "1"}},
		{"", []string{}},
	}
	for _, tc := range cases {
		got := splitArgs(tc.input)
		if len(got) != len(tc.expected) {
			t.Errorf("splitArgs(%q): expected %v, got %v", tc.input, tc.expected, got)
			continue
		}
		for i, v := range got {
			if v != tc.expected[i] {
				t.Errorf("splitArgs(%q)[%d]: expected %q, got %q", tc.input, i, tc.expected[i], v)
			}
		}
	}
}

func TestBuildPromptNoContext(t *testing.T) {
	p := buildPrompt("", "")
	if p == "" {
		t.Error("expected non-empty prompt")
	}
}

func TestBuildPromptWithContext(t *testing.T) {
	p := buildPrompt("myws", "dom1")
	if p == "" {
		t.Error("expected non-empty prompt")
	}
}
```

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/shell/ -v
```

- [ ] **Step 3: Implement `cmd/am/shell/shell.go`**

```go
package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func splitArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, ch := range input {
		switch {
		case inQuote:
			if ch == quoteChar {
				inQuote = false
			} else {
				current.WriteRune(ch)
			}
		case ch == '"' || ch == '\'':
			inQuote = true
			quoteChar = ch
		case ch == ' ':
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func buildPrompt(workspace, domain string) string {
	if workspace == "" {
		return "[not-configured] am> "
	}
	domainLabel := domain
	if domainLabel == "" {
		domainLabel = "(no domain)"
	}
	return fmt.Sprintf("[%s:%s] am> ", workspace, domainLabel)
}

func NewShellCmd(f *factory.Factory, parent *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:     "shell",
		Aliases: []string{"interactive"},
		Short:   "Start an interactive shell session",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			out := f.IOStreams.Out
			fmt.Fprintln(out, "\nGravitee AM CLI - Interactive Shell")
			fmt.Fprintln(out, "Type commands without the 'am' prefix. Type 'help' for available commands, 'exit' to quit.\n")

			scanner := bufio.NewScanner(os.Stdin)
			for {
				workspace := ""
				domain := ""
				if f.Config != nil {
					workspace = f.Config.CurrentContext
				}
				if f.Resolved != nil {
					domain = f.Resolved.Domain
				}
				fmt.Fprint(out, buildPrompt(workspace, domain))

				if !scanner.Scan() {
					fmt.Fprintln(out, "\nGoodbye!")
					return nil
				}
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}
				if line == "exit" || line == "quit" {
					fmt.Fprintln(out, "Goodbye!")
					return nil
				}
				if line == "clear" {
					fmt.Fprint(out, "\033[2J\033[H")
					continue
				}
				if line == "help" {
					_ = parent.Help()
					continue
				}

				args := splitArgs(line)
				parent.SetArgs(args)
				if err := parent.Execute(); err != nil {
					fmt.Fprintf(out, "Error: %v\n", err)
				}
			}
		},
	}
}
```

Note on wiring: `NewShellCmd` requires the parent `*cobra.Command` (the `am` command) so it can dispatch sub-commands. In `am.go`:

```go
shellCmd := shellcmd.NewShellCmd(f, cmd)
cmd.AddCommand(shellCmd)
```

- [ ] **Step 4: Wire in `am.go`** — see note above

- [ ] **Step 5: Run — verify pass**

```
go test ./cmd/am/shell/ -v
```

- [ ] **Step 6: Commit**

```bash
git add cmd/am/shell/ cmd/am/am.go
git commit -m "feat: add gio am shell interactive REPL"
```

---

## Task 10: `gio am test discover/login/client-credentials`

**Files:**
- Create: `cmd/am/oidctest/oidctest.go`, `cmd/am/oidctest/discover.go`, `cmd/am/oidctest/login.go`, `cmd/am/oidctest/clientcreds.go`, `cmd/am/oidctest/helpers_test.go`, `cmd/am/oidctest/oidctest_test.go`
- Modify: `cmd/am/am.go`

### Behaviour (test-oidc.ts)
- Gateway URL: `--gateway` flag or `AM_GATEWAY` env, else derive from management URL (same host, port 8092)
- `test discover`: GET `{gatewayUrl}{domainPath}/oidc/.well-known/openid-configuration`, print keys
- `test login`: ROPC flow — POST to token_endpoint with `grant_type=password`, Basic auth (confidential) or `client_id` in body (public)
- `test client-credentials`: POST `grant_type=client_credentials` with Basic auth
- Both login/cc: decode JWT id_token, validate exp/iss/aud, print decoded header+payload

These commands make direct HTTP calls (NOT via `f.Client`) since they target the gateway, not management API.

- [ ] **Step 1: Write the failing test**

Create `cmd/am/oidctest/helpers_test.go`:

```go
package oidctest

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts:       map[string]config.Context{"test": {URL: "http://am:8093", Token: "tok", Org: "DEFAULT", Env: "DEFAULT"}},
		CurrentContext: "test",
	}
	f := &factory.Factory{
		Config:   cfg,
		Resolved: &config.ResolvedContext{Name: "test", URL: "http://am:8093", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am"},
		Client:   c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
```

Create `cmd/am/oidctest/oidctest_test.go`:

```go
package oidctest

import (
	"testing"
)

func TestDeriveGatewayURL(t *testing.T) {
	cases := []struct {
		mgmtURL  string
		expected string
	}{
		{"http://am.example.com:8093", "http://am.example.com:8092"},
		{"https://am.example.com", "https://am.example.com:8092"},
		{"http://localhost:8093", "http://localhost:8092"},
	}
	for _, tc := range cases {
		got := deriveGatewayURL(tc.mgmtURL)
		if got != tc.expected {
			t.Errorf("deriveGatewayURL(%q) = %q, want %q", tc.mgmtURL, got, tc.expected)
		}
	}
}

func TestDecodeJWT(t *testing.T) {
	// A valid JWT with known payload (header.payload.signature)
	// eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJ1c2VyMSIsImV4cCI6OTk5OTk5OTk5OX0.sig
	token := "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJ1c2VyMSIsImV4cCI6OTk5OTk5OTk5OX0.sig"
	header, payload, err := decodeJWT(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if header["alg"] != "RS256" {
		t.Errorf("expected alg RS256, got %v", header["alg"])
	}
	if payload["sub"] != "user1" {
		t.Errorf("expected sub user1, got %v", payload["sub"])
	}
}

func TestTruncateToken(t *testing.T) {
	long := "abcdefghijklmnopqrstuvwxyz1234567890"
	result := truncateToken(long, 10)
	if len(result) <= 10 {
		t.Error("expected truncated string to be longer than limit due to suffix")
	}
	if result[:10] != "abcdefghij" {
		t.Errorf("unexpected prefix: %s", result[:10])
	}
}
```

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/oidctest/ -v
```

- [ ] **Step 3: Implement `cmd/am/oidctest/oidctest.go`**

```go
package oidctest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewTestCmd(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "OIDC testing utilities",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(newDiscoverCmd(f))
	cmd.AddCommand(newLoginCmd(f))
	cmd.AddCommand(newClientCredsCmd(f))
	return cmd
}

func deriveGatewayURL(mgmtURL string) string {
	parsed, err := url.Parse(mgmtURL)
	if err != nil {
		return "http://localhost:8092"
	}
	parsed.Host = parsed.Hostname() + ":8092"
	return strings.TrimRight(parsed.String(), "/")
}

func gatewayURL(flag, envVar, mgmtURL string) string {
	if flag != "" {
		return strings.TrimRight(flag, "/")
	}
	if envVar != "" {
		return strings.TrimRight(envVar, "/")
	}
	return deriveGatewayURL(mgmtURL)
}

func decodeJWT(token string) (header, payload map[string]interface{}, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("invalid JWT: expected 3 parts, got %d", len(parts))
	}
	decode := func(s string) (map[string]interface{}, error) {
		padded := s
		switch len(s) % 4 {
		case 2:
			padded += "=="
		case 3:
			padded += "="
		}
		b, err := base64.URLEncoding.DecodeString(padded)
		if err != nil {
			// try StdEncoding with +/
			b, err = base64.RawURLEncoding.DecodeString(s)
			if err != nil {
				return nil, err
			}
		}
		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, err
		}
		return m, nil
	}
	header, err = decode(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode header: %w", err)
	}
	payload, err = decode(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode payload: %w", err)
	}
	return header, payload, nil
}

func truncateToken(token string, maxLen int) string {
	if len(token) <= maxLen {
		return token
	}
	return token[:maxLen] + "...(truncated)"
}
```

- [ ] **Step 4: Implement `cmd/am/oidctest/discover.go`**

```go
package oidctest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newDiscoverCmd(f *factory.Factory) *cobra.Command {
	var gatewayFlag string

	return &cobra.Command{
		Use:   "discover",
		Short: "Fetch and display the OIDC discovery document",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			gw := gatewayURL(gatewayFlag, os.Getenv("AM_GATEWAY"), f.Resolved.URL)
			domainPath, err := fetchDomainPath(f)
			if err != nil {
				return err
			}
			discoveryURL := fmt.Sprintf("%s%s/oidc/.well-known/openid-configuration", gw, domainPath)
			data, err := httpGet(discoveryURL, "")
			if err != nil {
				return fmt.Errorf("OIDC discovery failed: %w", err)
			}
			var discovery map[string]interface{}
			if err := json.Unmarshal(data, &discovery); err != nil {
				return err
			}
			fmt.Fprintf(f.IOStreams.Out, "OIDC Discovery for %s%s\n\n", gw, domainPath)
			for key, val := range discovery {
				switch v := val.(type) {
				case []interface{}:
					fmt.Fprintf(f.IOStreams.Out, "%s:\n", key)
					for _, item := range v {
						fmt.Fprintf(f.IOStreams.Out, "  - %v\n", item)
					}
				default:
					fmt.Fprintf(f.IOStreams.Out, "%s: %v\n", key, v)
				}
			}
			return nil
		},
	}
}

func fetchDomainPath(f *factory.Factory) (string, error) {
	data, err := f.Client.Get(cmdutil.AMEnvPath(f, "domains/"+f.Resolved.Domain))
	if err != nil {
		return "", err
	}
	var domain map[string]interface{}
	if err := json.Unmarshal(data, &domain); err != nil {
		return "", err
	}
	if path, ok := domain["path"].(string); ok && path != "" {
		return path, nil
	}
	if hrid, ok := domain["hrid"].(string); ok && hrid != "" {
		return "/" + hrid, nil
	}
	return "/" + f.Resolved.Domain, nil
}

func httpGet(url, bearerToken string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	c := &http.Client{Timeout: 30 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
```

- [ ] **Step 5: Implement `cmd/am/oidctest/login.go`**

```go
package oidctest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newLoginCmd(f *factory.Factory) *cobra.Command {
	var gatewayFlag, app, secret, username, password, scope string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Test Resource Owner Password Credentials (ROPC) flow",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			if app == "" {
				return fmt.Errorf("--app is required")
			}
			if username == "" {
				return fmt.Errorf("--username is required")
			}
			if password == "" {
				return fmt.Errorf("--password is required")
			}
			gw := gatewayURL(gatewayFlag, os.Getenv("AM_GATEWAY"), f.Resolved.URL)
			domainPath, err := fetchDomainPath(f)
			if err != nil {
				return err
			}
			discovery, err := fetchDiscovery(gw, domainPath)
			if err != nil {
				return err
			}
			tokenEndpoint, _ := discovery["token_endpoint"].(string)
			if tokenEndpoint == "" {
				return fmt.Errorf("no token_endpoint in discovery")
			}

			params := url.Values{
				"grant_type": {"password"},
				"username":   {username},
				"password":   {password},
			}
			if scope != "" {
				params.Set("scope", scope)
			}

			headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
			if secret != "" {
				creds := base64.StdEncoding.EncodeToString([]byte(app + ":" + secret))
				headers["Authorization"] = "Basic " + creds
			} else {
				params.Set("client_id", app)
			}

			tokenResp, err := httpPost(tokenEndpoint, params.Encode(), headers)
			if err != nil {
				return err
			}
			printTokenResult(f, tokenResp, discovery, app)
			return nil
		},
	}
	cmd.Flags().StringVar(&gatewayFlag, "gateway", "", "Gateway URL")
	cmd.Flags().StringVar(&app, "app", "", "Application client ID (required)")
	cmd.Flags().StringVar(&secret, "secret", "", "Application client secret (omit for public clients)")
	cmd.Flags().StringVar(&username, "username", "", "Username (required)")
	cmd.Flags().StringVar(&password, "password", "", "Password (required)")
	cmd.Flags().StringVar(&scope, "scope", "", "Scopes to request (e.g. openid profile)")
	return cmd
}

func fetchDiscovery(gw, domainPath string) (map[string]interface{}, error) {
	discoveryURL := fmt.Sprintf("%s%s/oidc/.well-known/openid-configuration", gw, domainPath)
	data, err := httpGet(discoveryURL, "")
	if err != nil {
		return nil, fmt.Errorf("OIDC discovery failed: %w", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func httpPost(endpoint, body string, headers map[string]string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	c := &http.Client{Timeout: 30 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		_ = json.Unmarshal(respBody, &errResp)
		errStr, _ := errResp["error"].(string)
		desc, _ := errResp["error_description"].(string)
		return nil, fmt.Errorf("token request failed (HTTP %d): %s %s", resp.StatusCode, errStr, desc)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func printTokenResult(f *factory.Factory, tokenResp, discovery map[string]interface{}, clientID string) {
	out := f.IOStreams.Out
	accessToken, _ := tokenResp["access_token"].(string)
	fmt.Fprintf(out, "Access Token:  %s\n", truncateToken(accessToken, 40))
	fmt.Fprintf(out, "Token Type:    %v\n", tokenResp["token_type"])
	fmt.Fprintf(out, "Expires In:    %vs\n", tokenResp["expires_in"])
	fmt.Fprintf(out, "Scopes:        %v\n", tokenResp["scope"])

	if idToken, ok := tokenResp["id_token"].(string); ok {
		header, payload, err := decodeJWT(idToken)
		if err != nil {
			fmt.Fprintf(out, "\nCould not decode ID token: %v\n", err)
			return
		}
		fmt.Fprintln(out, "\nID Token (decoded):")
		fmt.Fprintln(out, "  Header:")
		for k, v := range header {
			fmt.Fprintf(out, "    %s: %v\n", k, v)
		}
		fmt.Fprintln(out, "  Payload:")
		for k, v := range payload {
			fmt.Fprintf(out, "    %s: %v\n", k, v)
		}
		fmt.Fprintln(out, "  Validation:")
		if iss, ok := payload["iss"].(string); ok {
			if discoveryIss, ok := discovery["issuer"].(string); ok {
				if iss == discoveryIss {
					fmt.Fprintln(out, "    ✓ Issuer matches discovery")
				} else {
					fmt.Fprintf(out, "    ✗ Issuer mismatch: %s vs %s\n", iss, discoveryIss)
				}
			}
		}
		aud := payload["aud"]
		audMatches := false
		switch a := aud.(type) {
		case string:
			audMatches = a == clientID
		case []interface{}:
			for _, item := range a {
				if s, ok := item.(string); ok && s == clientID {
					audMatches = true
				}
			}
		}
		if audMatches {
			fmt.Fprintln(out, "    ✓ Audience matches client_id")
		} else {
			fmt.Fprintf(out, "    ✗ Audience mismatch: %v vs %s\n", aud, clientID)
		}
	}
}
```

- [ ] **Step 6: Implement `cmd/am/oidctest/clientcreds.go`**

```go
package oidctest

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newClientCredsCmd(f *factory.Factory) *cobra.Command {
	var gatewayFlag, app, secret, scope string

	cmd := &cobra.Command{
		Use:   "client-credentials",
		Short: "Test client_credentials grant flow",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			if app == "" {
				return fmt.Errorf("--app is required")
			}
			if secret == "" {
				return fmt.Errorf("--secret is required")
			}
			gw := gatewayURL(gatewayFlag, os.Getenv("AM_GATEWAY"), f.Resolved.URL)
			domainPath, err := fetchDomainPath(f)
			if err != nil {
				return err
			}
			discovery, err := fetchDiscovery(gw, domainPath)
			if err != nil {
				return err
			}
			tokenEndpoint, _ := discovery["token_endpoint"].(string)
			if tokenEndpoint == "" {
				return fmt.Errorf("no token_endpoint in discovery")
			}

			params := url.Values{"grant_type": {"client_credentials"}}
			if scope != "" {
				params.Set("scope", scope)
			}
			creds := base64.StdEncoding.EncodeToString([]byte(app + ":" + secret))
			headers := map[string]string{
				"Content-Type":  "application/x-www-form-urlencoded",
				"Authorization": "Basic " + creds,
			}
			tokenResp, err := httpPost(tokenEndpoint, params.Encode(), headers)
			if err != nil {
				return err
			}
			printTokenResult(f, tokenResp, discovery, app)
			return nil
		},
	}
	cmd.Flags().StringVar(&gatewayFlag, "gateway", "", "Gateway URL")
	cmd.Flags().StringVar(&app, "app", "", "Application client ID (required)")
	cmd.Flags().StringVar(&secret, "secret", "", "Application client secret (required)")
	cmd.Flags().StringVar(&scope, "scope", "", "Scopes to request")
	return cmd
}
```

- [ ] **Step 7: Wire in `am.go`** — `cmd.AddCommand(oidctestcmd.NewTestCmd(f))`

- [ ] **Step 8: Run — verify pass**

```
go test ./cmd/am/oidctest/ -v
```

- [ ] **Step 9: Commit**

```bash
git add cmd/am/oidctest/ cmd/am/am.go
git commit -m "feat: add gio am test discover/login/client-credentials OIDC commands"
```

---

## Task 11: `gio am trace`

**Files:**
- Create: `cmd/am/trace/trace.go`, `cmd/am/trace/checks.go`, `cmd/am/trace/helpers_test.go`, `cmd/am/trace/trace_test.go`
- Modify: `cmd/am/am.go`

### Behaviour (trace.ts)
- `--user <idOrUsername>` and `--app <idOrName>` (both required)
- Fetch in parallel: user, application, idps, factors, flows, domain
- User resolution: try GET by ID first, then search by username
- App resolution: try GET by ID first, then search by name/clientId
- Run 7 checks: user-status, idp-match, grant-type, mfa, pre-login-flows, consent, token-config
- Output: table with ✓/⚠/✗ per check, verdict at bottom

- [ ] **Step 1: Write the failing test**

Create `cmd/am/trace/helpers_test.go`:

```go
package trace

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts:       map[string]config.Context{"test": {URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT"}},
		CurrentContext: "test",
	}
	f := &factory.Factory{
		Config:   cfg,
		Resolved: &config.ResolvedContext{Name: "test", URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am"},
		Client:   c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
```

Create `cmd/am/trace/trace_test.go`:

```go
package trace

import (
	"testing"
)

func TestCheckUserStatusEnabled(t *testing.T) {
	user := map[string]interface{}{"enabled": true, "accountNonLocked": true}
	step := checkUserStatus(user)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q: %s", step.Status, step.Detail)
	}
}

func TestCheckUserStatusDisabled(t *testing.T) {
	user := map[string]interface{}{"enabled": false}
	step := checkUserStatus(user)
	if step.Status != "block" {
		t.Errorf("expected block, got %q", step.Status)
	}
}

func TestCheckUserStatusLocked(t *testing.T) {
	user := map[string]interface{}{"enabled": true, "accountNonLocked": false}
	step := checkUserStatus(user)
	if step.Status != "block" {
		t.Errorf("expected block, got %q", step.Status)
	}
}

func TestCheckIdpMatchSuccess(t *testing.T) {
	user := map[string]interface{}{"source": "idp-1"}
	app := map[string]interface{}{
		"identityProviders": []interface{}{
			map[string]interface{}{"identity": "idp-1"},
		},
	}
	domainIdps := []map[string]interface{}{{"id": "idp-1", "name": "GitHub"}}
	step := checkIdpMatch(user, app, domainIdps)
	if step.Status != "ok" {
		t.Errorf("expected ok, got %q: %s", step.Status, step.Detail)
	}
}

func TestCheckIdpMatchFail(t *testing.T) {
	user := map[string]interface{}{"source": "idp-other"}
	app := map[string]interface{}{
		"identityProviders": []interface{}{
			map[string]interface{}{"identity": "idp-1"},
		},
	}
	domainIdps := []map[string]interface{}{{"id": "idp-1", "name": "GitHub"}}
	step := checkIdpMatch(user, app, domainIdps)
	if step.Status != "block" {
		t.Errorf("expected block, got %q", step.Status)
	}
}

func TestBuildVerdict_AllOk(t *testing.T) {
	steps := []TraceStep{
		{Status: "ok"}, {Status: "ok"},
	}
	v := buildVerdict(steps)
	if !v.CanAuthenticate {
		t.Error("expected can authenticate")
	}
	if v.Reason != "All checks passed" {
		t.Errorf("unexpected reason: %s", v.Reason)
	}
}

func TestBuildVerdict_Blocked(t *testing.T) {
	steps := []TraceStep{
		{Status: "ok"},
		{Status: "block", Detail: "User disabled"},
	}
	v := buildVerdict(steps)
	if v.CanAuthenticate {
		t.Error("expected cannot authenticate")
	}
}
```

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/trace/ -v
```

- [ ] **Step 3: Implement `cmd/am/trace/checks.go`**

```go
package trace

import "fmt"

type TraceStep struct {
	Phase  string
	Status string // "ok", "warn", "block"
	Label  string
	Detail string
}

type TraceVerdict struct {
	CanAuthenticate bool
	Reason          string
}

func checkUserStatus(user map[string]interface{}) TraceStep {
	enabled, _ := user["enabled"].(bool)
	locked, hasLocked := user["accountNonLocked"].(bool)

	if !enabled {
		return TraceStep{"user-status", "block", "User status", "User account is disabled"}
	}
	if hasLocked && !locked {
		return TraceStep{"user-status", "block", "User status", "User account is locked"}
	}
	credExp, hasCredExp := user["credentialNonExpired"].(bool)
	if hasCredExp && !credExp {
		return TraceStep{"user-status", "warn", "User status", "User credentials are expired"}
	}
	return TraceStep{"user-status", "ok", "User status", "enabled, account not locked"}
}

func checkIdpMatch(user, app map[string]interface{}, domainIdps []map[string]interface{}) TraceStep {
	appIdps := extractIdpIds(app)
	if len(appIdps) == 0 {
		return TraceStep{"idp-match", "block", "Identity source", "Application has no identity providers assigned"}
	}
	userSource, _ := user["source"].(string)
	if userSource == "" {
		return TraceStep{"idp-match", "warn", "Identity source", "User has no source identity provider set"}
	}
	for _, idpID := range appIdps {
		if idpID == userSource {
			sourceName := lookupIdpName(domainIdps, userSource)
			return TraceStep{"idp-match", "ok", "Identity source", fmt.Sprintf("%s matches app IdP", sourceName)}
		}
	}
	sourceName := lookupIdpName(domainIdps, userSource)
	appNames := make([]string, 0, len(appIdps))
	for _, id := range appIdps {
		appNames = append(appNames, lookupIdpName(domainIdps, id))
	}
	return TraceStep{"idp-match", "block", "Identity source",
		fmt.Sprintf("User's IdP '%s' not in app (app has: %v)", sourceName, appNames)}
}

func checkGrantTypes(app map[string]interface{}) TraceStep {
	oauth, _ := app["settings"].(map[string]interface{})
	oauthMap, _ := oauth["oauth"].(map[string]interface{})
	grantTypes, _ := oauthMap["grantTypes"].([]interface{})
	var userFacing []string
	for _, g := range grantTypes {
		if s, ok := g.(string); ok && (s == "password" || s == "authorization_code") {
			userFacing = append(userFacing, s)
		}
	}
	if len(userFacing) > 0 {
		return TraceStep{"grant-type", "ok", "Grant types", fmt.Sprintf("%v available", userFacing)}
	}
	allGrants := make([]string, 0, len(grantTypes))
	for _, g := range grantTypes {
		if s, ok := g.(string); ok {
			allGrants = append(allGrants, s)
		}
	}
	return TraceStep{"grant-type", "warn", "Grant types",
		fmt.Sprintf("No user-facing grant type (only: %v)", allGrants)}
}

func checkMfa(user map[string]interface{}, domainFactors []map[string]interface{}) TraceStep {
	if len(domainFactors) == 0 {
		return TraceStep{"mfa", "ok", "MFA", "MFA not required"}
	}
	userFactors, _ := user["factors"].([]interface{})
	if len(userFactors) > 0 {
		return TraceStep{"mfa", "ok", "MFA", fmt.Sprintf("MFA factor enrolled (%d)", len(userFactors))}
	}
	available := make([]string, 0, len(domainFactors))
	for _, f := range domainFactors {
		if name, ok := f["name"].(string); ok {
			available = append(available, name)
		}
	}
	return TraceStep{"mfa", "warn", "MFA",
		fmt.Sprintf("MFA required but user has no enrolled factor. Available: %v", available)}
}

func checkFlows(flows []map[string]interface{}) TraceStep {
	var policies []string
	for _, flow := range flows {
		flowType, _ := flow["type"].(string)
		if flowType != "ROOT" && flowType != "LOGIN" {
			continue
		}
		pre, _ := flow["pre"].([]interface{})
		for _, step := range pre {
			if sm, ok := step.(map[string]interface{}); ok {
				if name, ok := sm["name"].(string); ok && name != "" {
					policies = append(policies, name)
				}
			}
		}
	}
	if len(policies) > 0 {
		return TraceStep{"pre-login", "ok", "Pre-login flows",
			fmt.Sprintf("%d policies (%v)", len(policies), policies)}
	}
	return TraceStep{"pre-login", "ok", "Pre-login flows", "No pre-login policies"}
}

func checkConsent(app map[string]interface{}) TraceStep {
	advanced, _ := app["settings"].(map[string]interface{})
	adv, _ := advanced["advanced"].(map[string]interface{})
	skipConsent, _ := adv["skipConsent"].(bool)
	detail := "will be requested"
	if skipConsent {
		detail = "will be skipped"
	}
	return TraceStep{"consent", "ok", "Consent", detail}
}

func checkTokenConfig(app map[string]interface{}) TraceStep {
	settings, _ := app["settings"].(map[string]interface{})
	oauth, _ := settings["oauth"].(map[string]interface{})
	access := fmtTokenVal(oauth["accessTokenValiditySeconds"])
	refresh := fmtTokenVal(oauth["refreshTokenValiditySeconds"])
	id := fmtTokenVal(oauth["idTokenValiditySeconds"])
	return TraceStep{"token", "ok", "Token config",
		fmt.Sprintf("access=%ss, refresh=%ss, id=%ss", access, refresh, id)}
}

func buildVerdict(steps []TraceStep) TraceVerdict {
	for _, s := range steps {
		if s.Status == "block" {
			return TraceVerdict{false, s.Detail}
		}
	}
	for _, s := range steps {
		if s.Status == "warn" {
			return TraceVerdict{true, fmt.Sprintf("Likely yes, but warnings present")}
		}
	}
	return TraceVerdict{true, "All checks passed"}
}

func extractIdpIds(app map[string]interface{}) []string {
	var ids []string
	if idps, ok := app["identityProviders"].([]interface{}); ok {
		for _, idp := range idps {
			if m, ok := idp.(map[string]interface{}); ok {
				if id, ok := m["identity"].(string); ok {
					ids = append(ids, id)
				}
			}
		}
	}
	if idps, ok := app["identities"].([]interface{}); ok {
		for _, idp := range idps {
			if id, ok := idp.(string); ok {
				ids = append(ids, id)
			}
		}
	}
	return ids
}

func lookupIdpName(idps []map[string]interface{}, id string) string {
	for _, idp := range idps {
		if idp["id"] == id {
			if name, ok := idp["name"].(string); ok {
				return name
			}
		}
	}
	return id
}

func fmtTokenVal(v interface{}) string {
	if v == nil {
		return "default"
	}
	return fmt.Sprintf("%v", v)
}
```

- [ ] **Step 4: Implement `cmd/am/trace/trace.go`**

```go
package trace

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewTraceCmd(f *factory.Factory) *cobra.Command {
	var userArg, appArg string

	cmd := &cobra.Command{
		Use:   "trace",
		Short: "Trace the authentication path for a user and application",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMDomain(f); err != nil {
				return err
			}
			if userArg == "" || appArg == "" {
				return fmt.Errorf("--user and --app are required")
			}
			return runTrace(f, userArg, appArg)
		},
	}
	cmd.Flags().StringVar(&userArg, "user", "", "User to trace (username, email, or ID) (required)")
	cmd.Flags().StringVar(&appArg, "app", "", "Application (name, clientId, or ID) (required)")
	return cmd
}

func runTrace(f *factory.Factory, userArg, appArg string) error {
	// Resolve user
	user, err := resolveUser(f, userArg)
	if err != nil {
		return fmt.Errorf("user not found %q: %w", userArg, err)
	}
	// Resolve app
	app, err := resolveApp(f, appArg)
	if err != nil {
		return fmt.Errorf("application not found %q: %w", appArg, err)
	}

	// Fetch idps, factors, flows concurrently
	type result struct {
		data []map[string]interface{}
		err  error
	}
	idpCh := make(chan result, 1)
	factorCh := make(chan result, 1)
	flowCh := make(chan result, 1)

	var wg sync.WaitGroup
	fetch := func(path string, ch chan result) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := f.Client.Get(cmdutil.AMDomainPath(f, path))
			if err != nil {
				ch <- result{nil, err}
				return
			}
			var items []map[string]interface{}
			_ = json.Unmarshal(data, &items)
			ch <- result{items, nil}
		}()
	}
	fetch("identities", idpCh)
	fetch("factors", factorCh)
	fetch("flows", flowCh)
	wg.Wait()

	idps := (<-idpCh).data
	factors := (<-factorCh).data
	flows := (<-flowCh).data

	// Run checks
	steps := []TraceStep{
		checkUserStatus(user),
		checkIdpMatch(user, app, idps),
		checkGrantTypes(app),
		checkMfa(user, factors),
		checkFlows(flows),
		checkConsent(app),
		checkTokenConfig(app),
	}
	verdict := buildVerdict(steps)

	// Print
	out := f.IOStreams.Out
	userLabel := stringField(user, "email")
	if userLabel == "" {
		userLabel = stringField(user, "username")
	}
	appLabel := stringField(app, "name")
	fmt.Fprintf(out, "\nAuth flow trace: %s → %s\n\n", userLabel, appLabel)
	for _, step := range steps {
		icon := "✓"
		if step.Status == "warn" {
			icon = "⚠"
		} else if step.Status == "block" {
			icon = "✗"
		}
		fmt.Fprintf(out, "  %s %-16s %s\n", icon, step.Label, step.Detail)
	}
	fmt.Fprintln(out)
	if verdict.CanAuthenticate {
		fmt.Fprintf(out, "  Verdict: ✓ %s\n\n", verdict.Reason)
	} else {
		fmt.Fprintf(out, "  Verdict: ✗ %s\n\n", verdict.Reason)
	}
	return nil
}

func resolveUser(f *factory.Factory, userArg string) (map[string]interface{}, error) {
	// Try by ID first
	data, err := f.Client.Get(cmdutil.AMDomainPath(f, "users/"+userArg))
	if err == nil {
		var user map[string]interface{}
		if json.Unmarshal(data, &user) == nil && user["id"] != nil {
			return user, nil
		}
	}
	// Search by username
	data, err = f.Client.Get(cmdutil.AMDomainPath(f, "users?q="+userArg+"&page=0&size=1"))
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return resp.Data[0], nil
}

func resolveApp(f *factory.Factory, appArg string) (map[string]interface{}, error) {
	// Try by ID first
	data, err := f.Client.Get(cmdutil.AMDomainPath(f, "applications/"+appArg))
	if err == nil {
		var app map[string]interface{}
		if json.Unmarshal(data, &app) == nil && app["id"] != nil {
			return app, nil
		}
	}
	// Search by name
	data, err = f.Client.Get(cmdutil.AMDomainPath(f, "applications?q="+appArg+"&page=0&size=1"))
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil || len(resp.Data) == 0 {
		return nil, fmt.Errorf("application not found")
	}
	return resp.Data[0], nil
}

func stringField(m map[string]interface{}, key string) string {
	s, _ := m[key].(string)
	return s
}
```

- [ ] **Step 5: Wire in `am.go`** — `cmd.AddCommand(tracecmd.NewTraceCmd(f))`

- [ ] **Step 6: Run — verify pass**

```
go test ./cmd/am/trace/ -v
```

- [ ] **Step 7: Commit**

```bash
git add cmd/am/trace/ cmd/am/am.go
git commit -m "feat: add gio am trace auth path command"
```

---

## Task 12: `gio am support-dump`

**Files:**
- Create: `cmd/am/supportdump/supportdump.go`, `cmd/am/supportdump/redact.go`, `cmd/am/supportdump/helpers_test.go`, `cmd/am/supportdump/supportdump_test.go`
- Modify: `cmd/am/am.go`

### Behaviour (support-dump.ts)
- Collect from domain: domain config, apps, idps, certs, flows, factors, roles, scopes, groups, members, plus optionally users (PII, opt-in) and audits
- Collect platform plugins: identities, factors, certificates, resources, policies, reporters, bot-detections
- `--all-domains`: collect from all domains in environment
- `--no-audit`: skip audit logs
- `--audit-size N` (default 100): number of audit events
- `--include-users`: include user list (PII)
- `--no-redact`: disable secret redaction
- Redact keys matching: secret, password, private, credential, apiKey, api_key, token, *key (but NOT the safe set)
- Output: JSON to stdout or `-f <file>`

- [ ] **Step 1: Write the failing test**

Create `cmd/am/supportdump/helpers_test.go`:

```go
package supportdump

import (
	"bytes"

	"github.com/gravitee-io/gio-cli/internal/client"
	"github.com/gravitee-io/gio-cli/internal/config"
	"github.com/gravitee-io/gio-cli/internal/factory"
)

func newTestFactory(c client.GraviteeClient) (*factory.Factory, *bytes.Buffer) {
	out := &bytes.Buffer{}
	cfg := &config.Config{
		Contexts:       map[string]config.Context{"test": {URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT"}},
		CurrentContext: "test",
	}
	f := &factory.Factory{
		Config:   cfg,
		Resolved: &config.ResolvedContext{Name: "test", URL: "http://am", Token: "tok", Org: "DEFAULT", Env: "DEFAULT", Domain: "dom1", Type: "am"},
		Client:   c,
		IOStreams: factory.IOStreams{Out: out},
	}
	return f, out
}
```

Create `cmd/am/supportdump/supportdump_test.go`:

```go
package supportdump

import (
	"testing"
)

func TestRedactSecrets(t *testing.T) {
	input := map[string]interface{}{
		"name":         "my-cert",
		"clientSecret": "super-secret",
		"publicKey":    "pk-value",
	}
	result := redactSecrets(input)
	m := result.(map[string]interface{})
	if m["clientSecret"] != "[REDACTED]" {
		t.Errorf("expected clientSecret to be redacted, got %v", m["clientSecret"])
	}
	if m["publicKey"] != "pk-value" {
		t.Errorf("expected publicKey to be preserved, got %v", m["publicKey"])
	}
	if m["name"] != "my-cert" {
		t.Errorf("expected name to be preserved, got %v", m["name"])
	}
}

func TestShouldRedactKey(t *testing.T) {
	cases := []struct {
		key      string
		expected bool
	}{
		{"clientSecret", true},
		{"password", true},
		{"privateKey", true},
		{"apiKey", true},
		{"publicKey", false},
		{"tokenEndpoint", false},
		{"passwordPolicy", false},
		{"name", false},
	}
	for _, tc := range cases {
		got := shouldRedactKey(tc.key)
		if got != tc.expected {
			t.Errorf("shouldRedactKey(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}

func TestRedactSecretsNested(t *testing.T) {
	input := map[string]interface{}{
		"configuration": map[string]interface{}{
			"clientId":     "abc",
			"clientSecret": "secret123",
		},
	}
	result := redactSecrets(input)
	m := result.(map[string]interface{})
	conf := m["configuration"].(map[string]interface{})
	if conf["clientSecret"] != "[REDACTED]" {
		t.Errorf("expected nested secret to be redacted")
	}
	if conf["clientId"] != "abc" {
		t.Errorf("expected clientId to be preserved")
	}
}
```

- [ ] **Step 2: Run — verify fail**

```
go test ./cmd/am/supportdump/ -v
```

- [ ] **Step 3: Implement `cmd/am/supportdump/redact.go`**

```go
package supportdump

import (
	"regexp"
	"strings"
)

const redactPlaceholder = "[REDACTED]"

var redactPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)secret`),
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)private`),
	regexp.MustCompile(`(?i)credential`),
	regexp.MustCompile(`(?i)apiKey`),
	regexp.MustCompile(`(?i)api_key`),
	regexp.MustCompile(`(?i)token`),
	regexp.MustCompile(`(?i)key$`),
}

var safeKeys = map[string]bool{
	"tokenEndpoint":              true,
	"tokenExchangeSettings":      true,
	"tokenExpiresIn":             true,
	"passwordPolicy":             true,
	"passwordSettings":           true,
	"passwordPolicies":           true,
	"secretExpirationSettings":   true,
	"accessTokenValiditySeconds": true,
	"refreshTokenValiditySeconds": true,
	"idTokenValiditySeconds":     true,
	"publicKey":                  true,
	"publicKeys":                 true,
	"keyId":                      true,
}

func shouldRedactKey(key string) bool {
	if safeKeys[key] {
		return false
	}
	for _, p := range redactPatterns {
		if p.MatchString(key) {
			return true
		}
	}
	return false
}

func redactSecrets(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	switch v := obj.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			if shouldRedactKey(key) {
				if s, ok := val.(string); ok && s != "" {
					result[key] = redactPlaceholder
					continue
				}
			}
			result[key] = redactSecrets(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = redactSecrets(item)
		}
		return result
	default:
		return v
	}
}

func stringsJoin(strs []string, sep string) string {
	return strings.Join(strs, sep)
}
```

- [ ] **Step 4: Implement `cmd/am/supportdump/supportdump.go`**

```go
package supportdump

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewSupportDumpCmd(f *factory.Factory) *cobra.Command {
	var outputFile string
	var allDomains bool
	var noAudit bool
	var auditSize int
	var includeUsers bool
	var noRedact bool

	cmd := &cobra.Command{
		Use:   "support-dump",
		Short: "Generate a comprehensive support diagnostic dump",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := cmdutil.RequireAMContext(f); err != nil {
				return err
			}
			shouldRedact := !noRedact

			var domainIDs []string
			if allDomains {
				data, err := f.Client.Get(cmdutil.AMEnvPath(f, "domains?page=0&size=1000"))
				if err != nil {
					return err
				}
				var resp struct {
					Data []map[string]interface{} `json:"data"`
				}
				if err := json.Unmarshal(data, &resp); err != nil {
					return err
				}
				for _, d := range resp.Data {
					if id, ok := d["id"].(string); ok {
						domainIDs = append(domainIDs, id)
					}
				}
			} else {
				if err := cmdutil.RequireAMDomain(f); err != nil {
					return err
				}
				domainIDs = []string{f.Resolved.Domain}
			}

			output := map[string]interface{}{
				"_metadata": map[string]interface{}{
					"serverUrl":      f.Resolved.URL,
					"organizationId": f.Resolved.Org,
					"environmentId":  f.Resolved.Env,
					"secretsRedacted": shouldRedact,
					"includesUsers":  includeUsers,
					"includesAudit":  !noAudit,
					"domainCount":    len(domainIDs),
				},
			}

			if allDomains {
				var domains []interface{}
				for _, domainID := range domainIDs {
					sections, errs := collectDomain(f, domainID, !noAudit, auditSize, includeUsers)
					entry := map[string]interface{}{"domainId": domainID}
					for k, v := range sections {
						entry[k] = v
					}
					if len(errs) > 0 {
						entry["_errors"] = errs
					}
					domains = append(domains, entry)
				}
				output["domains"] = domains
			} else {
				sections, errs := collectDomain(f, domainIDs[0], !noAudit, auditSize, includeUsers)
				for k, v := range sections {
					output[k] = v
				}
				if len(errs) > 0 {
					output["_errors"] = errs
				}
			}

			if shouldRedact {
				output = redactSecrets(output).(map[string]interface{})
			}

			jsonBytes, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return err
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, jsonBytes, 0600); err != nil {
					return err
				}
				fmt.Fprintf(f.IOStreams.Out, "Support dump written to %s\n", outputFile)
			} else {
				fmt.Fprintln(f.IOStreams.Out, string(jsonBytes))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&outputFile, "file", "f", "", "Output file path (default: stdout)")
	cmd.Flags().BoolVar(&allDomains, "all-domains", false, "Dump all domains in the environment")
	cmd.Flags().BoolVar(&noAudit, "no-audit", false, "Skip audit logs")
	cmd.Flags().IntVar(&auditSize, "audit-size", 100, "Number of recent audit events to include")
	cmd.Flags().BoolVar(&includeUsers, "include-users", false, "Include user list (contains PII)")
	cmd.Flags().BoolVar(&noRedact, "no-redact", false, "Disable secret redaction")
	return cmd
}

func collectDomain(f *factory.Factory, domainID string, includeAudit bool, auditSize int, includeUsers bool) (map[string]interface{}, []string) {
	sections := make(map[string]interface{})
	var errs []string

	get := func(label, path string, dest *interface{}) {
		data, err := f.Client.Get(path)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", label, err))
			return
		}
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			errs = append(errs, fmt.Sprintf("%s: parse error", label))
			return
		}
		*dest = v
	}

	domainPath := func(p string) string {
		return cmdutil.AMDomainPathFor(f, domainID, p)
	}

	var domain, apps, idps, certs, flows, factors, roles, scopes, groups, members, users, audits interface{}
	get("domain", domainPath(""), &domain)
	get("applications", domainPath("applications?page=0&size=1000"), &apps)
	get("identities", domainPath("identities"), &idps)
	get("certificates", domainPath("certificates"), &certs)
	get("flows", domainPath("flows"), &flows)
	get("factors", domainPath("factors"), &factors)
	get("roles", domainPath("roles?page=0&size=1000"), &roles)
	get("scopes", domainPath("scopes?page=0&size=1000"), &scopes)
	get("groups", domainPath("groups?page=0&size=1000"), &groups)
	get("members", domainPath("members"), &members)

	if includeUsers {
		get("users", domainPath("users?page=0&size=1000"), &users)
	}
	if includeAudit {
		get("audits", domainPath(fmt.Sprintf("audits?page=0&size=%d", auditSize)), &audits)
	}

	setIfNotNil := func(key string, val interface{}) {
		if val != nil {
			sections[key] = val
		}
	}
	setIfNotNil("domain", domain)
	setIfNotNil("applications", extractData(apps))
	setIfNotNil("identityProviders", idps)
	setIfNotNil("certificates", certs)
	setIfNotNil("flows", flows)
	setIfNotNil("factors", factors)
	setIfNotNil("roles", extractData(roles))
	setIfNotNil("scopes", extractData(scopes))
	setIfNotNil("groups", extractData(groups))
	setIfNotNil("members", members)
	if includeUsers {
		setIfNotNil("users", extractData(users))
	}
	if includeAudit {
		setIfNotNil("recentAudits", extractData(audits))
	}

	return sections, errs
}

func extractData(v interface{}) interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		if data, ok := m["data"]; ok {
			return data
		}
	}
	return v
}
```

- [ ] **Step 5: Wire in `am.go`** — `cmd.AddCommand(supportdumpcmd.NewSupportDumpCmd(f))`

- [ ] **Step 6: Run — verify pass**

```
go test ./cmd/am/supportdump/ -v
```

- [ ] **Step 7: Commit**

```bash
git add cmd/am/supportdump/ cmd/am/am.go
git commit -m "feat: add gio am support-dump diagnostic command"
```

---

## Task 13: `gio am completion`

**Files:**
- Create: `cmd/am/completion.go`
- Modify: `cmd/am/am.go`

### Behaviour (completion.ts)
- Cobra has built-in shell completion — just expose it as `gio am completion <shell>`
- Supported: bash, zsh, fish, powershell

- [ ] **Step 1: Write the failing test** — add to `cmd/am/health_test.go`:

```go
func TestCompletion(t *testing.T) {
    f, out := newAMTestFactory(nil, false)
    cmd := newCompletionCmd(f)
    cmd.SetArgs([]string{"bash"})
    if err := cmd.Execute(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out.String(), "bash") && !strings.Contains(out.String(), "compgen") {
        // completion output varies by shell — just check it's non-empty
        if out.Len() == 0 {
            t.Error("expected non-empty completion output")
        }
    }
}
```

Actually cobra's `GenBashCompletion` writes to a writer. The easiest approach is to delegate to cobra's root command. Since this is the `am` sub-command, we need to pass the parent.

Simpler: just expose `cobra.Command` with `ValidArgs: []string{"bash", "zsh", "fish", "powershell"}` and generate manually.

- [ ] **Step 2: Implement `cmd/am/completion.go`**

```go
package am

import (
	"fmt"

	"github.com/gravitee-io/gio-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCompletionCmd(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Generate shell completion script",
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			out := f.IOStreams.Out
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(out)
			case "zsh":
				return root.GenZshCompletion(out)
			case "fish":
				return root.GenFishCompletion(out, true)
			case "powershell":
				return root.GenPowerShellCompletionWithDesc(out)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}
}
```

- [ ] **Step 3: Wire in `am.go`** — `cmd.AddCommand(newCompletionCmd(f))`

- [ ] **Step 4: Build check**

```
go build ./...
```

- [ ] **Step 5: Run all tests**

```
go test ./cmd/am/... ./internal/...
```

Expected: all PASS

- [ ] **Step 6: Commit**

```bash
git add cmd/am/completion.go cmd/am/am.go
git commit -m "feat: add gio am completion shell completion command"
```

---

## Final Wiring — Update `cmd/am/am.go`

After all tasks are done, `am.go` should import and wire all new packages:

```go
import (
    logscmd       "github.com/gravitee-io/gio-cli/cmd/am/logs"
    plugincmd     "github.com/gravitee-io/gio-cli/cmd/am/plugin"
    diffcmd       "github.com/gravitee-io/gio-cli/cmd/am/diff"
    lintcmd       "github.com/gravitee-io/gio-cli/cmd/am/lint"
    watchcmd      "github.com/gravitee-io/gio-cli/cmd/am/watch"
    shellcmd      "github.com/gravitee-io/gio-cli/cmd/am/shell"
    oidctestcmd   "github.com/gravitee-io/gio-cli/cmd/am/oidctest"
    tracecmd      "github.com/gravitee-io/gio-cli/cmd/am/trace"
    supportdumpcmd "github.com/gravitee-io/gio-cli/cmd/am/supportdump"
    // ... existing imports
)
```

And in `NewAMCmd`:
```go
cmd.AddCommand(newLogoutCmd(f))
cmd.AddCommand(newStatusCmd(f))
cmd.AddCommand(newDoctorCmd(f))
cmd.AddCommand(newCompletionCmd(f))
cmd.AddCommand(logscmd.NewLogsCmd(f))
cmd.AddCommand(plugincmd.NewPluginCmd(f))
cmd.AddCommand(diffcmd.NewDiffCmd(f))
cmd.AddCommand(lintcmd.NewLintCmd(f))
cmd.AddCommand(watchcmd.NewWatchCmd(f))
cmd.AddCommand(shellcmd.NewShellCmd(f, cmd))
cmd.AddCommand(oidctestcmd.NewTestCmd(f))
cmd.AddCommand(tracecmd.NewTraceCmd(f))
cmd.AddCommand(supportdumpcmd.NewSupportDumpCmd(f))
```

- [ ] **Final build + test**

```
go build ./...
go test ./...
```

- [ ] **Final commit**

```bash
git add cmd/am/am.go
git commit -m "feat: wire all remaining AM commands into gio am"
```
