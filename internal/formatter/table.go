package formatter

import "strings"

func FormatTables(input string) string {
	lines, trailing := splitLines(input)
	var out []string
	for i := 0; i < len(lines); {
		if i+1 < len(lines) && isTableSeparator(lines[i+1]) {
			end := i + 2
			for end < len(lines) && isTableRow(lines[end]) {
				end++
			}
			out = append(out, formatTableBlock(lines[i:end])...)
			i = end
			continue
		}
		out = append(out, lines[i])
		i++
	}
	return joinLines(out, trailing)
}

func formatTableBlock(block []string) []string {
	rows := make([][]string, len(block))
	widths := []int{}
	aligns := []string{}
	for i, line := range block {
		cells := splitTableRow(line)
		rows[i] = cells
		for len(widths) < len(cells) {
			widths = append(widths, 3)
		}
		for j, cell := range cells {
			width := displayWidth(cell)
			if width > widths[j] {
				widths[j] = width
			}
		}
	}
	for _, cell := range rows[1] {
		aligns = append(aligns, alignmentFor(cell))
	}
	for len(aligns) < len(widths) {
		aligns = append(aligns, "---")
	}

	formatted := make([]string, len(block))
	for i, row := range rows {
		if i == 1 {
			formatted[i] = formatSeparator(widths, aligns)
			continue
		}
		formatted[i] = formatDataRow(row, widths)
	}
	return formatted
}

func splitTableRow(line string) []string {
	trimmed := strings.TrimSpace(line)
	trimmed = strings.TrimPrefix(trimmed, "|")
	trimmed = strings.TrimSuffix(trimmed, "|")
	parts := strings.Split(trimmed, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func formatDataRow(row []string, widths []int) string {
	cells := make([]string, len(widths))
	for i := range widths {
		cell := ""
		if i < len(row) {
			cell = row[i]
		}
		cells[i] = " " + cell + strings.Repeat(" ", widths[i]-displayWidth(cell)) + " "
	}
	return "|" + strings.Join(cells, "|") + "|"
}

func formatSeparator(widths []int, aligns []string) string {
	cells := make([]string, len(widths))
	for i := range widths {
		switch aligns[i] {
		case ":--":
			cells[i] = " :" + strings.Repeat("-", widths[i]-1) + " "
		case "--:":
			cells[i] = " " + strings.Repeat("-", widths[i]-1) + ": "
		case ":-:":
			cells[i] = " :" + strings.Repeat("-", widths[i]-2) + ": "
		default:
			cells[i] = " " + strings.Repeat("-", widths[i]) + " "
		}
	}
	return "|" + strings.Join(cells, "|") + "|"
}

func isTableRow(line string) bool {
	return strings.Contains(line, "|")
}

func isTableSeparator(line string) bool {
	cells := splitTableRow(line)
	if len(cells) == 0 {
		return false
	}
	for _, cell := range cells {
		trimmed := strings.Trim(cell, " ")
		trimmed = strings.TrimPrefix(trimmed, ":")
		trimmed = strings.TrimSuffix(trimmed, ":")
		if len(trimmed) < 3 || strings.Trim(trimmed, "-") != "" {
			return false
		}
	}
	return true
}

func alignmentFor(cell string) string {
	trimmed := strings.TrimSpace(cell)
	left := strings.HasPrefix(trimmed, ":")
	right := strings.HasSuffix(trimmed, ":")
	switch {
	case left && right:
		return ":-:"
	case left:
		return ":--"
	case right:
		return "--:"
	default:
		return "---"
	}
}

func displayWidth(s string) int {
	width := 0
	for _, r := range s {
		if isWideRune(r) {
			width += 2
			continue
		}
		width++
	}
	return width
}

func isWideRune(r rune) bool {
	return (r >= 0x1100 && r <= 0x115F) ||
		(r >= 0x2329 && r <= 0x232A) ||
		(r >= 0x2E80 && r <= 0xA4CF) ||
		(r >= 0xAC00 && r <= 0xD7A3) ||
		(r >= 0xF900 && r <= 0xFAFF) ||
		(r >= 0xFE10 && r <= 0xFE19) ||
		(r >= 0xFE30 && r <= 0xFE6F) ||
		(r >= 0xFF00 && r <= 0xFF60) ||
		(r >= 0xFFE0 && r <= 0xFFE6)
}
