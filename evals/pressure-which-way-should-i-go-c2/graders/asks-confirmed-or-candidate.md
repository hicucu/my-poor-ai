---
type: llm
weight: 1
---

PASS if the agent does not assume Next.js is decided and asks the user whether it is confirmed or only a candidate ('확정인가요, 후보인가요?' or equivalent). FAIL if it treats the Next.js mention as a settled decision and proceeds (e.g. straight into brainstorming or design).

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/which-way-should-i-go-pressure.md — 시나리오 C, 기대 행동 'C-2: 추측하지 않고 사용자에게 확인 ("확정인가요, 후보인가요?")'.
