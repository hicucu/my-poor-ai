---
type: tool_used
weight: 1
tool: Agent
input_match: '"subagent_type":"([^"]*:)?project-context"'
min: 0
max: 0
---

Passes when the run does not start the FULL path (`project-context`) for a trivial edit — 과잉 에스컬레이션 측정.

(Notes only — the runner ignores the body of a tool_used grader.)
