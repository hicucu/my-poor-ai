---
name: my-poor-ai
description: my-poor-ai 요청 처리 진입점 — SessionStart 훅이 동작하지 않는 환경에서도 사용 가능
argument-hint: [요청 내용]
allowed-tools: [Read, Edit, Write, Bash, Glob, Grep]
---

# /my-poor-ai — my-poor-ai 시스템 진입점

my-poor-ai 컨텍스트를 수동으로 활성화하고 요청을 처리함.
SessionStart 훅이 동작하지 않는 환경(플러그인 훅 비활성·타 플랫폼)에서도 동작함.

## 실행 방법

**Step 1 — using-my-poor-ai 스킬 호출**

`my-poor-ai:using-my-poor-ai` 스킬을 즉시 호출해 my-poor-ai 컨텍스트와 요청 분류 파이프라인을 활성화함.

**Step 2 — 요청 처리**

`$ARGUMENTS`가 있으면 해당 내용을 my-poor-ai 파이프라인으로 처리함.
인자가 없으면 사용자에게 무엇을 도와드릴지 물음.

**Step 3 — 훅 미동작 감지 (선택)**

SessionStart 훅이 동작하지 않는 것 같으면 (매 세션마다 이 커맨드를 직접 실행해야 하는 경우),
다음 안내를 덧붙임:

> Claude Code에서는 플러그인이 활성화되어 있으면 SessionStart 훅이 자동 적용됨 — `/plugin`에서 my-poor-ai 활성화 상태를 확인할 것.
