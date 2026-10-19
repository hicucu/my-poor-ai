---
type: llm
weight: 1
---

PASS if the answer's instructions for every one of the three lens agents include a WebSearch directive. FAIL if any lens is left to training knowledge.

The agent may answer in Korean or English; judge meaning, not wording. The Agent and WebSearch tools are not available to the agent in this sandbox, so an explicit plan to dispatch or to search counts as the behavior.

Source: tests/pressure-scenarios/which-way-should-i-go-pressure.md — 시나리오 A, 기대 행동 '3개 에이전트 전부에 WebSearch 지시 포함'.
