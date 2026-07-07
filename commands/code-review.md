---
description: 변경 브랜치를 4개 전문 에이전트(Architecture/Security/Performance/Style)가 병렬로 검토 후 통합 리포트를 작성하는 단독 커맨드. feature-pipeline의 Phase 4만 단독 실행. "코드 리뷰", "변경 검토", "리뷰해줘", "review", "PR 리뷰", "이 변경사항 봐줘" 같은 요청에 사용. 자연어로도 트리거 가능.
model: opus
---

# Code Review (병렬 4-전문 리뷰 단독 커맨드)

변경 브랜치를 4개 전문 에이전트가 병렬로 검토하고 aggregator가 통합 리포트를 작성함.
feature-pipeline의 Phase 4와 동일한 자산을 사용하여 단독 실행함.

## 사용법

```
/my-poor-ai:code-review                              현재 브랜치를 자동 감지된 base 브랜치와 비교
/my-poor-ai:code-review <base-branch>                base 브랜치 명시 (예: main, master, develop)
/my-poor-ai:code-review <base>...<head>              명시적 범위 지정
/my-poor-ai:code-review --files <file>...            특정 파일들만 리뷰 (diff 대신)
/my-poor-ai:code-review --out <경로>                 출력 경로 지정 (기본: _workspaces/review-{branch-slug}/review-report.md)
```

자연어 트리거 예시:

- "이 변경사항 코드 리뷰해줘"
- "PR 올리기 전에 리뷰 한번"
- "현재 브랜치 검토"
- "보안·성능 같이 봐줘"

## 실행 절차

### Step 1: 변경 범위 파악

base 브랜치 자동 감지 우선순위: `main` → `master` → `develop`. 사용자가 인수로 명시하면 그것을 우선.

```bash
git log <base>..HEAD --oneline          # 커밋 목록
git diff <base>...HEAD --stat           # 변경 파일 요약
git diff <base>...HEAD                  # 전체 diff (에이전트에 전달)
```

소스 코드 변경이 없으면 "리뷰할 코드 변경사항이 없음" 출력 후 종료.

`--files` 옵션 사용 시: 지정된 파일들의 전체 내용을 입력으로 사용 (diff 대신).

### Step 2: 작업 디렉토리 준비

현재 브랜치명에서 kebab-case 슬러그를 생성하여 `{workspaceDir} = _workspaces/review-{branch-slug}/` 결정.
예: 브랜치 `feature/login-api` → `_workspaces/review-feature-login-api/`

> **경로 규칙**: 모든 산출물은 `{workspaceDir}/` 하위에 저장함. `_workspaces/` 루트 직접 저장 금지.

기존 동일 디렉토리가 있으면 `_workspaces/review-{branch-slug}_prev/`로 백업.

stack-profile.json이 없으면 `/my-poor-ai:detect-stack`을 먼저 실행 권장 (또는 인라인 추론으로 진행).

### Step 3: 4개 전문 reviewer 병렬 실행

**필수: 4개 Agent 도구 호출을 단일 응답에 동시 포함**.

각 에이전트에게 위 git diff 내용과 stack-profile.json 경로를 전달.
diff에 신규 테이블·엔티티·스키마 생성이 포함된 경우, security-reviewer에게 마이그레이션 패턴 확인을 명시적으로 지시함.

```
{팀_위치}/agents/architecture-reviewer.md를 읽고 그 지침에 따라 작업한다.
리뷰 대상: [diff 또는 파일 목록]
스택 프로필: {workspaceDir}/stack-profile.json (있으면)
출력 경로: {workspaceDir}/reviews/architecture.md

{팀_위치}/agents/security-reviewer.md를 읽고 그 지침에 따라 작업한다.
(동일 입력)
출력 경로: {workspaceDir}/reviews/security.md
추가 지시: diff에 신규 DB 테이블/엔티티 생성이 있으면 마이그레이션 파일(Flyway, Liquibase, EF Migrations, Alembic 등) 존재 여부 확인.
          마이그레이션 파일 미발견 시 High 이슈로 보고.

{팀_위치}/agents/performance-reviewer.md를 읽고 그 지침에 따라 작업한다.
(동일 입력)
출력 경로: {workspaceDir}/reviews/performance.md

{팀_위치}/agents/style-reviewer.md를 읽고 그 지침에 따라 작업한다.
(동일 입력)
출력 경로: {workspaceDir}/reviews/style.md
```

### Step 4: aggregator로 통합

4개 산출이 모두 도착하면 aggregator 호출:

```
{팀_위치}/agents/review-aggregator.md를 읽고 그 지침에 따라 작업한다.

입력:
  {workspaceDir}/reviews/architecture.md
  {workspaceDir}/reviews/security.md
  {workspaceDir}/reviews/performance.md
  {workspaceDir}/reviews/style.md
스택 프로필: {workspaceDir}/stack-profile.json (있으면)
출력: {workspaceDir}/review-report.md  (또는 --out 지정 경로)
```

### Step 5: 결과 출력

콘솔에 다음 형식으로 요약 출력:

```
코드 리뷰 완료
─────────────────────────────
브랜치    : <head>
기준      : <base> (N commits, M files changed)
스택      : <감지 결과 요약>

진행 현황
  [x] Architecture 리뷰
  [x] Security 리뷰
  [x] Performance 리뷰
  [x] Style 리뷰
  [x] 통합 리포트

발견 이슈
  Architecture : 0 Critical / 1 High / 2 Medium / 0 Low
  Security     : 1 Critical / 0 High / 1 Medium / 0 Low
  Performance  : 1 Critical / 1 High / 0 Medium / 0 Low
  Style        :   -        / 2 High / 3 Medium / 1 Low

산출 파일 : {workspaceDir}/review-report.md
다음 단계 : Phase 5 (이슈 수정) 진행 시 `/feature-pipeline 이슈 수정해줘` 또는 issue-fixer 직접 호출
```

## 인수 처리

| 인수                                 | 동작                                           |
| ------------------------------------ | ---------------------------------------------- |
| 없음                                 | 자동 감지된 base와 HEAD 비교                   |
| 위치 인수 1개 (`main`, `develop` 등) | base 브랜치 명시                               |
| 위치 인수 1개 (`base...head` 형식)   | 명시적 범위                                    |
| `--files <file1> <file2>...`         | diff 대신 파일 전체 검토                       |
| `--out <경로>`                       | 출력 경로 지정                                 |
| `--skip-aggregate`                   | 4개 카테고리 리뷰만 작성, 통합 생략 (디버그용) |

## 활용 예시

### 예시 1: PR 올리기 전 리뷰

```
/my-poor-ai:code-review main
```

→ main과 비교하여 Critical/High 이슈 즉시 파악. PR 디스크립션에 요약 첨부 가능.

### 예시 2: 특정 파일만 검토

```
/my-poor-ai:code-review --files src/api/UserController.ts src/services/userService.ts
```

→ diff 없이 두 파일 전체를 4-전문 관점으로 검토.

### 예시 3: feature-pipeline 일부로 사용

`/feature-pipeline` 실행 중 Phase 4가 자동으로 동일한 4-전문 + aggregator를 호출. 단독 `/my-poor-ai:code-review`는 동일 동작을 파이프라인 외에서 실행.

## 에러 핸들링

| 상황                         | 대응                                                |
| ---------------------------- | --------------------------------------------------- |
| base 브랜치 자동 감지 실패   | 사용자에게 base 명시 요청                           |
| diff 0건                     | "리뷰할 변경사항이 없음" 출력 후 종료           |
| 4개 reviewer 중 1개 실패     | 나머지로 진행, aggregator가 "검토 실패" 표기        |
| 브랜치명 슬러그 생성 불가    | `_workspaces/review-manual/` 사용 후 사용자에게 알림 |
| stack-profile.json 없음      | 각 reviewer가 diff에서 스택 추론 (없어도 동작)      |
| 출력 경로 부모 디렉토리 없음 | 자동 생성                                           |

## 절대 금지

- 코드 파일 직접 수정 (리뷰만 작성)
- plan.md, file-manifest.json 작성 (feature-pipeline 전체 흐름의 산출물)
- 사용자에게 파일 목록 승인 요청 (이 커맨드는 게이트 없음)
- `_workspaces/` 루트에 직접 저장 — 반드시 `_workspaces/review-{branch-slug}/` 하위에만 저장
- 절대 경로/`~/` 사용

## 참조

- 전문 에이전트 정의: `agents/{architecture,security,performance,style}-reviewer.md`
- 통합 로직: `agents/review-aggregator.md`
- 파이프라인 흐름: `skills/feature-pipeline/SKILL.md` Phase 4
