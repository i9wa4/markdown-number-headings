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
		return runDocumentFormat(args, stdin, stdout)
	}

	switch args[0] {
	case "help", "--help", "-h":
		writeUsage(stdout)
		return nil
	case "version", "--version":
		fmt.Fprintf(stdout, "mdfmt %s (%s)\n", version.Version, version.Commit)
		return nil
	case "format":
		return runDocumentFormat(args[1:], stdin, stdout)
	case "remove-numbers", "remove":
		return runTransform(args[1:], stdin, stdout, func(_ formatter.Options) formatter.Pass {
			return formatter.HeadingNumberRemoval()
		})
	case "spacing":
		return runSimple(args[1:], stdin, stdout, "mdfmt spacing", formatter.HeadingSpacing())
	case "table":
		return runSimple(args[1:], stdin, stdout, "mdfmt table", formatter.TableAlignment())
	default:
		return runDocumentFormat(args, stdin, stdout)
	}
}

func runDocumentFormat(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("mdfmt", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	shift := fs.Int("shift", 1, "heading level shift")
	noHeadingNumbering := fs.Bool("no-heading-numbering", false, "skip heading numbering")
	write := fs.Bool("write", false, "write files in place")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *shift < 0 || *shift > 5 {
		return fmt.Errorf("invalid shift %s: expected an integer from 0 to 5", strconv.Itoa(*shift))
	}
	pass := formatter.DocumentFormatting(formatter.Options{Shift: *shift})
	if *noHeadingNumbering {
		pass = formatter.DocumentFormattingWithoutHeadingNumbering()
	}
	return apply(fs.Args(), stdin, stdout, *write, func(input string) string {
		return formatter.Apply(input, pass)
	})
}

func runTransform(args []string, stdin io.Reader, stdout io.Writer, passForOptions func(formatter.Options) formatter.Pass) error {
	fs := flag.NewFlagSet("mdfmt", flag.ContinueOnError)
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
		return formatter.Apply(input, passForOptions(formatter.Options{Shift: *shift}))
	})
}

func runSimple(args []string, stdin io.Reader, stdout io.Writer, command string, pass formatter.Pass) error {
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	write := fs.Bool("write", false, "write files in place")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return apply(fs.Args(), stdin, stdout, *write, func(input string) string {
		return formatter.Apply(input, pass)
	})
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
	fmt.Fprint(w, `mdfmt formats Markdown headings, spacing, and tables.

Usage:
  mdfmt [--shift N] [--no-heading-numbering] [--write] [FILE...]
  mdfmt help
  mdfmt version

Examples:
  mdfmt < README.md
  mdfmt --write README.md docs/behavior-decisions.md
  mdfmt --no-heading-numbering --write README.md
`)
}
