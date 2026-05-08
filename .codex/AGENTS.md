# SDP Codex Instructions

You are operating in a repository that uses Spec-Driven Protocol (SDP).
SDP is a structured workflow for AI-assisted software development:
explicit scope, workstreams, quality gates, review, and evidence before ship.

## Quick Start

Read these in order:

1. `AGENTS.md`
2. `docs/reference/project-map.md`
3. `prompts/commands.yml`
4. `docs/reference/FALLBACK_MODE.md` if your Codex runtime cannot spawn subagents

## Main Commands

### Planning and analysis

- `@vision` — strategic product shaping
- `@feature` — feature planning
- `@idea` — requirements gathering
- `@design` — workstream design
- `@understand`, `@scout`, `@architect`, `@reality`, `@metrics` — repo analysis

### Execution

- `@build` — execute one scoped workstream
- `@oneshot` — end-to-end feature execution
- `@operate` / `@deploy` — release and operations work

### Bugs and review

- `@fix`, `@bugfix`, `@hotfix`, `@issue`, `@debug`
- `@review`, `@verify-workstream`, `@ci-triage`

### Coordination

- `@llm-council` — multi-model synthesis for hard decisions
- `@git-worktree` — safe parallel work setup
- `@parallel-dispatch` — parallel subagent delegation

## Quality Gates

Run the relevant gates before claiming a task is complete:

| Language | Build | Test | Lint |
|---|---|---|---|
| Go | `go build ./...` | `go test ./...` | `go vet ./...` |
| Python | `pip install .` | `pytest` | `ruff check .` |
| Node.js | `npm run build` | `npm test` | `npm run lint` |
| Rust | `cargo build` | `cargo test` | `cargo clippy` |
| Java | `mvn compile` | `mvn test` | `mvn checkstyle:check` |

## Operating Rules

- No code change without a clear scope.
- Prefer TDD for behavior changes.
- Do not hide broken assumptions. Call them out and resolve them.
- Use `prompts/commands.yml` as the canonical command mapping.
- Use `prompts/skills/` as the canonical skill source.

## Landing The Plane

Before ending a session:

1. Run the relevant quality gates.
2. Verify acceptance criteria with evidence.
3. Update docs if behavior or UX changed.
4. Commit and push from a harness that has git access if your Codex sandbox does not.

## Related Files

- `prompts/commands.yml`
- `prompts/skills/`
- `prompts/agents/`
- `docs/reference/FALLBACK_MODE.md`
