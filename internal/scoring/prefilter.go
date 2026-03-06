package scoring

import "strings"

// PreFilterResult indicates whether a library should be filtered out as deprecated.
type PreFilterResult struct {
	Deprecated bool
	Reason     string
}

// PreFilter checks if a library is deprecated before scoring.
// If deprecated, overall_score should be set to 0 and scoring skipped.
func PreFilter(data map[string]any) PreFilterResult {
	// Check GitHub archived status
	if archived, ok := data["archived"].(bool); ok && archived {
		return PreFilterResult{Deprecated: true, Reason: "GitHub repository is archived"}
	}

	// Check npm deprecated flag
	if deprecated, ok := data["deprecated"].(bool); ok && deprecated {
		return PreFilterResult{Deprecated: true, Reason: "Package is marked as deprecated on npm"}
	}

	// Check README for deprecation keywords
	if readmeDeprecated, ok := data["readme_has_deprecated"].(bool); ok && readmeDeprecated {
		return PreFilterResult{Deprecated: true, Reason: "README indicates project is deprecated or unmaintained"}
	}

	// Also check raw readme content if available
	if readme, ok := data["readme"].(string); ok {
		lower := strings.ToLower(readme)
		if strings.Contains(lower, "deprecated") || strings.Contains(lower, "unmaintained") || strings.Contains(lower, "no longer maintained") {
			return PreFilterResult{Deprecated: true, Reason: "README indicates project is deprecated or unmaintained"}
		}
	}

	return PreFilterResult{Deprecated: false}
}
