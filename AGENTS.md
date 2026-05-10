# Agent Instructions

## 1. Project Scope

This repository is a Go Markdown formatter project. Keep changes focused on
`mdfmt`, its formatter package, Nix workflow, GitHub Actions, and concise
repository documentation.

`CLAUDE.md` is a regular file whose body is exactly `@AGENTS.md`, so this file
is the single source of truth for agent guidance.

## 2. Development Rules

- Prefer existing Go, Nix, and GitHub Actions patterns already in the repo.
- Keep public docs and commit messages repo-relative; do not write
  machine-local absolute paths there.
- Stage newly added files before running full Nix checks.
- Run focused checks for the files touched, then run the broader checks when a
  change affects shared formatter behavior or workflow contracts.
- Do not push without an explicit instruction.

## 3. Formatter Rules

`mdfmt` is the preferred user-facing CLI. Default formatting includes heading
numbering, heading spacing normalization, and pipe table alignment. Use
`--no-heading-numbering` only when the task explicitly wants spacing and table
formatting without renumbering headings.

Document behavior changes in `docs/behavior-decisions.md` when they affect the
formatter contract.

## 4. Required Checks

For docs or workflow changes, run at least:

```sh
nix fmt
nix develop .#ci --command go test ./...
```

For release, dependency, Nix, CI, or shared formatter behavior changes, add:

```sh
nix flake check --print-build-logs
nix build --print-build-logs
```

Use `RELEASING.md` for release-specific verification.
