# Claude Instructions

## 1. Operating Contract

Follow `AGENTS.md` as the shared agent contract. This file keeps the same
project-specific guidance visible for Claude Code and compatible tools.

## 2. Repository Priorities

- Preserve the formatter-only scope of the project.
- Prefer `mdfmt` in user-facing examples and docs.
- Keep `README.md` concise; move contribution detail to `CONTRIBUTING.md` and
  release detail to `RELEASING.md`.
- Keep CI references aligned with `.github/workflows/ci.yml` and workflow name
  `ci`.

## 3. Verification

Use Nix for repeatable checks. For most docs and workflow tasks:

```sh
nix fmt
nix develop .#ci --command go test ./...
```

For changes with release, dependency, formatter, or CI risk:

```sh
nix flake check --print-build-logs
nix build --print-build-logs
```

When Markdown docs change, dogfood the formatter from the flake build where
practical:

```sh
./result/bin/mdfmt --write README.md CONTRIBUTING.md AGENTS.md CLAUDE.md RELEASING.md
```

## 4. Public Surface Hygiene

Public docs, commits, pull requests, and issue comments must use repo-relative
paths or stable URLs. Local absolute paths are acceptable only in private task
artifacts or local debug notes.
