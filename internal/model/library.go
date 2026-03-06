package model

import (
	"time"

	"github.com/google/uuid"
)

type Registry string

const (
	RegistryNPM    Registry = "npm"
	RegistryPyPI   Registry = "pypi"
	RegistryCrates Registry = "crates"
	RegistryGo     Registry = "go"
)

type Library struct {
	ID                uuid.UUID  `json:"id"`
	Name              string     `json:"name"`
	Registry          Registry   `json:"registry"`
	PackageName       string     `json:"packageName"`
	GithubRepo        string     `json:"githubRepo"`
	Description       string     `json:"description"`
	HomepageURL       *string    `json:"homepageUrl,omitempty"`
	License           string     `json:"license"`
	LatestVersion     string     `json:"latestVersion"`
	LatestVersionDate *time.Time `json:"latestVersionDate,omitempty"`
	Deprecated        bool       `json:"deprecated"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}
