# SDP Product Vision

> *What's actually new. What already works. No rebrand.*

---

## What Exists (and Works)

SDP already has a functioning engine. This is not aspirational — it's running in production today as a Claude Code skill system.

**Decomposition.** NL → atomic workstreams with dependency graph. Topological sort, parallel dispatch, circuit breakers, checkpoint recovery. Written in Go. Works.

**Verification.** TDD pipeline (red → green → refactor). Coverage enforcement ≥80%. Contract validation. Type checking. Static analysis via semgrep. Works.

**Multi-agent orchestration.** 19 agent roles: idea gathering, architecture design, build execution, adversarial review, deployment. Synthesis engine for conflict resolution. Works.

**Two modes.** `ship` (autonomous) and `drive` (human-in-the-loop). Progressive disclosure — start simple, deepen when needed. Works.

**Real-world traction.** Used in production-style workflows with strict quality and traceability requirements.

**This is not a prototype.** It's a working system with real operational usage.

---

## What's Missing (The Real Delta)

Everything above generates code, verifies it, and ships it. But it leaves no trace. When something breaks at 3 AM, you have:

- A git blame
- A PR that got approved
- No idea which model wrote it, from what spec, or what verification actually ran

**The gap is not in generation. The gap is in evidence.**

### 1. Model Provenance

**Today:** Code gets generated. Nothing records which model, which version, what parameters, what prompt produced it.

**Needed:** Every AI generation records:
- Model name and version (e.g., `claude-sonnet-4-20250514`)
- Prompt hash (not the prompt — privacy)
- Temperature and parameters
- Timestamp
- Which spec it was generated against
- Who initiated the generation

**Why this matters:** When a model version gets deprecated or a vulnerability is found in model-generated patterns, you need to know which code was affected. Without provenance, you're grepping git logs and guessing.

### 2. Evidence Log

**Today:** `@build` runs TDD, `@review` runs adversarial review. Results go into markdown reports in workstream files. Human-readable, but not machine-parseable, not linked, not queryable.

**Needed:** Single structured log in `.sdp/log/events.jsonl`. Provenance is not a separate subsystem — it's the `generation` event type. One log, multiple event types, one reader.

- JSON records following a published schema
- Four event types: `plan`, `generation` (with provenance), `verification`, `approval`
- Hash-chained for integrity (each record references the previous — corruption detection, **not** tamper-proof)
- Actual verification output (real pytest stdout, real mypy output — not "tests passed")
- Decision log (in drive mode: what the human decided and why)
- Committed to repo by default (evidence must survive laptop wipes and ephemeral CI runners)

**Why this matters:** "Tests passed" is an assertion. The actual pytest output is evidence. The difference matters when someone asks "what exactly was tested?" and you need the answer in 30 seconds, not 30 minutes of archaeology.

### 3. Forensic Trace (`sdp log trace`)

**Today:** Something breaks in production. You read the git log. You find a commit. You look at the PR. You find a workstream file that's already been archived. Maybe. If someone didn't clean up.

**Needed:** `sdp log trace <commit>` — one command that walks the evidence chain backwards from any commit:
- Which model generated this code
- From what specification
- What verification ran (with actual output)
- Who approved it
- When, and what they saw at the time

"Trace" not "incident" — because the query is forensic (reconstructing the chain), not operational (managing runtime). The log is the artifact; trace is the query.

**Why this matters:** This is the tool Georgy and Alexander need at 3 AM. The reconstruction command. The thing that turns a panicked incident into a structured investigation.

### 4. Protocol Schema

**Today:** SDP is a collection of skills and Go code. The "protocol" is implicit — embedded in skill instructions and agent behavior. There's no formal definition that another tool could implement.

**Needed:** A machine-readable JSON Schema that defines:
- `plan` event: what was intended
- `apply` event: what was generated and verified
- `evidence` record: the linked chain
- `incident` query: how to traverse the chain

**Why this matters:** Without a schema, SDP is a product. With a schema, SDP is a standard. Products get replaced. Standards get adopted.

---

## What Stays

The engine is good. The infrastructure is good.

- **Skills stay.** `@idea`, `@design`, `@build`, `@review`, `@deploy`.
- **Beads stays.** Issue tracking across sessions.
- **Decomposition stays.** NL → workstreams → dependency graph → parallel execution.
- **TDD pipeline stays.** Red → green → refactor.
- **Go engine stays.** Dispatcher, circuit breaker, synthesis engine.
- **Two modes stay.** Ship and drive.

---

## What Needs to Change in the Existing System

### 5. Acceptance Test Gate

The engine works. But it produces non-working products.

Real case: SDP built an ad server — 7,149 LOC, 88% coverage, Clean Architecture, all quality gates green. Basic features didn't work. `IsWithinBudget()` was never called. Weighted random always returned max weight. Frequency capping was missing. A vibe-coded alternative did the same in 1,005 LOC and it worked from `docker compose up`.

The difference: vibe-coders run the app after every change. 5-minute feedback loop. SDP builds 25 workstreams across architecture layers, then discovers nothing works. Days-long feedback loop.

**The fix:** After every `@build`, run the app and check if it does what it's supposed to. Not unit tests — an actual smoke test. Start the server, make a request, check the response. If the app doesn't work — the build failed, regardless of coverage.

This is the vibe-coder's feedback loop, formalized into the protocol.

Architecture depth, scope, decomposition order — those are user choices. The user told `@idea` to go deeper, and it did. That's a conversation, not a bug. But "does it work?" is not a choice. That's a gate.

---

## The Thesis

**Accountability for AI-generated code. Starting with "does it actually work?"**

The protocol must answer three questions:
1. **Does it work?** — acceptance test gate
2. **What happened?** — evidence log, provenance, forensic trace
3. **Who else is affected?** — scope collision, shared contracts, cross-branch integration

Without #1, #2 is just a detailed record of how you produced non-working software.

SDP is the neutral ledger. It records:
1. What was specified (the intent)
2. What was generated (model + params + output)
3. What was verified (actual tool output, not summaries)
4. **Whether it actually works** (acceptance test — new)
5. Who approved (the human decision)

---

## The UX Surface

Three entry points. That's it.

**`sdp plan`** — Show what you're about to do. Under the hood, this is `@idea` + `@design`. Decompose, show the units, estimate cost.

**`sdp apply`** — Do it. Under the hood, this is `@build` + `@review`. Each unit: generate → verify → record. **The "record" part is new.** Every generation emits a `generation` event. Every verification emits actual output. Every approval gets logged. All into one `.sdp/log/`.

**`sdp log trace`** — Trace. **This is entirely new.** Walk the evidence chain backwards from any commit.

The 19 skills are the implementation. The three commands are the interface. Users see `plan/apply/log`. Skills do the work. Evidence gets recorded at every step.

---

## The Five Bets

1. **AI will write most code by 2028.** Trajectory, not prediction.
2. **"What happened?" will be mandatory.** Courtrooms, audits, incident reviews will demand records.
3. **Provenance is permanent.** Regardless of model capability — the record of what generated what never stops being valuable.
4. **Forensics > Verification.** "Did it pass?" is table stakes. "What happened?" is the moat.
5. **Protocol > Product.** A schema that others implement beats a closed tool.

---

## Who This Is For

**Code that costs more to fix than to write.**

| Segment | What They Need |
|---------|---------------|
| Fintech | Forensic proof that payment code was verified, by which model |
| Healthcare | HIPAA audit trail for AI code provenance |
| Infrastructure | `terraform plan` for AI code — see before you execute |
| High-assurance SaaS | Policy: "all AI PRs need evidence chain" |
| Regulated industries | SOC2/DORA/ISO 27001 compliance export |

**Not for:** Landing pages, MVPs, disposable code. If rewriting is cheaper than recording — just ship.

---

## Kill Criteria

> **Primary:** After 100 incidents involving AI-generated code: if SDP evidence did not reduce MTTR compared to git-blame-only investigation → rethink the evidence layer.

> **Secondary:** After 500 SDP runs: if verification catch rate < 5% AND post-merge defect rate equals baseline → kill the verification product entirely.

Specific. Testable. Honest. Measured monthly by product owner.

---

## Summary: What Changes, What Stays

| Aspect | Status | Notes |
|--------|--------|-------|
| Decomposition engine | **Stays** | Works today |
| Verification (TDD, types, semgrep) | **Stays** | Works today |
| Multi-agent orchestration | **Stays** | Works today |
| Skills (idea, design, build, review) | **Stays** | Works today |
| Ship / Drive modes | **Stays** | Works today |
| Go engine (dispatcher, synthesis) | **Stays** | Works today |
| **Acceptance test gate** | **NEW** | "Does the app actually work?" — smoke test after every build |
| **Evidence log** | **NEW** | Single `.sdp/log/`, four event types, hash-chained |
| **Model provenance** | **NEW** | `generation` event type inside the evidence log |
| **`sdp log trace`** | **NEW** | Forensic trace from commit to spec |
| **Protocol schema** | **NEW** | JSON Schema for the evidence format |
| **Observability bridge** | **DESIGN** | Deploy markers, OTel span attributes |

**The real delta: a protocol that checks results, not just process. Plus a data layer that records what happened.**

---

*SDP Vision v4.0 — February 2026*
*Not a rebrand. A data layer.*
