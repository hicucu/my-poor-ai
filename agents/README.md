# Agents

**English** | [한국어](README.ko.md)

16 single-responsibility subagent definitions, each with an explicit input/output contract. Every agent file's YAML frontmatter (`name`, `description`, `model`, `tools`) is the source of truth; this page is a quick index. For call order, phase diagrams, and full I/O contracts, see [`AGENTS.md`](../AGENTS.md) at the repository root.

None of these agents are called directly by the user. The main agent spawns them via `subagent_type` on the `using-my-poor-ai` FULL path, `review-agent` spawns the review group, and the `sync-docs-from-diff` and `code-review` skills spawn their own.

## Common infrastructure (1)

| Agent | What it does |
| --- | --- |
| [`project-context.md`](project-context.md) | Captures project structure, stack, conventions, and recent commits before feature work starts, plus the shared `stack-profile.json`; cached for 24h. |

## docs-suite — `sync-docs-from-diff` group (5)

Analyzes a branch diff, proposes documentation patches (never edits files directly), and validates them after user approval.

| Agent | What it does |
| --- | --- |
| [`change-analyzer.md`](change-analyzer.md) | Analyzes the git diff/commits from a base branch to HEAD into a structured change-analysis report. |
| [`readme-updater.md`](readme-updater.md) | Proposes updates to the project root `README.md` when user-facing behavior has changed. |
| [`docs-updater.md`](docs-updater.md) | Proposes updates to everything under `./docs/` — guides, tutorials, API references, architecture docs. |
| [`inline-doc-updater.md`](inline-doc-updater.md) | Proposes updates to inline docs (component README, module notes) near the changed files. |
| [`doc-sync-validator.md`](doc-sync-validator.md) | After patches are applied, verifies nothing was missed, wording is consistent, and links/signatures still match. |

## Review group (6)

Spawned by `review-agent` (FULL path, Phase 4). `/my-poor-ai:code-review` reuses `architecture-reviewer`, `performance-reviewer`, and `review-aggregator`.

| Agent | What it does |
| --- | --- |
| [`architecture-reviewer.md`](architecture-reviewer.md) | Reviews changed code for layering violations, coupling, SRP, abstraction issues, and schema changes that ship without a migration. |
| [`security-reviewer.md`](security-reviewer.md) | Reviews changed code against OWASP Top 10, auth, input validation, and secrets exposure. `review-agent` uses it only when the native `/security-review` can't run. |
| [`performance-reviewer.md`](performance-reviewer.md) | Reviews changed code for N+1 queries, blocking I/O, memory leaks, and caching issues. |
| [`style-reviewer.md`](style-reviewer.md) | Reviews changed code for naming, duplication, magic numbers, and language idioms. |
| [`review-aggregator.md`](review-aggregator.md) | Merges the axis reports it receives into a single file-grouped `review-report.md`. |
| [`issue-fixer.md`](issue-fixer.md) | Fixes the review issues for a single file per the stack's conventions. |

## Subagent-driven flow group (4)

Spawned directly by the main agent via `subagent_type` on the `using-my-poor-ai` FULL path (complex features).

| Agent | What it does |
| --- | --- |
| [`brainstorming-agent.md`](brainstorming-agent.md) | Analyzes a requirement, compares 2–3 design options, and writes `design.md` for user approval. |
| [`planning-agent.md`](planning-agent.md) | Breaks an approved `design.md` down into TDD task specs (`specs/*.md`, `file-manifest.json`), checking the project's migration approach first when the design changes a DB schema. |
| [`developer-agent.md`](developer-agent.md) | Implements a single spec via TDD (RED-GREEN-REFACTOR) and commits. |
| [`review-agent.md`](review-agent.md) | Review orchestrator — runs the architecture, performance, and style reviewers alongside the native `/security-review` (falling back to `security-reviewer`), then the aggregator and the issue-fixer. |

---

> **`_shared/`** is not an agent. `_shared/implementation-conventions.md` holds the coding discipline, stack-convention matrix, and test conventions for `developer-agent`; `_shared/stack-profile.md` holds the stack-detection matrix and `stack-profile.json` schema for `project-context`; `_shared/review-orchestration.md` holds the review fan-out contract for `review-agent`. All are pulled in via an `@include` directive. See invariant 9 in [`AGENTS.md`](../AGENTS.md).
>
> Only the files listed under `agents` in `.claude-plugin/plugin.json` load as agents, so this README and `_shared/` stay out of the agent list (invariant 11).
