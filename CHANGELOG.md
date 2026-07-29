# Changelog

## 5.0.0

Aligns the plugin with what Claude Code actually is today. Four of these findings were things the plugin got wrong — a command writing a file nothing reads, a validator rejecting valid definitions, a hook missing two session events — and the rest are places where the plugin was carrying weight the harness now carries itself.

Every skill change in this release was developed against a measured baseline: a subagent run without the skill, its rationalizations recorded verbatim, then the same scenario re-run with the skill. One planned change was dropped because its baseline refused to reproduce.

### Breaking

- **`/my-poor-ai:generate-claudeignore` is gone.** Claude Code does not read `.claudeignore` — the command produced a file nothing consumes. Exclude files with `permissions.deny` rules in `settings.json` (`Read(./.env)`, `Read(./secrets/**)`) and `CLAUDE.md` memory files with `claudeMdExcludes`
- **`/my-poor-ai:code-review` produces a different set of files.** It no longer writes `reviews/style.md`, and `reviews/security.md` now comes from the bundled `/security-review` rather than a local agent. Anything reading the `reviews/` directory by filename needs updating
- **The `commands/` directory no longer exists.** All 11 commands live at `skills/<name>/SKILL.md`. Invocation is unchanged — `/my-poor-ai:code-review`, `/my-poor-ai:setup`, and the rest resolve exactly as before — but forks, in-flight pull requests, and links pointing at `commands/*.md` need to be repointed

### Added

- `--init-review-md` on `/my-poor-ai:code-review` — scaffolds a repository `REVIEW.md`. GitHub Code Review injects that file into every review agent as its highest-priority instruction block, so review rules land far more reliably there than in `CLAUDE.md`
- `agents/_shared/review-orchestration.md` — single source for the contract driving the four reviewers, the aggregator, and `issue-fixer`, following the `_shared/` pattern 4.2.0 established for the implementation workers
- `skills/using-my-poor-ai/router.md` — the compact routing block the SessionStart hook injects
- `README.md` (both languages) gained a section contrasting `/my-poor-ai:code-review` with the bundled `/code-review`: which axes each covers, what each produces, and that only the user can start the bundled one. They are namespaced, both should be run, and this plugin neither wraps nor replaces the bundled command

### Changed

- **`/my-poor-ai:code-review` now covers only the axes native review does not.** Architecture and performance run locally, the security axis is delegated to the bundled `/security-review`, and correctness is left to `/code-review`. Style is dropped as redundant — native review already reports reuse, simplification, and efficiency cleanups
  - `/code-review` and `/verify` are `disable-model-invocation` as of Claude Code v2.1.215, so no skill or command can run them on the user's behalf. The report now ends by telling the user to run `/code-review` and states plainly that the correctness axis is unreviewed until they do
  - `/security-review` takes no target argument and reviews the session's own repository and diff scope. When the target is a different repository, when `--files` selects specific files, or when the native skill cannot resolve a diff range, the security axis is recorded as not applicable rather than silently reviewing the wrong tree
- **The SessionStart hook injects 74% fewer bytes** (15,617 → 4,046) — a 79-line routing block instead of the full 282-line skill. Nothing was deleted; the full skill still loads on demand through the Skill tool, which is how progressive disclosure is meant to work. All five `tests/pipeline-triggering` prompts classify identically under both, including the two boundary cases, and under a five-way pressure set the agent skipped zero skill calls
- **`using-git-worktrees` rewritten around native tooling** (222 → 146 lines). The directory-selection priority list, the `git check-ignore` step, the legacy `~/.config/my-poor-ai/worktrees/` path, and the per-stack setup script are all standardized by native tooling now. The baseline showed the removed directory-priority clause was not merely obsolete — it *licensed* bypassing the native tool whenever a team document specified a different layout. Two traps that do cause real failures replaced it: the native base-branch default (branches from the default branch, not `HEAD`) and gitignored config missing from a fresh checkout (`.worktreeinclude`)
- **Review orchestration consolidated.** `review-agent` and `feature-pipeline` stay separate — they serve different purposes — but each now states only what is genuinely its own (`workspaceDir` naming, how the target is derived, invocation style, user gate, whether `issue-fixer` commits), tabulated in `AGENTS.md` invariant 10. This closed real drift: the migration-file check for new DB entities existed only in `review-agent`, so the `feature-pipeline` path never asked for it
- **All 11 commands moved into `skills/`.** Claude Code merged custom commands into skills; skills additionally support supporting-file directories, invocation control, and subagent execution, and two parallel trees meant two catalogs and two README pairs. `model`, `allowed-tools`, and `disable-model-invocation` all carry over as supported SKILL.md frontmatter
- `/my-poor-ai:setup`, `/my-poor-ai:codex-setup`, and `/my-poor-ai:graphify-setup` set `disable-model-invocation: true`. All three write outside the project — global `settings.json`, `~/.codex/config.toml`, package installs plus a git hook — and should run only when asked for by name
- `scripts/validate-agents.mjs` brought up to date with the frontmatter Claude Code actually supports, so it stops rejecting valid definitions: `model` accepts `fable` and pinned `claude-*` IDs; the tool whitelist covers the documented background-subagent set (`Skill`, `ToolSearch`, `EnterWorktree`, `Monitor`, `SendMessage`, …), which previously blocked any agent from declaring `Skill`; and `disallowedTools`, `skills`, `effort`, `isolation`, `background`, `maxTurns`, `memory`, `color` are accepted in a canonical order. The old fixed four-field order rejected `isolation: worktree` and every other modern field
- `scripts/validate-agents.mjs` gains checks in exchange: tools always stripped from subagents, fields ignored for plugin-distributed agents (`hooks`, `mcpServers`, `permissionMode`), unknown frontmatter keys, legacy `Task` references in prose, and — replacing the old command-catalog check — per-skill structure. The last one matters most: in a plugin skill the frontmatter `name` sets the command's final segment, so a mismatch against the directory name silently renames the command

### Fixed

- The SessionStart matcher was `startup|clear|compact`, missing `resume` and `fork`. Sessions restored with `--resume` / `--continue` / `/resume`, or branched with `/fork`, never received the context injection
- `using-my-poor-ai` stated the FULL path's first step two different ways — the decision graph said `brainstorming`, the phase table said `project-context`. The measured baseline flagged it as genuinely ambiguous; the router states it once, and Phase 0 is `project-context`
- Legacy `Task` tool references replaced with `Agent` across seven skill files (11 occurrences), including the Codex / Copilot / Gemini tool-mapping tables whose left column names the Claude Code tool. `Task` became `Agent` in v2.1.63
- Aggregator reports silently dropped issues belonging to no changed file — a missing migration, for instance. The aggregator now emits a `## 파일 미귀속 이슈` section, and the contract requires assigning such an issue to the file that should be created or reporting it explicitly
- `.codex-plugin/plugin.json` had drifted to 4.0.0 while the other manifests were at 4.2.0; this release resyncs all three
- `bump-version.sh --audit` was unusable at this version. It greps for the bare version string, so its noise scales with how common that string is — `4.2.0` matched almost nothing, `5.0.0` matched 21 unrelated npm dependency ranges under `examples/`. Those are standalone demo projects with independent versioning that must never be bumped alongside the plugin, so `examples` joins `audit.exclude`. The audit now reports one line, `CHANGELOG.md`, which is the accurate reminder that the changelog heading is written by hand

### Tests

- `tests/pressure-scenarios/worktree-native-preference-pressure.md` — measured baseline where the agent recognized `EnterWorktree` and explicitly rejected it, citing the team's onboarding document. Doubles as the regression test for the removed directory-priority clause
- `tests/pressure-scenarios/verification-pressure.md` scenario C — a negative-result regression test. `verification-before-completion` was deliberately left unchanged: across two measured baselines under five combined pressures, including a team lead delegating "merge if the three CI checks are green", the agent ran the app, caught a boot crash behind three green checks, noticed two of them were `echo` stubs, and refused to ship. With no reproducible failure there was no evidence to justify editing a tuned skill, so the scenario is recorded to catch the day that changes

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
