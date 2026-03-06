# 의사결정 로그 & 미해결 사항

## 확정된 결정

### 기술 스택
| 결정 | 선택 | 이유 | 날짜 |
|------|------|------|------|
| 백엔드 언어 | **Go** | goroutine 병렬 수집, 단일 바이너리 배포, 개발자 숙련도 | 2026-03-06 |
| HTTP 라우터 | **Chi** | 경량, net/http 호환, 미들웨어 체이닝 | 2026-03-06 |
| DB 접근 | **sqlc** | SQL → 타입 안전 Go 코드 생성, 복잡한 쿼리에 적합 | 2026-03-06 |
| 마이그레이션 | **golang-migrate** | SQL 직접 관리, CLI + 라이브러리 | 2026-03-06 |
| 작업 큐 | **asynq** | Redis 기반, Go 네이티브, 재시도/스케줄링 | 2026-03-06 |
| 프론트엔드 | **Next.js 15** | SSR/SSG, React Server Components | 2026-03-06 |
| 프론트엔드 타입 연동 | **OpenAPI codegen** | Go API → OpenAPI 스펙 → TS 타입 자동 생성 | 2026-03-06 |

### 데이터 모델
| 결정 | 내용 | 날짜 |
|------|------|------|
| 라이브러리 유니크 키 | `(registry, package_name)` 복합 유니크 | 2026-03-06 |
| 점수 메트릭 | 6개 (Agent Compatibility는 v2에서 재도입) | 2026-03-06 |
| Deprecation 처리 | 프리필터 (archived/deprecated → score=0) + deprecation_safety 메트릭 | 2026-03-06 |
| alternatives 방향성 | source_library_id → target_library_id + reason | 2026-03-06 |
| 부분 실패 정책 | stale-while-revalidate, >=80% 성공 기준, DLQ | 2026-03-06 |

### 제품
| 결정 | 내용 | 날짜 |
|------|------|------|
| MVP 레지스트리 | npm 우선, PyPI는 Phase 2 | 2026-03-06 |
| 초기 카테고리 | 15개 확정 (docs/PLAN.md 참조) | 2026-03-06 |
| 시드 데이터 | JSON 파일 + `make seed` 스크립트, awesome-nodejs + npm 다운로드 기반 | 2026-03-06 |
| GitHub PAT | MVP: 단일 PAT (5,000 req/hr), v1: GitHub App | 2026-03-06 |
| API 공개 | MVP부터 공개 읽기 전용, Rate limiting 100 req/min by IP | 2026-03-06 |
| /v1/recommend | Phase 1에 포함, 규칙 기반 키워드→카테고리 매핑 | 2026-03-06 |

---

## 미해결 사항 (Open Questions)

### Phase 1 구현에 영향 없음 — 구현하면서 자연스럽게 결정

이 항목들은 코드 작성을 시작하기 전에 결정할 필요가 없습니다.
구현 과정에서 데이터를 보며 판단하거나, 배포 시점에 결정하면 됩니다.

| # | 항목 | 왜 지금 결정 안 해도 되는가 | 결정 시점 |
|---|------|--------------------------|----------|
| 1 | "API Clarity" 정량화 방법 | 휴리스틱(타입 존재/의존성 수/README 예제)이 이미 정의됨. 실제 데이터로 채점해보고 조정 | Phase 1 구현 중 (1-4 단계) |
| 2 | 데이터 수집 주기 (주간/격주/월간) | 기본값 주간. 실제 GitHub API 사용량 측정 후 조정 | Phase 2 자동화 시 |
| 3 | PostgreSQL 호스팅 (Railway/Supabase/Neon) | 로컬은 Docker Compose. 프로덕션 배포 시 선택 | 배포 직전 |
| 4 | alternatives 수동 vs 자동 매핑 | MVP는 시드 데이터에 수동 매핑. 100개 규모에서 수동이 충분 | Phase 2-3 |
| 5 | /recommend 키워드 딕셔너리 범위 | 15개 카테고리에 대해 구현하면서 키워드 추가. 초기에 완벽할 필요 없음 | Phase 1 구현 중 (1-5 단계) |
| 6 | 배포 플랫폼 (Railway/Fly.io) | 코드와 무관. Go 바이너리는 어디든 배포 가능 | 배포 직전 |
| 7 | 모니터링/알림 채널 | v1 범위. MVP는 로그 + health 엔드포인트로 충분 | Phase 2 |
| 8 | 도메인명 | 배포 시 결정 | 배포 직전 |
| 9 | GitHub API ToS 준수 | 코딩에는 영향 없음. 파생 점수만 공개하는 구조이므로 리스크 낮음 | 퍼블릭 배포 전 |
| 10 | 데이터 재배포 정책 | 위와 동일, 퍼블릭 공개 전 확인 | 퍼블릭 배포 전 |

---

## 변경 이력

| 날짜 | 변경 내용 |
|------|----------|
| 2026-03-06 | 초기 계획 수립 (Planner-Architect-Critic 2회 반복 합의) |
| 2026-03-06 | 백엔드 TypeScript → Go 전환 결정 |
