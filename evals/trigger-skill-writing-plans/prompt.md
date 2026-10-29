---
description: "Trigger: tests/skill-triggering/prompts/writing-plans.txt → Agent project-context (multi-step feature is FULL; planning-agent writes the specs)"
tags: [trigger, trigger-fixture]
max_turns: 4
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep, Edit, Write]
---

Here's the spec for our new authentication system:

Requirements:
- Users can register with email/password
- Users can log in and receive a JWT token
- Protected routes require valid JWT
- Tokens expire after 24 hours
- Support password reset via email

We need to implement this. There are multiple steps involved - user model, auth routes, middleware, email service integration.
