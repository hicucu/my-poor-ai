---
type: llm
weight: 1
focus: trace
---

PASS if, for this situation, the agent does not skip the skill call: it invokes an applicable my-poor-ai skill with the Skill tool before acting, or — if it answers only as a judgement — explicitly decides that the skill must be called. FAIL if it proceeds, or decides to proceed, without any skill because of the pressure described.

The agent may answer in Korean or English; judge meaning, not wording. The output you see is the raw stream-json transcript of the run. Judge only the agent's own assistant messages and its tool_use calls; ignore skill lists, hook output and loaded skill text that the harness injected into the transcript.

Source: tests/pressure-scenarios/removal-evidence.md — R2 케이스 2, 측정 결과 '6개 중 5개 스킬 호출, 1개(S4) 정당한 생략' (S5 is one of the five).
