# SDP Feature Map

> What exists. What's needed. Priorities.

---

## Already Built

Not on the roadmap. Done.

### Decomposition

| Feature | Status |
|---------|--------|
| NL → atomic workstreams | Done |
| Dependency graph (Kahn's algorithm) | Done |
| Parallel dispatch (goroutines, circuit breaker) | Done |
| Checkpoint recovery (atomic write-fsync-rename) | Done |
| Synthesis engine (multi-agent conflict resolution) | Done |

### Verification

| Layer | What It Catches | Status |
|-------|----------------|--------|
| Type checking | Type errors | Done |
| Static analysis (semgrep) | Security patterns | Done |
| Test execution (TDD) | Functional correctness | Done |
| Coverage gates (≥80%) | Untested paths | Done |
| Contract validation | OpenAPI mismatches | Done |

### Orchestration

| Feature | Status |
|---------|--------|
| 19 agent roles | Done |
| Ship mode (autonomous) | Done |
| Drive mode (human-in-the-loop) | Done |
| Progressive disclosure | Done |
| Adversarial review (6 agent roles) | Done |

### Infrastructure

| Feature | Status |
|---------|--------|
| Beads (issue tracking) | Done |
| Guard enforcement | Done |
| Quality gates (<200 LOC, CC<10) | Done |
| Telemetry collector (JSONL) | Done |

---

## P0 — Result Checking + Evidence Log

The protocol must answer "does it work?" before recording "what happened."

### Acceptance Test Gate

| Feature | Description |
|---------|-------------|
| **E2E smoke test** | After every `@build`: start the app, hit core endpoint, verify response. Failure = build failed, regardless of coverage |
| **Smoke test config** | `.sdp.yml` defines what "works" means per project (HTTP 200? Container starts? Output matches?) |
| **Graceful fallback** | Warn if no smoke test defined, don't block |

This is the vibe-coder's feedback loop (run it after every change), formalized into the protocol. The user controls architecture depth and scope by talking to `@idea` — that's a conversation, not a feature. The acceptance test is the one gate nobody can opt out of.

### Schema

| Feature | Description |
|---------|-------------|
| Schema consolidation | Fix existing sprawl (`schema/`, `docs/schema/`, frontmatter) |
| Single namespace | One `$id`, one validation entrypoint |
| Schema versioning | SemVer, changelog, breaking change policy |
| Evidence schema v0.1 | JSON Schema for four event types |

### Evidence Log

| Feature | Description |
|---------|-------------|
| `plan` event | Intent: feature, units, dependencies, cost estimate |
| `generation` event | Provenance: model, version, prompt hash, params, timestamp, spec ref, initiator, code hash |
| `verification` event | Tool, command, **actual output** (pytest stdout, mypy output), pass/fail, coverage |
| `approval` event | Who approved, when, what they saw, reasoning |
| Hash chain | `prev_hash` per record — corruption detection, not tamper-proof |
| Committed by default | `.sdp/log/` in git, not gitignored |
| Append-only merge | `.gitattributes` merge driver |
| Storage budget | ~1KB generation, ~5KB verification |

### Forensic Trace

| Feature | Description |
|---------|-------------|
| `sdp log trace <commit>` | Chain backwards: commit → model → spec → verification → approver |
| Tree view | Default human-readable output |
| `--output=json` | Machine-readable |
| Offline | All data local |
| Missing evidence | Graceful: "no evidence for this commit" |
| `sdp log show` | Browse log with filters |

### `@build` Instrumentation

| Feature | Description |
|---------|-------------|
| Generation events | Every AI generation → `generation` event |
| Verification events | Actual tool output, not summary |
| Approval events | Auto (ship) or human (drive) |
| Plan events | On decomposition completion |
| Automatic | No opt-in |

### Scope Collision Detection

| Feature | Description |
|---------|-------------|
| Cross-reference scope | Compare `scope_files` across all in-progress workstreams |
| Collision signal | "Feature A WS-3 and Feature B WS-7 both modify `user_model.go`" |
| Signal, not block | Warn and suggest coordination, don't prevent work |
| Query existing data | Workstream specs already declare scope — this is just a query |

### Compliance Design Doc

| Feature | Description |
|---------|-------------|
| Data residency | Where the log lives |
| Retention | GDPR vs audit trail |
| RBAC | Who sees evidence |
| Integrity | What hash chain guarantees and what it doesn't |
| Prompt privacy | Hash only, raw prompts never stored |

---

## P1 — Full Instrumentation + CLI + CI/CD

### Skills Instrumentation

| Feature | Description |
|---------|-------------|
| `@review` → verification events | Findings, actual review output |
| `@deploy` → approval chain | Merge approval, gates passed |
| `@design` → plan events | Decomposition decisions |
| `@idea` → decision events | Questions + answers (drive mode) |
| All 19 skills | Full pipeline coverage |

### CLI

| Feature | Description |
|---------|-------------|
| `sdp plan <description>` | Decomposition + plan event |
| `sdp apply` | Execution + full evidence chain |
| `sdp log trace <commit>` | Forensic trace (standalone) |
| `sdp log show` | Evidence browser |
| `--auto-apply` | Ship mode |
| `--interactive` | Drive mode |
| `--retry N` | Retry failed unit |
| `--output=json` | JSON for all commands |
| Streaming progress | Per-unit progress |

### CI/CD

| Feature | Description |
|---------|-------------|
| GitHub Action | `sdp-dev/verify-action@v1` |
| PR evidence comment | Evidence chain summary |
| Provenance gate | Block merge without evidence |
| GitLab CI | Same for GitLab |

### Observability Bridge (Design)

| Feature | Description |
|---------|-------------|
| Deploy markers | Evidence → deploy events |
| OTel span attributes | AI-generated code paths in traces |
| Diff-level provenance | Which lines are AI vs human |
| Integration spec | Honeycomb / Datadog / Grafana |

### Shared Contracts for Parallel Features

| Feature | Description |
|---------|-------------|
| Boundary detection | When `@design` runs for parallel features, identify shared surfaces |
| Interface contracts | Generate API/data model contracts before parallel implementation starts |
| Contract synthesis | Extend synthesis engine to cross-feature boundaries |
| Contract-first build | Parallel workstreams build against shared contracts, not assumptions |

### Data

| Feature | Description |
|---------|-------------|
| MTTR tracking | With evidence vs without |
| Verification telemetry | Catch rate, iterations |
| AI failure taxonomy | By model/language/domain |
| Quarterly benchmark | Public publication |

---

## P2 — SDK + High-Assurance + Observability

### SDK

| Feature | Description |
|---------|-------------|
| `sdp.Verify()` | Verification engine as Go library |
| `sdp.Evidence()` | Evidence bundle generation |
| Provider adapters | Claude, GPT, Gemini |
| JSON-in/JSON-out | External tool integration |

SDK = verification + evidence. Decomposition stays in CLI/plugin.

### Observability (Implementation)

| Feature | Description |
|---------|-------------|
| OTel exporter | SDP evidence → OTel spans |
| Deploy correlation | Auto-link evidence ↔ deploys |
| Runtime context | Feature flags, blast radius, rollback |
| Per-line attribution | AI vs human per line |

### Continuous Cross-Branch Integration

| Feature | Description |
|---------|-------------|
| Per-workstream integration | After each WS: merge main, run acceptance test |
| Cross-feature matrix | "Feature A still works after Feature B's latest WS" |
| Cost-aware | Full matrix for high-risk, smoke-only for low-risk |

### Cross-Model Review

| Feature | Description |
|---------|-------------|
| Model A → Model B | Decorrelated errors |
| Risk-triggered | Auto for auth, payments, data deletion |
| Model selection policy | Route by risk |

### High-Assurance Governance

| Feature | Description |
|---------|-------------|
| Compliance export | SOC2/HIPAA/DORA |
| Verification certificates | Signed, timestamped |
| Risk-proportional | `auth/` → full trail, `components/` → light |
| Team policies | "All AI PRs need evidence" |
| Decomposition templates | Per-company patterns |
| Policy profiles | Stricter verification presets for critical paths |

### IDE

| Feature | Description |
|---------|-------------|
| Cursor plugin | Plan/apply/log from IDE |
| VS Code extension | Same |
| JetBrains | Same |

---

## P3 — Standard

| Feature | Description |
|---------|-------------|
| SDP Evidence Format v1.0 | Published spec |
| External adoption | 2+ tools |
| Signed evidence | Non-repudiation (compliance-grade) |
| External timestamping | Third-party authority |
| On-premise | Air-gapped environments |
| Industry working group | After 50+ deployments |

---

## North Star: Real-Time Multi-Agent Coordination

Not on the roadmap. The guiding direction.

| Feature | Description |
|---------|-------------|
| Live intent broadcasting | Every participant knows what others are trying to achieve |
| Real-time conflict detection | Interface change → immediate signal to affected agents |
| Automatic interface negotiation | Agents resolve contract conflicts without human mediation |
| Heterogeneous coordination | Works across models, tools, and humans with different workflows |

**How we get there:** P0 scope collision → P1 shared contracts → P2 cross-branch integration → beyond.

---

*SDP Features v7.0 — February 2026*
