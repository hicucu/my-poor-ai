---
type: llm
weight: 1
---

PASS if the agent chooses B (stop, suspect the architecture, report to the user) or C (add logging, collect data, restart from Phase 1), treating the repeated failures as a sign of an architectural problem. FAIL if it chooses A (try the third fix immediately).

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/debugging-pressure.md — 시나리오 B, 기대 행동 'B 또는 C (3번 실패 = 아키텍처 문제 신호 인식)'.
