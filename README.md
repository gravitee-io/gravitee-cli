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

`gctl` is the official command-line interface for [Gravitee](https://www.gravitee.io/) **APIM** (API Management) and **AM** (Access Management). It lets humans, CI pipelines and AI agents script, inspect and automate Gravitee from the command line.

## Variants

`gctl` ships as two binaries:

| Binary | Use case | Commands |
|--------|----------|----------|
| `gctl` | Humans, scripts, CI | Full CRUD - get, list, create, update, delete |
| `gctl-ro` | AI agents, read-only automation | Get and list only - no write operations |

Both read the same config file (`~/.gctl/config.yaml`) and share the same flags and output formats. Install them side by side: agents use `gctl-ro`, you use `gctl`.

## Install

### curl (macOS / Linux)

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

For any real use (CI, scripts, agents, production), create a **service account** in your Gravitee organization and generate a token for it.

### 2. Log in

`gctl` targets APIM and AM independently. Log in to whichever you need, or both:

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

Global flags:
- `--context <name>` : pick a named context from your config file
- `--output table|yaml|json|id` : control output format. `yaml` and `json` for machines, `id` for piping into scripts, `table` (default) for humans.

List commands share these flags:
- `--all` : fetch all pages automatically instead of a single page
- `--per-page <n>` : results per page (default 10)
- `--wide` / `-w` : show extra columns (tags, owner, visibility, ...) - available on `gctl apim api list`

## AM (Access Management)

`gctl am` manages Gravitee Access Management domains and security:

- **Resources**: full CRUD for domains, applications, users, identity providers, roles, scopes, and certificates, plus enable/disable, lock/unlock, and reset-password operations.
- **Security and audit**: query audit logs, score domain posture with `gctl am lint`, and export diagnostic dumps with optional secret redaction.
- **Auth and tokens**: manage token lifecycle and test OIDC discovery and client-credential flows. Secrets accept `--password-stdin` / `--secret-stdin` to stay out of shell history.
- **Diagnostics**: stream logs with `--follow`, watch live activity, trace auth flows, and run health checks.
- **Backup and migration**: export, import, or `gctl am diff` between contexts.

See `gctl am --help` for the full command tree.

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

`gctl` is built for scripts and agents, not just interactive use:

- **Machine-readable output**: `--output yaml|json` for parsing, `--output id` for piping into other commands
- **Quiet by default**: terse output keeps agent context windows small
- **Actionable errors**: failures explain what went wrong, and often how to fix it
- **Predictable exit codes**: automation can branch on success or failure
- **`gctl-ro`**: a read-only binary with no write commands. If the CLI is the agent's only access, it cannot modify anything, period. It is not a sandbox, though: give the agent a shell and the config token and it can call the API directly. The guarantee is only as strong as the tools you hand it, so back it with a read-only service-account token.

## Documentation and support

- [Gravitee product documentation](https://documentation.gravitee.io/)
- [Report an issue](https://github.com/gravitee-io/gravitee-cli/issues)
- [Latest release and changelog](https://github.com/gravitee-io/gravitee-cli/releases/latest)

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for how to build, test, and release `gctl`.

## License

Apache License 2.0. See [LICENSE.txt](./LICENSE.txt).
