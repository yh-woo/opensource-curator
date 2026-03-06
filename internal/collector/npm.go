package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// NpmCollector collects package metadata from the npm registry.
type NpmCollector struct {
	client *http.Client
}

func NewNpmCollector() *NpmCollector {
	return &NpmCollector{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (n *NpmCollector) Name() string { return "npm" }

func (n *NpmCollector) Collect(ctx context.Context, githubRepo string, registry string, packageName string) (*Result, error) {
	start := time.Now()

	if registry != "npm" {
		return nil, fmt.Errorf("npm collector only handles npm registry, got %s", registry)
	}

	// Fetch package metadata
	pkgData, err := n.fetchPackage(ctx, packageName)
	if err != nil {
		return &Result{Err: err, Duration: time.Since(start)}, err
	}

	// Fetch download counts
	downloads, err := n.fetchDownloads(ctx, packageName)
	if err != nil {
		slog.Warn("failed to fetch downloads", "package", packageName, "err", err)
		downloads = map[string]any{"downloads": 0}
	}

	// Extract info from package metadata
	data := map[string]any{
		"weekly_downloads": downloads["downloads"],
	}

	// Extract latest version info
	if distTags, ok := pkgData["dist-tags"].(map[string]any); ok {
		if latest, ok := distTags["latest"].(string); ok {
			data["latest_version"] = latest

			// Get version-specific info
			if versions, ok := pkgData["versions"].(map[string]any); ok {
				if versionData, ok := versions[latest].(map[string]any); ok {
					data["has_types"] = hasTypeScript(versionData)
					data["dependencies_count"] = countDeps(versionData, "dependencies")
					data["dev_dependencies_count"] = countDeps(versionData, "devDependencies")
					data["deprecated"] = versionData["deprecated"] != nil
				}
			}
		}
	}

	// Extract time info (publish dates)
	if timeInfo, ok := pkgData["time"].(map[string]any); ok {
		if created, ok := timeInfo["created"].(string); ok {
			data["created_at"] = created
		}
		if modified, ok := timeInfo["modified"].(string); ok {
			data["modified_at"] = modified
		}
		// Find latest version publish date
		if latest, ok := data["latest_version"].(string); ok {
			if publishDate, ok := timeInfo[latest].(string); ok {
				data["latest_version_date"] = publishDate
			}
		}
	}

	// License
	if license, ok := pkgData["license"].(string); ok {
		data["license"] = license
	}

	// Description
	if desc, ok := pkgData["description"].(string); ok {
		data["description"] = desc
	}

	return &Result{
		Data:     data,
		Duration: time.Since(start),
	}, nil
}

func (n *NpmCollector) fetchPackage(ctx context.Context, packageName string) (map[string]any, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", packageName)
	body, err := n.doRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse package data: %w", err)
	}
	return data, nil
}

func (n *NpmCollector) fetchDownloads(ctx context.Context, packageName string) (map[string]any, error) {
	url := fmt.Sprintf("https://api.npmjs.org/downloads/point/last-week/%s", packageName)
	body, err := n.doRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse downloads data: %w", err)
	}
	return data, nil
}

func (n *NpmCollector) doRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Politeness delay
	time.Sleep(500 * time.Millisecond)

	resp, err := n.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

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

func hasTypeScript(versionData map[string]any) bool {
	// Check for "types" or "typings" field
	if _, ok := versionData["types"]; ok {
		return true
	}
	if _, ok := versionData["typings"]; ok {
		return true
	}
	return false
}

func countDeps(versionData map[string]any, field string) int {
	if deps, ok := versionData[field].(map[string]any); ok {
		return len(deps)
	}
	return 0
}
