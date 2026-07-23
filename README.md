# My Poor AI — Engineering Discipline for AI Coding Agents

**English** | [한국어](README.ko.md)

[![validate](https://github.com/hicucu/my-poor-ai/actions/workflows/validate.yml/badge.svg)](https://github.com/hicucu/my-poor-ai/actions/workflows/validate.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**my-poor-ai** makes Claude Code follow real engineering process — test-driven development, root-cause debugging, design review, and verified completion — instead of vibe-coding. It routes every request through an orchestrator that picks the right pipeline, dispatches specialized subagents, and refuses to report "done" without proof.

## Honestly, I'm Not Even Sure How

Other than this paragraph, I didn't write any of this. It's all AI-made.
It started with an instructions.md file and one line: "end sentences with noun forms." Then I added "put commas in big numbers," and kept tacking on lines until the file got big and I split it up. Then I split those pieces further, said "turn this file into a skill" to get skills, "make a command" to get slash commands, and "do this, do that" to get subagents — and around then the first plugin I'd ever installed was `superpowers`, and seeing it made me go "ooh, I want a plugin too~" and ask for one, and that snowballed into this.
Up through roughly v2 I had it polish stuff I'd written myself — "do this, do that, polish this one like this and that one like that." At some point I stopped doing anything except saying "just do the thing~."
And it (my-poor-ai, or whatever came before it) just handled everything on its own, start to finish.
Even this version started the same way — with the previous one, I began with "let's level this up," and eventually it became just "I want the thing~," and that's how we got here.
It told me there was some whole checklist for open-sourcing this, so I said sure, go for it, and writing this paragraph at the end is the only thing I actually did. (Only the Korean text is mine — the English is whatever the AI translated on its own.)

## Why My Poor AI

AI coding agents are fast but undisciplined: they fix symptoms instead of root causes, skip tests under pressure, and declare victory without verification. my-poor-ai counters this with **19 skills** (process rules the agent must follow), **24 subagents** (single-responsibility workers), and **12 slash commands**, wired together by an orchestrator that classifies each request and enforces the matching pipeline.

## Quick Start

### 1. Register the marketplace (once)

```
/plugin marketplace add hicucu/my-poor-ai
```

### 2. Install the plugin

```
/plugin install my-poor-ai@hicucu
/reload-plugins
```

### 3. Register the session-start workflow

Register my-poor-ai so its workflow applies automatically at the start of each session. Choose the setup for your CLI.

#### Claude Code CLI

Once registered, the `using-my-poor-ai` skill context is injected automatically at every session start (`/clear`, `/compact`, new session) via a SessionStart hook. Choose either automatic or manual registration.

**Automatic registration**

```
/my-poor-ai:setup
```

The `my-poor-ai:setup` skill reads and updates `~/.claude/settings.json` directly.

**Manual registration**

Add the following to the `hooks` section of `~/.claude/settings.json`:

```json
"SessionStart": [
  {
    "hooks": [
      {
        "type": "command",
        "command": "bash \"${CLAUDE_PLUGIN_ROOT}/hooks/session-start\"",
        "timeout": 10000
      }
    ]
  }
]
```

#### Codex

Codex has no SessionStart hook, so add the workflow block below to your global instructions file `~/.codex/AGENTS.md`. (For agent / multi-agent registration, see `/my-poor-ai:codex-setup`; this block defines the skill-dispatch workflow to follow each session.)

```markdown
# my-poor-ai default workflow

On every user request, before acting, answering, or asking clarifying questions, first identify and follow any applicable installed skill. If a `my-poor-ai:*` skill might apply, consult that skill's guidance first.

## Request classification

- Bugs, test failures, unexpected behavior: use `my-poor-ai:systematic-debugging`, then `my-poor-ai:verification-before-completion`.
- Small, well-defined code changes: use `my-poor-ai:test-driven-development`, then `my-poor-ai:verification-before-completion`.
- New features, multi-file changes, design decisions, external dependencies, or ambiguous requirements: follow `my-poor-ai:brainstorming` → `my-poor-ai:writing-plans` → `my-poor-ai:subagent-driven-development` → `my-poor-ai:requesting-code-review` → `my-poor-ai:finishing-a-development-branch`.
- New project with no decided tech stack or approach: use `my-poor-ai:which-way-should-i-go` first.

## Working rules

- Before feature work, use `my-poor-ai:brainstorming`; do not implement before the design is approved.
- Before implementing or fixing a bug, use `my-poor-ai:test-driven-development`.
- Before claiming something is done, fixed, or passing, always confirm the real result with `my-poor-ai:verification-before-completion`.
- When there are two or more independent tasks, consider `my-poor-ai:dispatching-parallel-agents`.
- When work is finished or ready to merge, use `my-poor-ai:finishing-a-development-branch`.
- When receiving code review feedback, do not accept it blindly — verify it technically with `my-poor-ai:receiving-code-review`.
- When writing, editing, or reviewing GitHub Actions workflows, use `my-poor-ai:preventing-github-actions-loops`.
```

## How It Works

Every request is classified into one of three pipelines:

| Path       | Trigger                                  | Pipeline                                                        |
| ---------- | ---------------------------------------- | --------------------------------------------------------------- |
| **DEBUG**  | Bugs, errors, unexpected behavior        | GOAL.md → systematic-debugging → verification                   |
| **SIMPLE** | 1–2 files, no design decisions, < 10 min | GOAL.md → TDD → verification                                    |
| **FULL**   | New features, complex changes            | brainstorming → planning → parallel development → 4-way review  |

The FULL path runs a 5-phase multi-agent pipeline: a brainstorming agent produces a design document (user approval gate), a planning agent decomposes it into TDD task specs, developer agents implement specs in parallel groups, and a review orchestrator fans out four parallel reviewers (architecture / security / performance / style), aggregates findings, and dispatches parallel issue-fixers.

## Core Components

- **19 skills** — TDD, systematic debugging, brainstorming, plan writing, code review (giving and receiving), multi-agent pipelines, doc sync, worktree isolation, skill authoring, and more
- **24 subagents** — project-context capture, a 10-agent docs suite, a 9-agent feature pipeline, and a 4-agent subagent-driven flow; each with a single responsibility and an explicit I/O contract (see `AGENTS.md`)
- **12 slash commands** — `/my-poor-ai:code-review`, `/my-poor-ai:detect-stack`, `/my-poor-ai:roles`, session management, stack detection, and setup utilities
- **Session handoff** — `HANDOFF.md` records narrative context at spec/phase completion so a fresh session can pick up mid-pipeline; `GOAL.md` tracks goal and success criteria as a completion gate
- **Multi-platform** — Claude Code first; agent definitions auto-generated for Codex (`.codex/agents/`), tool mappings for Copilot CLI, Gemini CLI, and an OpenCode test suite

## Role Presets

Not sure which skill to start with? `/my-poor-ai:roles` maps common roles to skill bundles:

| Role          | Bundle                                                              |
| ------------- | ------------------------------------------------------------------- |
| **Architect** | brainstorming → writing-plans → socratic-plan-review                 |
| **Builder**   | test-driven-development → subagent-driven-development → finishing    |
| **Debugger**  | systematic-debugging → verification-before-completion                |
| **Reviewer**  | requesting-code-review / receiving-code-review / `/my-poor-ai:code-review` |
| **Docs**      | sync-docs-from-diff / generate-claude-instructions                   |

## Verified, Not Vibes

my-poor-ai applies its own discipline to itself:

- **CI-validated on every push** — `validate-agents.mjs` checks frontmatter contracts (name/model/tool whitelists), reference resolution, and code-fence balance across 100+ markdown files; `generate-codex-agents.mjs --check` blocks drift between the 24 agent definitions and their generated Codex mirrors
- **Behaviorally tested** — skills are validated with RED–GREEN–PRESSURE runs against live agents; the worktree-isolation skill's guidance held at **50/50 runs with zero failures** (20 GREEN + 20 PRESSURE + 10 full-skill-text)
- **Adversarial pressure scenarios** — a dedicated test suite verifies that discipline skills (TDD, debugging) hold up under time pressure, sunk cost, and authority pressure — the exact conditions where agents rationalize shortcuts

## Repository Structure

```
my-poor-ai/
├── .claude-plugin/        # marketplace + plugin manifests
├── .codex/agents/         # auto-generated Codex agent definitions (do not edit)
├── agents/                # 24 subagent definitions (single source of truth)
├── commands/              # 12 slash commands
├── hooks/                 # SessionStart hooks (Claude Code + Cursor)
├── skills/                # 19 skill directories
├── scripts/               # CI validators + Codex mirror generator
├── tests/                 # deterministic + LLM-behavioral + pressure-scenario suites
├── docs/                  # recommended MCP pairings
├── AGENTS.md              # agent I/O contracts and invariants
├── CLAUDE.md              # working agreement for AI agents in this repo
├── CHANGELOG.md
├── CONTRIBUTING.md        # how to contribute
└── SECURITY.md
```

## Pairs Well With

my-poor-ai is pure instruction — no bundled integrations. [docs/recommended-mcp.md](docs/recommended-mcp.md) lists MCP servers that strengthen specific pipeline phases (docs lookup for design, browser automation for verification, GitHub for review flow). Everything works without them.

Want proof instead of promises? Two examples were built **entirely unattended** by this pipeline, each with a full commit trail and verification in its `PROVENANCE.md`:

- [examples/go-fractals/](examples/go-fractals/) — a working Go CLI (Sierpinski + Mandelbrot ASCII renderers) built from a given plan: 10 tasks, one TDD commit each, plus fix commits its own tests and review phase caught
- [examples/svelte-todo/](examples/svelte-todo/) — a Svelte 5 todo app built from **seven requirement bullets and nothing else**: the design document itself was authored by the brainstorming phase, then planned, implemented (30/30 tests), and review-fixed; the pipeline's own design/specs/review artifacts are preserved in `pipeline-artifacts/`

To reproduce: `bash tests/subagent-driven-dev/run-test.sh go-fractals` (invokes the `claude` CLI; real tokens, 10–30+ minutes).

## Origin

> 시작은 "문장 종결 어미를 명사로 해"였다.
>
> It all started with: *"End your sentences with nouns."*

A one-line Korean style request turned into a repository-wide audit, dead-code sweeps, a single-source agent-definition pipeline, and this public release. Fitting, for a project about engineering discipline.

## Philosophy

Test-first development. Systematic process over guessing. Evidence-based completion. Delegate complexity to specialized agents.

## Contributing

Issues and pull requests are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md). All changes must pass the CI validators, and skill changes require behavioral pressure-testing (see `skills/writing-skills/`).

Most repository documentation is written in Korean (the maintainer's working language); English issues and PRs are fully welcome.

## License

MIT — see [LICENSE](LICENSE). Concept attribution for [Superpowers](https://github.com/obra/superpowers) in [NOTICE](NOTICE).
