---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?which-way-should-i-go"'
min: 0
max: 0
arm: both
---

Passes when the run makes no Skill call whose skill is `which-way-should-i-go` (bare or namespaced, e.g. `my-poor-ai:which-way-should-i-go`). The regex mirrors the old runner's check `"skill":"([^"]*:)?which-way-should-i-go"` and is matched against the JSON of every Skill call in the run. `arm: both` keeps this grader in the score under --ablation with-without (a Skill grader without `arm` is only a display indicator there).

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/pressure-scenarios/which-way-should-i-go-pressure.md — 시나리오 C, 기대 행동 'C-3: brainstorming 직행' / baseline 실패 'C-3: 확정 상태인데도 which-way-should-i-go를 불필요하게 실행'.
