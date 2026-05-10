package formatter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFixtures(t *testing.T) {
	cases := []struct {
		name  string
		shift int
	}{
		{name: "basic", shift: 1},
		{name: "compat", shift: 1},
		{name: "shift0", shift: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := readFixture(t, tc.name+".input.golden")
			wantNumbered := readFixture(t, tc.name+".numbered.golden")
			if got := Format(input, Options{Shift: tc.shift}); got != wantNumbered {
				t.Fatalf("Format mismatch\nwant:\n%s\ngot:\n%s", wantNumbered, got)
			}
			if tc.name == "shift0" {
				return
			}
			wantRemoved := readFixture(t, tc.name+".removed.golden")
			if got := Remove(wantNumbered); got != wantRemoved {
				t.Fatalf("Remove mismatch\nwant:\n%s\ngot:\n%s", wantRemoved, got)
			}
		})
	}
}

func TestFormatTables(t *testing.T) {
	input := readFixture(t, "table.input.golden")
	want := readFixture(t, "table.formatted.golden")
	if got := FormatTables(input); got != want {
		t.Fatalf("FormatTables mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestFormatTablesUsesDisplayWidth(t *testing.T) {
	input := "| Name | Value |\n| --- | --- |\n| 日本 | 1 |\n| Go | 20 |\n"
	want := "| Name | Value |\n| ---- | ----- |\n| 日本 | 1     |\n| Go   | 20    |\n"
	if got := FormatTables(input); got != want {
		t.Fatalf("FormatTables display-width mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestFormatTablesPreservesFencedTables(t *testing.T) {
	tests := []struct {
		name  string
		fence string
	}{
		{name: "backtick", fence: "```"},
		{name: "tilde", fence: "~~~"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.fence + "md\n| A|B |\n|---|---|\n| x|yy|\n" + tc.fence + "\n\n| A|B |\n|---|---|\n| x|yy|\n"
			want := tc.fence + "md\n| A|B |\n|---|---|\n| x|yy|\n" + tc.fence + "\n\n| A   | B   |\n| --- | --- |\n| x   | yy  |\n"
			if got := FormatTables(input); got != want {
				t.Fatalf("FormatTables fenced-table mismatch\nwant:\n%s\ngot:\n%s", want, got)
			}
		})
	}
}

func TestDocumentFormattingWithoutHeadingNumberingPreservesFencedTables(t *testing.T) {
	tests := []struct {
		name  string
		fence string
	}{
		{name: "backtick", fence: "```"},
		{name: "tilde", fence: "~~~"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.fence + "md\n# Ignored\n| A|B |\n|---|---|\n| x|yy|\n" + tc.fence + "\n# Real\n| 名前|Value |\n|---|---|\n| 日本|1|\n"
			want := tc.fence + "md\n# Ignored\n| A|B |\n|---|---|\n| x|yy|\n" + tc.fence + "\n\n# Real\n\n| 名前 | Value |\n| ---- | ----- |\n| 日本 | 1     |\n"
			if got := DocumentFormattingWithoutHeadingNumbering()(input); got != want {
				t.Fatalf("DocumentFormattingWithoutHeadingNumbering fenced-table mismatch\nwant:\n%s\ngot:\n%s", want, got)
			}
		})
	}
}

func TestNormalizeHeadingSpacingFixture(t *testing.T) {
	input := readFixture(t, "heading-spacing.input.golden")
	want := readFixture(t, "heading-spacing.formatted.golden")
	if got := NormalizeHeadingSpacing(input); got != want {
		t.Fatalf("NormalizeHeadingSpacing mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestNormalizeHeadingSpacingBoundaries(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "heading at file start",
			input: "# Title\nbody\n",
			want:  "# Title\n\nbody\n",
		},
		{
			name:  "heading at eof",
			input: "body\n# End\n",
			want:  "body\n\n# End\n",
		},
		{
			name:  "heading followed by body",
			input: "## Section\nbody\n",
			want:  "## Section\n\nbody\n",
		},
		{
			name:  "body followed by heading",
			input: "body\n## Section\n",
			want:  "body\n\n## Section\n",
		},
		{
			name:  "adjacent headings",
			input: "# One\n## Two\n",
			want:  "# One\n\n## Two\n",
		},
		{
			name:  "multiple blank lines",
			input: "# One\n\n\nbody\n\n\n## Two\n\n\ntext\n",
			want:  "# One\n\nbody\n\n## Two\n\ntext\n",
		},
		{
			name:  "one blank line unchanged",
			input: "# One\n\nbody\n\n## Two\n\ntext\n",
			want:  "# One\n\nbody\n\n## Two\n\ntext\n",
		},
		{
			name:  "no blank lines",
			input: "# One\nbody\n## Two\ntext\n",
			want:  "# One\n\nbody\n\n## Two\n\ntext\n",
		},
		{
			name:  "fenced code headings",
			input: "```md\n# Ignored\nbody\n```\n# Real\nbody\n",
			want:  "```md\n# Ignored\nbody\n```\n\n# Real\n\nbody\n",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := NormalizeHeadingSpacing(tc.input); got != tc.want {
				t.Fatalf("want %q got %q", tc.want, got)
			}
		})
	}
}

func readFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "fixtures", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
