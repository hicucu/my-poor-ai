---
type: llm
weight: 1
---

PASS if the answer says the controller reads the plan file once (a single time), not once per task. FAIL otherwise.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/claude-code/test-subagent-driven-development.sh — Test 4, assertion 'Read plan once'.
