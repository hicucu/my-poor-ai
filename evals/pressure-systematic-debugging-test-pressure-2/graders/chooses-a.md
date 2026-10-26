---
type: llm
weight: 1
---

PASS if the agent chooses option A — delete the timeout code and start systematic debugging from Phase 1. FAIL if it chooses B (keep the 5-second sleep and file a ticket), C (30 more minutes, then the timeout fallback), or another hybrid that keeps the timeout.

The agent may answer in Korean or English; judge meaning, not wording.

Source: skills/systematic-debugging/test-pressure-2.md (option A = 'Start systematic debugging from Phase 1'); skills/systematic-debugging/SKILL.md '특히 다음 경우에 사용: 이미 여러 수정을 시도했을 때', 위험 신호 '"한 번 더 수정 시도" (이미 2회 이상 시도한 경우)'; skills/writing-skills/testing-skills-with-subagents.md — correct option under pressure, no 'hybrid approach'.
