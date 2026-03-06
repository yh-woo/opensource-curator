package model

import (
	"time"

	"github.com/google/uuid"
)

type Score struct {
	ID                uuid.UUID `json:"id"`
	LibraryID         uuid.UUID `json:"libraryId"`
	ScoredAt          time.Time `json:"scoredAt"`
	OverallScore      float64   `json:"overall"`
	MaintenanceHealth float64   `json:"maintenanceHealth"`
	APIClarity        float64   `json:"apiClarity"`
	DocQuality        float64   `json:"docQuality"`
	SecurityPosture   float64   `json:"securityPosture"`
	CommunitySignal   float64   `json:"communitySignal"`
	DeprecationSafety float64   `json:"deprecationSafety"`
	ScoringVersion    string    `json:"version"`
}

type ScoreBreakdown struct {
	MaintenanceHealth float64 `json:"maintenanceHealth"`
	APIClarity        float64 `json:"apiClarity"`
	DocQuality        float64 `json:"docQuality"`
	SecurityPosture   float64 `json:"securityPosture"`
	CommunitySignal   float64 `json:"communitySignal"`
	DeprecationSafety float64 `json:"deprecationSafety"`
}
