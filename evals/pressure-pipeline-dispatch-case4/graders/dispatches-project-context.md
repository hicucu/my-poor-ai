---
type: tool_used
weight: 1
tool: Agent
input_match: '"subagent_type":"([^"]*:)?project-context"'
---

Passes when the run dispatches an Agent whose subagent_type is `project-context` (bare or namespaced, e.g. `my-poor-ai:project-context`) — FULL path Phase 0 in skills/using-my-poor-ai/router.md. Agent graders are scored in both arms under --ablation with-without.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/pressure-scenarios/pipeline-dispatch.md — Case 4, 기대 '복잡 경로' (router FULL Phase 0 = project-context) / 실패 '단순 경로 (설계 없이 바로 구현)'.
