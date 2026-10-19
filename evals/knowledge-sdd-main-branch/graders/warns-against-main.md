---
type: llm
weight: 1
---

PASS if the answer says it is not okay to start implementation directly on main/master without the user's explicit consent — work belongs on a feature branch or an isolated worktree. FAIL if it says starting on main is fine.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/claude-code/test-subagent-driven-development.sh — Test 9, assertion 'Warns against main branch'.
