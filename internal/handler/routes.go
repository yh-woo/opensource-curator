package handler

import (
	"encoding/json"
	"math"
	"math/big"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/woo/opensource-curator/internal/db"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	Q *db.Queries
}

func RegisterRoutes(r chi.Router, q *db.Queries, opts ...RouteOption) {
	h := &Handler{Q: q}
	var cfg routeConfig
	for _, o := range opts {
		o(&cfg)
	}
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", Health)
		r.Get("/libraries", h.ListLibraries)
		r.Get("/libraries/{id}", h.GetLibrary)
		r.Get("/libraries/{registry}/{packageName}", h.GetLibraryBySlug)
		r.Get("/categories", h.ListCategories)
		r.Get("/categories/{slug}", h.GetCategory)
		r.Get("/search", h.Search)
		r.Get("/recommend", h.Recommend)
		r.Get("/scoring/weights", GetScoringWeights)
		if cfg.collectHandler != nil {
			r.Post("/collect", cfg.collectHandler)
		}
	})
}

type routeConfig struct {
	collectHandler http.HandlerFunc
}

// RouteOption configures optional routes.
type RouteOption func(*routeConfig)

// WithCollectHandler adds the POST /v1/collect endpoint.
func WithCollectHandler(h http.HandlerFunc) RouteOption {
	return func(c *routeConfig) { c.collectHandler = h }
}

// Envelope is the standard API response wrapper.
type Envelope struct {
	Data        any      `json:"data"`
	Error       string   `json:"error,omitempty"`
	NextActions []Action `json:"next_actions,omitempty"`
}

// Action represents a HATEOAS-style next action for AI agents.
type Action struct {
	Rel         string `json:"rel"`
	Href        string `json:"href"`
	Description string `json:"description,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, Envelope{Error: msg})
}

// numericToFloat converts pgtype.Numeric to float64.
func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	if !f.Valid {
		return 0
	}
	return math.Round(f.Float64*100) / 100
}

// floatToNumeric converts float64 to pgtype.Numeric.
func floatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Valid = true
	n.Int = big.NewInt(int64(math.Round(f * 100)))
	n.Exp = -2
	return n
}

// LibraryResponse is the JSON shape returned to clients.
type LibraryResponse struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Registry          string         `json:"registry"`
	PackageName       string         `json:"packageName"`
	GithubRepo        string         `json:"githubRepo"`
	Description       string         `json:"description"`
	LatestVersion     string         `json:"latestVersion"`
	LatestVersionDate string         `json:"latestVersionDate,omitempty"`
	Deprecated        bool           `json:"deprecated"`
	Score             *ScoreResponse `json:"score,omitempty"`
}

type ScoreResponse struct {
	Overall   float64        `json:"overall"`
	Breakdown ScoreBreakdown `json:"breakdown"`
	Version   string         `json:"version"`
}

type ScoreBreakdown struct {
	MaintenanceHealth float64 `json:"maintenanceHealth"`
	APIClarity        float64 `json:"apiClarity"`
	DocQuality        float64 `json:"docQuality"`
	SecurityPosture   float64 `json:"securityPosture"`
	CommunitySignal   float64 `json:"communitySignal"`
	DeprecationSafety float64 `json:"deprecationSafety"`
}
