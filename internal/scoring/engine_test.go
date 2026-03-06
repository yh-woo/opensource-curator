package scoring

import (
	"testing"
	"time"
)

func TestPreFilter_Archived(t *testing.T) {
	result := PreFilter(map[string]any{"archived": true})
	if !result.Deprecated {
		t.Error("expected archived repo to be deprecated")
	}
}

func TestPreFilter_NpmDeprecated(t *testing.T) {
	result := PreFilter(map[string]any{"deprecated": true})
	if !result.Deprecated {
		t.Error("expected npm deprecated to be filtered")
	}
}

func TestPreFilter_ReadmeDeprecated(t *testing.T) {
	result := PreFilter(map[string]any{"readme_has_deprecated": true})
	if !result.Deprecated {
		t.Error("expected readme deprecated to be filtered")
	}
}

func TestPreFilter_HealthyLibrary(t *testing.T) {
	result := PreFilter(map[string]any{
		"archived":             false,
		"deprecated":           false,
		"readme_has_deprecated": false,
	})
	if result.Deprecated {
		t.Error("expected healthy library to pass pre-filter")
	}
}

func TestCalculate_DeprecatedLibrary(t *testing.T) {
	data := map[string]any{"archived": true}
	result := Calculate(data)

	if result.Overall != 0 {
		t.Errorf("expected overall=0 for deprecated, got %f", result.Overall)
	}
	if !result.Deprecated {
		t.Error("expected deprecated=true")
	}
	if result.DeprecationReason == "" {
		t.Error("expected deprecation reason")
	}
	if result.Version != Version {
		t.Errorf("expected version %s, got %s", Version, result.Version)
	}
}

func TestCalculate_HealthyLibrary(t *testing.T) {
	now := time.Now()
	data := map[string]any{
		"archived":             false,
		"deprecated":           false,
		"readme_has_deprecated": false,
		"stars":                5000.0,
		"forks":                500.0,
		"watchers":             100.0,
		"recent_commits":       15,
		"releases_count":       6,
		"contributors_count":   12,
		"pushed_at":            now.Add(-48 * time.Hour).Format(time.RFC3339),
		"latest_release_date":  now.Add(-30 * 24 * time.Hour).Format(time.RFC3339),
		"has_readme":           true,
		"readme_length":        3000,
		"readme_has_examples":  5,
		"has_types":            true,
		"dependencies_count":   3,
		"weekly_downloads":     100000.0,
		"license":              "MIT",
		"homepage":             "https://example.com",
		"has_wiki":             true,
	}

	result := Calculate(data)

	if result.Deprecated {
		t.Error("expected not deprecated")
	}
	if result.Overall <= 0 {
		t.Errorf("expected positive overall score, got %f", result.Overall)
	}
	if result.Overall > 100 {
		t.Errorf("expected overall <= 100, got %f", result.Overall)
	}

	// Verify individual metrics are computed
	if result.MaintenanceHealth <= 0 {
		t.Error("expected positive maintenance score")
	}
	if result.APIClarity <= 0 {
		t.Error("expected positive API clarity score")
	}
	if result.DocQuality <= 0 {
		t.Error("expected positive doc quality score")
	}
	if result.SecurityPosture <= 0 {
		t.Error("expected positive security score")
	}
	if result.CommunitySignal <= 0 {
		t.Error("expected positive community score")
	}
	if result.DeprecationSafety <= 0 {
		t.Error("expected positive deprecation safety score")
	}

	// Verify weights sum correctly (overall should be weighted average)
	expected := result.MaintenanceHealth*0.25 +
		result.APIClarity*0.20 +
		result.DocQuality*0.15 +
		result.SecurityPosture*0.15 +
		result.CommunitySignal*0.15 +
		result.DeprecationSafety*0.10
	expected = round2(expected)

	if result.Overall != expected {
		t.Errorf("overall %f != weighted average %f", result.Overall, expected)
	}
}

func TestWeightsSum(t *testing.T) {
	sum := WeightMaintenanceHealth + WeightAPIClarity + WeightDocQuality +
		WeightSecurityPosture + WeightCommunitySignal + WeightDeprecationSafety
	if sum != 1.0 {
		t.Errorf("weights sum to %f, expected 1.0", sum)
	}
}

func TestCalcMaintenanceHealth(t *testing.T) {
	now := time.Now()
	data := map[string]any{
		"recent_commits":     15,
		"releases_count":     5,
		"pushed_at":          now.Add(-24 * time.Hour).Format(time.RFC3339),
		"contributors_count": 25,
	}
	score := CalcMaintenanceHealth(data)
	if score != 100 {
		t.Errorf("expected 100 for highly active repo, got %f", score)
	}
}

func TestCalcMaintenanceHealth_Inactive(t *testing.T) {
	data := map[string]any{
		"recent_commits":     0,
		"releases_count":     0,
		"pushed_at":          "",
		"contributors_count": 0,
	}
	score := CalcMaintenanceHealth(data)
	if score != 0 {
		t.Errorf("expected 0 for inactive repo, got %f", score)
	}
}

func TestCommunitySignal_LogScale(t *testing.T) {
	s1 := scoreLogScale(100, 30, 100000)
	s2 := scoreLogScale(10000, 30, 100000)
	s3 := scoreLogScale(100000, 30, 100000)

	if s1 >= s2 || s2 >= s3 {
		t.Errorf("expected monotonically increasing: %f < %f < %f", s1, s2, s3)
	}
	if s3 != 30 {
		t.Errorf("expected max score at ref value, got %f", s3)
	}
}
