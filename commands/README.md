# Commands

**English** | [한국어](README.ko.md)

12 slash commands, invoked as `/my-poor-ai:{command}`. `/my-poor-ai:commands` is the catalog entry point; the rest are listed here for quick reference.

| Command | What it does |
| --- | --- |
| [`my-poor-ai.md`](my-poor-ai.md) | Entry point that routes any request into the right pipeline (DEBUG / SIMPLE / FULL) — usable without setup. |
| [`commands.md`](commands.md) | Lists available commands; the `/my-poor-ai:commands` catalog itself. |
| [`setup.md`](setup.md) | Registers the `SessionStart` hook in `~/.claude/settings.json` automatically. |
| [`codex-setup.md`](codex-setup.md) | Registers my-poor-ai's agents and multi-agent features in `~/.codex/config.toml`. |
| [`roles.md`](roles.md) | Role-preset catalog — routes a role name (Architect/Builder/Debugger/Reviewer/Docs) to its skill bundle. |
| [`code-review.md`](code-review.md) | Standalone 4-way parallel code review (architecture/security/performance/style) with an aggregated report. |
| [`detect-stack.md`](detect-stack.md) | Scans marker files to detect the tech stack and generate `stack-profile.json`, without running the full feature pipeline. |
| [`git-resume.md`](git-resume.md) | Reconstructs prior work context from commit history, given a time expression ("yesterday", "last week") or a commit hash. |
| [`generate-claudeignore.md`](generate-claudeignore.md) | Generates an optimized `.claudeignore` based on the detected stack and actual files, or merges missing entries into an existing one. |
| [`graphify-setup.md`](graphify-setup.md) | One-stop install/setup for a code-graph tool (`graphifyy` or `codegraph`) — package install, graph generation, Claude Code integration, git hook. |
| [`session-manager.md`](session-manager.md) | Analyzes all local Claude Code sessions with up to 10 parallel subagents; list, rename, or delete. |
| [`weekly-commits.md`](weekly-commits.md) | Prints this week's commits for a given GitHub ID/name as a markdown table, split by project for monorepos. |
