package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/woo/opensource-curator/internal/db"
)

func (h *Handler) ListLibraries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	rows, err := h.Q.ListLibraries(r.Context(), db.ListLibrariesParams{
		Column1: r.URL.Query().Get("registry"),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list libraries")
		return
	}

	libs := make([]LibraryResponse, 0, len(rows))
	for _, row := range rows {
		lib := LibraryResponse{
			ID:          uuidToString(row.ID),
			Name:        row.Name,
			Registry:    row.Registry,
			PackageName: row.PackageName,
			GithubRepo:  row.GithubRepo,
			Description: row.Description,
			LatestVersion: row.LatestVersion,
			Deprecated:  row.Deprecated,
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
		Data: libs,
		NextActions: []Action{
			{Rel: "self", Href: "/v1/libraries", Description: "List libraries"},
			{Rel: "categories", Href: "/v1/categories", Description: "Browse by category"},
		},
	})
}

func (h *Handler) GetLibrary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")
	id := chi.URLParam(r, "id")

	uuid, err := stringToUUID(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid library ID")
		return
	}

	lib, err := h.Q.GetLibrary(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "library not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get library")
		return
	}

	resp := toLibraryResponse(lib)
	h.attachScore(r.Context(), lib.ID, &resp)

	writeJSON(w, http.StatusOK, Envelope{
		Data: resp,
		NextActions: []Action{
			{Rel: "scores", Href: "/v1/libraries/" + id, Description: "Library detail"},
			{Rel: "categories", Href: "/v1/categories", Description: "Browse categories"},
		},
	})
}

func (h *Handler) GetLibraryBySlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")
	registry := chi.URLParam(r, "registry")
	packageName := chi.URLParam(r, "packageName")

	lib, err := h.Q.GetLibraryBySlug(r.Context(), db.GetLibraryBySlugParams{
		Registry:    registry,
		PackageName: packageName,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "library not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get library")
		return
	}

	resp := toLibraryResponse(lib)
	h.attachScore(r.Context(), lib.ID, &resp)

	writeJSON(w, http.StatusOK, Envelope{
		Data: resp,
		NextActions: []Action{
			{Rel: "self", Href: "/v1/libraries/" + registry + "/" + packageName},
			{Rel: "recommend", Href: "/v1/recommend", Description: "Get recommendations"},
		},
	})
}

func toLibraryResponse(lib db.Library) LibraryResponse {
	resp := LibraryResponse{
		ID:            uuidToString(lib.ID),
		Name:          lib.Name,
		Registry:      lib.Registry,
		PackageName:   lib.PackageName,
		GithubRepo:    lib.GithubRepo,
		Description:   lib.Description,
		LatestVersion: lib.LatestVersion,
		Deprecated:    lib.Deprecated,
	}
	if lib.LatestVersionDate.Valid {
		resp.LatestVersionDate = lib.LatestVersionDate.Time.Format("2006-01-02")
	}
	return resp
}

func (h *Handler) attachScore(ctx context.Context, libID pgtype.UUID, resp *LibraryResponse) {
	score, err := h.Q.GetLatestScore(ctx, libID)
	if err != nil {
		return
	}
	resp.Score = &ScoreResponse{
		Overall: numericToFloat(score.OverallScore),
		Breakdown: ScoreBreakdown{
			MaintenanceHealth: numericToFloat(score.MaintenanceHealth),
			APIClarity:        numericToFloat(score.ApiClarity),
			DocQuality:        numericToFloat(score.DocQuality),
			SecurityPosture:   numericToFloat(score.SecurityPosture),
			CommunitySignal:   numericToFloat(score.CommunitySignal),
			DeprecationSafety: numericToFloat(score.DeprecationSafety),
		},
		Version: score.ScoringVersion,
	}
}
