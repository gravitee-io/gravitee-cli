# AM CLI Implementation

This document describes the Gravitee Access Management (AM) CLI module implementation. The `am` branch adds comprehensive command-line interface for managing AM domains, applications, users, and security configurations.

## Overview

The AM CLI provides a complete set of tools for managing Gravitee Access Management domains and resources. It follows the same patterns as the existing APIM CLI, with support for multiple contexts, output formatting, and domain management.

## Architecture

- **Module**: `cmd/am/` - Main AM command package
- **Service Layer**: `internal/am/` - AM service abstraction
- **Config**: Extended `config.Context` with `Type: "am"` and `Domain` fields for domain-scoped operations
- **Authentication**: Username/password and token-based authentication
- **Output**: Table, JSON, and YAML output formats via `internal/printer`

## Implemented Commands

### Foundation (Commit: feat: add AM module foundation)

Established core infrastructure:
- Config extensions: `Type` and `Domain` fields for context
- Path helpers: `AMEnvPath`, `AMDomainPath`, `AMPath`
- Context validators: `RequireAMContext`, `RequireAMDomain`
- HTTP client `Patch` method for partial updates
- AM parent command registration

### Core CRUD Operations (Commit: feat: implement AM CRUD commands)

Complete CRUD for AM resources:
- **Domains**: list, get, create, update, delete, enable, disable, set (current)
- **Applications**: list, get, create, update, delete, with settings management
- **Users**: list, get, create, update, delete, lock, unlock, reset-password
- **Identity Providers**: list, get, create, update, delete
- **Roles**: list, get, create, update, delete
- **Scopes**: list, get, create, update, delete
- **Certificates**: list, get, create, update, delete

Each command supports domain-scoped queries and provides appropriate error handling.

### Operational Commands (Commit: feat: add AM operational commands)

Read-only and state-change operations:
- **Audit**: list, get - audit log queries with filtering (type, status, date range)
- **Health**: System health and status checks
- **WhoAmI**: Current authentication context
- **Status**: Domain status and configuration verification
- **Doctor**: Health and connectivity diagnostics
- **Token**: list, create, revoke - token lifecycle management
- **Factor**, **Flow**, **Group**: List and get operations
- **Logout**: Session termination with config reset
- **Domain Export/Import/Copy**: Backup and replication operations

### Advanced Features (Commit: feat: add AM advanced commands)

Specialized operational tools:
- **Logs**: Real-time log streaming with `--follow` support
- **Watch**: Live dashboard for domain activity
- **Diff**: Compare domain configurations across contexts
- **Lint**: Security rule validation (14 built-in rules) with scoring
- **Trace**: Authentication path analysis with step-by-step diagnostics
- **Support-Dump**: Diagnostic data export (optional redaction for secrets)
- **Completion**: Shell completion script generation
- **Shell**: Interactive REPL for exploring AM resources
- **OIDC Test**: OpenID Connect discovery and client credential flow testing
- **Plugins**: Plugin discovery and schema inspection

### Test Coverage (Commit: test: add command-level coverage)

Command-level integration tests for:
- **Audit**: Event parsing, formatting, column rendering, pagination (totalCount and incomplete-page termination)
- **Trace**: All 7 check functions (grant types, MFA, flows, consent, token config) + full e2e with fake client
- **Support-Dump**: Single-domain dump with explicit API routing and error detection
- **Lint**: Critical finding detection and score output
- **Diff**: Cross-context diff with httptest servers, scope additions detection

Test helpers use `factory.Factory` with `config.ProductConfig` for AM service mocking. Tests verify both happy paths and edge cases.

## Configuration

Contexts for AM operations require:
```yaml
contexts:
  my-am:
    url: https://am.example.com
    token: <personal-access-token>
    type: am
    domain: my-domain  # AM domain ID
```

## Usage Examples

```bash
# Login and set domain
gio am login -u admin -p password
gio am set domain my-domain

# Audit trail
gio am audit list --type USER_LOGIN --status SUCCESS

# Trace authentication
gio am trace --user john --app myapp

# Security linting
gio am lint

# Compare environments
gio am diff --from staging --to production

# Export for backup
gio am support-dump --domain prod-domain > backup.json
```

## Implementation Notes

- All commands require a domain context (`RequireAMDomain`)
- Pagination defaults to 10 items; use `--all` to fetch all results
- Output format defaults to table; use `--format json` for scripting
- Error handling includes helpful diagnostics (connect failures, missing domains, invalid tokens)
- Concurrent API calls where possible (e.g., trace diagnostic checks)
- YAML-based config storage with automatic validation

## Testing Strategy

Tests use:
- **FakeClient**: In-process mock for fast unit tests
- **Factory Helper**: Standard `newTestFactory` pattern across all packages
- **HTTPTest Servers**: For commands creating independent HTTP clients (diff)
- **Coverage**: Both happy paths and error conditions
- **Edge Cases**: Pagination termination, incomplete pages, missing optional data

## Future Enhancements

- Batch operations (apply multiple changes atomically)
- Policy and flow templating
- Advanced filtering and query syntax
- Custom output field selection
- Webhook integration for events
