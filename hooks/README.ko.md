# Hooks

[English](README.md) | **한국어**

모든 세션 시작(신규 세션, `--resume` / `/resume`, `/clear`, `/compact`, `/fork`) 시 압축 라우팅 블록(`skills/using-my-poor-ai/router.md`)을 주입하는 Claude Code / Cursor `SessionStart` 훅. Claude Code는 플러그인이 활성화되어 있으면 `hooks.json`을 자동으로 로드하므로 별도 등록 절차가 없음. `using-my-poor-ai` 스킬 전문은 Skill 도구로 필요할 때 로드됨.

| 파일 | 역할 |
| --- | --- |
| [`hooks.json`](hooks.json) | Claude Code `SessionStart` 훅 매니페스트 — `startup\|resume\|clear\|compact\|fork`에 매칭, `run-hook.mjs session-start` 실행 |
| [`hooks-cursor.json`](hooks-cursor.json) | Cursor용 `sessionStart` 훅 매니페스트. 2026-07 기준 Cursor 훅 형식으로 작성, 현행 릴리스에서 미검증 |
| [`run-hook.mjs`](run-hook.mjs) | 크로스플랫폼 러너 — bash(Windows는 Git Bash, Unix는 시스템 bash) 위치를 찾아 지정된 훅 스크립트 실행 |
| [`session-start`](session-start) | 실제 훅 스크립트 — `router.md`를 `hookSpecificOutput.additionalContext`(Claude Code), `additional_context`(Cursor), 최상위 `additionalContext`(Copilot CLI 등)로 출력 |
