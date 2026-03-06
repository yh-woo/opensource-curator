-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- libraries
CREATE TABLE libraries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    registry TEXT NOT NULL CHECK (registry IN ('npm', 'pypi', 'crates', 'go')),
    package_name TEXT NOT NULL,
    github_repo TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    homepage_url TEXT,
    license TEXT NOT NULL DEFAULT '',
    latest_version TEXT NOT NULL DEFAULT '',
    latest_version_date TIMESTAMPTZ,
    deprecated BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    metadata JSONB NOT NULL DEFAULT '{}',
    UNIQUE(registry, package_name)
);

-- categories
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    parent_id UUID REFERENCES categories(id)
);

-- library_categories (M:N)
CREATE TABLE library_categories (
    library_id UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (library_id, category_id)
);

-- scores
CREATE TABLE scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
    scored_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    overall_score NUMERIC(5,2) NOT NULL,
    maintenance_health NUMERIC(5,2) NOT NULL,
    api_clarity NUMERIC(5,2) NOT NULL,
    doc_quality NUMERIC(5,2) NOT NULL,
    security_posture NUMERIC(5,2) NOT NULL,
    community_signal NUMERIC(5,2) NOT NULL,
    deprecation_safety NUMERIC(5,2) NOT NULL,
    scoring_version TEXT NOT NULL,
    raw_data JSONB NOT NULL DEFAULT '{}'
);

-- alternatives
CREATE TABLE alternatives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_library_id UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
    target_library_id UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
    relationship TEXT NOT NULL CHECK (relationship IN ('replacement', 'complement', 'fork')),
    reason TEXT NOT NULL DEFAULT ''
);

-- collection_runs
CREATE TABLE collection_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'running'
        CHECK (status IN ('running', 'completed', 'completed_with_errors', 'failed')),
    libraries_processed INT NOT NULL DEFAULT 0,
    libraries_succeeded INT NOT NULL DEFAULT 0,
    libraries_failed INT NOT NULL DEFAULT 0,
    success_rate NUMERIC(5,2) NOT NULL DEFAULT 0,
    errors JSONB NOT NULL DEFAULT '[]',
    trigger TEXT NOT NULL CHECK (trigger IN ('scheduled', 'manual'))
);

-- collection_run_items
CREATE TABLE collection_run_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL REFERENCES collection_runs(id) ON DELETE CASCADE,
    library_id UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
    collector_type TEXT NOT NULL CHECK (collector_type IN ('github', 'npm', 'pypi', 'security')),
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'success', 'failed', 'skipped')),
    error_message TEXT,
    duration_ms INT NOT NULL DEFAULT 0,
    completed_at TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_scores_library_latest ON scores(library_id, scored_at DESC);
CREATE INDEX idx_scores_overall_rank ON scores(overall_score DESC);
CREATE INDEX idx_lib_categories_cat ON library_categories(category_id);
CREATE INDEX idx_collection_items_run ON collection_run_items(run_id, status);
CREATE INDEX idx_libraries_registry ON libraries(registry);
CREATE INDEX idx_libraries_deprecated ON libraries(deprecated) WHERE deprecated = true;
