# Markdown Behavior Decisions

Compatibility mode is the default release path. Behavior changes that would
alter existing output should be shipped behind a documented future release
plan, not silently folded into compatibility work.

## 1. Current Decisions

- Strict ATX spacing: deferred. Headings without a space after the marker are
  still recognized and normalized to one space.
- Setext headings: deferred. Only ATX headings are formatted.
- Skipped-level zero segments: preserved. A direct h2 to h4 transition renders
  as `1.0.1`.
- Indented headings: up to three leading spaces are accepted, matching
  CommonMark indentation tolerance.
- Prefix stripping: decimal prefixes such as `1.`, `1.2.`, and `1.2.3` are
  replaced by format mode and removed by remove mode.
- Seven-or-more hash lines: ignored rather than treated as h6.
- Heading spacing: the formatter normalizes zero or multiple blank lines around
  ATX headings to exactly one blank line when a neighboring line exists. It
  skips document start/end boundaries and ignores headings inside fenced code
  blocks.

## 2. Future Release Path

Future behavior changes should be tracked as focused issues, include fixtures
for the old and new behavior, and document whether the change is a compatibility
break or an opt-in mode.
