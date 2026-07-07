# Provenance

This example makes a stronger claim than [go-fractals](../go-fractals/PROVENANCE.md): here the pipeline was **not** given a design or a plan. The only human input was [REQUIREMENTS.md](REQUIREMENTS.md) — seven feature bullets. **The design document itself was authored by the pipeline's brainstorming phase**, then decomposed, implemented, and reviewed, all unattended.

## How it was generated

A fresh git repo containing only `REQUIREMENTS.md` was created, then a single headless prompt ran the FULL pipeline: brainstorming → design.md → planning (task specs + file manifest) → TDD implementation per spec → 4-way parallel review → issue fixes. Approval gates were self-reviewed (unattended CI mode); in interactive use a human approves at each gate.

## Pipeline-authored artifacts (preserved verbatim)

[pipeline-artifacts/](pipeline-artifacts/) contains what the pipeline wrote for itself:

- `design.md` — the design document produced by the brainstorming phase (no human wrote or edited this)
- `specs/` + `file-manifest.json` — the planning phase's task decomposition
- `HANDOFF.md` — session-handoff narrative maintained during implementation
- `review-report.md` — the aggregated 4-reviewer report that drove the final fix commit

## Commit trail (from the pipeline run, 2026-07-07)

```
9d3f4d1 fix: review 이슈 수정 - localStorage 로드 스키마 검증 및 스토어 수명주기 정리
0dcad9d feat: wire App integrating store and all todo components
a17b58a feat: TodoList component rendering items with empty state
bea81ca feat: TodoFooter with remaining count and clear-completed
47b9056 feat: TodoFilter component with All/Active/Completed
9343762 feat: TodoItem component with toggle and delete
6661580 feat: TodoInput component with add-on-enter and blank guard
06578b0 feat: reactive todo store with localStorage persistence
9df84e6 chore: scaffold Vite Svelte5 TS project with Vitest test pipeline
bf830dc Raw requirement only — no design, no plan (human input ends here)
```

The final commit (`9d3f4d1`) was produced by the review phase — the 4-reviewer report flagged a localStorage schema-validation gap and store lifecycle issue, and the issue-fixer resolved them. Discipline working, not decoration.

## Verification (run on the generated code, unmodified)

```
$ npm test
Test Files  8 passed (8)
     Tests  30 passed (30)

$ npm run build
✓ built in 230ms   (dist/assets/index-*.js 42.47 kB │ gzip: 16.41 kB)
```

The `.git` of the pipeline run is not carried over; the log above is the record of it.
