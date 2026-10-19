---
type: llm
weight: 1
---

PASS if it treats the 4 modified files as meeting the 3-or-more-files threshold and runs socratic-plan-review. FAIL if it decides the plan is below the threshold or skips the review.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/socratic-plan-review-pressure.md — 시나리오 C, 회귀 확인 포인트 '발동 임계값: 수정 파일 4개가 description("3+ files")과 본문("3개 이상") 양쪽에서 일관되게 걸리는가'.
