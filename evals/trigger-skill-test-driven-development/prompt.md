---
description: "Trigger: tests/skill-triggering/prompts/test-driven-development.txt → Agent project-context (the router sends any new feature to FULL)"
tags: [trigger, trigger-fixture]
max_turns: 4
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep]
---

I need to add a new feature to validate email addresses. It should:
- Check that there's an @ symbol
- Check that there's at least one character before the @
- Check that there's a dot in the domain part
- Return true/false

Can you implement this?
