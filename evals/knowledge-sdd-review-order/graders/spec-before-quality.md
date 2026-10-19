---
type: llm
weight: 1
---

PASS if the answer states that spec compliance review comes first and code quality review comes after it (code quality review starts only once spec compliance passes). FAIL if the order is reversed or left ambiguous.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/claude-code/test-subagent-driven-development.sh — Test 2, assertion 'Spec compliance before code quality'.
