---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?subagent-driven-development"'
---

Passes when the run makes a Skill call whose skill is `subagent-driven-development` (bare or namespaced, e.g. `my-poor-ai:subagent-driven-development`). The regex mirrors the old runner's check `"skill":"([^"]*:)?subagent-driven-development"` and is matched against the JSON of every Skill call in the run.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/explicit-skill-requests/prompts/skip-formalities.txt; not in run-all.sh; prompt names the skill; check in tests/explicit-skill-requests/run-test.sh
