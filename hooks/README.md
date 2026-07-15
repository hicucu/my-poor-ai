# Hooks

**English** | [한국어](README.ko.md)

Claude Code / Cursor `SessionStart` hooks that inject the `using-my-poor-ai` skill context automatically at the start of every session (`/clear`, `/compact`, new session). Registered manually or via `/my-poor-ai:setup`.

| File | What it does |
| --- | --- |
| [`hooks.json`](hooks.json) | Claude Code `SessionStart` hook manifest — matches `startup\|clear\|compact` and runs `run-hook.mjs session-start`. |
| [`hooks-cursor.json`](hooks-cursor.json) | Cursor's equivalent `sessionStart` hook manifest. |
| [`run-hook.mjs`](run-hook.mjs) | Cross-platform runner — locates `bash` (Git Bash on Windows, system bash on Unix) and executes the named hook script. |
| [`session-start`](session-start) | The actual hook script — injects the `using-my-poor-ai` skill content into context, and warns if a legacy `~/.config/my-poor-ai/skills` directory is still present. |
