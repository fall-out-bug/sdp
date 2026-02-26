# Development Guide

Get started with SDP development in under 15 minutes.

---

## Prerequisites

| Tool | Version | How to Check |
|------|---------|--------------|
| **Go** | 1.24+ | `go version` (see `.go-version`) |
| **Git** | 2.30+ | `git --version` |
| **golangci-lint** | latest | `golangci-lint version` (optional) |
| **Beads CLI** | latest | `bd --version` (for task tracking) |

### Installing Prerequisites

```bash
# Go (macOS)
brew install go

# golangci-lint (macOS)
brew install golangci-lint

# Beads CLI (macOS)
brew tap beads-dev/tap && brew install beads
```

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/fall-out-bug/sdp.git
cd sdp

# Build the CLI
cd sdp-plugin
CGO_ENABLED=0 go build -o sdp ./cmd/sdp

# Run tests
go test ./...

# Verify installation
./sdp version
```

Expected time: **5 minutes**

---

## Project Structure

```
sdp/
├── sdp-plugin/           # Go CLI implementation
│   ├── cmd/sdp/          # CLI entry point
│   └── internal/         # Core logic
├── prompts/              # Canonical skill/agent definitions
│   ├── skills/           # AI skill prompts
│   └── agents/           # Multi-agent definitions
├── docs/                 # Documentation
│   ├── PROTOCOL.md       # Core specification
│   └── workstreams/      # Feature planning
└── .claude/              # Claude Code integration
    ├── skills -> ../prompts/skills
    └── agents -> ../prompts/agents
```

---

## Running Tests

### All Tests

```bash
cd sdp-plugin
go test ./...
```

### With Coverage

```bash
go test -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out  # Summary
go tool cover -html=coverage.out  # Browser report
```

### Specific Test Suites

```bash
# Smoke tests (fast, critical path)
go test -run TestSmoke ./...

# Unit tests only
go test -short ./...

# Integration tests
go test -run Integration ./...

# Single package
go test -v ./internal/guard/...
```

### Coverage Threshold

CI enforces **80% coverage**. Local runs display warnings only.

```bash
# Check coverage locally
go tool cover -func=coverage.out | grep total
```

---

## Linting

```bash
# Run all linters
golangci-lint run ./...

# With auto-fix
golangci-lint run --fix ./...

# Specific linter
golangci-lint run --disable-all --enable=errcheck ./...
```

### Pre-commit Checks

```bash
# Run guard checks (same as CI)
./sdp guard check

# Check staged files only
./sdp guard check --staged
```

---

## IDE Setup

### VS Code

1. Install [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)
2. Install Go tools: `Cmd+Shift+P` → `Go: Install/Update Tools`
3. Enable gopls language server

Recommended `settings.json`:

```json
{
  "go.testFlags": ["-v"],
  "go.coverOnSave": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package"
}
```

### GoLand

1. Open `sdp-plugin/` as Go module
2. Set Go SDK to 1.24
3. Enable golangci-lint in Settings → Go → Linter

---

## Using SDP Skills

SDP includes AI-powered skills for common workflows. See [CLAUDE.md](CLAUDE.md) for full usage.

### Example Workflow

```bash
# 1. Create a new feature
@feature "Add user authentication"

# 2. Design workstreams
@design idea-auth

# 3. Execute workstreams
@build 00-001-01

# 4. Review quality
@review F01

# 5. Deploy
@deploy F01
```

### Available Skills

| Skill | Purpose |
|-------|---------|
| `@feature` | Plan new feature |
| `@build` | Execute workstream |
| `@review` | Quality review |
| `@deploy` | Merge to production |
| `@debug` | Debug issues |
| `@issue` | Route bugs |

---

## Task Tracking with Beads

```bash
# Find available work
bd ready

# Start working on issue
bd update sdp-xxx --status=in_progress

# Close when done
bd close sdp-xxx

# Sync to remote
bd sync
```

---

## Troubleshooting

### Build Fails: "go: updates to go.mod needed"

```bash
go mod tidy
```

### Tests Fail: "permission denied"

```bash
chmod +x ./sdp
```

### golangci-lint Timeout

```bash
golangci-lint run --timeout=5m ./...
```

### Coverage Below 80%

Check uncovered lines:

```bash
go tool cover -html=coverage.out
```

---

## Next Steps

1. Read [PROTOCOL.md](docs/PROTOCOL.md) for specification
2. Review [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines
3. Check [CLAUDE.md](CLAUDE.md) for AI integration guide

---

**Version:** 1.0.0
**Updated:** 2026-02-13

---

## Repository Hardening Notes

### Changes Applied (2026-02-14)

| Area | Change |
|------|--------|
| Coverage threshold | 64% → 80% |
| Evidence logs | Removed from git tracking |
| Agent files | Added debugger.md, visionary.md, fixer.md |
| Git ignore | Added events.jsonl patterns |
| Quality gates | Aligned CI with local config |

### Evidence Log Policy

Evidence logs (`.sdp/log/events.jsonl`) are auto-generated session data:
- **Do not commit** - added to `.gitignore`
- **Retention**: 90 days
- **Max size**: 10MB with auto-compaction

### Dual Module Structure

SDP uses two Go modules (see ADR-001):
- `go.mod` (root) - Core modules
- `sdp-plugin/go.mod` - CLI implementation

Both use `github.com/fall-out-bug/sdp` module path. Work in `sdp-plugin/` for CLI development.

### Deferred Items

| Item | Reason |
|------|--------|
| go.work | Same module name conflict (ADR-001) |
| Release process | Requires version tagging workflow |
| src/sdp LOC violations | Legacy code, separate module |
