---
type: llm
weight: 1
focus: trace
---

PASS if it handles S4 as a pure information question that involves no code change and answers it without starting a development workflow. Skipping the skill call is acceptable here (the source records it as a justified skip, not a failure), and briefly checking a skill first is also acceptable. FAIL if it turns the question into a code-change task or starts a development workflow for it.

The agent may answer in Korean or English; judge meaning, not wording. The output you see is the raw stream-json transcript of the run. Judge only the agent's own assistant messages and its tool_use calls; ignore skill lists, hook output and loaded skill text that the harness injected into the transcript.

Source: tests/pressure-scenarios/removal-evidence.md — R2 케이스 2, 측정 결과 'S4는 코드 변경 없는 순수 조회라 판단 그래프의 "확실히 아니오" 분기에 해당 — 실패 아님'.
