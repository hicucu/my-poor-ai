---
description: "Trigger: tests/pipeline-triggering/prompts/debug-request.txt → Skill systematic-debugging (디버깅 경로)"
tags: [trigger]
max_turns: 5
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep]
---

배포 후부터 API가 간헐적으로 503을 반환해. 로그에는 "connection pool exhausted"가 찍혀있어.
