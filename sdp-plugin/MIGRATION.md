# Migration Guide: Python SDP → Claude Plugin

Migrating from the Python `sdp` CLI tool to the Claude Plugin? This guide shows you what's changed and how to migrate.

## What's Different?

| Aspect | Python SDP | Claude Plugin |
|--------|-----------|---------------|
| **Installation** | `pip install sdp` | Copy prompts to `.claude/` |
| **Dependencies** | Python 3.10+, Click, PyYAML | None (prompts only) |
| **Validation** | pytest, mypy, ruff (tools) | AI analysis (prompts) |
| **Languages** | Python only | Python, Java, Go (any) |
| **Binary** | Required (sdp CLI) | Optional (Go binary) |
| **Speed** | Fast (tool execution) | Slower (AI analysis) |
| **Flexibility** | Python-specific | Language-agnostic |

## Breaking Changes

### 1. Quality Check Method

**OLD (Python SDP):**
```bash
# Tool-based validation (fast)
sdp quality check --module src/
# Runs: pytest --cov, mypy --strict, ruff check
```

**NEW (Plugin):**
```bash
# AI-based validation (language-agnostic)
@review
# Claude reads code, analyzes patterns
```

**Impact:**
- ⚠️ Slower validation (AI analysis vs tool execution)
- ✅ Works with Java, Go (not just Python)
- ✅ More flexible (understands context)
- ✅ No tool dependencies

### 2. CLI Commands

**OLD (Python SDP):**
```bash
sdp workstream create WS-001-01
sdp workstream verify WS-001-01
sdp quality check --module src/
sdp hooks install
```

**NEW (Plugin + Optional Binary):**
```bash
# Option 1: Use Claude skills directly
@design feature-name
@build 00-001-01
@review F01

# Option 2: Use Go binary (optional)
./sdp init
./sdp doctor
./sdp hooks install
```

**Impact:**
- ✅ More intuitive workflow (@design, @build, @review)
- ✅ Binary is optional (not required)
- ✅ Binary provides convenience commands

### 3. Workstream ID Format

**OLD (Python SDP):**
```
WS-FFF-SS format
Example: WS-001-01
```

**NEW (Plugin):**
```
PP-FFF-SS format (project ID prefix)
Example: 00-001-01

Legacy format still supported!
```

**Impact:**
- ✅ Backward compatible (WS-001-01 still works)
- ✅ Multi-project support (01-001-01, 02-001-01)
- ⚠️ New projects should use PP-FFF-SS

## Migration Steps

### Step 1: Install Plugin

```bash
# Clone plugin repository
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp

# Copy prompts to your project
cp -r ~/.claude/sdp/prompts/* .claude/

# Verify installation
ls .claude/skills/
# Should show: feature.md, design.md, build.md, review.md, deploy.md
```

### Step 2: Verify Workstreams

**Your existing workstreams still work!** No migration needed.

```bash
# Check existing workstreams
ls docs/workstreams/backlog/*.md

# All WS-FFF-SS files are compatible
# PP-FFF-SS format also supported
```

**Example:** Both formats work:
```
✅ WS-001-01 (legacy format)
✅ 00-001-01 (new format with project ID)
```

### Step 3: Test Quality Gates

```bash
# Open Claude Code in your project

# Run review skill
@review F01

# Expected behavior:
# 1. AI reads your code
# 2. Runs validators:
#    - Coverage: Maps tests to functions
#    - Architecture: Checks imports
#    - Errors: Finds unsafe patterns
#    - Complexity: Counts lines
# 3. Produces verdict: ✅ PASS or ❌ FAIL
```

### Step 4: Optional Go Binary

**If you liked the Python SDP CLI:**

```bash
# Download Go binary (macOS arm64 example)
curl -L https://github.com/ai-masters/sdp/releases/latest/download/sdp-darwin-arm64 -o sdp
chmod +x sdp

# Run familiar commands
./sdp init --project-type=python
./sdp doctor
./sdp hooks install
```

**Available platforms:**
- macOS ARM64 (Apple Silicon)
- macOS AMD64 (Intel)
- Linux AMD64
- Windows AMD64

## Compatibility Matrix

| Feature | Python SDP | Plugin | Notes |
|---------|-----------|--------|-------|
| **Workstream format** | ✅ WS-FFF-SS | ✅ WS-FFF-SS, PP-FFF-SS | Backward compatible |
| **@feature skill** | ❌ No | ✅ Yes | Progressive disclosure |
| **@design skill** | ❌ No | ✅ Yes | Interactive planning |
| **@build skill** | ✅ sdp build | ✅ @build | Same TDD workflow |
| **@review skill** | ✅ sdp review | ✅ @review | Now uses AI validators |
| **@deploy skill** | ✅ sdp deploy | ✅ @deploy | Same GitFlow workflow |
| **Multi-agent mode** | ✅ Yes | ✅ Yes | Unchanged |
| **Git hooks** | ✅ Yes | ✅ Yes (via binary) | Same functionality |
| **CLI commands** | ✅ sdp * | ✅ ./sdp * (optional) | Binary provides same commands |
| **Python validation** | ✅ Fast (tools) | ⚠️ Slower (AI) | Tradeoff: flexibility vs speed |
| **Java validation** | ❌ No | ✅ Yes | New capability |
| **Go validation** | ❌ No | ✅ Yes | New capability |
| **Dependencies** | ❌ Python required | ✅ None | Major improvement |

## Rollback Plan

Need to rollback to Python SDP? No problem:

```bash
# 1. Remove plugin
rm -rf .claude/

# 2. Reinstall Python SDP
pip install sdp

# 3. Your workstreams are unchanged
# (They're just markdown files)
```

**Your work is safe!** Workstreams are markdown files, compatible with both versions.

## Common Questions

### Q: Do I need to rewrite my workstreams?

**A:** No! Workstreams are markdown files in PP-FFF-SS format. Both versions read the same format.

### Q: Will my existing tests pass?

**A:** Yes! The plugin uses the same test commands (pytest, mvn test, go test).

### Q: What happens to my quality gates?

**A:** They still work, but now use AI validation instead of tools. Same thresholds (≥80% coverage, etc.).

### Q: Is AI validation as good as tools?

**A:** It's different:
- **Tools:** Faster, more precise, but Python-only
- **AI:** Slower, more flexible, language-agnostic

For Python projects with tools configured: Use tools
For Java/Go projects: Use AI validators

### Q: Can I keep using Python SDP?

**A:** Yes! Python SDP will be maintained for 6 months (until 2026-08-03). Both versions can coexist.

### Q: Do I need the Go binary?

**A:** No! The binary is optional convenience. Prompts work standalone.

### Q: What if I don't like the plugin?

**A:** Rollback is simple (see above). Your work is safe in markdown files.

## Migration Example

Let's migrate a real project:

### Before (Python SDP)

```bash
# Installation
pip install sdp

# Create workstream
sdp workstream create WS-001-01 \
  --goal "Add user authentication" \
  --ac "AC1: Create User entity" \
  --ac "AC2: Add login endpoint"

# Execute workstream
sdp build WS-001-01
# Runs: pytest, mypy, ruff

# Quality check
sdp quality check --module src/
# Runs: pytest --cov, mypy --strict, ruff check

# Deploy
sdp deploy WS-001-01
```

### After (Plugin)

```bash
# Installation (no pip)
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp
cp -r ~/.claude/sdp/prompts/* .claude/

# Create workstream (interactive)
@feature "Add user authentication"
# Claude asks deep questions about requirements

# Plan workstreams
@design feature-auth
# Claude decomposes into workstreams

# Execute workstream
@build 00-001-01
# Runs: pytest tests/ -v
# AI validates: coverage, architecture, errors, complexity

# Quality check
@review F01
# AI runs all validators

# Deploy
@deploy F01
```

**Key differences:**
- ✅ No `pip install` required
- ✅ Interactive requirements gathering (better specs)
- ✅ Language-agnostic (works with Java, Go)
- ⚠️ AI validation (slower but flexible)

## Tips for Smooth Migration

1. **Start with new projects** - Try the plugin on a new feature first
2. **Keep Python SDP installed** - Run both in parallel during transition
3. **Use AI validators for new languages** - Java/Go support is excellent
4. **Keep tools for Python** - If you have pytest/mypy/ruff, keep using them
5. **Read the tutorial** - [TUTORIAL.md](docs/TUTORIAL.md) has comprehensive examples

## Need Help?

- **Documentation:** [docs/](docs/)
- **Tutorial:** [docs/TUTORIAL.md](docs/TUTORIAL.md)
- **Examples:** [docs/examples/](docs/examples/)
- **Issues:** [GitHub Issues](https://github.com/ai-masters/sdp/issues)
- **Discussions:** [GitHub Discussions](https://github.com/ai-masters/sdp/discussions)

## Timeline

- **2026-02-03:** Plugin v1.0.0 released
- **2026-02-06:** Python SDP archived (code moved to `src/sdp-deprecated/`)
  - CLI entry points removed from `pyproject.toml`
  - 252 Python files (25K LOC) archived
  - Go SDP becomes the primary implementation
- **2026-08-03:** Python SDP EOL (end of 6-month grace period)
- **After 2026-08-03:** Python SDP community-supported only

**Recommendation:** Migrate to Go SDP + Plugin immediately. Python version is deprecated.

## Archival Status

**Python SDP Code:**
- ✅ Archived to `src/sdp-deprecated/` (2026-02-06)
- ✅ CLI entry points removed
- ✅ Poetry package no longer installs commands
- ⚠️ Available for reference only (not maintained)

**Go SDP:**
- ✅ Primary implementation (v1.0.0-go)
- ✅ All 4 parity commands implemented
- ✅ Active development
- ✅ Multi-language support (Python, Java, Go)

**Migration Path:**
```bash
# Old (deprecated)
pip install sdp  # Still works but deprecated

# New (recommended)
go install github.com/fall-out-bug/sdp/sdp-plugin/cmd/sdp@latest
# OR download binary from releases
```
