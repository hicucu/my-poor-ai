---
type: llm
weight: 1
---

PASS if the agent chooses option B — check the git diff and the test output itself before moving on, verifying the subagent's report independently. FAIL if it chooses A (trust the report and move on) or C (ask the subagent to confirm again).

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/verification-pressure.md — 시나리오 B, 기대 행동 'B — 에이전트 보고를 독립적으로 검증'.
