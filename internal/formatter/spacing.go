package formatter

import "strings"

func NormalizeHeadingSpacing(input string) string {
	lines, trailing := splitLines(input)
	if len(lines) == 0 {
		return joinLines(lines, trailing)
	}

	headingLines := mapHeadingLines(lines)
	out := make([]string, 0, len(lines))
	prevHasContent := false
	prevContentHeading := false

	for i := 0; i < len(lines); {
		if isBlankLine(lines[i]) {
			start := i
			for i < len(lines) && isBlankLine(lines[i]) {
				i++
			}
			next := i
			nextHeading := next < len(lines) && headingLines[next]

			switch {
			case prevHasContent && next < len(lines) && (prevContentHeading || nextHeading):
				if len(out) == 0 || !isBlankLine(out[len(out)-1]) {
					out = append(out, "")
				}
			case prevHasContent && next == len(lines) && prevContentHeading:
				continue
			case !prevHasContent && nextHeading:
				continue
			default:
				out = append(out, lines[start:i]...)
			}
			continue
		}

		if prevHasContent && (prevContentHeading || headingLines[i]) && (len(out) == 0 || !isBlankLine(out[len(out)-1])) {
			out = append(out, "")
		}
		out = append(out, lines[i])
		prevHasContent = true
		prevContentHeading = headingLines[i]
		i++
	}

	return joinLines(out, trailing)
}

func mapHeadingLines(lines []string) []bool {
	headingLines := make([]bool, len(lines))
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
		_, headingLines[i] = parseHeading(line)
	}

	return headingLines
}

func isBlankLine(line string) bool {
	return strings.TrimSpace(line) == ""
}
