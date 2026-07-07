# Provenance

Every line of Go code in this directory was written **unattended** by the my-poor-ai `subagent-driven-development` pipeline. No human edited the implementation.

## How it was generated

```bash
bash tests/subagent-driven-dev/run-test.sh go-fractals
```

The inputs were [design.md](design.md) and [plan.md](plan.md) (kept here verbatim). The pipeline scaffolded an empty git repo, then executed the plan task-by-task — one TDD cycle and one commit per task — and finished with a review pass that produced its own fix commit.

## Commit trail (from the pipeline run, 2026-07-07)

```
bf4b314 Fix: validate mandelbrot --char length for consistency with sierpinski
9069870 Task 10: README
39a5f9f Task 9: fix temp dir leak in TestMain
be6c7d6 Task 9: Integration Tests
ed95900 Task 8: scope SilenceUsage to validation errors only
4bcda7f Task 8: Input Validation and Error Handling
ce791aa Task 7: Character Set Configuration
4d6a93b Task 6: Mandelbrot CLI Integration
d7f14e1 Task 5: add Render()-level known-point tests
614ba67 Task 5: Mandelbrot Algorithm
6bdf9ef Task 4: Sierpinski CLI Integration
31de48e Task 3: fix blank rows for non-power-of-two sizes
2cd6093 Task 3: Sierpinski Algorithm
93c1b16 Task 2: Cobra root command with help listing subcommands
f0de11f Task 1: Project setup with go.mod, directory structure, and minimal main
e0736b6 Initial project setup with design and plan
```

Note the non-task commits: `31de48e`, `d7f14e1`, `39a5f9f`, `ed95900`, and the final `bf4b314` are fixes the pipeline's own test cycles and review phase caught — the discipline working as designed, not a hand-polished afterthought.

## Verification (run on the generated code, unmodified)

```
$ go test ./...
ok  github.com/my-poor-ai-test/fractals/cmd/fractals
ok  github.com/my-poor-ai-test/fractals/internal/cli
ok  github.com/my-poor-ai-test/fractals/internal/mandelbrot
ok  github.com/my-poor-ai-test/fractals/internal/sierpinski

$ go vet ./...   # clean
```

```
$ go run ./cmd/fractals sierpinski --size 16
               *
              * *
             *   *
            * * * *
           *       *
          * *     * *
         *   *   *   *
        * * * * * * * *
       *               *
      * *             * *
     *   *           *   *
    * * * *         * * * *
   *       *       *       *
  * *     * *     * *     * *
 *   *   *   *   *   *   *   *
* * * * * * * * * * * * * * * *
```

The `.git` directory of the pipeline run is not carried over (this repo has its own history); the log above is the record of it.
