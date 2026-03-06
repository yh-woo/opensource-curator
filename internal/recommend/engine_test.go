package recommend

import (
	"testing"
)

func TestMatch_HTTPClient(t *testing.T) {
	result := Match("make http requests in typescript", "lightweight")
	if len(result.MatchedCategories) == 0 {
		t.Fatal("expected at least one matched category")
	}
	found := false
	for _, cat := range result.MatchedCategories {
		if cat == "http-client" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected http-client in matched categories, got %v", result.MatchedCategories)
	}
	if result.SortField != "dependencies_count" {
		t.Errorf("expected sort by dependencies_count for lightweight, got %s", result.SortField)
	}
}

func TestMatch_Database(t *testing.T) {
	result := Match("need a database orm for postgres", "stable")
	found := false
	for _, cat := range result.MatchedCategories {
		if cat == "orm-database" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected orm-database, got %v", result.MatchedCategories)
	}
	if result.SortField != "maintenance_health" {
		t.Errorf("expected maintenance_health for stable, got %s", result.SortField)
	}
}

func TestMatch_AI(t *testing.T) {
	result := Match("build an ai chatbot with llm", "")
	found := false
	for _, cat := range result.MatchedCategories {
		if cat == "ai-ml-sdk" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected ai-ml-sdk, got %v", result.MatchedCategories)
	}
	if result.SortField != "overall_score" {
		t.Errorf("expected overall_score as default sort, got %s", result.SortField)
	}
}

func TestMatch_NoMatch(t *testing.T) {
	result := Match("xyzzy foobar baz", "")
	if len(result.MatchedCategories) != 0 {
		t.Errorf("expected no match, got %v", result.MatchedCategories)
	}
}

func TestMatch_PreferSecure(t *testing.T) {
	result := Match("authentication library", "secure")
	if result.SortField != "security_posture" {
		t.Errorf("expected security_posture for secure, got %s", result.SortField)
	}
}

func TestMatch_Reason(t *testing.T) {
	result := Match("http request library", "")
	if result.MatchReason == "" {
		t.Error("expected non-empty match reason")
	}
}
