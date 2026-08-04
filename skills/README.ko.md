# Skills

[English](README.md) | **한국어**

스킬 디렉토리 28개. 각 스킬 `SKILL.md`의 frontmatter `description`이 실제 트리거 조건이며, 이 문서는 개발 단계별로 묶은 빠른 색인임. 수정 전에는 `skills/writing-skills/` 참조 필수.

두 종류가 함께 있으며, 차이는 **누가 시작하는가**임:

- **프로세스 스킬** (아래 앞 섹션들) — 요청 내용이 맞으면 Claude가 자동 로드하며, `/my-poor-ai:{이름}`으로 직접 호출도 가능함
- **커맨드** (아래 [커맨드](#커맨드) 섹션) — 사용자가 `/my-poor-ai:{이름}`으로 호출함. 프로젝트 외부에 부작용이 있는 것은 `disable-model-invocation: true`로 Claude의 자동 발동을 차단함

둘은 같은 파일 형식임. Claude Code가 커스텀 커맨드를 스킬로 통합하여 `commands/deploy.md`와 `skills/deploy/SKILL.md`가 동일하게 `/deploy`를 만들며, 스킬 쪽이 참조 파일 번들·호출 제어·서브에이전트 실행을 추가로 지원함. 이 플러그인이 전부 `skills/` 하위로 두는 이유임.

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
| [`using-my-poor-ai`](using-my-poor-ai/) | 요청 분류와 파이프라인 계약. 옆의 `router.md`가 SessionStart 훅이 주입하는 압축 블록임 |

## 커맨드

`/my-poor-ai:{이름}` 형태로 호출함. `/my-poor-ai:commands`가 카탈로그 진입점임.

| 스킬 | 역할 | 모델 자동 발동 |
| --- | --- | --- |
| [`my-poor-ai`](my-poor-ai/) | 요청을 알맞은 파이프라인(DEBUG/SIMPLE/FULL)으로 라우팅하는 진입점 — setup 없이도 사용 가능 | 가능 |
| [`commands`](commands/) | 사용 가능한 목록 표시, 카탈로그 자체 | 가능 |
| [`roles`](roles/) | 역할 프리셋 카탈로그 — 역할명(architect/builder/debugger/reviewer/docs)을 스킬 번들로 라우팅 | 가능 |
| [`code-review`](code-review/) | 아키텍처·성능 축 검토 + 네이티브 `/security-review` 통합, `REVIEW.md` 초안 생성 | 가능 |
| [`detect-stack`](detect-stack/) | 마커 파일 스캔으로 기술 스택 감지, `stack-profile.json` 생성 | 가능 |
| [`git-resume`](git-resume/) | 자연어 시간 표현 또는 commit hash로 과거 작업 맥락 복원 | 가능 |
| [`session-manager`](session-manager/) | 로컬 Claude Code 세션 목록 조회·이름 변경·삭제 | 가능 |
| [`setup`](setup/) | `SessionStart` 훅을 `~/.claude/settings.json`에 등록 | **차단** — 프로젝트 외부 기록 |
| [`codex-setup`](codex-setup/) | 이 플러그인의 에이전트를 `~/.codex/config.toml`에 등록 | **차단** — 프로젝트 외부 기록 |
