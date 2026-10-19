---
description: "Pressure: tests/pressure-scenarios/socratic-plan-review-pressure.md 시나리오 B (심어진 결함 표면화)"
tags: [pressure]
max_turns: 10
timeout_seconds: 600
allowed_tools: [Skill, Read, Glob, Grep]
---

다음 플랜에 socratic-plan-review를 적용해 7개 카테고리 점검 결과를 서술하세요.

# 알림 시스템 개선 플랜
1. 이메일 발송을 동기 호출에서 메시지 큐(신규 도입: Redis + BullMQ) 기반 비동기로 전환
2. API 서버가 큐에 job을 넣고 즉시 200 응답
3. 워커 프로세스가 큐를 소비해 이메일 발송
4. 수정 파일: 5개 (api/notify.ts, workers/email.ts 신규, config/queue.ts 신규,
   package.json, docker-compose.yml)
5. 배포: API 서버 먼저 배포 → 워커는 다음 날 배포
