# Contributing to my-poor-ai

Issues and pull requests are welcome. This document explains how to contribute without breaking the contracts that hold the system together.

## Ground Rules

1. **Skills are code, not prose.** Every skill file shapes agent behavior. Carefully tuned content (Red Flags tables, rationalization lists, key phrases) must not be changed without behavioral evidence — see `skills/writing-skills/SKILL.md` for the RED–GREEN–REFACTOR methodology.
2. **`agents/*.md` is the single source of truth.** Never edit `.codex/agents/*.toml` by hand — they are generated. After changing any `agents/*.md`, run `node scripts/generate-codex-agents.mjs` and commit the regenerated tomls.
3. **Korean documentation convention.** Repository docs are written in Korean with noun-ending sentences (명사형 종결: -함/-됨/-임). English contributions are welcome in issues and PR descriptions; if you change Korean docs, keep the ending style consistent. The README and this file are English.
4. **YAGNI.** Don't add features or refactors beyond the scope of the issue you're addressing. Propose discovered improvements as separate issues.

## Before You Open a PR

Run the validators locally — CI runs the same checks on every push:

```bash
node scripts/validate-agents.mjs           # frontmatter, references, code fences, plugin.json agent list
node scripts/generate-codex-agents.mjs --check   # Codex mirror sync
bash scripts/bump-version.sh --check       # manifest versions in sync
claude plugin validate .                   # manifest and component schema (needs the claude CLI; not run in CI)
```

`claude plugin validate` reports frontmatter warnings for the non-agent files under `agents/` (`README.md`, `README.ko.md`, `_shared/*.md`). It scans `agents/` regardless of the manifest's `agents` list; the runtime loads only the listed files, so those warnings are expected. It also notes that the root `CLAUDE.md` is not loaded as plugin context — that file is the working agreement for contributors to this repository, not plugin content.

For skill changes, run the relevant behavioral tests (they invoke the `claude` CLI and consume tokens):

```bash
claude plugin eval . --tag trigger --ablation none --trust-plugin --no-publish   # see evals/README.md for every tag
```

See `tests/README.md` for the full suite catalog (deterministic, LLM-behavioral, pressure scenarios).

## What Makes a Good PR

- One concern per PR — a skill fix, an agent contract change, or a doc sync, not all three
- For skill changes: include before/after behavioral evidence (pressure-test transcripts or run counts)
- For agent changes: keep the I/O contract table in `AGENTS.md` in sync
- For routing changes: `skills/using-my-poor-ai/router.md` (injected at session start) and `skills/using-my-poor-ai/SKILL.md` both state the routing rules — change them in the same PR
- Passing CI (`validate` workflow)

## Adding New Components

| Component | Guide |
| --------- | ----- |
| New skill | `skills/writing-skills/SKILL.md` — includes testing methodology |
| New agent | `AGENTS.md` "공통 에이전트 작성 규칙" — the invariants there (single responsibility, relative paths, no hardcoding, validator pass, Codex mirror regeneration, registration in the `agents` list of `.claude-plugin/plugin.json`, …) |
| New command | Add `skills/<name>/SKILL.md` (frontmatter `name` must match the directory) and register it in `skills/README.md` and `skills/README.ko.md` |

## Code of Conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md).
