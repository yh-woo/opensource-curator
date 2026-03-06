package scoring

import "time"

// CalcSecurityPosture scores security posture (0-100).
// MVP: based on dependency count and license (full CVE analysis in Phase 2).
func CalcSecurityPosture(data map[string]any) float64 {
	var total float64

	// Dependency count as proxy for attack surface: 0-50
	depCount := toInt(data["dependencies_count"])
	switch {
	case depCount == 0:
		total += 50
	case depCount <= 3:
		total += 40
	case depCount <= 10:
		total += 30
	case depCount <= 20:
		total += 20
	default:
		total += 10
	}

	// License presence: 0-25
	license := toString(data["license"])
	switch license {
	case "MIT", "Apache-2.0", "BSD-2-Clause", "BSD-3-Clause", "ISC":
		total += 25
	case "":
		total += 0
	default:
		total += 15
	}

	// Recency of updates: 0-25
	if pushedAt := toString(data["pushed_at"]); pushedAt != "" {
		if t, err := time.Parse(time.RFC3339, pushedAt); err == nil {
			months := time.Since(t).Hours() / (24 * 30)
			switch {
			case months <= 3:
				total += 25
			case months <= 6:
				total += 15
			case months <= 12:
				total += 5
			}
		}
	}

	return clamp(total, 0, 100)
}
