# Releasing

## 1. Release Model

Releases are published by `.github/workflows/release.yml`. The workflow can run
from a matching tag or from `workflow_dispatch`, and it invokes GoReleaser
through the Nix `cd` development shell.

Prefer tag-based releases. Use the manual workflow only when the intended
release ref and GoReleaser inputs are clear.

## 2. Pre-Release Checklist

Start from the commit that should be released and verify the worktree is clean:

```sh
git status --short
```

Run the local checks:

```sh
nix fmt
nix flake check --print-build-logs
nix build --print-build-logs
nix develop .#ci --command go test ./...
nix develop .#ci --command govulncheck ./...
nix develop .#ci --command gitleaks detect --verbose --redact
nix develop .#cd --command goreleaser check
```

For docs changes, dogfood the built formatter before the final checks:

```sh
./result/bin/mdfmt --write README.md CONTRIBUTING.md AGENTS.md CLAUDE.md RELEASING.md
```

Confirm `README.md`, `CONTRIBUTING.md`, and this file contain no machine-local
absolute paths.

## 3. Tag Release

Create an annotated semantic version tag from the checked commit:

```sh
git tag -a vX.Y.Z -m "vX.Y.Z"
git push origin main
git push origin vX.Y.Z
```

Tags matching `v[0-9]*` start the release workflow. Do not push the tag until
the pre-release checklist passes.

## 4. Manual Release

The release workflow also supports `workflow_dispatch`. Use the Actions UI for
manual runs when a rerun or controlled publication is needed. Prefer running it
from the release tag ref so GoReleaser and version metadata match the intended
release.

## 5. Version Metadata

The Nix package version comes from the current ref when it is a semantic
version tag. Otherwise it falls back to a development version based on the git
revision.

GoReleaser injects version, commit, date, and tree-state metadata through
linker flags. The CLI exposes the result through:

```sh
mdfmt version
```

## 6. Expected Artifacts

GoReleaser follows `.goreleaser.yaml`. The current release archives build the
configured `markdown-formatter` binary for darwin and linux on amd64 and arm64,
include `README.md` and `LICENSE`, and publish a checksums file.

After the workflow succeeds, verify the GitHub Release contains the expected
archives and checksums. Download one archive if needed and confirm the binary
runs `version`.
