-- name: GetLibrary :one
SELECT * FROM libraries WHERE id = $1;

-- name: GetLibraryBySlug :one
SELECT * FROM libraries WHERE registry = $1 AND package_name = $2;

-- name: ListLibraries :many
SELECT l.*, s.overall_score, COALESCE(s.scoring_version, '') AS scoring_version
FROM libraries l
LEFT JOIN LATERAL (
    SELECT overall_score, scoring_version
    FROM scores
    WHERE library_id = l.id
    ORDER BY scored_at DESC
    LIMIT 1
) s ON true
WHERE ($1::text = '' OR l.registry = $1)
  AND ($2::boolean IS NULL OR l.deprecated = $2)
ORDER BY COALESCE(s.overall_score, 0) DESC
LIMIT $3 OFFSET $4;

-- name: ListLibrariesByCategory :many
SELECT l.*, s.overall_score, COALESCE(s.scoring_version, '') AS scoring_version
FROM libraries l
JOIN library_categories lc ON lc.library_id = l.id
JOIN categories c ON c.id = lc.category_id
LEFT JOIN LATERAL (
    SELECT overall_score, scoring_version
    FROM scores
    WHERE library_id = l.id
    ORDER BY scored_at DESC
    LIMIT 1
) s ON true
WHERE c.slug = $1
ORDER BY COALESCE(s.overall_score, 0) DESC
LIMIT $2 OFFSET $3;

-- name: SearchLibraries :many
SELECT * FROM libraries
WHERE name ILIKE '%' || $1 || '%'
   OR description ILIKE '%' || $1 || '%'
   OR package_name ILIKE '%' || $1 || '%'
LIMIT $2 OFFSET $3;

-- name: CreateLibrary :one
INSERT INTO libraries (name, registry, package_name, github_repo, description, homepage_url, license)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateLibrary :exec
UPDATE libraries SET
    name = $2,
    description = $3,
    homepage_url = $4,
    license = $5,
    latest_version = $6,
    latest_version_date = $7,
    deprecated = $8,
    metadata = $9,
    updated_at = now()
WHERE id = $1;

-- name: GetCategory :one
SELECT * FROM categories WHERE slug = $1;

-- name: ListCategories :many
SELECT c.*, COUNT(lc.library_id) AS library_count
FROM categories c
LEFT JOIN library_categories lc ON lc.category_id = c.id
GROUP BY c.id
ORDER BY c.name;

-- name: CreateCategory :one
INSERT INTO categories (slug, name, description, parent_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: LinkLibraryCategory :exec
INSERT INTO library_categories (library_id, category_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: GetLatestScore :one
SELECT * FROM scores
WHERE library_id = $1
ORDER BY scored_at DESC
LIMIT 1;

-- name: GetScoreHistory :many
SELECT * FROM scores
WHERE library_id = $1
ORDER BY scored_at DESC
LIMIT $2;

-- name: CreateScore :one
INSERT INTO scores (
    library_id, overall_score,
    maintenance_health, api_clarity, doc_quality,
    security_posture, community_signal, deprecation_safety,
    scoring_version, raw_data
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: CreateCollectionRun :one
INSERT INTO collection_runs (trigger) VALUES ($1) RETURNING *;

-- name: UpdateCollectionRun :exec
UPDATE collection_runs SET
    completed_at = now(),
    status = $2,
    libraries_processed = $3,
    libraries_succeeded = $4,
    libraries_failed = $5,
    success_rate = $6,
    errors = $7
WHERE id = $1;

-- name: CreateCollectionRunItem :one
INSERT INTO collection_run_items (run_id, library_id, collector_type)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCollectionRunItem :exec
UPDATE collection_run_items SET
    status = $2,
    error_message = $3,
    duration_ms = $4,
    completed_at = now()
WHERE id = $1;

-- name: ListAllLibraries :many
SELECT * FROM libraries ORDER BY name;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;
