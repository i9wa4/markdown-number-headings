package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDispatchHelpVersionAndErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr string
	}{
		{name: "help", args: []string{"help"}, wantOut: "Usage:"},
		{name: "version", args: []string{"version"}, wantOut: "mdfmt"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			err := Run(tc.args, strings.NewReader(""), &stdout, &stderr)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(stdout.String(), tc.wantOut) {
				t.Fatalf("stdout missing %q:\n%s", tc.wantOut, stdout.String())
			}
		})
	}
}

func TestFormatRemoveAndTableStdin(t *testing.T) {
	tests := []struct {
		args  []string
		input string
		want  string
	}{
		{args: nil, input: "## One\n### Two\n", want: "## 1. One\n\n### 1.1. Two\n"},
		{args: []string{"format"}, input: "## One\n### Two\n", want: "## 1. One\n\n### 1.1. Two\n"},
		{args: []string{"format", "--shift", "0"}, input: "# One\n", want: "# 1. One\n"},
		{args: []string{"--no-heading-numbering"}, input: "# One\nbody\n", want: "# One\n\nbody\n"},
		{args: []string{"remove-numbers"}, input: "## 1. One\n", want: "## One\n"},
		{args: []string{"remove"}, input: "## 1. One\n", want: "## One\n"},
		{args: []string{"table"}, input: "| A|B |\n|---|---|\n| x|yy|\n", want: "| A   | B   |\n| --- | --- |\n| x   | yy  |\n"},
		{args: []string{"spacing"}, input: "# One\nbody\n## Two\n", want: "# One\n\nbody\n\n## Two\n"},
	}
	for _, tc := range tests {
		var stdout, stderr bytes.Buffer
		if err := Run(tc.args, strings.NewReader(tc.input), &stdout, &stderr); err != nil {
			t.Fatalf("%v stderr=%s", err, stderr.String())
		}
		if stdout.String() != tc.want {
			t.Fatalf("want %q got %q", tc.want, stdout.String())
		}
	}
}

func TestDefaultFormatRunsAllPasses(t *testing.T) {
	input := "## One\nbody\n| 名前|Value |\n|---|---|\n| 日本|1|\n"
	want := "## 1. One\n\nbody\n| 名前 | Value |\n| ---- | ----- |\n| 日本 | 1     |\n"
	var stdout, stderr bytes.Buffer
	if err := Run(nil, strings.NewReader(input), &stdout, &stderr); err != nil {
		t.Fatalf("%v stderr=%s", err, stderr.String())
	}
	if stdout.String() != want {
		t.Fatalf("want %q got %q", want, stdout.String())
	}
}

func TestDefaultFormatNormalizesFinalNewline(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run(nil, strings.NewReader("## Final"), &stdout, &stderr); err != nil {
		t.Fatalf("%v stderr=%s", err, stderr.String())
	}
	if stdout.String() != "## 1. Final\n" {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
}

func TestFileWriteMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "doc.md")
	if err := os.WriteFile(path, []byte("## One\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"--write", path}, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatalf("%v stderr=%s", err, stderr.String())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "## 1. One\n" {
		t.Fatalf("unexpected file content: %q", string(data))
	}
}

func TestInvalidShift(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Run([]string{"format", "--shift", "-1"}, strings.NewReader(""), &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "invalid shift") {
		t.Fatalf("expected invalid shift error, got %v", err)
	}
}
