package handler

import (
	"net/http"

	"github.com/woo/opensource-curator/internal/db"
	"github.com/woo/opensource-curator/internal/recommend"
)

type RecommendResponse struct {
	Data        []RecommendItem `json:"data"`
	Query       RecommendQuery  `json:"query"`
	NextActions []Action        `json:"next_actions,omitempty"`
}

type RecommendItem struct {
	Rank        int             `json:"rank"`
	Library     LibraryResponse `json:"library"`
	MatchReason string          `json:"matchReason"`
}

type RecommendQuery struct {
	Task              string   `json:"task"`
	Prefer            string   `json:"prefer,omitempty"`
	MatchedCategories []string `json:"matchedCategories"`
}

func (h *Handler) Recommend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")
	task := r.URL.Query().Get("task")
	prefer := r.URL.Query().Get("prefer")

	if task == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'task' is required")
		return
	}

	result := recommend.Match(task, prefer)

	var items []RecommendItem
	seen := make(map[string]bool)

	for _, catSlug := range result.MatchedCategories {
		rows, err := h.Q.ListLibrariesByCategory(r.Context(), db.ListLibrariesByCategoryParams{
			Slug:   catSlug,
			Limit:  10,
			Offset: 0,
		})
		if err != nil {
			continue
		}

		for _, row := range rows {
			id := uuidToString(row.ID)
			if seen[id] {
				continue
			}
			seen[id] = true

			lib := LibraryResponse{
				ID:            id,
				Name:          row.Name,
				Registry:      row.Registry,
				PackageName:   row.PackageName,
				GithubRepo:    row.GithubRepo,
				Description:   row.Description,
				LatestVersion: row.LatestVersion,
				Deprecated:    row.Deprecated,
			}
			if row.LatestVersionDate.Valid {
				lib.LatestVersionDate = row.LatestVersionDate.Time.Format("2006-01-02")
			}
			if row.OverallScore.Valid {
				lib.Score = &ScoreResponse{
					Overall: numericToFloat(row.OverallScore),
					Version: row.ScoringVersion,
				}
			}

			items = append(items, RecommendItem{
				Library:     lib,
				MatchReason: result.MatchReason,
			})
		}
	}

	// Assign ranks
	for i := range items {
		items[i].Rank = i + 1
	}
	if items == nil {
		items = []RecommendItem{}
	}

	resp := RecommendResponse{
		Data: items,
		Query: RecommendQuery{
			Task:              task,
			Prefer:            prefer,
			MatchedCategories: result.MatchedCategories,
		},
		NextActions: []Action{
			{Rel: "categories", Href: "/v1/categories", Description: "Browse all categories"},
			{Rel: "search", Href: "/v1/search", Description: "Search libraries by name"},
		},
	}

	writeJSON(w, http.StatusOK, resp)
}
