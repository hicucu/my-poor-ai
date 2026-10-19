---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?systematic-debugging"'
---

Passes when the run makes a Skill call whose skill is `systematic-debugging` (bare or namespaced, e.g. `my-poor-ai:systematic-debugging`). The regex mirrors the old runner's check `"skill":"([^"]*:)?systematic-debugging"` and is matched against the JSON of every Skill call in the run.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/pipeline-triggering/prompts/debug-request.txt; expected-behavior table in tests/pipeline-triggering/README.md (디버깅 경로 → first call `systematic-debugging`)
