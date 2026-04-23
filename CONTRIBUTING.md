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

# Contributing to gio

Thanks for taking the time to contribute. This guide walks through reporting issues, setting up your environment, and the conventions we expect for changes and releases.

## Reporting a bug or suggesting an enhancement

Open an [issue](https://github.com/gravitee-io-labs/gio-cli/issues/new/choose) and follow the template that matches your request. Be sure to include the `gio` version (`gio version`), the Gravitee version you are targeting, and a minimal reproducer.

## Submitting a pull request

If the change addresses an existing issue, link it in the pull request description.

### Conventional commits are mandatory

All commit messages **must** follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) spec. GoReleaser parses commit prefixes to build the release changelog:

| Prefix               | Changelog section      |
|----------------------|------------------------|
| `feat:`              | Features               |
| `fix:`               | Bug fixes              |
| anything else        | Other changes          |
| `docs:`, `chore:`, `test:`, `ci:`, `style:`, `refactor:` | Excluded from changelog |

Use scopes when useful (`feat(apim): ...`, `fix(am): ...`). Breaking changes: append `!` (`feat!: ...`).

### Before pushing

Run the full lint locally so CI does not bounce your PR:

```bash
make lint            # golangci-lint + license headers
make add-license     # stamp headers on any new file
```

## Setting up your environment

### Prerequisites

- **Go** matching `go.mod` (currently 1.26+)
- **Docker** (required to run the end-to-end test suite, which spins up APIM/AM via `docker compose`)
- **Git** with the `v*` tag permission on the remote (for maintainers who cut releases)

### Install tooling

All release and lint tools are pinned Go modules under `hack/tools/` and installed into `bin/` (git-ignored). Install them once:

```bash
make install-tools
```

To bump a tool, edit the matching `hack/tools/<tool>/go.mod` and run:

```bash
make reinstall-tools
```

### Why pinned tool modules?

Every contributor and CI runs the exact same version of each tool (`addlicense`, `goreleaser`, ...). Bumps are tracked via `go.mod` diffs, integrity is guaranteed by `go.sum`, and there is no system-wide install required. See `hack/tools/addlicense/go.mod` and `hack/tools/goreleaser/go.mod` for examples.

## Building locally

```bash
make build     # builds ./dist/gio for your current platform
./dist/gio --help
```

Each `.mk` under `hack/make/` is self-sufficient: IntelliJ's gutter run and `make -f hack/make/build.mk build` both work without loading the root `Makefile` first.

## Testing

### Unit tests

```bash
make test          # runs go test ./...
make test-cover    # with coverage, generates cover.html
```

### End-to-end tests

End-to-end tests boot a full APIM stack via Docker Compose (`e2e/docker-compose.yml`) and exercise `gio` against it. They are gated behind the `e2e` build tag to stay out of `go test ./...`.

```bash
make e2e-up        # start the e2e infra in the background, wait for healthy
make test-e2e      # run the e2e tests against the running infra
make e2e-down      # stop the infra and remove volumes
```

CI uses the one-shot target which always tears down:

```bash
make e2e
```

## Linting and license headers

Every source and config file carries an Apache 2.0 header. The CI job `licenses` gates the build and fails on any missing header.

```bash
make lint-licenses   # check (fails if a file is missing its header)
make add-license     # stamp any missing header
```

The ignore list lives in `hack/make/lint.mk` (`.idea/`, `bin/`, `dist/`, `hack/tools/`, `hack/license.go.txt`).

## Release process

Releases are cut by pushing a Git tag that matches `v*`. GoReleaser takes over from there.

### Local validation (any time)

```bash
make release-check      # validate .goreleaser.yaml
make release-snapshot   # full cross-platform build + archives + checksums into dist/ (no publish)
```

The same snapshot runs automatically on every PR and merge via the `release-dry-run` job in `.github/workflows/ci.yml`, so release-time surprises surface before tagging.

### Pre-flight checks

Before tagging:

1. You are on the commit you want to ship (usually `main`, working tree clean).
2. CI is green on that commit, including the `release-dry-run` job.

### Cutting a real release

Create a Git tag matching `v*` (for example `v0.1.0`) on the target commit, either via `git tag` + `git push` or from **Releases** → **Draft a new release** in the GitHub UI. `.github/workflows/release.yml` picks up the tag, runs `goreleaser release --clean`, and publishes the Release with its grouped changelog, archives and `checksums.txt`.

If you use the UI, leave the release title and description empty: GoReleaser overwrites them.

If the workflow fails mid-run (transient GitHub API error for example), re-run it from the **Actions** tab: the publish step is idempotent. Only drop and re-tag if the code itself is broken.

### Pre-release (release candidate)

Tag with a pre-release suffix to publish a **pre-release** that is not marked as "Latest" on the repository home:

```bash
git tag v0.1.0-rc1
git push origin v0.1.0-rc1
```

GoReleaser detects the `-rc1` suffix automatically.

### Rolling back a broken release

```bash
gh release delete v0.1.0 --cleanup-tag --yes
```

This removes the GitHub Release and the Git tag.

## Project structure

```
.
|-- cmd/                 # Cobra command tree (login, apim, am, version, ...)
|-- internal/            # Non-public packages (client, config, cmdutil, apim, am, ...)
|-- e2e/                 # End-to-end test suite (tag `e2e`) + docker-compose stack
|-- hack/
|   |-- make/            # Included .mk files (build, test, lint, tool)
|   `-- tools/           # Pinned Go tool modules (addlicense, goreleaser, ...)
|-- LICENSE_TEMPLATE.txt # License header template used by addlicense
|-- .goreleaser.yaml     # Release pipeline config
|-- .github/workflows/   # CI (ci.yml), E2E (e2e.yml) and release (release.yml)
`-- LICENSE.txt
```

To add a new APIM or AM command, start from an existing one under `cmd/apim/` or `cmd/am/` and mirror the pattern (Cobra command, thin wrapper over an `internal/...` service).