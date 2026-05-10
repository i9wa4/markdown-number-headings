# markdown-formatter

`markdown-formatter` is a small Go CLI for Markdown heading numbering, heading
number removal, and pipe table alignment. It replaces the original
`markdown-number-headings` editor behavior with a formatter that can run in
shell pipelines, CI, and editor integrations.

Implementation status: the Go CLI, pure formatter package, Nix development
environment, CI, and release metadata are present. The compatibility mode keeps
the current heading behavior stable while future Markdown behavior changes are
documented separately in `docs/behavior-decisions.md`.

## Install

From this repository:

```sh
nix build
./result/bin/markdown-formatter version
```

With Go:

```sh
go install github.com/i9wa4/markdown-formatter@latest
```

## Usage

Number headings from h2 by default:

```sh
markdown-formatter format < README.md
markdown-formatter format --write README.md
```

Remove heading numbers:

```sh
markdown-formatter remove --write README.md
```

Start numbering from h1 with `--shift 0`:

```sh
markdown-formatter format --shift 0 --write README.md
```

Align Markdown pipe tables:

```sh
markdown-formatter table --write README.md
```

Multiple file paths are accepted. Without file paths, commands read from stdin
and write to stdout. `--write` updates files in place and requires file paths.

## Formatter Contract

- ATX headings with one through six `#` characters are supported.
- Default `--shift 1` leaves h1 unnumbered and starts numbering at h2.
- `--shift 0` starts numbering at h1.
- Skipped heading levels preserve zero segments, for example h2 followed by h4
  becomes `1.0.1`.
- Existing decimal prefixes are replaced in format mode and removed in remove
  mode.
- Fenced code blocks using backticks or tildes are ignored.
- ATX headings without a space after `#` are accepted for compatibility and
  normalized to a single space.
- Lines beginning with seven or more `#` characters are not treated as headings.
- Setext headings and broader Markdown AST rewrites are non-goals for the MVP.

## Vim Migration

The repository is formatter-only. Vim or Neovim integrations should call the
CLI instead of duplicating formatter logic.

Existing command mapping:

- `NumberHeader` maps to `markdown-formatter format --write <file>`.
- `RemoveNumbers` maps to `markdown-formatter remove --write <file>`.
- Header level shift maps to `markdown-formatter format --shift N`.
- The default behavior still starts numbering at h2.

Denops retirement should happen after editor wrappers or user configuration
examples point at this CLI.

## Development

```sh
nix develop
nix fmt
nix flake check
nix build
```

The flake exposes Go, formatter, workflow, Nix, Markdown, YAML, and security
checks where they are useful for this project.

## Release

Releases are tag-based or manual through `.github/workflows/release.yml`.
GoReleaser builds darwin and linux archives for amd64 and arm64, includes
`README.md` and `LICENSE`, and publishes checksums.

Before publishing:

```sh
nix flake check
nix build
goreleaser check
```
