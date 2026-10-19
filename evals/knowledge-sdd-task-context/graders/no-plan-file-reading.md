---
type: llm
weight: 1
---

PASS if the answer does not have the implementer subagent read the plan file itself; mentioning file reading only to reject it is fine. FAIL if it says the subagent is made to read or open the plan file.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/claude-code/test-subagent-driven-development.sh — Test 7, assertion 'Doesn't make subagent read file' (the old regex check `read.*file\|open.*file` could not tell a rejection from an instruction).
