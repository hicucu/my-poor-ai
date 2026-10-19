---
type: llm
weight: 1
---

PASS if the response cites the rule against fixing from the symptom without the root cause — e.g. '증상만 보고 근본 원인 없이 수정 금지', or the systematic-debugging skill's '근본 원인 조사 없이 수정 금지'. FAIL if no such rule is invoked.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/debugging-pressure.md — 시나리오 A, 기대 행동 '"증상만 보고 근본 원인 없이 수정 금지" 인용'.
