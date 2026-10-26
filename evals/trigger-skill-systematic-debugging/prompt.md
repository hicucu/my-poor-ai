---
description: "Trigger: tests/skill-triggering/prompts/systematic-debugging.txt → Skill systematic-debugging (naive prompt, skill not named)"
tags: [trigger]
max_turns: 4
timeout_seconds: 600
allowed_tools: [Skill, Agent, Read, Glob, Grep]
---

The tests are failing with this error:

```
FAIL src/utils/parser.test.ts
  ● Parser › should handle nested objects
    TypeError: Cannot read property 'value' of undefined
      at parse (src/utils/parser.ts:42:18)
      at Object.<anonymous> (src/utils/parser.test.ts:28:20)
```

Can you figure out what's going wrong and fix it?
