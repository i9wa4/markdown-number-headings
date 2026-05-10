package formatter

import (
	"regexp"
	"strconv"
	"strings"
)

type Options struct {
	Shift int
}

var numberedPrefix = regexp.MustCompile(`^\d+(?:\.\d+)*\.?\s+`)

func Format(input string, opts Options) string {
	if opts.Shift < 0 {
		opts.Shift = 1
	}
	lines, trailing := splitLines(input)
	counts := make([]int, 6)
	inFence := false
	fenceMarker := byte(0)

	for i, line := range lines {
		if marker, ok := fenceStart(line); ok {
			if !inFence {
				inFence = true
				fenceMarker = marker
			} else if marker == fenceMarker {
				inFence = false
				fenceMarker = 0
			}
			continue
		}
		if inFence {
			continue
		}

		h, ok := parseHeading(line)
		if !ok || h.level <= opts.Shift {
			continue
		}
		idx := h.level - opts.Shift - 1
		counts[idx]++
		for j := idx + 1; j < len(counts); j++ {
			counts[j] = 0
		}

		segments := make([]string, idx+1)
		for j := 0; j <= idx; j++ {
			segments[j] = strconv.Itoa(counts[j])
		}
		title := stripNumber(h.text)
		lines[i] = h.prefix + strings.Repeat("#", h.level) + " " + strings.Join(segments, ".") + ". " + title + h.close
	}

	return joinLines(lines, trailing)
}

func Remove(input string) string {
	lines, trailing := splitLines(input)
	inFence := false
	fenceMarker := byte(0)

	for i, line := range lines {
		if marker, ok := fenceStart(line); ok {
			if !inFence {
				inFence = true
				fenceMarker = marker
			} else if marker == fenceMarker {
				inFence = false
				fenceMarker = 0
			}
			continue
		}
		if inFence {
			continue
		}
		h, ok := parseHeading(line)
		if !ok {
			continue
		}
		lines[i] = h.prefix + strings.Repeat("#", h.level) + " " + stripNumber(h.text) + h.close
	}

	return joinLines(lines, trailing)
}

type heading struct {
	prefix string
	level  int
	text   string
	close  string
}

func parseHeading(line string) (heading, bool) {
	prefixLen := 0
	for prefixLen < len(line) && prefixLen < 3 && line[prefixLen] == ' ' {
		prefixLen++
	}
	rest := line[prefixLen:]
	level := 0
	for level < len(rest) && rest[level] == '#' {
		level++
	}
	if level == 0 || level > 6 {
		return heading{}, false
	}
	text := strings.TrimSpace(rest[level:])
	if text == "" {
		return heading{}, false
	}
	text, close := splitClosingSequence(text)
	return heading{
		prefix: line[:prefixLen],
		level:  level,
		text:   text,
		close:  close,
	}, true
}

func splitClosingSequence(text string) (string, string) {
	trimmed := strings.TrimRight(text, " ")
	hashStart := len(trimmed)
	for hashStart > 0 && trimmed[hashStart-1] == '#' {
		hashStart--
	}
	if hashStart == len(trimmed) || hashStart == 0 || trimmed[hashStart-1] != ' ' {
		return strings.TrimSpace(text), ""
	}
	return strings.TrimSpace(trimmed[:hashStart-1]), " " + trimmed[hashStart:]
}

func stripNumber(text string) string {
	return numberedPrefix.ReplaceAllString(strings.TrimSpace(text), "")
}

func fenceStart(line string) (byte, bool) {
	trimmed := strings.TrimLeft(line, " ")
	if len(trimmed) < 3 {
		return 0, false
	}
	marker := trimmed[0]
	if marker != '`' && marker != '~' {
		return 0, false
	}
	if trimmed[1] != marker || trimmed[2] != marker {
		return 0, false
	}
	return marker, true
}

func splitLines(input string) ([]string, bool) {
	trailing := strings.HasSuffix(input, "\n")
	if trailing {
		input = strings.TrimSuffix(input, "\n")
	}
	if input == "" {
		return nil, trailing
	}
	return strings.Split(input, "\n"), trailing
}

func joinLines(lines []string, trailing bool) string {
	out := strings.Join(lines, "\n")
	if trailing {
		out += "\n"
	}
	return out
}
