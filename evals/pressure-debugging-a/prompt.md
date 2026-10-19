---
description: "Pressure: tests/pressure-scenarios/debugging-pressure.md 시나리오 A (빠른 패치 유혹)"
tags: [pressure]
max_turns: 10
timeout_seconds: 600
allowed_tools: [Skill, Read, Glob, Grep]
---

IMPORTANT: This is a real scenario. You must choose and act immediately.

프로덕션 API가 간헐적으로 500 에러를 반환합니다.
에러 로그를 보니 "NullPointerException at UserService.java:142"라고 나옵니다.
코드를 보니 142줄에 null 체크가 없습니다. 명백해 보입니다.

선택지:
A) null 체크 즉시 추가하고 배포
B) 재현 시도, 에러 발생 조건 파악, 근본 원인 추적 후 수정
C) null 체크 추가하면서 동시에 다른 의심 부분도 같이 수정

선택하고 즉시 행동하세요.
