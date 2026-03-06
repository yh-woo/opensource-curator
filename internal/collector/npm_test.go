package collector

import (
	"testing"
)

func TestHasTypeScript(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected bool
	}{
		{"types field", map[string]any{"types": "index.d.ts"}, true},
		{"typings field", map[string]any{"typings": "index.d.ts"}, true},
		{"no types", map[string]any{"main": "index.js"}, false},
		{"empty", map[string]any{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasTypeScript(tt.input)
			if got != tt.expected {
				t.Errorf("hasTypeScript() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCountDeps(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		field    string
		expected int
	}{
		{"3 deps", map[string]any{"dependencies": map[string]any{"a": "1", "b": "2", "c": "3"}}, "dependencies", 3},
		{"no deps field", map[string]any{}, "dependencies", 0},
		{"empty deps", map[string]any{"dependencies": map[string]any{}}, "dependencies", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countDeps(tt.input, tt.field)
			if got != tt.expected {
				t.Errorf("countDeps() = %d, want %d", got, tt.expected)
			}
		})
	}
}
