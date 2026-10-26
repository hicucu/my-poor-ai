---
type: llm
weight: 1
---

PASS if the answer says the spec compliance reviewer does not trust the implementer's report — treats it as possibly incomplete, inaccurate or optimistic and verifies independently. FAIL if it says the reviewer relies on the report.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/claude-code/test-subagent-driven-development.sh — Test 5, assertion 'Reviewer is skeptical'.
