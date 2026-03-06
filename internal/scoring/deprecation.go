package scoring

import "time"

// CalcDeprecationSafety scores how far a library is from deprecation (0-100, higher=safer).
// Only called for non-deprecated libraries (pre-filter passed).
// Components: recent commit activity (40), recent release (30), bus factor (30).
func CalcDeprecationSafety(data map[string]any) float64 {
	var total float64

	// Recent commit activity (1 year): 0-40
	if pushedAt := toString(data["pushed_at"]); pushedAt != "" {
		if t, err := time.Parse(time.RFC3339, pushedAt); err == nil {
			months := time.Since(t).Hours() / (24 * 30)
			switch {
			case months <= 6:
				total += 40
			case months <= 12:
				total += 20
			}
		}
	}

	// Recent release (2 years): 0-30
	if releaseDate := toString(data["latest_release_date"]); releaseDate != "" {
		if t, err := time.Parse(time.RFC3339, releaseDate); err == nil {
			months := time.Since(t).Hours() / (24 * 30)
			switch {
			case months <= 12:
				total += 30
			case months <= 24:
				total += 15
			}
		}
	}

	// Bus factor (contributors): 0-30
	contributors := toInt(data["contributors_count"])
	switch {
	case contributors >= 3:
		total += 30
	case contributors >= 2:
		total += 20
	case contributors >= 1:
		total += 10
	}

	return clamp(total, 0, 100)
}
