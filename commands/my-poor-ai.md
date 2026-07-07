---
description: my-poor-ai 요청 처리 진입점 — setup 없이도 사용 가능
argument-hint: [요청 내용]
allowed-tools: [Read, Edit, Write, Bash, Glob, Grep]
---

# /my-poor-ai — my-poor-ai 시스템 진입점

my-poor-ai 컨텍스트를 수동으로 활성화하고 요청을 처리함.
SessionStart 훅이 등록되지 않은 환경에서도 동작함.

## 실행 방법

**Step 1 — using-my-poor-ai 스킬 호출**

`my-poor-ai:using-my-poor-ai` 스킬을 즉시 호출해 my-poor-ai 컨텍스트와 요청 분류 파이프라인을 활성화함.

**Step 2 — 요청 처리**

`$ARGUMENTS`가 있으면 해당 내용을 my-poor-ai 파이프라인으로 처리함.
인자가 없으면 사용자에게 무엇을 도와드릴지 물음.

**Step 3 — setup 미완료 감지 (선택)**

SessionStart 훅이 등록되지 않은 것 같으면 (매 세션마다 이 커맨드를 직접 실행해야 하는 경우),
다음 안내를 덧붙임:

> 💡 `/my-poor-ai:setup` 을 실행하면 매 세션 시작 시 my-poor-ai가 자동 활성화됨.
