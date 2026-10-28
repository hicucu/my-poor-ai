---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill"\s*:\s*"'
min: 0
max: 0
---

Passes when the run makes no Skill call at all. The request is a pure information question with no code change, bug, or feature work, so routing through a my-poor-ai skill is an over-trigger. Measures the cost side of the routing gate; the positive `trigger-*` cases measure the miss side.

(Notes only — the runner ignores the body of a tool_used grader.)
