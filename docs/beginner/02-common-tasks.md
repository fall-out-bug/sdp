# Common SDP Tasks

Quick reference for common SDP workflows and tasks.

---

## Table of Contents

- [Feature Development](#feature-development)
- [Bug Fixes](#bug-fixes)
- [Quality & Testing](#quality--testing)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)

---

## Feature Development

### Create a New Feature

```bash
# Start with feature command
@feature "Add user authentication"

# Plan workstreams
@design idea-user-auth

# Execute first workstream
@build WS-001-01

# Continue with remaining workstreams
@build WS-001-02
@build WS-001-03

# Review feature quality
@review F001

# Deploy to production
@deploy F001
```

**Or execute autonomously:**
```bash
@oneshot F001
```

**See Also:**
- [01-first-feature.md](01-first-feature.md) - Hands-on tutorial
- [PROTOCOL.md](../../PROTOCOL.md) - Complete specification

---

### Check Workstream Status

```bash
# Check specific workstream
sdp tier metrics WS-001-01

# Check all workstreams
sdp tier metrics

# View workstream index
cat docs/workstreams/INDEX.md
```

---

## Bug Fixes

### Report and Fix Bugs

```bash
# Report bug (routes to appropriate fix command)
@issue "Login fails on Firefox"

# For P0 (emergency) issues
@hotfix "Critical security vulnerability"

# For P1/P2 (quality) issues
@bugfix "Incorrect totals calculation"
```

**Severity Classification:**
- **P0** - Critical security, data loss, production down
- **P1** - Major functionality broken
- **P2** - Minor issues, workarounds available
- **P3** - Cosmetic, nice to have

**See Also:**
- [GLOSSARY.md](../GLOSSARY.md#severity-classification) - Severity definitions
- [/debug skill](../../.claude/skills/debug/SKILL.md) - Systematic debugging

---

## Quality & Testing

### Run Quality Checks

```bash
# Check code quality
@review F001

# Run tests
pytest tests/unit/ -v

# Check coverage
pytest --cov=src/sdp --cov-report=term-missing

# Type checking
mypy src/sdp/ --strict

# Linting
ruff check src/sdp/
```

---

### Debug Failing Tests

```bash
# Use debug skill
/debug "Test fails unexpectedly"

# Or use debug runbook
# See: docs/runbooks/debug-runbook.md
```

**See Also:**
- [03-troubleshooting.md](03-troubleshooting.md) - Common issues and solutions
- [runbooks/debug-runbook.md](../runbooks/debug-runbook.md) - Debug workflow

---

## Deployment

### Deploy to Production

```bash
# Deploy feature
@deploy F001

# This will:
# 1. Verify quality gates pass
# 2. Create git tag
# 3. Push to production branch
# 4. Trigger deployment pipeline
```

**Pre-deployment Checklist:**
- [ ] All workstreams completed
- [ ] All tests passing
- [ ] Coverage ≥80%
- [ ] Code review approved
- [ ] Documentation updated

**See Also:**
- [github-integration/SETUP.md](../github-integration/SETUP.md) - GitHub setup
- [PROTOCOL.md#deployment](../../PROTOCOL.md#deployment) - Deployment workflow

---

## Troubleshooting

### "I'm stuck. Where do I start?"

1. **[Glossary](../GLOSSARY.md)** - Look up unfamiliar terms
2. **[Site Map](../SITEMAP.md)** - Find specific documentation
3. **[03-troubleshooting.md](03-troubleshooting.md)** - Common issues
4. **[Runbooks](../runbooks/)** - Step-by-step procedures

---

### "Command fails with error"

1. Check error message for specific guidance
2. See [03-troubleshooting.md](03-troubleshooting.md) for common errors
3. Use `/debug` for systematic debugging
4. Check [GitHub Issues](https://github.com/fall-out-bug/sdp/issues)

---

### "Tests are failing"

```bash
# Run specific test with output
pytest tests/unit/test_module.py -v

# Run with coverage to see what's missing
pytest tests/unit/test_module.py --cov=module --cov-report=term-missing

# Use debug skill
/debug "Test fails unexpectedly"
```

**See Also:**
- [runbooks/test-runbook.md](../runbooks/test-runbook.md) - Test workflow
- [runbooks/debug-runbook.md](../runbooks/debug-runbook.md) - Debug workflow

---

## Quick Reference

### Essential Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `@feature` | Create feature | `@feature "Add comments"` |
| `@design` | Plan workstreams | `@design idea-comments` |
| `@build` | Execute workstream | `@build WS-001-01` |
| `@review` | Quality check | `@review F001` |
| `@deploy` | Deploy to production | `@deploy F001` |
| `@oneshot` | Autonomous execution | `@oneshot F001` |
| `/debug` | Debug issues | `/debug "Test fails"` |
| `@issue` | Report bug | `@issue "Login fails"` |

---

### Key Files

| File | Purpose |
|------|---------|
| **START_HERE.md** | Welcome page |
| **README.md** | Project overview |
| **PROTOCOL.md** | Complete specification |
| **GLOSSARY.md** | 150+ term reference |
| **SITEMAP.md** | Documentation index |

---

### Workstream Status

| Status | Meaning |
|--------|---------|
| **backlog** | Planned, not started |
| **active** | Currently in progress |
| **completed** | Finished and verified |
| **blocked** | Waiting for dependency |

---

## Learning Paths

### New to SDP?
1. [00-quick-start.md](00-quick-start.md) - 5-minute overview
2. [01-first-feature.md](01-first-feature.md) - Hands-on tutorial
3. This file - Common tasks
4. [03-troubleshooting.md](03-troubleshooting.md) - Handle issues

### Ready to Dive Deeper?
1. [PROTOCOL.md](../../PROTOCOL.md) - Complete specification
2. [reference/](../reference/) - Lookup documentation
3. [internals/](../internals/) - Architecture and extending

---

**Need Help?**
- [Glossary](../GLOSSARY.md) - Look up terms
- [Site Map](../SITEMAP.md) - Find documentation
- [GitHub Issues](https://github.com/fall-out-bug/sdp/issues) - Report issues

---

**Version:** SDP v0.5.0
**Updated:** 2026-01-29
