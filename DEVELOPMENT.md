# Development Guide

Get started with SDP development in under 15 minutes.

---

## Prerequisites

| Tool | Version | How to Check |
|------|---------|--------------|
| **Go** | 1.26+ | `go version` (see `.go-version`) |
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
./sdp --version
```

Expected time: **5 minutes**

---

## Project Structure

```
sdp/
├── sdp-plugin/           # Go CLI implementation
│   ├── cmd/sdp/          # CLI entry point
│   └── internal/         # Core logic
├── prompts/              # Canonical prompt, skill, and agent definitions
│   ├── skills/           # AI skill prompts
│   ├── commands/         # Harness command adapters
│   └── agents/           # Multi-agent definitions
├── docs/                 # Onboarding, protocol, and reference docs
├── .claude/              # Claude adapter
├── .cursor/              # Cursor adapter
├── .opencode/            # OpenCode adapter
└── .codex/               # Codex adapter
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
golangci-lint run --enable-only errcheck ./...

# Modern Go audit snapshot
golangci-lint run --enable-only modernize --issues-exit-code 0 ./...
```

## Go Style

For Go code, follow `@go-modern` and prefer behavior-preserving stdlib idioms that match the module Go version.

- Prefer `slices.SortFunc` over `sort.Slice`
- Prefer `strings.Cut` over `strings.SplitN(..., 2)` or manual `strings.Index` slicing
- Prefer `strings.CutPrefix` and `strings.CutSuffix` over prefix or suffix checks plus trim
- Prefer `slices.Contains`, `maps.Copy`, and `maps.Clone` over handwritten helper loops
- Prefer `any` over `interface{}` for locals, parameters, and variadics when behavior does not change
- Use `golangci-lint` or `staticcheck` instead of `golint`

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
2. Set Go SDK to 1.26
3. Enable golangci-lint in Settings → Go → Linter

---

## Using SDP During Development

SDP currently has two working surfaces:

- **Local Mode:** `sdp init`, `sdp doctor`, `sdp plan`, `sdp apply`, `sdp verify`, `sdp status`, `sdp next`
- **Prompt surfaces:** harness-native commands installed through `.claude/`, `.cursor/`, `.opencode/`, or `.codex/`

For most contributors, the simplest local check is:

```bash
cd sdp-plugin
go run ./cmd/sdp init --help
go run ./cmd/sdp doctor --help
go run ./cmd/sdp plan --help
```

If you use prompt surfaces while developing:

- edit canonical prompt source in `prompts/`
- treat Beads-backed queue flows as advanced, not required
- remember that `sdp deploy` only records approval after merge

---

## Task Tracking with Beads

```bash
# Find available work
bd ready

# Start working on issue
bd update sdp-xxx --status in_progress

# Close when done
bd close sdp-xxx
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
