package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// NpmSearchResult represents a single result from npm search API.
type NpmSearchResult struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Links       struct {
		Repository string `json:"repository"`
	} `json:"links"`
}

type npmSearchResponse struct {
	Objects []struct {
		Package NpmSearchResult `json:"package"`
	} `json:"objects"`
	Total int `json:"total"`
}

// NpmDiscoverer finds popular npm packages by keyword search.
type NpmDiscoverer struct {
	client *http.Client
}

func NewNpmDiscoverer() *NpmDiscoverer {
	return &NpmDiscoverer{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// SearchByKeywords queries npm registry for packages matching keywords.
// Returns up to `size` results sorted by popularity.
func (d *NpmDiscoverer) SearchByKeywords(keywords string, size int) ([]NpmSearchResult, error) {
	u := fmt.Sprintf(
		"https://registry.npmjs.org/-/v1/search?text=%s&popularity=1.0&size=%d",
		url.QueryEscape(keywords), size,
	)

	resp, err := d.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("npm search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("npm search returned %d", resp.StatusCode)
	}

	var result npmSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode npm search: %w", err)
	}

	var packages []NpmSearchResult
	for _, obj := range result.Objects {
		packages = append(packages, obj.Package)
	}
	return packages, nil
}

// ExtractGithubRepo extracts "owner/repo" from a GitHub URL.
func ExtractGithubRepo(repoURL string) string {
	if repoURL == "" {
		return ""
	}
	// Handle various formats:
	// https://github.com/owner/repo
	// https://github.com/owner/repo.git
	// git+https://github.com/owner/repo.git
	// git://github.com/owner/repo.git
	repoURL = strings.TrimPrefix(repoURL, "git+")
	repoURL = strings.TrimPrefix(repoURL, "git://")
	repoURL = strings.TrimSuffix(repoURL, ".git")

	// Parse as URL
	parsed, err := url.Parse(repoURL)
	if err != nil {
		return ""
	}

	if !strings.Contains(parsed.Host, "github.com") {
		return ""
	}

	path := strings.TrimPrefix(parsed.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) >= 2 {
		return parts[0] + "/" + parts[1]
	}
	return ""
}
