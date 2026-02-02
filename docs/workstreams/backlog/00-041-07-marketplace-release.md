# 00-041-07: Marketplace Release

> **Feature:** F041 - Claude Plugin Distribution
> **Status:** backlog
> **Size:** SMALL
> **Created:** 2026-02-02

## Goal

Publish SDP plugin to Claude Plugin Marketplace and create migration guide for Python SDP users.

## Acceptance Criteria

- AC1: Plugin package validated against Claude schema
- AC2: Marketplace listing created (README, description, screenshots)
- AC3: Version 1.0.0 released and tagged
- AC4: Migration guide for Python SDP users (MIGRATION.md)
- AC5: Installation instructions tested on fresh project

## Scope

### Input Files
- `sdp-plugin/` (complete plugin package)
- Test results from WS-00-041-06
- Existing Python SDP documentation

### Output Files
- `sdp-plugin/plugin.json` (FINAL version)
- `sdp-plugin/README.md` (marketplace description)
- `sdp-plugin/MIGRATION.md` (migration guide)
- `sdp-plugin/CHANGELOG.md` (v1.0.0 release notes)
- Git tag v1.0.0

### Out of Scope
- Modifying plugin functionality (all previous WS complete)

## Implementation Steps

### Step 1: Finalize plugin.json

**File: sdp-plugin/plugin.json**

```json
{
  "name": "sdp",
  "version": "1.0.0",
  "displayName": "Spec-Driven Protocol (SDP)",
  "description": "Workstream-driven development with TDD, clean architecture, and quality gates for AI agents. Language-agnostic support for Python, Java, Go projects.",
  "author": "MSU AI Masters",
  "license": "MIT",
  "homepage": "https://github.com/ai-masters/sdp",
  "repository": "https://github.com/ai-masters/sdp-plugin",
  "categories": ["development", "testing", "documentation", "workflow"],
  "keywords": ["tdd", "clean-architecture", "workstream", "quality-gates", "ai-agents"],
  "permissions": [
    "Read(*)",
    "Write(*)",
    "Edit(*)",
    "Bash(git status, git log, git diff, git add, git commit)",
    "Bash(pytest, mvn test, go test)",
    "Grep(*)",
    "Glob(*)"
  ],
  "prompts": {
    "skills": "prompts/skills/*.md",
    "agents": "prompts/agents/*.md",
    "validators": "prompts/validators/*.md"
  },
  "binaries": {
    "optional": true,
    "url": "https://github.com/ai-masters/sdp/releases/latest/download/sdp-{platform}-{arch}",
    "sha256sum": "https://github.com/ai-masters/sdp/releases/latest/download/sdp-checksums.txt",
    "platforms": ["darwin-arm64", "darwin-amd64", "linux-amd64", "windows-amd64"]
  },
  "documentation": {
    "tutorial": "docs/TUTORIAL.md",
    "examples": "docs/examples/*/",
    "migration": "MIGRATION.md"
  }
}
```

Validate:
```bash
cat sdp-plugin/plugin.json | python -m json.tool
# Expected: Valid JSON, no syntax errors
```

### Step 2: Marketplace README

**File: sdp-plugin/README.md**

```markdown
# Spec-Driven Protocol (SDP) 🚀

**Workstream-driven development for AI agents with multi-language support.**

## Features

✅ **TDD Discipline** - Red → Green → Refactor cycle enforced by prompts
✅ **Clean Architecture** - Layer separation validated by AI
✅ **Quality Gates** - Coverage ≥80%, type safety, error handling
✅ **Multi-Language** - Python, Java, Go support
✅ **No Installation Required** - Prompts work standalone
✅ **Optional Binary** - Go CLI for init, doctor, hooks

## Quick Start

```bash
# 1. Install plugin
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp
cp -r ~/.claude/sdp/prompts/* .claude/

# 2. Start development
@feature "Add user authentication"
@design feature-auth
@build 00-001-01
```

## Languages Supported

| Language | Tests | Coverage | Type Check | Lint |
|----------|-------|----------|------------|------|
| Python   | pytest | pytest-cov | mypy | ruff |
| Java     | Maven/Gradle | JaCoCo | javac | checkstyle |
| Go       | go test | go tool cover | go vet | golint |

## Workflow

1. **`@feature`** - Gather requirements (interactive interview)
2. **`@design`** - Plan workstreams (dependencies, scope)
3. **`@build`** - Execute workstream (TDD cycle)
4. **`@review`** - Quality check (AI validators)
5. **`@deploy`** - Deploy to production

## Documentation

- [Full Tutorial](docs/TUTORIAL.md)
- [Python Examples](docs/examples/python/)
- [Java Examples](docs/examples/java/)
- [Go Examples](docs/examples/go/)
- [Migration Guide](MIGRATION.md)

## Migration from Python SDP

If you're using the Python `sdp` CLI tool:

✅ **Your existing workstreams still work** (prompts are compatible)
✅ **Git hooks continue to work** (use Go binary for convenience)
⚠️ **Quality checks now use AI validation** (no Python required)
📖 **See [MIGRATION.md](MIGRATION.md) for details**

## License

MIT © MSU AI Masters
```

### Step 3: Migration Guide

**File: sdp-plugin/MIGRATION.md**

```markdown
# Migration: Python SDP → Claude Plugin

## What's Different

| Aspect | Python SDP | Claude Plugin |
|--------|-----------|---------------|
| **Installation** | `pip install sdp` | Copy prompts to .claude/ |
| **Dependencies** | Python 3.10+, Click, PyYAML | None (prompts only) |
| **Validation** | pytest, mypy, ruff (tools) | AI analysis (prompts) |
| **Languages** | Python only | Python, Java, Go |
| **Binary** | Required (sdp CLI) | Optional (Go binary) |

## Breaking Changes

### Quality Checks

**OLD (Python SDP):**
```bash
sdp quality check --module src/
# Runs pytest, mypy, ruff (tool-based)
```

**NEW (Plugin):**
```bash
@review
# Uses AI validators (reads code, analyzes)
```

**Impact:**
- Slower validation (AI analysis vs tool execution)
- Language-agnostic (works with Java, Go)
- More flexible (AI can understand context)

### CLI Commands

**OLD (Python SDP):**
```bash
sdp workstream create WS-001-01
sdp workstream verify WS-001-01
```

**NEW (Plugin + Optional Binary):**
```bash
# Option 1: Use Claude skills
@design feature-name
@build 00-001-01

# Option 2: Use Go binary (optional)
sdp init
sdp doctor
sdp hooks install
```

## Migration Steps

### Step 1: Install Plugin

```bash
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp
cp -r ~/.claude/sdp/prompts/* .claude/
```

### Step 2: Verify Workstreams

```bash
# Your existing workstreams work as-is
ls docs/workstreams/backlog/*.md
# All PP-FFF-SS files compatible
```

### Step 3: Test Quality Gates

```bash
# Test @review skill
@review F01

# Expected:
# - AI validators run
# - Coverage, architecture, errors, complexity checked
# - PASS/FAIL verdict
```

### Step 4: Optional Go Binary

```bash
# If you prefer CLI commands
curl -L https://github.com/ai-masters/sdp/releases/latest/download/sdp-darwin-arm64 -o sdp
chmod +x sdp

./sdp init
./sdp doctor
./sdp hooks install
```

## Compatibility Matrix

| Feature | Python SDP | Plugin |
|---------|-----------|--------|
| Workstream format (PP-FFF-SS) | ✅ Yes | ✅ Yes |
| @feature, @design, @build | ✅ Yes | ✅ Yes |
| Multi-agent coordination | ✅ Yes | ✅ Yes |
| Beads integration | ✅ Yes | ⚠️ Planned |
| Git hooks (pre-commit) | ✅ Yes | ✅ Yes (via binary) |
| CLI commands (sdp *) | ✅ Yes | ✅ Yes (via binary) |
| Python quality checks | ✅ Fast (tools) | ⚠️ Slower (AI) |
| Java quality checks | ❌ No | ✅ Yes |
| Go quality checks | ❌ No | ✅ Yes |

## Rollback Plan

If you need to rollback to Python SDP:

```bash
# 1. Uninstall plugin
rm -rf .claude/sdp

# 2. Reinstall Python SDP
pip install sdp

# 3. Your workstreams are unchanged
# (They're just markdown files)
```

## Questions?

- **Documentation:** [docs/](docs/)
- **Issues:** [GitHub Issues](https://github.com/ai-masters/sdp/issues)
- **Discussion:** [GitHub Discussions](https://github.com/ai-masters/sdp/discussions)
```

### Step 4: Release Notes

**File: sdp-plugin/CHANGELOG.md**

```markdown
# Changelog

## [1.0.0] - 2026-02-15

### Added
- 🎉 Language-agnostic Claude Plugin
- ✅ Support for Python, Java, Go projects
- 🤖 AI-based validation (coverage, architecture, errors, complexity)
- 📚 18 skills for workflow automation (@feature, @design, @build, @review, @deploy)
- 👥 11 agents for multi-agent coordination
- 🔧 Optional Go binary for convenience (init, doctor, hooks)
- 📖 Language-specific tutorials (Python, Java, Go)

### Changed
- 🔄 Prompts now work without Python dependencies
- 🔄 Quality gates use AI validation (language-agnostic)
- 🔄 @build skill auto-detects project type
- 🔄 @review skill uses structured AI validators

### Migration
- 📝 Migration guide for Python SDP users
- ✅ Backward compatible with existing workstreams
- ⚠️ Quality checks: Tool-based → AI-based (slower but flexible)

### Documentation
- 📖 Full tutorial with language examples
- 📖 Quick start guides for Python, Java, Go
- 📖 Migration guide from Python SDP

### Deprecated
- ⚠️ Python SDP CLI (deprecated, will be maintained for 6 months)

### Removed
- ❌ Python 3.10+ dependency
- ❌ Click, PyYAML dependencies
- ❌ Language-specific quality checks (pytest, mypy, ruff)
```

### Step 5: Tag and Release

```bash
# 1. Create git tag
git tag -a v1.0.0 -m "Release v1.0.0: Claude Plugin Distribution"

# 2. Push tag
git push origin v1.0.0

# 3. Create GitHub Release
# (via GitHub UI or gh CLI)
gh release create v1.0.0 \
  --title "v1.0.0: Claude Plugin Distribution" \
  --notes-file sdp-plugin/CHANGELOG.md

# 4. Attach binaries (if using Go binary)
gh release upload v1.0.0 \
  bin/sdp-darwin-arm64 \
  bin/sdp-darwin-amd64 \
  bin/sdp-linux-amd64 \
  bin/sdp-windows-amd64.exe

# 5. Submit to Claude Plugin Marketplace
# (via Claude API or web interface)
# Submit sdp-plugin/plugin.json
```

## Verification

```bash
# Test 1: Fresh install (user perspective)
mkdir /tmp/test-sdp-fresh
cd /tmp/test-sdp-fresh
cp -r /path/to/sdp-plugin/prompts/* .claude/

claude "@feature 'Test feature'"
# Expected: Feature workflow executes

# Test 2: Migration from Python SDP
cd existing-python-sdp-project
cp -r /path/to/sdp-plugin/prompts/* .claude/
claude "@build 00-001-01"
# Expected: Existing workstream works

# Test 3: Documentation accuracy
# Follow README quick start
# Expected: All commands work

# Test 4: Plugin validation
cat sdp-plugin/plugin.json | python -m json.tool
# Expected: Valid JSON

# Test 5: Tag created
git tag -l "v1.0.0"
# Expected: Tag exists
```

## Quality Gates

- plugin.json valid JSON schema
- README has all sections (features, quick start, languages)
- MIGRATION.md has comparison table and rollback steps
- CHANGELOG.md has v1.0.0 release notes
- Git tag v1.0.0 created and pushed
- GitHub release created with binaries

## Dependencies

- 00-041-06 (Cross-Language Validation) - needs validated plugin

## Blocks

None (final workstream)
