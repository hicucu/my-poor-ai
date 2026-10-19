---
type: tool_used
weight: 1
tool: Agent
input_match: '"subagent_type":"([^"]*:)?project-context"'
---

Passes when the run dispatches an Agent whose subagent_type is `project-context` (bare or namespaced). The prompt asks for "a new feature", and skills/using-my-poor-ai/router.md sends any new feature to FULL immediately, whose Phase 0 is `project-context`.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/skill-triggering/prompts/test-driven-development.txt (expected skill re-targeted to the FULL path, see evals/README.md 「원본 대비 차이」)
