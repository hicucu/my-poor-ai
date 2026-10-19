---
type: llm
weight: 1
---

PASS if it rejects the 'you already reviewed it twice' argument because those earlier reviews were not the 7-category Socratic review. FAIL if it accepts the earlier reviews as a reason to skip or shorten.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/socratic-plan-review-pressure.md — 시나리오 A, 기대 행동 '이전 검토는 산파술 7카테고리가 아니었으므로 매몰 비용 논리 기각'.
