# F054: SDP Evidence Layer — Workstream Summary

> Feature: sdp-ejzz | 11 workstreams | ~19h total

## Dependency Graph

```
           ┌─ 01 Schema ─── 03 Evidence Schema ─── 04 Evidence Writer ─┬─ 05 @build Instrumentation ─┐
           │                  (6 event types)                           ├─ 07 sdp log trace ──────────┤
Parallel → ├─ 02 Config ─── 06 Acceptance Test Gate ───────────────────┤                              ├─► 09 Dogfooding
           │                                                            ├─ 10 Decision Capture ────────┤
           └─ 08 Scope Collision ───────────────────────────────────────┘    └─ 11 Lessons & Warnings ─┘
```

## Execution Order (3 parallel lanes)

| Wave | Lane A (Evidence) | Lane B (Acceptance) | Lane C (Collision) |
|------|-------------------|--------------------|--------------------|
| 1    | 01 Schema Consolidation | 02 Project Config | 08 Scope Collision |
| 2    | 03 Evidence Schema (6 types) | 06 Acceptance Test Gate | — |
| 3    | 04 Evidence Writer | — | — |
| 4    | 05 @build Instrumentation | — | — |
| 4    | 07 sdp log trace | — | — |
| 4    | 10 Decision Capture & Search | — | — |
| 5    | 11 Lessons & "We Already Tried" | — | — |
| 6    | **09 Dogfooding & Integration** (all converge) | | |

**Critical path:** 01 → 03 → 04 → 10 → 11 → 09 (6 workstreams, ~11h)

## Workstream Registry

| WS ID | Title | Est. | Depends On | Beads ID |
|-------|-------|------|------------|----------|
| 00-054-01 | Schema Consolidation | 1h | — | sdp-kaox |
| 00-054-02 | Project Config File | 1.5h | — | sdp-d03u |
| 00-054-03 | Evidence Schema v0.1 (6 event types) | 2h | 01 | sdp-jgd0 |
| 00-054-04 | Evidence Log Writer | 2h | 03 | sdp-uoud |
| 00-054-05 | @build Instrumentation | 2h | 04 | sdp-7ybe |
| 00-054-06 | Acceptance Test Gate | 1.5h | 02 | sdp-6lfu |
| 00-054-07 | sdp log trace | 1.5h | 04 | sdp-2tnp |
| 00-054-08 | Scope Collision Detection | 1.5h | — | sdp-3c2f |
| 00-054-09 | Dogfooding & Integration | 2h | 05,06,07,08,11 | sdp-63lk |
| **00-054-10** | **Decision Capture & Search** | **2h** | 04 | TBD |
| **00-054-11** | **Lessons & "We Already Tried" Warnings** | **2h** | 10 | TBD |

## Key Files (New)

```
schema/
  evidence.schema.json          # WS-03 (6 event types: plan, generation, verification, approval, decision, lesson)
  config.schema.json            # WS-02
sdp-plugin/internal/
  evidence/
    types.go                    # WS-03
    writer.go                   # WS-04
    reader.go                   # WS-04
    emitter.go                  # WS-05, WS-10
    tracer.go                   # WS-07
    lessons.go                  # WS-11
    warnings.go                 # WS-11
  config/
    project.go                  # WS-02
  acceptance/
    runner.go                   # WS-06
  collision/
    detector.go                 # WS-08
sdp-plugin/cmd/sdp/
  log.go                        # WS-07
  acceptance.go                 # WS-06
  collision.go                  # WS-08
  decisions.go                  # WS-10 (updated)
.sdp/
  config.yml                    # WS-09
  log/events.jsonl              # WS-09
```

## Key Design: Memory IS the Evidence Log

Old F051 had separate `.sdp-memory/` with 14 workstreams (9 weeks).
New approach: decisions and lessons are event types in the evidence log.

| Old F051 | New F054 |
|----------|----------|
| `.sdp-memory/decisions.jsonl` | `decision` events in `.sdp/log/events.jsonl` |
| `.sdp-memory/lessons.jsonl` | `lesson` events in `.sdp/log/events.jsonl` |
| `.sdp-memory/sessions.jsonl` | Telemetry + evidence events (already tracked) |
| `sdp memory search` | `sdp log show --type=decision --search` |
| `sdp memory lessons` | `sdp log show --type=lesson` |
| Separate similarity engine | Keyword + tag matching on existing evidence |
| 14 WS, 9 weeks | 2 WS, ~4h |

## Execution

```bash
# Option A: Manual (recommended for first time)
@build 00-054-01   # Start with schema
@build 00-054-02   # Parallel with 01
@build 00-054-08   # Parallel with 01 and 02
# ... continue by dependency order

# Option B: Autonomous
@oneshot F054
```
