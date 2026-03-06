package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/woo/opensource-curator/internal/db"
	"github.com/woo/opensource-curator/internal/pipeline"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://curator:curator@localhost:5432/curator?sslmode=disable"
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		slog.Warn("GITHUB_TOKEN not set - GitHub API rate limit will be 60 requests/hour")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		slog.Info("shutting down collection...")
		cancel()
	}()

	pool, err := db.NewPool(ctx, dbURL)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	queries := db.New(pool)
	runner := pipeline.NewRunner(queries, githubToken)

	slog.Info("starting collection pipeline")
	if err := runner.RunAll(ctx); err != nil {
		slog.Error("collection failed", "err", err)
		os.Exit(1)
	}
	slog.Info("collection pipeline completed")
}
