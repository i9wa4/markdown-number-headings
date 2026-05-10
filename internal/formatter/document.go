package formatter

type document struct {
	lines       []string
	trailing    bool
	fenced      []bool
	fencedDirty bool
}

type documentPass func(*document)

func applyDocument(input string, passes ...documentPass) string {
	doc := newDocument(input)
	for _, pass := range passes {
		if pass == nil {
			continue
		}
		pass(doc)
	}
	return doc.String()
}

func newDocument(input string) *document {
	lines, trailing := splitLines(input)
	doc := &document{
		lines:    lines,
		trailing: trailing,
	}
	doc.refreshFences()
	return doc
}

func (doc *document) String() string {
	return joinLines(doc.lines, doc.trailing)
}

func (doc *document) refreshFences() {
	doc.fenced = mapFenceLines(doc.lines)
	doc.fencedDirty = false
}

func (doc *document) setLines(lines []string) {
	doc.lines = lines
	doc.fencedDirty = true
}

func (doc *document) inFence(index int) bool {
	if doc.fencedDirty {
		doc.refreshFences()
	}
	return index >= 0 && index < len(doc.fenced) && doc.fenced[index]
}

func normalizeFinalNewlinePass() documentPass {
	return func(doc *document) {
		for len(doc.lines) > 0 && isBlankLine(doc.lines[len(doc.lines)-1]) {
			doc.lines = doc.lines[:len(doc.lines)-1]
			doc.fencedDirty = true
		}
		doc.trailing = len(doc.lines) > 0
	}
}

func mapFenceLines(lines []string) []bool {
	fenced := make([]bool, len(lines))
	inFence := false
	fenceMarker := byte(0)

	for i, line := range lines {
		if marker, ok := fenceStart(line); ok {
			if !inFence {
				inFence = true
				fenceMarker = marker
				fenced[i] = true
				continue
			}
			if marker == fenceMarker {
				fenced[i] = true
				inFence = false
				fenceMarker = 0
				continue
			}
		}
		if inFence {
			fenced[i] = true
		}
	}

	return fenced
}
