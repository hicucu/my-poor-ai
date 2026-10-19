# AGENTS.md — my-poor-ai 에이전트 명세

이 문서는 my-poor-ai에 포함된 16개 에이전트의 역할·입출력 계약·호출 관계 기술.
사용법(스킬 트리거, 커맨드)은 `CLAUDE.md` 참조.

---

## 에이전트 풀 개요

| 그룹                   | 에이전트 수 | 소속 스킬·커맨드                           | 호출 방식                                                                 |
| ---------------------- | ----------- | ------------------------------------------ | ------------------------------------------------------------------------- |
| 공통 인프라            | 1개         | `using-my-poor-ai`                         | 복잡 경로 Phase 0, 신규 요구사항 수신 시 호출                             |
| docs-suite             | 5개         | `sync-docs-from-diff`                      | 스킬 오케스트레이터 → 병렬 fan-out                                        |
| 리뷰 그룹              | 6개         | `review-agent` · `/my-poor-ai:code-review` | `review-agent`가 스폰 / 커맨드는 아키텍처·성능 리뷰어와 aggregator만 스폰 |
| subagent-driven 플로우 | 4개         | `using-my-poor-ai` (복잡 경로 FULL)        | main agent가 직접 스폰                                                    |

전체 합계: **1 + 5 + 6 + 4 = 16개**. 모든 에이전트 정의는 `agents/*.md`에 위치하며, 로드 대상은 `.claude-plugin/plugin.json`의 `agents` 목록이 정함 (불변식 11).

**호출 방식**: `subagent_type`으로 스폰함 — 예: `Agent(subagent_type="my-poor-ai:review-agent")`. frontmatter의 `model`·`tools`가 그대로 적용됨. 에이전트가 등록되지 않은 환경(타 플랫폼)에서만 범용 에이전트에 `${CLAUDE_PLUGIN_ROOT}/agents/<name>.md`를 읽혀 같은 입력으로 실행함 — 이때 읽은 파일 본문의 변수는 치환되지 않으므로 프롬프트에 플러그인 루트 절대 경로를 함께 전달함. subagent-driven 플로우 4개는 이름에 `-agent` 접미사가 붙음.

**모델 배정**: 각 에이전트 파일의 frontmatter `model` 필드가 정본임 (이 문서에는 개별 값을 중복 기재하지 않음). 배정 원칙: 탐색·문서 제안·파일 단위 수정 워커(project-context, 문서 업데이터 3종, issue-fixer)는 haiku, 판단 중심 리뷰·계획·합성은 역할 무게에 따라 sonnet/opus. 추론 강도(`effort`)는 diff 리뷰어 4종(architecture·security·performance·style)만 frontmatter에 `high`로 고정함 — 지정하지 않은 서브에이전트는 보통 부모 세션의 강도를 물려받으므로, max 세션에서는 리뷰어도 max로 돎. 2026-10-19 측정(결함 11개를 심은 diff, security-reviewer 1회씩): 상속(max) 90,892 토큰·496초, `high` 21,439 토큰·37초, 탐지는 둘 다 11/11. 나머지 에이전트는 고정하지 않고 호출 시 Agent의 `effort` 파라미터로 지정할 수 있음 (`subagent-driven-development` 「모델 선택」 참조).

---

## 공통 인프라 에이전트 (1개)

### project-context.md

모든 기능개발 시작 전 프로젝트 상태 캡처.

**호출 시점**: using-my-poor-ai 복잡 경로 Phase 0, 신규 요구사항 수신
**입력**: `CWD`, `mode` (full | stack-only), `output` 경로
**출력**: `_workspaces/project-context.md` + `_workspaces/stack-profile.json` (감지 매트릭스·스키마는 `_shared/stack-profile.md`)

```
탐색 범위 (mode별):
  full:       스택 감지 + 구조 트리 + 진입점 + 컨벤션 + 최근 커밋 5개
  stack-only: 마커 파일 재스캔 (스택·의존성 변경 감지)
```

**완료 조건**: 두 파일 생성 완료 후 오케스트레이터에 경로 보고

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

## 리뷰 에이전트 그룹 (6개)

`review-agent`(복잡 경로 Phase 4)가 부리는 리뷰어·aggregator·issue-fixer. `/my-poor-ai:code-review`는 이 중 아키텍처·성능 리뷰어와 aggregator만 씀 (불변식 10).

### review-agent 흐름

```
Step 3a (단일 응답 동시 호출)
├── architecture-reviewer.md   → _workspaces/review-{branch-slug}/reviews/architecture.md
├── /security-review (Skill)   → _workspaces/review-{branch-slug}/reviews/security.md
│     쓸 수 없거나 적용 불가면 security-reviewer.md로 폴백
├── performance-reviewer.md    → _workspaces/review-{branch-slug}/reviews/performance.md
└── style-reviewer.md          → _workspaces/review-{branch-slug}/reviews/style.md
         ↓
Step 3b (순차)
└── review-aggregator.md       → _workspaces/review-{branch-slug}/review-report.md
         ↓ Critical/High가 있으면 자동 진행
Step 5 (파일별 병렬 — 단일 응답 동시 호출)
└── issue-fixer.md × M         → 수정 + 커밋
```

### 에이전트별 입력·출력 계약

| 에이전트              | 주요 입력                                      | 주요 출력               |
| --------------------- | ---------------------------------------------- | ----------------------- |
| architecture-reviewer | 변경 diff 또는 파일, stack-profile.json        | reviews/architecture.md |
| security-reviewer     | 변경 diff 또는 파일, stack-profile.json — 폴백 | reviews/security.md     |
| performance-reviewer  | 변경 diff 또는 파일, stack-profile.json        | reviews/performance.md  |
| style-reviewer        | 변경 diff 또는 파일, stack-profile.json        | reviews/style.md        |
| review-aggregator     | 호출자가 전달한 축별 reviews/\*.md             | review-report.md        |
| issue-fixer           | review-report.md 파일 섹션, stack-profile.json | 수정된 파일 (+ 커밋)    |

---

## subagent-driven 플로우 에이전트 그룹 (4개)

`using-my-poor-ai`의 복잡 경로(FULL)에서 main agent가 `subagent_type`으로 직접 스폰하는 오케스트레이션 세트.

| 에이전트            | 역할                                                                                                     | 주요 입력               | 주요 출력                          |
| ------------------- | -------------------------------------------------------------------------------------------------------- | ----------------------- | ---------------------------------- |
| brainstorming-agent | 요구사항 분석 + 설계안 2~3개 비교 후 design.md 작성                                                       | 요구사항, 프로젝트 경로 | `design.md` + STATUS               |
| planning-agent      | 승인된 design.md를 TDD 태스크 단위 스펙으로 분해 (스키마 변경 시 마이그레이션 패턴 확인)                   | `design.md`             | `specs/*.md`, `file-manifest.json` |
| developer-agent     | 단일 스펙을 TDD(RED-GREEN-REFACTOR)로 구현 + 커밋 + HANDOFF 갱신                                          | `spec-{x}.md`           | 코드 + 커밋 + STATUS               |
| review-agent        | 리뷰 오케스트레이터 — 리뷰어 3종·네이티브 `/security-review`(불가 시 security-reviewer)·aggregator·issue-fixer를 부림 | branch-slug, base       | `review-report.md` + STATUS        |

리뷰·수정은 리뷰 그룹 에이전트가 맡으며, 부리는 계약은 `_shared/review-orchestration.md`가 정본임 (불변식 10).

### 산출물 경로 정책

| 용도               | 경로                                                               |
| ------------------ | ------------------------------------------------------------------ |
| 프로젝트 공용      | `_workspaces/project-context.md`, `_workspaces/stack-profile.json` |
| 복잡 경로 산출물   | `_workspaces/{branch-slug}/`                                       |
| 세션 인계 문서     | `_workspaces/{branch-slug}/HANDOFF.md`                             |
| 리뷰 산출물        | `_workspaces/review-{branch-slug}/`                                |
| 절대 경로 사용     | 금지 (`~/`, `/` 시작 경로)                                         |

### 산출물 파일 가이드

모든 산출물은 `_workspaces/` 하위에 생성되는 작업용 문서이며 저장소에 커밋하지 않음 — 대상 프로젝트의 `.gitignore`에 `_workspaces/`를 넣을 것을 권장함. 대상 프로젝트가 `_workspaces/`를 추적하는 경우에만 설계·인계 문서를 함께 커밋함.

**공용 (브랜치 무관)**

| 파일                 | 생성 주체                                   | 역할                                   | 수명       |
| -------------------- | ------------------------------------------- | -------------------------------------- | ---------- |
| `project-context.md` | project-context                             | 스택·아키텍처·구조 분석                | 24h 캐시   |
| `stack-profile.json` | project-context (`_shared/stack-profile.md`) | 언어·프레임워크·테스트 프레임워크 감지 | 24h 재사용 |

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

| 파일               | 생성 주체          | 역할            |
| ------------------ | ------------------ | --------------- |
| `review-report.md` | review-aggregator  | 축별 리뷰 통합  |

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
그룹 간 호출: 없음 (docs-suite ↔ 리뷰 그룹 직접 참조 금지)
그룹 내 호출: 스킬 오케스트레이터만 호출, 에이전트 간 직접 호출 없음
  (예외: review-agent는 오케스트레이터 에이전트로서 리뷰어·review-aggregator·issue-fixer를 스폰하고
   네이티브 /security-review를 Skill로 호출)
외부 의존성: git CLI (모든 에이전트)
```

---

## 공통 에이전트 작성 규칙

새 에이전트 추가 시 준수해야 할 불변식:

1. **단일 책임**: 에이전트 1개 = 산출물 타입 1개
2. **경로**: CWD 기준 상대 경로만 사용
3. **하드코딩 금지**: 레포 경로·사용자 환경 경로 미포함
4. **파일 수정 금지**: sync-docs-from-diff 그룹의 Phase 2 에이전트는 proposals/ 에 패치만 생성, 원본 파일 직접 수정 금지
5. **검증 통과**: 에이전트·스킬·커맨드 수정 후 `node scripts/validate-agents.mjs` 통과 필수 — frontmatter(도구 실명·name=파일명·필드 순서), 참조 해소, 코드펜스 균형, `plugin.json` `agents` 목록 동기화, 레거시 표기(구 플러그인 루트 플레이스홀더·서브에이전트에 없는 백그라운드 파라미터) 잔존을 검사하며 CI(push)에서도 실행됨
6. **항상 exit 0**: 에이전트 오류가 파이프라인 전체를 차단하지 않도록 1회 재시도 후 실패 보고
7. **산출물 명사형 종결**: 에이전트가 생성하는 한글 문서 문장은 명사형으로 종결 (CLAUDE.md 문서 작성 규칙 참조)
8. **Codex 미러 자동 생성**: `.codex/agents/*.toml`은 `agents/*.md`에서 자동 생성되는 파생물 — 수동 편집 금지. `agents/*.md` 수정 후 `node scripts/generate-codex-agents.mjs` 재실행 필수 (CI가 `--check`로 동기화 검증)
9. **공유 지침 모듈**: developer-agent가 쓰는 코딩 규율·스택 컨벤션·테스트 컨벤션 매트릭스는 `agents/_shared/implementation-conventions.md`, project-context가 쓰는 감지 매트릭스·`stack-profile.json` 스키마는 `agents/_shared/stack-profile.md`, review-agent가 쓰는 리뷰 계약은 `agents/_shared/review-orchestration.md`에 두고, 각 에이전트가 `@include: _shared/…`(agents/ 기준 상대 경로)로 참조함. Claude 런타임은 `${CLAUDE_PLUGIN_ROOT}/agents/_shared/…`를 Read로 읽고, Codex 미러는 생성기가 `@include` 줄을 대상 파일 내용으로 인라인 확장함 — 따라서 `_shared/` 수정 후 codex 재생성 필수. validate가 `@include` 해소를, `--check`가 미러 동기화를 검사함. `_shared/`는 실행 에이전트가 아니므로 frontmatter·toml 미러·`plugin.json` `agents` 목록 대상이 아님

10. **리뷰 오케스트레이션 계약**: 리뷰어·`review-aggregator`·`issue-fixer`를 부리는 계약(축 목록, 보안 축의 네이티브 우선·폴백, 병렬 디스패치 규칙, 리뷰어 공통 입력, aggregator 입출력, 심각도별 후속 조치, 파일 단위 수정 규칙, 절대 금지)은 `agents/_shared/review-orchestration.md`가 정본이며 `agents/review-agent.md`가 `@include`로 참조함. review-agent가 정하는 값: `{workspaceDir}` = `_workspaces/review-{branch-slug}/`, 리뷰 대상 = 브랜치 diff `{base}...HEAD`, 사용자 게이트 없음(자동 진행), issue-fixer `commit: true`.

    단독 커맨드 `/my-poor-ai:code-review`는 이 계약을 따르지 않음 — 보안·정확성 축을 네이티브 `/security-review`·`/code-review`에 위임하고 스타일 축을 제거한 구성임(직접 수행은 아키텍처·성능 2축). review-agent는 무인 실행이라 스타일 축까지 4축을 유지하고 보안 축만 네이티브 우선으로 둠. **드리프트가 아니라 문서화된 의도적 분기임**

11. **에이전트 등록**: `.claude-plugin/plugin.json`의 `agents` 목록이 기본 스캔(`agents/` 하위 모든 `.md` 재귀 로드)을 대체함 — 새 에이전트는 이 목록에 추가해야 로드되며, README·`_shared/` 같은 비에이전트 문서는 넣지 않음. 목록이 없으면 비에이전트 문서도 frontmatter 없이 파일명으로 로드되어 호출 가능한 가짜 에이전트가 됨. `validate-agents.mjs`가 목록과 frontmatter 보유 파일의 일치를 검사함. `claude plugin validate`는 목록과 무관하게 `agents/` 하위 전체를 검사하므로 비에이전트 문서의 frontmatter 경고는 남음 — 로드에는 영향 없음 (2026-10-19 실측)

12. **핸드오프 파일 명시**: 서브에이전트가 다음 단계로 넘기는 `.md` 산출물(리뷰어 축 파일·`review-report.md`·`project-context.md`·`01_change_analysis.md`·제안 패치·`02_validation_report.md`)은 에이전트 정의와 스폰 프롬프트 모두에 "다음 단계가 읽는 입력 파일 — 반드시 작성"이라고 명시함. Claude Code 서브에이전트 하네스는 보고서·요약·분석용 `.md` 작성을 금지하고 결과를 최종 메시지로 반환하게 하며, 다른 도구의 입력으로 쓰는 파일만 예외로 둠 — 명시가 없으면 파일이 만들어지지 않아 다음 단계가 읽을 것이 없음 (2026-10-19 실측, 2.1.303: security-reviewer 명시 없이 2/2 미작성, 명시 후 1/1 작성). 오케스트레이터는 산출 파일이 없으면 반환된 텍스트를 그 경로에 저장한 뒤 다음 단계를 호출함. 핸드오프 파일명은 `REPORT`·`SUMMARY`·`FINDINGS`·`ANALYSIS`로 시작하지 않음 — 서브에이전트의 Write가 이 이름을 차단함
