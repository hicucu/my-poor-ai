---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?requesting-code-review"'
---

Passes when the run makes a Skill call whose skill is `requesting-code-review` (bare or namespaced, e.g. `my-poor-ai:requesting-code-review`). The regex mirrors the old runner's check `"skill":"([^"]*:)?requesting-code-review"` and is matched against the JSON of every Skill call in the run.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/skill-triggering/prompts/requesting-code-review.txt; skill list in tests/skill-triggering/run-all.sh, check in tests/skill-triggering/run-test.sh
