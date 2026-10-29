---
description: "Trigger: tests/skill-triggering/prompts/requesting-code-review.txt → Skill requesting-code-review (naive prompt, skill not named)"
tags: [trigger, trigger-fixture]
max_turns: 4
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep, Edit, Write]
---

I just finished implementing the user authentication feature. All the code is committed. Can you review the changes before I merge to main?

The commits are between abc123 and def456.
