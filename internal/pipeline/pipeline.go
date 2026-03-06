package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/woo/opensource-curator/internal/collector"
	"github.com/woo/opensource-curator/internal/db"
	"github.com/woo/opensource-curator/internal/scoring"
)

// Runner orchestrates data collection and scoring for all libraries.
type Runner struct {
	queries *db.Queries
	github  *collector.GitHubCollector
	npm     *collector.NpmCollector
}

func NewRunner(queries *db.Queries, githubToken string) *Runner {
	return &Runner{
		queries: queries,
		github:  collector.NewGitHubCollector(githubToken),
		npm:     collector.NewNpmCollector(),
	}
}

// RunAll collects data and scores all libraries.
func (r *Runner) RunAll(ctx context.Context) error {
	run, err := r.queries.CreateCollectionRun(ctx, "manual")
	if err != nil {
		return err
	}
	slog.Info("collection run started", "run_id", run.ID)

	libs, err := r.queries.ListAllLibraries(ctx)
	if err != nil {
		return err
	}

	var succeeded, failed int32
	for i, lib := range libs {
		slog.Info("processing library",
			"index", i+1,
			"total", len(libs),
			"name", lib.Name,
			"registry", lib.Registry,
		)

		err := r.processLibrary(ctx, run.ID, lib)
		if err != nil {
			slog.Error("failed to process library", "name", lib.Name, "err", err)
			failed++
		} else {
			succeeded++
		}
	}

	total := succeeded + failed
	var successRate pgtype.Numeric
	if total > 0 {
		rate := float64(succeeded) / float64(total) * 100
		successRate = floatToNumeric(rate)
	}

	err = r.queries.UpdateCollectionRun(ctx, db.UpdateCollectionRunParams{
		ID:                 run.ID,
		Status:             "completed",
		LibrariesProcessed: total,
		LibrariesSucceeded: succeeded,
		LibrariesFailed:    failed,
		SuccessRate:        successRate,
	})
	if err != nil {
		slog.Error("failed to update collection run", "err", err)
	}

	slog.Info("collection run completed",
		"total", total,
		"succeeded", succeeded,
		"failed", failed,
	)
	return nil
}

func (r *Runner) processLibrary(ctx context.Context, runID pgtype.UUID, lib db.Library) error {
	mergedData := make(map[string]any)

	// GitHub collection
	ghResult, err := r.collectWithTracking(ctx, runID, lib.ID, r.github, lib)
	if err != nil {
		slog.Warn("github collection failed", "name", lib.Name, "err", err)
	} else {
		for k, v := range ghResult.Data {
			mergedData[k] = v
		}
	}

	// npm collection (only for npm packages)
	if lib.Registry == "npm" {
		npmResult, err := r.collectWithTracking(ctx, runID, lib.ID, r.npm, lib)
		if err != nil {
			slog.Warn("npm collection failed", "name", lib.Name, "err", err)
		} else {
			for k, v := range npmResult.Data {
				mergedData[k] = v
			}
		}
	}

	if len(mergedData) == 0 {
		return nil
	}

	// Calculate score
	result := scoring.Calculate(mergedData)

	// Store raw data as JSON
	rawData, _ := json.Marshal(mergedData)

	// Save score
	_, err = r.queries.CreateScore(ctx, db.CreateScoreParams{
		LibraryID:         lib.ID,
		OverallScore:      floatToNumeric(result.Overall),
		MaintenanceHealth: floatToNumeric(result.MaintenanceHealth),
		ApiClarity:        floatToNumeric(result.APIClarity),
		DocQuality:        floatToNumeric(result.DocQuality),
		SecurityPosture:   floatToNumeric(result.SecurityPosture),
		CommunitySignal:   floatToNumeric(result.CommunitySignal),
		DeprecationSafety: floatToNumeric(result.DeprecationSafety),
		ScoringVersion:    result.Version,
		RawData:           rawData,
	})
	if err != nil {
		return err
	}

	// Update library metadata
	latestVersion := toString(mergedData["latest_version"])
	var latestVersionDate pgtype.Timestamptz
	if dateStr := toString(mergedData["latest_version_date"]); dateStr != "" {
		if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
			latestVersionDate = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}

	err = r.queries.UpdateLibrary(ctx, db.UpdateLibraryParams{
		ID:                lib.ID,
		Name:              lib.Name,
		Description:       lib.Description,
		HomepageUrl:       lib.HomepageUrl,
		License:           firstNonEmpty(toString(mergedData["license"]), lib.License),
		LatestVersion:     latestVersion,
		LatestVersionDate: latestVersionDate,
		Deprecated:        result.Deprecated,
		Metadata:          rawData,
	})
	if err != nil {
		slog.Warn("failed to update library metadata", "name", lib.Name, "err", err)
	}

	slog.Info("scored library",
		"name", lib.Name,
		"overall", result.Overall,
		"deprecated", result.Deprecated,
	)
	return nil
}

func (r *Runner) collectWithTracking(ctx context.Context, runID, libID pgtype.UUID, c collector.Collector, lib db.Library) (*collector.Result, error) {
	item, err := r.queries.CreateCollectionRunItem(ctx, db.CreateCollectionRunItemParams{
		RunID:         runID,
		LibraryID:     libID,
		CollectorType: c.Name(),
	})
	if err != nil {
		return nil, err
	}

	result, collectErr := c.Collect(ctx, lib.GithubRepo, lib.Registry, lib.PackageName)

	status := "completed"
	var errMsg pgtype.Text
	var durationMs int32
	if collectErr != nil {
		status = "failed"
		errMsg = pgtype.Text{String: collectErr.Error(), Valid: true}
	}
	if result != nil {
		durationMs = int32(result.Duration.Milliseconds())
	}

	_ = r.queries.UpdateCollectionRunItem(ctx, db.UpdateCollectionRunItemParams{
		ID:           item.ID,
		Status:       status,
		ErrorMessage: errMsg,
		DurationMs:   durationMs,
	})

	if collectErr != nil {
		return nil, collectErr
	}
	return result, nil
}

func floatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%.2f", f))
	return n
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
