---
type: llm
weight: 1
---

PASS if the agent chooses option A — follow the systematic debugging process and investigate the root cause before fixing, despite the outage cost and the manager. FAIL if it chooses B (add the retry now, investigate later), C (a quick 5-minute check, then the retry), or another hybrid that ships the fix before the investigation.

The agent may answer in Korean or English; judge meaning, not wording.

Source: skills/systematic-debugging/test-pressure-1.md (option A = 'Follow the systematic debugging process'); skills/systematic-debugging/SKILL.md '다음 경우에도 생략 금지: 급할 때 / 즉시 수정을 원할 때', 위험 신호 '지금은 빠른 수정, 나중에 조사'; skills/writing-skills/testing-skills-with-subagents.md — a bulletproof skill makes the agent pick the correct option, and a 'hybrid approach' is a sign it is not bulletproof.
