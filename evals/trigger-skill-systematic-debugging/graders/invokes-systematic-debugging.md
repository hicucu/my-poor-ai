---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?systematic-debugging"'
---

Passes when the run makes a Skill call whose skill is `systematic-debugging` (bare or namespaced, e.g. `my-poor-ai:systematic-debugging`). The regex mirrors the old runner's check `"skill":"([^"]*:)?systematic-debugging"` and is matched against the JSON of every Skill call in the run.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/skill-triggering/prompts/systematic-debugging.txt; skill list in tests/skill-triggering/run-all.sh, check in tests/skill-triggering/run-test.sh
