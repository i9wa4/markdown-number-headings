# Contributing

## 1. Local Workflow

Use the Nix development shell for repository work:

```sh
nix develop
```

Run the focused local checks before committing:

```sh
nix fmt
nix flake check --print-build-logs
nix build --print-build-logs
```

For a fast Go-only loop while editing formatter behavior, use:

```sh
nix develop .#ci --command go test ./...
```

Stage newly added files before running the full flake check. The pre-commit
check surface is part of `nix flake check`, and staged files make the local
result match the CI workflow more closely.

## 2. Formatter Development

`mdfmt` is the preferred public CLI. The default formatter pipeline runs
heading numbering, heading spacing normalization, and pipe table alignment.

Useful local commands:

```sh
go test ./...
nix build --print-build-logs
./result/bin/mdfmt --write README.md CONTRIBUTING.md AGENTS.md CLAUDE.md RELEASING.md
```

Use `--no-heading-numbering` when checking spacing and table formatting without
renumbering headings. Keep behavior decisions that affect Markdown output in
`docs/behavior-decisions.md`.

## 3. Nix And CI

The flake owns the reproducible toolchain, formatter checks, release tooling,
and CI entry points. Keep `go.mod`, `flake.nix`, and the Go override aligned
when changing Go versions.

The primary CI workflow is `.github/workflows/ci.yml` with workflow name `ci`.
It runs Nix checks, Nix build, Go vulnerability checks, and secret scanning.

## 4. Commit Expectations

Use conventional commit messages such as `docs: update release guide` or
`fix(formatter): preserve fenced code blocks`. Keep docs, behavior changes,
workflow changes, and dependency changes separated when the review surface is
meaningfully different.

Do not include machine-local absolute paths in commits, public docs, issues, or
pull request text. Use repo-relative paths such as `README.md` and
`.github/workflows/ci.yml`.

## 5. Releases

Release work is documented in `RELEASING.md`. Run the pre-release checklist
there before creating a release tag or using the manual release workflow.
