# Opensource Curator - TODO & Roadmap

## Priority: High

### Data Quality
- [x] **GITHUB_TOKEN 발급 및 적용** - 완료 (2026-03-20). 77개 라이브러리 전체 GitHub 메트릭 반영 점수 수집 성공
- [ ] **점수 알고리즘 고도화** - 현재 npm 데이터만 반영. GitHub stars, issues, commit frequency, contributor 수 등 반영 시 정확도 대폭 향상
- [ ] **LLM 기반 deprecated 판독** - 현재 GitHub archived + npm deprecated 플래그만 신뢰. README 문맥 분석으로 정확도 향상 가능 (키워드 매칭은 오탐 다수로 비활성화됨)
- [ ] **라이브러리 커버리지 확장** - 현재 77개 라이브러리. Python(PyPI), Rust(crates.io) 등 다른 생태계 추가

### Frontend
- [ ] **i18n 다국어 지원** - next-intl 기반, 키 참조 번역 시스템. 초기 언어: 영어/한국어. `[locale]/` 라우트 구조, messages/*.json 번역 파일, 언어 전환 UI 필요
- [ ] **검색 기능 개선** - 현재 단순 ILIKE 검색. Fuzzy search, 자동완성, 필터링 추가
- [ ] **반응형 모바일 최적화** - 카드 레이아웃 모바일 확인 및 개선
- [ ] **점수 트렌드 차트** - 시간에 따른 점수 변화 그래프 (collection_runs 데이터 활용)

## Priority: Medium

### AI Agent Integration
- [ ] **MCP 서버 구현** - AI 에이전트가 직접 호출할 수 있는 MCP(Model Context Protocol) 서버
- [ ] **OpenAPI 스펙 생성** - API 문서 자동화, AI 에이전트가 스키마를 읽고 활용
- [ ] **비교 API** - `GET /v1/compare?libs=axios,undici,got` 형태로 라이브러리 직접 비교
- [ ] **AI Agent 전용 응답 포맷** - 에이전트가 파싱하기 쉬운 structured output

### Backend
- [ ] **캐싱 레이어** - Redis 캐시로 API 응답 속도 개선
- [ ] **Rate limiting** - API abuse 방지
- [ ] **수집 파이프라인 에러 핸들링** - 개별 라이브러리 실패 시 전체 실패 방지 (현재 부분 구현)
- [ ] **npm 다운로드 트렌드 분석** - 주간/월간 다운로드 추이로 성장세 파악

### Scoring
- [ ] **TypeScript 타입 품질 점수** - `@types/` 존재 여부, 빌트인 타입 지원
- [ ] **의존성 트리 분석** - 의존성 깊이, 번들 사이즈 영향
- [ ] **보안 취약점 연동** - Snyk/OSV API로 실시간 CVE 체크
- [ ] **Breaking change 빈도** - semver major 버전 히스토리 분석

## Priority: Low

### DevOps
- [ ] **Docker Compose 프로덕션 설정** - Nginx reverse proxy, SSL, health checks
- [ ] **CI/CD 파이프라인** - GitHub Actions로 테스트/빌드/배포 자동화
- [ ] **모니터링** - Prometheus metrics, Grafana 대시보드

### UX
- [ ] **다크/라이트 테마 토글** - 현재 다크 전용, 라이트 테마 옵션 추가
- [ ] **카테고리 아이콘** - 각 카테고리별 아이콘으로 시각적 구분
- [ ] **즐겨찾기/북마크** - 사용자별 관심 라이브러리 저장 (localStorage)
- [ ] **공유 기능** - 라이브러리 비교 결과 URL 공유

### Content
- [ ] **라이브러리별 AI 사용 가이드** - 에이전트가 해당 라이브러리를 잘 쓰려면 알아야 할 팁
- [ ] **카테고리별 선택 가이드** - "HTTP 클라이언트를 고를 때 AI 에이전트는..."
- [ ] **주간 큐레이션 뉴스레터** - 새로 뜨는 라이브러리, 점수 변동 알림

---

*Last updated: 2026-03-20*
