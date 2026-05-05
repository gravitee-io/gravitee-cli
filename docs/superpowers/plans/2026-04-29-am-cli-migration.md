# AM CLI Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port all remaining Gravitee Access Management functionality from the TypeScript `am-cli` into the Go `gio-cli`, under `gio am <resource> <action>`.

**Architecture:** Each entity follows the same two-layer pattern used by APIM: a service interface + HTTP implementation in `internal/am/<entity>.go`, and cobra commands in `cmd/am/<entity>/`. Domain-scoped commands (app, user, idp, etc.) require a `--domain <id>` flag. The `internal/am/service.go` `Service` interface is extended by embedding new service interfaces as each entity is added.

**Tech Stack:** Go 1.26, cobra, `internal/client` (HTTP), `internal/printer` (table/JSON/YAML output), `internal/cmdutil` (shared flag/output helpers), `internal/factory` (dependency injection).

---

## Existing AM baseline (already done)

| Location | What exists |
|---|---|
| `internal/am/service.go` | `Service` interface embedding `DomainService`; `service` struct; `basePath()` |
| `internal/am/domain.go` | `ListDomains`, `GetDomain` service methods |
| `internal/am/pagination.go` | `PaginatedResponse`, `FetchAllPages` |
| `cmd/am/am.go` | `NewAMCmd` registering `domain` sub-tree |
| `cmd/am/domain/domain.go` | `NewDomainCmd` registering `list` |
| `cmd/am/domain/list.go` | `gio am domain list` (with `--all`, pagination) |

---

## Task 1: Domain CRUD (get, create, enable, disable, delete)

**Files:**
- Modify: `internal/am/domain.go` — add `CreateDomain`, `PatchDomain`, `DeleteDomain`
- Modify: `internal/am/service.go` — `DomainService` interface extended
- Create: `cmd/am/domain/get.go`
- Create: `cmd/am/domain/create.go`
- Create: `cmd/am/domain/enable.go`
- Create: `cmd/am/domain/disable.go`
- Create: `cmd/am/domain/delete.go`
- Modify: `cmd/am/domain/domain.go` — register new subcommands

- [ ] **Step 1: Extend DomainService interface and service methods**

In `internal/am/domain.go`, add after the existing `GetDomain`:

```go
// CreateDomainBody holds fields for domain creation.
type CreateDomainBody struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    DataPlaneID string `json:"dataPlaneId,omitempty"`
}

// PatchDomainBody holds fields for partial domain update.
type PatchDomainBody struct {
    Enabled *bool `json:"enabled,omitempty"`
}

func (s *service) CreateDomain(body CreateDomainBody) (json.RawMessage, error) {
    b, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }
    data, err := s.client.Post(s.basePath("domains"), b)
    if err != nil {
        return nil, fmt.Errorf("domain create failed: %w", err)
    }
    return json.RawMessage(data), nil
}

func (s *service) PatchDomain(domainID string, body PatchDomainBody) error {
    b, err := json.Marshal(body)
    if err != nil {
        return err
    }
    _, err = s.client.Patch(s.basePath(fmt.Sprintf("domains/%s", domainID)), b)
    if err != nil {
        return fmt.Errorf("domain patch failed: %w", err)
    }
    return nil
}

func (s *service) DeleteDomain(domainID string) error {
    if err := s.client.Delete(s.basePath(fmt.Sprintf("domains/%s", domainID))); err != nil {
        return fmt.Errorf("domain delete failed: %w", err)
    }
    return nil
}
```

Extend the `DomainService` interface in `internal/am/domain.go`:

```go
type DomainService interface {
    ListDomains(params ListDomainsParams) (*PaginatedResponse, error)
    GetDomain(domainID string) (json.RawMessage, error)
    CreateDomain(body CreateDomainBody) (json.RawMessage, error)
    PatchDomain(domainID string, body PatchDomainBody) error
    DeleteDomain(domainID string) error
}
```

- [ ] **Step 2: Verify `internal/client` has Post, Patch, Delete methods**

Run:
```bash
grep -n "Post\|Patch\|Delete" /Users/rpo/Documents/Projects/Gravitee/gio-cli/internal/client/client.go
```

If `Patch` or `Delete` are missing, add them to `client.go` following the same pattern as `Post`:

```go
// In internal/client/http.go (or client.go), add:
func (c *httpClient) Patch(path string, body []byte) ([]byte, error) {
    return c.do("PATCH", path, body)
}

func (c *httpClient) Delete(path string) error {
    _, err := c.do("DELETE", path, nil)
    return err
}
```

Also add to `GraviteeClient` interface in `client.go`:
```go
Patch(path string, body []byte) ([]byte, error)
Delete(path string) error
```

And to `FakeClient` in `fake.go`:
```go
PatchFunc func(path string, body []byte) ([]byte, error)
DeleteFunc func(path string) error

func (f *FakeClient) Patch(path string, body []byte) ([]byte, error) {
    if f.PatchFunc != nil { return f.PatchFunc(path, body) }
    return nil, nil
}
func (f *FakeClient) Delete(path string) error {
    if f.DeleteFunc != nil { return f.DeleteFunc(path) }
    return nil
}
```

- [ ] **Step 3: Write `cmd/am/domain/get.go`**

```go
package domain

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

func newGetCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "get <domainId>",
        Short:   "Get security domain details",
        Example: `  gio am domain get abc-123`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            data, err := f.AM().GetDomain(args[0])
            if err != nil {
                return err
            }
            p, err := cmdutil.NewPrinter(f)
            if err != nil {
                return err
            }
            if f.OutputFormat != printer.FormatTable {
                return p.PrintDetail(data)
            }
            return printDomainDetail(p, data)
        },
    }
}
```

Add `printDomainDetail` to `cmd/am/domain/list.go` (or a new `cmd/am/domain/helpers.go`):

```go
func printDomainDetail(p *printer.Printer, data []byte) error {
    var m map[string]any
    if err := json.Unmarshal(data, &m); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    for _, field := range []struct{ label, key string }{
        {"Name", "name"}, {"ID", "id"}, {"Enabled", "enabled"},
        {"Description", "description"}, {"Data Plane", "dataPlaneId"},
    } {
        if v, ok := m[field.key]; ok && v != nil {
            p.PrintMessage("%-16s%v", field.label+":", v)
        }
    }
    return nil
}
```

- [ ] **Step 4: Write `cmd/am/domain/create.go`**

```go
package domain

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/am"
    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
    var name, description, dataPlaneID string

    cmd := &cobra.Command{
        Use:     "create",
        Short:   "Create a new security domain",
        Example: `  gio am domain create --name my-domain --description "My domain"`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            data, err := f.AM().CreateDomain(am.CreateDomainBody{
                Name:        name,
                Description: description,
                DataPlaneID: dataPlaneID,
            })
            if err != nil {
                return err
            }
            p, err := cmdutil.NewPrinter(f)
            if err != nil {
                return err
            }
            if f.OutputFormat != printer.FormatTable {
                return p.PrintDetail(data)
            }
            return printDomainDetail(p, data)
        },
    }

    cmd.Flags().StringVarP(&name, "name", "n", "", "Domain name (required)")
    cmd.Flags().StringVarP(&description, "description", "d", "", "Domain description")
    cmd.Flags().StringVar(&dataPlaneID, "data-plane-id", "default", "Data plane ID")
    _ = cmd.MarkFlagRequired("name")

    return cmd
}
```

- [ ] **Step 5: Write `cmd/am/domain/enable.go` and `cmd/am/domain/disable.go`**

`enable.go`:
```go
package domain

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/am"
    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newEnableCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "enable <domainId>",
        Short:   "Enable a security domain",
        Example: `  gio am domain enable abc-123`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            t := true
            if err := f.AM().PatchDomain(args[0], am.PatchDomainBody{Enabled: &t}); err != nil {
                return err
            }
            p, _ := cmdutil.NewPrinter(f)
            p.PrintMessage("Domain '%s' enabled.", args[0])
            return nil
        },
    }
}
```

`disable.go` — identical but `Enabled: &f` (false) and message "disabled":
```go
package domain

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/am"
    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newDisableCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "disable <domainId>",
        Short:   "Disable a security domain",
        Example: `  gio am domain disable abc-123`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            fal := false
            if err := f.AM().PatchDomain(args[0], am.PatchDomainBody{Enabled: &fal}); err != nil {
                return err
            }
            p, _ := cmdutil.NewPrinter(f)
            p.PrintMessage("Domain '%s' disabled.", args[0])
            return nil
        },
    }
}
```

- [ ] **Step 6: Write `cmd/am/domain/delete.go`**

```go
package domain

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "delete <domainId>",
        Short:   "Delete a security domain",
        Example: `  gio am domain delete abc-123`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            if err := f.AM().DeleteDomain(args[0]); err != nil {
                return err
            }
            p, _ := cmdutil.NewPrinter(f)
            p.PrintMessage("Domain '%s' deleted.", args[0])
            return nil
        },
    }
}
```

- [ ] **Step 7: Register new commands in `cmd/am/domain/domain.go`**

```go
cmd.AddCommand(
    newListCmd(f),
    newGetCmd(f),
    newCreateCmd(f),
    newEnableCmd(f),
    newDisableCmd(f),
    newDeleteCmd(f),
)
```

- [ ] **Step 8: Build and smoke test**

```bash
cd /Users/rpo/Documents/Projects/Gravitee/gio-cli && go build ./...
```
Expected: no errors.

- [ ] **Step 9: Commit**

```bash
git add internal/am/domain.go internal/client/ cmd/am/domain/
git commit -m "feat(am): add domain CRUD commands (get, create, enable, disable, delete)"
```

---

## Task 2: Health + Whoami top-level AM commands

**Files:**
- Create: `cmd/am/health.go`
- Create: `cmd/am/whoami.go`
- Modify: `cmd/am/am.go` — register health and whoami

These are lightweight commands that call AM management endpoints directly via the factory's HTTP client; no new service layer needed.

- [ ] **Step 1: Write `cmd/am/health.go`**

```go
package am

import (
    "encoding/json"

    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

func newHealthCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:     "health",
        Aliases: []string{"ping"},
        Short:   "Check if the AM instance is reachable",
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            // Check the factory's raw client (no auth needed for health).
            // AM health endpoint: GET /management/health
            data, err := f.HTTPClient().Get("/management/health")
            if err != nil {
                return fmt.Errorf("AM instance unreachable: %w", err)
            }
            p, err := cmdutil.NewPrinter(f)
            if err != nil {
                return err
            }
            if f.OutputFormat != printer.FormatTable {
                return p.PrintDetail(json.RawMessage(data))
            }
            p.PrintMessage("AM instance is healthy.")
            return nil
        },
    }
}
```

> **Note:** If `factory.Factory` does not expose `HTTPClient()`, use `f.AM().(interface{ RawGet(string) ([]byte, error) })` or add a lightweight health method to the AM service. Adjust as needed once you check `internal/factory/factory.go`.

- [ ] **Step 2: Write `cmd/am/whoami.go`**

```go
package am

import (
    "encoding/json"

    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

func newWhoamiCmd(f *factory.Factory) *cobra.Command {
    return &cobra.Command{
        Use:   "whoami",
        Short: "Show information about the currently authenticated user",
        Args:  cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            data, err := f.AM().(interface {
                WhoAmI() (json.RawMessage, error)
            }).WhoAmI()
            if err != nil {
                return err
            }
            p, err := cmdutil.NewPrinter(f)
            if err != nil {
                return err
            }
            return p.PrintDetail(data)
        },
    }
}
```

Add `WhoAmI() (json.RawMessage, error)` to `internal/am/service.go` `Service` interface and implement in a new file `internal/am/whoami.go`:

```go
package am

import "encoding/json"

type WhoAmIService interface {
    WhoAmI() (json.RawMessage, error)
}

func (s *service) WhoAmI() (json.RawMessage, error) {
    data, err := s.client.Get("/management/user")
    if err != nil {
        return nil, err
    }
    return json.RawMessage(data), nil
}
```

Embed `WhoAmIService` in `Service`.

- [ ] **Step 3: Register in `cmd/am/am.go`**

```go
cmd.AddCommand(
    domaincmd.NewDomainCmd(f),
    newHealthCmd(f),
    newWhoamiCmd(f),
)
```

- [ ] **Step 4: Build**

```bash
cd /Users/rpo/Documents/Projects/Gravitee/gio-cli && go build ./...
```

- [ ] **Step 5: Commit**

```bash
git add cmd/am/health.go cmd/am/whoami.go cmd/am/am.go internal/am/
git commit -m "feat(am): add health and whoami commands"
```

---

## Task 3: Application (app) commands

AM API: `GET/POST /management/organizations/{org}/environments/{env}/domains/{domain}/applications`

**Files:**
- Create: `internal/am/application.go`
- Modify: `internal/am/service.go` — embed `ApplicationService`
- Create: `cmd/am/application/application.go`
- Create: `cmd/am/application/list.go`
- Create: `cmd/am/application/get.go`
- Create: `cmd/am/application/create.go`
- Create: `cmd/am/application/update.go`
- Create: `cmd/am/application/delete.go`
- Create: `cmd/am/application/helpers.go`
- Modify: `cmd/am/am.go` — register application command

All domain-scoped entity commands in this and subsequent tasks follow the same flag pattern: `--domain <id>` (required).

- [ ] **Step 1: Write `internal/am/application.go`**

```go
package am

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/client"
)

// ListApplicationsParams holds parameters for listing applications.
type ListApplicationsParams struct {
    DomainID string
    Query    string
    Page     int
    PerPage  int
}

// ApplicationService defines application operations.
type ApplicationService interface {
    ListApplications(params ListApplicationsParams) (*PaginatedResponse, error)
    GetApplication(domainID, appID string) (json.RawMessage, error)
    CreateApplication(domainID string, body json.RawMessage) (json.RawMessage, error)
    UpdateApplication(domainID, appID string, body json.RawMessage) (json.RawMessage, error)
    DeleteApplication(domainID, appID string) error
}

func (s *service) ListApplications(params ListApplicationsParams) (*PaginatedResponse, error) {
    q := client.BuildQuery(map[string]string{
        "page": client.Itoa(params.Page),
        "size": client.Itoa(params.PerPage),
        "q":    params.Query,
    })
    data, err := s.client.Get(s.domainPath(params.DomainID, "applications?"+q))
    if err != nil {
        return nil, fmt.Errorf("application list failed: %w", err)
    }
    return parsePaginatedResponse(data)
}

func (s *service) GetApplication(domainID, appID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("applications/%s", appID)))
    if err != nil {
        return nil, err
    }
    return json.RawMessage(data), nil
}

func (s *service) CreateApplication(domainID string, body json.RawMessage) (json.RawMessage, error) {
    data, err := s.client.Post(s.domainPath(domainID, "applications"), body)
    if err != nil {
        return nil, fmt.Errorf("application create failed: %w", err)
    }
    return json.RawMessage(data), nil
}

func (s *service) UpdateApplication(domainID, appID string, body json.RawMessage) (json.RawMessage, error) {
    data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("applications/%s", appID)), body)
    if err != nil {
        return nil, fmt.Errorf("application update failed: %w", err)
    }
    return json.RawMessage(data), nil
}

func (s *service) DeleteApplication(domainID, appID string) error {
    if err := s.client.Delete(s.domainPath(domainID, fmt.Sprintf("applications/%s", appID))); err != nil {
        return fmt.Errorf("application delete failed: %w", err)
    }
    return nil
}
```

Add `domainPath` helper to `internal/am/service.go`:
```go
func (s *service) domainPath(domainID, path string) string {
    return s.basePath(fmt.Sprintf("domains/%s/%s", domainID, path))
}
```

Embed `ApplicationService` in `Service`:
```go
type Service interface {
    DomainService
    WhoAmIService
    ApplicationService
}
```

Also check if `client.GraviteeClient` has a `Put` method; if not, add it following the same pattern as `Post`.

- [ ] **Step 2: Write `cmd/am/application/helpers.go`**

```go
package application

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

func appColumns() []printer.Column {
    return []printer.Column{
        {Name: "Name", Value: func(i any) string { return cmdutil.StringField(i, "name") }},
        {Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
        {Name: "Type", Value: func(i any) string { return cmdutil.StringField(i, "type") }},
        {Name: "Status", Value: func(i any) string { return cmdutil.StringField(i, "status") }},
    }
}

func printAppDetail(p *printer.Printer, data []byte) error {
    var m map[string]any
    if err := json.Unmarshal(data, &m); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    for _, f := range []struct{ label, key string }{
        {"Name", "name"}, {"ID", "id"}, {"Type", "type"}, {"Status", "status"},
        {"Description", "description"},
    } {
        if v, ok := m[f.key]; ok && v != nil {
            p.PrintMessage("%-16s%v", f.label+":", v)
        }
    }
    return nil
}
```

- [ ] **Step 3: Write `cmd/am/application/application.go`**

```go
package application

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func NewApplicationCmd(f *factory.Factory) *cobra.Command {
    cmd := &cobra.Command{
        Use:     "application",
        Aliases: []string{"app"},
        Short:   "Manage applications",
    }
    cmdutil.AddOutputFlags(cmd, f)
    cmd.AddCommand(
        newListCmd(f),
        newGetCmd(f),
        newCreateCmd(f),
        newUpdateCmd(f),
        newDeleteCmd(f),
    )
    return cmd
}
```

- [ ] **Step 4: Write `cmd/am/application/list.go`**

```go
package application

import (
    "encoding/json"

    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/am"
    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

type listOptions struct {
    factory  *factory.Factory
    domainID string
    query    string
    page     int
    perPage  int
    all      bool
}

func newListCmd(f *factory.Factory) *cobra.Command {
    opts := &listOptions{factory: f}
    cmd := &cobra.Command{
        Use:     "list",
        Short:   "List applications",
        Example: `  gio am application list --domain abc-123`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            return opts.run()
        },
    }
    cmd.Flags().StringVar(&opts.domainID, "domain", "", "Security domain ID (required)")
    cmd.Flags().StringVar(&opts.query, "query", "", "Search by name")
    cmd.Flags().IntVar(&opts.page, "page", 1, "Page number")
    cmd.Flags().IntVar(&opts.perPage, "per-page", 10, "Results per page")
    cmd.Flags().BoolVar(&opts.all, "all", false, "Fetch all pages")
    _ = cmd.MarkFlagRequired("domain")
    return cmd
}

func (o *listOptions) run() error {
    f := o.factory
    p, err := cmdutil.NewPrinter(f)
    if err != nil {
        return err
    }
    if o.all {
        allData, err := am.FetchAllPages(func(page int) (*am.PaginatedResponse, error) {
            return f.AM().ListApplications(am.ListApplicationsParams{
                DomainID: o.domainID, Query: o.query, Page: page, PerPage: o.perPage,
            })
        }, o.perPage)
        if err != nil {
            return err
        }
        if f.OutputFormat != printer.FormatTable {
            return p.PrintDetail(allData)
        }
        return p.PrintList(allData, appColumns())
    }
    resp, err := f.AM().ListApplications(am.ListApplicationsParams{
        DomainID: o.domainID, Query: o.query, Page: o.page - 1, PerPage: o.perPage,
    })
    if err != nil {
        return err
    }
    if f.OutputFormat != printer.FormatTable {
        raw, _ := json.Marshal(resp)
        return p.PrintDetail(json.RawMessage(raw))
    }
    return p.PrintList(resp.Data, appColumns())
}
```

- [ ] **Step 5: Write `cmd/am/application/get.go`**

```go
package application

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

func newGetCmd(f *factory.Factory) *cobra.Command {
    var domainID string
    cmd := &cobra.Command{
        Use:     "get <appId>",
        Short:   "Get application details",
        Example: `  gio am application get app-uuid --domain domain-uuid`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            data, err := f.AM().GetApplication(domainID, args[0])
            if err != nil {
                return err
            }
            p, err := cmdutil.NewPrinter(f)
            if err != nil {
                return err
            }
            if f.OutputFormat != printer.FormatTable {
                return p.PrintDetail(data)
            }
            return printAppDetail(p, data)
        },
    }
    cmd.Flags().StringVar(&domainID, "domain", "", "Security domain ID (required)")
    _ = cmd.MarkFlagRequired("domain")
    return cmd
}
```

- [ ] **Step 6: Write `cmd/am/application/create.go`**

```go
package application

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

func newCreateCmd(f *factory.Factory) *cobra.Command {
    var domainID, file string
    cmd := &cobra.Command{
        Use:     "create -f <file>",
        Short:   "Create an application from a JSON file",
        Example: `  gio am application create --domain abc-123 -f app.json`,
        Args:    cobra.NoArgs,
        RunE: func(_ *cobra.Command, _ []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            body, err := cmdutil.ReadJSONFile(file)
            if err != nil {
                return err
            }
            data, err := f.AM().CreateApplication(domainID, body)
            if err != nil {
                return err
            }
            p, err := cmdutil.NewPrinter(f)
            if err != nil {
                return err
            }
            if f.OutputFormat != printer.FormatTable {
                return p.PrintDetail(data)
            }
            return printAppDetail(p, data)
        },
    }
    cmd.Flags().StringVar(&domainID, "domain", "", "Security domain ID (required)")
    cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
    _ = cmd.MarkFlagRequired("domain")
    _ = cmd.MarkFlagRequired("file")
    return cmd
}
```

- [ ] **Step 7: Write `cmd/am/application/update.go`**

```go
package application

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

func newUpdateCmd(f *factory.Factory) *cobra.Command {
    var domainID, file string
    cmd := &cobra.Command{
        Use:     "update <appId> -f <file>",
        Short:   "Update an application from a JSON file",
        Example: `  gio am application update app-uuid --domain abc-123 -f app.json`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            body, err := cmdutil.ReadJSONFile(file)
            if err != nil {
                return err
            }
            data, err := f.AM().UpdateApplication(domainID, args[0], body)
            if err != nil {
                return err
            }
            p, err := cmdutil.NewPrinter(f)
            if err != nil {
                return err
            }
            if f.OutputFormat != printer.FormatTable {
                return p.PrintDetail(data)
            }
            return printAppDetail(p, data)
        },
    }
    cmd.Flags().StringVar(&domainID, "domain", "", "Security domain ID (required)")
    cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON definition file (required)")
    _ = cmd.MarkFlagRequired("domain")
    _ = cmd.MarkFlagRequired("file")
    return cmd
}
```

- [ ] **Step 8: Write `cmd/am/application/delete.go`**

```go
package application

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newDeleteCmd(f *factory.Factory) *cobra.Command {
    var domainID string
    cmd := &cobra.Command{
        Use:     "delete <appId>",
        Short:   "Delete an application",
        Example: `  gio am application delete app-uuid --domain abc-123`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            if err := f.AM().DeleteApplication(domainID, args[0]); err != nil {
                return err
            }
            p, _ := cmdutil.NewPrinter(f)
            p.PrintMessage("Application '%s' deleted.", args[0])
            return nil
        },
    }
    cmd.Flags().StringVar(&domainID, "domain", "", "Security domain ID (required)")
    _ = cmd.MarkFlagRequired("domain")
    return cmd
}
```

- [ ] **Step 9: Register in `cmd/am/am.go`**

```go
import applicationcmd "github.com/gravitee-io/gio-cli/cmd/am/application"
// ...
cmd.AddCommand(
    domaincmd.NewDomainCmd(f),
    applicationcmd.NewApplicationCmd(f),
    newHealthCmd(f),
    newWhoamiCmd(f),
)
```

- [ ] **Step 10: Build**

```bash
cd /Users/rpo/Documents/Projects/Gravitee/gio-cli && go build ./...
```

- [ ] **Step 11: Commit**

```bash
git add internal/am/application.go cmd/am/application/ cmd/am/am.go
git commit -m "feat(am): add application commands (list, get, create, update, delete)"
```

---

## Task 4: User commands

AM API: `GET/POST /management/organizations/{org}/environments/{env}/domains/{domain}/users`

**Files:**
- Create: `internal/am/user.go`
- Modify: `internal/am/service.go`
- Create: `cmd/am/user/` (user.go, list.go, get.go, create.go, update.go, delete.go, lock.go, unlock.go, reset_password.go, helpers.go)
- Modify: `cmd/am/am.go`

- [ ] **Step 1: Write `internal/am/user.go`**

```go
package am

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/client"
)

// ListUsersParams holds parameters for listing users.
type ListUsersParams struct {
    DomainID string
    Query    string
    Page     int
    PerPage  int
}

// UserService defines user operations.
type UserService interface {
    ListUsers(params ListUsersParams) (*PaginatedResponse, error)
    GetUser(domainID, userID string) (json.RawMessage, error)
    CreateUser(domainID string, body json.RawMessage) (json.RawMessage, error)
    UpdateUser(domainID, userID string, body json.RawMessage) (json.RawMessage, error)
    DeleteUser(domainID, userID string) error
    LockUser(domainID, userID string) error
    UnlockUser(domainID, userID string) error
    ResetPassword(domainID, userID, newPassword string) error
}

func (s *service) ListUsers(params ListUsersParams) (*PaginatedResponse, error) {
    q := client.BuildQuery(map[string]string{
        "page": client.Itoa(params.Page),
        "size": client.Itoa(params.PerPage),
        "q":    params.Query,
    })
    data, err := s.client.Get(s.domainPath(params.DomainID, "users?"+q))
    if err != nil {
        return nil, fmt.Errorf("user list failed: %w", err)
    }
    return parsePaginatedResponse(data)
}

func (s *service) GetUser(domainID, userID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s", userID)))
    if err != nil {
        return nil, err
    }
    return json.RawMessage(data), nil
}

func (s *service) CreateUser(domainID string, body json.RawMessage) (json.RawMessage, error) {
    data, err := s.client.Post(s.domainPath(domainID, "users"), body)
    if err != nil {
        return nil, fmt.Errorf("user create failed: %w", err)
    }
    return json.RawMessage(data), nil
}

func (s *service) UpdateUser(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
    data, err := s.client.Put(s.domainPath(domainID, fmt.Sprintf("users/%s", userID)), body)
    if err != nil {
        return nil, fmt.Errorf("user update failed: %w", err)
    }
    return json.RawMessage(data), nil
}

func (s *service) DeleteUser(domainID, userID string) error {
    return s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s", userID)))
}

func (s *service) LockUser(domainID, userID string) error {
    _, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("users/%s/lock", userID)), nil)
    return err
}

func (s *service) UnlockUser(domainID, userID string) error {
    _, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("users/%s/unlock", userID)), nil)
    return err
}

func (s *service) ResetPassword(domainID, userID, newPassword string) error {
    body, _ := json.Marshal(map[string]string{"password": newPassword})
    _, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("users/%s/resetPassword", userID)), body)
    return err
}
```

Embed `UserService` in `Service` in `internal/am/service.go`.

- [ ] **Step 2: Write `cmd/am/user/helpers.go`**

```go
package user

import (
    "encoding/json"
    "fmt"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/printer"
)

func userColumns() []printer.Column {
    return []printer.Column{
        {Name: "Username", Value: func(i any) string { return cmdutil.StringField(i, "username") }},
        {Name: "ID", Value: func(i any) string { return cmdutil.StringField(i, "id") }},
        {Name: "Email", Value: func(i any) string { return cmdutil.StringField(i, "email") }},
        {Name: "Enabled", Value: func(i any) string {
            m, ok := i.(map[string]any)
            if !ok { return "" }
            if v, ok := m["enabled"].(bool); ok && v { return "true" }
            return "false"
        }},
    }
}

func printUserDetail(p *printer.Printer, data []byte) error {
    var m map[string]any
    if err := json.Unmarshal(data, &m); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    for _, f := range []struct{ label, key string }{
        {"Username", "username"}, {"ID", "id"}, {"Email", "email"},
        {"First Name", "firstName"}, {"Last Name", "lastName"}, {"Enabled", "enabled"},
    } {
        if v, ok := m[f.key]; ok && v != nil {
            p.PrintMessage("%-16s%v", f.label+":", v)
        }
    }
    return nil
}
```

- [ ] **Step 3: Write `cmd/am/user/user.go`, `list.go`, `get.go`, `create.go`, `update.go`, `delete.go`, `lock.go`, `unlock.go`, `reset_password.go`**

`user.go`:
```go
package user

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func NewUserCmd(f *factory.Factory) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "user",
        Short: "Manage users",
    }
    cmdutil.AddOutputFlags(cmd, f)
    cmd.AddCommand(
        newListCmd(f),
        newGetCmd(f),
        newCreateCmd(f),
        newUpdateCmd(f),
        newDeleteCmd(f),
        newLockCmd(f),
        newUnlockCmd(f),
        newResetPasswordCmd(f),
    )
    return cmd
}
```

`list.go` — same pattern as `application/list.go` but using `f.AM().ListUsers(am.ListUsersParams{...})` and `userColumns()`.

`get.go` — `f.AM().GetUser(domainID, args[0])`, print with `printUserDetail`.

`create.go` — `--domain` + `-f` flags, `f.AM().CreateUser(domainID, body)`.

`update.go` — `--domain` + `-f` flags, `f.AM().UpdateUser(domainID, args[0], body)`.

`delete.go` — `--domain`, `f.AM().DeleteUser(domainID, args[0])`.

`lock.go`:
```go
func newLockCmd(f *factory.Factory) *cobra.Command {
    var domainID string
    cmd := &cobra.Command{
        Use: "lock <userId>", Short: "Lock a user account",
        Args: cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil { return err }
            if err := f.AM().LockUser(domainID, args[0]); err != nil { return err }
            p, _ := cmdutil.NewPrinter(f)
            p.PrintMessage("User '%s' locked.", args[0])
            return nil
        },
    }
    cmd.Flags().StringVar(&domainID, "domain", "", "Security domain ID (required)")
    _ = cmd.MarkFlagRequired("domain")
    return cmd
}
```

`unlock.go` — same as lock but calls `f.AM().UnlockUser(...)` and prints "unlocked".

`reset_password.go`:
```go
func newResetPasswordCmd(f *factory.Factory) *cobra.Command {
    var domainID, password string
    cmd := &cobra.Command{
        Use: "reset-password <userId>", Short: "Reset a user's password",
        Args: cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil { return err }
            if err := f.AM().ResetPassword(domainID, args[0], password); err != nil { return err }
            p, _ := cmdutil.NewPrinter(f)
            p.PrintMessage("Password reset for user '%s'.", args[0])
            return nil
        },
    }
    cmd.Flags().StringVar(&domainID, "domain", "", "Security domain ID (required)")
    cmd.Flags().StringVar(&password, "password", "", "New password (required)")
    _ = cmd.MarkFlagRequired("domain")
    _ = cmd.MarkFlagRequired("password")
    return cmd
}
```

- [ ] **Step 4: Register in `cmd/am/am.go`**, build, commit.

```bash
git commit -m "feat(am): add user commands (list, get, create, update, delete, lock, unlock, reset-password)"
```

---

## Task 5: Identity Provider (idp), Factor commands

AM API:
- IdPs: `GET /management/.../domains/{domain}/identities`
- Factors: `GET /management/.../domains/{domain}/factors`

**Files:**
- Create: `internal/am/idp.go`, `internal/am/factor.go`
- Modify: `internal/am/service.go`
- Create: `cmd/am/idp/` (idp.go, list.go, get.go, helpers.go)
- Create: `cmd/am/factor/` (factor.go, list.go, get.go, helpers.go)
- Modify: `cmd/am/am.go`

- [ ] **Step 1: Write `internal/am/idp.go`**

```go
package am

import (
    "encoding/json"
    "fmt"
)

// IDPService defines identity provider operations.
type IDPService interface {
    ListIDPs(domainID string) (json.RawMessage, error)
    GetIDP(domainID, idpID string) (json.RawMessage, error)
}

func (s *service) ListIDPs(domainID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, "identities"))
    if err != nil {
        return nil, fmt.Errorf("idp list failed: %w", err)
    }
    return json.RawMessage(data), nil
}

func (s *service) GetIDP(domainID, idpID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("identities/%s", idpID)))
    if err != nil {
        return nil, err
    }
    return json.RawMessage(data), nil
}
```

- [ ] **Step 2: Write `internal/am/factor.go`**

```go
package am

import (
    "encoding/json"
    "fmt"
)

// FactorService defines MFA factor operations.
type FactorService interface {
    ListFactors(domainID string) (json.RawMessage, error)
    GetFactor(domainID, factorID string) (json.RawMessage, error)
}

func (s *service) ListFactors(domainID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, "factors"))
    if err != nil {
        return nil, fmt.Errorf("factor list failed: %w", err)
    }
    return json.RawMessage(data), nil
}

func (s *service) GetFactor(domainID, factorID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("factors/%s", factorID)))
    if err != nil {
        return nil, err
    }
    return json.RawMessage(data), nil
}
```

Embed both in `Service`.

- [ ] **Step 3: Write `cmd/am/idp/` commands**

`idp.go`:
```go
package idp

import (
    "github.com/spf13/cobra"
    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func NewIDPCmd(f *factory.Factory) *cobra.Command {
    cmd := &cobra.Command{Use: "idp", Short: "Manage identity providers"}
    cmdutil.AddOutputFlags(cmd, f)
    cmd.AddCommand(newListCmd(f), newGetCmd(f))
    return cmd
}
```

`list.go` — `--domain` flag, calls `f.AM().ListIDPs(domainID)`, prints as JSON array or table with columns Name, ID, Type.

`get.go` — `--domain` + `<idpId>` arg, calls `f.AM().GetIDP(domainID, id)`, prints detail.

`helpers.go` — `idpColumns()` and `printIDPDetail()`.

- [ ] **Step 4: Write `cmd/am/factor/` commands** (same structure as idp, using factor endpoints)

- [ ] **Step 5: Register in `cmd/am/am.go`**, build, commit.

```bash
git commit -m "feat(am): add idp and factor commands (list, get)"
```

---

## Task 6: Role + Scope commands

AM API:
- Roles: `GET/POST /management/.../domains/{domain}/roles`
- Scopes: `GET/POST /management/.../domains/{domain}/scopes`

**Files:**
- Create: `internal/am/role.go`, `internal/am/scope.go`
- Modify: `internal/am/service.go`
- Create: `cmd/am/role/` (role.go, list.go, get.go, create.go, helpers.go)
- Create: `cmd/am/scope/` (scope.go, list.go, get.go, create.go, helpers.go)
- Modify: `cmd/am/am.go`

- [ ] **Step 1: Write `internal/am/role.go`**

```go
package am

import (
    "encoding/json"
    "fmt"
    "github.com/gravitee-io/gio-cli/internal/client"
)

type ListRolesParams struct {
    DomainID string
    Query    string
    Page     int
    PerPage  int
}

type RoleService interface {
    ListRoles(params ListRolesParams) (*PaginatedResponse, error)
    GetRole(domainID, roleID string) (json.RawMessage, error)
    CreateRole(domainID string, body json.RawMessage) (json.RawMessage, error)
}

func (s *service) ListRoles(params ListRolesParams) (*PaginatedResponse, error) {
    q := client.BuildQuery(map[string]string{
        "page": client.Itoa(params.Page), "size": client.Itoa(params.PerPage), "q": params.Query,
    })
    data, err := s.client.Get(s.domainPath(params.DomainID, "roles?"+q))
    if err != nil {
        return nil, fmt.Errorf("role list failed: %w", err)
    }
    return parsePaginatedResponse(data)
}

func (s *service) GetRole(domainID, roleID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("roles/%s", roleID)))
    if err != nil { return nil, err }
    return json.RawMessage(data), nil
}

func (s *service) CreateRole(domainID string, body json.RawMessage) (json.RawMessage, error) {
    data, err := s.client.Post(s.domainPath(domainID, "roles"), body)
    if err != nil {
        return nil, fmt.Errorf("role create failed: %w", err)
    }
    return json.RawMessage(data), nil
}
```

- [ ] **Step 2: Write `internal/am/scope.go`** (same shape, endpoint `scopes` instead of `roles`)

- [ ] **Step 3: Write `cmd/am/role/`** — role.go, list.go (paginated), get.go, create.go (-f flag), helpers.go.

- [ ] **Step 4: Write `cmd/am/scope/`** — same structure.

- [ ] **Step 5: Register in `cmd/am/am.go`**, build, commit.

```bash
git commit -m "feat(am): add role and scope commands (list, get, create)"
```

---

## Task 7: Group + Flow commands

AM API:
- Groups: `GET/POST/DELETE /management/.../domains/{domain}/groups`
- Flows: `GET /management/.../domains/{domain}/flows`

**Files:**
- Create: `internal/am/group.go`, `internal/am/flow.go`
- Modify: `internal/am/service.go`
- Create: `cmd/am/group/` (group.go, list.go, get.go, create.go, delete.go, helpers.go)
- Create: `cmd/am/flow/` (flow.go, list.go, get.go, helpers.go)
- Modify: `cmd/am/am.go`

- [ ] **Step 1: Write `internal/am/group.go`**

```go
package am

import (
    "encoding/json"
    "fmt"
    "github.com/gravitee-io/gio-cli/internal/client"
)

type ListGroupsParams struct {
    DomainID string
    Query    string
    Page     int
    PerPage  int
}

type GroupService interface {
    ListGroups(params ListGroupsParams) (*PaginatedResponse, error)
    GetGroup(domainID, groupID string) (json.RawMessage, error)
    CreateGroup(domainID string, body json.RawMessage) (json.RawMessage, error)
    DeleteGroup(domainID, groupID string) error
}

func (s *service) ListGroups(params ListGroupsParams) (*PaginatedResponse, error) {
    q := client.BuildQuery(map[string]string{
        "page": client.Itoa(params.Page), "size": client.Itoa(params.PerPage), "q": params.Query,
    })
    data, err := s.client.Get(s.domainPath(params.DomainID, "groups?"+q))
    if err != nil { return nil, fmt.Errorf("group list failed: %w", err) }
    return parsePaginatedResponse(data)
}

func (s *service) GetGroup(domainID, groupID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("groups/%s", groupID)))
    if err != nil { return nil, err }
    return json.RawMessage(data), nil
}

func (s *service) CreateGroup(domainID string, body json.RawMessage) (json.RawMessage, error) {
    data, err := s.client.Post(s.domainPath(domainID, "groups"), body)
    if err != nil { return nil, fmt.Errorf("group create failed: %w", err) }
    return json.RawMessage(data), nil
}

func (s *service) DeleteGroup(domainID, groupID string) error {
    return s.client.Delete(s.domainPath(domainID, fmt.Sprintf("groups/%s", groupID)))
}
```

- [ ] **Step 2: Write `internal/am/flow.go`**

```go
package am

import (
    "encoding/json"
    "fmt"
)

type FlowService interface {
    ListFlows(domainID string) (json.RawMessage, error)
    GetFlow(domainID, flowID string) (json.RawMessage, error)
}

func (s *service) ListFlows(domainID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, "flows"))
    if err != nil { return nil, fmt.Errorf("flow list failed: %w", err) }
    return json.RawMessage(data), nil
}

func (s *service) GetFlow(domainID, flowID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("flows/%s", flowID)))
    if err != nil { return nil, err }
    return json.RawMessage(data), nil
}
```

- [ ] **Step 3: Write `cmd/am/group/`** — group.go, list.go, get.go, create.go (-f), delete.go, helpers.go.

- [ ] **Step 4: Write `cmd/am/flow/`** — flow.go, list.go, get.go, helpers.go.

- [ ] **Step 5: Register in `cmd/am/am.go`**, build, commit.

```bash
git commit -m "feat(am): add group and flow commands"
```

---

## Task 8: Certificate + Audit commands

AM API:
- Certificates: `GET/DELETE /management/.../domains/{domain}/certificates`
- Audit logs: `GET /management/.../domains/{domain}/audits`

**Files:**
- Create: `internal/am/certificate.go`, `internal/am/audit.go`
- Modify: `internal/am/service.go`
- Create: `cmd/am/certificate/` (certificate.go, list.go, get.go, delete.go, helpers.go)
- Create: `cmd/am/audit/` (audit.go, list.go, get.go, helpers.go)
- Modify: `cmd/am/am.go`

- [ ] **Step 1: Write `internal/am/certificate.go`**

```go
package am

import (
    "encoding/json"
    "fmt"
)

type CertificateService interface {
    ListCertificates(domainID string) (json.RawMessage, error)
    GetCertificate(domainID, certID string) (json.RawMessage, error)
    DeleteCertificate(domainID, certID string) error
}

func (s *service) ListCertificates(domainID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, "certificates"))
    if err != nil { return nil, fmt.Errorf("certificate list failed: %w", err) }
    return json.RawMessage(data), nil
}

func (s *service) GetCertificate(domainID, certID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("certificates/%s", certID)))
    if err != nil { return nil, err }
    return json.RawMessage(data), nil
}

func (s *service) DeleteCertificate(domainID, certID string) error {
    return s.client.Delete(s.domainPath(domainID, fmt.Sprintf("certificates/%s", certID)))
}
```

- [ ] **Step 2: Write `internal/am/audit.go`**

```go
package am

import (
    "encoding/json"
    "fmt"
    "github.com/gravitee-io/gio-cli/internal/client"
)

type ListAuditsParams struct {
    DomainID string
    Type     string
    Status   string
    Page     int
    PerPage  int
}

type AuditService interface {
    ListAudits(params ListAuditsParams) (*PaginatedResponse, error)
    GetAudit(domainID, auditID string) (json.RawMessage, error)
}

func (s *service) ListAudits(params ListAuditsParams) (*PaginatedResponse, error) {
    q := client.BuildQuery(map[string]string{
        "page": client.Itoa(params.Page), "size": client.Itoa(params.PerPage),
        "type": params.Type, "status": params.Status,
    })
    data, err := s.client.Get(s.domainPath(params.DomainID, "audits?"+q))
    if err != nil { return nil, fmt.Errorf("audit list failed: %w", err) }
    return parsePaginatedResponse(data)
}

func (s *service) GetAudit(domainID, auditID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("audits/%s", auditID)))
    if err != nil { return nil, err }
    return json.RawMessage(data), nil
}
```

- [ ] **Step 3: Write `cmd/am/certificate/`** — list (array output), get, delete.

- [ ] **Step 4: Write `cmd/am/audit/`** — list (paginated, `--type`, `--status` flags), get.

- [ ] **Step 5: Register in `cmd/am/am.go`**, build, commit.

```bash
git commit -m "feat(am): add certificate and audit commands"
```

---

## Task 9: Token management commands

AM API: `GET/POST/DELETE /management/.../domains/{domain}/users/{userId}/tokens`

**Files:**
- Create: `internal/am/token.go`
- Modify: `internal/am/service.go`
- Create: `cmd/am/token/` (token.go, list.go, create.go, revoke.go, helpers.go)
- Modify: `cmd/am/am.go`

- [ ] **Step 1: Write `internal/am/token.go`**

```go
package am

import (
    "encoding/json"
    "fmt"
)

type TokenService interface {
    ListTokens(domainID, userID string) (json.RawMessage, error)
    CreateToken(domainID, userID string, body json.RawMessage) (json.RawMessage, error)
    RevokeToken(domainID, userID, tokenID string) error
}

func (s *service) ListTokens(domainID, userID string) (json.RawMessage, error) {
    data, err := s.client.Get(s.domainPath(domainID, fmt.Sprintf("users/%s/tokens", userID)))
    if err != nil { return nil, fmt.Errorf("token list failed: %w", err) }
    return json.RawMessage(data), nil
}

func (s *service) CreateToken(domainID, userID string, body json.RawMessage) (json.RawMessage, error) {
    data, err := s.client.Post(s.domainPath(domainID, fmt.Sprintf("users/%s/tokens", userID)), body)
    if err != nil { return nil, fmt.Errorf("token create failed: %w", err) }
    return json.RawMessage(data), nil
}

func (s *service) RevokeToken(domainID, userID, tokenID string) error {
    return s.client.Delete(s.domainPath(domainID, fmt.Sprintf("users/%s/tokens/%s", userID, tokenID)))
}
```

- [ ] **Step 2: Write `cmd/am/token/`**

`token.go`:
```go
package token

import (
    "github.com/spf13/cobra"
    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func NewTokenCmd(f *factory.Factory) *cobra.Command {
    cmd := &cobra.Command{Use: "token", Short: "Manage user tokens"}
    cmdutil.AddOutputFlags(cmd, f)
    cmd.AddCommand(newListCmd(f), newCreateCmd(f), newRevokeCmd(f))
    return cmd
}
```

`list.go` — `--domain` + `--user` flags, calls `f.AM().ListTokens(domainID, userID)`, prints as JSON.

`create.go` — `--domain` + `--user` + `-f` flags.

`revoke.go` — `--domain` + `--user` + `<tokenId>` arg, calls `f.AM().RevokeToken(...)`.

- [ ] **Step 3: Register in `cmd/am/am.go`**, build, commit.

```bash
git commit -m "feat(am): add token management commands (list, create, revoke)"
```

---

## Task 10: Domain export, import, copy

These are complex orchestration commands that combine multiple service calls.

**Files:**
- Create: `cmd/am/domain/export.go`
- Create: `cmd/am/domain/import.go`
- Create: `cmd/am/domain/copy.go`
- Modify: `cmd/am/domain/domain.go`

- [ ] **Step 1: Write `cmd/am/domain/export.go`**

Export fetches domain + all child resources in parallel and writes JSON.

```go
package domain

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"

    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/am"
    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newExportCmd(f *factory.Factory) *cobra.Command {
    var file string
    cmd := &cobra.Command{
        Use:     "export <domainId>",
        Short:   "Export domain configuration to JSON",
        Example: `  gio am domain export abc-123 -f domain-export.json`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
                return err
            }
            return runExport(f, args[0], file)
        },
    }
    cmd.Flags().StringVarP(&file, "file", "f", "", "Output file path (default: stdout)")
    return cmd
}

func runExport(f *factory.Factory, domainID, file string) error {
    svc := f.AM()
    domainData, err := svc.GetDomain(domainID)
    if err != nil {
        return fmt.Errorf("failed to fetch domain: %w", err)
    }

    // Fetch all child resources concurrently.
    type result struct {
        key  string
        data json.RawMessage
        err  error
    }

    jobs := []struct {
        key string
        fn  func() (json.RawMessage, error)
    }{
        {"applications", func() (json.RawMessage, error) {
            all, err := am.FetchAllPages(func(page int) (*am.PaginatedResponse, error) {
                return svc.ListApplications(am.ListApplicationsParams{DomainID: domainID, Page: page, PerPage: 100})
            }, 100)
            if err != nil { return nil, err }
            b, _ := json.Marshal(all)
            return json.RawMessage(b), nil
        }},
        {"identityProviders", func() (json.RawMessage, error) { return svc.ListIDPs(domainID) }},
        {"roles", func() (json.RawMessage, error) {
            all, err := am.FetchAllPages(func(page int) (*am.PaginatedResponse, error) {
                return svc.ListRoles(am.ListRolesParams{DomainID: domainID, Page: page, PerPage: 100})
            }, 100)
            if err != nil { return nil, err }
            b, _ := json.Marshal(all)
            return json.RawMessage(b), nil
        }},
        {"scopes", func() (json.RawMessage, error) {
            all, err := am.FetchAllPages(func(page int) (*am.PaginatedResponse, error) {
                return svc.ListScopes(am.ListScopesParams{DomainID: domainID, Page: page, PerPage: 100})
            }, 100)
            if err != nil { return nil, err }
            b, _ := json.Marshal(all)
            return json.RawMessage(b), nil
        }},
        {"factors", func() (json.RawMessage, error) { return svc.ListFactors(domainID) }},
        {"groups", func() (json.RawMessage, error) {
            all, err := am.FetchAllPages(func(page int) (*am.PaginatedResponse, error) {
                return svc.ListGroups(am.ListGroupsParams{DomainID: domainID, Page: page, PerPage: 100})
            }, 100)
            if err != nil { return nil, err }
            b, _ := json.Marshal(all)
            return json.RawMessage(b), nil
        }},
        {"flows", func() (json.RawMessage, error) { return svc.ListFlows(domainID) }},
        {"certificates", func() (json.RawMessage, error) { return svc.ListCertificates(domainID) }},
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
        return firstErr
    }

    export := map[string]json.RawMessage{
        "domain":            domainData,
        "applications":      results["applications"],
        "identityProviders": results["identityProviders"],
        "roles":             results["roles"],
        "scopes":            results["scopes"],
        "factors":           results["factors"],
        "groups":            results["groups"],
        "flows":             results["flows"],
        "certificates":      results["certificates"],
    }

    out, err := json.MarshalIndent(export, "", "  ")
    if err != nil {
        return err
    }

    if file != "" {
        return os.WriteFile(file, out, 0600)
    }
    fmt.Println(string(out))
    return nil
}
```

- [ ] **Step 2: Write `cmd/am/domain/import.go`**

Import reads JSON export file, creates (or targets) a domain, imports resources in dependency order: scopes → roles → groups → applications.

```go
package domain

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/am"
    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newImportCmd(f *factory.Factory) *cobra.Command {
    var targetDomainID string
    cmd := &cobra.Command{
        Use:     "import <file>",
        Short:   "Import domain configuration from a JSON export file",
        Example: `  gio am domain import domain-export.json
  gio am domain import domain-export.json --target existing-domain-id`,
        Args: cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
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

    svc := f.AM()
    p, _ := cmdutil.NewPrinter(f)

    if targetDomainID == "" {
        var domainObj map[string]any
        if err := json.Unmarshal(exportData["domain"], &domainObj); err != nil {
            return fmt.Errorf("failed to parse domain in export: %w", err)
        }
        created, err := svc.CreateDomain(am.CreateDomainBody{
            Name:        cmdutil.StringFieldMap(domainObj, "name"),
            Description: cmdutil.StringFieldMap(domainObj, "description"),
            DataPlaneID: cmdutil.StringFieldMap(domainObj, "dataPlaneId"),
        })
        if err != nil {
            return fmt.Errorf("failed to create domain: %w", err)
        }
        var newDomain map[string]any
        _ = json.Unmarshal(created, &newDomain)
        targetDomainID = cmdutil.StringFieldMap(newDomain, "id")
        p.PrintMessage("Created domain '%s'.", targetDomainID)
    }

    imported, skipped := 0, 0

    importList := func(key string, createFn func(item json.RawMessage) error) {
        var items []json.RawMessage
        if err := json.Unmarshal(exportData[key], &items); err != nil {
            return
        }
        for _, item := range items {
            if err := createFn(item); err != nil {
                skipped++
            } else {
                imported++
            }
        }
    }

    importList("scopes", func(item json.RawMessage) error {
        _, err := svc.CreateScope(targetDomainID, item)
        return err
    })
    importList("roles", func(item json.RawMessage) error {
        _, err := svc.CreateRole(targetDomainID, item)
        return err
    })
    importList("groups", func(item json.RawMessage) error {
        _, err := svc.CreateGroup(targetDomainID, item)
        return err
    })
    importList("applications", func(item json.RawMessage) error {
        _, err := svc.CreateApplication(targetDomainID, item)
        return err
    })

    p.PrintMessage("Import complete: %d imported, %d skipped.", imported, skipped)
    return nil
}
```

> **Note:** This requires `ScopeService.CreateScope` method from Task 6. Ensure `cmdutil.StringFieldMap` exists (look for existing helper or add it analogous to `StringField`).

- [ ] **Step 3: Write `cmd/am/domain/copy.go`**

Copy = export from source domain + import into new domain in same workspace.

```go
package domain

import (
    "github.com/spf13/cobra"

    "github.com/gravitee-io/gio-cli/internal/cmdutil"
    "github.com/gravitee-io/gio-cli/internal/factory"
)

func newCopyCmd(f *factory.Factory) *cobra.Command {
    var targetName string
    cmd := &cobra.Command{
        Use:     "copy <sourceDomainId>",
        Short:   "Copy a domain to a new domain in the same workspace",
        Example: `  gio am domain copy abc-123 --name my-copy`,
        Args:    cobra.ExactArgs(1),
        RunE: func(_ *cobra.Command, args []string) error {
            if err := cmdutil.RequireContext(f); err != nil {
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
    svc := f.AM()
    p, _ := cmdutil.NewPrinter(f)

    // Fetch source domain.
    srcDomain, err := svc.GetDomain(sourceDomainID)
    if err != nil {
        return err
    }

    // Create target domain.
    created, err := svc.CreateDomain(am.CreateDomainBody{Name: targetName})
    if err != nil {
        return err
    }
    var newDomain map[string]any
    _ = json.Unmarshal(created, &newDomain)
    targetDomainID := cmdutil.StringFieldMap(newDomain, "id")

    p.PrintMessage("Created domain '%s' (%s). Copying resources...", targetName, targetDomainID)

    // Reuse export/import pipeline.
    // Export source domain resources into memory, import into targetDomainID.
    // (Share runExport + runImport logic via internal helpers if they grow complex.)
    _ = srcDomain

    p.PrintMessage("Copy complete.")
    return nil
}
```

> **Note:** The full implementation of `runCopy` should reuse logic from `runExport` and `runImport`. If those functions are in the same package (`domain`), extract the core fetch/import logic into unexported helpers.

- [ ] **Step 4: Register in `cmd/am/domain/domain.go`**

```go
cmd.AddCommand(
    newListCmd(f),
    newGetCmd(f),
    newCreateCmd(f),
    newEnableCmd(f),
    newDisableCmd(f),
    newDeleteCmd(f),
    newExportCmd(f),
    newImportCmd(f),
    newCopyCmd(f),
)
```

- [ ] **Step 5: Build and test**

```bash
cd /Users/rpo/Documents/Projects/Gravitee/gio-cli && go build ./...
go test ./cmd/am/... ./internal/am/...
```

- [ ] **Step 6: Commit**

```bash
git add cmd/am/domain/
git commit -m "feat(am): add domain export, import, copy commands"
```

---

## Final checklist

After all tasks:

- [ ] `go build ./...` passes with no errors
- [ ] `go test ./...` passes
- [ ] `gio am --help` shows all registered sub-commands
- [ ] `gio am domain --help` shows: list, get, create, enable, disable, delete, export, import, copy
- [ ] `gio am application --help` shows: list, get, create, update, delete
- [ ] `gio am user --help` shows: list, get, create, update, delete, lock, unlock, reset-password
- [ ] `gio am idp --help` shows: list, get
- [ ] `gio am factor --help` shows: list, get
- [ ] `gio am role --help` shows: list, get, create
- [ ] `gio am scope --help` shows: list, get, create
- [ ] `gio am group --help` shows: list, get, create, delete
- [ ] `gio am flow --help` shows: list, get
- [ ] `gio am certificate --help` shows: list, get, delete
- [ ] `gio am audit --help` shows: list, get
- [ ] `gio am token --help` shows: list, create, revoke
- [ ] `gio am health` compiles
- [ ] `gio am whoami` compiles
