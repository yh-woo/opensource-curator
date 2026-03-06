package collector

import (
	"context"
	"time"
)

// Result holds the raw data collected from an external source.
type Result struct {
	Data      map[string]any
	Duration  time.Duration
	Err       error
}

// Collector defines the interface for data collection from external APIs.
type Collector interface {
	Name() string
	Collect(ctx context.Context, githubRepo string, registry string, packageName string) (*Result, error)
}
