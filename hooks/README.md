# Hooks

**English** | [한국어](README.ko.md)

Claude Code / Cursor `SessionStart` hooks that inject the compact routing block (`skills/using-my-poor-ai/router.md`) at the start of every session — new session, `--resume` / `/resume`, `/clear`, `/compact`, `/fork`. Claude Code loads `hooks.json` automatically while the plugin is enabled; there is no registration step. The full `using-my-poor-ai` skill loads on demand through the Skill tool.

| File | What it does |
| --- | --- |
| [`hooks.json`](hooks.json) | Claude Code `SessionStart` hook manifest — matches `startup\|resume\|clear\|compact\|fork` and runs `run-hook.mjs session-start`. |
| [`hooks-cursor.json`](hooks-cursor.json) | Cursor's equivalent `sessionStart` hook manifest. Written against Cursor's 2026-07 hook format; not verified on current releases. |
| [`run-hook.mjs`](run-hook.mjs) | Cross-platform runner — locates `bash` (Git Bash on Windows, system bash on Unix) and executes the named hook script. |
| [`session-start`](session-start) | The actual hook script — emits `router.md` as `hookSpecificOutput.additionalContext` (Claude Code), `additional_context` (Cursor), or top-level `additionalContext` (Copilot CLI and others). |
