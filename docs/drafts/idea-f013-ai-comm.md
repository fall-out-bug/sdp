# AI-Human Communication Enhancement

> **Feature ID:** F013
> **Status:** Implemented
> **Created:** 2026-01-26
> **Links:** [Plan](../plans/2025-01-26-ai-human-comm-design.md)

---

## Problem

SDP lacks structured AI-human communication patterns:
- No unified entry point for feature development
- No machine-readable intent specification
- Workstreams lack execution metadata (prerequisites, graphs, steps)
- No systematic debugging guidance
- TDD discipline not enforced
- No deep thinking framework for complex decisions

## Users

1. **Solo developers** building AI-assisted projects
2. **Small teams** (2-5 engineers) with AI workflows
3. **DevOps engineers** integrating SDP into CI/CD

## Success Criteria

- [ ] Time from idea to running code: <1 hour
- [ ] New user onboarding: <30 min to first WS
- [ ] All features have machine-readable intent.json
- [ ] Workstreams include execution graphs
- [ ] TDD cycle enforced automatically

## Goals

### Primary Goals

1. **Unified entry point** — `/feature` command for progressive disclosure workflow
2. **Intent schema** — Machine-readable specification with validation
3. **Enhanced workstreams** — Prerequisites, execution graphs, file inventory
4. **TDD enforcement** — Red→Green→Refactor cycle as internal discipline
5. **Systematic debugging** — Scientific method for root cause analysis
6. **Deep thinking** — Parallel expert agents for complex decisions

### Non-Goals

- Real-time collaboration (multiplayer)
- Enterprise SSO
- Language-agnostic (Python-first, extensible)

## Technical Approach

### Architecture

Layered architecture for AI-human communication:

```
Layer 1: Foundation (can parallelize)
├── intent schema + validation
├── /tdd skill
└── /debug skill

Layer 2: Core Skills (depends on Layer 1)
├── /feature skill
└── @idea enhancement

Layer 3: Planning (depends on Layer 2)
├── @design enhancement
└── /build skill

Layer 4: Execution (depends on Layer 3)
├── @oneshot enhancement
└── /think skill
```

### Storage

- JSON schema: `docs/schema/intent.schema.json`
- Intent files: `docs/intent/{slug}.json`
- Product vision: `PRODUCT_VISION.md` (root)

### Failure Mode

Graceful degradation — if schema file missing, fallback to embedded schema

### Auth Method

None (local CLI)

## Implementation

### Modules Created

| Module | Purpose | LOC | Coverage |
|--------|---------|-----|----------|
| `src/sdp/schema/` | Intent validation | 201 | 86% |
| `src/sdp/tdd/` | TDD cycle runner | 107 | 100% |
| `src/sdp/feature/` | Product vision | 117 | 98% |
| `src/sdp/design/` | Dependency graph | 121 | 98% |

### Skills Created/Enhanced

| Skill | Status | Purpose |
|-------|--------|---------|
| `/feature` | NEW | Unified entry point |
| `/tdd` | NEW | TDD discipline (internal) |
| `/debug` | REIMPL | Systematic debugging |
| `/think` | NEW | Deep thinking with parallel agents |
| `/build` | NEW | Workstream execution |
| `@idea` | ENHANCED | + vision + schema |
| `@design` | ENHANCED | + execution graphs |
| `@oneshot` | ENHANCED | + checkpoint/resume |

## Tradeoffs

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| DX vs Control | Prioritize DX | Friction kills adoption |
| Speed vs Quality | Both | TDD + quality gates |
| Simple vs Expressive | Progressive | Simple entry, power when needed |

## Open Questions

None — all requirements implemented.

---

**Next Step:** Transition to `@design` for workstream decomposition (completed retroactively).
