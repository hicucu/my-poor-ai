---
type: llm
weight: 1
focus: trace
---

PASS if the agent proceeds on the ordinary feature-addition path — the FULL path (in the router: the project-context agent as Phase 0, then brainstorming-agent) or brainstorming. FAIL if it runs a direction survey, or skips the feature workflow and implements directly.

The agent may answer in Korean or English; judge meaning, not wording. The output you see is the raw stream-json transcript of the run. Judge only the agent's own assistant messages and its tool_use calls; ignore skill lists, hook output and loaded skill text that the harness injected into the transcript.

Source: tests/pressure-scenarios/which-way-should-i-go-pressure.md — 시나리오 E, 기대 행동 (스킬 유무 무관) '일반 기능 추가 경로(FULL 경로 또는 brainstorming)로 진행'.
