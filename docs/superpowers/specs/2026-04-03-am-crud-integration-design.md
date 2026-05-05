# AM CRUD Integration into gio-cli

**Date:** 2026-04-03
**Status:** Approved
**Scope:** Tier 1 CRUD — domain, app (+ settings), user, idp, role, scope, certificate

---

## Context

The `gio` CLI currently covers Gravitee API Management (APIM). This design adds Access Management (AM) support as `gio am <resource> <operation>`, migrating CRUD functionality from the existing TypeScript am-cli (v0.1) into Go.

**Source of truth for AM API:** Java source in `gravitee-am-management-api` (JAX-RS resources).
**Source of truth for UX:** TypeScript am-cli at `am-tooling/am-cli` (Commander.js).

---

## Decision Log

| Decision | Choice | Rationale |
|---|---|---|
| Domain context | Part of unified config context (field `type` + `domain`) | Consistent with how APIM handles org/env |
| Command hierarchy | Flat: `gio am domain list`, `gio am user list` | Matches APIM convention, simpler UX |
| Folder structure | `cmd/am/` with sub-folders per resource | Clean separation from APIM commands |
| Infrastructure reuse | Same `GraviteeClient`, config, factory; add AM path helpers + Patch | DRY, no duplication |
| Create/Update UX | Inline flags + interactive prompts + `-f <file>` | 1:1 with am-cli v0.1 |
| Login flow | AM-style: username/password -> JWT, or direct token | AM has no PATs; matches am-cli v0.1 |
| Config | Unified contexts with `type` field, backward compatible | Single config file, one `gio config` |
| Pagination | 0-based `--page`, `--size`, `--all`, `-q` | Matches AM API (0-based) |
| Approach | Bottom-up: infra -> domain (reference) -> remaining 6 resources | Proven pattern, early e2e validation |
| Interactive prompts | `github.com/AlecAivazis/survey/v2` | Closest to inquirer.js from am-cli |

---

## Architecture

### Config Changes (`internal/config/config.go`)

```go
type Context struct {
    URL      string `json:"url"`
    Token    string `json:"token"`
    Org      string `json:"org,omitempty"`
    Env      string `json:"env,omitempty"`
    ReadOnly bool   `json:"readOnly,omitempty"`
    Type     string `json:"type,omitempty"`     // "apim" (default) or "am"
    Domain   string `json:"domain,omitempty"`   // AM domain ID or HRID
}

type ResolvedContext struct {
    Name     string
    URL      string
    Token    string
    Org      string
    Env      string
    ReadOnly bool
    Type     string
    Domain   string
}
```

Backward compatibility: contexts without `type` are treated as `"apim"`. Zero breaking changes.

### Client Changes (`internal/client/`)

Add `Patch` method to `GraviteeClient` interface:

```go
type GraviteeClient interface {
    Get(path string) ([]byte, error)
    Post(path string, body interface{}) ([]byte, error)
    Put(path string, body interface{}) ([]byte, error)
    Patch(path string, body interface{}) ([]byte, error)  // new
    Delete(path string) error
}
```

Update `HTTPClient` and `FakeClient` accordingly.

### Path Helpers (`internal/cmdutil/`)

```go
// AMEnvPath builds: /management/organizations/{org}/environments/{env}/{path}
// Used for domain-level operations (list/create domains).
func AMEnvPath(f *factory.Factory, path string) string

// AMDomainPath builds: /management/organizations/{org}/environments/{env}/domains/{domain}/{path}
// Used for resources under a domain (users, apps, roles, etc.).
func AMDomainPath(f *factory.Factory, path string) string

// RequireAMContext validates that the active context is type=am.
func RequireAMContext(f *factory.Factory) error

// RequireAMDomain validates type=am AND domain is set.
func RequireAMDomain(f *factory.Factory) error
```

---

## Command Structure

### Folder Layout

```
cmd/am/
    am.go                    # NewAMCmd — parent "gio am"
    login.go                 # gio am login
    login_test.go
    set.go                   # gio am set domain <id>
    set_test.go
    domain/
        domain.go            # NewDomainCmd
        list.go
        get.go
        create.go
        update.go
        delete.go
        enable.go
        disable.go
        helpers_test.go
        list_test.go
        get_test.go
        create_test.go
        ...
    app/
        app.go
        list.go
        get.go
        create.go
        update.go
        delete.go
        settings.go          # view/update OAuth2 settings
        helpers_test.go
        ...
    user/
        user.go
        list.go
        get.go
        create.go
        update.go
        delete.go
        lock.go
        unlock.go
        reset_password.go
        helpers_test.go
        ...
    idp/
        idp.go
        list.go, get.go, create.go, update.go, delete.go
        helpers_test.go
    role/
        role.go
        list.go, get.go, create.go, update.go, delete.go
        helpers_test.go
    scope/
        scope.go
        list.go, get.go, create.go, update.go, delete.go
        helpers_test.go
    certificate/
        certificate.go
        list.go, get.go, create.go, update.go, delete.go
        helpers_test.go
```

Registration in `cmd/root.go`:

```go
import amcmd "github.com/gravitee-io/gio-cli/cmd/am"
// ...
cmd.AddCommand(amcmd.NewAMCmd(f))
```

---

## Login & Auth

### `gio am login`

**Interactive:**
```
gio am login
-> URL: http://localhost:8093
-> Username: admin
-> Password: ****
Context 'localhost-am' saved and set as current.
```

**Non-interactive:**
```
gio am login --url http://localhost:8093 --username admin --password admin
gio am login --url http://localhost:8093 --token eyJhbG...
```

**Mechanism:**
- Credentials: POST `/management/auth/token` with `Authorization: Basic base64(user:pass)`
- Returns JWT token, saved in context with `type: "am"`
- Token mode (`--token`): skip POST, save directly
- Context auto-named from hostname (e.g. `localhost-am`), `--context` overrides

### `gio am set domain <idOrHrid>`

- Fetches domain list, matches by ID or HRID
- Saves `domain` field in active AM context
- `gio am set domain --clear` removes it

### Environment Variables (CI/CD)

| Variable | Overrides |
|---|---|
| `AM_URL` | context URL |
| `AM_TOKEN` | context token |
| `AM_DOMAIN` | context domain |
| `AM_ORG` | context org (default: DEFAULT) |
| `AM_ENV` | context env (default: DEFAULT) |

---

## CRUD Operations

### Standard Pattern

| Operation | HTTP | Path | Input |
|---|---|---|---|
| list | GET | `/{resource}` | `--page`, `--size`, `--all`, `-q` |
| get | GET | `/{resource}/{id}` | positional `<id>` |
| create | POST | `/{resource}` | inline flags OR interactive prompts OR `-f <file>` |
| update | PUT/PATCH | `/{resource}/{id}` | inline flags OR `-f <file>` |
| delete | DELETE | `/{resource}/{id}` | positional `<id>`, `--force` |

**Create/Update modes (priority order):**
1. Inline flags provided -> use them, non-interactive
2. No required flags -> interactive prompts (survey/v2)
3. `-f <file>` -> read JSON from file

### API Paths

| Resource | Base path | Requires domain? |
|---|---|---|
| domain | `AMEnvPath("domains")` | no |
| app | `AMDomainPath("applications")` | yes |
| user | `AMDomainPath("users")` | yes |
| idp | `AMDomainPath("identities")` | yes |
| role | `AMDomainPath("roles")` | yes |
| scope | `AMDomainPath("scopes")` | yes |
| certificate | `AMDomainPath("certificates")` | yes |

### Special Operations

| Command | Method | Endpoint |
|---|---|---|
| `domain enable <id>` | PATCH | `/domains/{id}` body: `{"enabled": true}` |
| `domain disable <id>` | PATCH | `/domains/{id}` body: `{"enabled": false}` |
| `user lock <id>` | POST | `/users/{id}/lock` |
| `user unlock <id>` | POST | `/users/{id}/unlock` |
| `user reset-password <id>` | POST | `/users/{id}/resetPassword` body: `{"password": "..."}` |
| `app settings <id>` | GET/PATCH | `/applications/{id}` (view or update OAuth2 settings) |

### Inline Flags per Resource

**domain create:**
- `--name` (required/prompted), `--description`, `--data-plane-id` (default: "default")

**app create:**
- `--name` (required/prompted), `--type` (prompted: web/native/browser/service/resource_server), `--description`, `--redirect-uris`, `--idp`

**app settings:**
- `--grant-types`, `--response-types`, `--redirect-uris`, `--post-logout-uris`, `--token-lifetime`, `--refresh-token-lifetime`, `--id-token-lifetime`, `--enhance-scopes`

**user create:**
- `--username` (required/prompted), `--email`, `--firstName`, `--lastName`, `--password` (prompted if missing), `--preRegistration`

**role create:**
- `--name` (required/prompted), `--type` (prompted: DOMAIN/APPLICATION), `--description`

**scope create:**
- `--key` (required/prompted), `--name` (required/prompted), `--description`

**idp/certificate create:**
- `-f <file>` (required) — these have complex plugin configuration schemas

### Table Columns (list output)

| Resource | Columns |
|---|---|
| domain | name, hrid, id, enabled, description |
| app | name, type, clientId, id, enabled, description |
| user | username, email, firstName, lastName, id, enabled, accountNonLocked |
| idp | name, type, id, external |
| role | name, assignableType, id, description |
| scope | key, name, id, description |
| certificate | name, type, id, status |

### Pagination

AM API uses 0-based pagination. Response format:
```json
{"data": [...], "currentPage": 0, "totalCount": N}
```

Flags: `--page` (0-based, default 0), `--size` (default 20), `--all` (fetch all pages), `-q` (search query).

`--all` fetches pages in a loop with size=100, safety limit at page 1000.

---

## Testing

### Pattern (same as APIM)

- `FakeClient` with injectable functions (add `PatchFunc`)
- `helpers_test.go` per AM package with `newTestFactory` setting `type: "am"`, `domain: "test-domain"`
- Coverage: success paths, error responses, read-only mode, output formats

### Scope

- Non-interactive paths tested via flags (interactive prompts tested manually)
- `RequireAMContext` / `RequireAMDomain` validation
- AM path helpers unit tests
- Login: JWT token save, direct token save
- Pagination: 0-based query params
- Config: backward compat (no type = apim), type=am with domain

---

## Implementation Order (Bottom-Up)

1. **Infrastructure** — config type/domain fields, Patch on client, AM path helpers, RequireAMContext/RequireAMDomain
2. **Login & set domain** — `gio am login`, `gio am set domain`
3. **Domain CRUD** — full reference implementation (list/get/create/update/delete/enable/disable)
4. **App CRUD + settings** — list/get/create/update/delete/settings
5. **User CRUD + actions** — list/get/create/update/delete/lock/unlock/reset-password
6. **IdP CRUD** — list/get/create/update/delete
7. **Role CRUD** — list/get/create/update/delete
8. **Scope CRUD** — list/get/create/update/delete
9. **Certificate CRUD** — list/get/create/update/delete

Each step includes tests.

---

## Out of Scope (Tier 2+)

- Domain export/import/copy
- OIDC testing (test discover, test login, test client-credentials)
- Diagnostics (health, doctor, support-dump, logs)
- Interactive shell mode
- Smart resolution (name/email -> UUID lookup)
- group, factor, flow, audit, form, email resources
- bot-detection, device-identifier, password-policy, theme, i18n, extension-grant, reporter, resource, auth-device-notifier, authorization-engine
