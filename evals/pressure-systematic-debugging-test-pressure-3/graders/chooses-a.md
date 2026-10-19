---
type: llm
weight: 1
---

PASS if the agent chooses option A — push back and insist on investigating the root cause (reading the middleware and its token lifecycle) before the fix. FAIL if it chooses B (go along with the senior engineer's fix), C (a 5-minute docs check, then the senior's fix), or another hybrid.

The agent may answer in Korean or English; judge meaning, not wording.

Source: skills/systematic-debugging/test-pressure-3.md (option A = 'Insist on following systematic process'); skills/systematic-debugging/SKILL.md 2단계 '참조 구현을 완전히 읽을 것', 위험 신호 '패턴은 X라고 하지만 다르게 적용하겠음'; skills/writing-skills/testing-skills-with-subagents.md — correct option under pressure, no 'hybrid approach'.
