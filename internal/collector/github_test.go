package collector

import (
	"testing"
)

func TestCountCodeBlocks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"no blocks", "hello world", 0},
		{"one block", "```js\nconsole.log('hi')\n```", 1},
		{"two blocks", "```js\nfoo\n```\ntext\n```ts\nbar\n```", 2},
		{"empty", "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countCodeBlocks(tt.input)
			if got != tt.expected {
				t.Errorf("countCodeBlocks() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestContainsDeprecated(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"no match", "this is a great library", false},
		{"deprecated", "This library is DEPRECATED", true},
		{"unmaintained", "⚠️ This project is unmaintained", true},
		{"no longer maintained", "No longer maintained. Use X instead.", true},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsDeprecated(tt.input)
			if got != tt.expected {
				t.Errorf("containsDeprecated() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExtractLicense(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected string
	}{
		{"MIT", map[string]any{"license": map[string]any{"spdx_id": "MIT"}}, "MIT"},
		{"nil license", map[string]any{"license": nil}, ""},
		{"no license key", map[string]any{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractLicense(tt.input)
			if got != tt.expected {
				t.Errorf("extractLicense() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDecodeBase64Content(t *testing.T) {
	// "Hello World" in base64
	got, err := decodeBase64Content("SGVsbG8gV29ybGQ=")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Hello World" {
		t.Errorf("got %q, want %q", got, "Hello World")
	}

	// With newlines (GitHub style)
	got, err = decodeBase64Content("SGVs\nbG8g\nV29ybGQ=")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Hello World" {
		t.Errorf("got %q, want %q", got, "Hello World")
	}
}
