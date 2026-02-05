# Python SDP Deprecation: Migration Guide

**Published:** 2026-02-05
**Effective:** 2026-02-05
**Maintenance Ends:** 2026-08-03

---

## Executive Summary

The Python SDP implementation is **deprecated** and succeeded by the [SDP Plugin](https://github.com/ai-masters/sdp-plugin). This guide explains why, what changed, and how to migrate.

**Key Points:**
- Your existing workstreams are **compatible** with the plugin
- Migration is **simple** (copy prompts, no code changes)
- The plugin is **language-agnostic** (Python, Java, Go, any)
- Python SDP enters **maintenance mode** (bug fixes only)

---

## Why the Change?

### Limitations of Python SDP

The Python implementation has fundamental limitations:

1. **Language Lock-in** - Works only with Python projects
2. **Heavy Dependencies** - Requires Python 3.10+, Poetry, pytest, mypy, ruff
3. **Tool-Based Validation** - Fast but inflexible (Python-specific tools)
4. **Installation Friction** - `pip install sdp` + Python environment setup
5. **Maintenance Burden** - Python packaging, PyPI, version conflicts

### Advantages of SDP Plugin

The Go-based plugin addresses these limitations:

| Aspect | Python SDP | SDP Plugin | Benefit |
|--------|-----------|------------|---------|
| **Languages** | Python only | Python, Java, Go, any | One tool for all projects |
| **Installation** | `pip install sdp` | Copy prompts | No dependencies |
| **Dependencies** | Python, Poetry, tools | None (optional Go binary) | Zero install friction |
| **Validation** | pytest, mypy, ruff | AI analysis | Language-agnostic |
| **Quality Gates** | Tool-based (fast) | AI-based (flexible) | Understands context |
| **Maintenance** | Python packaging | Go binary + prompts | Simpler deployment |
| **Extensibility** | Python code only | Edit prompts | Customize workflows |

### The Real Value Proposition

**SDP Plugin separates concerns:**

```
┌─────────────────────────────────────────┐
│  SDP Plugin (Go + Prompts)              │
├─────────────────────────────────────────┤
│  • Prompts: Workflow orchestration      │
│  • Binary: Convenience (init, doctor)   │
│  • No language-specific logic           │
└─────────────────────────────────────────┘
           │
           │ Uses language-specific tools
           ▼
┌─────────────────────────────────────────┐
│  Your Project (Python/Java/Go)          │
├─────────────────────────────────────────┤
│  • Python: pytest, mypy, ruff           │
│  • Java: Maven, JaCoCo, checkstyle      │
│  • Go: go test, go vet, golint          │
└─────────────────────────────────────────┘
```

**Python SDP tightly coupled:**

```
┌─────────────────────────────────────────┐
│  Python SDP (Everything in Python)      │
├─────────────────────────────────────────┤
│  • Workflow orchestration (Python)      │
│  • Validation logic (Python)            │
│  • pytest/mypy/ruff integration (Py)    │
│  • Hard to extend to Java/Go            │
└─────────────────────────────────────────┘
```

---

## Feature Parity Matrix

### What Works the Same

| Feature | Python SDP | SDP Plugin | Status |
|---------|-----------|------------|--------|
| **Workstream format** | ✅ PP-FFF-SS | ✅ PP-FFF-SS | Identical |
| **Markdown files** | ✅ Yes | ✅ Yes | Compatible |
| **@feature skill** | ✅ Yes | ✅ Yes | Same prompts |
| **@design skill** | ✅ Yes | ✅ Yes | Same prompts |
| **@build skill** | ✅ Yes | ✅ Yes | Same TDD cycle |
| **@deploy skill** | ✅ Yes | ✅ Yes | Same GitFlow |
| **Multi-agent mode** | ✅ Yes | ✅ Yes | Unchanged |
| **Beads integration** | ✅ Yes | ✅ Yes | Compatible |
| **Git hooks** | ✅ Yes | ✅ Yes | Same functionality |
| **Documentation** | ✅ Extensive | ✅ Extensive | Similar coverage |

### What's Better

| Feature | Python SDP | SDP Plugin | Improvement |
|---------|-----------|------------|-------------|
| **Language support** | Python only | Python, Java, Go, any | **Multi-language** |
| **Installation** | `pip install sdp` | Copy prompts | **No dependencies** |
| **Dependencies** | Python, Poetry, tools | None (optional binary) | **Zero friction** |
| **@review skill** | Tool-based | AI-based validators | **Language-agnostic** |
| **@debug skill** | Tool-based | AI-based | **Understands context** |
| **Customization** | Edit Python code | Edit prompts | **Easier to extend** |
| **Updates** | Repackage Python | Copy new prompts | **Instant updates** |

### What's Different (Tradeoffs)

| Feature | Python SDP | SDP Plugin | Tradeoff |
|---------|-----------|------------|----------|
| **Validation speed** | Fast (tools) | Slower (AI) | Speed vs flexibility |
| **Validation precision** | Exact (tool output) | Approximate (AI) | Precision vs agnosticism |
| **Setup time** | 5-10 min (pip) | 1-2 min (copy) | Minor difference |
| **Learning curve** | Moderate | Lower | No Python knowledge needed |

### What's Missing in Plugin

| Feature | Python SDP | SDP Plugin | Mitigation |
|---------|-----------|------------|------------|
| **@oneshot skill** | ✅ Yes | ❌ No | Use manual @build for each WS |
| **Telegram notifications** | ✅ Yes | ❌ No | Roadmap item |
| **Checkpoint system** | ✅ Yes | ❌ No | Use git branches instead |
| **Progressive disclosure (@feature)** | ✅ Yes | ✅ Yes | **Already implemented** |

**Note:** The plugin is actively developed. Missing features are on the roadmap.

---

## Migration Steps

### Step 1: Install SDP Plugin (2 minutes)

```bash
# Clone the plugin repository
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp

# Copy prompts to your project
cd /path/to/your/project
cp -r ~/.claude/sdp/prompts/* .claude/

# Verify installation
ls .claude/skills/
# Should see: feature.md, design.md, build.md, review.md, deploy.md, etc.
```

**That's it!** No `pip install`, no Python environment, no dependencies.

### Step 2: Verify Your Workstreams (1 minute)

**Your existing workstreams work as-is!** No conversion needed.

```bash
# Check your workstreams
ls docs/workstreams/backlog/*.md

# Example: All these formats work:
# ✅ WS-001-01 (legacy format)
# ✅ 00-001-01 (new format with project ID)
# ✅ 01-042-07 (multi-project format)
```

**Why compatible?** Workstreams are markdown files with frontmatter. Both versions read the same format.

### Step 3: Test the Plugin (5 minutes)

```bash
# Open Claude Code in your project
cd /path/to/your/project
claude

# Test 1: Create a new feature
@feature "Test migration"
# Expected: Claude interviews you (same as Python SDP)

# Test 2: Plan workstreams
@design test-migration
# Expected: Claude decomposes into workstreams (same as Python SDP)

# Test 3: Execute a workstream
@build 00-001-01
# Expected: TDD cycle (Red → Green → Refactor)

# Test 4: Quality check
@review F01
# Expected: AI validators analyze code
```

**What's different:**
- `@review` now uses AI validators (slower but language-agnostic)
- No `@oneshot` skill (use manual @build for each WS)
- No Telegram notifications (roadmap item)

### Step 4: Update Documentation (Optional)

If you have project-specific documentation referencing Python SDP:

```bash
# Find references to Python SDP
grep -r "pip install sdp" docs/
grep -r "sdp-cli" docs/
grep -r "Python SDP" docs/

# Update to reference SDP Plugin
# Old: "pip install sdp"
# New: "See https://github.com/ai-masters/sdp-plugin"
```

### Step 5: Uninstall Python SDP (Optional)

```bash
# Remove Python package
pip uninstall sdp

# Remove virtual environment (if isolated)
rm -rf ~/.venv/sdp

# Your workstreams are safe (they're markdown files)
```

**Recommendation:** Keep Python SDP installed during transition period. Run both in parallel to verify compatibility.

---

## Migration Example

Let's migrate a real project: `hw_checker` (Python project)

### Before (Python SDP)

```bash
# 1. Install Python SDP
pip install sdp

# 2. Initialize project
sdp init --project-type=python

# 3. Create feature
sdp feature create "Add password reset" \
  --user-story "User can reset password via email"

# 4. Plan workstreams (manual)
# Create docs/workstreams/backlog/WS-001-01.md
# Create docs/workstreams/backlog/WS-001-02.md

# 5. Execute workstream
sdp build WS-001-01
# Runs: pytest tests/ -v, mypy src/, ruff check

# 6. Quality check
sdp quality check --module src/
# Runs: pytest --cov, mypy --strict, ruff check

# 7. Deploy
sdp deploy WS-001-01
```

### After (SDP Plugin)

```bash
# 1. Install plugin (no pip)
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp
cp -r ~/.claude/sdp/prompts/* .claude/

# 2. Initialize (optional binary)
curl -L https://github.com/ai-masters/sdp/releases/latest/download/sdp-darwin-arm64 -o sdp
chmod +x sdp
./sdp init --project-type=python

# 3. Create feature (interactive)
@feature "Add password reset"
# Claude interviews you:
# - Who are the users?
# - What's the success metric?
# - What are the edge cases?
# Result: Comprehensive spec in docs/drafts/

# 4. Plan workstreams (interactive)
@design feature-password-reset
# Claude explores codebase, asks:
# - Email service available?
# - Token storage approach?
# - JWT vs sessions?
# Result: 00-001-01, 00-001-02, etc.

# 5. Execute workstream
@build 00-001-01
# Runs: pytest tests/ -v (same command)
# AI validates: coverage, architecture, errors

# 6. Quality check
@review F01
# AI runs validators:
# - /coverage-validator (maps tests to code)
# - /architecture-validator (checks imports)
# - /error-validator (finds unsafe patterns)
# - /complexity-validator (counts lines)

# 7. Deploy
@deploy F01
# Same GitFlow workflow
```

**Key Differences:**
- ✅ No `pip install` required
- ✅ Interactive requirements gathering (better specs)
- ✅ Language-agnostic (works with Java, Go)
- ⚠️ AI validation (slower but flexible)
- ❌ No `@oneshot` (use manual @build)

---

## Common Questions

### Q: Do I need to rewrite my workstreams?

**A:** No! Workstreams are markdown files in PP-FFF-SS format. Both versions read the same format.

**Example:**
```markdown
---
ws_id: 00-001-01
feature: F001
status: backlog
---

## WS-00-001-01: User Entity

### Goal
Create User domain entity with validation.
```

This file works with **both** Python SDP and SDP Plugin.

### Q: Will my existing tests pass?

**A:** Yes! The plugin uses the same test commands:

```bash
# Python
pytest tests/ -v

# Java
mvn test

# Go
go test ./...
```

**What changed:** Validation is now AI-based (not tool-based), but test execution is identical.

### Q: What happens to my quality gates?

**A:** They still work, but use AI validation:

| Quality Gate | Python SDP | SDP Plugin |
|--------------|-----------|------------|
| **Coverage ≥80%** | pytest-cov (exact) | AI maps tests → code (approx) |
| **Type hints** | mypy --strict (exact) | AI checks annotations (approx) |
| **Error handling** | ruff (pattern) | AI finds unsafe patterns (context) |
| **File size <200 LOC** | Python script | AI counts lines (approx) |

**Tradeoff:** Precision vs flexibility. For Python projects with tools configured, you can still run `pytest --cov`, `mypy --strict`, `ruff check` manually.

### Q: Is AI validation as good as tools?

**A:** It's different:

**Tools (pytest, mypy, ruff):**
- ✅ Faster (seconds vs minutes)
- ✅ More precise (exact numbers)
- ✅ Reliable (deterministic)
- ❌ Python-only

**AI Validators:**
- ✅ Language-agnostic (works with Java, Go)
- ✅ Understands context (business logic)
- ✅ Flexible (customizable prompts)
- ⚠️ Slower (AI analysis takes time)
- ⚠️ Approximate (not exact numbers)

**Recommendation:**
- For **Python projects** with tools configured: Use tools (run manually)
- For **Java/Go projects**: Use AI validators (only option)
- For **mixed-language projects**: Use AI validators (unified approach)

### Q: Can I keep using Python SDP?

**A:** Yes, but with limitations:

**Maintenance Period:** 2026-02-05 to 2026-08-03 (6 months)

**During maintenance period:**
- ✅ Bug fixes (critical issues)
- ✅ Security patches (vulnerabilities)
- ❌ New features (deprecated)
- ❌ Enhancements (use plugin instead)

**After 2026-08-03:**
- ❌ No updates (community-supported only)
- ✅ Still works (no breaking changes)
- ⚠️ May become incompatible with future Claude Code versions

**Recommendation:** Migrate to plugin during the 6-month overlap period.

### Q: Do I need the Go binary?

**A:** No! The binary is optional convenience.

**Plugin without binary:**
```bash
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp
cp -r ~/.claude/sdp/prompts/* .claude/
# Done! Skills work standalone
```

**Plugin with binary:**
```bash
curl -L https://github.com/ai-masters/sdp/releases/latest/download/sdp-darwin-arm64 -o sdp
chmod +x sdp

# Optional convenience commands:
./sdp init --project-type=python
./sdp doctor
./sdp hooks install
```

**Binary provides:**
- Project initialization (scaffolding)
- Health checks (Python version, tools, etc.)
- Git hooks installation
- Workstream verification

**All binary features are optional.** Skills work standalone.

### Q: What if I don't like the plugin?

**A:** Rollback is simple:

```bash
# 1. Remove plugin prompts
rm -rf .claude/

# 2. Reinstall Python SDP
pip install sdp

# 3. Your workstreams are unchanged
# (They're just markdown files in docs/workstreams/)
```

**Your work is safe!** Workstreams are markdown files, compatible with both versions.

### Q: What about @oneshot and autonomous execution?

**A:** Not yet implemented in the plugin (roadmap item).

**Workaround:**
```bash
# Old way (Python SDP)
@oneshot F01
# Executes all workstreams autonomously

# New way (Plugin)
@build 00-001-01
@build 00-001-02
@build 00-001-03
# Manual execution for each workstream
```

**Why removed:** `@oneshot` complexity wasn't worth the maintenance burden. Most users prefer manual @build for visibility.

### Q: What about Telegram notifications?

**A:** Not yet implemented in the plugin (roadmap item).

**Workaround:** Use Claude Code's built-in notifications or desktop alerts.

### Q: Will the plugin get feature parity?

**A:** Yes, the plugin is actively developed:

**Short-term (next 3 months):**
- ✅ Multi-language examples (Java, Go)
- ✅ Enhanced AI validators
- 🔄 @oneshot skill (under consideration)
- 🔄 Telegram notifications (under consideration)

**Long-term:**
- Custom validator prompts
- Language-specific workflows
- Integration with other AI-IDEs (Cursor, OpenCode)

**Track progress:** [GitHub Issues](https://github.com/ai-masters/sdp-plugin/issues)

---

## Rollback Plan

If you need to rollback to Python SDP:

### Step 1: Remove Plugin

```bash
# Remove plugin prompts
rm -rf .claude/

# Remove Go binary (if installed)
rm sdp
```

### Step 2: Reinstall Python SDP

```bash
# Install from source
git clone https://github.com/fall-out-bug/sdp.git
cd sdp
pip install -e .

# Or from PyPI (if available)
pip install sdp-cli
```

### Step 3: Verify Workstreams

```bash
# Your workstreams are unchanged
ls docs/workstreams/backlog/*.md

# All workstreams work with Python SDP
sdp build 00-001-01
```

### Step 4: Continue Development

```bash
# Use Python SDP commands
sdp feature create "New feature"
sdp build 00-002-01
sdp quality check --module src/
sdp deploy 00-002-01
```

**Your work is safe!** No data loss, no conversion needed.

---

## Timeline

| Date | Milestone | Status |
|------|-----------|--------|
| 2026-02-03 | SDP Plugin v1.0.0 released | ✅ Complete |
| 2026-02-05 | Python SDP deprecation announced | ✅ Complete |
| 2026-02-05 to 2026-08-03 | Maintenance period (bug fixes only) | 🔄 In progress |
| 2026-08-03 | Maintenance ends (community-supported) | ⏳ Future |
| After 2026-08-03 | Plugin recommended, Python SDP legacy | ⏳ Future |

**Recommendation:** Migrate to plugin during the 6-month overlap period (2026-02-05 to 2026-08-03).

---

## Comparison Table

### Installation

| Aspect | Python SDP | SDP Plugin |
|--------|-----------|------------|
| **Install time** | 5-10 min | 1-2 min |
| **Commands** | `pip install sdp` | `git clone + cp` |
| **Dependencies** | Python 3.10+, Poetry, pytest, mypy, ruff | None (optional Go binary) |
| **Disk space** | ~50 MB (Python env) | ~1 MB (prompts) |
| **Setup friction** | Medium (Python env) | Low (copy files) |

### Features

| Feature | Python SDP | SDP Plugin |
|---------|-----------|------------|
| **Workstreams** | ✅ Yes | ✅ Yes |
| **@feature skill** | ✅ Yes | ✅ Yes |
| **@design skill** | ✅ Yes | ✅ Yes |
| **@build skill** | ✅ Yes | ✅ Yes |
| **@review skill** | ✅ Tool-based | ✅ AI-based |
| **@deploy skill** | ✅ Yes | ✅ Yes |
| **@debug skill** | ✅ Tool-based | ✅ AI-based |
| **@oneshot skill** | ✅ Yes | ❌ No |
| **Multi-agent** | ✅ Yes | ✅ Yes |
| **Beads integration** | ✅ Yes | ✅ Yes |
| **Git hooks** | ✅ Yes | ✅ Yes |
| **Telegram** | ✅ Yes | ❌ No |
| **Progressive disclosure** | ✅ Yes | ✅ Yes |

### Validation

| Quality Gate | Python SDP | SDP Plugin |
|--------------|-----------|------------|
| **Coverage ≥80%** | pytest-cov (exact) | AI maps tests → code |
| **Type hints** | mypy --strict (exact) | AI checks annotations |
| **Error handling** | ruff (pattern) | AI finds unsafe patterns |
| **File size <200 LOC** | Python script | AI counts lines |
| **Clean Architecture** | Import analysis | AI checks imports |
| **Languages** | Python only | Python, Java, Go, any |
| **Speed** | Fast (seconds) | Slower (minutes) |
| **Precision** | Exact (numbers) | Approximate (context) |

### Workflow

| Workflow | Python SDP | SDP Plugin |
|----------|-----------|------------|
| **Create feature** | `sdp feature create` | `@feature` (interactive) |
| **Plan workstreams** | Manual markdown | `@design` (interactive) |
| **Execute workstream** | `sdp build WS-ID` | `@build WS-ID` |
| **Quality check** | `sdp quality check` | `@review` (AI validators) |
| **Deploy** | `sdp deploy WS-ID` | `@deploy F-ID` |
| **Autonomous execution** | `@oneshot F-ID` | Manual @build only |

### Maintenance

| Aspect | Python SDP | SDP Plugin |
|--------|-----------|------------|
| **Status** | Deprecated | Active development |
| **Bug fixes** | Until 2026-08-03 | Ongoing |
| **New features** | No | Yes (roadmap) |
| **Security updates** | Until 2026-08-03 | Ongoing |
| **Community support** | After 2026-08-03 | Ongoing |
| **Long-term viability** | Low | High |

---

## Decision Matrix

### Should You Migrate?

**Migrate if:**
- ✅ You work with multiple languages (Python, Java, Go)
- ✅ You want zero-dependency setup
- ✅ You value language-agnostic validation
- ✅ You want active development and new features
- ✅ You want to customize workflows (edit prompts)

**Stay with Python SDP if:**
- ⚠️ You only work with Python
- ⚠️ You need fast validation (tools vs AI)
- ⚠️ You rely on @oneshot autonomous execution
- ⚠️ You need Telegram notifications
- ⚠️ You're in maintenance mode (no new features)

**Both during transition:**
- ✅ Keep Python SDP installed
- ✅ Test plugin on new features
- ✅ Compare workflows
- ✅ Choose based on experience

---

## Resources

### SDP Plugin

- **Repository:** [https://github.com/ai-masters/sdp-plugin](https://github.com/ai-masters/sdp-plugin)
- **Documentation:** [https://github.com/ai-masters/sdp-plugin/blob/main/docs/TUTORIAL.md](https://github.com/ai-masters/sdp-plugin/blob/main/docs/TUTORIAL.md)
- **Migration Guide:** [https://github.com/ai-masters/sdp-plugin/blob/main/MIGRATION.md](https://github.com/ai-masters/sdp-plugin/blob/main/MIGRATION.md)
- **Examples:** [https://github.com/ai-masters/sdp-plugin/tree/main/docs/examples](https://github.com/ai-masters/sdp-plugin/tree/main/docs/examples)

### Python SDP (Legacy)

- **Repository:** [https://github.com/fall-out-bug/sdp](https://github.com/fall-out-bug/sdp)
- **Documentation:** [docs/PROTOCOL.md](PROTOCOL.md)
- **Support:** Community (after 2026-08-03)

### Getting Help

- **Plugin Issues:** [https://github.com/ai-masters/sdp-plugin/issues](https://github.com/ai-masters/sdp-plugin/issues)
- **Plugin Discussions:** [https://github.com/ai-masters/sdp-plugin/discussions](https://github.com/ai-masters/sdp-plugin/discussions)
- **Python SDP Issues:** [https://github.com/fall-out-bug/sdp/issues](https://github.com/fall-out-bug/sdp/issues) (during maintenance period)

---

## Appendix: Feature Checklist

Use this checklist to verify migration:

### Pre-Migration Checklist

- [ ] Read this guide completely
- [ ] Reviewed SDP Plugin documentation
- [ ] Tested plugin on sample project
- [ ] Backed up workstreams (git commit)
- [ ] Notified team (if applicable)

### Migration Checklist

- [ ] Installed SDP Plugin (`git clone + cp`)
- [ ] Verified prompts in `.claude/skills/`
- [ ] Tested `@feature` skill
- [ ] Tested `@design` skill
- [ ] Tested `@build` skill
- [ ] Tested `@review` skill (AI validators)
- [ ] Verified existing workstreams work
- [ ] Updated project documentation (if applicable)

### Post-Migration Checklist

- [ ] All workstreams execute successfully
- [ ] Quality gates pass (coverage, type hints, etc.)
- [ ] Team trained on plugin (if applicable)
- [ ] Python SDP uninstalled (optional)
- [ ] Feedback provided to plugin maintainers

---

## Conclusion

The Python SDP served the community well, but the SDP Plugin represents the next evolution:

**Key improvements:**
- ✅ Multi-language support (Python, Java, Go, any)
- ✅ Zero dependencies (no Python required)
- ✅ Language-agnostic validation (AI-based)
- ✅ Simpler installation (copy prompts)
- ✅ Active development (new features)

**Migration path:**
- ✅ Simple (copy prompts, no code changes)
- ✅ Safe (workstreams are compatible)
- ✅ Reversible (rollback is easy)
- ✅ Supported (6-month maintenance period)

**Recommendation:** Migrate to the [SDP Plugin](https://github.com/ai-masters/sdp-plugin) during the 6-month overlap period (2026-02-05 to 2026-08-03).

**Thank you** for using Python SDP! We hope you enjoy the plugin's improvements.

---

**Last Updated:** 2026-02-05
**Next Review:** 2026-08-03 (end of maintenance period)
