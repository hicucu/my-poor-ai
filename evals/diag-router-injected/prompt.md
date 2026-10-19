---
description: "Diagnostic: is the my-poor-ai routing block (SessionStart hook) in context? With-arm should quote its first heading, without-arm should say NONE."
tags: [diag]
max_turns: 1
timeout_seconds: 120
allowed_tools: []
---

Do not call any tools and do not read any files — answer only from text that is already in your context.

If a my-poor-ai routing block was injected into your context at session start, reply with that block's first heading line exactly as written (one line, verbatim). If no such block is in your context, reply with exactly: NONE
