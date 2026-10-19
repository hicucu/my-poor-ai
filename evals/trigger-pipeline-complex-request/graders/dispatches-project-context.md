---
type: tool_used
weight: 1
tool: Agent
input_match: '"subagent_type":"([^"]*:)?project-context"'
---

Passes when the run dispatches an Agent whose subagent_type is `project-context` (bare or namespaced, e.g. `my-poor-ai:project-context`) — FULL path Phase 0 in skills/using-my-poor-ai/router.md. Agent graders are scored in both arms under --ablation with-without.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/pipeline-triggering/prompts/complex-request.txt; expected-behavior table in tests/pipeline-triggering/README.md (복잡 경로 → first call `project-context`)
