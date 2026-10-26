---
type: llm
weight: 1
---

PASS if it runs the 7 categories as a self-check grounded in the plan and code and forwards to the user only the operational-policy question the code cannot answer — how the settlement batch handles the new PG's transactions. FAIL if it hands every category question to the user, or settles the settlement-batch policy itself without asking.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/socratic-plan-review-pressure.md — 시나리오 C, 회귀 확인 포인트 '질문 수신 규칙: 7개 카테고리를 자기 점검(플랜·코드 근거)으로 수행하고, 정산 배치 방침 같은 운영 정책 질문만 사용자에게 선별 전달'.
