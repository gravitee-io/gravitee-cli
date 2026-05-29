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

# gctl - the Gravitee CLI

`gctl` is the official command-line interface for [Gravitee](https://www.gravitee.io/) **APIM** (API Management) and **AM** (Access Management). It lets platform teams, integrators and AI agents script, inspect and automate Gravitee control planes from any terminal.

## Variants

Two binaries are available:

| Binary | Use case | Commands |
|--------|----------|----------|
| `gctl` | Humans, scripts, CI | Full CRUD - get, list, create, update, delete |
| `gctl-ro` | AI agents, read-only automation | Get and list only - no write operations |

Both read the same config file (`~/.gctl/config.yaml`) and support the same flags and output formats. Install both on the same machine if you want agents to use `gctl-ro` while you use `gctl` - they coexist without conflict.

## Install

### curl (macOS / Linux, and Windows via Git Bash / WSL)

```bash
curl -fsSL https://raw.githubusercontent.com/gravitee-io/gravitee-cli/main/install.sh | sh
```

Install the read-only variant (`gctl-ro`) instead:

```bash
curl -fsSL https://raw.githubusercontent.com/gravitee-io/gravitee-cli/main/install.sh | GCTL_BIN=gctl-ro sh
```

### Manual (any OS, including native Windows)

Download the archive for your platform from the [latest release](https://github.com/gravitee-io/gravitee-cli/releases/latest), extract it, and move `gctl` (and optionally `gctl-ro`) into a directory on your `PATH`. Run `gctl version` to verify.

## Quickstart

### 1. Generate a service-account token (recommended)

For anything beyond local experimentation (CI, scripts, AI agents, production automation), create a **service account** in your Gravitee organization and generate a token for it.

### 2. Log in

`gctl` targets APIM and AM independently. Log in to the product you need (both is fine):

```bash
gctl login apim
gctl login am
```

Both commands support three ways to authenticate.

#### Paste the `curl` from Gravitee (easiest)

When you generate a token in APIM or AM, the UI displays a ready-to-run `curl` command. Paste it **entirely** at the `URL` prompt and `gctl` extracts the URL, token, organization and environment for you:

```
$ gctl login apim
URL (or paste full curl command): curl -H 'Authorization: Bearer gioat_abc...' https://apim.example.com/management/v2/organizations/DEFAULT/environments/DEFAULT/apis
```

No manual copy/paste of each field, no format mistakes.

#### Flags (for scripts)

```bash
gctl login apim \
  --url https://apim.example.com \
  --token gioat_abc... \
  --org DEFAULT \
  --env DEFAULT \
  --context production
```

#### Environment variables (for CI)

When these variables are set, `gctl` bypasses the config file entirely:

| Variable         | Scope  |
|------------------|--------|
| `GCTL_APIM_URL`   | APIM   |
| `GCTL_APIM_TOKEN` | APIM   |
| `GCTL_AM_URL`     | AM     |
| `GCTL_AM_TOKEN`   | AM     |
| `GCTL_ORG`        | shared |
| `GCTL_ENV`        | shared |

### 3. Run your first command

```bash
gctl apim api list
gctl am domain list
```

## Core commands

```bash
gctl --help         # top-level overview
gctl apim --help    # APIs, applications, plans, subscriptions, ...
gctl am --help      # domains, applications, users, identity providers, ...
```

Global flags worth knowing:
- `--context <name>` : pick a named context from your config file
- `--output table|yaml|json|id` : control output format. `yaml` and `json` for machines, `id` for piping into scripts, `table` (default) for humans.

List commands share a common set of flags:
- `--all` : fetch all pages automatically instead of a single page
- `--per-page <n>` : results per page (default 10)
- `--wide` / `-w` : show extra columns (tags, owner, visibility, ...) - available on `gctl apim api list`

## AM (Access Management)

Full-featured CLI for managing Gravitee Access Management domains and security:

**Domain and resource management** — Full CRUD for domains, applications, users, identity providers, roles, scopes, and certificates. Additional operations: enable/disable domains, manage application settings, lock/unlock users, reset passwords, and switch between domains with `gctl am set domain <name>`.

**Security and auditing** — Query audit logs with filtering by type, status, and date range. Run built-in security rule validation (`gctl am lint`) to score domain security posture. Export diagnostic dumps with optional secret redaction for compliance reporting.

**Authentication and tokens** — Manage token lifecycle (create, list, revoke), lock/unlock users, reset passwords, and test OpenID Connect discovery and client credential flows with integrated OIDC testing. Bearer tokens are masked in list output; passwords and client secrets accept `--password-stdin` / `--secret-stdin` to avoid leaking via process listings or shell history.

**Monitoring and diagnostics** — Stream logs in real-time with `--follow`, watch live domain activity on a dashboard, trace authentication flows step-by-step to diagnose issues, and run health checks.

**Backup and migration** — Export domain configurations for backup, import into another domain or environment, or copy between contexts with `gctl am diff` to preview changes.

**Scripting and exploration** — Generate shell completions, use an interactive REPL shell to explore resources, and inspect plugin schemas.

## Configuration

After a successful login, `gctl` persists your contexts in `~/.gctl/config.yaml` with file mode `0600` (owner read/write only). Example:

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

`gctl` is designed to be script- and agent-friendly:

- **YAML output** via `--output yaml` for deterministic, diff-able responses
- **Terse defaults**, minimal noise to keep small context windows usable
- **Actionable error messages** with hints when a call fails
- **Clean exit codes** so automation can branch on failure
- **`gctl-ro`** for agents that only need to read - structurally prevents any write operation, no prompt engineering required

## Documentation and support

- [Gravitee product documentation](https://documentation.gravitee.io/)
- [Report an issue](https://github.com/gravitee-io/gravitee-cli/issues)
- [Latest release and changelog](https://github.com/gravitee-io/gravitee-cli/releases/latest)

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for how to build, test, and release `gctl`.

## License

Apache License 2.0. See [LICENSE.txt](./LICENSE.txt).
