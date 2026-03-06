package scoring

import (
	"time"
)

// CalcMaintenanceHealth scores maintenance health (0-100).
// Components: recent commits (30), release frequency (25), issue response (25), PR merge speed (20).
func CalcMaintenanceHealth(data map[string]any) float64 {
	var total float64

	// Recent commits in last 90 days: 0-30
	total += scoreRecentCommits(toInt(data["recent_commits"]))

	// Release frequency (releases in the last year): 0-25
	total += scoreReleaseFrequency(toInt(data["releases_count"]))

	// Pushed at recency (proxy for activity): 0-25
	total += scorePushedAt(toString(data["pushed_at"]))

	// Contributors count (proxy for PR activity): 0-20
	total += scoreContributors(toInt(data["contributors_count"]))

	return clamp(total, 0, 100)
}

func scoreRecentCommits(count int) float64 {
	switch {
	case count >= 10:
		return 30
	case count >= 5:
		return 20
	case count >= 1:
		return 10
	default:
		return 0
	}
}

func scoreReleaseFrequency(count int) float64 {
	switch {
	case count >= 4:
		return 25
	case count >= 2:
		return 15
	case count >= 1:
		return 10
	default:
		return 0
	}
}

func scorePushedAt(pushedAt string) float64 {
	if pushedAt == "" {
		return 0
	}
	t, err := time.Parse(time.RFC3339, pushedAt)
	if err != nil {
		return 0
	}
	days := time.Since(t).Hours() / 24
	switch {
	case days <= 7:
		return 25
	case days <= 30:
		return 20
	case days <= 90:
		return 15
	case days <= 180:
		return 10
	case days <= 365:
		return 5
	default:
		return 0
	}
}

func scoreContributors(count int) float64 {
	switch {
	case count >= 20:
		return 20
	case count >= 10:
		return 15
	case count >= 5:
		return 10
	case count >= 2:
		return 5
	default:
		return 0
	}
}
