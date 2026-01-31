# ADR-007: Simplified Workstream Template (Contracts Only)

**Status:** Accepted  
**Date:** 2026-01-30  
**Context:** F032 - SDP Protocol Enforcement Overhaul

## Context

Original workstream template was ~256 lines with:
- Full code examples (not contracts)
- Implementation patterns
- Extensive DO/DON'T sections
- Project-specific conventions

**Problem:** Agents copy-paste code instead of implementing. Workstreams become implementation specs, not execution specs.

## Decision

Adopt simplified 72-line template with:

1. **Contracts only** - function/class signatures with `raise NotImplementedError`
2. **Max 80 lines total** - forces clarity, prevents over-specification
3. **Max 30 lines per code block** - prevents full implementations
4. **Required sections only** - Goal, AC, Contract, Scope, Verification
5. **No DO/DON'T** - patterns belong in CODE_PATTERNS.md, not WS

## Rationale

**Why contracts, not implementations?**
- Forces agent to think, not copy-paste
- Reduces WS file size → AI can read more WS in context
- Separates "what" (WS) from "how" (agent decision)

**Why 80 lines?**
- Fits on one screen without scrolling
- Forces clarity and brevity
- Aligns with "workstream = one atomic task"

**Why remove DO/DON'T?**
- Duplicates CODE_PATTERNS.md and .cursorrules
- Adds noise to WS file
- Should be in project-wide docs, not per-WS

**Why max 30 lines per code block?**
- Prevents hiding full implementations in "examples"
- Contract signatures rarely exceed 30 lines
- If it does, decompose into smaller contracts

## Consequences

**Positive:**
- Shorter WS files → more WS in AI context window
- Agents implement instead of copy-paste
- Clearer separation: WS = spec, agent = implementation
- Faster WS authoring (less to write)

**Negative:**
- Less guidance for complex implementations
  - **Mitigation:** Use @think skill for complex architecture
- May require more back-and-forth with user for clarification
  - **Mitigation:** EnterPlanMode for interactive planning

**Migration:**
- Existing WS using old template still work (backward compatible)
- New WS should use v2 template
- No forced migration required

## Validation

Automated checks:
```bash
sdp ws validate {ws_file}
```

Checks:
- Total lines ≤ 150 (with buffer for future additions)
- Code blocks ≤ 30 lines
- Required sections present
- Large code blocks without `NotImplementedError` trigger warnings

## Alternatives Considered

1. **Keep existing template, add "minimal implementation" guidance**
   - Rejected: Agents ignore guidance, need structural enforcement

2. **Remove code examples entirely**
   - Rejected: Some contracts are complex (async, generics, protocols)

3. **50-line limit instead of 80**
   - Rejected: Too tight for contracts with multiple classes/functions

## References

- Original template: `templates/workstream.md` (256 lines)
- New template: `templates/workstream-v2.md` (72 lines)
- Validator: `src/sdp/validators/ws_template_checker.py`
- Related: ADR-004 (Unified Progressive Consensus)
