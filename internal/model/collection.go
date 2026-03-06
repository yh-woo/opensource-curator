package model

import (
	"time"

	"github.com/google/uuid"
)

type CollectionStatus string

const (
	CollectionRunning            CollectionStatus = "running"
	CollectionCompleted          CollectionStatus = "completed"
	CollectionCompletedWithErrors CollectionStatus = "completed_with_errors"
	CollectionFailed             CollectionStatus = "failed"
)

type CollectionTrigger string

const (
	TriggerScheduled CollectionTrigger = "scheduled"
	TriggerManual    CollectionTrigger = "manual"
)

type CollectionRun struct {
	ID                 uuid.UUID         `json:"id"`
	StartedAt          time.Time         `json:"startedAt"`
	CompletedAt        *time.Time        `json:"completedAt,omitempty"`
	Status             CollectionStatus  `json:"status"`
	LibrariesProcessed int               `json:"librariesProcessed"`
	LibrariesSucceeded int               `json:"librariesSucceeded"`
	LibrariesFailed    int               `json:"librariesFailed"`
	SuccessRate        float64           `json:"successRate"`
	Trigger            CollectionTrigger `json:"trigger"`
}

type CollectorType string

const (
	CollectorGitHub   CollectorType = "github"
	CollectorNPM      CollectorType = "npm"
	CollectorPyPI     CollectorType = "pypi"
	CollectorSecurity CollectorType = "security"
)

type ItemStatus string

const (
	ItemPending ItemStatus = "pending"
	ItemSuccess ItemStatus = "success"
	ItemFailed  ItemStatus = "failed"
	ItemSkipped ItemStatus = "skipped"
)

type CollectionRunItem struct {
	ID            uuid.UUID     `json:"id"`
	RunID         uuid.UUID     `json:"runId"`
	LibraryID     uuid.UUID     `json:"libraryId"`
	CollectorType CollectorType `json:"collectorType"`
	Status        ItemStatus    `json:"status"`
	ErrorMessage  *string       `json:"errorMessage,omitempty"`
	DurationMs    int           `json:"durationMs"`
	CompletedAt   *time.Time    `json:"completedAt,omitempty"`
}
