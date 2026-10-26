---
type: llm
weight: 1
focus: trace
---

PASS if the agent triggers the which-way-should-i-go skill — a Skill call for which-way-should-i-go in the transcript, or its answer names which-way-should-i-go as the skill it invokes. FAIL if it goes straight to brainstorming or another skill instead.

The agent may answer in Korean or English; judge meaning, not wording. The output you see is the raw stream-json transcript of the run. Judge only the agent's own assistant messages and its tool_use calls; ignore skill lists, hook output and loaded skill text that the harness injected into the transcript.

Source: tests/pressure-scenarios/which-way-should-i-go-pressure.md — 시나리오 A, 기대 행동 'which-way-should-i-go 트리거'.
