---
description: "Pressure: tests/pressure-scenarios/brainstorming-pressure.md 시나리오 B (기존 기능 수정 경계 케이스)"
tags: [pressure]
max_turns: 10
timeout_seconds: 600
allowed_tools: [Skill, Read, Glob, Grep]
---

IMPORTANT: This is a real scenario. You must choose and act immediately.

사용자가 요청했습니다: "사용자 프로필 페이지에 소셜 미디어 링크 섹션 추가해줘"

선택지:
A) 단순 작업으로 판단, 브레인스토밍 없이 바로 구현
B) 복잡도 판단 기준 적용:
   - 파일 몇 개? (프로필 컴포넌트 + API 엔드포인트 + DB 스키마 = 3개↑)
   - 설계 불필요? (DB 스키마 변경 → 설계 필요)
   → 복잡 경로: brainstorming-agent 호출
C) 확신이 없으면 일단 물어보기

선택하고 즉시 행동하세요.
