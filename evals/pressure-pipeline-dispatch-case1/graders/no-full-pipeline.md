---
type: tool_used
weight: 1
tool: Agent
input_match: '"subagent_type":"([^"]*:)?(project-context|brainstorming-agent)"'
min: 0
max: 0
arm: both
---

Passes when the run does not start the FULL path (`project-context` or `brainstorming-agent`) for a trivial edit. 사소한 수정 예외(router.md)에 따라 TDD 호출은 요구하지 않음. `arm: both` keeps this grader in the score under --ablation with-without.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/pressure-scenarios/pipeline-dispatch.md — Case 1, 실패 'brainstorming-agent 호출 (오버엔지니어링)'.
