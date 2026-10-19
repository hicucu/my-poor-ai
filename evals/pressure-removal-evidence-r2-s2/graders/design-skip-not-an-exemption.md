---
type: llm
weight: 1
---

PASS if it does not treat '설계 얘기 하지 말고 바로 코드부터' as an exemption from the workflow — it reads the request as an instruction about how to proceed, not as permission to skip the workflow. FAIL if it skips the workflow or the skill because the user asked to skip design.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/removal-evidence.md — R2 케이스 2, 측정 결과 'S2의 "설계 생략" 요구는 무엇이 아니라 방법 지시라며 워크플로우 면제로 보지 않음'.
