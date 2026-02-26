# Reference Documentation

Quick lookup guides for SDP commands, configuration, and quality standards.

**CLI reference:** [../CLI_REFERENCE.md](../CLI_REFERENCE.md) — key `sdp` commands.

---

## Contents

- [Commands](#commands)
- [Quality Gates](#quality-gates)
- [Configuration](#configuration)
- [Error Handling](#error-handling)
- [Skills](#skills)

---

## Commands

### SDP CLI Commands

**Core Commands:**
- `@feature` - Create new feature
- `@design` - Plan workstreams
- `@build` - Execute workstream
- `@review` - Quality review
- `@deploy` - Deploy to production
- `@oneshot` - Autonomous execution

**Utility Commands:**
- `/debug` - Systematic debugging
- `@issue` - Bug routing
- `@hotfix` - Emergency fix (P0)
- `@bugfix` - Quality fix (P1/P2)

**See:** [../CLI_REFERENCE.md](../CLI_REFERENCE.md) — SDP CLI commands

---

## Quality Gates

### Mandatory Checks

Every workstream must pass:

1. **TDD** - Tests written before implementation
2. **Coverage ≥80%** - Test coverage threshold
3. **mypy --strict** - Full type hint compliance
4. **ruff** - Code linting
5. **File Size <200 LOC** - Keep code focused
6. **No Bare Exceptions** - Explicit error handling

**See:** [quality-gates.md](quality-gates.md) - Complete quality standards

---

## Configuration

### Quality Gate Configuration

**File:** `quality-gate.toml`

```toml
[coverage]
minimum = 80

[complexity]
max_cc = 10

[file_size]
max_lines = 200

[type_hints]
strict_mode = true
```

**See:** [../PROTOCOL.md](../PROTOCOL.md) — Protocol and config

---

## Error Handling

### SDP Error Framework

Structured errors with:
- Category classification
- Remediation steps
- Documentation links
- Context information

**Error Types:**
- `BeadsNotFoundError` - Task not found
- `CoverageTooLowError` - Coverage below threshold
- `QualityGateViolationError` - Quality gate failed
- `WorkstreamValidationError` - Validation failed
- `ConfigurationError` - Config invalid
- `DependencyNotFoundError` - Dependency missing
- `HookExecutionError` - Hook failed
- `TestFailureError` - Tests failed
- `BuildValidationError` - Build check failed
- `ArtifactValidationError` - Artifact invalid

**See:** [skills.md](skills.md) — Skill contracts and error handling

---

## Skills

### Available Skills

**Feature Development:**
- `feature` - Unified entry point
- `idea` - Requirements gathering
- `design` - Workstream planning
- `build` - Execute workstream
- `review` - Quality check
- `deploy` - Production deployment

**Utilities:**
- `oneshot` - Autonomous execution
- `debug` - Systematic debugging
- `issue` - Bug routing
- `hotfix` - Emergency fix
- `bugfix` - Quality fix

**Internal:**
- `tdd` - TDD enforcement (automatic)

**See:** [skills.md](skills.md) — Skill catalog

---

## Quick Find

### Looking For...

| Need | Doc |
|------|-----|
| Command syntax | [../CLI_REFERENCE.md](../CLI_REFERENCE.md) |
| Quality standards | [build-spec.md](build-spec.md) |
| Skill details | [skills.md](skills.md) |
| Design workflow | [design-spec.md](design-spec.md) |
| Review workflow | [review-spec.md](review-spec.md) |

---

**Version:** SDP v0.9.8 | **Updated:** 2026-02-26
