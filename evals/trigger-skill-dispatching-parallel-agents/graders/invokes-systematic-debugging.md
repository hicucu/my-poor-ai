---
type: tool_used
weight: 1
tool: Skill
input_match: '"skill":"([^"]*:)?(systematic-debugging|dispatching-parallel-agents)"'
---

Passes when the run makes a Skill call whose skill is `systematic-debugging` (bare or namespaced). Failing tests are a bug report, and skills/using-my-poor-ai/router.md classifies bugs, errors, and unexpected behavior as DEBUG before anything else.

2026-10-29: `dispatching-parallel-agents`도 통과로 인정함 — 원본 기대 스킬이며, Opus 5.5는 픽스처 위에서 3/3 이를 호출함. 실패 테스트를 DEBUG로 보는 라우터 판정(`systematic-debugging`)도 계속 통과.

(Notes only — the runner ignores the body of a tool_used grader.)

Source: tests/skill-triggering/prompts/dispatching-parallel-agents.txt (expected skill re-targeted to the DEBUG path, see evals/README.md 「원본 대비 차이」)
