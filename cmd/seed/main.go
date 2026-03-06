package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SeedCategory struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SeedLibrary struct {
	Name        string   `json:"name"`
	Registry    string   `json:"registry"`
	PackageName string   `json:"package_name"`
	GithubRepo  string   `json:"github_repo"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://curator:curator@localhost:5432/curator?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "err", err)
		os.Exit(1)
	}

	if err := seedCategories(ctx, pool); err != nil {
		slog.Error("failed to seed categories", "err", err)
		os.Exit(1)
	}

	if err := seedLibraries(ctx, pool); err != nil {
		slog.Error("failed to seed libraries", "err", err)
		os.Exit(1)
	}

	fmt.Println("seed data loaded successfully")
}

func seedCategories(ctx context.Context, pool *pgxpool.Pool) error {
	data, err := os.ReadFile("seed/categories.json")
	if err != nil {
		return fmt.Errorf("read categories.json: %w", err)
	}

	var categories []SeedCategory
	if err := json.Unmarshal(data, &categories); err != nil {
		return fmt.Errorf("parse categories.json: %w", err)
	}

	for _, c := range categories {
		_, err := pool.Exec(ctx,
			`INSERT INTO categories (slug, name, description)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (slug) DO UPDATE SET name = $2, description = $3`,
			c.Slug, c.Name, c.Description,
		)
		if err != nil {
			return fmt.Errorf("insert category %s: %w", c.Slug, err)
		}
		slog.Info("seeded category", "slug", c.Slug)
	}

	slog.Info("categories seeded", "count", len(categories))
	return nil
}

func seedLibraries(ctx context.Context, pool *pgxpool.Pool) error {
	data, err := os.ReadFile("seed/libraries.json")
	if err != nil {
		return fmt.Errorf("read libraries.json: %w", err)
	}

	var libraries []SeedLibrary
	if err := json.Unmarshal(data, &libraries); err != nil {
		return fmt.Errorf("parse libraries.json: %w", err)
	}

	for _, lib := range libraries {
		var libID string
		err := pool.QueryRow(ctx,
			`INSERT INTO libraries (name, registry, package_name, github_repo, description)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (registry, package_name) DO UPDATE SET
			   name = $1, github_repo = $4, description = $5
			 RETURNING id`,
			lib.Name, lib.Registry, lib.PackageName, lib.GithubRepo, lib.Description,
		).Scan(&libID)
		if err != nil {
			return fmt.Errorf("insert library %s: %w", lib.Name, err)
		}

		for _, catSlug := range lib.Categories {
			_, err := pool.Exec(ctx,
				`INSERT INTO library_categories (library_id, category_id)
				 SELECT $1::uuid, id FROM categories WHERE slug = $2
				 ON CONFLICT DO NOTHING`,
				libID, catSlug,
			)
			if err != nil {
				return fmt.Errorf("link library %s to category %s: %w", lib.Name, catSlug, err)
			}
		}

		slog.Info("seeded library", "name", lib.Name, "registry", lib.Registry)
	}

	slog.Info("libraries seeded", "count", len(libraries))
	return nil
}
