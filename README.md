# Opensource Curator

AI-agent-optimized open-source library curation service. Scores libraries on 6 metrics weighted for AI agent usability, using real-time data from GitHub and npm.

## Why

AI agents pick libraries based on training-data bias, not actual quality. This project fixes that by computing scores from live metadata — so agents (and humans) can make informed choices.

## Scoring Metrics

| Metric | Weight | What it measures |
|--------|--------|------------------|
| Maintenance Health | 25% | Commit frequency, issue response time, release cadence |
| API Clarity | 20% | TypeScript types, export structure, naming conventions |
| Doc Quality | 15% | README completeness, examples, API reference |
| Security Posture | 15% | Vulnerability history, dependency count, audit status |
| Community Signal | 15% | Stars, contributors, download trends |
| Deprecation Safety | 10% | Deprecation markers, successor availability |

## Tech Stack

- **Backend**: Go, Chi router, sqlc, pgx/v5, PostgreSQL 16, Redis 7
- **Frontend**: Next.js 15, Tailwind CSS v4, TypeScript
- **Worker**: asynq (Redis-based task queue, daily scheduled collection)
- **Infra**: Docker Compose

## Quick Start

```bash
# 1. Start PostgreSQL + Redis
make dev-infra

# 2. Run migrations
make migrate-up

# 3. Seed 77 libraries across 15 categories
make seed

# 4. Collect scores (optional: set GITHUB_TOKEN for better data)
export GITHUB_TOKEN=ghp_...
make collect

# 5. Start API server (port 8080)
make dev-api

# 6. Start frontend (port 3000)
make dev-web
```

## API

```
GET /v1/health
GET /v1/libraries
GET /v1/libraries/:id
GET /v1/categories
GET /v1/categories/:slug
GET /v1/search?q=http+client
GET /v1/recommend?task=make+http+requests&prefer=lightweight
POST /v1/collect              # trigger collection run
```

All responses include `next_actions` (HATEOAS) for AI agent workflow guidance.

## Project Structure

```
cmd/
  api/          # HTTP server
  collect/      # CLI collection runner
  seed/         # DB seeder
  worker/       # asynq worker (scheduled collection)
internal/
  collector/    # GitHub + npm data collectors
  db/           # sqlc generated queries
  handler/      # HTTP handlers
  model/        # Domain types
  pipeline/     # Collection orchestrator
  recommend/    # Recommendation engine
  scoring/      # 6-metric scoring engine
  worker/       # Task definitions
migrations/     # PostgreSQL schema
seed/           # Category + library seed data
web/            # Next.js frontend
```

## License

MIT
