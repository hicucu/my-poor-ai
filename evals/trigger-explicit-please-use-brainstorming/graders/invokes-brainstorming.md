---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?brainstorming"'
---

Passes when the run makes a Skill call whose skill is `brainstorming` (bare or namespaced, e.g. `my-poor-ai:brainstorming`). The regex mirrors the old runner's check `"skill":"([^"]*:)?brainstorming"` and is matched against the JSON of every Skill call in the run.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/explicit-skill-requests/prompts/please-use-brainstorming.txt; run-all.sh Test 3; check in tests/explicit-skill-requests/run-test.sh
