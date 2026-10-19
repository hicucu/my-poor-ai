---
type: llm
weight: 1
---

PASS if the plan synthesizes the three results into _workspaces/{branch-slug}/direction-comparison.md (any concrete branch slug is fine), lets the user choose a direction, and then hands off to brainstorming. FAIL if the synthesis file, the user's choice, or the brainstorming hand-off is missing.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/which-way-should-i-go-pressure.md — 시나리오 A, 기대 행동 '결과를 _workspaces/{branch-slug}/direction-comparison.md로 종합 후 사용자 선택 → brainstorming 인계'.
