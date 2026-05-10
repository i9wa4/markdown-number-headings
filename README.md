# markdown-formatter

`markdown-formatter` provides `mdfmt`, a small Markdown formatter for heading
numbers, heading spacing, and pipe tables.

## 1. Install

```sh
go install github.com/i9wa4/markdown-formatter/cmd/mdfmt@latest
```

From this repository:

```sh
nix build
./result/bin/mdfmt version
```

## 2. Usage

Format stdin:

```sh
mdfmt < README.md
```

Update files in place:

```sh
mdfmt --write README.md docs/*.md
```

Skip heading numbering:

```sh
mdfmt --no-heading-numbering --write README.md
```

In Vim, filter the current buffer through the formatter:

```vim
:%!mdfmt
```

## 3. Formatting

| Pass            | Behavior                                     |
| --------------- | -------------------------------------------- |
| Heading numbers | Numbers ATX headings from h2 by default.     |
| Heading spacing | Keeps one blank line around headings.        |
| Tables          | Aligns pipe tables, including CJK width 2.   |
| EOF newline     | Keeps exactly one final newline.             |

Useful options:

- `--shift 0` starts heading numbering from h1.
- `--no-heading-numbering` keeps spacing and table formatting, but does not add
  or update heading numbers.
- `--write` updates files in place.

## 4. Development

```sh
nix develop
nix fmt
nix flake check --print-build-logs
nix build --print-build-logs
```

See `CONTRIBUTING.md` and `RELEASING.md` for project workflow and release
steps.
