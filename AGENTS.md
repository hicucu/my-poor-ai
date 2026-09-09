# AGENTS.md — my-poor-ai 에이전트 명세

이 문서는 my-poor-ai에 포함된 19개 에이전트의 역할·입출력 계약·호출 관계 기술.
사용법(스킬 트리거, 커맨드)은 `CLAUDE.md` 참조.

---

## 에이전트 풀 개요

| 그룹             | 에이전트 수 | 소속 스킬                                             | 호출 방식                                                  |
| ---------------- | ----------- | ----------------------------------------------------- | ---------------------------------------------------------- |
| 공통 인프라      | 1개         | `using-my-poor-ai`                                    | using-my-poor-ai 복잡 경로 Phase 0, 신규 요구사항 수신 시 호출 |
| docs-suite       | 5개         | `sync-docs-from-diff`                                 | 스킬 오케스트레이터 → 병렬 fan-out                         |
| feature-pipeline | 9개         | `feature-pipeline`                                    | Phase 1→5 순차 (Phase 내 일부 병렬)                        |
| subagent-driven 플로우 | 4개   | `using-my-poor-ai` (복잡 경로 FULL)                        | main agent가 `subagent_type`으로 직접 스폰                 |

전체 합계: **1 + 5 + 9 + 4 = 19개**. 모든 에이전트 정의는 `agents/*.md`에 위치함.

**호출 방식 2가지**: ① 스킬/커맨드가 에이전트 정의 파일을 읽혀 지침으로 주입 (`{팀_위치}/agents/X.md`), ② `subagent_type`으로 직접 스폰. 리뷰어·aggregator·issue-fixer·project-context는 두 방식 모두에서 동일 계약으로 동작함. subagent-driven 플로우 4개는 `-agent` 접미사를 유지해 유사 역할의 feature-pipeline 세트(feature-planner/file-developer)와 구분함.

**모델 배정**: 각 에이전트 파일의 frontmatter `model` 필드가 정본임 (이 문서에는 개별 값을 중복 기재하지 않음). 배정 원칙: 파일 단위 팬아웃 워커(file-developer, test-writer, issue-fixer 등)는 haiku, 판단 중심 리뷰·계획·합성은 역할 무게에 따라 sonnet/opus.

---

## 공통 인프라 에이전트 (1개)

### project-context.md

모든 기능개발 시작 전 프로젝트 상태 캡처.

**호출 시점**: using-my-poor-ai 복잡 경로 Phase 0, 신규 요구사항 수신
**입력**: `CWD`, `mode` (full | stack-only), `output` 경로
**출력**: `_workspaces/project-context.md`

```
탐색 범위 (mode별):
  full:       스택 감지 + 구조 트리 + 진입점 + 컨벤션 + 최근 커밋 5개
  stack-only: 마커 파일 재스캔 (스택·의존성 변경 감지)
```

**완료 조건**: `_workspaces/project-context.md` 생성 완료 후 오케스트레이터에 경로 보고

---

## docs-suite 에이전트 그룹 (5개)

### sync-docs-from-diff 서브그룹 (5개)

```
Phase 1 (순차)
└── change-analyzer.md
        → _workspaces/01_change_analysis.json
        → _workspaces/01_change_analysis.md

Phase 2 (병렬, 파일 수정 없음)
├── readme-updater.md      → _workspaces/proposals/readme/
├── docs-updater.md        → _workspaces/proposals/docs/
└── inline-doc-updater.md  → _workspaces/proposals/inline/

Phase 3 (오케스트레이터 직접 수행 — 에이전트 없음)
    사용자 승인 후 Edit으로 패치 적용
    → _workspaces/03_apply_log.md

Phase 4 (순차)
└── doc-sync-validator.md  → _workspaces/02_validation_report.md
```

**에이전트별 입력 계약**:

| 에이전트           | 입력                                                    | 출력                                        |
| ------------------ | ------------------------------------------------------- | ------------------------------------------- |
| change-analyzer    | `BASE_BRANCH`, `WORKSPACE_DIR`, `PROJECT_ROOT`          | `01_change_analysis.{json,md}`              |
| readme-updater     | `WORKSPACE_DIR`, `PROJECT_ROOT`, `01_change_analysis.*` | `proposals/readme/*.patch.md` + `_index.md` |
| docs-updater       | 위와 동일                                               | `proposals/docs/*.patch.md` + `_index.md`   |
| inline-doc-updater | 위와 동일                                               | `proposals/inline/*.patch.md` + `_index.md` |
| doc-sync-validator | `WORKSPACE_DIR`, `PROJECT_ROOT`                         | `02_validation_report.md`                   |

**영역 분리 원칙** (각 에이전트가 담당하는 파일 범위):

| 에이전트           | 담당 영역                                             | 제외                |
| ------------------ | ----------------------------------------------------- | ------------------- |
| readme-updater     | `<PROJECT_ROOT>/README.md` 단일 파일                  | docs/, 인라인       |
| docs-updater       | `<PROJECT_ROOT>/docs/` 하위 전체 (locale README 포함) | 루트 README, 인라인 |
| inline-doc-updater | 변경 파일의 부모 dir 1~2단계 `*.md`                   | 루트 README, docs/  |

---

## feature-pipeline 에이전트 그룹 (9개)

### Phase 흐름

```
Phase 1 (순차)
└── feature-planner.md
    입력: 사용자 기능 요구사항
    출력: stack-profile.json, plan.md, file-manifest.json
    위치: _workspaces/{workspaceName}/
         ↓ 사용자 승인 게이트
Phase 2 (파일별 병렬, run_in_background: true)
└── file-developer.md × N
    입력: stack-profile.json, file-manifest.json의 각 파일 항목
    출력: 구현된 소스 파일
         ↓
Phase 3 (순차)
└── test-writer.md
    입력: stack-profile.json, file-manifest.json의 businessLogicFiles
    출력: 단위테스트 파일
         ↓
Phase 4a (4개 병렬, run_in_background: true)
├── architecture-reviewer.md  → _workspaces/review-{branch-slug}/reviews/architecture.md
├── security-reviewer.md      → _workspaces/review-{branch-slug}/reviews/security.md
├── performance-reviewer.md   → _workspaces/review-{branch-slug}/reviews/performance.md
└── style-reviewer.md         → _workspaces/review-{branch-slug}/reviews/style.md
         ↓
Phase 4b (순차)
└── review-aggregator.md      → _workspaces/review-{branch-slug}/review-report.md
         ↓ 사용자 확인 게이트
Phase 5 (파일별 병렬, run_in_background: true)
└── issue-fixer.md × M
    입력: review-report.md의 파일별 이슈 + stack-profile.json
    출력: 수정된 소스 파일
```

### 에이전트별 입력·출력 계약

| 에이전트              | 주요 입력                                      | 주요 출력                                       |
| --------------------- | ---------------------------------------------- | ----------------------------------------------- |
| feature-planner       | 사용자 요구사항, CWD                           | stack-profile.json, plan.md, file-manifest.json |
| file-developer        | stack-profile.json, 파일 항목 1개              | 구현 파일                                       |
| test-writer           | stack-profile.json, businessLogicFiles         | 테스트 파일                                     |
| architecture-reviewer | 변경 파일, stack-profile.json                  | reviews/architecture.md                         |
| security-reviewer     | 변경 파일, stack-profile.json                  | reviews/security.md                             |
| performance-reviewer  | 변경 파일, stack-profile.json                  | reviews/performance.md                          |
| style-reviewer        | 변경 파일, stack-profile.json                  | reviews/style.md                                |
| review-aggregator     | 4개 reviews/\*.md                              | review-report.md                                |
| issue-fixer           | review-report.md 파일 섹션, stack-profile.json | 수정된 파일                                     |

---

## subagent-driven 플로우 에이전트 그룹 (4개)

`using-my-poor-ai`의 복잡 경로(FULL)에서 main agent가 `subagent_type`으로 직접 스폰하는 오케스트레이션 세트.

| 에이전트            | 역할                                                                                          | 주요 입력               | 주요 출력                          |
| ------------------- | --------------------------------------------------------------------------------------------- | ----------------------- | ---------------------------------- |
| brainstorming-agent | 요구사항 분석 + 설계안 2~3개 비교 후 design.md 작성                                            | 요구사항, 프로젝트 경로 | `design.md` + STATUS               |
| planning-agent      | 승인된 design.md를 TDD 태스크 단위 스펙으로 분해                                               | `design.md`             | `specs/*.md`, `file-manifest.json` |
| developer-agent     | 단일 스펙을 TDD(RED-GREEN-REFACTOR)로 구현 + 커밋 + HANDOFF 갱신                               | `spec-{x}.md`           | 코드 + 커밋 + STATUS               |
| review-agent        | 리뷰 오케스트레이터 — reviewer 4종·review-aggregator·issue-fixer를 `subagent_type`으로 스폰    | branch-slug, base       | `review-report.md` + STATUS        |

Phase 4~5의 리뷰·수정은 feature-pipeline 세트와 동일한 에이전트(architecture/security/performance/style-reviewer, review-aggregator, issue-fixer)를 재사용하며, 부리는 계약은 `_shared/review-orchestration.md`가 정본임 (불변식 10).

### 산출물 경로 정책

| 용도             | 경로                               |
| ---------------- | ---------------------------------- |
| 기능 개발 산출물 | `_workspaces/{workspaceName}/`      |
| 세션 인계 문서   | `_workspaces/{branch-slug}/HANDOFF.md` |
| 리뷰 산출물      | `_workspaces/review-{branch-slug}/` |
| 절대 경로 사용   | 금지 (`~/`, `/` 시작 경로)         |

### 산출물 파일 가이드

모든 산출물은 `_workspaces/` 하위에 생성되며 `.gitignore` 대상임 (저장소에 커밋되지 않는 작업용 문서).

**공용 (브랜치 무관)**

| 파일                 | 생성 주체                      | 역할                                   | 수명     |
| -------------------- | ------------------------------ | -------------------------------------- | -------- |
| `project-context.md` | project-context          | 스택·아키텍처·구조 분석                | 24h 캐시 |
| `stack-profile.json` | feature-planner                | 언어·프레임워크·테스트 프레임워크 감지 | 재생성   |

**복잡 경로(FULL) — `{branch-slug}/`**

| 파일                    | 생성 주체                | 역할                                     |
| ----------------------- | ------------------------ | ---------------------------------------- |
| `design.md`             | brainstorming            | 검증된 설계 — 사용자 승인 게이트         |
| `specs/spec-{a,b,…}.md` | writing-plans / planning | TDD 태스크 단위 구현 스펙                |
| `file-manifest.json`    | planning                 | 파일 단위 작업 + developmentOrder(병렬)  |
| `pipeline-state.md`     | 오케스트레이터           | Phase 완료 상태 체크박스 — 재개용        |
| `HANDOFF.md`            | 구현 주체                | spec/phase 완료 시 서술형 인계 맥락      |

**단순·디버깅 경로(SIMPLE·DEBUG) — `{branch-slug}/`**

| 파일      | 생성 주체            | 역할                                  |
| --------- | -------------------- | ------------------------------------- |
| `GOAL.md` | 작업자(소프트 게이트) | 목표 + 성공 기준 — 완료 검증 시 대조 |

**리뷰 — `review-{branch-slug}/`**

| 파일               | 생성 주체          | 역할                |
| ------------------ | ------------------ | ------------------- |
| `review-report.md` | review-aggregator  | 4종 병렬 리뷰 통합  |

**문서 동기화 (sync-docs-from-diff)**

| 파일                                          | 생성 주체        | 역할                |
| --------------------------------------------- | ---------------- | ------------------- |
| `01_change_analysis.json`                     | change-analyzer  | diff 기반 변경 분석 |
| `proposals/{inline,readme,docs}/*.patch.md`   | 문서 업데이터    | 문서 변경 제안      |
| `02_validation_report.md`                     | doc-sync-validator | 적용 일관성 검증  |
| `03_apply_log.md`                             | 오케스트레이터   | 승인·적용 로그      |

---

## 에이전트 간 의존성

```
그룹 간 호출: 없음 (docs-suite ↔ feature-pipeline 직접 참조 금지)
그룹 내 호출: 스킬 오케스트레이터만 호출, 에이전트 간 직접 호출 없음
  (예외: review-agent는 오케스트레이터 에이전트로서 reviewer 4종·review-aggregator·issue-fixer를 스폰)
외부 의존성: git CLI (모든 에이전트), pip/npm (project-setup 커맨드만)
```

---

## 공통 에이전트 작성 규칙

새 에이전트 추가 시 준수해야 할 불변식:

1. **단일 책임**: 에이전트 1개 = 산출물 타입 1개
2. **경로**: CWD 기준 상대 경로만 사용
3. **하드코딩 금지**: 레포 경로·사용자 환경 경로 미포함
4. **파일 수정 금지**: sync-docs-from-diff 그룹의 Phase 2 에이전트는 proposals/ 에 패치만 생성, 원본 파일 직접 수정 금지
5. **검증 통과**: 에이전트·스킬·커맨드 수정 후 `node scripts/validate-agents.mjs` 통과 필수 — frontmatter(도구 실명·name=파일명·필드 순서), 참조 해소, 코드펜스 균형을 검사하며 CI(push)에서도 실행됨
6. **항상 exit 0**: 에이전트 오류가 파이프라인 전체를 차단하지 않도록 1회 재시도 후 실패 보고
7. **산출물 명사형 종결**: 에이전트가 생성하는 한글 문서 문장은 명사형으로 종결 (CLAUDE.md 문서 작성 규칙 참조)
8. **Codex 미러 자동 생성**: `.codex/agents/*.toml`은 `agents/*.md`에서 자동 생성되는 파생물 — 수동 편집 금지. `agents/*.md` 수정 후 `node scripts/generate-codex-agents.mjs` 재실행 필수 (CI가 `--check`로 동기화 검증)
9. **공유 지침 모듈**: 구현 워커(developer-agent·file-developer)가 공통으로 쓰는 코딩 규율·스택 컨벤션 매트릭스는 `agents/_shared/implementation-conventions.md`에 단일 소스로 두고, 각 워커가 `@include: _shared/implementation-conventions.md`(agents/ 기준 상대 경로)로 참조함. Claude 런타임은 `{팀_위치}/agents/_shared/…`를 Read로 읽고, Codex 미러는 생성기가 `@include` 줄을 대상 파일 내용으로 인라인 확장함 — 따라서 `_shared/` 수정 후 codex 재생성 필수. validate가 `@include` 해소를, `--check`가 미러 동기화를 검사함. `_shared/`는 실행 에이전트가 아니므로 frontmatter·toml 미러 대상이 아님

10. **리뷰 오케스트레이션 단일 소스**: 리뷰어 4종·`review-aggregator`·`issue-fixer`를 부리는 계약(축 목록, 병렬 디스패치 규칙, 리뷰어 공통 입력, aggregator 입출력, 심각도별 후속 조치, 파일 단위 수정 규칙, 절대 금지)은 `agents/_shared/review-orchestration.md`가 정본임. `agents/review-agent.md`는 `@include`로, `skills/feature-pipeline/SKILL.md`는 Read 지시로 참조함(스킬은 codex 미러 대상이 아니라 `@include` 확장이 일어나지 않으므로).

    두 파이프라인은 **용도가 달라 각각 유지**하며, 아래 항목만 호출자가 정함 — 그 외를 호출자 문서에 다시 기술하지 않음:

    | 호출자 책임 | review-agent | feature-pipeline |
    | --- | --- | --- |
    | `{workspaceDir}` | `_workspaces/review-{branch-slug}/` | `_workspaces/{workspaceName}/` |
    | 리뷰 대상 | 브랜치 diff `{base}...HEAD` | Phase 2~3 변경 파일 목록 |
    | 호출 방식 | `subagent_type` 스폰 | 에이전트 정의 주입 |
    | 사용자 게이트 | 없음 (자동 진행) | 있음 (Phase 4 직후) |
    | issue-fixer 커밋 | `commit: true` | 커밋 안 함 |

    단독 커맨드 `/my-poor-ai:code-review`는 이 계약을 따르지 않음 — 보안 축을 네이티브 `/security-review`에 위임하고 스타일 축을 제거한 2축 구성임. 대화형 단독 실행은 네이티브 리뷰와의 중복 회피가 우선이고, 파이프라인 2종은 무인 실행이라 4축 자체 수행을 유지함. **드리프트가 아니라 문서화된 의도적 분기임**
