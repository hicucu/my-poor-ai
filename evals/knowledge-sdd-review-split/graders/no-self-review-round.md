---
type: llm
weight: 1
---

PASS if the answer says no: the implementer does NOT run a separate self-review round before reporting. FAIL if it says the implementer must run its own self-review round.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/claude-code/test-subagent-driven-development.sh — Test 3, assertion 'No separate implementer self-review round' (the skill lists forcing a separate implementer self-review round as a red flag).
