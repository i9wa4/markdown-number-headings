package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/i9wa4/markdown-formatter/internal/formatter"
	"github.com/i9wa4/markdown-formatter/internal/version"
)

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		writeUsage(stderr)
		return fmt.Errorf("missing command: use format, remove, table, help, or version")
	}

	switch args[0] {
	case "help", "--help", "-h":
		writeUsage(stdout)
		return nil
	case "version", "--version":
		fmt.Fprintf(stdout, "markdown-formatter %s (%s)\n", version.Version, version.Commit)
		return nil
	case "format":
		return runTransform(args[1:], stdin, stdout, formatter.Format)
	case "remove":
		return runTransform(args[1:], stdin, stdout, func(input string, _ formatter.Options) string {
			return formatter.Remove(input)
		})
	case "table":
		return runTable(args[1:], stdin, stdout)
	default:
		writeUsage(stderr)
		return fmt.Errorf("unknown command %q: use format, remove, table, help, or version", args[0])
	}
}

func runTransform(args []string, stdin io.Reader, stdout io.Writer, transform func(string, formatter.Options) string) error {
	fs := flag.NewFlagSet("markdown-formatter", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	shift := fs.Int("shift", 1, "heading level shift")
	write := fs.Bool("write", false, "write files in place")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *shift < 0 || *shift > 5 {
		return fmt.Errorf("invalid shift %s: expected an integer from 0 to 5", strconv.Itoa(*shift))
	}
	return apply(fs.Args(), stdin, stdout, *write, func(input string) string {
		return transform(input, formatter.Options{Shift: *shift})
	})
}

func runTable(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("markdown-formatter table", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	write := fs.Bool("write", false, "write files in place")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return apply(fs.Args(), stdin, stdout, *write, formatter.FormatTables)
}

func apply(paths []string, stdin io.Reader, stdout io.Writer, write bool, transform func(string) string) error {
	if len(paths) == 0 {
		if write {
			return fmt.Errorf("--write requires at least one file path")
		}
		data, err := io.ReadAll(stdin)
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}
		_, err = io.WriteString(stdout, transform(string(data)))
		return err
	}

	var combined string
	for i, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		result := transform(string(data))
		if write {
			if err := os.WriteFile(path, []byte(result), 0o644); err != nil {
				return fmt.Errorf("write %s: %w", path, err)
			}
			continue
		}
		if i > 0 && combined != "" && combined[len(combined)-1] != '\n' {
			combined += "\n"
		}
		combined += result
	}
	if write {
		return nil
	}
	_, err := io.WriteString(stdout, combined)
	return err
}

func writeUsage(w io.Writer) {
	fmt.Fprint(w, `markdown-formatter formats Markdown headings and tables.

Usage:
  markdown-formatter format [--shift N] [--write] [FILE...]
  markdown-formatter remove [--write] [FILE...]
  markdown-formatter table [--write] [FILE...]
  markdown-formatter version
  markdown-formatter help

Examples:
  markdown-formatter format < README.md
  markdown-formatter format --shift 0 --write docs/page.md
  markdown-formatter remove --write docs/page.md
  markdown-formatter table --write README.md
`)
}
