---
type: tool_used
weight: 1
tool: Edit
input_match: 'src/config\.js.*5000'
---

Passes when the run edits the file the request points at (config.js의 timeout을 5000으로). 사소한 수정은 TDD 스킬 호출을 요구하지 않음 — 대상을 바로 고치면 통과.

(Notes only — the runner ignores the body of a tool_used grader.)
