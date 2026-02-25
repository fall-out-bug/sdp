# Architectural Decisions

**Generated:** 2026-02-25

**Total:** 5 decisions

## 1. JWT with refresh tokens

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

## 2. Evidence Layer — Singleton Writer + flock (D1)

**Date:** 2026-02-25
**Type:** explicit
**Maker:** audit-gaps-design

### Question

How to prevent evidence loss and inter-process race in hash chain?

### Decision

Singleton Writer via `sync.Once`, `syscall.Flock` on `.lock` file, `slog.Error` on Emit failure. Append-to-file with re-read-under-lock.

### Rationale

In-process mutex is per-instance; `NewWriter` per call caused inter-process races. Flock provides advisory locking. Singleton holds state in memory.

### Alternatives Considered

- Per-call Writer (current) — races
- Advisory lock only — no singleton
- Sync-by-default Emit — deferred to v1.0 (breaking)

---

## 3. in-toto Schema — Pragmatic v1 (D2)

**Date:** 2026-02-25
**Type:** explicit
**Maker:** audit-gaps-design

### Question

How to align in-toto predicate with supply-chain standards?

### Decision

`StatementInTotoV1`, create `coding-workflow-predicate.schema.json`, add `version: "1.0"`, `digestSet: {"sha256":"..."}`. Keep snake_case — document as intentional.

### Rationale

v0.1 deprecated; goreleaser references predicate schema. lowerCamelCase only needed for registry; cosmetic churn deferred.

### Alternatives Considered

- Full lowerCamelCase rename — too much churn
- Keep v0.1 — deprecated

---

## 4. Inverted Architecture — Enforcement in CLI/Schemas (D3)

**Date:** 2026-02-25
**Type:** explicit
**Maker:** audit-gaps-design

### Question

Where does SDP enforce correctness: prompts vs tooling?

### Decision

Enforcement in CLI, schemas, hooks. Prompts focus on judgment; validation via JSON schemas (review-verdict, ws-verdict). Hallucination detection in guard/constraints.

### Rationale

Audit criticisms about prompt engineering often belong in tooling. Output validation, retry, token budget — already in CLI. Chain-of-thought and few-shot — selective prompt improvements.

### Alternatives Considered

- Prompt-only enforcement — brittle
- No schemas — no machine validation

---

## 5. Go Code Quality — Prioritized Plan (D4)

**Date:** 2026-02-25
**Type:** explicit
**Maker:** audit-gaps-design

### Question

How to address audit findings in Go code?

### Decision

P0: Mock removal (WorkstreamRunner interface), evidence slog.Error + singleton Writer. P1: Context propagation, wire coverage checker. P2: Config timeouts, slog migration.

### Rationale

Credibility fixes first; standards and contracts second; polish last.

### Alternatives Considered

- Big-bang refactor — risky
- Ignore P0 — credibility kill

---

