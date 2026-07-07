---
name: sync-docs-from-diff
description: "Analyzes the commits and diff from the develop (or specified base) branch up to the current branch to synchronize README.md / ./docs/ / inline documentation (*.md) near the changed code. Operates in review-then-apply mode — shows the proposed changes for each file and applies them after user approval. Trigger phrases: \"update docs compared to develop\", \"reflect branch changes in README/docs\", \"update documentation to match committed changes\", \"sync docs with this PR's changes\", \"docs sync\", and follow-up requests (\"redo just the docs part\", \"re-run README only\", \"restart from analysis\", \"re-run the validator only\"). Does not trigger for simple README writing, creating new documents, or change requests written without a diff."
---

# Sync Docs From Diff — 오케스트레이터

브랜치 변경 사항(`<base>..HEAD`)을 분석하여 프로젝트 내 문서들을 동기화함. **사용자 검토 후 적용 모드** — 모든 파일 변경은 사용자 승인 후에만 반영됨.

## 실행 모드

서브 에이전트 + 파이프라인 (필요 시 부분 재실행). 모델 전략: 분석·검증 에이전트(`change-analyzer`, `doc-sync-validator`)는 `model: "opus"`, 문서 업데이트 에이전트 3개(`inline-doc-updater`, `readme-updater`, `docs-updater`)는 `model: "haiku"`.

## 핵심 원칙

1. **읽기 → 제안 → 검토 → 적용** 순서. 어떤 에이전트도 사용자 승인 전 파일을 수정하지 않음.
2. 산출물은 `<CWD>/_workspaces/` 하위에 단일 디렉토리로 모이며, 후속 호출 시 재사용/누적됨.
3. 베이스 브랜치는 `develop > main > master` 순으로 자동 감지. 사용자가 인자로 명시하면 그 값 사용.
4. 영역 분리: README(루트) / docs(디렉토리) / inline(코드 인근) — 한 파일은 한 에이전트만 담당.

## Phase 0: 컨텍스트 확인 및 모드 판정

작업을 시작하기 전 다음을 확인함.

1. **CWD가 git 작업 트리인지 확인** — `git rev-parse --show-toplevel` 성공 여부. 실패 시 사용자에게 git 저장소가 아니라고 보고하고 종료.
2. **PROJECT_ROOT 결정** — 위 명령 결과를 절대 경로로 저장. 모든 후속 에이전트에 전달.
3. **WORKSPACE 결정** — `<CWD>/_workspaces/` (CWD ≠ PROJECT_ROOT인 경우에도 호출 위치 기준).
4. **재실행 모드 판정**:
   - `_workspaces/` 부재 → **초기 실행**
   - `_workspaces/01_change_analysis.json` 존재 + 사용자가 "전체 다시" 요청 → 기존을 `_workspaces_prev/`로 이동 후 **새 실행**
   - `_workspaces/` 존재 + 사용자가 부분 영역 지칭("docs만", "README만", "validator만", "분석부터 다시") → **부분 재실행** (해당 에이전트만 다시 호출)
   - `_workspaces/` 존재 + 사용자가 "이어서" / 별도 지시 없음 → 가장 마지막에 미완료된 Phase부터 재개

판정 결과를 한 줄로 사용자에게 보고하고 진행함 (예: `초기 실행 시작함. base=develop`).

## Phase 1: 변경 분석 (change-analyzer)

1. 베이스 브랜치 결정:
   - 사용자가 인자로 base 지정 → 그대로 사용
   - 아니면 `git rev-parse --verify develop` → 실패 시 `main` → 실패 시 `master` 순서로 시도
   - 모두 실패하면 사용자에게 베이스 브랜치를 묻고 입력 받음
2. `Agent` 호출:
   - `subagent_type`: `general-purpose`
   - `model`: `"opus"`
   - `description`: "변경 분석 보고서 생성"
   - `prompt`: 에이전트 정의(`{팀_위치}/agents/change-analyzer.md`)를 따르라는 지시 + `BASE_BRANCH`, `WORKSPACE_DIR`, `PROJECT_ROOT` 변수 전달
3. 산출물 검증: `_workspaces/01_change_analysis.json`과 `_workspaces/01_change_analysis.md` 모두 존재 확인.
4. 분석에 변경이 0건이면(빈 commits/files) "변경 없음, 종료" 보고 후 종료.

## Phase 2: 업데이트 제안 생성 (병렬, 파일 수정 없음)

3개 updater를 동시 호출함. 모두 `run_in_background: true`, `model: "haiku"`.

| 에이전트             | 영역                                                                           | 산출물                         |
| -------------------- | ------------------------------------------------------------------------------ | ------------------------------ |
| `inline-doc-updater` | 변경 파일의 부모 dir 1~2단계의 `*.md` (루트 README와 docs/ 제외)               | `_workspaces/proposals/inline/` |
| `readme-updater`     | `<PROJECT_ROOT>/README.md` 1개                                                 | `_workspaces/proposals/readme/` |
| `docs-updater`       | `<PROJECT_ROOT>/docs/` 하위 전체 (locale README 포함, 예: `docs/ko/README.md`) | `_workspaces/proposals/docs/`   |

각 호출 prompt에 다음을 포함:

- 자신의 에이전트 정의(`{팀_위치}/agents/<name>.md`)를 따르라는 지시
- 환경 변수: `WORKSPACE_DIR`, `PROJECT_ROOT`
- 강조: "파일을 직접 수정하지 말고 반드시 `proposals/<영역>/` 하위에 패치만 작성하라"

세 에이전트가 모두 완료될 때까지 대기. 한 에이전트가 실패해도 나머지는 진행 — 실패한 영역은 보고서에서 "패치 생성 실패, 사용자 수동 검토 필요"로 표기.

## Phase 3: 사용자 검토 및 적용

이 Phase는 **오케스트레이터(메인 Claude)가 직접 수행** — 별도 에이전트를 호출하지 않음. 상세 절차는 `references/review-protocol.md` 참조.

핵심 흐름:

1. 3개 영역의 `_index.md`를 읽고 사용자에게 한눈에 요약 제시 — 각 영역의 패치 개수와 영향 요약.
2. 패치 파일을 영역별·문서별로 순회하며 사용자에게 Before/After 제시 후 `AskUserQuestion`으로 승인/거부/보류 묻기.
3. 승인된 패치는 즉시 원본 파일에 `Edit`로 반영. 매 적용마다 `_workspaces/03_apply_log.md`에 기록.
4. 거부/보류는 로그에 기록하되 파일 미수정.
5. 사용자가 "이 영역 전체 일괄 승인" 또는 "전체 일괄 승인"을 요청하면 묻지 않고 모두 적용 (각 적용은 그래도 로그 기록).

> 적용은 항상 Edit/Write 도구로 직접 수행 — 외부 도구나 셸 명령으로 우회하지 않음.

## Phase 4: 검증 (doc-sync-validator)

**선행 조건**: `_workspaces/03_apply_log.md`가 존재해야 함. 해당 파일이 없으면(사용자가 모든 패치를 거부해 적용 내역이 0건인 경우 포함) 검증 에이전트를 호출하지 않고 "적용된 변경 없음, 검증 생략" 보고 후 Phase 5로 진행함.

`Agent` 호출:

- `subagent_type`: `general-purpose`, `model`: `"opus"`
- `description`: "문서 동기화 검증 보고서"
- prompt: `{팀_위치}/agents/doc-sync-validator.md`를 따르라는 지시 + `WORKSPACE_DIR`, `PROJECT_ROOT`

검증 보고서(`_workspaces/02_validation_report.md`)에 WARN/FAIL이 있으면 사용자에게 그 부분을 강조 보고하고 부분 재실행 옵션을 제시.

## Phase 5: 요약 보고

다음 형식으로 사용자에게 종합 보고:

```
# 문서 동기화 완료

- 베이스: <base> (<sha>) → HEAD (<sha>)
- 분석된 commit: N개, 변경 파일: M개
- 적용된 패치: README X / docs Y / inline Z (총 W개)
- 보류/거부: K개
- 검증 결과: PASS / WARN(상세) / FAIL(상세)

## 다음 단계 권장
- (검증 WARN이 있으면 권장 행동)
- 변경 파일을 git diff로 확인 후 commit 권장
```

## 에러 핸들링

| 상황                                | 처리                                                                               |
| ----------------------------------- | ---------------------------------------------------------------------------------- |
| git 저장소 아님                     | 즉시 종료, 사용자에게 안내                                                         |
| 베이스 브랜치 부재                  | 자동 감지 실패 시 사용자에게 입력 요청                                             |
| change-analyzer 실패                | 1회 재시도 → 재실패 시 종료, 사용자 보고                                           |
| updater 1개 실패                    | 나머지 2개는 계속 진행, 실패 영역은 보고서에 표시                                  |
| 패치 적용 시 Edit 실패(파일 변경됨) | 해당 패치 보류 처리, 사용자에게 "파일이 다른 도구로 변경되었으니 재분석 권장" 보고 |
| validator WARN/FAIL                 | 종료하지 않고 보고. 부분 재실행 안내.                                              |

## 후속 작업/부분 재실행 매핑

| 사용자 표현            | 실행할 Phase                                               |
| ---------------------- | ---------------------------------------------------------- |
| "전체 다시"            | Phase 1~5 (이전 \_workspaces는 \_workspaces_prev로 이동)     |
| "분석부터 다시"        | Phase 1~5                                                  |
| "README만 재실행"      | Phase 2의 readme-updater + Phase 3(README만) + Phase 4     |
| "docs만 다시"          | Phase 2의 docs-updater + Phase 3(docs만) + Phase 4         |
| "inline만"             | Phase 2의 inline-doc-updater + Phase 3(inline만) + Phase 4 |
| "검증만 다시"          | Phase 4만                                                  |
| "이어서" / 미완료 재개 | 마지막 미완료 Phase부터                                    |

## 데이터 흐름

```
[Phase 1] change-analyzer
    └─> _workspaces/01_change_analysis.{json,md}

[Phase 2] (병렬, 모두 위 분석을 입력으로)
    ├─> inline-doc-updater  ─> _workspaces/proposals/inline/*.patch.md + _index.md
    ├─> readme-updater      ─> _workspaces/proposals/readme/*.patch.md + _index.md
    └─> docs-updater        ─> _workspaces/proposals/docs/*.patch.md   + _index.md

[Phase 3] 오케스트레이터 (사용자 검토 후 Edit 적용)
    └─> _workspaces/03_apply_log.md  +  실제 파일 수정

[Phase 4] doc-sync-validator
    └─> _workspaces/02_validation_report.md

[Phase 5] 사용자 보고 (오케스트레이터 직접)
```

## 테스트 시나리오

### 시나리오 1: 정상 흐름 (초기 실행)

- 사용자: "develop과 비교해서 문서 동기화"
- 기대: Phase 0 모드=초기실행 → 베이스=develop 자동감지 → Phase 1 분석(N commits) → Phase 2 3개 updater 병렬 → Phase 3 패치별 사용자 승인 → Phase 4 검증 PASS → Phase 5 요약

### 시나리오 2: 변경 없음

- 사용자: "docs sync" (HEAD == base)
- 기대: Phase 1에서 빈 분석 감지 → Phase 2~5 생략 → "변경 없음" 보고 후 즉시 종료

### 시나리오 3: 부분 재실행 (README만)

- 사용자: "README만 다시 갱신해줘"
- 기대: Phase 0에서 \_workspaces 존재 확인 → Phase 1 스킵(분석 재사용) → Phase 2의 readme-updater만 재실행 → Phase 3에서 README 패치만 검토 → Phase 4 재실행 → 보고

### 시나리오 4: 베이스 브랜치 자동 감지 실패

- 상황: develop/main/master 모두 없음
- 기대: 사용자에게 베이스 브랜치 입력 요청 → 입력값으로 진행

### 시나리오 5: validator WARN

- 상황: README와 docs의 시그니처가 미세하게 다른 채로 적용됨
- 기대: Phase 4 보고서에 WARN 기록 → Phase 5에서 강조 보고 + "readme-updater 재실행 권장" 제안

## 참고

- `{팀_위치}/agents/*.md` — 각 에이전트의 입력/출력 프로토콜과 검증 체크리스트 (`{팀_위치}`는 my-poor-ai 플러그인 루트의 절대 경로)
- `references/review-protocol.md` — Phase 3 사용자 검토 절차 상세
