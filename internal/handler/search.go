package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/woo/opensource-curator/internal/db"
)

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=3600")
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	rows, err := h.Q.SearchLibraries(r.Context(), db.SearchLibrariesParams{
		Column1: pgtype.Text{String: q, Valid: true},
		Limit:   20,
		Offset:  0,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "search failed")
		return
	}

	libs := make([]LibraryResponse, 0, len(rows))
	for _, row := range rows {
		lib := toLibraryResponse(row)
		h.attachScore(r.Context(), row.ID, &lib)
		libs = append(libs, lib)
	}

	writeJSON(w, http.StatusOK, Envelope{
		Data: libs,
		NextActions: []Action{
			{Rel: "recommend", Href: "/v1/recommend?task=" + q, Description: "Get AI-optimized recommendations for this query"},
		},
	})
}
