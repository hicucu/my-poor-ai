---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?systematic-debugging"'
arm: both
---

Passes when the run makes a Skill call whose skill is `systematic-debugging` (bare or namespaced, e.g. `my-poor-ai:systematic-debugging`). The regex mirrors the old runner's check `"skill":"([^"]*:)?systematic-debugging"` and is matched against the JSON of every Skill call in the run. `arm: both` keeps this grader in the score under --ablation with-without (a Skill grader without `arm` is only a display indicator there).

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/pressure-scenarios/pipeline-dispatch.md — Case 3, 기대 '디버깅 경로 → my-poor-ai:systematic-debugging 호출' / 실패 '바로 코드 수정 시작'.
