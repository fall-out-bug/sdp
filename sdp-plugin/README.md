# SDP Plugin - Go CLI

Go binary for the Spec-Driven Protocol. It provides CLI commands for planning, execution, evidence, telemetry, quality gates, diagnostics, and long-running orchestration.

## Build

```bash
cd sdp-plugin
CGO_ENABLED=0 go build -o sdp ./cmd/sdp
```

## Commands

Use `sdp --help` for the full tree. The main command groups are:

- `sdp init`, `sdp doctor`, `sdp health`, `sdp status`, `sdp next`, `sdp demo`
- `sdp parse`, `sdp plan`, `sdp build`, `sdp apply`, `sdp orchestrate`, `sdp verify`, `sdp tdd`, `sdp deploy`
- `sdp guard ...`, `sdp quality ...`, `sdp drift detect`, `sdp diagnose`, `sdp watch`
- `sdp log ...`, `sdp decisions ...`, `sdp checkpoint ...`, `sdp coordination ...`
- `sdp telemetry ...`, `sdp metrics ...`, `sdp memory ...`, `sdp beads ...`, `sdp session ...`, `sdp task create`

See `../docs/CLI_REFERENCE.md` for the current command map.

## Core Libraries

The CLI is backed by reusable Go packages that map to user-visible subsystems:

- `internal/evidence` - evidence event builders, reader/writer chain, export, filtering, tracing
- `internal/telemetry` - opt-in local telemetry collection, analysis, export, upload packaging
- `internal/quality` - coverage, complexity, file size, and type checks for Python, Go, and Java
- `internal/watcher` - file watching, debounce, include/exclude filtering, quality violation tracking
- `internal/doctor` - environment checks, deep diagnostics, workstream dependency checks, config sanity
- `internal/orchestrator`, `internal/checkpoint`, `internal/session`, `internal/worktree` - long-running feature execution, resume, and worktree/session state

See `docs/reference/LIBRARY_CAPABILITIES.md` for the package capability map.

## Skills and Agents

Prompt-based skills and agents are in the parent directory:
- `../.claude/skills/` - workflow skills
- `../.claude/agents/` - multi-agent definitions
- `../.cursor/` - Cursor IDE integration
- `../.opencode/` - OpenCode integration

Prompt copies for the Go plugin:
- `prompts/skills/` - Skill prompts
- `prompts/agents/` - Agent prompts
- `prompts/validators/` - Quality validators

## Telemetry

SDP collects opt-in anonymized telemetry stored locally:

- Command invocations and duration
- Success/failure rates and quality gate results
- No PII, no file paths, no code content
- Local only (`~/.config/sdp/telemetry.jsonl` with config in `~/.config/sdp/telemetry.json`), auto-cleanup after 90 days

```bash
sdp telemetry status          # Check consent, event count, path
sdp telemetry consent         # Show consent options
sdp telemetry enable          # Opt-in
sdp telemetry disable         # Opt-out
sdp telemetry analyze         # Summarize local telemetry
sdp telemetry export json     # Export telemetry_export.json
sdp telemetry upload --format json
```

## Language Support

| Language | Tests | Coverage | Type Check | Lint |
|----------|-------|----------|------------|------|
| Python | pytest | pytest-cov | mypy | ruff |
| Go | go test | go tool cover | go vet | golangci-lint |
| Java | Maven/Gradle | JaCoCo | javac | checkstyle |
| TypeScript | jest/vitest | c8/istanbul | tsc | eslint |

## Documentation

- `../docs/CLI_REFERENCE.md` - current CLI surface
- `docs/reference/LIBRARY_CAPABILITIES.md` - package and subsystem capability map
- `docs/TELEMETRY_HOWTO.md` - telemetry behavior, paths, export/upload workflow
- `docs/quality-gates.md` - language-specific quality gate behavior

## License

MIT

## Version

0.9.0
