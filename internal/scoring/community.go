package scoring

import "math"

// CalcCommunitySignal scores community signal (0-100).
// Components: stars (30), downloads (30), contributors (20), dependents (20).
func CalcCommunitySignal(data map[string]any) float64 {
	var total float64

	// Stars (log scale): 0-30
	stars := toFloat(data["stars"])
	total += scoreLogScale(stars, 30, 100000)

	// Weekly downloads (log scale): 0-30
	downloads := toFloat(data["weekly_downloads"])
	total += scoreLogScale(downloads, 30, 10000000)

	// Contributors: 0-20
	contributors := toInt(data["contributors_count"])
	switch {
	case contributors >= 10:
		total += 20
	case contributors >= 5:
		total += 15
	case contributors >= 2:
		total += 10
	case contributors >= 1:
		total += 5
	}

	// Forks as proxy for dependents: 0-20
	forks := toFloat(data["forks"])
	total += scoreLogScale(forks, 20, 10000)

	return clamp(total, 0, 100)
}

// scoreLogScale maps a value to a score using log scale.
// maxScore is the maximum points, refValue is the "100%" reference value.
func scoreLogScale(value float64, maxScore float64, refValue float64) float64 {
	if value <= 0 {
		return 0
	}
	ratio := math.Log10(value+1) / math.Log10(refValue+1)
	score := ratio * maxScore
	if score > maxScore {
		score = maxScore
	}
	return score
}
