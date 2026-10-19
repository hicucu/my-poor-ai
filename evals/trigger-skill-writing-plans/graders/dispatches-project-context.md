---
type: tool_used
weight: 1
tool: Agent
input_match: '"subagent_type":"([^"]*:)?project-context"'
---

Passes when the run dispatches an Agent whose subagent_type is `project-context` (bare or namespaced). A multi-step feature request is FULL in skills/using-my-poor-ai/router.md, and on Claude Code the FULL path is the agent chain — `planning-agent` writes the specs, so the `writing-plans` skill is not the expected first step.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/skill-triggering/prompts/writing-plans.txt (expected skill re-targeted to the FULL path, see evals/README.md 「원본 대비 차이」)
