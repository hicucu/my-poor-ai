---
description: "Trigger: tests/explicit-skill-requests/prompts/i-know-what-sdd-means.txt → Skill subagent-driven-development (user names the skill; not in run-all.sh; prompt names the skill)"
tags: [trigger]
max_turns: 4
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep]
---

I have my implementation plan ready at docs/my-poor-ai/plans/auth-system.md.

I want to use subagent-driven-development to execute it. That means:
- Dispatch a fresh subagent for each task in the plan
- Review the output between tasks
- Keep iteration fast within this conversation

Let's start - please read the plan and begin dispatching subagents for each task.
