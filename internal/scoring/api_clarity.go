package scoring

// CalcAPIClarity scores API clarity (0-100).
// Components: type safety (40), dependency count (20), code examples (20), doc ratio (20).
func CalcAPIClarity(data map[string]any) float64 {
	var total float64

	// TypeScript types / Python type hints: 0-40
	total += scoreTypeSafety(toBool(data["has_types"]))

	// Dependency count (fewer = better): 0-20
	total += scoreDependencyCount(toInt(data["dependencies_count"]))

	// README code examples: 0-20
	total += scoreCodeExamples(toInt(data["readme_has_examples"]))

	// API documentation ratio (heuristic: README length as proxy): 0-20
	total += scoreDocRatio(toInt(data["readme_length"]))

	return clamp(total, 0, 100)
}

func scoreTypeSafety(hasTypes bool) float64 {
	if hasTypes {
		return 40
	}
	return 0
}

func scoreDependencyCount(count int) float64 {
	switch {
	case count <= 3:
		return 20
	case count <= 10:
		return 15
	case count <= 20:
		return 10
	default:
		return 5
	}
}

func scoreCodeExamples(codeBlockCount int) float64 {
	switch {
	case codeBlockCount >= 5:
		return 20
	case codeBlockCount >= 3:
		return 15
	case codeBlockCount >= 1:
		return 10
	default:
		return 0
	}
}

func scoreDocRatio(readmeLength int) float64 {
	switch {
	case readmeLength >= 5000:
		return 20
	case readmeLength >= 2000:
		return 15
	case readmeLength >= 500:
		return 10
	default:
		return 5
	}
}
