---
name: executing-plans
description: Use when executing a written implementation plan in a separate session with review checkpoints
---

# 플랜 실행

## 개요

플랜을 불러와서 비판적으로 검토하고, 모든 작업을 실행한 뒤 완료 시 보고.

**시작 시 공지:** "executing-plans 스킬을 사용하여 이 플랜을 구현함."

**참고:** subagent 접근 권한이 있을 때 my-poor-ai가 훨씬 더 잘 작동한다는 점을 사용자에게 알릴 것. subagent를 지원하는 플랫폼(Claude Code 또는 Codex 등)에서 실행하면 작업 품질이 현저히 높아짐. subagent를 사용할 수 있다면 이 스킬 대신 my-poor-ai:subagent-driven-development를 사용할 것.

## 프로세스

### 1단계: 플랜 불러오기 및 검토

1. 플랜 파일 읽기
2. 비판적으로 검토 — 플랜에 대한 질문이나 우려 사항 파악
3. 우려 사항이 있을 경우: 시작 전 사용자에게 제기
4. 우려 사항 없을 경우: TodoWrite 생성 후 진행

### 2단계: 작업 실행

각 작업에 대해:

1. in_progress로 표시
2. 각 단계를 정확히 따름 (플랜은 작은 단위의 단계로 구성됨)
3. 명시된 검증 수행
4. completed로 표시

**HANDOFF.md 갱신 (복잡 경로 전용):**

`_workspaces/{branch-slug}/HANDOFF.md`는 세션 인계 전용 문서임. 갱신 시점:

- 작업(task) 완료 → TodoWrite 체크만. HANDOFF 건드리지 않음.
- phase 완료(스펙에 Phase 그룹이 있을 때) → 본문 갱신 + 인계 로그 1줄.
- spec 완료 → 본문 갱신 + 인계 로그 1줄, 해당 마일스톤 커밋에 `HANDOFF.md` 포함.

HANDOFF.md가 없으면 첫 갱신 시 표준 템플릿(my-poor-ai:writing-plans 참조)으로 먼저 생성. 본문은 덮어쓰기, 인계 로그는 최근 5개만 유지. 단순·디버깅 경로에는 만들지 않음.

### 3단계: 개발 완료

모든 작업 완료 및 검증 후:

- 공지: "finishing-a-development-branch 스킬을 사용하여 이 작업을 완료함."
- **필수 서브 스킬:** my-poor-ai:finishing-a-development-branch 사용
- 해당 스킬에 따라 테스트 검증, 옵션 제시, 선택 실행

## 도움 요청 시점

**다음 상황에서 즉시 실행 중단:**

- 블로커 발생 (누락된 의존성, 테스트 실패, 지시 사항 불명확)
- 플랜에 시작을 막는 심각한 공백 존재
- 지시 사항 이해 불가
- 검증이 반복적으로 실패

**추측하지 말고 명확한 설명을 요청할 것.**

## 이전 단계로 돌아가는 경우

**다음 상황에서 검토(1단계)로 복귀:**

- 피드백에 따라 사용자가 플랜을 업데이트한 경우
- 근본적인 접근 방식 재검토가 필요한 경우

**블로커를 강제로 통과하지 말 것** — 중단하고 질문할 것.

## 주의 사항

- 먼저 플랜을 비판적으로 검토
- 플랜 단계를 정확히 따를 것
- 검증 단계 건너뛰기 금지
- 플랜에서 지시할 때 해당 스킬 참조
- 블로커 발생 시 중단, 추측 금지
- 사용자의 명시적 동의 없이 main/master 브랜치에서 구현 시작 금지

## 통합

**필수 워크플로우 스킬:**

- **my-poor-ai:using-git-worktrees** — 격리된 작업 공간 확보 (생성 또는 기존 환경 확인)
- **my-poor-ai:writing-plans** — 이 스킬이 실행할 플랜 생성
- **my-poor-ai:finishing-a-development-branch** — 모든 작업 완료 후 개발 마무리
