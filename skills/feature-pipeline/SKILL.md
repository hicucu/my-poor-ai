---
name: feature-pipeline
description: "Must be used when adding a feature to any project (regardless of language or framework) — a 5-phase multi-agent pipeline that automatically handles everything from requirements analysis to implementation, testing, review, and issue fixing. Use this skill for new feature development requests such as \"add a feature\", \"implement this\", \"build me ~\", \"in this code, ~\", \"add an endpoint\", \"button click\", \"modal\", \"form handling\", \"API integration\". Also applies to follow-up requests like \"just fix the issues\" or \"redo from the review step\"."
---

# Feature Pipeline — 스택 무관 기능 추가 파이프라인 오케스트레이터

임의 프로젝트(언어·프레임워크 무관)에 기능을 추가하는 5단계 파이프라인을 조율함.

스택별 분기는 모두 서브 에이전트 프롬프트 내부에서 처리함. 이 SKILL.md는 스택 어휘를 일절 사용하지 않음.

## 실행 모드: 서브 에이전트 (팬아웃)

파일별 병렬 개발은 파일 간 의존성이 없으므로 서브 에이전트 패턴 사용.
각 서브 에이전트는 독립 실행 후 결과를 오케스트레이터에 반환.

## Phase 0: 컨텍스트 판정

`_workspaces/` 하위에서 기존 작업 디렉토리 탐색 (`_workspaces/*/stack-profile.json` 패턴):

| 상태                                               | 처리                                                                               |
| -------------------------------------------------- | ---------------------------------------------------------------------------------- |
| 미존재                                             | **초기 실행** — Phase 1부터 전체 실행. `{workspaceDir}` 는 feature-planner가 결정  |
| 존재 + 새 요구사항                                 | **새 실행** — 기존 디렉토리를 `_workspaces_prev/`로 이동 후 초기화                  |
| 존재 + "이슈만 수정", "리뷰부터 다시" 등 후속 요청 | **부분 재실행** — `{workspaceDir}`를 stack-profile.json에서 읽어 해당 Phase만 실행 |

> **경로 규칙**: 모든 산출물은 `_workspaces/{workspaceName}/` 하위에 저장함.
> `{workspaceName}`은 feature-planner가 요구사항에서 생성한 kebab-case 슬러그 (예: `feature-login`, `order-service`).
> `_workspaces/` 루트 직접 저장 금지.

부분 재실행 판정 기준:

- `{workspaceDir}/review-report.md` 존재 + "이슈만 수정" 요청 → Phase 5만 실행
- `{workspaceDir}/plan.md` + `{workspaceDir}/file-manifest.json` 존재 + "개발만 다시" 요청 → Phase 2~5 실행
- `{workspaceDir}/stack-profile.json` 존재 + 동일 프로젝트 + 요구사항 일부 변경 → Phase 1.1부터 재실행, 스택 감지 생략

## Phase 1: 계획 수립 + 스택 감지

`{팀_위치}/agents/feature-planner.md`를 읽고 feature-planner 에이전트 실행.

```
{팀_위치}/agents/feature-planner.md를 읽고 그 지침에 따라 작업한다.

요구사항: [사용자 요청 원문]
프로젝트 경로: [프로젝트 루트 경로]
컨텍스트: [관련 파일 경로 (있으면)]
출력:
  - {workspaceDir}/stack-profile.json (Phase 1.0, workspaceDir 포함)
  - {workspaceDir}/plan.md            (Phase 1.1, - [ ] 체크박스 포함)
  - {workspaceDir}/file-manifest.json (Phase 1.1)
```

feature-planner는 내부적으로 두 단계를 수행:

1. **Phase 1.0 — 스택 감지**: 프로젝트 루트의 마커 파일을 스캔하여 `{workspaceDir}/stack-profile.json` 생성. 미감지 시 `primary: "unknown"`, `fallbackUsed: true`로 표기.
2. **Phase 1.1 — 요구사항 분해**: 스택 프로필 + 코드베이스 탐색 결과를 바탕으로 파일별 구현 명세 작성.

**→ 사용자 확인 필수**:

- stack-profile.json의 `primary`/`subtype`/`testFramework`가 적절한지 사용자에게 제시
- `fallbackUsed: true`인 경우 강제 중단 후 사용자에게 스택 명시 요청
- file-manifest.json의 파일 목록과 개발 순서를 사용자에게 제시하고 승인 받음
- 수정 요청 시 계획 재작성 후 재승인

## Phase 2: 파일별 병렬 개발

`{workspaceDir}/file-manifest.json`의 `developmentOrder` 그룹 순서로 실행.
같은 그룹의 파일은 `run_in_background: true`로 동시 실행.

```
# 각 파일별 에이전트 프롬프트
{팀_위치}/agents/file-developer.md를 읽고 그 지침에 따라 작업한다.

파일 경로: [file.path]
작업 유형: [file.action]  # create | modify
구현 명세: [file.spec]
의존 파일: [file.dependencies]  # 먼저 읽을 파일 목록
스택 프로필: {workspaceDir}/stack-profile.json
TDD 필수: my-poor-ai:test-driven-development 스킬을 사용하여 RED-GREEN-REFACTOR 사이클로 구현.
          코드 작성 전 반드시 실패하는 테스트를 먼저 작성하고 실패를 확인한 후 최소 코드로 통과시킬 것.
          테스트 없이 작성된 코드는 삭제 후 재작성.
```

한 그룹의 모든 에이전트 완료 후 다음 그룹 실행. 그룹 간 의존성 위반 방지.
**Phase 2 완료 후**: `{workspaceDir}/plan.md`의 `- [ ] Phase 2` 항목을 `- [x] Phase 2`로 업데이트.

> **부분 성공 정책**: Phase의 일부 파일만 성공한 경우 `- [~] Phase N (N/M 완료)` 로 표기하고 실패 파일 목록을 plan.md에 기록함. 전체 성공 시에만 `- [x]` 로 표기.

## Phase 3: 단위테스트 작성

`{workspaceDir}/file-manifest.json`의 `businessLogicFiles` 대상.

```
{팀_위치}/agents/test-writer.md를 읽고 그 지침에 따라 작업한다.

대상 파일: [businessLogicFiles 목록]
프로젝트 경로: [루트]
스택 프로필: {workspaceDir}/stack-profile.json
```

UI/뷰 계층 파일(`uiMarkers` 경로)은 테스트 대상에서 제외. 순수 비즈니스 로직만 대상.
**Phase 3 완료 후**: `{workspaceDir}/plan.md`의 `- [ ] Phase 3` 항목을 `- [x] Phase 3`로 업데이트.

## Phase 4: 코드 리뷰 (4개 전문 병렬 + 통합)

Phase 2~3에서 생성/수정된 전체 파일 대상. 4개 전문 reviewer를 병렬 실행한 후 aggregator로 통합 보고서 생성.

### Phase 4a: 4개 전문 reviewer 병렬 실행

```
# 4개 에이전트를 단일 응답에 동시에 호출 (run_in_background: true)

{팀_위치}/agents/architecture-reviewer.md를 읽고 그 지침에 따라 작업한다.
리뷰 대상: [Phase 2~3 변경 파일 전체 목록]
스택 프로필: {workspaceDir}/stack-profile.json
출력 경로: {workspaceDir}/reviews/architecture.md

{팀_위치}/agents/security-reviewer.md를 읽고 그 지침에 따라 작업한다.
리뷰 대상: [동일]
스택 프로필: {workspaceDir}/stack-profile.json
출력 경로: {workspaceDir}/reviews/security.md

{팀_위치}/agents/performance-reviewer.md를 읽고 그 지침에 따라 작업한다.
리뷰 대상: [동일]
스택 프로필: {workspaceDir}/stack-profile.json
출력 경로: {workspaceDir}/reviews/performance.md

{팀_위치}/agents/style-reviewer.md를 읽고 그 지침에 따라 작업한다.
리뷰 대상: [동일]
스택 프로필: {workspaceDir}/stack-profile.json
출력 경로: {workspaceDir}/reviews/style.md
```

### Phase 4b: aggregator로 통합

4개 산출이 모두 도착하면 aggregator 호출:

```
{팀_위치}/agents/review-aggregator.md를 읽고 그 지침에 따라 작업한다.

입력:
  {workspaceDir}/reviews/architecture.md
  {workspaceDir}/reviews/security.md
  {workspaceDir}/reviews/performance.md
  {workspaceDir}/reviews/style.md
스택 프로필: {workspaceDir}/stack-profile.json
출력: {workspaceDir}/review-report.md  (issue-fixer 입력 호환)
```

**→ 사용자 확인**: aggregator가 작성한 review-report.md를 사용자에게 제시. Critical이 없으면 사용자가 완료 처리 가능. Critical/High가 1건이라도 있으면 Phase 5 권장.
**Phase 4 완료 후**: `{workspaceDir}/plan.md`의 `- [ ] Phase 4` 항목을 `- [x] Phase 4`로 업데이트.

## Phase 5: 이슈 수정

`{workspaceDir}/review-report.md` 기준, 파일 단위로 병렬 수정. `run_in_background: true`.

```
# 각 이슈 파일별 에이전트 프롬프트
{팀_위치}/agents/issue-fixer.md를 읽고 그 지침에 따라 작업한다.

파일 경로: [이슈 대상 파일]
이슈 목록: [해당 파일의 리뷰 이슈 발췌]
스택 프로필: {workspaceDir}/stack-profile.json
```

완료 후 수정 결과 요약 보고 (Critical/High/Medium/Low 처리 건수, 미적용 항목 사유 포함).
**Phase 5 완료 후**: `{workspaceDir}/plan.md`의 `- [ ] Phase 5` 항목을 `- [x] Phase 5`로 업데이트.

## 에러 핸들링

| 상황                                  | 대응                                                               |
| ------------------------------------- | ------------------------------------------------------------------ |
| feature-planner 실패                  | 에러 내용 보고, 1회 재시도                                         |
| stack-profile fallbackUsed=true       | 사용자에게 스택 명시 요청 후 profile 수동 보정, Phase 1.1부터 재개 |
| file-developer 1개 실패               | 나머지 완료 후 해당 파일 재시도, 재실패 시 수동 처리 안내          |
| test-writer 실패 (지원 스택 미스매치) | 테스트 없이 진행, 보고서에 누락 명시                               |
| code-reviewer 실패                    | 리뷰 없이 완료 처리, 사용자에게 알림                               |

## 사용자 입력 어휘 매핑

| 사용자 표현                              | 처리 모드                                                           |
| ---------------------------------------- | ------------------------------------------------------------------- |
| "X 기능 추가해줘", "엔드포인트 만들어줘" | 초기 실행 (전체 5 Phase)                                            |
| "이슈만 수정해줘", "리뷰대로 고쳐줘"     | Phase 5만 실행                                                      |
| "리뷰 다시 받아줘"                       | Phase 2 구현 파일은 기존 유지, Phase 3(테스트)·Phase 4(리뷰) 재실행 |
| "코드 고치고 리뷰"                       | Phase 2 파일 수정 후 Phase 3부터 재실행                             |
| "처음부터 다시", "다른 요구사항으로"     | `_workspaces_prev/` 백업 후 초기 실행                                |
| "스택 잘못 감지됐어, X로 바꿔줘"         | stack-profile.json 수동 보정 후 Phase 1.1부터 재개                  |

## 테스트 시나리오

**정상 흐름 (스택 자동 감지)**: "주문 생성 기능 + 재고 검증 로직 추가"
→ Phase 1.0: 프로젝트 매니페스트·의존성 스캔 → stack-profile.json
→ Phase 1.1: 요청 처리 / 스키마 / 비즈니스 로직 계층 분해
→ 사용자 승인
→ Phase 2: 3개 파일 병렬 개발
→ Phase 3: service 단위테스트 작성
→ Phase 4: 코드 리뷰
→ Phase 5: 이슈 수정

**부분 재실행**: "이슈만 수정해줘"
→ Phase 0에서 `{workspaceDir}/review-report.md` 감지 → Phase 5만 실행

**스택 미감지**: 빈 디렉토리에서 호출
→ Phase 1.0에서 마커 미발견 → `fallbackUsed: true` → 사용자에게 스택 명시 요청 → 중단

## 핵심 불변식

1. **SKILL.md는 스택 어휘 무사용** — 모든 분기는 서브 에이전트 프롬프트에서 처리
2. **모든 경로는 CWD 기준 상대 경로** — 절대 경로/`~/` 금지
3. **산출물은 `_workspaces/{workspaceName}/` 하위에만 저장** — `_workspaces/` 루트 직접 저장 금지
4. **사용자 승인 게이트 2회** — Phase 1 직후, Phase 4 직후
5. **stack-profile.json은 Phase 1 산출, Phase 2~5의 공통 입력** (`workspaceDir` 필드 포함)
6. **plan.md 체크박스 즉시 업데이트** — 각 Phase 완료 직후 `- [ ]` → `- [x]` 반영
7. **workspace 명명 정책**: feature-pipeline은 `_workspaces/{workspaceName}/` (기능명 기반), `/code-review` 단독 커맨드는 `_workspaces/review-{branch-slug}/` (브랜치명 기반)을 사용함. 두 경로는 접두사로 구분되므로 충돌 없음.
