package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/woo/opensource-curator/internal/discovery"
)

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

	// Ensure new categories exist
	if err := ensureCategories(ctx, pool); err != nil {
		slog.Error("failed to ensure categories", "err", err)
		os.Exit(1)
	}

	disc := discovery.NewNpmDiscoverer()
	totalNew := 0

	for catSlug, keywords := range discovery.CategoryKeywords {
		// Get category ID
		var catID string
		err := pool.QueryRow(ctx,
			`SELECT id FROM categories WHERE slug = $1`, catSlug,
		).Scan(&catID)
		if err != nil {
			slog.Warn("category not found, skipping", "slug", catSlug, "err", err)
			continue
		}

		for _, kw := range keywords {
			results, err := disc.SearchByKeywords(kw, 20)
			if err != nil {
				slog.Warn("search failed", "keywords", kw, "err", err)
				continue
			}

			for _, pkg := range results {
				githubRepo := discovery.ExtractGithubRepo(pkg.Links.Repository)

				// Upsert library
				var libID string
				err := pool.QueryRow(ctx,
					`INSERT INTO libraries (name, registry, package_name, github_repo, description)
					 VALUES ($1, 'npm', $2, $3, $4)
					 ON CONFLICT (registry, package_name) DO UPDATE SET
					   github_repo = COALESCE(NULLIF($3, ''), libraries.github_repo),
					   description = COALESCE(NULLIF($4, ''), libraries.description)
					 RETURNING id`,
					pkg.Name, pkg.Name, githubRepo, pkg.Description,
				).Scan(&libID)
				if err != nil {
					slog.Warn("failed to upsert library", "name", pkg.Name, "err", err)
					continue
				}

				// Link to category
				_, err = pool.Exec(ctx,
					`INSERT INTO library_categories (library_id, category_id)
					 VALUES ($1::uuid, $2::uuid)
					 ON CONFLICT DO NOTHING`,
					libID, catID,
				)
				if err != nil {
					slog.Warn("failed to link category", "library", pkg.Name, "category", catSlug, "err", err)
				}
			}

			// Rate limit: npm search API (generous to avoid 429)
			time.Sleep(2 * time.Second)
		}

		slog.Info("discovered libraries for category", "category", catSlug)
	}

	// Count total libraries
	var total int
	pool.QueryRow(ctx, `SELECT COUNT(*) FROM libraries`).Scan(&total)

	// Count new (unscored) libraries
	pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM libraries l
		 WHERE NOT EXISTS (SELECT 1 FROM scores s WHERE s.library_id = l.id)`,
	).Scan(&totalNew)

	slog.Info("discovery completed",
		"total_libraries", total,
		"unscored", totalNew,
	)
	fmt.Printf("Discovery complete. Total: %d libraries (%d new, unscored)\n", total, totalNew)
	fmt.Println("Run 'make collect' to score new libraries.")
}

// New categories added by discovery
var newCategories = []struct {
	Slug        string
	Name        string
	Description string
}{
	{"email", "Email", "Email sending and templating libraries"},
	{"websocket", "WebSocket", "WebSocket and real-time communication"},
	{"graphql", "GraphQL", "GraphQL servers, clients, and tools"},
	{"markdown-templating", "Markdown / Templating", "Markdown parsers and template engines"},
	{"crypto-security", "Crypto / Security", "Encryption, hashing, and security utilities"},
	{"config-env", "Config / Environment", "Configuration management and environment variables"},
	{"monitoring", "Monitoring", "APM, error tracking, and metrics"},
	{"bundler-build", "Bundler / Build", "Build tools and bundlers"},
}

func ensureCategories(ctx context.Context, pool *pgxpool.Pool) error {
	for _, c := range newCategories {
		_, err := pool.Exec(ctx,
			`INSERT INTO categories (slug, name, description)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (slug) DO NOTHING`,
			c.Slug, c.Name, c.Description,
		)
		if err != nil {
			return fmt.Errorf("insert category %s: %w", c.Slug, err)
		}
	}
	return nil
}
