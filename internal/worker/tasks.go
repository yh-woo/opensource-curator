package worker

import (
	"context"
	"log/slog"

	"github.com/hibiken/asynq"

	"github.com/woo/opensource-curator/internal/db"
	"github.com/woo/opensource-curator/internal/pipeline"
)

const (
	TypeCollectAll = "collect:all"
)

// CollectAllHandler processes the collect-all task.
type CollectAllHandler struct {
	runner *pipeline.Runner
}

func NewCollectAllHandler(queries *db.Queries, githubToken string) *CollectAllHandler {
	return &CollectAllHandler{
		runner: pipeline.NewRunner(queries, githubToken),
	}
}

func (h *CollectAllHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	slog.Info("starting scheduled collection")
	if err := h.runner.RunAll(ctx); err != nil {
		slog.Error("scheduled collection failed", "err", err)
		return err
	}
	slog.Info("scheduled collection completed")
	return nil
}
