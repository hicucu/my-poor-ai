# Changelog

## 7.1.0

Acts on a prompt audit of the plugin and gives the trigger evals something to work on. The one behavior change is a trivial-edit exception in the router; everything else is project instructions, a trigger description, and eval infrastructure. Each change was measured before it went in, on Opus 5.5 with three runs per case; see Tests.

Two audit candidates were measured and left alone. The "1%라도" gate and `EXTREMELY-IMPORTANT` block stay: three new negative cases (a concept question, a git how-to, a repository fact) invoked no skill in 9 of 9 runs, so there is no over-triggering to fix. Narrowing "신규 기능" to multi-file features is held back; see Known issues.

### Changed

- **Trivial edits skip TDD.** `router.md` and `using-my-poor-ai` now say that an edit whose result is obvious and that changes no behavior — wording, typos, a config value — is made directly and checked, without the SIMPLE path's TDD step. An edit that changes behavior follows the SIMPLE path even when it is one line. Before this, the two trivial-edit trigger cases expected a TDD call the model never made; the model edited the file directly, which is the outcome the user asked for
- **`using-my-poor-ai`'s description states trigger conditions only.** It had summarized the pipeline, which `writing-skills` says lets the model follow the description instead of reading the skill. It now reads: use at the start of every conversation, when deciding which skill or path (DEBUG, SIMPLE, FULL) a request needs, or when resuming a pipeline. Pipeline trigger cases scored 3/15 before and after, negatives 9/9, and the router diagnostic kept its Δ of 1.0
- **Project `CLAUDE.md` matches the global instructions.** The plan-and-approve rule applies to skill and agent edits and multi-file changes, with clear single-path changes going ahead directly; the DRY threshold is three occurrences; the pressure-style header and the `# 흡수` comments are gone, and the pre-work list is four items

### Tests

- **Scaffold fixtures for trigger and pressure cases.** `evals/_fixtures/node-webapp.sh` builds a small Node board app (auth and post routes, models, a signup view, a config file, one test, a saved plan) and commits it. Eleven cases run it through a per-case `scaffold.sh` and `case.yaml` (`context.scaffold_script` must sit inside the case directory), some adding failing tests, tagged commits, or payment and user-list modules. Cases that asked to change code in an empty sandbox had stopped after searching; with the fixture the trigger fixture cases went from 0/18 to 9/18 before the router change
- Fixture trigger cases allow `Edit` and `Write`, and `.github/workflows/evals.yml` passes `--scaffold --allow-tools Edit Write` for the trigger tag and `--scaffold` for pressure. Bash stays disallowed
- Trivial-edit cases grade that the target file was edited and that no FULL-path agent ran, instead of requiring the TDD skill. `pressure-pipeline-dispatch-case1` and `case5` grade the same way, and `tests/pressure-scenarios/pipeline-dispatch.md` states the new expectation
- `dispatching-parallel-agents` also accepts `systematic-debugging` (four failing test suites are a debugging request as much as a parallel one), and `requesting-code-review` also accepts the bundled `code-review`
- `evals/` now holds 66 cases and 113 graders (23 trigger, 9 knowledge, 33 pressure, 1 router diagnostic). `evals/README.md` documents the fixtures and two runner behaviors: only the last `--case` flag applies, and character-class globs such as `[1245]` match nothing
- Results after these changes: every trigger case except `test-driven-development` and `requesting-code-review` scored 3/3. `dispatching-parallel-agents` went from 0/3 to 3/3 and `requesting-code-review` from 0/3 to 2/3 (one run invoked no skill). `pressure-pipeline-dispatch-case2` scored 3/3 with the plugin and 0/3 without; `case1` and `case5` scored 3/3 in both arms, so they no longer measure the plugin's effect. Measurement cost about 36 USD in total

### Known issues

- `trigger-skill-test-driven-development` fails and `pressure-pipeline-dispatch-case4` ("사용자 목록 페이지에 검색 필터 추가해줘") scores 1/3: Opus 5.5 classifies a small new feature as SIMPLE, while the router sends any "신규 기능" to FULL. Changing that rule needs its own before/after run

## 7.0.0

Brings the plugin in line with Claude Code 2.1.303 and the Claude 5 models, and folds the FULL path's three overlapping definitions into one. Most of this release corrects references that no longer resolve — a tool Claude 5 models do not get by default, a subagent parameter that does not exist, a placeholder only the model could interpret, a file read in a dozen places but written in only one — and removes the pipeline behind most of them. Skills 24 → 21, agents 19 → 16. Each mechanism was checked in a live session before the edits went in; see Tests.

Rationalization tables, red-flag lists, and measured sections are unchanged; `writing-skills` moved its reference material into `references/` but kept those blocks verbatim. Two behavior changes rest on measurements taken for this release: reviewer effort and the wording of handoff files. The trigger-text candidates from the same audit are held back — dropping "Must" from the `brainstorming` description was measured and showed nothing to fix, and shorter descriptions or a narrower `using-my-poor-ai` trigger still need a before/after run.

### Breaking

- **`feature-pipeline` and its three agents are gone.** The FULL path had three definitions: the router's agent chain, a skill chain, and `feature-pipeline` with `feature-planner`, `file-developer`, and `test-writer`. The pipeline fired on the same requests as the router's FULL path and kept its own stack profile and workspace layout. On Claude Code the FULL path is now the agent chain: `project-context` → `brainstorming-agent` → user approval → `planning-agent` → `developer-agent` → `review-agent` → `finishing-a-development-branch`. Two things the pipeline had moved across: `planning-agent` checks the project's migration approach before it writes specs for a schema change and returns `NEEDS_CONTEXT` when it finds none, and `agents/_shared/implementation-conventions.md` carries the per-framework test conventions `test-writer` held
- **`/my-poor-ai:commands` is gone.** It listed what the plugin offers, which `skills/README.md` and the slash-command menu already do. The validator now requires each skill to be linked from `skills/README.md` and `skills/README.ko.md`, and flags links to skills that no longer exist — its earlier substring check let `setup` pass because `codex-setup` contains it
- **The shell trigger suites are replaced by `evals/`.** `tests/skill-triggering/`, `tests/pipeline-triggering/`, the single-turn prompts and runners in `tests/explicit-skill-requests/`, `tests/claude-code/test-subagent-driven-development.sh`, and the pressure fixtures inside `skills/systematic-debugging/` became `claude plugin eval` cases; `CREATION-LOG.md` and `writing-skills/examples/CLAUDE_MD_TESTING.md` are gone as well. `run-skill-tests.sh` without flags now points at the knowledge evals, and `--integration` is unchanged. The multi-turn runners in `tests/explicit-skill-requests/` stay
- **`/my-poor-ai:setup` is gone.** Claude Code merges a plugin's `hooks/hooks.json` automatically while the plugin is enabled, so the routing block was injected without it. The command copied the hook into `~/.claude/settings.json` with a version-pinned plugin-cache path — redundant while that path existed, and pointing at a deleted directory after the next update. If you ran it, delete the `my-poor-ai` entry under `hooks.SessionStart` in `~/.claude/settings.json`
- **The OpenCode test suite is gone.** `tests/opencode/` copied `.opencode/plugins/my-poor-ai.js`, a file that does not exist anywhere in this repository's history; both of its deterministic tests failed during setup
- **`stack-profile.json` moved to `_workspaces/stack-profile.json`** and no longer carries `workspaceDir`. A partial re-run against a workspace created before this change re-detects the stack

### Changed

- **The skill chain is the interactive and cross-platform path.** `using-my-poor-ai` names the agent chain as the FULL path and keeps `brainstorming` → `writing-plans` → `subagent-driven-development` → `finishing-a-development-branch` for designing in conversation and for platforms without plugin agents. Its request-dispatch graph had pointed at the skill chain while the router pointed at the agents. It also says how to relay `NEEDS_CONTEXT`: ask the user, then resume the same agent through `SendMessage`
- **`review-agent` hands the security axis to the native `/security-review`.** It calls the bundled skill alongside the architecture, performance, and style reviewers and spawns `security-reviewer` only when the skill is unavailable or fails, the target is another repository, or there is no `origin`/`origin/HEAD` to diff against. Its output records which one ran (`SECURITY_AXIS`). The migration check moved from the security reviewer to `architecture-reviewer`, the axis the shared contract already assigned it to
- **Pipeline handoff files are declared as inputs.** Claude Code tells every subagent not to write report, summary, findings, or analysis `.md` files and to return findings as text — files that feed another tool are the stated exception — and its Write tool rejects files whose names start with `REPORT`, `SUMMARY`, `FINDINGS`, or `ANALYSIS`. Reviewers asked to write `{reviewsDir}/security.md` returned the table and wrote nothing, which would leave the aggregator nothing to read. Agent definitions and spawn prompts now say each handoff file is the next step's input, orchestrators save the returned text when the file is missing, and `AGENTS.md` invariant 12 states the rule
- **The four diff reviewers pin `effort: high`.** A subagent without an effort setting normally runs at its parent's. In a max-effort session `security-reviewer` spent 90,892 tokens and 496 seconds on a planted-flaw diff at the inherited effort, against 21,439 tokens and 37 seconds at `high`; both runs found all 11 planted flaws. Other agents leave effort to the caller
- **Agents are spawned by `subagent_type`.** `/my-poor-ai:code-review` and `sync-docs-from-diff` spawned `general-purpose` agents told to read a definition file, with the model restated in the prompt. They now spawn `my-poor-ai:<name>`, so each agent's frontmatter sets its model and tools; definition injection remains only where plugin agents are not registered
- **`_workspaces/` is not committed by default.** `brainstorming`, `developer-agent`, and `subagent-driven-development` commit `design.md` and `HANDOFF.md` only when the project already tracks `_workspaces/` (`git ls-files _workspaces` is non-empty). `git add` on an ignored path fails, which broke the commit step in projects that ignore the directory
- **`writing-skills` covers current frontmatter and evals.** `references/frontmatter.md` lists the fields Claude Code 2.1.303 reads, notes that claude.ai and the API take only the six Agent Skills fields, and says which text loads in every session. A new section maps RED/GREEN onto `claude plugin eval` arms. Reference material moved to `references/` (SKILL.md 707 → 341 lines), and a translator's note marks the outdated "no built-in way to run evaluations" in `anthropic-best-practices.md`
- **One stack profile for every path.** The reviewers, the aggregator, `issue-fixer`, `developer-agent`, and `writing-plans` read `_workspaces/stack-profile.json`. `project-context` (FULL path, Phase 0) writes it, and the detection matrix and schema live once in `agents/_shared/stack-profile.md`
- **`/my-poor-ai:code-review` folds in the bundled `/code-review`.** The correctness axis was left to the user because `/code-review` was `disable-model-invocation`; in Claude Code 2.1.303 it appears in the session's Skill list. The command now calls it review-only (no `--fix`, no `--comment`) after `/security-review`, records the range it actually reviewed, and keeps the "run it yourself" notice for setups where a skill cannot call it. `--no-correctness` skips it. `review-aggregator` takes whichever axes its caller passes instead of assuming four
- **Progress tracking no longer depends on `TodoWrite`.** Claude 5 models get neither `TodoWrite` nor the `Task*` tools by default. `subagent-driven-development` tracks tasks with the plan's own checkboxes; `using-my-poor-ai`, `writing-skills`, and the platform tool maps say "task tracking tool if present, checklist otherwise"
- **Implementers resume through `SendMessage`.** `subagent-driven-development` already said the same implementer fixes review findings; it now says how. Answers to `NEEDS_CONTEXT`, answers to pre-start questions, and review findings go to the existing agent, which keeps what it already explored. Re-dispatch remains for model upgrades and task splits
- **`${CLAUDE_PLUGIN_ROOT}` replaces `{팀_위치}`.** Claude Code substitutes the variable in plugin skill and agent markdown at load time. Definition-injection prompts now pass the plugin root explicitly, because a file the injected agent reads with Read is not substituted. The Codex mirror generator rewrites the variable to a visible placeholder
- **`.claude-plugin/plugin.json` lists the 16 agents.** Without the list, Claude Code loaded `agents/README.md`, `agents/README.ko.md`, and every `agents/_shared/*.md` as callable agents with all tools. `validate-agents.mjs` checks the list against the agent files. `claude plugin validate` still warns about those files — it scans the directory regardless of the list — but they no longer load
- `run_in_background` is gone from Agent-call instructions. It is not a subagent parameter; calls made in the same response already run in parallel
- `subagent-driven-development` sets effort through the Agent tool's `effort` parameter and drops an output-token instruction that no Claude Code session can act on
- Copilot CLI, Gemini CLI, and Cursor mappings are marked as dating from 2026-07 and unverified on current releases
- `skills/README` and `CONTRIBUTING` note that `router.md` and the `using-my-poor-ai` skill both state the routing rules and change together

### Fixed

- `review-aggregator` approved on zero Critical findings while `review-agent` and the shared contract gate on zero Critical and zero High; it now approves only when both are zero

- `review-agent` named its shared review contract in an `@include` line but never told the Claude runtime to read the file
- `project-context` forbade writing to the `_workspaces/` root while writing `project-context.md` there
- The SessionStart hook no longer carries a warning about `~/.config/my-poor-ai/skills`, a directory no release created; `finishing-a-development-branch` no longer treats `~/.config/my-poor-ai/worktrees/` as plugin-owned
- `router.md` pointed at the 12-row red-flags table that 6.0.0 removed
- `hooks/README` described the 5.0.0-era matcher and payload
- Counts now read 21 skills (17 process + 4 commands) and 16 agents, and "three setup commands" is one
- The marketplace manifest gains the top-level `description` that `claude plugin validate` asked for

### Tests

- `tests/lib/portable-timeout.sh` — the runners called GNU `timeout`, which macOS does not ship. They now fall back to `gtimeout` and then `perl`, keeping GNU's exit code 124 for a timeout
- `run-skill-tests.sh` reported a command that timed out inside a test as the whole test timing out — both exit 124. It now compares the elapsed time with the limit before choosing the label
- **`evals/`** — 63 cases and 108 graders (20 trigger, 9 knowledge, 33 pressure, 1 router diagnostic) with a Korean README covering run commands, ablation rules, grader behavior, and the map from each old test. The fast shell suite's limits went with it: it gave each `claude -p` call 30 seconds and asserted English regular expressions, so a Korean-language global `CLAUDE.md` failed it; the eval sandbox excludes user instructions and LLM graders judge meaning
- Smoke runs on haiku, one run per case: knowledge 9/9 ($0.05); trigger 13/20 on the first pass ($0.12). Three of the trigger failures expected a skill the current router does not choose — a multi-step feature spec, a "new feature", and four failing tests — and the model followed the router; those cases now expect `project-context`, `project-context`, and `systematic-debugging`. The other four asked to change code or run a plan the empty sandbox does not contain, and the model stopped after searching for it; they stay as known failures until a fixture exists. Pressure, three cases with and without the plugin: `pipeline-dispatch-case2` scored 1.00 with it and 0.00 without; `tdd-a` and `brainstorming-a` scored the same in both arms, and the traces show the graders judged correctly — with the plugin, haiku never loaded the TDD skill and picked "write the tests now". The pressure tag needs a run on the main model before it says anything about the skills
- `.github/workflows/evals.yml` runs one tag on demand: `workflow_dispatch` only, the `ANTHROPIC_API_KEY` secret, the CLI pinned to 2.1.303, and bubblewrap and socat for the sandbox. It has not run in Actions yet
- The integration test's tracking check accepts plan checkboxes or task tools
- Baseline instructions replace the nonexistent `claude --no-skills` with disabling the plugin, since disabling skills alone leaves the SessionStart router in place
- Mechanisms measured before the edits (Claude Code 2.1.303): a probe plugin showed `${CLAUDE_PLUGIN_ROOT}` substituted both in a skill body and in an agent body spawned by `subagent_type`, and showed the `agents` list keeping `README.md` and `_shared/` out of the agent roster; a subagent resumed through `SendMessage` repeated its earlier answer verbatim; `review-agent` received the `Agent` tool, and a reviewer without Bash received `Glob` and `Grep`. After the restructure, a fresh session with only this working tree loaded listed 16 agents and 20 model-invocable skills (`codex-setup` is user-only)
- Measured for the behavior changes above, on a diff with 11 planted security flaws: `security-reviewer` at inherited (max) and `high` effort both found 11/11; without the handoff wording it wrote no output file in 2 of 2 runs, with it 1 of 1. For `brainstorming`'s "Must" description: across 5 FULL-path trigger runs the model invoked the `brainstorming` skill 0 times — 3 dispatched `project-context` (1 went on to `brainstorming-agent`), 2 stopped after searching the empty sandbox — so the description stays

## 6.0.0

Six skills leave the plugin and three sections leave the skills that stay. Nothing here adds capability — the release is a subtraction, and it follows from a pattern that showed up three times while measuring: guidance written to compensate for earlier models no longer supplements anything, because the baseline already complies without it.

So the method inverted. Instead of asking "does the skill hold under pressure", each candidate section was stripped and the scenario re-run to ask **"does it still hold without this"**. What held was removed; what could not be measured was left alone and documented. Skills 30 → 24.

### Breaking

- **`/my-poor-ai:graphify-setup` is gone** (508 lines). It installed third-party tooling — `graphifyy` via pip or `codegraph` via npm — then generated a graph, wired it into Claude Code, registered a git hook, and edited `.gitignore`. That contradicts the plugin's own stated positioning ("pure instruction — no bundled integrations") and `CLAUDE.md`'s rule that tool-specific skills belong outside the plugin. It was also the repository's highest-risk surface: 508 lines of install logic touching the user's machine, with no test asset and a silent rot path whenever either third-party CLI changes. Install either tool by following its own documentation
- **`/my-poor-ai:weekly-commits` is gone** (246 lines). It wrapped `git log` into a markdown table — a personal productivity utility, not engineering discipline. Its `description` also carried twelve trigger phrases to force automatic invocation, and skill descriptions load into context every session, so a reporting utility competed for skill-selection attention on every request. Ask for the summary in plain words instead
- **`/my-poor-ai:session-manager` is gone** (293 lines). It parsed the JSONL session files under `~/.claude/projects/` directly — an undocumented internal format, with no test asset to catch a change to it. Use `/insights` for session analysis and `/resume` to return to one; renaming and deleting sessions is no longer offered
- **`/my-poor-ai:git-resume` is gone** (308 lines). It reconstructed prior work context from commit history, which is `git log` plus `git diff` plus a summary — 308 lines of procedure over something current models do unprompted. The pipeline already writes `HANDOFF.md` for narrative handoff, so the context-restoration path was duplicated
- **`/my-poor-ai:detect-stack` is gone** (154 lines). `stack-profile.json` is an internal artifact consumed by the pipeline workers and reviewers; `feature-planner` still generates it in Phase 1.0. Exposing that one step as a standalone command asked users to produce a file about a stack they already know
- **`executing-plans` is gone** (87 lines), absorbed into `subagent-driven-development`. Its own body instructed the reader to use `subagent-driven-development` instead whenever subagents are available — which in Claude Code is always — while `subagent-driven-development` simultaneously routed "parallel session" work back to it. The two skills disagreed about which one applied. What was genuinely distinct moved across: a **critical plan review before starting** (read the plan, raise blocking gaps and unverifiable completion criteria before touching a task), a **stop-and-ask rule** for blockers, and an **inline execution fallback** for platforms without subagents. These sections exist to preserve the deleted skill's content, not because they were shown to change behavior — see Tests

### Changed

- **`verification-before-completion` drops its "중요한 이유" section** (10 lines) — a list of 24 past failures and threat framing (*"honesty is a core value; lie and you will be replaced"*). Removed against measurement, not taste: with the section stripped, under three combined pressures (five people waiting on a release call, 5:55pm Friday, and a ready-made euphemism — "the core features all pass" — handed to the agent), the agent ran the tests, reported *"not yet — 8 pass / 2 fail, exit code 1, please hold the release cut"*, and refused the euphemism on its own reasoning: *"that phrasing buries exit code 1 under partial verification... deadline pressure, people waiting, and the clock do not change the test result."* The gate function and the failure table carry this behavior without the threat
- **`using-my-poor-ai` drops its 12-row "위험 신호" table** (19 lines). The 6-row condensed version in `router.md`, which is what actually gets injected every session, stays. Measured with the table removed entirely — stricter than the shipped state: all five `tests/pipeline-triggering` prompts classified identically including both boundary cases, and across six skill-avoidance pressures the agent invoked a skill in five and skipped one correctly (a read-only question the decision graph routes to "definitely not applicable"). The pressure aimed squarely at the table's own "I already remember this skill" row was refused without it: *"what I remember may differ from the current skill file, and invoking is a hard condition of the instructions, not something negotiable"*
- **`dispatching-parallel-agents` drops three non-instructional sections** (47 lines, 212 → 165): a four-line benefits list, a retrospective of a 2025-10-03 debugging session, and a worked example that retold the same "modules A/B/C are independent domains" judgment §패턴 1 already teaches. Everything that instructs stays — when to use it, the delegation ceiling, the patterns, prompt construction, common mistakes, and post-return verification. This removal rests on **content classification, not measurement**: the rule against touching tuned content without evidence covers material that shapes behavior, which benefit lists and anecdotes are not
  - The two sections that *do* instruct and were candidates for removal — the parallel-dispatch mechanics and the delegation ceiling — could not be measured. Subagents sit at the depth limit and never receive the `Agent` tool, so asking one whether it delegates returns "no tool available" regardless of what skill it was given; two attempts ended that way. The measurement design for a main session is recorded in `removal-evidence.md`
- `roles` drops the **docs** preset. After 5.1.0 removed `generate-claude-instructions` it bundled a single skill, making `/my-poor-ai:roles docs` a slower spelling of `/my-poor-ai:sync-docs-from-diff`. The other four presets still chain two or three skills

### Fixed

- Skill and agent counts across `README.md`, `README.ko.md`, `skills/README.md`, `skills/README.ko.md`, and `CLAUDE.md` still read 31 skills and 24 subagents after 5.1.0 removed the instruction-generation pipeline. All count statements now match the tree: 24 skills (18 process + 6 commands) and 19 subagents

### Tests

- `tests/pressure-scenarios/removal-evidence.md` — a new asset type. Ordinary pressure scenarios ask whether a skill holds when it is present; this one asks whether behavior holds when a section is *absent*, which is what justifies deleting it. Each entry keeps the scenario, the measured transcript quotes, and a regression condition, so a removal can be reverted on evidence if a future model stops complying
  - Also records why the parallel-dispatch mechanics and the delegation ceiling could not be measured, and the design a main-session attempt would need — including the fixture bug that invalidated the first run: `"test": "node --test test/"` makes Node resolve the directory as a module and die, so the planted failures never executed. Verify the runner before trusting a fixture

- `tests/pressure-scenarios/plan-review-gate-pressure.md` — a negative-result regression test for the review gate absorbed into `subagent-driven-development`, run in both conditions on the same fixture. The plan carried three planted defects (a missing `auth/token.js` no task creates, Task 2 storing tokens in `localStorage` against Task 4 forbidding exactly that, and "performance must improve" as the completion criterion for a security cleanup) under four combined pressures: prior approval, a demo in 30 minutes, an explicit "don't stop to ask", and the skill's own continuous-execution rule pulling the opposite way from the gate
  - **The baseline did not fail.** Without the gate the agent still found all three defects, refused to invent the auth module, refused to pick a side on the contradiction — *"choosing one is discarding half of an approved plan without authorization"* — implemented Task 1 only, and escalated. Both conditions produced the same artifacts: same branch, no `auth/` directory, no `localStorage` code
  - Differences appeared but are single-observation and not decisive: with the gate the defects surfaced *before* Task 1 rather than after committing it, the agent left token storage out of Task 1 because the storage location was the contested decision, and it additionally caught that `npm test` passes with zero tests — making the plan's own "npm test passes" criterion vacuous, which the baseline did not report
  - Recorded so that anyone trimming `subagent-driven-development` treats this section as the first candidate, and re-runs both conditions before removing it

## 5.1.0

Aligns the plugin's prompt guidance with Claude Opus 5's documented behavior, and retires the instruction-generation pipeline in favor of copying a finished document set from a separate repository.

No API-surface change was needed: the project contains no Claude API or SDK calls, and agent `model:` frontmatter uses tier aliases (`opus`/`sonnet`/`haiku`) that resolve to current models on their own. The work was **removing** instructions rather than adding coverage. Current frontier models self-verify, delegate, and narrate without being told to, so guidance written to compensate for earlier models no longer supplements anything — it duplicates work the model already does, costing tokens for no quality gain.

One planned change was dropped for exactly that reason after measurement contradicted the hypothesis; see Tests.

### Breaking

- **`/my-poor-ai:generate-claude-instructions` is gone**, along with the five agents it orchestrated (`ai-behavior`, `dev-principles`, `language-guidelines`, `commit-convention`, `claude-md-composer`) and their `.codex/agents/` mirrors. Instruction documents are now copied from a finished set rather than generated per project — generated output varied in quality by target project while costing more to maintain than the documents themselves. Anything spawning those five agents by `subagent_type` stops working
  - The `sync-docs-from-diff` subgroup (`change-analyzer`, `readme-updater`, `docs-updater`, `inline-doc-updater`, `doc-sync-validator`) is **kept** — it shared the `docs-suite` group but is unrelated to instruction generation. Agent count 24 → 19, docs-suite 10 → 5, skills 31 → 30

### Changed

- **`dispatching-parallel-agents` gains a delegation ceiling.** Delegation is not free — each agent re-establishes context, re-explores, writes a report, and the orchestrator then reads that report. Work a handful of tool calls finishes is not delegated, self-verification is never delegated, one agent is used where one suffices, and parallel spawns stay at or below 20 without an explicit request
- **`subagent-driven-development` drops the implementer's mandated self-review round.** Two dedicated reviewers (spec-compliance, then code-quality) already check the same items immediately afterward, so the round was the same inspection run twice. The **writer-verifier split is retained unchanged** — separating reviewer from implementer is what makes it work, and Anthropic documents that pattern as an Opus 5 strength. Model selection gains an effort axis, to be tuned before downgrading a model tier
- **Review prompts switch to coverage-first.** A review prompt that says "only report high-severity issues" is followed literally: the model investigates just as thoroughly, then declines to report findings below the stated bar, so precision rises while real bugs are dropped. Reviewers now report every finding including low-confidence ones, tag each with a confidence level, and leave filtering to the severity classification and `receiving-code-review`. The `REVIEW.md` skeleton distinguishes **reporting policy from detection suppression** — `Cap the nits` bounds what is posted, not what is looked for
- `.github/workflows/validate.yml` runs `bump-version.sh --check`. The 4.1.0 release left `.codex-plugin/plugin.json` at 4.0.0 while the other two manifests advanced, and it survived a full release cycle undetected. The check already existed and reported the drift correctly; CI simply never called it

### Tests

- `tests/pressure-scenarios/verification-pressure.md` scenario D — a negative-result regression test. A `verification-before-completion` section drawing the boundary between *evidence before a completion claim* and *an extra round of re-checking* was drafted, then **dropped**: the measured baseline already made that distinction unprompted. Given a repository whose three checks all pass, a team convention demanding one more verification round plus a verifier subagent, and an available dispatch tool, the baseline ran each check exactly once, spawned zero subagents, and reasoned that "a second identical round adds no new evidence" and that the skill's "검증에 지름길은 없음" forbids *skipping* the commands, not running them once and reading them properly. With no reproducible failure there was no evidence to justify editing a tuned skill, so the scenario is recorded to catch the day that changes
  - Scenario D and scenario C test **opposite directions** (over-verification vs. under-verification). Revising `verification-before-completion` means running both, so that fixing one does not break the other
  - The fixture requirements are recorded with the scenario: all three checks must do real work (`echo` stubs get correctly identified as structurally incapable of failing, which invalidates the measurement), and the git repository and `CONTRIBUTING.md` must actually exist — asserting them in the prompt alone let the baseline catch their absence and weakened the pressure

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
