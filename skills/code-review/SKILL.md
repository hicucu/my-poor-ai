---
name: code-review
description: 변경 브랜치를 아키텍처·성능 관점으로 검토하고 네이티브 /security-review 결과와 통합하는 커맨드. 네이티브 /code-review가 덮지 않는 축만 담당함. "코드 리뷰", "변경 검토", "리뷰해줘", "아키텍처 봐줘", "성능 검토", "PR 리뷰" 같은 요청에 사용.
model: opus
---

# Code Review (아키텍처·성능 축 + 네이티브 통합)

변경 브랜치를 **아키텍처·성능** 관점으로 검토하고, 네이티브 `/security-review` 결과와 함께 통합 리포트를 작성함.

## 축 분담

my-poor-ai는 네이티브 리뷰가 덮지 않는 축만 담당함. 같은 일을 두 번 하지 않음.

| 축 | 담당 | 호출 방법 |
| --- | --- | --- |
| 정확성(버그·회귀) + 재사용·단순화·효율 | 네이티브 `/code-review` | **사용자만 호출 가능** — 이 커맨드가 대신 실행할 수 없음 |
| 보안 취약점 | 네이티브 `/security-review` | 이 커맨드가 Skill 도구로 호출 |
| 아키텍처(레이어·의존성·SRP·결합도) | my-poor-ai `architecture-reviewer` | 이 커맨드가 스폰 |
| 성능(N+1·인덱스·블로킹·캐싱·렌더링) | my-poor-ai `performance-reviewer` | 이 커맨드가 스폰 |

> **`/code-review`를 대신 실행할 수 없는 이유**: 번들 `/code-review`와 `/verify`는 `disable-model-invocation`이라 사용자가 직접 입력할 때만 실행됨. 긴 실행·토큰 소비 시점의 통제권을 사용자에게 두기 위한 설계임. 따라서 이 커맨드는 정확성 축을 **직접 구현하지 않고**, 리포트 말미에 `/code-review` 실행을 안내함.

## 사용법

```
/my-poor-ai:code-review                      현재 브랜치를 자동 감지된 base와 비교
/my-poor-ai:code-review <base-branch>        base 브랜치 명시 (main, master, develop)
/my-poor-ai:code-review <base>...<head>      명시적 범위
/my-poor-ai:code-review --files <file>...    특정 파일들만 (diff 대신)
/my-poor-ai:code-review --out <경로>         출력 경로 지정
/my-poor-ai:code-review --no-security        네이티브 보안 리뷰 생략 (아키텍처·성능만)
/my-poor-ai:code-review --init-review-md     리뷰 리포트 대신 REVIEW.md 초안 생성 (아래 참조)
```

## 실행 절차

### Step 1: 변경 범위 파악

base 자동 감지: `main` → `master` → `develop`. 인수로 명시하면 그것을 우선.

```bash
git log <base>..HEAD --oneline
git diff <base>...HEAD --stat
git diff <base>...HEAD
```

소스 변경이 없으면 "리뷰할 코드 변경사항이 없음" 출력 후 종료.
`--files` 사용 시 지정 파일 전체 내용을 입력으로 사용.

### Step 2: 작업 디렉토리 준비

브랜치명 kebab-case 슬러그로 `{workspaceDir} = _workspaces/review-{branch-slug}/` 결정.
기존 동일 디렉토리는 `_workspaces/review-{branch-slug}_prev/`로 백업.

> 모든 산출물은 `{workspaceDir}/` 하위. `_workspaces/` 루트 직접 저장 금지.

### Step 3: 3개 리뷰를 단일 응답에서 동시 착수

**아키텍처·성능 2개 Agent 호출과 `/security-review` Skill 호출을 같은 응답에 포함함.**

```
{팀_위치}/agents/architecture-reviewer.md를 읽고 그 지침에 따라 작업한다.
리뷰 대상: [diff 또는 파일 목록]
스택 프로필: {workspaceDir}/stack-profile.json (있으면)
출력 경로: {workspaceDir}/reviews/architecture.md

{팀_위치}/agents/performance-reviewer.md를 읽고 그 지침에 따라 작업한다.
(동일 입력)
출력 경로: {workspaceDir}/reviews/performance.md
```

보안 축은 **에이전트를 스폰하지 않고** 네이티브 스킬을 호출함:

```
Skill 도구로 security-review 호출 → 결과를 {workspaceDir}/reviews/security.md에 저장
```

`--no-security` 지정 시 이 호출을 생략하고 리포트에 "보안 축 생략됨"을 명시함.

> **네이티브 보안 리뷰의 적용 범위 제약** (2026-07-28 실측): `/security-review`는 대상 경로 인수를 받지 않고 **현재 세션의 저장소와 그 자체 diff 범위**로만 동작함. 따라서 다음 경우에는 보안 축이 **적용 불가**이며, 조건을 억지로 맞추지 않고 리포트에 "보안 축 적용 불가(사유)"로 기록함:
>
> - 리뷰 대상이 세션 작업 디렉토리의 저장소가 아닌 경우 — 다른 저장소를 검토한 리포트는 미수행보다 위험한 거짓 안전 신호임
> - `--files`로 파일을 지정한 경우 — 네이티브는 diff 기준이라 파일 선택을 받지 않음
> - 원격이 없거나 `origin/HEAD`가 없어 네이티브가 diff 범위를 계산하지 못하는 경우
>
> 이 경우 사용자에게 "대상 저장소에서 직접 `/security-review` 실행"을 안내함.

> diff에 신규 DB 테이블·엔티티 생성이 있으면 마이그레이션 파일(Flyway, Liquibase, EF Migrations, Alembic 등) 존재 여부를 **아키텍처 리뷰어에게** 확인 지시함. 스키마 변경의 배포 가능성은 아키텍처 축임.

### Step 4: aggregator로 통합

```
{팀_위치}/agents/review-aggregator.md를 읽고 그 지침에 따라 작업한다.

입력:
  {workspaceDir}/reviews/architecture.md
  {workspaceDir}/reviews/performance.md
  {workspaceDir}/reviews/security.md   (--no-security면 생략)
스택 프로필: {workspaceDir}/stack-profile.json (있으면)
출력: {workspaceDir}/review-report.md  (또는 --out 지정 경로)
```

### Step 5: 결과 출력

```
코드 리뷰 완료 (아키텍처·성능 + 보안)
─────────────────────────────
브랜치    : <head>
기준      : <base> (N commits, M files changed)

발견 이슈
  Architecture : 0 Critical / 1 High / 2 Medium / 0 Low
  Performance  : 1 Critical / 1 High / 0 Medium / 0 Low
  Security     : 1 Critical / 0 High / 1 Medium / 0 Low   (네이티브 /security-review)

산출 파일 : {workspaceDir}/review-report.md

⚠ 정확성 축은 아직 검토되지 않음
  버그·회귀·엣지케이스와 재사용·단순화 정리는 네이티브 리뷰가 담당함.
  이어서 직접 실행할 것:  /code-review
  (이 커맨드는 /code-review를 대신 실행할 수 없음 — 사용자 전용 스킬)
```

## `--init-review-md` 모드

GitHub Code Review(관리형 PR 리뷰)를 쓰는 저장소라면, 리뷰 관점을 리포지토리 루트 `REVIEW.md`에 두는 것이 정본임. `REVIEW.md`는 리뷰 파이프라인 **모든 에이전트의 시스템 프롬프트에 최우선 블록으로 주입**되므로, 같은 규칙을 `CLAUDE.md`에 두는 것보다 훨씬 강하게 반영됨.

이 모드는 리뷰를 실행하지 않고 저장소에 맞춘 `REVIEW.md` 초안만 작성함.

절차:

1. 스택·디렉토리 구조·기존 `CLAUDE.md`를 읽어 저장소 성격 파악
2. 아래 골격을 저장소 실정에 맞게 채움
3. `REVIEW.md`가 이미 있으면 덮어쓰지 않고 차이만 제안

골격:

```markdown
# Review instructions

## What Important means here
<이 저장소에서 🔴 Important로 올릴 결함 부류를 명시 — 기본 보정은 프로덕션 코드 기준임>

## Cap the nits
<한 리뷰당 **게시할** 🟡 Nit 상한과, 초과분 처리 방법. 상한은 게시량 제한이며 탐지 억제가 아님 — 초과분은 찾지 않는 것이 아니라 묶어서 요약하거나 생략함>

## Do not report
<CI가 이미 강제하는 것(lint·format·타입), 생성 코드 경로, lockfile, 벤더 디렉토리. 여기 적는 것은 **결함 부류**여야 함 — "확신이 낮은 것", "경미한 것" 같은 확신도·중요도 기준을 적지 말 것>

## Always check
<이 저장소에서 매 PR 확인할 규칙 — 예: 신규 API 라우트에 통합 테스트 존재>

## Verification bar
<근거 요구 수준 — 예: 동작 주장은 file:line 인용 필요, 명명에서 추론 금지>
```

작성 원칙: **길수록 희석됨.** 리뷰 동작을 바꾸는 지시만 남기고, 일반 프로젝트 맥락은 `CLAUDE.md`에 둠. `REVIEW.md`는 `@` import가 전개되지 않으므로 규칙을 파일 안에 직접 씀.

> **보고 정책과 탐지 억제를 구분할 것.** 리뷰어는 "경미한 것은 올리지 말 것", "확신이 낮으면 생략할 것" 류의 지시를 문자 그대로 따르며, 그 결과 **찾는 것 자체를 줄임** — 정밀도는 오르지만 실제 결함이 조용히 누락됨. `Cap the nits`·`Do not report`는 *무엇을 게시할지*를 정하는 것이고, *무엇을 찾을지*를 정하는 것이 아님. 탐지는 전부, 게시는 상한선 — 이 순서를 REVIEW.md 문면에서 지킬 것.

## 인수 처리

| 인수 | 동작 |
| --- | --- |
| 없음 | 자동 감지된 base와 HEAD 비교 |
| 위치 인수 1개 (`main` 등) | base 브랜치 명시 |
| 위치 인수 1개 (`base...head`) | 명시적 범위 |
| `--files <file>...` | diff 대신 파일 전체 검토 |
| `--out <경로>` | 출력 경로 지정 |
| `--no-security` | 네이티브 보안 리뷰 생략 |
| `--init-review-md` | 리뷰 미실행, `REVIEW.md` 초안만 작성 |
| `--skip-aggregate` | 카테고리별 리뷰만, 통합 생략 (디버그용) |

## 에러 핸들링

| 상황 | 대응 |
| --- | --- |
| base 자동 감지 실패 | 사용자에게 base 명시 요청 |
| diff 0건 | "리뷰할 변경사항이 없음" 출력 후 종료 |
| `/security-review` 호출 실패·미가용 | 아키텍처·성능으로 진행, 리포트에 "보안 축 미수행" 명시 |
| 리뷰 대상이 세션 cwd 저장소가 아님 | 보안 축 "적용 불가"로 기록 + 대상 저장소에서 직접 실행 안내 |
| 육안으로 보안 결함이 보임 (축은 미수행) | "완전성이 보장되지 않는 부수 관찰"로 **분리 표기**. 보안 축 수행으로 승격 금지 |
| 리뷰어 1개 실패 | 나머지로 진행, aggregator가 "검토 실패" 표기 |
| 브랜치명 슬러그 생성 불가 | `_workspaces/review-manual/` 사용 후 알림 |
| stack-profile.json 없음 | 각 리뷰어가 diff에서 스택 추론 |

## 절대 금지

- **보안 리뷰를 자체 에이전트로 재구현** — 네이티브 `/security-review`를 호출함. 호출이 실패했다면 실패를 보고하지, 직접 구현으로 대체하지 않음
- **정확성·버그 축을 이 커맨드에서 검토했다고 보고** — 그 축은 `/code-review` 담당이며 미수행 상태를 반드시 명시함
- 코드 파일 직접 수정 (리뷰만 작성)
- 사용자에게 파일 목록 승인 요청 (이 커맨드는 게이트 없음)
- `_workspaces/` 루트 직접 저장, 절대 경로·`~/` 사용

## 참조

- 전문 에이전트: `agents/architecture-reviewer.md`, `agents/performance-reviewer.md`
- 통합 로직: `agents/review-aggregator.md`
- 네이티브 리뷰 축 분담: 이 문서 상단 "축 분담" 표

> **파이프라인과의 차이**: `feature-pipeline` Phase 4와 `review-agent`는 아직 보안·스타일을 자체 에이전트로 수행함. 이 커맨드만 먼저 네이티브 위임으로 전환된 상태이며, 파이프라인 정렬은 별도 과제임.
