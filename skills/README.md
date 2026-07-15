# Skills

**English** | [한국어](README.ko.md)

20 skill directories — process rules the agent must follow, triggered automatically by request content or invoked explicitly. Each skill's `SKILL.md` frontmatter `description` is the actual trigger condition; this page is a quick index grouped by development phase. See `skills/writing-skills/` before editing any of these.

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
| [`generate-claude-instructions`](generate-claude-instructions/) | Orchestrator that generates `CLAUDE.md` and its reference docs (`DEVELOPMENT.md`, `LANGUAGE_GUIDELINES.md`, `AI_BEHAVIOR.md`, `COMMIT_CONVENTION.md`). |
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
| [`using-my-poor-ai`](using-my-poor-ai/) | Loaded at the start of every conversation; establishes how to find and invoke skills. |
