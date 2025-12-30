# Core Concepts

This document explains the key concepts behind AI-assisted software development workflows.

## Why Structure Matters

AI coding assistants are powerful, but without structure they can:
- Miss requirements
- Make inconsistent architectural decisions
- Skip testing
- Create technical debt

Structure provides guardrails that improve quality without slowing you down.

## The Three Modes

| Mode | Complexity | When to Use |
|------|------------|-------------|
| [Solo](modes/solo/) | Low | Quick fixes, small features |
| [Structured](modes/structured/) | Medium | Features needing documentation |
| [Multi-Agent](modes/multi-agent/) | High | Large projects, teams, compliance |

## Key Concepts

### 1. Roles in Software Development

Even when one person (or AI) does everything, these responsibilities exist:

| Role | Question | Output |
|------|----------|--------|
| **Analyst** | What do we need? | Requirements |
| **Architect** | How should it work? | Design, ADRs |
| **Developer** | How do we build it? | Code, Tests |
| **QA** | Does it work correctly? | Test results |
| **DevOps** | How do we deliver it? | Deployment |

In **Solo mode**, one AI handles all roles implicitly.
In **Structured mode**, you guide the AI through phases.
In **Multi-Agent mode**, different sessions handle different roles.

See [concepts/roles/](concepts/roles/) for details.

### 2. Artifacts

Artifacts are the outputs that document decisions and enable collaboration:

| Artifact | Purpose | When Needed |
|----------|---------|-------------|
| `spec.md` | What to build | Structured, Multi-Agent |
| `design.md` | How to build | Structured, Multi-Agent |
| `ADR` | Why decisions were made | Any architectural decision |
| `test-report.md` | Quality verification | Formal projects |

See [concepts/artifacts/](concepts/artifacts/) for details.

### 3. Architecture Decision Records (ADR)

ADRs document WHY you made architectural choices:

```markdown
# ADR-0001: Use PostgreSQL for persistence

## Status
Accepted

## Context
We need a database. Options: PostgreSQL, MySQL, MongoDB.

## Decision
PostgreSQL for ACID compliance and JSON support.

## Consequences
+ Strong data integrity
+ Good tooling ecosystem
- Requires more setup than SQLite
```

See [concepts/adr/](concepts/adr/) for template and examples.

### 4. Clean Architecture

Dependencies point inward:

```
┌─────────────────────────────────────┐
│         Presentation                │  ← Controllers, Views, API
├─────────────────────────────────────┤
│         Infrastructure              │  ← Database, External APIs
├─────────────────────────────────────┤
│         Application                 │  ← Use Cases, Services
├─────────────────────────────────────┤
│         Domain                      │  ← Entities, Business Rules
└─────────────────────────────────────┘
```

**Rule**: Inner layers never depend on outer layers.

See [concepts/clean-architecture/](concepts/clean-architecture/) for details.

### 5. Quality Gates

Non-negotiable rules that prevent common problems:

| Gate | Rule | Why |
|------|------|-----|
| No silent failures | `except: pass` forbidden | Errors must be visible |
| Test coverage | ≥80% for new code | Prevents regressions |
| Architecture | No layer violations | Maintains structure |
| Documentation | English only | Team accessibility |

## Choosing the Right Mode

```
Start here
    │
    ▼
Is the task < 2 hours and < 10 files?
    │
    ├─ Yes → Solo Mode
    │
    └─ No
        │
        ▼
    Do you need documentation or ADRs?
        │
        ├─ Yes → Structured Mode
        │
        └─ No
            │
            ▼
        Parallel work or audit trail needed?
            │
            ├─ Yes → Multi-Agent Mode
            │
            └─ No → Structured Mode
```

## Learning Path

1. **Start with Solo** - Learn basic AI prompting
2. **Try Structured** - Learn phased development
3. **Explore Multi-Agent** - Understand team dynamics

## Next Steps

- [Solo Mode Guide](modes/solo/)
- [Structured Mode Guide](modes/structured/)
- [Multi-Agent Mode Guide](modes/multi-agent/)
- [Concepts Deep Dive](concepts/)
