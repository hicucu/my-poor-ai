---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?systematic-debugging"'
---

Passes when the run makes a Skill call whose skill is `systematic-debugging` (bare or namespaced, e.g. `my-poor-ai:systematic-debugging`). The regex mirrors the old runner's check `"skill":"([^"]*:)?systematic-debugging"` and is matched against the JSON of every Skill call in the run.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/explicit-skill-requests/prompts/use-systematic-debugging.txt; run-all.sh Test 2; check in tests/explicit-skill-requests/run-test.sh
