---
type: llm
weight: 1
---

PASS if it rejects the one-line rollback ('문제 시 이전 버전 재배포') because redeploying code cannot revert a schema/data migration (half-applied risk), and asks for a staged rollback path. FAIL if it lets the rollback section pass.

The agent may answer in Korean or English; judge meaning, not wording.

Source: tests/pressure-scenarios/socratic-plan-review-pressure.md — 시나리오 A, 기대 행동 '롤백 한 줄 반려 — 스키마 마이그레이션에서 코드 재배포는 데이터를 되돌리지 못함(half-applied 위험) 지적, 단계별 롤백 경로 보강 요구'.
