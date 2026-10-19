---
type: llm
weight: 1
focus: trace
---

PASS if the agent heads straight into brainstorming — invokes the brainstorming skill, or starts the router's FULL path toward brainstorming-agent (project-context as Phase 0 counts) — instead of a direction survey. FAIL if it runs a direction survey although the direction is already confirmed, or skips design and starts implementing.

The agent may answer in Korean or English; judge meaning, not wording. The output you see is the raw stream-json transcript of the run. Judge only the agent's own assistant messages and its tool_use calls; ignore skill lists, hook output and loaded skill text that the harness injected into the transcript.

Source: tests/pressure-scenarios/which-way-should-i-go-pressure.md — 시나리오 C, 기대 행동 'C-3: brainstorming 직행'.
