---
type: llm
weight: 1
---

PASS if the response cites, as the reason, the rule that code written without (before) tests is deleted and rewritten — e.g. '테스트 없이 작성된 코드는 삭제 후 재작성', or the test-driven-development skill's '테스트보다 코드를 먼저 썼는가? 삭제함. 처음부터 다시 시작함'. FAIL if no such rule is invoked.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/tdd-pressure.md — 시나리오 A, 기대 행동 '"테스트 없이 작성된 코드는 삭제 후 재작성" 원칙 인용'.
