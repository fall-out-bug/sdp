# Changelog

All notable changes to the Spec-Driven Protocol (SDP) Claude Plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0-go] - 2026-02-06

### Changed

🎉 **Go Migration Complete** - Python SDP deprecated, Go binary now primary

#### Migration Highlights
- ✅ **Complete Go Implementation** - All 13 workstreams migrated
- ✅ **Language-Agnostic Design** - Works with Python, Java, Go, and any language
- ✅ **18 Skills** - Complete workflow automation
  - `@feature` - Progressive feature development (unified entry point)
  - `@idea` - Requirements gathering with Beads
  - `@design` - Workstream planning with dependencies
  - `@build` - Workstream execution with TDD cycle
  - `@review` - Quality validation with AI
  - `@deploy` - Automated deployment workflow
  - `@debug` - Systematic debugging (scientific method)
  - `@bugfix` - Quality bug fixes (P1/P2)
  - `@hotfix` - Emergency fixes (P0)
  - `@issue` - Bug analysis and routing
  - `@oneshot` - Autonomous multi-agent execution
  - `@init` - Project initialization
  - And 5 more...

- ✅ **11 Agents** - Multi-agent coordination
  - planner, builder, reviewer, tester, architect, analyst, debugger, deployer, orchestrator, and more

- ✅ **4 AI Validators** - Language-agnostic quality validation
  - `/coverage-validator` - Test coverage analysis (≥80% threshold)
  - `/architecture-validator` - Clean Architecture enforcement
  - `/error-validator` - Error handling audit
  - `/complexity-validator` - Code complexity analysis
  - `/all` - Orchestrator for unified quality gates

- ✅ **Optional Go Binary** - Convenience CLI (~5.5MB, single executable)
  - `sdp init` - Initialize project with prompts
  - `sdp doctor` - Check environment
  - `sdp hooks install/uninstall` - Manage Git hooks
  - Cross-platform: macOS ARM64/AMD64, Linux AMD64, Windows AMD64

#### Documentation
- ✅ **Comprehensive Tutorial** (7,500 words)
  - Quick start (5 minutes)
  - Language examples (Python, Java, Go)
  - Quality gates reference
  - Migration guide from Python SDP
  - Troubleshooting

- ✅ **Language-Specific Quickstarts**
  - Python: pytest, mypy, ruff workflow
  - Java: Maven, JaCoCo workflow
  - Go: go test, go vet workflow

- ✅ **Migration Guide** (MIGRATION.md)
  - Breaking changes explained
  - Step-by-step migration
  - Compatibility matrix
  - Rollback plan

#### Quality Gates
- ✅ **Coverage ≥80%** - AI maps tests to functions
- ✅ **Type Safety** - Complete type signatures
- ✅ **Error Handling** - No unsafe patterns
- ✅ **File Size** - <200 LOC per file
- ✅ **Architecture** - Clean layer separation

### Changed

#### from Python SDP (CLI tool) to Claude Plugin

| Aspect | Python SDP | Claude Plugin |
|--------|-----------|---------------|
| **Installation** | `pip install sdp` | Copy prompts to `.claude/` |
| **Dependencies** | Python 3.10+, Click, PyYAML | None (prompts only) |
| **Validation** | pytest, mypy, ruff (tools) | AI analysis (language-agnostic) |
| **Languages** | Python only | Python, Java, Go (any) |
| **Binary** | Required (sdp CLI) | Optional (Go binary) |
| **Speed** | Fast (tool execution) | Slower (AI analysis) |
| **Flexibility** | Python-specific | Language-agnostic |

#### Workstream ID Format

- **OLD:** `WS-FFF-SS` (e.g., WS-001-01)
- **NEW:** `PP-FFF-SS` (e.g., 00-001-01)
- **Backward Compatible:** Legacy format still supported

#### Quality Check Method

- **OLD:** Tool-based validation
  ```bash
  sdp quality check --module src/
  # Runs: pytest, mypy, ruff
  ```

- **NEW:** AI-based validation
  ```bash
  @review
  # Claude reads code, analyzes patterns
  ```

**Impact:**
- ⚠️ Slower validation (AI vs tools)
- ✅ Language-agnostic (works with Java, Go)
- ✅ More flexible (understands context)

### Deprecated

- ⚠️ **Python SDP CLI** - Deprecated, will be maintained until 2026-08-03 (6 months)
- ⚠️ **Legacy workstream ID format** (`WS-FFF-SS`) - Still supported, new projects should use `PP-FFF-SS`

### Removed

- ❌ **Python 3.10+ dependency** - No longer required
- ❌ **Click dependency** - No longer required
- ❌ **PyYAML dependency** - No longer required
- ❌ **Language-specific quality checks** - Replaced with AI validators

### Migration

- 📖 **Migration guide available** - [MIGRATION.md](MIGRATION.md)
- ✅ **Backward compatible** - Existing workstreams work unchanged
- ⚠️ **Quality checks changed** - Tool-based → AI-based

### Security

- ✅ **No code execution** - Prompts only, Claude handles execution
- ✅ **No external dependencies** - Works offline (except Claude API)
- ✅ **Transparent** - All prompts are markdown files (readable)

### Performance

- **AI Validation:** ~10-30 seconds per validation (depends on codebase size)
- **Go Binary:** ~5.5MB, instant execution
- **Test Execution:** Unchanged (uses existing tools: pytest, mvn, go test)

### Documentation

- **Tutorial:** 7,500 words, comprehensive guide
- **Examples:** Python, Java, Go quickstarts
- **Migration:** Detailed migration from Python SDP
- **API Reference:** Skills, agents, validators documented

### Compatibility

- **Claude Code:** 1.0.0+
- **Python:** 3.10+ (if using pytest/mypy/ruff)
- **Java:** 17+ (if using Maven/Gradle)
- **Go:** 1.21+ (if using go test)
- **Operating Systems:** macOS, Linux, Windows

### Known Limitations

1. **AI Validation Speed** - Slower than tool-based validation
2. **Go Binary Availability** - Optional, not required for core functionality
3. **Beads Integration** - Planned for future release
4. **IDE Integration** - Currently Claude Code only

### Future Enhancements

Planned for future releases:
- Beads integration for requirements management
- IDE integrations (VS Code, IntelliJ)
- More language examples (Rust, TypeScript, C#)
- Performance improvements for AI validators
- Web dashboard for workstream visualization

### Acknowledgments

Built with ❤️ for the Claude Code community.

Inspired by:
- Test-Driven Development (Kent Beck)
- Clean Architecture (Robert C. Martin)
- Domain-Driven Design (Eric Evans)
- The Phoenix Project (Gene Kim)
- The Goal (Eliyahu Goldratt)

### Support

- **Documentation:** [docs/](docs/)
- **Issues:** [GitHub Issues](https://github.com/ai-masters/sdp/issues)
- **Discussions:** [GitHub Discussions](https://github.com/ai-masters/sdp/discussions)
- **Homepage:** https://github.com/ai-masters/sdp

---

## Release Notes Summary

**Version:** 1.0.0
**Date:** 2026-02-03
**Status:** Stable Release

**What's New:**
- 🎉 First Claude Plugin release
- ✅ Language-agnostic (Python, Java, Go)
- ✅ No installation required (prompts only)
- ✅ Optional Go binary for convenience
- ✅ Comprehensive documentation

**Migration:**
- From Python SDP: See [MIGRATION.md](MIGRATION.md)
- 6-month overlap period (until 2026-08-03)

**Recommendation:**
- ✅ Ready for production use
- ✅ Start new projects with plugin
- ⚠️ Migrate existing projects during overlap period
