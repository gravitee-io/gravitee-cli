<!--
Copyright (C) 2015 The Gravitee team (http://gravitee.io)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./.assets/gravitee-logo-dark.svg">
    <img src="./.assets/gravitee-logo-light.svg" alt="Gravitee" width="400">
  </picture>
</p>

# gio - the Gravitee CLI

`gio` is the official command-line interface for [Gravitee](https://www.gravitee.io/) **APIM** (API Management) and **AM** (Access Management). It lets platform teams, integrators and AI agents script, inspect and automate Gravitee control planes from any terminal.

## Install

See the [latest release](https://github.com/gravitee-io-labs/gio-cli/releases/latest) for archives and install instructions.

## Quickstart

### 1. Generate a service-account token (recommended)

For anything beyond local experimentation (CI, scripts, AI agents, production automation), create a **service account** in your Gravitee organization and generate a token for it.

### 2. Log in

`gio` targets APIM and AM independently. Log in to the product you need (both is fine):

```bash
gio login apim
gio login am
```

Both commands support three ways to authenticate.

#### Paste the `curl` from Gravitee (easiest)

When you generate a token in APIM or AM, the UI displays a ready-to-run `curl` command. Paste it **entirely** at the `URL` prompt and `gio` extracts the URL, token, organization and environment for you:

```
$ gio login apim
URL (or paste full curl command): curl -H 'Authorization: Bearer gioat_abc...' https://apim.example.com/management/v2/organizations/DEFAULT/environments/DEFAULT/apis
```

No manual copy/paste of each field, no format mistakes.

#### Flags (for scripts)

```bash
gio login apim \
  --url https://apim.example.com \
  --token gioat_abc... \
  --org DEFAULT \
  --env DEFAULT \
  --context production
```

#### Environment variables (for CI)

When these variables are set, `gio` bypasses the config file entirely:

| Variable         | Scope  |
|------------------|--------|
| `GIO_APIM_URL`   | APIM   |
| `GIO_APIM_TOKEN` | APIM   |
| `GIO_AM_URL`     | AM     |
| `GIO_AM_TOKEN`   | AM     |
| `GIO_ORG`        | shared |
| `GIO_ENV`        | shared |

### 3. Run your first command

```bash
gio apim api list
gio am domain list
```

## Core commands

```bash
gio --help         # top-level overview
gio apim --help    # APIs, applications, plans, subscriptions, ...
gio am --help      # domains, applications, users, identity providers, ...
```

Global flags worth knowing:
- `--context <name>` : pick a named context from your config file
- `--output table|yaml|json|id` : control output format. `yaml` and `json` for machines, `id` for piping into scripts, `table` (default) for humans.

## AM (Access Management)

Full-featured CLI for managing Gravitee Access Management domains and security:

**Domain and resource management** — Full CRUD for domains, applications, users, identity providers, roles, scopes, and certificates. Additional operations: enable/disable domains, manage application settings, lock/unlock users, reset passwords, and switch between domains with `gio am set domain <name>`.

**Security and auditing** — Query audit logs with filtering by type, status, and date range. Run built-in security rule validation (`gio am lint`) to score domain security posture. Export diagnostic dumps with optional secret redaction for compliance reporting.

**Authentication and tokens** — Manage token lifecycle (create, list, revoke), lock/unlock users, reset passwords, and test OpenID Connect discovery and client credential flows with integrated OIDC testing. Bearer tokens are masked in list output; passwords and client secrets accept `--password-stdin` / `--secret-stdin` to avoid leaking via process listings or shell history.

**Monitoring and diagnostics** — Stream logs in real-time with `--follow`, watch live domain activity on a dashboard, trace authentication flows step-by-step to diagnose issues, and run health checks.

**Backup and migration** — Export domain configurations for backup, import into another domain or environment, or copy between contexts with `gio am diff` to preview changes.

**Scripting and exploration** — Generate shell completions, use an interactive REPL shell to explore resources, and inspect plugin schemas.

## Configuration

After a successful login, `gio` persists your contexts in `~/.gio/config.yaml` with file mode `0600` (owner read/write only). Example:

```yaml
current: default
contexts:
  default:
    org: DEFAULT
    env: DEFAULT
    apim:
      url: https://apim.example.com
      token: gioat_abc...
    am:
      url: https://am.example.com
      token: gioat_def...
```

Switch between contexts with `--context <name>` on any command, or by editing the `current` field.

## AI-agent friendly

`gio` is designed to be script- and agent-friendly:

- **YAML output** via `--output yaml` for deterministic, diff-able responses
- **Terse defaults**, minimal noise to keep small context windows usable
- **Actionable error messages** with hints when a call fails
- **Clean exit codes** so automation can branch on failure

## Documentation and support

- [Gravitee product documentation](https://documentation.gravitee.io/)
- [Report an issue](https://github.com/gravitee-io-labs/gio-cli/issues)
- [Latest release and changelog](https://github.com/gravitee-io-labs/gio-cli/releases/latest)

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for how to build, test, and release `gio`.

## License

Apache License 2.0. See [LICENSE.txt](./LICENSE.txt).
