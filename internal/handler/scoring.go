package handler

import "net/http"

// ScoringWeights represents the current algorithm weights.
type ScoringWeights struct {
	MaintenanceHealth float64 `json:"maintenanceHealth"`
	APIClarity        float64 `json:"apiClarity"`
	DocQuality        float64 `json:"docQuality"`
	SecurityPosture   float64 `json:"securityPosture"`
	CommunitySignal   float64 `json:"communitySignal"`
	DeprecationSafety float64 `json:"deprecationSafety"`
	Version           string  `json:"version"`
}

func GetScoringWeights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=86400")
	writeJSON(w, http.StatusOK, Envelope{
		Data: ScoringWeights{
			MaintenanceHealth: 0.25,
			APIClarity:        0.20,
			DocQuality:        0.15,
			SecurityPosture:   0.15,
			CommunitySignal:   0.15,
			DeprecationSafety: 0.10,
			Version:           "1.0.0",
		},
	})
}
