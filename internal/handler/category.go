package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/woo/opensource-curator/internal/db"
)

type CategoryResponse struct {
	ID           string            `json:"id"`
	Slug         string            `json:"slug"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	LibraryCount int64             `json:"libraryCount,omitempty"`
	Libraries    []LibraryResponse `json:"libraries,omitempty"`
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")

	rows, err := h.Q.ListCategories(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list categories")
		return
	}

	cats := make([]CategoryResponse, 0, len(rows))
	for _, row := range rows {
		cats = append(cats, CategoryResponse{
			ID:           uuidToString(row.ID),
			Slug:         row.Slug,
			Name:         row.Name,
			Description:  row.Description,
			LibraryCount: row.LibraryCount,
		})
	}

	writeJSON(w, http.StatusOK, Envelope{
		Data: cats,
		NextActions: []Action{
			{Rel: "recommend", Href: "/v1/recommend", Description: "Get AI-optimized recommendations"},
		},
	})
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")
	slug := chi.URLParam(r, "slug")

	cat, err := h.Q.GetCategory(r.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "category not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get category")
		return
	}

	libRows, err := h.Q.ListLibrariesByCategory(r.Context(), db.ListLibrariesByCategoryParams{
		Slug:   slug,
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list libraries")
		return
	}

	libs := make([]LibraryResponse, 0, len(libRows))
	for _, row := range libRows {
		lib := LibraryResponse{
			ID:            uuidToString(row.ID),
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
		libs = append(libs, lib)
	}

	writeJSON(w, http.StatusOK, Envelope{
		Data: CategoryResponse{
			ID:          uuidToString(cat.ID),
			Slug:        cat.Slug,
			Name:        cat.Name,
			Description: cat.Description,
			Libraries:   libs,
		},
		NextActions: []Action{
			{Rel: "categories", Href: "/v1/categories", Description: "All categories"},
		},
	})
}
