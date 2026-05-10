package formatter

import "strings"

func NormalizeHeadingSpacing(input string) string {
	return applyDocument(input, headingSpacingPass())
}

func headingSpacingPass() documentPass {
	return func(doc *document) {
		if len(doc.lines) == 0 {
			return
		}

		headingLines := mapHeadingLines(doc)
		out := make([]string, 0, len(doc.lines))
		prevHasContent := false
		prevContentHeading := false

		for i := 0; i < len(doc.lines); {
			if doc.inProtectedBlock(i) {
				out = append(out, doc.lines[i])
				prevHasContent = true
				prevContentHeading = false
				i++
				continue
			}

			if isBlankLine(doc.lines[i]) {
				start := i
				for i < len(doc.lines) && isBlankLine(doc.lines[i]) {
					i++
				}
				next := i
				nextHeading := next < len(doc.lines) && headingLines[next]

				switch {
				case prevHasContent && next < len(doc.lines) && (prevContentHeading || nextHeading):
					if len(out) == 0 || !isBlankLine(out[len(out)-1]) {
						out = append(out, "")
					}
				case prevHasContent && next == len(doc.lines) && prevContentHeading:
					continue
				case !prevHasContent && nextHeading:
					continue
				default:
					out = append(out, doc.lines[start:i]...)
				}
				continue
			}

			if prevHasContent && (prevContentHeading || headingLines[i]) && (len(out) == 0 || !isBlankLine(out[len(out)-1])) {
				out = append(out, "")
			}
			out = append(out, doc.lines[i])
			prevHasContent = true
			prevContentHeading = headingLines[i]
			i++
		}

		doc.setLines(out)
	}
}

func mapHeadingLines(doc *document) []bool {
	headingLines := make([]bool, len(doc.lines))
	for i, line := range doc.lines {
		if doc.inProtectedBlock(i) {
			continue
		}
		_, headingLines[i] = parseHeading(line)
	}

	return headingLines
}

func isBlankLine(line string) bool {
	return strings.TrimSpace(line) == ""
}
