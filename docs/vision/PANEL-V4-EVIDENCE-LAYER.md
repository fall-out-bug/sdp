# Panel Review V4: The Evidence Layer

> Four experts reviewed VISION v3.0, ROADMAP v3.0, FEATURES v3.0.
> Question: Is the evidence layer real or overthinking?

---

## Panel

| Expert | Perspective | Verdict |
|--------|------------|---------|
| **Marty Cagan** (SVPG) | Product strategy | **Build it — as a feature, not a religion** |
| **Andrej Karpathy** (ex-Tesla AI) | AI/ML practitioner | **Rethink it** |
| **Mitchell Hashimoto** (HashiCorp) | Infrastructure tooling | **Build it** |
| **Charity Majors** (Honeycomb) | Observability/ops | **Build it** |

**Score: 3 Build, 1 Rethink, 0 Kill.**

---

## Consensus Points (All Four Agree)

1. **The existing system is real.** Not a prototype, not a deck. Working engine with paying customers.
2. **The provenance gap is real.** AI-generated code without records is a growing problem.
3. **The kill criteria need fixing.** "500 runs, catch rate < 5%" is either too vague (what is "catch"?), measures the wrong thing (MTTR not catch rate), or is toothless (who enforces?).
4. **Six weeks is aggressive** for the full Phase 1 scope.

---

## Key Disagreements

### Forensics vs Verification (Karpathy vs Majors)

**Karpathy:** "Forensics > Verification is wrong. Catching problems before deployment beats forensics after the incident. Building forensics on top of verification admits your verification isn't good enough."

**Majors:** "You're saying the point isn't to prevent bad code, it's to reconstruct what happened after it breaks. That's observability thinking. I like it."

**Resolution:** Both are right. Verification prevents. Forensics reconstructs. The question is where to invest marginal effort. Karpathy says: make verification better. Majors says: build the reconstruction tool.

### Provenance Value (Karpathy vs Hashimoto)

**Karpathy:** "Provenance is metadata theater. Knowing 'this was Claude 3.5 with temperature 0.3' doesn't change the remediation decision. You'd just look at the code."

**Hashimoto:** "Provenance as a separate subsystem is wrong — fold it into the evidence log as an event type. But the data itself is valuable."

**Resolution:** Provenance alone is low-signal. Provenance *linked to* verification output and spec is higher-signal. Don't oversell model identity as the key insight.

### Protocol vs Product (Cagan vs all)

**Cagan:** "You don't have a standards play. You have an enterprise SaaS with a compliance moat. Standards emerge from shared pain across an industry. You have three customers."

**Hashimoto/Majors:** Didn't challenge the protocol bet directly but focused on practical value over standardization.

**Resolution:** Build the schema for internal consistency first. If it becomes a standard, that's a result, not a goal.

---

## Actionable Feedback (Synthesized)

### Must Fix Before Building

1. **Schema sprawl** (Hashimoto): Existing schemas in `schema/`, `docs/schema/`, and frontmatter are already out of sync. Consolidate before adding a new protocol schema.
2. **Evidence storage** (Majors): `.sdp/log/` gitignored by default = evidence lost when laptops are wiped or CI runners are ephemeral. Must be committable or durable.
3. **Hash chain honesty** (Hashimoto): Hash-chaining detects corruption, doesn't prevent tampering. Don't call it "tamper-proof" in front of the bank's security team.
4. **Kill criteria** (all four): Redefine. Cagan: "who enforces?" Karpathy: "measure what the evidence told you that git blame didn't." Majors: "MTTR with vs without SDP evidence."

### Architecture Changes

5. **Fold provenance into evidence log** (Hashimoto): One log, multiple event types, one reader. Not `.sdp/provenance/` + `.sdp/log/`.
6. **`sdp incident` → `sdp log show --trace`** (Hashimoto): Incident implies runtime. Trace implies forensics. Log as the single queryable artifact.
7. **Observability bridge** (Majors): Evidence log needs to connect to production — deploy markers, OpenTelemetry span attributes for AI-generated code paths. Without this, it's a filing cabinet nobody opens during fires.
8. **Runtime context** (Majors): Evidence captures generation-time context but not runtime context (blast radius, feature flags, rollback path).

### Strategic Changes

9. **Ship as feature, not thesis** (Cagan): Evidence layer is what the bank needs to go to production. Frame it that way. Drop "neutral ledger" and "five bets" from customer-facing materials.
10. **Thin audit trail first** (Karpathy): Model ID + pass/fail + timestamp in one week. Full forensic chain can come later if the thin version proves useful.
11. **Invest in verification** (Karpathy): The existing TDD pipeline is the real product. Make it better before building a forensics layer for when it fails.
12. **Diff-level provenance** (Majors): Track which *lines* were AI-generated vs human-edited. Currently only at unit/commit level.

### Timeline Reality

13. **Scope Phase 1 tighter** (Hashimoto): Schema + log + incident query in 6 weeks is realistic. Provenance instrumentation of all 19 skills → Phase 2.
14. **One customer, one deliverable** (Cagan): What does the bank need to go to production? Ship exactly that in 4 weeks.

---

## Adjusted Recommendations

### Phase 1 (4 weeks, not 6)

1. **Consolidate schemas** (week 1) — fix existing sprawl first
2. **Evidence log** (weeks 1-3) — single `.sdp/log/`, JSON, hash-chained, provenance as event type
3. **`sdp incident` / trace** (weeks 2-4) — query the log from a commit
4. **Instrument `@build` only** (weeks 3-4) — one skill, not all 19

### Phase 1.5 (weeks 5-8)

5. Instrument remaining skills
6. Compliance design doc
7. `sdp plan`/`sdp apply` thin wrappers
8. Deploy markers / OTel bridge design

### Revised Kill Criteria

> After 100 incidents involving AI-generated code: if SDP evidence did not reduce mean time to root cause (MTTR) compared to git-blame-only investigation → rethink the approach.

---

*Panel Review V4 — February 2026*
*Verdict: Build the evidence layer. But smaller, faster, and more honest.*
