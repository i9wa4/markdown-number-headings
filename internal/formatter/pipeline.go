package formatter

type Pass func(string) string

func Apply(input string, passes ...Pass) string {
	for _, pass := range passes {
		if pass == nil {
			continue
		}
		input = pass(input)
	}
	return input
}

func HeadingNumbering(opts Options) Pass {
	return func(input string) string {
		return applyDocument(input, headingNumberingPass(opts))
	}
}

func DocumentFormatting(opts Options) Pass {
	return func(input string) string {
		return applyDocument(input, headingNumberingPass(opts), tableAlignmentPass(), headingSpacingPass(), normalizeFinalNewlinePass())
	}
}

func DocumentFormattingWithoutHeadingNumbering() Pass {
	return func(input string) string {
		return applyDocument(input, tableAlignmentPass(), headingSpacingPass(), normalizeFinalNewlinePass())
	}
}

func HeadingNumberRemoval() Pass {
	return func(input string) string {
		return applyDocument(input, headingNumberRemovalPass())
	}
}

func TableAlignment() Pass {
	return func(input string) string {
		return applyDocument(input, tableAlignmentPass())
	}
}

func HeadingSpacing() Pass {
	return func(input string) string {
		return applyDocument(input, headingSpacingPass())
	}
}
