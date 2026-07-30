# Skills

**English** | [한국어](README.ko.md)

31 skill directories. Each skill's `SKILL.md` frontmatter `description` is the actual trigger condition; this page is a quick index grouped by development phase. See `skills/writing-skills/` before editing any of these.

Two kinds live here, and the difference is who starts them:

- **Process skills** (the first sections below) — Claude loads these automatically when the request matches, and you can also invoke them with `/my-poor-ai:{name}`.
- **Task commands** ([Commands](#commands) below) — you invoke these with `/my-poor-ai:{name}`. The ones with side effects outside the project set `disable-model-invocation: true` so Claude never starts them on its own.

Both are the same file format. Claude Code merged custom commands into skills, so `commands/deploy.md` and `skills/deploy/SKILL.md` produce the same `/deploy`; skills additionally support a directory of supporting files, invocation control, and subagent execution. This plugin keeps everything under `skills/` for that reason.

## Design & planning

| Skill | What it does |
| --- | --- |
| [`brainstorming`](brainstorming/) | Must-use before any creative work — explores intent, requirements, and design through dialogue before implementation. |
| [`which-way-should-i-go`](which-way-should-i-go/) | Compares 2–3 generationally-proven approaches via parallel web research before deciding what to even brainstorm. |
| [`writing-plans`](writing-plans/) | Writes a comprehensive, bite-sized implementation plan assuming a context-free, taste-questionable engineer. |
| [`socratic-plan-review`](socratic-plan-review/) | Validates a complex plan via 7 categories of structured Socratic questioning before execution. |

## Implementation

| Skill | What it does |
| --- | --- |
| [`test-driven-development`](test-driven-development/) | RED → GREEN → REFACTOR discipline — write the failing test before the implementation. |
| [`subagent-driven-development`](subagent-driven-development/) | Executes a plan of independent tasks by dispatching a fresh subagent per task, with a 2-stage review after each. |
| [`executing-plans`](executing-plans/) | Loads a written plan in a separate session, critiques it, executes all tasks, and reports on completion. |
| [`using-git-worktrees`](using-git-worktrees/) | Ensures an isolated workspace via native worktree tooling, falling back to a manual git worktree. |
| [`dispatching-parallel-agents`](dispatching-parallel-agents/) | Delegates 2+ independent tasks to isolated subagents with precisely constructed instructions and context. |

## Completion & review

| Skill | What it does |
| --- | --- |
| [`verification-before-completion`](verification-before-completion/) | Requires running and inspecting verification commands before claiming work is complete, fixed, or passing. |
| [`finishing-a-development-branch`](finishing-a-development-branch/) | Presents structured merge/PR/cleanup options once implementation is done and tests pass. |
| [`requesting-code-review`](requesting-code-review/) | Dispatches a reviewer subagent early and often, with context built specifically for evaluation. |
| [`receiving-code-review`](receiving-code-review/) | Requires technical verification of review feedback rather than blind or performative agreement. |

## Debugging

| Skill | What it does |
| --- | --- |
| [`systematic-debugging`](systematic-debugging/) | Root-causes bugs, test failures, and unexpected behavior methodically instead of guessing fixes. |

## Documentation & CI safety

| Skill | What it does |
| --- | --- |
| [`sync-docs-from-diff`](sync-docs-from-diff/) | Analyzes a branch diff to propose README/docs/inline doc updates, applied only after user review. |
| [`preventing-github-actions-loops`](preventing-github-actions-loops/) | Detects and prevents self-triggering GitHub Actions workflow cycles. |

## Pipeline

| Skill | What it does |
| --- | --- |
| [`feature-pipeline`](feature-pipeline/) | Stack-agnostic 5-phase multi-agent pipeline: plan → implement → test → review → fix. |

## Meta

| Skill | What it does |
| --- | --- |
| [`writing-skills`](writing-skills/) | TDD applied to process documents — how to author, edit, and validate a skill. |
| [`using-my-poor-ai`](using-my-poor-ai/) | Request classification and pipeline contracts. `router.md` beside it is the compact block the SessionStart hook injects. |

## Commands

Invoked as `/my-poor-ai:{name}`. `/my-poor-ai:commands` is the catalog entry point.

| Skill | What it does | Model-invocable |
| --- | --- | --- |
| [`my-poor-ai`](my-poor-ai/) | Entry point that routes any request into the right pipeline (DEBUG / SIMPLE / FULL) — usable without setup. | Yes |
| [`commands`](commands/) | Lists what is available; the catalog itself. | Yes |
| [`roles`](roles/) | Role-preset catalog — routes a role name (architect / builder / debugger / reviewer / docs) to its skill bundle. | Yes |
| [`code-review`](code-review/) | Reviews the architecture and performance axes, folds in native `/security-review`, and can scaffold a `REVIEW.md`. | Yes |
| [`detect-stack`](detect-stack/) | Scans marker files to detect the tech stack and generate `stack-profile.json`. | Yes |
| [`git-resume`](git-resume/) | Reconstructs prior work context from commit history, given a time expression or a commit hash. | Yes |
| [`session-manager`](session-manager/) | Lists, renames, or deletes local Claude Code sessions. | Yes |
| [`weekly-commits`](weekly-commits/) | Prints this week's commits for a given author as a markdown table. | Yes |
| [`setup`](setup/) | Registers the `SessionStart` hook in `~/.claude/settings.json`. | **No** — writes outside the project |
| [`codex-setup`](codex-setup/) | Registers this plugin's agents in `~/.codex/config.toml`. | **No** — writes outside the project |
| [`graphify-setup`](graphify-setup/) | Installs and configures a code-graph tool (`graphifyy` or `codegraph`). | **No** — installs packages, adds a git hook |
