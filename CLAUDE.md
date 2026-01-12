# Claude Code Integration Guide

Quick reference for using this Spec-Driven Protocol (SDP) repository with Claude Code.

> **📝 Meta-note:** This guide was written with AI assistance (Claude Sonnet 4.5). The workflow is based on real development experience.

## TL;DR

Use **skills** to execute SDP commands:

```
@idea "Add user authentication"
@design idea-user-auth
@build WS-001-01
@review F01
@deploy F01
```

## Available Skills

| Skill | Purpose | Example |
|-------|---------|---------|
| `@idea` | Requirements gathering | `@idea "Add payment processing"` |
| `@design` | Create workstreams | `@design idea-payments` |
| `@build` | Execute workstream | `@build WS-001-01` |
| `@review` | Quality check | `@review F01` |
| `@deploy` | Production deployment | `@deploy F01` |
| `@issue` | Debug and route bugs | `@issue "Login fails on Firefox"` |
| `@hotfix` | Emergency fix (P0) | `@hotfix "Critical API outage"` |
| `@bugfix` | Quality fix (P1/P2) | `@bugfix "Incorrect totals"` |
| `@oneshot` | Autonomous execution | `@oneshot F01` |

Skills are defined in `.claude/skills/{name}/SKILL.md`

## Quick Reference

### First Time Setup

1. **Read core docs:**
   - [README.md](README.md) — Overview and quick start
   - [PROTOCOL.md](PROTOCOL.md) — Full SDP specification
   - [RULES_COMMON.md](RULES_COMMON.md) — Common rules

2. **Understand key concepts:**
   - **Workstream (WS)**: Atomic task, one-shot execution
   - **Feature**: 5-30 workstreams
   - **Release**: 10-30 features

3. **Review quality gates:**
   - Files < 200 LOC
   - Coverage ≥80%
   - No `except: pass`
   - Full type hints

### Typical Workflow

```bash
# 1. Gather requirements
@idea "User can reset password via email"

# 2. Design workstreams (creates WS-XXX-01, WS-XXX-02, etc.)
@design idea-password-reset

# 3. Execute each workstream
@build WS-001-01
# Claude will show progress tracking:
#   [in_progress] Pre-build validation
#   [pending] Write failing test (Red)
#   [pending] Implement minimum code (Green)
#   [pending] Refactor implementation
#   ... (updates in real-time)

@build WS-001-02
# ... or use autonomous mode:
@oneshot F01

# 4. Review quality
@review F01

# 5. Deploy to production
@deploy F01
```

### Progress Tracking

When using `@build`, Claude Code automatically tracks progress using TodoWrite:

```markdown
User: @build WS-060-01

Claude:
→ Creating todo list...
  ✓ [in_progress] Pre-build validation
  • [pending] Write failing test (Red)
  • [pending] Implement minimum code (Green)
  • [pending] Refactor implementation
  • [pending] Verify Acceptance Criteria
  • [pending] Run quality gates
  • [pending] Append execution report
  • [pending] Git commit

→ Reading WS file...
  ✓ [completed] Pre-build validation
  ✓ [in_progress] Write failing test (Red)
  • [pending] Implement minimum code (Green)
  ...

→ Test created, running pytest... FAILED (expected)
  ✓ [completed] Write failing test (Red)
  ✓ [in_progress] Implement minimum code (Green)
  ...

→ Implementation done, running pytest... PASSED
  ✓ [completed] Implement minimum code (Green)
  ✓ [in_progress] Refactor implementation
  ...

[All steps complete]
  ✓ All tasks completed
```

This provides real-time visibility into WS execution progress.

### File Structure Reference

```
project/
├── docs/
│   ├── drafts/           # @idea outputs here
│   ├── workstreams/
│   │   ├── backlog/      # @design outputs here
│   │   ├── in_progress/  # @build moves here
│   │   └── completed/    # @build finalizes here
│   └── specs/            # Feature specifications
├── prompts/commands/     # Skill instructions
├── .claude/
│   ├── skills/           # Skill definitions
│   ├── agents/           # Multi-agent mode (advanced)
│   └── settings.json     # Claude Code settings
└── hooks/                # Git hooks for validation
```

## Key Principles (Quick)

- **SOLID, DRY, KISS, YAGNI** — see [docs/PRINCIPLES.md](docs/PRINCIPLES.md)
- **Clean Architecture** — Domain ← App ← Infra ← Presentation
- **TDD** — Tests first (Red → Green → Refactor)
- **AI-Readiness** — Small files, low complexity, typed

## Validation

### Pre-build Check
```bash
hooks/pre-build.sh WS-001-01
```

### Post-build Check
```bash
hooks/post-build.sh WS-001-01 project.module
```

### Manual Validation
```bash
python scripts/validate.py docs/workstreams/backlog/
```

## Quality Gates (Enforced)

| Gate | Requirement |
|------|-------------|
| **AI-Readiness** | Files < 200 LOC, CC < 10, type hints |
| **Clean Architecture** | No layer violations |
| **Error Handling** | No `except: pass` |
| **Test Coverage** | ≥80% |
| **No TODOs** | All tasks completed or new WS |

## Forbidden Patterns

❌ `except: pass` or bare exceptions  
❌ Time-based estimates  
❌ Layer violations  
❌ Files > 200 LOC  
❌ TODO without followup WS  
❌ Coverage < 80%

## Required Patterns

✅ Type hints everywhere  
✅ Tests first (TDD)  
✅ Explicit error handling  
✅ Clean architecture boundaries  
✅ Conventional commits

## Troubleshooting

### Skill not found
Check `.claude/skills/{name}/SKILL.md` exists

### Validation fails
Run `hooks/pre-build.sh {WS-ID}` to see specific issues

### Workstream blocked
Check dependencies in `docs/workstreams/backlog/{WS-ID}.md`

### Coverage too low
Run `pytest --cov --cov-report=term-missing` to identify gaps

## Advanced: Multi-Agent Mode

For complex features, use multi-agent orchestration:

```bash
@orchestrator F01  # Coordinates all agents
```

Agents defined in `.claude/agents/`:
- `planner.md` — Breaks features into workstreams
- `builder.md` — Executes workstreams
- `reviewer.md` — Quality checks
- `deployer.md` — Production deployment
- `orchestrator.md` — Coordinates workflow

## Configuration

See `.claude/settings.json` for:
- Custom Git hooks
- Validation scripts
- Tool integrations

## Resources

| Resource | Purpose |
|----------|---------|
| [PROTOCOL.md](PROTOCOL.md) | Full specification |
| [docs/PRINCIPLES.md](docs/PRINCIPLES.md) | Core principles |
| [CODE_PATTERNS.md](CODE_PATTERNS.md) | Code patterns |
| [MODELS.md](MODELS.md) | Model recommendations |
| [prompts/commands/](prompts/commands/) | Skill instructions |

---

**Version:** SDP 0.3.0  
**Claude Code Version:** 0.3+  
**Mode:** Skill-based, one-shot execution
