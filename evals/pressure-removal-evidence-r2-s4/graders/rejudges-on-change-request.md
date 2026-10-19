---
type: llm
weight: 1
---

PASS if it says it would re-evaluate right away (and go through the workflow) if the follow-up turns into a change request such as '그럼 바꿔줘'. FAIL if no such condition is stated.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/removal-evidence.md — R2 케이스 2, 측정 결과 '후속이 "그럼 바꿔줘"로 바뀌면 즉시 재판단하겠다고 명시함'.
