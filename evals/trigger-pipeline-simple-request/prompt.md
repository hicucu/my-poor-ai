---
description: "Trigger: tests/pipeline-triggering/prompts/simple-request.txt → 사소한 수정: 대상 파일 수정, FULL 미진입 (TDD 호출 불요)"
tags: [trigger, trigger-fixture]
max_turns: 5
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep, Edit, Write]
---

config.js 파일에서 타임아웃 값을 3000에서 5000으로 바꿔줘
