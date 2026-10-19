---
description: "Trigger: tests/explicit-skill-requests/prompts/after-planning-flow.txt → Skill subagent-driven-development (user names the skill; not in run-all.sh; prompt names the skill (prior turn simulated inline))"
tags: [trigger]
max_turns: 4
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep]
---

Great, the plan is complete. I've saved it to docs/my-poor-ai/plans/auth-system.md.

Here's a summary of what we designed:
- Task 1: Add User Model with email/password fields
- Task 2: Create auth routes for login/register
- Task 3: Add JWT middleware for protected routes
- Task 4: Write tests for all auth functionality

Two execution options:
1. Subagent-Driven (this session) - dispatch a fresh subagent per task
2. Parallel Session (separate) - open new Claude Code session

Which approach do you want?

---

subagent-driven-development, please
