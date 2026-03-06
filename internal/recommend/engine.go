package recommend

import (
	"strings"
)

// MatchResult holds the recommendation engine's output.
type MatchResult struct {
	MatchedCategories []string `json:"matchedCategories"`
	SortField         string   `json:"sortField"`
	MatchReason       string   `json:"matchReason"`
}

// Match takes a task description and optional preference, returns matched categories.
func Match(task string, prefer string) MatchResult {
	task = strings.ToLower(strings.TrimSpace(task))

	matched := matchCategories(task)

	sortField := "overall_score" // default
	if field, ok := preferSortFields[strings.ToLower(prefer)]; ok {
		sortField = field
	}

	reason := buildReason(task, matched)

	return MatchResult{
		MatchedCategories: matched,
		SortField:         sortField,
		MatchReason:       reason,
	}
}

func matchCategories(task string) []string {
	scores := make(map[string]int)

	for category, keywords := range categoryKeywords {
		for _, kw := range keywords {
			if strings.Contains(task, kw) {
				scores[category]++
			}
		}
	}

	if len(scores) == 0 {
		return []string{}
	}

	// Find the best-matching categories (top matches)
	maxScore := 0
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}

	var result []string
	for cat, score := range scores {
		// Include categories that scored at least half of the max
		if score*2 >= maxScore {
			result = append(result, cat)
		}
	}

	return result
}

func buildReason(task string, matched []string) string {
	if len(matched) == 0 {
		return "No category matched for task: " + task
	}

	// Find which keywords matched
	var matchedKeywords []string
	for _, cat := range matched {
		for _, kw := range categoryKeywords[cat] {
			if strings.Contains(task, kw) {
				matchedKeywords = append(matchedKeywords, kw)
			}
		}
	}

	return "Matched categories [" + strings.Join(matched, ", ") +
		"] from keywords [" + strings.Join(matchedKeywords, ", ") + "] in task"
}
