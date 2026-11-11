# 실행 모드 — 서브에이전트 / 에이전트팀

복잡 경로 Phase 3(개발)·Phase 4(리뷰)와 `subagent-driven-development`의 실행 모드 정본임. 두 스킬은 모드 판정 요약만 두고 상세는 이 문서를 따름.

> 근거: Claude Code Agent Teams 공식 문서(https://code.claude.com/docs/en/agent-teams), 2026-11-11 확인. 실험 기능이므로 도구 이름·제약이 바뀔 수 있음 — 다르면 공식 문서를 우선함.

## 모드 판정

| 조건 | 모드 |
| --- | --- |
| `printenv CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` 결과가 `1`이 아님 | 서브에이전트 — 질문하지 않음 |
| 사용자에게 질문할 수단이 없음 (`-p`·SDK 등 비대화형 세션) | 서브에이전트 — 이 환경에서는 팀원이 스폰되지 않음 |
| 위 둘 다 아님 (팀 기능 활성 + 대화형) | **사용자에게 질문** — 서브에이전트 / 에이전트팀 중 선택 |
| `pipeline-state.md`에 `mode`가 이미 기록됨 | 그 값을 재사용 — 다시 묻지 않음 |

- 질문 시점: 복잡 경로는 **Phase 3 진입 직전 1회** (Phase 4까지 같은 모드), `subagent-driven-development`는 첫 작업 배포 전 1회
- 질문에 담을 내용: 에이전트팀은 팀원 간 직접 메시지로 검토 왕복이 빠른 대신 토큰 비용이 크고, 세션 재개(`/resume`) 시 팀원이 복원되지 않음
- 선택 결과는 `pipeline-state.md`의 `mode:` 줄에 기록함 (`subagent` 또는 `team`)

WHY: 팀 기능은 기본 비활성 실험 기능이고 비용·재개 특성이 서브에이전트와 다름. 비활성 환경에서 묻는 것은 의미 없는 질문이고, 활성 환경에서 임의로 고르는 것은 비용 결정을 사용자 대신 내리는 것임.

## 공통 규칙 (팀 모드)

- **팀 리드는 메인 세션뿐임.** 서브에이전트·팀원은 팀을 만들거나 이끌 수 없고(중첩 팀 불가), 세션당 팀은 1개임
- 팀원 스폰: `Agent` 호출에 `name`을 지정함. 역할 정의는 기존 에이전트를 그대로 씀 — `Agent(subagent_type="my-poor-ai:developer-agent", name="dev-spec-a", ...)`
- 팀원에게는 에이전트 정의의 `skills` frontmatter가 적용되지 않음 — 필요한 지침은 스폰 프롬프트에 직접 담음
- 통신: `SendMessage`(이름 지정, 1명씩). 진행 추적: 공유 태스크 도구(`TaskCreate`·`TaskList`·`TaskUpdate`)가 있으면 사용, 없으면 메시지로 조율
- 동시 팀원은 **최대 5명** (공식 권장 3~5명) — 초과분은 다음 배치로 미룸
- 같은 파일을 두 팀원이 동시에 편집하지 않음 — 덮어쓰기 발생
- Phase(또는 SDD 실행) 종료 시 팀원 전원에게 종료를 요청함
- 세션 재개 시 팀원은 복원되지 않음 — `pipeline-state.md`의 미완료 작업 기준으로 팀원을 새로 스폰함

## 복잡 경로 Phase 3 — 개발 (팀 모드)

1. `file-manifest.json`의 `developmentOrder` 그룹 순서를 그대로 따름 — 그룹 간 의존성 규칙은 서브에이전트 모드와 동일
2. 그룹 내 spec마다 `developer-agent` 팀원 1명 (`name="dev-{spec 이름}"`). 프롬프트는 서브에이전트 모드의 developer-agent 입력과 동일함
3. 같은 파일을 건드리는 spec이 한 그룹에 있으면 **같은 팀원에게 순차 배정**함
4. 팀원의 `STATUS: NEEDS_CONTEXT`·질문은 리드가 사용자에게 중계하고 `SendMessage`로 같은 팀원에게 답함
5. 그룹의 모든 팀원이 DONE이면 종료 요청 후 다음 그룹 진행

## 복잡 경로 Phase 4 — 리뷰 (팀 모드)

`review-agent`는 서브에이전트라 팀을 만들 수 없으므로, **메인 세션이 리드로서 리뷰를 직접 조율함.** 리뷰 축·산출 파일·aggregator·issue-fixer 계약은 `agents/_shared/review-orchestration.md`를 그대로 따르고, 차이는 아래뿐임:

| 항목 | 서브에이전트 모드 | 팀 모드 |
| --- | --- | --- |
| 조율 주체 | `review-agent` | 메인 세션 (리드) |
| 리뷰어 3종 | 서브에이전트 병렬 | 팀원 (`name="rev-architecture"` 등). 보안 축은 리드가 `/security-review` 호출 (폴백 조건이면 `security-reviewer` 팀원) |
| 교차 검증 | 없음 | 4축 산출 후, 각 리뷰어가 다른 축 산출 파일을 읽고 반박·중복을 해당 리뷰어에게 `SendMessage` — 1라운드만, 결과는 각자 산출 파일에 반영 |
| aggregator·issue-fixer | review-agent가 호출 | 리드가 호출 (issue-fixer 파일 단위 병렬·중복 금지 규칙 동일) |

## 흔한 실수

| 실수 | 결과 | 대안 |
| --- | --- | --- |
| 팀 모드에서 `review-agent`에 리뷰 위임 | 서브에이전트는 팀을 못 만들어 서브에이전트 모드로 조용히 퇴행 | 리드가 직접 조율 |
| 팀 기능 비활성인데 모드 질문 | 고를 수 없는 선택지 제시 | 묻지 않고 서브에이전트 |
| 활성 환경에서 묻지 않고 팀 모드 선택 | 사용자 동의 없는 비용 증가 | Phase 3 직전 1회 질문 |
| 재개 시 이전 팀원에게 `SendMessage` | 팀원이 존재하지 않음 | 미완료 작업 기준 재스폰 |
