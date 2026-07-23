# Agents

**English** | [한국어](README.ko.md)

24 single-responsibility subagent definitions, each with an explicit input/output contract. Every agent file's YAML frontmatter (`name`, `description`, `model`, `tools`) is the source of truth; this page is a quick index. For call order, phase diagrams, and full I/O contracts, see [`AGENTS.md`](../AGENTS.md) at the repository root.

None of these agents are called directly by the user — they're invoked by a skill orchestrator (`feature-pipeline`, `generate-claude-instructions`, `sync-docs-from-diff`) or spawned by the main agent via `subagent_type` in the `using-my-poor-ai` FULL path.

## Common infrastructure (1)

| Agent | What it does |
| --- | --- |
| [`project-context.md`](project-context.md) | Captures project structure, stack, conventions, and recent commits before feature work starts; cached for 24h. |

## docs-suite — `generate-claude-instructions` group (5)

Phase 1 runs the first four in parallel; Phase 2 runs the composer to synthesize their output into `CLAUDE.md`.

| Agent | What it does |
| --- | --- |
| [`dev-principles.md`](dev-principles.md) | Writes `DEVELOPMENT.md` — SOLID, TDD, clean code, security, and performance principles. |
| [`language-guidelines.md`](language-guidelines.md) | Writes `LANGUAGE_GUIDELINES.md` — one section per detected language/framework. |
| [`ai-behavior.md`](ai-behavior.md) | Writes `AI_BEHAVIOR.md` — response format, workflow, tool use, and self-verification principles. |
| [`commit-convention.md`](commit-convention.md) | Writes `COMMIT_CONVENTION.md` based on the project's commitlint config or Conventional Commits. |
| [`claude-md-composer.md`](claude-md-composer.md) | Reads the four docs above and synthesizes a concise, pointer-style `CLAUDE.md`. |

## docs-suite — `sync-docs-from-diff` group (5)

Analyzes a branch diff, proposes documentation patches (never edits files directly), and validates them after user approval.

| Agent | What it does |
| --- | --- |
| [`change-analyzer.md`](change-analyzer.md) | Analyzes the git diff/commits from a base branch to HEAD into a structured change-analysis report. |
| [`readme-updater.md`](readme-updater.md) | Proposes updates to the project root `README.md` when user-facing behavior has changed. |
| [`docs-updater.md`](docs-updater.md) | Proposes updates to everything under `./docs/` — guides, tutorials, API references, architecture docs. |
| [`inline-doc-updater.md`](inline-doc-updater.md) | Proposes updates to inline docs (component README, module notes) near the changed files. |
| [`doc-sync-validator.md`](doc-sync-validator.md) | After patches are applied, verifies nothing was missed, wording is consistent, and links/signatures still match. |

## feature-pipeline group (9)

The 5-phase pipeline behind the `feature-pipeline` skill: plan → implement → test → review → fix.

| Agent | What it does |
| --- | --- |
| [`feature-planner.md`](feature-planner.md) | Detects the stack and decomposes a feature request into a file-level plan (`stack-profile.json`, `plan.md`, `file-manifest.json`). |
| [`file-developer.md`](file-developer.md) | Implements or modifies a single file per its spec, language/framework-agnostic. |
| [`test-writer.md`](test-writer.md) | Writes unit tests for business-logic files, auto-selecting the detected test framework. |
| [`architecture-reviewer.md`](architecture-reviewer.md) | Reviews changed code for layering violations, coupling, SRP, and abstraction issues. |
| [`security-reviewer.md`](security-reviewer.md) | Reviews changed code against OWASP Top 10, auth, input validation, and secrets exposure. |
| [`performance-reviewer.md`](performance-reviewer.md) | Reviews changed code for N+1 queries, blocking I/O, memory leaks, and caching issues. |
| [`style-reviewer.md`](style-reviewer.md) | Reviews changed code for naming, duplication, magic numbers, and language idioms. |
| [`review-aggregator.md`](review-aggregator.md) | Merges the four reviewer reports into a single file-grouped `review-report.md`. |
| [`issue-fixer.md`](issue-fixer.md) | Fixes the review issues for a single file per the stack's conventions. |

## subagent-driven flow group (4)

Spawned directly by the main agent via `subagent_type` on the `using-my-poor-ai` FULL path (complex features). Reuses the four reviewers, `review-aggregator`, and `issue-fixer` from the feature-pipeline group above for its review phase.

| Agent | What it does |
| --- | --- |
| [`brainstorming-agent.md`](brainstorming-agent.md) | Analyzes a requirement, compares 2–3 design options, and writes `design.md` for user approval. |
| [`planning-agent.md`](planning-agent.md) | Breaks an approved `design.md` down into TDD task specs (`specs/*.md`, `file-manifest.json`). |
| [`developer-agent.md`](developer-agent.md) | Implements a single spec via TDD (RED-GREEN-REFACTOR) and commits. |
| [`review-agent.md`](review-agent.md) | Review orchestrator — spawns the 4 reviewers, the aggregator, and the issue-fixer as its own flow. |

---

> **`_shared/`** is not an agent. `_shared/implementation-conventions.md` holds the coding discipline + stack-convention matrix shared by the two implementation workers (`developer-agent`, `file-developer`), pulled in via an `@include` directive. See invariant 9 in [`AGENTS.md`](../AGENTS.md).
