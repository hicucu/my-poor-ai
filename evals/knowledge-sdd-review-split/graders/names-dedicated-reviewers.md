---
type: llm
weight: 1
---

PASS if the answer names the dedicated reviewers as the ones who do the review: a spec compliance reviewer and a code quality reviewer (if it states an order, spec compliance must come first). FAIL if it does not name both dedicated reviewers.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/claude-code/test-subagent-driven-development.sh — Test 3, assertion 'Dedicated reviewers named'.
