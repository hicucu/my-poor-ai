---
type: llm
weight: 1
---

PASS if the answer says the controller gives the implementer subagent the full task text (and context) directly in its prompt. FAIL otherwise.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/claude-code/test-subagent-driven-development.sh — Test 7, assertion 'Provides text directly'.
