# SDP Roadmap

> What we build and in what order.

---

## Current State (February 2026)

### Works

- **Claude Code plugin**: 19 skills/agents (idea, design, build, review, deploy, vision, reality)
- **Go engine**: decomposition, dependency graph, parallel dispatcher, synthesis engine, circuit breaker
- **Verification**: TDD (red → green → refactor), coverage ≥80%, contract validation, types, semgrep
- **Telemetry**: append-only JSONL collector, atomic write-fsync-rename checkpoints
- **Enterprise**: 2 contracts (bank, airline), 1 evaluating (marketplace)

### Broken

- **Result checking** — gates verify coverage/types/lint, but not "does the app actually work"
- **Parallel coordination** — no cross-feature scope collision detection, no shared contracts
- No structured evidence log — results stay in markdown reports
- No model provenance — unknown which model generated what
- No forensic trace — no way to reconstruct the chain from commit to spec
- No protocol schema — "protocol" is informal, lives in skill code

---

## P0 — Acceptance Test + Evidence Log

### Acceptance Test Gate

**The one thing that's actually broken in SDP.**

Real case: SDP built an ad server — 7,149 LOC, 88% coverage, Clean Architecture, all quality gates green. Basic features didn't work. Weighted random always returned max weight. Budget enforcement was never called. Frequency capping was missing entirely.

Meanwhile, vibe-coded alternative: 1,005 LOC, works from `docker compose up`.

The problem: SDP has a 5-minute feedback loop for code quality (types, tests, coverage) and a **days-long** feedback loop for "does the app work." Vibe-coders run the app after every change. SDP builds 25 workstreams, then discovers nothing works.

**Fix:** After every workstream, check if the app still works.

- E2E smoke test after `@build`: start the app, hit the core endpoint, verify the response
- Not unit tests — an actual "run it and check" verification
- Failure = build failed, regardless of coverage
- Configurable per project: `.sdp.yml` defines what "works" means (HTTP 200? Container starts? Output matches?)
- Graceful when no smoke test defined: warn, don't block
- This is the vibe-coder's feedback loop, formalized

Everything else — adaptive architecture, scope limiting, value-first decomposition — those are nice-to-have. The user can always say "I don't care about architecture." The user can't say "I don't care if it works."

### Schema Consolidation

Existing schemas diverged across `schema/`, `docs/schema/`, frontmatter — and don't match. Fix before adding new ones.

- Audit all `.schema.json`, frontmatter formats, validation paths
- Single `$id` namespace, one validation entrypoint
- Migrate stale schemas (`WS-` → `PP-FFF-SS`)
- Versioning: SemVer, changelog

### Evidence Schema v0.1

JSON Schema for the evidence log. Criterion: another tool can produce SDP-compatible evidence from the schema alone.

- `plan` — intent: feature description, units, dependencies, cost estimate
- `generation` — model: name, version, prompt hash, temperature, params, timestamp, spec reference, initiator, code hash
- `verification` — check: tool, command, **actual output** (real pytest stdout, not "tests passed"), pass/fail, coverage, duration
- `approval` — decision: who approved, when, what they saw, why (in drive mode)
- Hash chain: each record contains `prev_hash` (corruption detection, **not** tamper-proof)
- Version 0.1, breaking changes expected

### Evidence Log

Single `.sdp/log/events.jsonl`.

- Append-only JSONL (extends existing telemetry primitive)
- One record per event, following published schema
- Provenance is not a separate subsystem — it's the `generation` event type
- One log, many event types, one reader
- Committed by default (not gitignored — evidence must survive laptop wipes and ephemeral CI runners)
- `.gitattributes` with append-only merge driver
- Budget: ~1KB per generation, ~5KB per verification (with actual output)

### `sdp log trace <commit>`

New command: walk the evidence chain backwards from a commit.

- Output: tree — model → spec → verification output → approver
- `--output=json` for tooling
- Works offline (all data is local)
- Graceful on missing evidence ("no SDP evidence for this commit")

### `@build` Instrumentation

First skill. Not all 19 at once.

- `@build` emits `generation` event on every AI generation
- `@build` emits `verification` event with **actual tool output**
- `@build` emits `approval` event (auto-approve in ship, human-approve in drive)
- `plan` event when decomposition completes
- Emission is automatic — no opt-in

### Scope Collision Detection

When 10 agents and 5 humans work on 5 features in parallel — who's touching what?

- Each workstream already declares `scope_files` in its spec
- Cross-reference scope declarations across all in-progress workstreams
- If two parallel workstreams (across different features) modify the same files — signal, not block
- "Feature A WS-3 and Feature B WS-7 both modify `user_model.go`. Coordinate."
- Cheap to build: it's a query over existing data

### Compliance Design Doc

Not code — a document for enterprise conversations.

- Data residency, retention, RBAC, integrity guarantees
- Prompt privacy (hash only, never raw prompts)
- What the hash chain guarantees and what it doesn't

---

## P1 — Full Instrumentation + CLI + CI/CD

**Remaining skills instrumentation.**
- `@review` → verification events (findings, actual output)
- `@deploy` → approval chain (who approved merge, what gates passed)
- `@design` → plan events (decomposition decisions, dependency graph)
- `@idea` → decision events (questions and answers in drive mode)
- All 19 skills produce evidence

**CLI: plan / apply / log.**

```bash
sdp plan "Add OAuth2 login"            # Show decomposition, emit plan event
sdp apply                               # Execute with evidence recording
sdp plan "Add auth" --auto-apply        # Ship mode
sdp plan "Add auth" --interactive       # Drive mode
sdp apply --retry 3                     # Retry failed unit only
sdp log trace abc1234                   # Forensic trace
sdp log show                            # Browse evidence
sdp log show --unit 00-001-03           # Evidence for specific unit
```

- `sdp plan` → `@idea` + `@design` + plan event
- `sdp apply` → `@build` + `@review` + full evidence chain
- `sdp log trace` → reads `.sdp/log/` + git metadata
- `sdp log show` → log browser with filters
- `--output=json` for all commands
- Streaming progress per unit

**GitHub Action.**

```yaml
- uses: sdp-dev/verify-action@v1
  with:
    gates: [types, semgrep, tests, evidence]
```

- Verification on every AI-generated PR
- Evidence chain summary in PR comment
- Provenance gate: block merge if AI code lacks evidence
- Free tier: 50 runs/month

**Observability bridge (design).**
- Deploy markers: tie evidence records to deploy events
- OTel span attributes: mark AI-generated code paths in traces
- Diff-level provenance: which *lines* are AI-generated vs human-edited
- Integration spec: how SDP connects to Honeycomb / Datadog / Grafana

**Shared contracts for parallel features.**
- When `@design` runs for multiple features in parallel, identify shared boundaries
- Generate interface contracts (API surfaces, data models, function signatures) before implementation starts
- Synthesis engine already exists — extend it to cross-feature boundary detection
- "Auth and Payments both need User — here's the contract both must respect"

**Data collection.**
- Instrument every run: catch rate, iteration count, model performance
- AI failure taxonomy: what goes wrong, by model/language/domain
- MTTR tracking: time to root cause with evidence vs without
- Target: "AI Code Quality Benchmark" — quarterly publication

---

## P2 — SDK + Enterprise + Observability

**Verification + Evidence SDK.**

```go
result := sdp.Verify(files, sdp.Gates{"types", "tests", "semgrep"})
evidence := sdp.Evidence(result)  // Schema-compliant evidence bundle
```

- `sdp.Verify()` — verification engine as library
- `sdp.Evidence()` — evidence bundle generation
- Provider adapters: Claude, GPT, Gemini
- JSON-in/JSON-out API
- SDK = verification + evidence only. Decomposition stays in CLI/plugin.

**Observability integration (implementation).**
- OTel exporter: SDP evidence → OpenTelemetry spans
- Deploy correlation: auto-link evidence to deploy events
- Runtime context: feature flags, blast radius, rollback path
- Per-line attribution: AI vs human per line

**Cross-model review.**
- Model A generates, Model B reviews (decorrelated errors)
- Auto-trigger for high-risk code (auth, payments, data deletion)
- Model selection policy

**Continuous cross-branch integration.**
- After each workstream: merge main into feature branch, run acceptance test
- Catches integration breaks per-workstream, not per-PR (when it's too late)
- Cross-feature test matrix: "Feature A still works after Feature B's latest workstream"
- Cost-aware: full matrix for high-risk features, smoke-only for low-risk

**Enterprise features.**
- Compliance export: SOC2/HIPAA/DORA-ready evidence format
- Verification certificates (signed, timestamped)
- Risk-proportional verification: `auth/` → full trail, `components/` → light
- Team policies: "all AI PRs need evidence chain"
- Shared decomposition templates
- Billing/metering

**IDE.**
- Cursor plugin
- VS Code extension
- JetBrains plugin

---

## P3 — Standard

- SDP Evidence Format v1.0 — published, auditor-reviewed
- Adoption: 2+ tools beyond SDP
- Signed evidence records (cryptographic non-repudiation — compliance-grade)
- External timestamping (third-party timestamp authority)
- On-premise for air-gapped environments
- Industry working group (only after 50+ deployments)

---

## Success Metrics

| Metric | P0 | P1 | P2 | P3 |
|--------|----|----|----|----|
| Acceptance test pass rate | Baseline | 90% | 95% | Published benchmark |
| Evidence records | 100 (dogfood) | 5,000 | 50,000 | 500,000 |
| Skills instrumented | 1 (`@build`) | All 19 | All + SDK | All + external |
| MTTR improvement | Baseline | 20% ↓ | 40% ↓ | Benchmark published |
| Repos with evidence | 3 | 200 | 1,000 | 5,000+ |

---

## Kill Criteria

> **Primary:** After 100 incidents involving AI-generated code: if SDP evidence did not reduce MTTR compared to git-blame-only — rethink the approach.

> **Secondary:** After 500 SDP runs: if catch rate < 5% AND defect rate equals baseline — kill the product entirely.

> **New:** After 50 SDP builds: if acceptance test pass rate is not measurably higher than baseline vibe-coding — the orchestration isn't adding value.

---

## Risks

| Risk | Mitigation |
|------|-----------|
| SDP produces non-working code with passing gates | Acceptance test gate (P0) |
| Over-architecture for small projects | Adaptive architecture (P0) |
| Scope creep through deep questioning | Scope limiter + ship gate (P0) |
| Evidence log adds friction | Automatic emission, no opt-in |
| Models improve, decomposition unnecessary | Provenance valuable regardless. Decomposition is a technique, not the core |
| Hash chain oversold to enterprise | Honest labeling: corruption detection, not non-repudiation |
| Evidence never consulted during incidents | Observability bridge: deploy markers, OTel |
| Schemas diverge again | Single namespace, CI validation |
| Cursor/GitHub ships native verification | Protocol + evidence log is the moat |
| Log bloats repo | Budget per event, compression, archival policy |

---

## North Star: Real-Time Multi-Agent Coordination

Not on the roadmap. Not buildable today. But this is where the world is going, and we're thinking about it.

**The problem:** 10 agents and 5 humans build 5 features in parallel. Agent A changes the User interface. Agent B is mid-build using the old User interface. Agent B doesn't know. Nobody knows until merge time, which is the worst possible moment to find out.

**The vision:** Real-time coordination where every participant — agent or human — has live awareness of what others are doing. Not just scope (which files) but intent (what they're trying to achieve) and state (what they've changed so far). When Agent A changes an interface, Agent B gets a signal immediately, not at merge time.

This is a distributed systems problem. Eventual consistency, conflict resolution, intent broadcasting across heterogeneous participants (different models, different tools, humans with different workflows). It's hard. Nobody has solved it.

**How we get there incrementally:**
- P0: Scope collision detection (static, at plan time)
- P1: Shared contracts (at design time, before parallel work starts)
- P2: Continuous cross-branch integration (at build time, after each workstream)
- Beyond: Live intent broadcasting, real-time conflict detection, automatic interface negotiation

Each step is useful on its own. Together, they build toward a system where parallel AI-assisted development doesn't collapse into merge hell.

---

*SDP Roadmap v7.0 — February 2026*
