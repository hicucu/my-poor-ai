# Hooks

[English](README.md) | **한국어**

모든 세션 시작(`/clear`, `/compact`, 신규 세션) 시 `using-my-poor-ai` 스킬 컨텍스트를 자동 주입하는 Claude Code / Cursor `SessionStart` 훅. 수동 등록 또는 `/my-poor-ai:setup`으로 등록함.

| 파일 | 역할 |
| --- | --- |
| [`hooks.json`](hooks.json) | Claude Code `SessionStart` 훅 매니페스트 — `startup\|clear\|compact`에 매칭, `run-hook.mjs session-start` 실행 |
| [`hooks-cursor.json`](hooks-cursor.json) | Cursor용 `sessionStart` 훅 매니페스트 |
| [`run-hook.mjs`](run-hook.mjs) | 크로스플랫폼 러너 — bash(Windows는 Git Bash, Unix는 시스템 bash) 위치를 찾아 지정된 훅 스크립트 실행 |
| [`session-start`](session-start) | 실제 훅 스크립트 — `using-my-poor-ai` 스킬 내용을 컨텍스트에 주입, 레거시 `~/.config/my-poor-ai/skills` 디렉터리가 남아있으면 경고 |
