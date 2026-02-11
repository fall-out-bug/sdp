# SDP Plugin — Go CLI

**Go binary for the Spec-Driven Protocol. Provides CLI commands for health checks, workstream parsing, guard enforcement, and telemetry.**

## Build

```bash
cd sdp-plugin
go build -o sdp ./cmd/sdp
```

## Commands

```bash
sdp doctor                    # Health check (hooks, config, deps)
sdp status                    # Show project state
sdp init                      # Initialize SDP in a new project
sdp guard activate 00-001-01  # Enforce edit scope
sdp guard check <file>        # Verify file is in scope
sdp parse <ws-file>           # Parse workstream file
sdp verify <ws-id>            # Verify workstream completion
sdp tdd <ws-id>               # Run TDD cycle
sdp telemetry status          # Telemetry status
```

## Skills and Agents

Prompt-based skills and agents are in the parent directory:
- `../.claude/skills/` — 24 workflow skills
- `../.claude/agents/` — 23 multi-agent definitions
- `../.cursor/` — Cursor IDE integration
- `../.opencode/` — OpenCode integration

Prompt copies for the Go plugin:
- `prompts/skills/` — Skill prompts
- `prompts/agents/` — Agent prompts
- `prompts/validators/` — Quality validators

## Telemetry

SDP collects **opt-in anonymized telemetry** stored locally:

- Command invocations and duration
- Success/failure rates and quality gate results
- **No PII**, no file paths, no code content
- Local only (`~/.sdp/telemetry.jsonl`), auto-cleanup after 90 days

```bash
sdp telemetry status   # Check
sdp telemetry enable   # Opt-in
sdp telemetry disable  # Opt-out
sdp telemetry clear    # Delete data
```

## Language Support

| Language | Tests | Coverage | Type Check | Lint |
|----------|-------|----------|------------|------|
| Python | pytest | pytest-cov | mypy | ruff |
| Go | go test | go tool cover | go vet | golangci-lint |
| Java | Maven/Gradle | JaCoCo | javac | checkstyle |
| TypeScript | jest/vitest | c8/istanbul | tsc | eslint |

## License

MIT

## Version

0.9.0
