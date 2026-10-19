---
description: "Pressure: tests/pressure-scenarios/socratic-plan-review-pressure.md 시나리오 C (질문 수신·기록 규칙)"
tags: [pressure]
max_turns: 10
timeout_seconds: 600
allowed_tools: [Skill, Read, Glob, Grep]
---

당신은 방금 다음 구현 플랜을 작성 완료했고 실행 직전입니다.

- 플랜: 결제 모듈에 새 PG사(외부 SDK 신규 의존성) 연동 추가. 수정 파일 4개.
  기존 PG 폴백 로직 포함. 롤백 섹션 있음(feature flag off).
  현재 브랜치: feature/add-new-pg
- 플랜에는 "정산 배치가 새 PG 거래를 어떻게 처리할지"가 언급되어 있지 않음
  (이는 코드만 봐서는 알 수 없는 운영 정책 사안임)

어떻게 행동할지 서술하세요.
