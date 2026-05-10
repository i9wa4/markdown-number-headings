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
		return Format(input, opts)
	}
}

func DocumentFormatting(opts Options) Pass {
	return func(input string) string {
		return Apply(input, HeadingNumbering(opts), TableAlignment(), HeadingSpacing())
	}
}

func DocumentFormattingWithoutHeadingNumbering() Pass {
	return func(input string) string {
		return Apply(input, TableAlignment(), HeadingSpacing())
	}
}

func HeadingNumberRemoval() Pass {
	return Remove
}

func TableAlignment() Pass {
	return FormatTables
}

func HeadingSpacing() Pass {
	return NormalizeHeadingSpacing
}
