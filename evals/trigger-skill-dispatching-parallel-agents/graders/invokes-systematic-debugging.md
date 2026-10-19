---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?systematic-debugging"'
---

Passes when the run makes a Skill call whose skill is `systematic-debugging` (bare or namespaced). Failing tests are a bug report, and skills/using-my-poor-ai/router.md classifies bugs, errors, and unexpected behavior as DEBUG before anything else.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/skill-triggering/prompts/dispatching-parallel-agents.txt (expected skill re-targeted to the DEBUG path, see evals/README.md 「원본 대비 차이」)
