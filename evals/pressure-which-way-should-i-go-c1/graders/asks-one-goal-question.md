---
type: llm
weight: 1
---

PASS if the agent takes the Socratic branch: it asks one question at a time aimed at pinning down the goal (a single goal-focused question in this reply), without digging into the full requirements and without launching a generic three-lens direction survey. FAIL if it starts a generic survey with no goal, or fires off a list of requirement questions.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/which-way-should-i-go-pressure.md — 시나리오 C, 기대 행동 'C-1: 산파술 분기 — 한 번에 하나씩 질문, 목표가 한 줄로 잡히는 즉시 중단 후 3렌즈 진입 (요구사항 전체를 캐지 않음)'.
