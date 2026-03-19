package collector

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// GitHubCollector collects repository metadata from the GitHub API.
type GitHubCollector struct {
	token  string
	client *http.Client
}

func NewGitHubCollector(token string) *GitHubCollector {
	return &GitHubCollector{
		token: token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (g *GitHubCollector) Name() string { return "github" }

func (g *GitHubCollector) Collect(ctx context.Context, githubRepo string, registry string, packageName string) (*Result, error) {
	start := time.Now()

	parts := strings.SplitN(githubRepo, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid github_repo format: %s (expected owner/repo)", githubRepo)
	}
	owner, repo := parts[0], parts[1]

	// Fetch repo metadata
	repoData, err := g.fetchRepo(ctx, owner, repo)
	if err != nil {
		return &Result{Err: err, Duration: time.Since(start)}, err
	}

	// Fetch recent commits (last 90 days)
	commits, err := g.fetchRecentCommits(ctx, owner, repo)
	if err != nil {
		slog.Warn("failed to fetch commits", "repo", githubRepo, "err", err)
		commits = []map[string]any{}
	}

	// Fetch releases
	releases, err := g.fetchReleases(ctx, owner, repo)
	if err != nil {
		slog.Warn("failed to fetch releases", "repo", githubRepo, "err", err)
		releases = []map[string]any{}
	}

	// Fetch contributors
	contributors, err := g.fetchContributors(ctx, owner, repo)
	if err != nil {
		slog.Warn("failed to fetch contributors", "repo", githubRepo, "err", err)
		contributors = []map[string]any{}
	}

	data := map[string]any{
		"stars":             repoData["stargazers_count"],
		"forks":             repoData["forks_count"],
		"watchers":          repoData["watchers_count"],
		"open_issues":       repoData["open_issues_count"],
		"archived":          repoData["archived"],
		"license":           extractLicense(repoData),
		"description":       repoData["description"],
		"homepage":          repoData["homepage"],
		"default_branch":    repoData["default_branch"],
		"created_at":        repoData["created_at"],
		"updated_at":        repoData["updated_at"],
		"pushed_at":         repoData["pushed_at"],
		"has_issues":        repoData["has_issues"],
		"has_wiki":          repoData["has_wiki"],
		"size":              repoData["size"],
		"recent_commits":    len(commits),
		"releases_count":    len(releases),
		"contributors_count": len(contributors),
	}

	// Extract latest release info
	if len(releases) > 0 {
		data["latest_release_tag"] = releases[0]["tag_name"]
		data["latest_release_date"] = releases[0]["published_at"]
	}

	// Check README existence via repo size and has_wiki (heuristic)
	readme, err := g.fetchReadme(ctx, owner, repo)
	if err == nil {
		data["readme_length"] = len(readme)
		data["has_readme"] = true
		data["readme_has_examples"] = countCodeBlocks(readme)
		data["readme_has_deprecated"] = containsDeprecated(readme)
	} else {
		data["readme_length"] = 0
		data["has_readme"] = false
	}

	return &Result{
		Data:     data,
		Duration: time.Since(start),
	}, nil
}

func (g *GitHubCollector) doRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limited (status %d)", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found: %s", url)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return body, nil
}

func (g *GitHubCollector) fetchRepo(ctx context.Context, owner, repo string) (map[string]any, error) {
	body, err := g.doRequest(ctx, fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo))
	if err != nil {
		return nil, err
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (g *GitHubCollector) fetchRecentCommits(ctx context.Context, owner, repo string) ([]map[string]any, error) {
	since := time.Now().AddDate(0, 0, -90).Format(time.RFC3339)
	body, err := g.doRequest(ctx, fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/commits?since=%s&per_page=100", owner, repo, since))
	if err != nil {
		return nil, err
	}
	var data []map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (g *GitHubCollector) fetchReleases(ctx context.Context, owner, repo string) ([]map[string]any, error) {
	body, err := g.doRequest(ctx, fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/releases?per_page=10", owner, repo))
	if err != nil {
		return nil, err
	}
	var data []map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (g *GitHubCollector) fetchContributors(ctx context.Context, owner, repo string) ([]map[string]any, error) {
	body, err := g.doRequest(ctx, fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contributors?per_page=100", owner, repo))
	if err != nil {
		return nil, err
	}
	var data []map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (g *GitHubCollector) fetchReadme(ctx context.Context, owner, repo string) (string, error) {
	body, err := g.doRequest(ctx, fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/readme", owner, repo))
	if err != nil {
		return "", err
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	// README is base64-encoded in the API response
	if content, ok := data["content"].(string); ok {
		decoded, err := decodeBase64Content(content)
		if err != nil {
			return "", err
		}
		return decoded, nil
	}
	return "", fmt.Errorf("no content in readme response")
}

func decodeBase64Content(s string) (string, error) {
	// GitHub returns base64 with newlines, strip them
	cleaned := strings.ReplaceAll(s, "\n", "")
	b, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func extractLicense(repo map[string]any) string {
	if lic, ok := repo["license"].(map[string]any); ok {
		if spdx, ok := lic["spdx_id"].(string); ok {
			return spdx
		}
	}
	return ""
}

func countCodeBlocks(readme string) int {
	count := 0
	for i := 0; i < len(readme)-2; i++ {
		if readme[i] == '`' && readme[i+1] == '`' && readme[i+2] == '`' {
			count++
			i += 2
		}
	}
	return count / 2 // opening + closing = 1 block
}

// containsDeprecated is intentionally disabled.
// README keyword matching produces too many false positives.
// We rely solely on GitHub archived flag and npm deprecated field.
// TODO: Replace with LLM-based contextual analysis.
func containsDeprecated(_ string) bool {
	return false
}
