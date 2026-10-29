---
description: "Trigger: tests/skill-triggering/prompts/dispatching-parallel-agents.txt → Skill systematic-debugging (failing tests are DEBUG first)"
tags: [trigger, trigger-fixture]
max_turns: 4
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep, Edit, Write]
---

I have 4 independent test failures happening in different modules:

1. tests/auth/login.test.ts - "should redirect after login" is failing
2. tests/api/users.test.ts - "should return user list" returns 500
3. tests/components/Button.test.tsx - snapshot mismatch
4. tests/utils/date.test.ts - timezone handling broken

These are unrelated issues in different parts of the codebase. Can you investigate all of them?
