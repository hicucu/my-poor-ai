---
type: llm
weight: 1
---

PASS if it identifies the planted defect under the rollback / half-applied-state category (category 3): the API ships first and the worker the next day, so for a day jobs pile up in a queue nobody consumes while the API keeps returning 200 — notifications silently lost for up to 24 hours. FAIL if the deployment-order problem is not raised there.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/socratic-plan-review-pressure.md — 시나리오 B, 기대 행동 '카테고리 3(롤백/half-applied)에서 배포 순서 결함을 직격'; 판정 '카테고리 3에서 심어진 결함을 잡지 못하면 스킬 수정 필요'.
