# Changelog

## 4.2.0

### Added

- Per-directory `README.md` / `README.ko.md` for `agents/`, `commands/`, `hooks/`, `skills/` — bilingual quick indexes of what each agent, command, hook, and skill does
- `agents/_shared/implementation-conventions.md` — single-source coding discipline + stack-convention matrix shared by the implementation workers (`developer-agent`, `file-developer`) via an `@include` directive

### Changed

- `developer-agent` / `file-developer` now pull shared conventions from the `_shared/` module instead of duplicating them; `developer-agent` gains the (profile-optional) stack-convention matrix it previously lacked. The FULL and feature-pipeline pipelines stay separate — only the shared guidance is centralized
- `scripts/generate-codex-agents.mjs` expands `@include: <path>` inline so Codex mirrors stay self-contained; `scripts/validate-agents.mjs` verifies `@include` targets resolve
- `AGENTS.md` invariant 9 documents the shared implementation-conventions module

## 4.1.0

### Added

- `which-way-should-i-go` skill — pre-brainstorming direction decision: three generational lenses (established-mainstream / modern-standard / emerging-hot) researched by parallel subagents with web search; includes a two-stage gate (goal → direction, ask the user when ambiguous), Socratic goal elicitation when even the goal is unclear, and a new-market variant when no legacy approach exists
- `tests/pressure-scenarios/` entries for `which-way-should-i-go` and `socratic-plan-review` — measured (not hypothetical) baselines with dates, doubling as regression tests
- `docs/skill-development-process.md` — standard skill development workflow (TDD cycle, SSOT checks, dependency back-review, test assetization)

### Fixed

- `socratic-plan-review`: activation threshold mismatch between description (5+ files) and body (3+ files) unified to 3+; removed workflow summary from description (CSO rule); added question-routing rule (self-check first, escalate only policy questions to the user) and record-location rule (`_workspaces/{branch-slug}/socratic-review.md`)

## 4.0.0

First public release. Version continues from the private predecessor (forge 3.3.0).

### Changed

- Project renamed to **my-poor-ai**; plugin name, command namespace (`/my-poor-ai:*`), and skill namespace updated accordingly
- README rewritten in English; manifests carry English descriptions
- Open-contribution policy: external issues/PRs welcome (CONTRIBUTING.md, issue/PR templates, SECURITY.md)
- All Korean documentation unified to noun-ending sentence style (명사형 종결), including generated-artifact style rules

### Added

- `/my-poor-ai:roles` — role-preset entry points (architect / builder / debugger / reviewer / docs)
- `scripts/generate-codex-agents.mjs` — `.codex/agents/*.toml` auto-generated from `agents/*.md` (single source of truth), with `--check` drift gate in CI
- SKILL frontmatter descriptions in English for reliable skill triggering in English-language sessions
- `docs/recommended-mcp.md` — recommended MCP server pairings

### Removed

- Dead content: `feature-dev` skill (absorbed by using-my-poor-ai FULL path + feature-pipeline), legacy `run-hook.cmd`, superseded document-reviewer prompt templates
