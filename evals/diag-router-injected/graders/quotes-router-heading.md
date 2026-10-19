---
type: regex
weight: 1
pattern: "my-poor-ai 라우터"
---

Passes when the final message contains the routing block's first heading text (`# my-poor-ai 라우터`, the first line of skills/using-my-poor-ai/router.md, injected by hooks/session-start). Under --ablation with-without the without arm should answer NONE and fail.

(Notes only — the pattern comes from the frontmatter.)
