# Security Policy

## Reporting a Vulnerability

If you discover a security issue in my-poor-ai (e.g., a hook or script that could execute unintended commands, a prompt-injection vector in a skill/agent definition, or unsafe file operations in the test suites), please report it privately via [GitHub Security Advisories](https://github.com/hicucu/my-poor-ai/security/advisories/new) rather than opening a public issue.

Please include:

- The affected file(s) and the scenario in which the issue triggers
- A minimal reproduction (prompt, command, or environment)
- The impact you believe it has

You can expect an initial response within a week. Once fixed, the advisory will be published with credit unless you request otherwise.

## Scope Notes

- Hooks (`hooks/`) and scripts (`scripts/`) run on the installer's machine — issues here are highest priority.
- LLM-behavioral test runners use `--dangerously-skip-permissions` by design and are documented as such in `tests/README.md`; reports about that flag itself are out of scope, but escapes from the temp-directory sandbox are in scope.
- Skill/agent markdown is injected into LLM context; injection vectors that redirect agent behavior against the user's intent are in scope.
