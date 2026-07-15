# Skills

[English](README.md) | **한국어**

에이전트가 반드시 따라야 할 프로세스 규칙 20개. 요청 내용에 따라 자동 트리거되거나 명시적으로 호출됨. 각 스킬 `SKILL.md`의 frontmatter `description`이 실제 트리거 조건이며, 이 문서는 개발 단계별로 묶은 빠른 색인임. 수정 전에는 `skills/writing-skills/` 참조 필수.

## 설계·계획

| 스킬 | 역할 |
| --- | --- |
| [`brainstorming`](brainstorming/) | 창작 작업(기능 추가·기능 수정) 전 필수 — 구현 전 대화로 의도·요구사항·설계를 탐색 |
| [`which-way-should-i-go`](which-way-should-i-go/) | 무엇을 브레인스토밍할지조차 정하기 어려울 때, 검증된 접근법 2~3가지를 병렬 웹 조사로 비교 |
| [`writing-plans`](writing-plans/) | 컨텍스트 없는 엔지니어를 가정하고 한 입 크기 작업으로 나눈 종합 구현 플랜 작성 |
| [`socratic-plan-review`](socratic-plan-review/) | 실행 전 7개 카테고리의 구조화된 소크라테스식 질문으로 복잡한 플랜 검증 |

## 구현

| 스킬 | 역할 |
| --- | --- |
| [`test-driven-development`](test-driven-development/) | RED → GREEN → REFACTOR — 구현 전 실패하는 테스트부터 작성 |
| [`subagent-driven-development`](subagent-driven-development/) | 독립 태스크로 구성된 플랜을 태스크별 신규 서브에이전트로 실행, 매 태스크 후 2단계 검토 |
| [`executing-plans`](executing-plans/) | 별도 세션에서 작성된 플랜을 불러와 비판적으로 검토, 전체 실행 후 완료 보고 |
| [`using-git-worktrees`](using-git-worktrees/) | 네이티브 worktree 도구로 격리된 작업 공간 보장, 없으면 수동 git worktree로 대체 |
| [`dispatching-parallel-agents`](dispatching-parallel-agents/) | 독립적인 태스크 2개 이상을 정밀하게 구성한 지시·컨텍스트로 격리된 서브에이전트에 위임 |

## 완료·리뷰

| 스킬 | 역할 |
| --- | --- |
| [`verification-before-completion`](verification-before-completion/) | 완료·수정·통과를 주장하기 전, 검증 명령 실행과 결과 확인 필수 |
| [`finishing-a-development-branch`](finishing-a-development-branch/) | 구현·테스트 통과 후 merge/PR/cleanup 구조화된 선택지 제시 |
| [`requesting-code-review`](requesting-code-review/) | 리뷰 전용으로 구성된 컨텍스트를 갖춘 리뷰어 subagent를 일찍, 자주 파견 |
| [`receiving-code-review`](receiving-code-review/) | 리뷰 피드백을 맹목적 수용이 아닌 기술적 검증 대상으로 취급 |

## 디버깅

| 스킬 | 역할 |
| --- | --- |
| [`systematic-debugging`](systematic-debugging/) | 버그·테스트 실패·예상치 못한 동작을 임의 수정 대신 근본 원인 추적으로 해결 |

## 문서화·CI 안전

| 스킬 | 역할 |
| --- | --- |
| [`generate-claude-instructions`](generate-claude-instructions/) | `CLAUDE.md`와 참조 문서(`DEVELOPMENT.md`, `LANGUAGE_GUIDELINES.md`, `AI_BEHAVIOR.md`, `COMMIT_CONVENTION.md`)를 생성하는 오케스트레이터 |
| [`sync-docs-from-diff`](sync-docs-from-diff/) | 브랜치 diff를 분석해 README/docs/인라인 문서 갱신을 제안, 사용자 승인 후에만 적용 |
| [`preventing-github-actions-loops`](preventing-github-actions-loops/) | GitHub Actions 워크플로우의 자기 트리거 순환(무한루프) 탐지·방지 |

## 파이프라인

| 스킬 | 역할 |
| --- | --- |
| [`feature-pipeline`](feature-pipeline/) | 스택 무관 5단계 멀티에이전트 파이프라인: 계획 → 구현 → 테스트 → 리뷰 → 수정 |

## 메타

| 스킬 | 역할 |
| --- | --- |
| [`writing-skills`](writing-skills/) | 프로세스 문서에 적용한 TDD — 스킬 작성·수정·검증 방법 |
| [`using-my-poor-ai`](using-my-poor-ai/) | 모든 대화 시작 시 로딩 — 스킬을 찾고 호출하는 방법 확립 |
