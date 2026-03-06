package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/hibiken/asynq"

	"github.com/woo/opensource-curator/internal/db"
	"github.com/woo/opensource-curator/internal/worker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://curator:curator@localhost:5432/curator?sslmode=disable"
	}

	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	githubToken := os.Getenv("GITHUB_TOKEN")

	ctx := context.Background()
	pool, err := db.NewPool(ctx, dbURL)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("connected to database")

	queries := db.New(pool)

	// Create asynq server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 1,
			Queues:      map[string]int{"default": 1},
		},
	)

	// Register task handlers
	mux := asynq.NewServeMux()
	mux.Handle(worker.TypeCollectAll, worker.NewCollectAllHandler(queries, githubToken))

	// Schedule periodic collection (every 24 hours)
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisAddr},
		nil,
	)
	_, err = scheduler.Register("0 3 * * *", asynq.NewTask(worker.TypeCollectAll, nil))
	if err != nil {
		slog.Error("failed to register scheduled task", "err", err)
		os.Exit(1)
	}

	// Start scheduler in background
	go func() {
		if err := scheduler.Run(); err != nil {
			slog.Error("scheduler error", "err", err)
		}
	}()

	slog.Info("worker starting", "redis", redisAddr, "schedule", "daily at 03:00")

	if err := srv.Run(mux); err != nil {
		slog.Error("worker error", "err", err)
		os.Exit(1)
	}
}
