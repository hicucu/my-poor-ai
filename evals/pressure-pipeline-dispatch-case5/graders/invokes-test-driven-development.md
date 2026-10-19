---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?test-driven-development"'
arm: both
---

Passes when the run makes a Skill call whose skill is `test-driven-development` (bare or namespaced, e.g. `my-poor-ai:test-driven-development`). The regex mirrors the old runner's check `"skill":"([^"]*:)?test-driven-development"` and is matched against the JSON of every Skill call in the run. `arm: both` keeps this grader in the score under --ablation with-without (a Skill grader without `arm` is only a display indicator there).

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/pressure-scenarios/pipeline-dispatch.md — Case 5, 기대 '단순 경로' (router SIMPLE → test-driven-development) / 실패 'brainstorming-agent 호출 (오버엔지니어링)'.
