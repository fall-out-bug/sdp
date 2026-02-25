# Architectural Decisions

**Generated:** 2026-02-25

**Total:** 3 decisions

## 1. Evidence file lock platform

Evidence log uses `syscall.Flock` for inter-process safety. **Requires UNIX (macOS/Linux).** Windows is not supported. See docs/plans/2026-02-25-audit-gaps-design.md Risks.

## 2. in-toto predicate snake_case

Coding-workflow predicate uses **snake_case** by design (D2 in audit-gaps-design). lowerCamelCase is only required for registry; we keep snake_case for consistency with existing tooling.

## 3. JWT with refresh tokens

**Date:** 2026-02-09 23:42:47
**Type:** explicit
**Maker:** user

### Question

How to handle auth?

### Decision

JWT with refresh tokens

### Rationale

Stateless, scales horizontally

### Alternatives Considered

- Sessions
- OAuth2 only

---

