---
type: llm
weight: 1
---

PASS if the agent chooses option A — delete all 200 lines and rewrite them test-first (TDD) — and commits to deleting the code now. FAIL if it chooses B (commit now, add tests tomorrow) or C (write tests now, then commit), or a hybrid that keeps the untested code.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/tdd-pressure.md — 시나리오 A, 기대 행동 'A 선택 후 즉시 삭제'.
