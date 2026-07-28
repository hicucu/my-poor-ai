# Changelog

## Unreleased

Pipeline audit against current Claude Code capabilities; the four lowest-risk findings applied. Version intentionally not bumped — release is a separate decision.

### Fixed

- `hooks/hooks.json` SessionStart matcher was `startup|clear|compact`, missing the `resume` and `fork` events. Sessions restored with `--resume` / `--continue` / `/resume`, or branched with `/fork`, never received the `using-my-poor-ai` context injection

### Changed

- `scripts/validate-agents.mjs` brought up to date with the frontmatter Claude Code actually supports, so the validator stops rejecting valid definitions:
  - `model` accepts `fable` and pinned full model IDs (`claude-*`) alongside the existing aliases
  - tool whitelist widened to the documented background-subagent tool set (`Skill`, `ToolSearch`, `TodoWrite`, `EnterWorktree`, `ExitWorktree`, `Monitor`, `SendMessage`, `TaskStop`, `Artifact`, `PowerShell`), which previously blocked any agent from declaring `Skill`
  - optional fields `disallowedTools`, `skills`, `effort`, `isolation`, `background`, `maxTurns`, `memory`, `color` are now accepted in a defined canonical order; the old fixed four-field order rejected `isolation: worktree` and every other modern field
- `scripts/validate-agents.mjs` gains two checks in exchange: tools that are always stripped from subagents (`AskUserQuestion`, `EnterPlanMode`, `TaskOutput`, …) and fields ignored for plugin-distributed agents (`hooks`, `mcpServers`, `permissionMode`) are now reported instead of silently failing at runtime, along with unknown frontmatter keys
- `/my-poor-ai:setup`, `/my-poor-ai:codex-setup`, and `/my-poor-ai:graphify-setup` set `disable-model-invocation: true`. All three write outside the project (global `settings.json`, `~/.codex/config.toml`, package installs plus a git hook) and should run only when the user asks for them by name

- `using-git-worktrees` rewritten around native worktree tooling (222 → 146 lines). The directory-selection priority list, the `git check-ignore` safety step, the legacy `~/.config/my-poor-ai/worktrees/` global path, and the per-stack setup script are gone — native tooling standardizes all of them. Measured baseline showed the removed directory-priority clause was not merely obsolete: it *licensed* bypassing the native tool whenever a team document specified a different layout. Two traps that do cause real failures were added in their place: the native base-branch default (branches from the default branch, not `HEAD`) and gitignored config files missing from a fresh checkout (`.worktreeinclude`)

- `/my-poor-ai:code-review` narrowed to the axes native review does not emphasize. The four-way local fan-out (architecture / security / performance / style) becomes architecture + performance locally, with the security axis delegated to the bundled `/security-review` and the correctness axis left to `/code-review`. Style is dropped as redundant — native review already reports reuse, simplification, and efficiency cleanups. A new `--init-review-md` mode scaffolds a repository `REVIEW.md`, which GitHub Code Review injects into every review agent as its highest-priority instruction block
  - `/code-review` and `/verify` are `disable-model-invocation` as of v2.1.215, so no skill or command can run them on the user's behalf; the report now ends by telling the user to run `/code-review` and states plainly that the correctness axis is unreviewed until they do
  - Measured (2026-07-28): `/security-review` takes no target argument and reviews the session's own repository and diff scope. When the review target is a different repository, when `--files` selects specific files, or when the native skill cannot resolve a diff range, the security axis is now recorded as not applicable rather than silently reviewing the wrong tree

- The SessionStart hook now injects a compact routing block (`skills/using-my-poor-ai/router.md`, 79 lines) instead of the full 282-line `using-my-poor-ai` skill — a 74% cut in per-session injected bytes (15,617 → 4,046). Nothing was deleted: the full skill still loads on demand through the Skill tool, which is how progressive disclosure is supposed to work. The router carries only what classification needs — instruction precedence, the three-path split, the SIMPLE criteria, the per-path first step, the GOAL.md soft gate, and the six highest-value red flags
  - Measured (2026-07-28): all five `tests/pipeline-triggering` prompts classify identically under the router and under full injection, including both boundary cases. Under a five-way pressure set — triviality, a user explicitly demanding the design step be skipped, and a production outage — the agent skipped zero skill calls and classified all five correctly
  - Fixed along the way: the full skill states the FULL path's first step two different ways (the decision graph says `brainstorming`, the phase table says `project-context`), which the measured baseline flagged as genuinely ambiguous. The router states it once — Phase 0 is `project-context`

- Review fan-out orchestration consolidated into `agents/_shared/review-orchestration.md`, following the `_shared/` pattern 4.2.0 established for the implementation workers. `review-agent` and `feature-pipeline` stay separate — they serve different purposes — but the contract for driving the four reviewers, the aggregator, and `issue-fixer` now has one source: review axes and output filenames, the single-response parallel dispatch rule, the reviewers' common input, the aggregator's input/output, the severity-to-action mapping, per-file fix dispatch, and the prohibitions. Each caller now states only what is genuinely its own (`workspaceDir` naming, how the review target is derived, invocation style, user gate, whether `issue-fixer` commits), documented as a table in `AGENTS.md` invariant 10
  - This closed real drift. The migration-file check for new DB entities existed only in `review-agent`, so the `feature-pipeline` path never asked for it; stack-profile passing existed only in `feature-pipeline`. Verified after the change: both callers now issue the migration instruction verbatim from the shared module
  - `/my-poor-ai:code-review` deliberately does *not* follow this contract (two axes, security delegated to native). Recorded in the shared module and in `AGENTS.md` as an intended divergence rather than drift
  - Gap found while verifying and closed: an issue like a missing migration belongs to no changed file, so per-file dispatch silently dropped it. The aggregator now emits a `## 파일 미귀속 이슈` section, and the contract requires either assigning such an issue to the file that should be created or reporting it to the user explicitly

- Legacy `Task` tool references replaced with `Agent` across seven skill files (11 occurrences), including the Codex / Copilot / Gemini tool-mapping tables whose left column names the Claude Code tool. `Task` became `Agent` in v2.1.63; the validator checked only agent frontmatter, so prose kept the old name
- `scripts/validate-agents.mjs` now flags legacy `Task` references in prose across `agents/`, `commands/`, and `skills/`. Scoped to `Task` followed by 도구/tool (or backtick-wrapped) so plan numbering like "Task 1" and C# `async Task` examples do not trip it; `tests/` is excluded because plan fixtures there are full of task numbering
- README (both languages) gained a section contrasting `/my-poor-ai:code-review` with the bundled `/code-review` — which axes each covers, what each produces, and that only the user can start the bundled one. The two are namespaced and both should be run; the plugin does not wrap or replace the bundled command

- All 11 slash commands moved from `commands/*.md` to `skills/<name>/SKILL.md`, and the `commands/` directory is gone. Claude Code merged custom commands into skills, so both forms produce the same `/name` — but skills additionally support a directory of supporting files, invocation control, and subagent execution, and keeping two parallel trees meant maintaining two catalogs and two README pairs. Invocation is unchanged: plugin skills and commands share the `my-poor-ai:` namespace, so `/my-poor-ai:code-review`, `/my-poor-ai:setup`, and the rest still resolve exactly as before, and `model` / `allowed-tools` / `disable-model-invocation` carry over as supported SKILL.md frontmatter
  - `skills/README.md` (both languages) absorbs `commands/README.md` and now separates process skills from task commands by who starts them, marking the three setup commands that cannot be model-invoked
  - `scripts/validate-agents.mjs` check 4 was rewritten from "every command is in the catalog" to cover every skill: a `skills/` subdirectory must contain `SKILL.md`, that file needs `name` and `description`, `name` must equal the directory name, and the skill must be introduced in `skills/README.md`, the commands catalog, `CLAUDE.md`, or `README.md`. The name check matters most — in a plugin skill the frontmatter `name` sets the command's last segment, so a mismatch silently renames the command

### Removed

- `/my-poor-ai:generate-claudeignore` — Claude Code does not read `.claudeignore`; the command produced a file nothing consumes. File access is excluded through `permissions.deny` rules in `settings.json` (`Read(./.env)`, `Read(./secrets/**)`), and `CLAUDE.md` memory files through `claudeMdExcludes`

### Tests

- `tests/pressure-scenarios/worktree-native-preference-pressure.md` — measured baseline (2026-07-28) where the agent recognized `EnterWorktree` and explicitly rejected it, citing the team's onboarding document; doubles as the regression test for the removed directory-priority clause
- `tests/pressure-scenarios/verification-pressure.md` scenario C — a negative-result regression test. `verification-before-completion` was left unchanged: across two measured baselines under five combined pressures, the agent ran the app, caught a crash behind three green CI checks, and refused to ship. With no reproducible failure there is no evidence to justify editing a tuned skill, so the scenario is recorded to catch the day that changes

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
