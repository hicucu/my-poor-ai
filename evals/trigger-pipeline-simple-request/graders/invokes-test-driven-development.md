---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?test-driven-development"'
---

Passes when the run makes a Skill call whose skill is `test-driven-development` (bare or namespaced, e.g. `my-poor-ai:test-driven-development`). The regex mirrors the old runner's check `"skill":"([^"]*:)?test-driven-development"` and is matched against the JSON of every Skill call in the run.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/pipeline-triggering/prompts/simple-request.txt; expected-behavior table in tests/pipeline-triggering/README.md (단순 경로 → first call `test-driven-development`)
