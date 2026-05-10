package formatter

import "strings"

func FormatTables(input string) string {
	return applyDocument(input, tableAlignmentPass())
}

func tableAlignmentPass() documentPass {
	return func(doc *document) {
		var out []string
		for i := 0; i < len(doc.lines); {
			if doc.inProtectedBlock(i) {
				out = append(out, doc.lines[i])
				i++
				continue
			}
			if i+1 < len(doc.lines) && isTableRow(doc.lines[i]) && !doc.inProtectedBlock(i+1) && isTableSeparator(doc.lines[i+1]) {
				end := i + 2
				for end < len(doc.lines) && !doc.inProtectedBlock(end) && isTableRow(doc.lines[end]) {
					end++
				}
				out = append(out, formatTableBlock(doc.lines[i:end])...)
				i = end
				continue
			}
			out = append(out, doc.lines[i])
			i++
		}
		doc.lines = out
	}
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
	parts := splitTableCells(trimmed)
	if len(parts) > 0 && strings.TrimSpace(parts[0]) == "" {
		parts = parts[1:]
	}
	if len(parts) > 0 && strings.TrimSpace(parts[len(parts)-1]) == "" {
		parts = parts[:len(parts)-1]
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func splitTableCells(line string) []string {
	parts := []string{}
	var cell strings.Builder
	codeTicks := 0

	for i := 0; i < len(line); {
		switch line[i] {
		case '\\':
			cell.WriteByte(line[i])
			i++
			if i < len(line) {
				cell.WriteByte(line[i])
				i++
			}
		case '`':
			ticks := countBackticks(line[i:])
			cell.WriteString(line[i : i+ticks])
			if codeTicks == ticks {
				codeTicks = 0
			} else if codeTicks == 0 {
				codeTicks = ticks
			}
			i += ticks
		case '|':
			if codeTicks == 0 {
				parts = append(parts, cell.String())
				cell.Reset()
				i++
				continue
			}
			cell.WriteByte(line[i])
			i++
		default:
			cell.WriteByte(line[i])
			i++
		}
	}
	parts = append(parts, cell.String())
	return parts
}

func countBackticks(s string) int {
	count := 0
	for count < len(s) && s[count] == '`' {
		count++
	}
	return count
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
	return len(splitTableCells(strings.TrimSpace(line))) > 1
}

func isTableSeparator(line string) bool {
	if !isTableRow(line) {
		return false
	}
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
