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
- **stack-profile.json**: 테스트 실행 명령·프레임워크·파일 명명 규칙을 결정. 본문(SKILL.md)은 스택 어휘를 쓰지 않고, 구체 명령은 이 프로필에서 가져옴.

## 실행 계약 (능력)

- implementer는 **테스트를 실제 실행할 수 있어야** 함. RED 검증(실패 직접 확인)과 GREEN 검증(통과 확인)은 명령 실행이 전제임.
- 도구가 제한되어 테스트를 실행할 수 없는 작업자에게는 TDD를 위임하지 않음 — 실행 없는 RED/GREEN은 검증이 아님.

## 순서 계약 (모순 해결)

- implementer는 **테스트 선행**으로 동작함: RED(실패 테스트) → GREEN(최소 구현) → REFACTOR.
- "구현 먼저, 테스트 나중"은 **금지**함 — 테스트 후행은 안티패턴이며, 즉시 통과하는 테스트는 아무것도 증명하지 못함.
- 이 계약은 `subagent-driven-development`의 implementer 프롬프트와 일치해야 함. 프롬프트가 "구현 → 테스트" 순서이면 TDD가 깨지므로, 프롬프트를 테스트 선행으로 맞춤.

## 네임스페이스

- 파이프라인·문서·프롬프트에서 이 skill 참조는 항상 **`my-poor-ai:test-driven-development`** 로 함 (구 네임스페이스 형태 금지).

## harness 대체 메모 (권고, 침습 변경 아님)

`example-project/harness`의 `feature-pipeline`은 "TDD"를 표방하지만 실제로는 TDD 정합이 아님:

- Phase 2 `file-developer`는 도구에 테스트 실행 수단이 없어 RED/GREEN 검증이 불가능함 (실행 계약 위반).
- 실제 테스트는 Phase 3 `test-writer`가 **구현 이후**에 작성함 (순서 계약 위반 = 테스트 후행).

→ 이 구조는 my-poor-ai spec-driven 파이프라인으로 **대체**하면 해소됨. harness 자체(에이전트 도구·Phase 구성)는 이 작업에서 침습적으로 바꾸지 않음.
