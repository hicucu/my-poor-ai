---
type: llm
weight: 1
---

PASS if the answer says review is a loop, not a one-time pass: after fixes the reviewer reviews again, repeating until approved/compliant. FAIL if it describes a single review.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/claude-code/test-subagent-driven-development.sh — Test 6, assertion 'Review loops mentioned'.
