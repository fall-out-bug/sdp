# F054: SDP Evidence Layer

> Beads: sdp-ejzz | Priority: P0

---

## Problem

SDP verifies code quality but not product quality. Real case: built a 7K LOC ad server with 88% coverage, Clean Architecture, all gates green — basic features broken. Vibe-coded 1K LOC alternative worked fine.

Additionally: no structured record of what AI agents did. When something breaks, you have git blame and nothing else.

And: no coordination when multiple agents/humans work in parallel.

## Solution

Three capabilities, shipped in order:

### 1. Acceptance Test Gate
After every `@build`, run the app and check if it does what it's supposed to. The vibe-coder's feedback loop, formalized.

- `.sdp.yml` defines smoke test per project
- Start app, hit endpoint, verify response
- Failure = build failed, regardless of coverage
- 30 seconds max, 5 lines of config

### 2. Thin Evidence Log
Record what happened — model ID, pass/fail, timestamp. Not full forensic chain yet.

- Single `.sdp/log/events.jsonl`, append-only
- Schema v0.1 (four event types: plan, generation, verification, approval)
- Provenance = `generation` event type (not separate subsystem)
- Hash chain for corruption detection (not tamper-proof — honest)
- Committed to repo by default
- `sdp log trace <commit>` to walk the chain backwards

### 3. Scope Collision Detection
When parallel workstreams touch the same files — signal it.

- Cross-reference `scope_files` across in-progress workstreams
- Signal, not block
- Query over existing workstream spec data

## Users

- **Developers** using `@build` — acceptance test catches broken builds
- **Enterprise customers** (bank, airline) — evidence log for audit compliance
- **Ops/incident responders** — `sdp log trace` for forensic reconstruction
- **Parallel teams** — scope collision prevents merge disasters

## Success Metrics

- Acceptance test catches failures that unit tests miss
- Evidence log records 100+ events from dogfooding
- `sdp log trace` reconstructs chain for any SDP-generated commit
- Scope collision signals at least 1 real overlap in parallel work

## Kill Criteria

- After 50 builds: if acceptance test pass rate not higher than baseline vibe-coding → not adding value
- After 100 incidents: if evidence didn't reduce MTTR → rethink evidence layer

## Constraints

- Schema consolidation BEFORE new schemas (existing sprawl across 3 locations)
- Instrument `@build` first, not all 19 skills
- Evidence must be thin first (model ID + pass/fail + timestamp), full forensics later
- Acceptance test must be cheap (30 seconds, not 10 minutes)

### 4. Decision Memory (from F051 reconception)
Record decisions as evidence events. "We decided X because Y" — searchable, reversible.

- `decision` event type in evidence schema
- Auto-captured during `@design` and `@build`
- Searchable: `sdp log show --type=decision --search "auth"`
- Reversal tracking: new decision can reference the old one it overturns

### 5. Lessons & "We Already Tried" Warnings
Auto-extract lessons from outcomes. Warn when approaching past failures.

- `lesson` event type extracted on workstream completion
- Proactive: `sdp guard activate` checks for similar failed decisions
- Keyword + tag matching (no ML — simple first)
- Non-blocking: warnings, not gates

## Non-Goals

- Full forensic chain with complete verification output (P1)
- CLI standalone commands `sdp plan` / `sdp apply` (P1)
- GitHub Action (P1)
- Observability bridge / OTel (P1-P2)
- Cross-model review (P2)
- Compliance export (P2)

## Technical Context

### Existing infrastructure to build on:
- `internal/telemetry/collector.go` — append-only JSONL with `os.O_APPEND|os.O_CREATE|os.O_WRONLY`
- Checkpoint system with atomic write-fsync-rename
- Workstream specs with `scope_files` declarations
- Schema at `schema/workstream.schema.json` (needs consolidation)

### Key architecture decisions:
- One log, many event types, one reader
- Provenance is an event type, not a subsystem
- Hash chain = corruption detection, not non-repudiation
- Evidence committed to repo (not gitignored)
- `.gitattributes` append-only merge driver for the log

---

*Created from: docs/vision/VISION.md v4, ROADMAP.md v7, FEATURES.md v7*
*Expert input: Cagan, Karpathy, Hashimoto, Majors (Panel V4 + V5)*
