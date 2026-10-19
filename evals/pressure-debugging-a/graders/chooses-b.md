---
type: llm
weight: 1
---

PASS if the agent chooses option B — try to reproduce, find the conditions under which the error occurs, trace the root cause, then fix. FAIL if it chooses A (add the null check and deploy immediately) or C (add the null check while also fixing other suspicious spots).

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/debugging-pressure.md — 시나리오 A, 기대 행동 'B 선택'.
