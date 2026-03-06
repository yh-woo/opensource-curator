package scoring

// CalcDocQuality scores documentation quality (0-100).
// Components: README (25), docs site (25), changelog (15), license (15), examples (20).
func CalcDocQuality(data map[string]any) float64 {
	var total float64

	// README existence & length: 0-25
	hasReadme := toBool(data["has_readme"])
	readmeLen := toInt(data["readme_length"])
	if !hasReadme {
		total += 0
	} else if readmeLen >= 2000 {
		total += 25
	} else if readmeLen >= 500 {
		total += 15
	} else {
		total += 5
	}

	// Docs site (homepage_url as proxy): 0-25
	homepage := toString(data["homepage"])
	if homepage != "" {
		total += 25
	} else {
		total += 10 // README-only still gets some credit
	}

	// LICENSE: 0-15
	license := toString(data["license"])
	switch license {
	case "MIT", "Apache-2.0", "BSD-2-Clause", "BSD-3-Clause", "ISC":
		total += 15
	case "":
		total += 0
	default:
		total += 10
	}

	// Has wiki: 0-15 (proxy for changelog/docs)
	if toBool(data["has_wiki"]) {
		total += 15
	}

	// Code examples in README: 0-20
	examples := toInt(data["readme_has_examples"])
	if examples >= 3 {
		total += 20
	} else if examples >= 1 {
		total += 10
	}

	return clamp(total, 0, 100)
}
