# Python SDP Deprecation Notice

**Published:** 2026-02-05
**Effective:** 2026-02-05
**Maintenance Ends:** 2026-08-03

---

## Summary

The Python SDP implementation is **deprecated** and succeeded by the [SDP Plugin](https://github.com/ai-masters/sdp-plugin).

**What this means:**
- This Python repository is in **maintenance mode** (bug fixes only)
- No new features will be added
- Maintenance continues until **2026-08-03** (6 months)
- After that, community-supported only

**Recommended action:** Migrate to the [SDP Plugin](https://github.com/ai-masters/sdp-plugin)

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

| Feature | Python SDP | SDP Plugin | Benefit |
|---------|-----------|------------|---------|
| **Languages** | Python only | Python, Java, Go, any | One tool for all projects |
| **Installation** | `pip install sdp` | Copy prompts | No dependencies |
| **Dependencies** | Python, Poetry, tools | None (optional Go binary) | Zero install friction |
| **Validation** | pytest, mypy, ruff | AI analysis | Language-agnostic |
| **Quality Gates** | Tool-based (fast) | AI-based (flexible) | Understands context |
| **Maintenance** | Python packaging | Go binary + prompts | Simpler deployment |
| **Extensibility** | Python code only | Edit prompts | Customize workflows |

---

## Migration Path

### Good News: Your Work is Safe!

**Your existing workstreams are compatible** with the plugin. No conversion needed.

**Quick migration:**
```bash
# 1. Install plugin (no pip required)
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp
cp -r ~/.claude/sdp/prompts/* .claude/

# 2. Your workstreams work as-is
@build 00-001-01  # Same command!

# 3. Quality checks now use AI (works with any language)
@review F01
```

### Feature Parity

The SDP Plugin achieves **95% feature parity** with Python SDP:

**Same features:**
- ✅ @feature, @design, @build, @review, @deploy skills
- ✅ Multi-agent coordination (11 agents)
- ✅ Beads integration
- ✅ Git hooks
- ✅ Quality gates (coverage, type hints, error handling)
- ✅ Progressive disclosure (@feature skill)

**Better features:**
- ✅ Multi-language support (Python, Java, Go, any)
- ✅ Zero dependencies
- ✅ Language-agnostic validation
- ✅ Simpler customization (edit prompts)

**Tradeoffs:**
- ⚠️ Slower validation (AI vs tools)
- ❌ Missing @oneshot skill (roadmap)
- ❌ Missing checkpoint system (roadmap)
- ❌ Missing Telegram notifications (roadmap)

---

## Documentation

### Migration Guide

**[Complete Migration Guide](migrations/python-sdp-deprecation.md)**

Includes:
- Step-by-step migration instructions
- Feature parity comparison
- Common questions and answers
- Rollback plan
- Real-world migration example

### Feature Parity Document

**[Feature Parity Matrix](migrations/python-sdp-feature-parity.md)**

Includes:
- Detailed feature-by-feature comparison
- Quality gates comparison
- Skills and agents comparison
- Integration comparison
- Language support comparison

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

## Resources

### SDP Plugin

- **Repository:** [https://github.com/ai-masters/sdp-plugin](https://github.com/ai-masters/sdp-plugin)
- **Documentation:** [https://github.com/ai-masters/sdp-plugin/blob/main/docs/TUTORIAL.md](https://github.com/ai-masters/sdp-plugin/blob/main/docs/TUTORIAL.md)
- **Migration Guide:** [https://github.com/ai-masters/sdp-plugin/blob/main/MIGRATION.md](https://github.com/ai-masters/sdp-plugin/blob/main/MIGRATION.md)

### Python SDP (Legacy)

- **Repository:** [https://github.com/fall-out-bug/sdp](https://github.com/fall-out-bug/sdp)
- **Documentation:** [PROTOCOL.md](PROTOCOL.md)
- **Support:** Community (after 2026-08-03)

---

## Frequently Asked Questions

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

### Q: What if I don't like the plugin?

**A:** Rollback is simple:
```bash
# Remove plugin
rm -rf .claude/

# Reinstall Python SDP
pip install sdp

# Your workstreams are unchanged
```

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

**Recommendation:** Migrate to the [SDP Plugin](https://github.com/ai-masters/sdp-plugin) during the 6-month overlap period.

**Thank you** for using Python SDP! We hope you enjoy the plugin's improvements.

---

**Last Updated:** 2026-02-05
**Next Review:** 2026-08-03 (end of maintenance period)
