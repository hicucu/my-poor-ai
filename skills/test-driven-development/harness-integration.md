# harness/파이프라인 통합 — test-driven-development

**이 참조를 불러올 때:** TDD가 단독 대화가 아니라 my-poor-ai 파이프라인(spec → 구현)의 한 단계로 실행될 때.

## 위치

TDD는 my-poor-ai spec-driven 워크플로우의 **구현 단계 규율**임.

```
brainstorming → writing-plans(spec) → subagent-driven-development(구현) → 코드 리뷰
                                              └─ 각 작업에서 test-driven-development 적용
```

`writing-plans`가 만든 spec을 받아, implementer가 각 작업을 RED-GREEN-REFACTOR로 구현함.

## 입력 계약

implementer는 다음을 입력으로 받음:

- **spec**: 구현할 작업 명세 (`writing-plans` 산출, `_workspaces/{branch-slug}/specs/`)
- **`_workspaces/stack-profile.json`**: 테스트 실행 명령·프레임워크·파일 명명 규칙을 결정. 복잡 경로에서는 Phase 0 `project-context`가 생성함 (없으면 프로젝트 매니페스트에서 판단). 본문(SKILL.md)은 스택 어휘를 쓰지 않고, 구체 명령은 이 프로필에서 가져옴.

## 실행 계약 (능력)

- implementer는 **테스트를 실제 실행할 수 있어야** 함. RED 검증(실패 직접 확인)과 GREEN 검증(통과 확인)은 명령 실행이 전제임.
- 도구가 제한되어 테스트를 실행할 수 없는 작업자에게는 TDD를 위임하지 않음 — 실행 없는 RED/GREEN은 검증이 아님.

## 순서 계약 (모순 해결)

- implementer는 **테스트 선행**으로 동작함: RED(실패 테스트) → GREEN(최소 구현) → REFACTOR.
- "구현 먼저, 테스트 나중"은 **금지**함 — 테스트 후행은 안티패턴이며, 즉시 통과하는 테스트는 아무것도 증명하지 못함.
- 이 계약은 `subagent-driven-development`의 implementer 프롬프트와 일치해야 함. 프롬프트가 "구현 → 테스트" 순서이면 TDD가 깨지므로, 프롬프트를 테스트 선행으로 맞춤.

## 네임스페이스

- 파이프라인·문서·프롬프트에서 이 skill 참조는 항상 **`my-poor-ai:test-driven-development`** 로 함 (구 네임스페이스 형태 금지).

## 에이전트 체인에서의 적용

복잡 경로의 에이전트 체인(`planning-agent` → `developer-agent`)도 같은 두 계약을 따름:

- **실행 계약** — `developer-agent`는 Bash를 보유해 테스트를 실제로 실행함.
- **순서 계약** — `planning-agent`의 스펙은 태스크를 "실패하는 테스트 작성 → 실패 확인 → 최소 구현 → 통과 확인" 순으로 둠.
