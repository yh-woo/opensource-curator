# opensource-curator 종합 계획서

**작성일:** 2026-03-06
**상태:** 합의 완료 (Planner-Architect-Critic 2회 반복)
**복잡도:** HIGH

---

## 1. 프로젝트 컨텍스트

### 문제 정의
AI 에이전트(Claude, GPT 등)가 오픈소스 라이브러리를 추천할 때, 학습 데이터 편향("자주 언급된 것 = 좋은 것")에 의존하여:
- 인기 있지만 deprecated된 라이브러리 추천
- 더 나은 최신 대안 누락
- 니치하지만 우수한 솔루션 무시
- 실시간 품질 시그널(유지보수 활동, 보안 이슈, 커뮤니티 건강도) 부재

### 비전
AI 에이전트 소비에 최적화된 오픈소스 라이브러리 큐레이션 서비스. 인간 개발자 브라우징이 아닌, AI 에이전트가 프로그래매틱하게 조회하여 더 나은 라이브러리 추천을 할 수 있도록 한다.

---

## 2. 기술 스택

### Backend (Go)
| 기술 | 선정 이유 |
|------|-----------|
| **Go** | goroutine 기반 병렬 수집에 자연스러움, 단일 바이너리 배포, 낮은 메모리 사용, 개발자 숙련도 최고 |
| **Chi** | 경량 HTTP 라우터, 미들웨어 체이닝, 표준 `net/http` 호환 |
| **sqlc** | SQL → 타입 안전 Go 코드 생성, SQL 직접 작성으로 복잡한 랭킹 쿼리에 최적 |
| **golang-migrate** | SQL 마이그레이션 관리, CLI + 라이브러리 양용 |
| **asynq** | Redis 기반 작업 큐, 재시도/스케줄링/동시성 제어, Go 네이티브 |
| **PostgreSQL** | 복잡한 쿼리(카테고리별 랭킹, 시계열 점수 추적), JSONB, 전문 검색 |
| **Redis** | asynq 작업 큐 백엔드 + API 응답 캐시 |

### Frontend (별도 프로젝트)
| 기술 | 선정 이유 |
|------|-----------|
| **Next.js 15 (App Router)** | SSR/SSG로 SEO 최적화, React Server Components |
| **Tailwind CSS v4** | 유틸리티 퍼스트, 빠른 UI 개발 |
| **shadcn/ui** | Radix UI 기반 컴포넌트, 커스터마이징 자유도 |
| **TanStack Table** | 라이브러리 비교 테이블, 정렬/필터링/페이지네이션 |

### 인프라 & 도구
| 기술 | 선정 이유 |
|------|-----------|
| **Docker Compose** | 로컬 개발 환경(PostgreSQL, Redis) |
| **Makefile** | 빌드/테스트/마이그레이션/시드 명령 통합 |
| **go test + testify** | Go 표준 테스트 + assertion 라이브러리 |
| **GitHub Actions** | CI/CD, 스케줄드 워크플로우 |
| **openapi-typescript-codegen** | Go API의 OpenAPI 스펙에서 프론트엔드 타입 자동 생성 |

---

## 3. 프로젝트 구조

```
opensource-curator/
├── cmd/
│   ├── api/                    # API 서버 진입점
│   │   └── main.go
│   ├── worker/                 # 큐 워커 진입점
│   │   └── main.go
│   └── seed/                   # 시드 데이터 로더
│       └── main.go
├── internal/
│   ├── handler/                # HTTP 핸들러 (Chi 라우트)
│   │   ├── library.go
│   │   ├── category.go
│   │   ├── recommend.go
│   │   ├── health.go
│   │   └── middleware.go       # rate limit, cache headers, CORS
│   ├── collector/              # 데이터 수집기
│   │   ├── github.go
│   │   ├── npm.go
│   │   ├── pypi.go             # Phase 2
│   │   ├── security.go         # Phase 2
│   │   └── collector.go        # 공통 인터페이스
│   ├── scoring/                # 점수 산정 엔진 (정식 위치)
│   │   ├── engine.go           # 통합 스코어러
│   │   ├── prefilter.go        # deprecation 프리필터
│   │   ├── maintenance.go
│   │   ├── api_clarity.go
│   │   ├── doc_quality.go
│   │   ├── security.go
│   │   ├── community.go
│   │   ├── deprecation.go
│   │   └── weights.go          # 가중치 설정
│   ├── recommend/              # 추천 엔진
│   │   ├── engine.go           # 키워드 → 카테고리 매핑
│   │   └── keywords.go         # 키워드 딕셔너리
│   ├── model/                  # 도메인 타입
│   │   ├── library.go
│   │   ├── score.go
│   │   ├── category.go
│   │   └── collection.go
│   ├── db/                     # 데이터베이스
│   │   ├── query/              # sqlc 생성 코드
│   │   ├── query.sql           # sqlc 쿼리 정의
│   │   └── db.go               # 연결 관리
│   └── queue/                  # asynq 작업 정의
│       ├── tasks.go
│       └── handlers.go
├── migrations/                 # SQL 마이그레이션 파일
│   ├── 001_initial.up.sql
│   └── 001_initial.down.sql
├── seed/                       # 시드 데이터 (JSON)
│   ├── categories.json
│   └── libraries.json
├── web/                        # Next.js 프론트엔드
│   ├── app/
│   ├── components/
│   ├── lib/
│   ├── package.json
│   └── next.config.js
├── docs/                       # 프로젝트 문서 (이 파일들)
│   ├── PLAN.md
│   └── DECISIONS.md
├── docker-compose.yml
├── Makefile
├── go.mod
├── sqlc.yaml
└── .github/
    └── workflows/
        ├── ci.yml
        └── data-pipeline.yml
```

### 빌드/실행 명령 (Makefile)
```makefile
make dev-api        # API 서버 실행 (air 핫리로드)
make dev-worker     # 워커 실행
make dev-web        # Next.js 프론트엔드 실행
make dev-infra      # Docker Compose (PostgreSQL + Redis)
make migrate-up     # 마이그레이션 실행
make migrate-down   # 마이그레이션 롤백
make seed           # 시드 데이터 적재
make sqlc           # sqlc 코드 생성
make test           # 전체 테스트
make build          # 바이너리 빌드 (api + worker)
```

---

## 4. 데이터 모델

### 핵심 엔티티

```sql
-- libraries
CREATE TABLE libraries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,                          -- 표시명
    registry TEXT NOT NULL CHECK (registry IN ('npm', 'pypi', 'crates', 'go')),
    package_name TEXT NOT NULL,                  -- 레지스트리 내 패키지명
    github_repo TEXT NOT NULL,                   -- owner/repo
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
```

### 인덱스
```sql
CREATE INDEX idx_scores_library_latest ON scores(library_id, scored_at DESC);
CREATE INDEX idx_scores_overall_rank ON scores(overall_score DESC);
CREATE INDEX idx_lib_categories_cat ON library_categories(category_id);
CREATE INDEX idx_collection_items_run ON collection_run_items(run_id, status);
```

### 초기 카테고리 (15개)

| slug | name | 설명 |
|------|------|------|
| `http-client` | HTTP Client | HTTP 요청 라이브러리 |
| `orm-database` | ORM / Database | 데이터베이스 ORM 및 쿼리 빌더 |
| `auth` | Authentication | 인증/인가 라이브러리 |
| `testing` | Testing | 테스트 프레임워크 및 유틸리티 |
| `cli-framework` | CLI Framework | CLI 앱 제작 프레임워크 |
| `logging` | Logging | 로깅 라이브러리 |
| `file-processing` | File Processing | 파일 처리/변환 |
| `validation` | Validation | 데이터 검증/스키마 라이브러리 |
| `date-time` | Date / Time | 날짜/시간 처리 |
| `state-management` | State Management | 상태 관리 라이브러리 |
| `web-framework` | Web Framework | 웹 서버/프레임워크 |
| `api-client` | API Client | 외부 API 클라이언트 SDK |
| `caching` | Caching | 캐싱 솔루션 |
| `queue-messaging` | Queue / Messaging | 메시지 큐 및 이벤트 처리 |
| `ai-ml-sdk` | AI / ML SDK | AI/ML 관련 SDK 및 도구 |

---

## 5. 데이터 수집 파이프라인

### 수집기 인터페이스 (Go)

```go
type Collector interface {
    Name() string
    Collect(ctx context.Context, lib model.Library) (*CollectionResult, error)
}
```

### 수집기별 데이터

**GitHubCollector** (GitHub REST + GraphQL API)
- 스타/포크 수, 워쳐 수
- 최근 커밋 이력 (90일)
- 이슈/PR 응답 시간 (중앙값)
- 릴리스 빈도, 컨트리뷰터 수 & 활동 분포
- README 존재 여부 & 길이, LICENSE 타입
- 아카이브 여부 (deprecated 감지)

**NpmCollector** (npm registry API + bundlephobia API)
- 주간 다운로드 수 & 추세
- 의존성 수, TypeScript 타입 포함 여부, 번들 사이즈
- 최신 버전 & 게시일, deprecated 플래그, dependents 수

**PyPICollector** (Phase 2 — PyPI JSON API + pypistats API)
- 다운로드 수, Python 버전 호환성, 타입 힌트 지원
- 최신 버전 & 게시일, classifiers, 의존성 수

**SecurityCollector** (Phase 2 — OSV API)
- 알려진 취약점 (CVE), 심각도, 마지막 보안 업데이트

### 파이프라인 흐름

```
[Scheduler (cron/asynq)] --> [Orchestrator]
                                   |
                   ┌───────────────┼───────────────┐
                   v               v               v
           [GitHub Task]    [NPM Task]      [PyPI Task]
                   |               |               |
                   v               v               v
           (goroutine)      (goroutine)     (goroutine)
                   |               |               |
                   └───────┬───────┘               |
                           v                       v
                   [Security Task] <───────────────┘
                           |
                           v
                   [Scoring Task] --> DB 저장 + 캐시 무효화
```

> Go의 goroutine으로 수집기를 병렬 실행. asynq는 전체 파이프라인 스케줄링/재시도/DLQ 관리.

### Rate Limit 전략
- **GitHub MVP:** 단일 PAT (5,000 req/hr). GraphQL 선호
- **GitHub v1:** GitHub App으로 업그레이드
- npm/PyPI: politeness delay 500ms
- 모든 수집기: 지수 백오프 재시도 (3회), circuit breaker

### 부분 실패 정책

| 항목 | 정책 |
|------|------|
| 개별 라이브러리 실패 | 해당 라이브러리 scoring 건너뛰고 이전 유효 점수 유지 (stale-while-revalidate) |
| 성공 기준 | >= 80% 성공: `completed_with_errors`, < 80%: `failed` |
| 실패 추적 | `collection_run_items` 테이블에 수집기별 상태/에러/소요시간 기록 |
| Dead Letter Queue | 3회 연속 실패한 라이브러리는 DLQ로 이동 |

---

## 6. 점수/랭킹 알고리즘

### Deprecation 프리필터 (가중 평균 이전)

```go
func PreFilter(lib model.Library) (deprecated bool, reason string) {
    if lib.Archived         { return true, "GitHub repository archived" }
    if lib.NpmDeprecated    { return true, "npm deprecated flag" }
    if containsDeprecated(lib.Readme) { return true, "README indicates deprecated" }
    return false, ""
}
// deprecated인 경우: overall_score = 0, 가중 평균 건너뜀
```

### 6개 메트릭 (각 0-100, higher=better)

#### 6.1 유지보수 건강도 (Maintenance Health) — 25%
- 최근 커밋 (90일 내): 0-30점
- 릴리스 빈도: 0-25점
- 이슈 응답 시간 (중앙값): 0-25점
- PR 머지 속도: 0-20점

#### 6.2 API 명확성 (API Clarity) — 20%
- TypeScript 타입 / Python 타입 힌트: 0-40점
- 의존성 수 (적을수록 높은 점수): 0-20점
- README 코드 예제 수: 0-20점
- API 문서화 비율 (휴리스틱): 0-20점

#### 6.3 문서 품질 (Documentation Quality) — 15%
- README 존재 & 길이: 0-25점
- 전용 문서 사이트: 0-25점
- CHANGELOG: 0-15점
- LICENSE: 0-15점
- 예제/튜토리얼 디렉토리: 0-20점

#### 6.4 보안 태세 (Security Posture) — 15%
- 알려진 취약점: 0-50점
- 마지막 보안 업데이트: 0-25점
- 의존성 깊이: 0-25점

#### 6.5 커뮤니티 시그널 (Community Signal) — 15%
- GitHub 스타 (로그 스케일): 0-30점
- 주간 다운로드 (로그 스케일): 0-30점
- 컨트리뷰터 수: 0-20점
- dependents 수: 0-20점

#### 6.6 Deprecation Safety (폐기 안전도) — 10%
- 최근 커밋 활동 (1년 기준): 0-40점
- 최근 릴리스 (2년 기준): 0-30점
- 유지보수자 수 (bus factor): 0-30점

### 전체 점수 계산

```go
overall = maintenance*0.25 + apiClarity*0.20 + docQuality*0.15 +
          security*0.15 + community*0.15 + deprecationSafety*0.10
```

### 투명성
- 모든 원시 데이터는 `scores.raw_data` JSONB에 저장
- 알고리즘 버전(`scoring_version`)으로 변화 추적
- API 응답에 점수 분해(breakdown) 항상 포함

---

## 7. API 설계

### 기본 원칙
- JSON 응답, 일관된 봉투(envelope) 형식
- 버전: URL 경로 (`/v1/...`)
- 공개 읽기 전용 API (MVP부터), Rate Limiting 100 req/min by IP
- cursor 기반 페이지네이션
- `next_actions` HATEOAS 필드 (에이전트 워크플로우 안내)

### HTTP 캐시

| 엔드포인트 유형 | Cache-Control | 비고 |
|----------------|---------------|------|
| 목록/카테고리 | `max-age=3600` | 수집 완료 시 무효화 |
| 점수 | `max-age=86400` | 주간 갱신 |
| 추천 | `max-age=3600` | 점수 기반 |
| 헬스 | `no-cache` | 항상 최신 |

### 엔드포인트

```
Phase 1 (MVP):
  GET /v1/libraries                           카테고리별 상위 라이브러리
  GET /v1/libraries/:id                       단일 라이브러리 상세
  GET /v1/libraries/:registry/:package_name   slug 기반 접근 (에이전트 친화적)
  GET /v1/categories                          전체 카테고리
  GET /v1/categories/:slug                    카테고리별 랭킹
  GET /v1/search?q=...                        텍스트 검색
  GET /v1/recommend?task=...&prefer=...       AI 에이전트 전용 추천 (핵심)
  GET /v1/health                              헬스 체크
  GET /v1/scoring/weights                     점수 가중치 공개

Phase 2:
  GET /v1/libraries/:id/scores                점수 이력
  GET /v1/compare?ids=...                     라이브러리 비교
```

### /v1/recommend MVP 구현 전략
1. `task` 파라미터에서 키워드 추출
2. 키워드 → 카테고리 딕셔너리로 매칭
3. `prefer` 파라미터에 따른 정렬 (lightweight/stable/secure/popular)
4. 매칭된 카테고리에서 top-N 반환

### 응답 예시
```json
{
  "data": {
    "id": "uuid",
    "name": "axios",
    "registry": "npm",
    "packageName": "axios",
    "githubRepo": "axios/axios",
    "latestVersion": "1.7.2",
    "deprecated": false,
    "score": {
      "overall": 82.5,
      "breakdown": {
        "maintenanceHealth": 85,
        "apiClarity": 90,
        "docQuality": 80,
        "securityPosture": 75,
        "communitySignal": 95,
        "deprecationSafety": 88
      },
      "version": "1.0.0"
    }
  },
  "next_actions": [
    { "rel": "compare", "href": "/v1/compare?ids=uuid1,uuid2" },
    { "rel": "category", "href": "/v1/categories/http-client" }
  ]
}
```

---

## 8. 프론트엔드

### 페이지
```
/                       랜딩 (소개, 주요 통계, 인기 카테고리)
/categories             전체 카테고리 그리드
/categories/[slug]      카테고리 내 랭킹 테이블
/library/[id]           라이브러리 상세 (점수 분해, 차트, 대안)
/compare                라이브러리 비교 (Phase 2)
/search                 검색 결과
/about                  프로젝트 소개, 알고리즘 설명
/api-docs               API 문서
```

### Go API와 타입 연동
- Go API에서 OpenAPI 스펙 자동 생성 (chi-openapi 또는 swag)
- `openapi-typescript-codegen`으로 프론트엔드 타입 자동 생성
- Next.js Server Components에서 Go API 직접 호출 (SSR)

---

## 9. 배포

### MVP
```
GitHub Actions (CI/CD)
        |
        v
[Docker Build] --> [Railway / Fly.io]
        |                    |
        ├── api (Go binary)  ├── PostgreSQL (managed)
        ├── worker (Go binary)└── Redis (managed)
        └── web (Next.js)
```
- Go 바이너리: ~15MB Docker 이미지 (scratch/distroless 기반)
- 비용: 월 $5-20 수준

---

## 10. 구현 로드맵

### Phase 1: MVP
**목표:** 핵심 파이프라인 작동, npm 라이브러리 100개 큐레이션

| 단계 | 작업 | 수락 기준 |
|------|------|-----------|
| 1-1 | 프로젝트 스캐폴딩 | Go 모듈 초기화, cmd/api + cmd/worker + internal/ 구조, Makefile, Docker Compose, `make build` 통과 |
| 1-2 | DB 스키마 & 마이그레이션 | golang-migrate 마이그레이션 파일, sqlc 쿼리 정의 및 코드 생성, `make migrate-up` 성공 |
| 1-2.5 | 시드 데이터 | 15개 카테고리 + 카테고리당 5-10개 npm 라이브러리 JSON, `make seed` 실행하여 DB 적재 성공 |
| 1-3 | GitHub + npm 수집기 | goroutine 병렬 수집, rate limit, 에러 핸들링, collection_run_items 기록 |
| 1-4 | 점수 산정 엔진 | deprecation 프리필터, 6개 메트릭, 전체 점수 산출, `go test` 통과 |
| 1-5 | REST API 핵심 | Chi 라우터, 9개 엔드포인트 (/recommend 포함), rate limiting, Cache-Control, next_actions, OpenAPI |
| 1-6 | 프론트엔드 기본 UI | 카테고리 목록, 라이브러리 랭킹 테이블, 상세 페이지 |

### Phase 2: v1
**목표:** PyPI 지원, 비교 기능, 자동 파이프라인

| 단계 | 작업 |
|------|------|
| 2-1 | PyPI 수집기 |
| 2-2 | 보안 수집기 (OSV) |
| 2-3 | 비교 API (/compare) |
| 2-4 | 프론트엔드 고도화 (차트, 비교, 검색) |
| 2-5 | asynq 스케줄드 파이프라인 (주간 자동 수집) |
| 2-6 | GitHub App 업그레이드 |

### Phase 3: v2
- MCP 서버 구현
- LLM 기반 /recommend 고도화
- Agent Compatibility 메트릭 재도입
- Go/Rust 레지스트리 추가
- 사용자 기여 시스템, 대시보드

---

## 11. 리스크

| 리스크 | 심각도 | 완화 |
|--------|--------|------|
| GitHub API Rate Limit | HIGH | MVP: 단일 PAT + GraphQL. v1: GitHub App |
| 점수 알고리즘 신뢰도 | MEDIUM | 버전 관리, 원시 데이터 공개, 커뮤니티 피드백 |
| "API 명확성" 주관성 | MEDIUM | v1 휴리스틱, v2 LLM 분석 |
| Deprecated 미감지 | MEDIUM | 프리필터 3중 시그널 + deprecation_safety 메트릭 |
| 파이프라인 부분 실패 | MEDIUM | stale-while-revalidate, 80% 기준, DLQ |
| GitHub API ToS | MEDIUM | 파생 점수만 공개, 출처 명시, 초기 법적 검토 |

---

## 12. 가드레일

### Must Have
- 모든 API 응답에 점수 breakdown 포함
- 점수 알고리즘 버전 추적
- Deprecation 프리필터
- `(registry, package_name)` 복합 유니크
- `internal/scoring/`이 점수 산정의 단일 정식 위치
- Rate Limiting + Cache-Control (MVP부터)
- `next_actions` 에이전트 워크플로우 안내
- 핵심 로직 테스트 커버리지 80%+

### Must NOT Have
- 사용자 인증 (MVP)
- 실시간 수집 (배치 기반만)
- 자체 패키지 호스팅
- Agent Compatibility 메트릭 (MVP, v2에서 재도입)

---

## 13. 성공 기준

1. **MVP:** npm 100개+ 큐레이션, 15개 카테고리, API 200ms 이하, /recommend MVP 동작
2. **v1:** npm + PyPI 500개+, 비교 API, 주간 자동 수집, 성공률 80%+
3. **품질:** /recommend에서 deprecated 추천률 0%, top-3 정확도 80%+
