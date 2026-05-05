# AM CLI Parity Gaps — v0.1 (TypeScript) vs gio-cli (Go)

**Date:** 2026-04-03
**Status:** Gap analysis — ready for prioritization

---

## Current State

Go gio-cli implements **core CRUD** for 7 resources + login + set domain.
TypeScript am-cli v0.1 has **27 command families** covering CRUD, diagnostics, testing, analysis, and utilities.

---

## Tier A — Critical for adoption (session & config management)

These are blockers — without them, daily use is painful.

| # | Feature | TS Command | What it does | Effort |
|---|---------|-----------|--------------|--------|
| A1 | **logout** | `am logout [--all]` | Clear token for current workspace or all | S |
| A2 | **whoami** | `am whoami` | GET `/user` — show authenticated user info (username, email, roles) | S |
| A3 | **status** | `am status` | Display current context: workspace, URL, domain, auth status, token expiry | S |
| A4 | **config management** | `am config *` | set-workspace, use-workspace, delete-workspace, list, current, set-default-output, path | M |
| A5 | **env var overrides** | (implicit) | AM_URL, AM_TOKEN, AM_DOMAIN, AM_ORG, AM_ENV override config — needed for CI/CD | M |
| A6 | **interactive login** | `am login` (no flags) | Prompt for URL, username, password when flags missing | S |
| A7 | **interactive delete confirmation** | `am * delete` | Use survey.Confirm instead of just printing "--force hint" | S |

---

## Tier B — Important features (diagnostics & monitoring)

| # | Feature | TS Command | What it does | Effort |
|---|---------|-----------|--------------|--------|
| B1 | **health** | `am health [--gateway <url>]` | Check AM Management API reachability, optionally gateway | S |
| B2 | **doctor** | `am doctor` | Full diagnostic: config file, workspace, auth, domain, connectivity, env vars | M |
| B3 | **audit list** | `am audit list [--type] [--status] [--from] [--to]` | List audit events with filtering | M |
| B4 | **audit get** | `am audit get <id>` | Get single audit event | S |
| B5 | **logs** | `am logs [-f] [--type] [--status]` | Tail audit logs in real-time with polling | L |

---

## Tier C — Remaining resources (CRUD for factor, flow, group)

| # | Feature | TS Command | Operations | Effort |
|---|---------|-----------|-----------|--------|
| C1 | **factor** | `am factor *` | list, get (MFA factors) | S |
| C2 | **flow** | `am flow *` | list, get (auth flows) | S |
| C3 | **group** | `am group *` | list, get, create, delete | M |

---

## Tier D — Advanced domain operations

| # | Feature | TS Command | What it does | Effort |
|---|---------|-----------|--------------|--------|
| D1 | **domain export** | `am domain export [-f file]` | Export full domain config (apps, idps, certs, roles, scopes, factors, groups, flows) to JSON | L |
| D2 | **domain import** | `am domain import <file> [--target] [--dry-run]` | Import domain from export file, create resources in order (scopes→roles→groups→apps) | L |
| D3 | **domain copy** | `am domain copy --to <workspace> [--dry-run]` | Copy domain across workspaces (export from source, import into target) | L |

---

## Tier E — OIDC Testing

| # | Feature | TS Command | What it does | Effort |
|---|---------|-----------|--------------|--------|
| E1 | **test discover** | `am test discover` | Fetch and display OIDC discovery document from gateway | S |
| E2 | **test login (ROPC)** | `am test login --app <id> --username <u> --password <p>` | Test password grant flow, decode JWT, validate issuer/expiry | M |
| E3 | **test client-credentials** | `am test client-credentials --app <id> --secret <s>` | Test client_credentials grant, decode token | M |

---

## Tier F — Analysis & validation tools

| # | Feature | TS Command | What it does | Effort |
|---|---------|-----------|--------------|--------|
| F1 | **diff** | `am diff --workspace <target>` | Compare domain config between workspaces (apps, idps, roles, scopes, factors, certs, flows, groups) | L |
| F2 | **lint** | `am lint [--ci]` | Security audit with 13 rules (implicit grant, PKCE, token lifetime, cert expiry, redirect URIs, etc.) | L |
| F3 | **trace** | `am trace --app <id> --user <id>` | Trace auth path for user+app combination — shows IdP chain, MFA factors, flows | L |

---

## Tier G — Utilities

| # | Feature | TS Command | What it does | Effort |
|---|---------|-----------|--------------|--------|
| G1 | **shell** | `am shell` | Interactive REPL with tab completion, command history | XL |
| G2 | **completion** | `am completion bash\|zsh\|fish` | Generate shell completions (cobra has this built-in!) | S |
| G3 | **support-dump** | `am support-dump [-f file] [--all-domains] [--no-redact]` | Generate diagnostic dump (config, domains, apps, audits) with secret redaction | L |
| G4 | **token** | `am token create/list/revoke` | Manage service account tokens for users | M |
| G5 | **plugin** | `am plugin list <type> \| schema <type> <id> \| create <type> <id>` | Plugin discovery, schema inspection, interactive creation from schema | L |
| G6 | **watch** | `am watch [--interval]` | Live dashboard with event stats, success rates, top errors | L |

---

## Tier H — Cross-cutting quality improvements

These apply across all existing commands.

| # | Feature | What it does | Where | Effort |
|---|---------|--------------|-------|--------|
| H1 | **Smart ID resolution** | Accept name/HRID/email besides UUID for get/update/delete | app get, user get, domain set, etc. | M |
| H2 | **Retry with backoff** | Auto-retry on 502/503/504 + timeouts | internal/client | M |
| H3 | **Token expiry warning** | Warn 5 min before token expires, error on expired | internal/config or middleware | S |
| H4 | **Spinners** | Show progress spinner during long operations | create, delete, export, import | S |
| H5 | **Compact output format** | Add `compact` output alongside table/json/yaml | internal/printer | S |

---

## Implementation Priority Recommendation

### Phase 1 — Make it usable daily (Tier A)
A1-A7: logout, whoami, status, config, env vars, interactive login, interactive delete

### Phase 2 — Diagnostics & monitoring (Tier B + C)
B1-B5: health, doctor, audit, logs
C1-C3: factor, flow, group resources

### Phase 3 — Domain operations & testing (Tier D + E)
D1-D3: domain export/import/copy
E1-E3: OIDC testing

### Phase 4 — Analysis & power tools (Tier F + G)
F1-F3: diff, lint, trace
G1-G6: shell, completion, support-dump, token, plugin, watch

### Phase 5 — Polish (Tier H)
H1-H5: smart resolution, retry, token expiry, spinners, compact output

---

## Effort Legend

| Size | Description | Estimate |
|------|------------|----------|
| S | Single file, < 100 lines | 1 task |
| M | 2-3 files, some logic | 2-3 tasks |
| L | Multi-file, complex logic | 4-6 tasks |
| XL | Major feature, new architecture | 8+ tasks |
