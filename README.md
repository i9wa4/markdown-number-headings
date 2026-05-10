# markdown-formatter

`markdown-formatter` provides `mdfmt`, a small Go CLI for Markdown heading
numbering, heading spacing, and pipe table alignment. It replaces the original
`markdown-number-headings` editor behavior with a formatter that can run in
shell pipelines, CI, and editor integrations.

Implementation status: the Go CLI, pure formatter package, Nix development
environment, CI, and release metadata are present. The compatibility mode keeps
the current heading behavior stable while future Markdown behavior changes are
documented separately in `docs/behavior-decisions.md`.

## 1. Project Docs

- `CONTRIBUTING.md` covers local development, checks, formatter workflow, and
  commit expectations.
- `RELEASING.md` covers tag-based and manual releases, GoReleaser, Nix checks,
  version metadata, and expected artifacts.
- `AGENTS.md` and `CLAUDE.md` keep AI-agent guidance aligned with this repo's
  workflow.

## 2. Install

From this repository:

```sh
nix build
./result/bin/mdfmt version
```

With Go:

```sh
go install github.com/i9wa4/markdown-formatter/cmd/mdfmt@latest
```

## 3. Usage

Format Markdown with the full formatter pipeline:

```sh
mdfmt < README.md
mdfmt --write README.md
```

Skip heading numbering while keeping spacing and table formatting:

```sh
mdfmt --no-heading-numbering --write README.md
```

Start numbering from h1 with `--shift 0`:

```sh
mdfmt --shift 0 --write README.md
```

Multiple file paths are accepted. Without file paths, `mdfmt` reads from stdin
and writes to stdout. `--write` updates files in place and requires file paths.

| Invocation                     | Purpose                                                |
| ------------------------------ | ------------------------------------------------------ |
| `mdfmt`                        | Run all formatter passes on stdin and write to stdout. |
| `mdfmt --write FILE...`        | Run all formatter passes and update files in place.    |
| `mdfmt --no-heading-numbering` | Format spacing and tables without heading numbering.   |
| `mdfmt --shift 0 --write FILE` | Start heading numbering from h1.                       |

The old `markdown-formatter` binary name and focused subcommands remain
compatibility paths where available, but new usage should prefer `mdfmt`.
For existing scripts that only remove heading numbers, use
`mdfmt remove-numbers --write FILE`.

## 4. Formatter Contract

- ATX headings with one through six `#` characters are supported.
- Default formatting runs heading numbering, table alignment, and heading
  spacing normalization.
- Default `--shift 1` leaves h1 unnumbered and starts numbering at h2.
- `--shift 0` starts numbering at h1.
- `--no-heading-numbering` skips only the heading-numbering pass.
- Skipped heading levels preserve zero segments, for example h2 followed by h4
  becomes `1.0.1`.
- Existing decimal prefixes are replaced in format mode and removed in remove
  mode.
- Fenced code blocks using backticks or tildes are ignored.
- Heading spacing normalization ensures exactly one blank line above and below
  ATX headings outside fenced code blocks where a neighboring line exists.
- Heading spacing normalization does not add a blank line before a heading at
  the start of a document or after a heading at the end of a document.
- ATX headings without a space after `#` are accepted for compatibility and
  normalized to a single space.
- Lines beginning with seven or more `#` characters are not treated as headings.
- Setext headings and broader Markdown AST rewrites are non-goals for the MVP.

## 5. Vim Migration

The repository is formatter-only. Vim or Neovim integrations should call the
CLI instead of duplicating formatter logic.

Existing command mapping:

- `NumberHeader` maps to `mdfmt --write <file>`.
- Formatting without heading numbering maps to
  `mdfmt --no-heading-numbering --write <file>`.
- Header level shift maps to `mdfmt --shift N --write <file>`.
- The default behavior still starts numbering at h2.

Denops retirement should happen after editor wrappers or user configuration
examples point at this CLI.

## 6. Development

```sh
nix develop
nix fmt
nix flake check --print-build-logs
nix build --print-build-logs
```

The flake exposes Go, formatter, workflow, Nix, Markdown, YAML, and security
checks where they are useful for this project.

The pre-commit check surface dogfoods this project by running the flake-built
`mdfmt --write` hook on Markdown files. It does not depend on a globally
installed `mdfmt` binary, and no standalone lint command is part of the public
CLI.

See `CONTRIBUTING.md` for contribution workflow, commit expectations, and
focused local commands.

## 7. Release

Releases are tag-based or manual through `.github/workflows/release.yml`.
GoReleaser publishes platform archives and checksums. See `RELEASING.md` for
the pre-release checklist, version metadata behavior, and artifact
expectations.
