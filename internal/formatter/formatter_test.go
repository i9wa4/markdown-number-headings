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

func readFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "fixtures", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
