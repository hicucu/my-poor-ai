---
description: "Trigger: tests/explicit-skill-requests/prompts/claude-suggested-it.txt → Skill subagent-driven-development (user names the skill; not in run-all.sh; prompt names the skill (prior turn simulated inline))"
tags: [trigger]
max_turns: 4
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep]
---

[Previous assistant message]:
Plan complete and saved to docs/my-poor-ai/plans/auth-system.md.

Two execution options:
1. Subagent-Driven (this session) - I dispatch a fresh subagent per task, review between tasks, fast iteration within this conversation
2. Parallel Session (separate) - Open a new Claude Code session with the execute-plan skill, batch execution with review checkpoints

Which approach do you want to use for implementation?

[Your response]:
subagent-driven-development, please
