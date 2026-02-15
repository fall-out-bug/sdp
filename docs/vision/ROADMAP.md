# SDP Unified Roadmap

> Wearing a tie, but t-shirts on Fridays.
> Discovery + Delivery + Evidence. From toy to sports car.

---

## Product Map

```
                    DISCOVERY                          DELIVERY
             ┌────────────────────┐           ┌────────────────────┐
  @vision ───┤  strategy          │           │   @build ──── TDD  │
  (arch?) ───┤  architecture gap  │ contracts │   @review ── self- │
  @reality ──┤  codebase truth    │◄─────────►│       │    check   │
  @feature ──┤  @idea + @design   │           │   @deploy ── prod  │
             │  (orchestrator)    │           │   @oneshot ── orch │
             └────────┬───────────┘           └──────────┬─────────┘
                      │         EVIDENCE                 │
                      └───►  sdp trace  ◄────────────────┘
             ┌──────────────────────────────────────────────┐
             │  multi-agent synthesis │ beads │ Go engine    │
             │  dependency graph │ parallel dispatch         │
             │  circuit breaker │ checkpoint recovery        │
             └──────────────────────────────────────────────┘
```

**Main flow:** `@feature → @oneshot → sdp trace` (plan it → ship it → prove it)

---

## Strategic Vector (Discovery + Delivery)

SDP's next evolution is not a pivot and not a vendor lock-in.

- `Discovery`: owner gives vision + constraints, agent team drives structured discovery to evidence-backed decisions.
- `Delivery`: approved decisions become executable delivery contracts, then autonomous build/review/deploy with proof.
- `Evidence`: one continuous chain from discovery assumptions to delivery verification.
- `Runtime-agnostic`: SDP stays a control protocol over agent runtimes, not tied to a single provider/tool.

### Operating Modes

| Mode | Goal | Guardrails | Use When |
|------|------|------------|----------|
| `Explore Path` | Max speed and learning | Minimal gates, capped cost, explicit uncertainty | Early discovery, low-stakes experiments |
| `Commit Path` | Max confidence and reliability | Full evidence + quality + approval gates | Delivery to prod, external launch, money/safety risk |

### Month Focus (Near-Term)

1. Formalize discovery handoff contracts and budget gates.
2. Wire discovery outputs directly into delivery contracts.
3. Keep commit-path quality high while reducing overhead via explore-path routing.

---

## Adoption Model (L0-L2, MIT)

SDP ships as a progressive stack. Teams adopt in slices.

| Level | Scope | Install | License |
|-------|-------|---------|---------|
| `L0` | Protocol: prompts, guides, templates, schemas | Prompt pack for any agent runtime | MIT |
| `L1` | Safety: hooks, guard, traces, provenance | CLI safety bundle (brew) | MIT |
| `L2` | Orchestrator: plan/apply/log, dispatcher, checkpoints | CLI core bundle (brew) | MIT |

Each level is useful standalone. Upgrade path is additive: `L0 → L1 → L2`.

See [LAYERED-ADOPTION.md](LAYERED-ADOPTION.md) for packaging details.
See [SPECTRUM.md](SPECTRUM.md) for the progressive disclosure model (T-shirt → Smart casual → Tie).

---

## Already Built

Not on the roadmap. Working. This is the foundation.

### Discovery

| Feature | Details | Status |
|---------|---------|--------|
| @vision | Product strategy, 7 expert agents, panel synthesis | Done |
| @reality | Codebase analysis, 8 expert agents, health/debt/gaps | Done |
| @idea | Requirements gathering, deep interview, progressive disclosure (12-27 questions) | Done |
| @design | Architecture + workstream decomposition, dependency graph | Done |
| @feature | Discovery orchestrator (@idea + @design) | Done |

### Delivery

| Feature | Details | Status |
|---------|---------|--------|
| @build | TDD cycle (red → green → refactor), single workstream | Done |
| @review | Adversarial review, 6 agent roles (QA, Security, DevOps, SRE, TechLead, Docs) | Done |
| @deploy | Production deployment, branch merge | Done |
| @oneshot | Delivery orchestrator (@build + @review + @deploy for all WS) | Done |
| Ship mode | Autonomous execution, no human stops | Done |
| Drive mode | Human-in-the-loop, decide at every fork | Done |

### Engine

| Feature | Details | Status |
|---------|---------|--------|
| NL → atomic workstreams | Natural language decomposition | Done |
| Dependency graph | Kahn's algorithm, topological sort | Done |
| Parallel dispatch | Goroutines, 4.96x speedup | Done |
| Circuit breaker | Fault tolerance for parallel execution | Done |
| Checkpoint recovery | Atomic write-fsync-rename, resume after crash | Done |
| Synthesis engine | Multi-agent conflict resolution (unanimous → expertise → quality → merge → escalate) | Done |

### Verification

| Feature | Details | Status |
|---------|---------|--------|
| TDD pipeline | Red → green → refactor enforcement | Done |
| Coverage gates | ≥80% enforced | Done |
| Type checking | mypy --strict | Done |
| Static analysis | Semgrep security patterns | Done |
| Contract validation | OpenAPI mismatches | Done |
| Quality gates | <200 LOC, CC<10, no `except: pass` | Done |

### Infrastructure

| Feature | Details | Status |
|---------|---------|--------|
| Beads | Issue tracking across sessions, hash-based IDs | Done |
| Guard enforcement | Active workstream required for edits | Done |
| Telemetry collector | Append-only JSONL, atomic checkpoints | Done |
| Progressive disclosure | Reduce question fatigue, 3-question cycles | Done |

### Delivered Features

| Feature | Workstreams | Status |
|---------|-------------|--------|
| F052: Multi-Agent System | 26 WS (00-052-00..25) | Done |
| 200+ workstreams completed | — | Done |

---

## M1: "T-shirt" — Evidence Foundation

> The protocol starts recording. User workflow unchanged.
> **Timeline:** 4-6 weeks (Mar-Apr 2026). **Level:** Solo developer.

### Features → Workstreams

| Feature | Description | Est. | WS |
|---------|-------------|------|----|
| **F054**: Evidence Layer | Evidence writer, 6 event types, `.sdp/log/events.jsonl`, hash chain | ~19h | 11 WS |
| **F063**: Guardian Hooks | Hooks, guardrails, acceptance test integration | ~8h | 4 WS |
| **F064**: Unified Task Resolver | Internal: resolve tasks consistently across skills | P0 | 4 WS |
| **F067**: Repo Hardening | Schema cleanup, alignment, tech debt from audits | P0 | 14 WS |
| **F068**: UX Foundation | First-run experience, README rewrite, time-to-first-value | P0 | 5 WS |
| **F070**: Failure & Recovery UX | Fast diagnosis, safe recovery, low frustration | P0 | 5 WS |
| **F075**: Self-Healing Doctor | `sdp doctor --repair --deep`, guided fix-it output, config migrations | P0 | 6 WS |
| **F076**: Guided Onboarding Wizard | Interactive + non-interactive onboarding with safe defaults and preflight checks | P0 | 5 WS |
| **NEW**: Acceptance Test Gate | Smoke test after @build, `.sdp.yml` config, graceful fallback | P0 | TBD |

### M1 Deliverables

**Evidence Log**

| Feature | Details | Priority |
|---------|---------|----------|
| **Evidence writer** | New `evidence/` package in Go. Append-only `.sdp/log/events.jsonl` | Must |
| **`decision` event** | Discovery: question asked, answer given, rationale | Must |
| **`plan` event** | Intent: feature, units, dependencies, cost estimate | Must |
| **`generation` event** | Provenance: model, version, prompt hash, params, timestamp, spec ref, code hash | Must |
| **`verification` event** | Tool, command, actual output (pytest stdout, not "passed"), pass/fail, coverage | Must |
| **`acceptance` event** | App-level smoke test result, command, timing | Must |
| **`approval` event** | Who approved, when, what they saw, reasoning (ship=auto, drive=human) | Must |
| **@build instrumentation** | @build auto-emits generation + verification + approval events | Must |
| **@idea instrumentation** | @idea auto-emits decision events (questions + answers) | Should |
| **@design instrumentation** | @design auto-emits plan events (decomposition, dependencies) | Should |
| **Committed by default** | `.sdp/log/` in git, not gitignored | Must |
| **Storage budget** | ~1KB per generation, ~5KB per verification (with actual output) | Must |

**Acceptance Test Gate**

| Feature | Details | Priority |
|---------|---------|----------|
| **Smoke test runner** | After every @build: start app, hit endpoint, check response, shut down | Must |
| **`.sdp.yml` config** | `acceptance.command`, `acceptance.timeout` (30s default) | Must |
| **Graceful fallback** | No `.sdp.yml` = warning, don't block | Must |
| **Evidence integration** | Result → `acceptance` event in evidence log | Must |
| **30 second max** | If acceptance test takes >30s, it won't be run. Fast or nothing. | Must |

**Forensic Trace**

| Feature | Details | Priority |
|---------|---------|----------|
| **`sdp trace <commit>`** | Walk evidence chain backwards: decisions → model → spec → tests → approval | Must |
| **Tree output** | Human-readable tree view (default) | Must |
| **`--json`** | Machine-readable output | Should |
| **Offline** | All data local, no network required | Must |
| **Graceful** | "No SDP evidence for this commit" when data missing | Must |

**README & Positioning**

| Feature | Details | Priority |
|---------|---------|----------|
| **README rewrite** | Discovery + Delivery + Evidence. Main flow. Enter at any level. Under the Hood. | Must |
| **Spectrum reference** | Link to SPECTRUM.md — from vibe coding to enterprise | Should |
| **Positioning** | "Structured AI development — every decision recorded, every generation traceable" | Must |

**Operational DX Hardening**

| Feature | Details | Priority |
|---------|---------|----------|
| **Self-healing doctor** | Add `--repair` and `--deep` modes with actionable fix-it guidance | Must |
| **Config migration in doctor** | Auto-migrate stale config keys with backup and explicit change log | Must |
| **First-run wizard** | Interactive setup: environment, hooks, telemetry preference, quality defaults | Must |
| **Headless onboarding** | Non-interactive flags for CI/automation onboarding | Must |
| **Onboarding preflight** | Validate prerequisites before writing config; fail with clear remediation | Should |

### M1 Validation

- Does the acceptance test catch bugs that unit tests miss?
- Does the evidence log slow down @build?
- Is `sdp trace` useful on a real project?
- Does the new README reduce time to "aha moment"?
- Does self-healing doctor reduce setup/support time?
- Can new users reach first successful build in < 10 minutes?

---

## M2: "Smart Casual" — Protocol & Coordination

> The protocol becomes formal. Teams coordinate.
> **Timeline:** 4-6 weeks after M1 (May-Jun 2026). **Level:** Small team (2-5).

### Features → Workstreams

| Feature | Description | Est. | WS |
|---------|-------------|------|----|
| **F060**: Shared Contracts | Boundary detection, interface contracts, contract synthesis | ~6h | 3 WS |
| **F071**: Team UX & Collaboration | Predictable handoff, multi-actor coordination | P1 | 5 WS |
| **F073**: Trust & Explainability | Decisions explainable, traceable, auditable | P0 | 5 WS |
| **F077**: Runtime Hooks Platform | Event-driven lifecycle hooks (`command:*`, `session:*`, `gateway:*`) | P1 | 5 WS |
| **F078**: Resilient Runtime Layer | Retry/fallback substrate, degraded mode, per-step retry semantics | P0 | 4 WS |
| **NEW**: Schema Consolidation | Audit schemas, single namespace, SemVer, evidence schema v0.1 | Must | TBD |
| **NEW**: Review as Protocol Self-Check | @review: drift detection, quality regression, architecture drift | Must | TBD |
| **NEW**: Architecture Awareness | @reality → @design wiring, architecture-aware decomposition | Must | TBD |
| **NEW**: Scope Collision Detection | Cross-reference scope_files, signal not block | Must | TBD |
| **NEW**: Agentic Discovery Contracts | Owner intake schema, handoff schema, budget gate protocol | Must | TBD |
| **NEW**: Discovery → Delivery Contract Bridge | Convert discovery outputs into executable delivery contracts | Must | TBD |
| **NEW**: Explore/Commit Path Policy | Route work by risk: fast learning path vs full assurance path | Must | TBD |

### M2 Deliverables

**Schema**

| Feature | Details | Priority |
|---------|---------|----------|
| **Schema consolidation** | Audit `schema/`, `docs/schema/`, frontmatter. Fix divergence. | Must |
| **Single namespace** | One `$id`, one validation entrypoint | Must |
| **Evidence schema v0.1** | JSON Schema for all 6 event types. Published in repo. | Must |
| **Schema versioning** | SemVer, changelog, breaking change policy | Must |
| **Validation on write** | Evidence writer validates against schema before appending | Should |
| **Migrate stale schemas** | `WS-` → `PP-FFF-SS` format | Should |
| **Interop criterion** | Another tool can produce SDP-compatible evidence from schema alone | Must |

**Review as Protocol Self-Check**

| Feature | Details | Priority |
|---------|---------|----------|
| **Requirements drift detection** | @review compares spec vs implementation: "spec says X, code does Y" | Must |
| **Quality regression** | "Coverage dropped from 87% to 72% since last WS" | Must |
| **Architecture drift** | Implementation violates declared architecture patterns | Should |
| **Findings → evidence** | New `review` event type with specific findings | Must |

**Architecture Awareness (Gap Fix)**

| Feature | Details | Priority |
|---------|---------|----------|
| **@reality → @design wiring** | @reality architectural findings formally feed into @design | Must |
| **Architecture-aware decomposition** | @design considers existing architecture when creating workstreams | Must |
| **NOT a new agent** | Wire existing tools, don't create @architect | Must |

**Scope Collision Detection**

| Feature | Details | Priority |
|---------|---------|----------|
| **Cross-reference scope** | Compare `scope_files` across all in-progress workstreams | Must |
| **Collision signal** | "Feature A WS-3 and Feature B WS-7 both modify `user_model.go`" | Must |
| **Signal, not block** | Warn + suggest coordination, don't prevent work | Must |
| **Query existing data** | Workstream specs already declare scope — just a query | Must |

**Contracts v0.1**

| Feature | Details | Priority |
|---------|---------|----------|
| **Boundary detection** | When @design runs for parallel features, identify shared surfaces | Must |
| **Interface contracts** | Generate API/data model contracts before parallel implementation | Must |
| **Contract synthesis** | Extend synthesis engine to cross-feature boundaries | Should |
| **One runtime test per contract** | Real HTTP request, not static analysis | Must |

**Runtime Hardening**

| Feature | Details | Priority |
|---------|---------|----------|
| **Hook runtime events** | Extend hooks beyond git hooks to runtime lifecycle events | Must |
| **Hook packaging model** | Discover/install hook packs with metadata and eligibility checks | Should |
| **Retry policy core** | Standard retry policy for transient failures with bounded backoff | Must |
| **Fallback execution path** | Retry with fallback execution options when primary path fails | Must |
| **Degraded mode behavior** | Continue in reduced capability mode instead of full stop where safe | Should |

**Agentic Discovery Contracts**

| Feature | Details | Priority |
|---------|---------|----------|
| **Owner intake contract** | Canonical inputs: vision, target outcome, constraints, budget cap, no-go rules | Must |
| **Handoff schema enforcement** | Required fields: objective, assumptions, evidence, confidence, open questions, next action | Must |
| **Budget gate protocol** | No paid/public action without explicit owner grant and expiry | Must |
| **Evidence quality levels** | A/B/C evidence quality classification wired into decision rules | Must |

**Discovery → Delivery Contract Bridge**

| Feature | Details | Priority |
|---------|---------|----------|
| **Contract artifact set** | Discovery outputs normalized into delivery-ready contracts and acceptance criteria | Must |
| **Dependency extraction** | Discovery dependencies become delivery graph inputs automatically | Must |
| **Contradiction blocking** | Delivery start blocked when discovery assumptions and evidence conflict | Should |
| **Trace continuity** | Every delivery workstream links back to discovery decision IDs | Must |

**Explore/Commit Path Policy**

| Feature | Details | Priority |
|---------|---------|----------|
| **Risk routing policy** | Route tasks to `Explore Path` or `Commit Path` using explicit risk class | Must |
| **Explore defaults** | Fast iteration with capped spend and lighter gates | Must |
| **Commit defaults** | Full verification + approval + evidence required | Must |
| **Path override audit** | Manual overrides allowed but always logged with rationale | Should |

### M2 Validation

- Is the schema stable enough for external tool interop?
- Does scope collision prevent real conflicts or just noise?
- Does @reality → @design wiring improve decomposition quality?
- Do contracts + runtime test catch real bugs?
- Do runtime hooks stay observable and predictable under load?
- Does retry/fallback reduce flaky failures without masking real defects?
- Does Explore/Commit routing reduce delivery overhead without lowering commit-path quality?
- Does discovery-to-delivery traceability make incident/root-cause analysis faster?

---

## M3: "Blazer" — Open Protocol

> The protocol goes beyond any single tool/runtime. First external users.
> **Timeline:** 4-6 weeks after M2 (Jul-Aug 2026). **Level:** Multi-tool team.
> **Timing:** EU AI Act full enforcement — August 2026.

### Features → Workstreams

| Feature | Description | Est. | WS |
|---------|-------------|------|----|
| **F057**: CLI plan/apply/log | `sdp plan`, `sdp apply`, `sdp log trace/show` | ~8h | 4 WS |
| **F058**: CI/CD GitHub Action | `sdp-dev/verify-action@v1`, PR comments, provenance gate | ~5h | 3 WS |
| **F069**: Next-Step Engine | `sdp next`, guided flow after commands and errors | P0 | 5 WS |
| **F072**: Interop & Migration | Import/export, migration from adjacent tools | P1 | 5 WS |
| **F074**: Layered OSS Packaging | `sdp init`, L0-L2 profiles, progressive adoption | P0 | 6 WS |
| **F079**: Model Routing & Economics | Profile rotation, cooldowns, model fallback, cost-aware routing policy | P0 | 5 WS |
| **NEW**: MCP Server | `sdp-mcp` — SDP exposed via Model Context Protocol | Must | TBD |
| **NEW**: Runtime-Agnostic Agent Adapter | Stable adapter interface for different agent runtimes/providers | Must | TBD |

### M3 Deliverables

**Tool Agnosticism**

| Feature | Details | Priority |
|---------|---------|----------|
| **MCP server** | `sdp-mcp` — exposes SDP via Model Context Protocol | Must |
| **MCP tools** | `sdp_plan`, `sdp_apply`, `sdp_trace`, `sdp_status`, `sdp_next` | Must |
| **Cursor support** | Via MCP — works out of the box | Must |
| **VS Code support** | Via MCP | Must |
| **Windsurf support** | Via MCP | Should |
| **Evidence across tools** | Same format regardless of source. `generation` event records tool origin. | Must |
| **Agent adapter interface** | Runtime/provider adapters behind one SDP execution contract | Must |

**Project Init**

| Feature | Details | Priority |
|---------|---------|----------|
| **`sdp init`** | Scaffold `.sdp/`, config, skills for target tool | Must |
| **`--ai claude`** | Generate Claude Code skills | Must |
| **`--ai cursor`** | Generate Cursor skills/rules | Must |
| **`--ai vscode`** | Generate VS Code configuration | Should |
| **Progressive defaults** | Minimal config by default, expand as needed | Must |

**CLI Surface**

| Feature | Details | Priority |
|---------|---------|----------|
| **`sdp status`** | Dashboard: features, workstreams, coverage, evidence, blocked items | Must |
| **`sdp next`** | Recommend next workstream by dependency + priority | Should |
| **`sdp log show`** | Evidence browser with filters (feature, model, date, author) | Must |
| **`sdp log show --unit 00-001-03`** | Evidence for specific workstream | Must |
| **`--output=json`** | JSON for all commands | Must |
| **Streaming progress** | Per-unit progress during execution | Should |

**CI/CD**

| Feature | Details | Priority |
|---------|---------|----------|
| **GitHub Action** | `sdp-dev/verify-action@v1` | Must |
| **PR evidence comment** | Evidence chain summary in PR | Must |
| **Provenance gate** | Block merge if AI code lacks evidence | Should |
| **GitLab CI** | Same for GitLab | Later |

**Model Routing & Economics**

| Feature | Details | Priority |
|---------|---------|----------|
| **Profile rotation** | Rotate credentials/profile order per provider with session stickiness | Must |
| **Cooldown/backoff policy** | Cool down failing profiles/models with bounded exponential backoff | Must |
| **Model fallback chain** | Fail over to configured fallback models when provider/profile is exhausted | Must |
| **Cost-aware routing** | Route by risk tier + budget policy, not only availability | Should |
| **Status/usage surface** | Show active model/profile, fallback state, and usage/cost summary in CLI status | Should |

**Outreach**

| Feature | Details | Priority |
|---------|---------|----------|
| **Article** | "Why your AI code needs both discovery AND delivery traces" | Must |
| **Demo on real project** | Consulting case or public project | Must |
| **3 external users** | Real teams using SDP | Must |

### M3 Validation

- Do people install SDP via `sdp init`?
- Is the MCP server stable in Cursor / VS Code?
- Do external users find evidence valuable?
- Is time to "aha moment" < 5 minutes?
- Does cost-aware model routing reduce spend without hurting pass rate?

**Kill:** If < 10 external users after 3 months with MCP — reach strategy failed.

---

## M4: "Tie" — Governance & Compliance

> Evidence becomes a legal artifact.
> **Timeline:** 4-6 weeks after M3 (Sep-Oct 2026). **Level:** Enterprise / regulated.

### Features → Workstreams

| Feature | Description | Est. | WS |
|---------|-------------|------|----|
| **F055**: Compliance Design Doc | Data residency, retention, RBAC, integrity, prompt privacy | ~2.5h | 2 WS |
| **F056**: Full Skills Instrumentation | @review, @deploy, @design, @idea → evidence | ~7h | 4 WS |
| **F059**: Observability Bridge Design | Deploy markers, OTel span attributes, integration spec | ~3.5h | 2 WS |
| **F061**: Data Collection & Benchmark | MTTR tracking, catch rate, AI failure taxonomy | ~5.5h | 3 WS |

### M4 Deliverables

**Full Instrumentation**

| Feature | Details | Priority |
|---------|---------|----------|
| **@vision → strategic decisions** | What market, what users, what positioning → evidence | Must |
| **@design → decomposition rationale** | Why this architecture, why these dependencies → evidence | Must |
| **@review → findings** | Drift detected, quality regression, architecture violation → evidence | Must |
| **@deploy → approval chain** | Who approved merge, what gates passed → evidence | Must |
| **All 19 skills produce evidence** | Full pipeline coverage | Must |

**Evidence Integrity**

| Feature | Details | Priority |
|---------|---------|----------|
| **Hash chain** | `prev_hash` per record — corruption detection, NOT tamper-proof | Must |
| **Append-only merge** | `.gitattributes` with merge driver | Must |
| **Honest labeling** | Documentation: what hash chain guarantees and what it doesn't | Must |

**Compliance**

| Feature | Details | Priority |
|---------|---------|----------|
| **Compliance design doc** | Data residency, retention, RBAC, integrity, prompt privacy | Must |
| **EU AI Act mapping** | Article-by-article → SDP capability | Must |
| **SOC2 mapping** | Trust Services Criteria → SDP evidence | Should |
| **DORA mapping** | ICT risk management → SDP audit trail | Should |

**Observability Bridge (Design)**

| Feature | Details | Priority |
|---------|---------|----------|
| **OTel span attributes** | Mark AI-generated code paths in traces | Must |
| **Deploy markers** | Tie evidence records to deploy events | Must |
| **Integration spec** | How SDP connects to Honeycomb / Datadog / Grafana | Should |

**Data & Analytics**

| Feature | Details | Priority |
|---------|---------|----------|
| **MTTR tracking** | Time to root cause with evidence vs without | Must |
| **Verification telemetry** | Catch rate, iteration count, model performance | Should |
| **AI failure taxonomy** | What goes wrong, by model/language/domain | Should |

### M4 Validation

- Does the compliance mapping resonate with real CISO/compliance officers?
- Does full instrumentation slow the workflow?
- Is the hash chain actually useful or just overhead?

---

## Horizon: Sports Car

Not on timeline. Features for future milestones, triggered by demand.

### SDK & Integration

| Feature | Details | Trigger |
|---------|---------|---------|
| `sdp.Verify()` | Verification engine as Go library | External integrators ask |
| `sdp.Evidence()` | Evidence bundle generation | External integrators ask |
| Provider adapters | Claude, GPT, Gemini model adapters | Multi-model demand |
| JSON-in/JSON-out API | External tool integration surface | SDK users |

### Observability (Full)

| Feature | Details | Trigger |
|---------|---------|---------|
| OTel exporter | SDP evidence → OpenTelemetry spans | Production deployments |
| Deploy correlation | Auto-link evidence ↔ deploy events | Ops teams request |
| Runtime context | Feature flags, blast radius, rollback path | Incident response need |
| Per-line attribution | AI vs human per line | Audit requirement |
| Diff-level provenance | Which lines are AI-generated vs human-edited | Audit requirement |

### Multi-Agent Coordination

| Feature | Details | Trigger |
|---------|---------|---------|
| Continuous cross-branch integration | After each WS: merge main, run acceptance test | Multiple parallel features |
| Cross-feature test matrix | "Feature A still works after Feature B's latest WS" | Large projects |
| Cost-aware testing | Full matrix for high-risk, smoke-only for low-risk | Scale |
| Cross-model review | Model A generates, Model B reviews (decorrelated errors) | High-risk code |
| Model selection policy | Route by risk level | Enterprise policy |

### Enterprise

| Feature | Details | Trigger |
|---------|---------|---------|
| Compliance export | SOC2/HIPAA/DORA-ready evidence format | Design partners validate |
| Verification certificates | Signed, timestamped | Legal requirement |
| Signed evidence | Cryptographic non-repudiation (compliance-grade) | Audit demand |
| External timestamping | Third-party timestamp authority | Legal validity |
| Risk-proportional verification | `auth/` → full trail, `components/` → light | Enterprise config |
| Team policies | "All AI PRs need evidence chain" | Team governance |
| Decomposition templates | Per-company shared patterns | Enterprise onboarding |
| On-premise | Air-gapped environments | Defense/banking |

### IDE Plugins

| Feature | Details | Trigger |
|---------|---------|---------|
| Cursor plugin | Native plan/apply/trace from Cursor | MCP insufficient |
| VS Code extension | Same | MCP insufficient |
| JetBrains plugin | Same | JetBrains user demand |

### Standard

| Feature | Details | Trigger |
|---------|---------|---------|
| SDP Evidence Format v1.0 | Published, auditor-reviewed spec | Schema stable 6+ months |
| External adoption | 2+ tools beyond SDP produce SDP evidence | Community traction |
| Industry working group | Formal standardization effort | After 50+ deployments |

### Research Track

| Question | Approach |
|----------|----------|
| How to hide complexity without sacrificing quality? | Auto-detection (solo/team/enterprise), progressive skill depth, adaptive architecture, "friday mode" |
| Can SDP work for vibe coders? | `@quick "fix the bug"` — skip decomposition, direct TDD, evidence in background |
| Adaptive architecture for project scale? | Small → flat, large → Clean Architecture, user always overrides |

---

## North Star: Real-Time Multi-Agent Coordination

Not on the roadmap. The guiding direction.

**The problem:** 10 agents and 5 humans build 5 features in parallel. Agent A changes the User interface. Agent B is mid-build using the old User interface. Nobody knows until merge time — the worst possible moment to find out.

**The vision:** Real-time coordination where every participant — agent or human — has live awareness of what others are doing. Not just scope (which files) but intent (what they're trying to achieve) and state (what they've changed so far).

**Path (each step is useful on its own):**

```
M2: Scope collision ──► M2: Shared contracts ──► Horizon: Cross-branch ──► North Star
    (static, plan time)     (design time)          (build time)            (real-time)
```

---

## Feature → Milestone Map

Complete map: which feature lives where. **NEW** = no workstreams yet.

| Feature | Name | Milestone | WS | Status |
|---------|------|-----------|-----|--------|
| F054 | Evidence Layer | **M1** | 11 | Backlog |
| F063 | Guardian Hooks | **M1** | 4 | Backlog |
| F064 | Unified Task Resolver | **M1** | 4 | Backlog |
| F067 | Repo Hardening | **M1** | 14 | Backlog |
| F068 | UX Foundation & First-Run | **M1** | 5 | Backlog |
| F070 | Failure & Recovery UX | **M1** | 5 | Backlog |
| F075 | Self-Healing Doctor | **M1** | 6 | Backlog |
| F076 | Guided Onboarding Wizard | **M1** | 5 | Backlog |
| NEW | Acceptance Test Gate | **M1** | TBD | Not planned |
| F060 | Shared Contracts | **M2** | 3 | Backlog |
| F071 | Team UX & Collaboration | **M2** | 5 | Backlog |
| F073 | Trust & Explainability | **M2** | 5 | Backlog |
| F077 | Runtime Hooks Platform | **M2** | 5 | Backlog |
| F078 | Resilient Runtime Layer | **M2** | 4 | Backlog |
| NEW | Schema Consolidation | **M2** | TBD | Not planned |
| NEW | Review Protocol Self-Check | **M2** | TBD | Not planned |
| NEW | Architecture Awareness | **M2** | TBD | Not planned |
| NEW | Scope Collision Detection | **M2** | TBD | Not planned |
| NEW | Agentic Discovery Contracts | **M2** | TBD | Not planned |
| NEW | Discovery → Delivery Contract Bridge | **M2** | TBD | Not planned |
| NEW | Explore/Commit Path Policy | **M2** | TBD | Not planned |
| F057 | CLI plan/apply/log | **M3** | 4 | Backlog |
| F058 | CI/CD GitHub Action | **M3** | 3 | Backlog |
| F069 | Next-Step Engine | **M3** | 5 | Backlog |
| F072 | Interop & Migration | **M3** | 5 | Backlog |
| F074 | Layered OSS Packaging | **M3** | 6 | Backlog |
| F079 | Model Routing & Economics | **M3** | 5 | Backlog |
| NEW | MCP Server | **M3** | TBD | Not planned |
| NEW | Runtime-Agnostic Agent Adapter | **M3** | TBD | Not planned |
| F055 | Compliance Design Doc | **M4** | 2 | Backlog |
| F056 | Full Skills Instrumentation | **M4** | 4 | Backlog |
| F059 | Observability Bridge Design | **M4** | 2 | Backlog |
| F061 | Data Collection & Benchmark | **M4** | 3 | Backlog |

**Total backlog:** ~159 workstreams + NEW features TBD

---

## Timeline

```
2026
───────────────────────────────────────────────────────────────────

 Mar─Apr    M1 "T-shirt"
            F054 F063 F064 F067 F068 F070 F075 F076 + Acceptance Test
            Evidence log · self-healing doctor · onboarding wizard · smoke test
            ══════════════════════

 May─Jun    M2 "Smart Casual"
            F060 F071 F073 F077 F078 + Schema + Review + Arch + Collision + D2D Bridge
            Formal schema · runtime hooks · retry/fallback · discovery contracts · delivery bridge
                   ══════════════════════

 Jul─Aug    M3 "Blazer"                                EU AI Act ←
            F057 F058 F069 F072 F074 F079 + MCP + Adapters full
            CLI · CI/CD · model router · sdp init · MCP   enforcement
            First external users
                          ══════════════════════

 Sep─Oct    M4 "Tie"
            F055 F056 F059 F061
            Full instrumentation · Hash chain
            Compliance mapping · OTel design
                                 ══════════════════════

 Nov+       Horizon
            SDK · OTel · Cross-model review · IDE plugins
            Enterprise features · Standard
                                        ═══════════════════
```

---

## Metrics

| Metric | M1 | M2 | M3 | M4 | Horizon |
|--------|----|----|----|----|---------|
| Evidence records | 100 | 500 | 2K | 10K | 100K |
| External users | 0 | 1-2 | 10+ | 30+ | 100+ |
| Supported tools | 1 | 1 | 4+ | 4+ | 4+ |
| Skills instrumented | 3 | 5 | 5 | All 19 | All + SDK |
| Acceptance test pass rate | Baseline | 90%+ | 90%+ | 90%+ | Published |
| Schema version | Draft | v0.1 | v0.2 | v0.5 | v1.0 |
| Repos with evidence | 1 | 3 | 20 | 100 | 1000 |
| Commit-path defect reduction vs baseline | Baseline | 20%+ | 25%+ | 30%+ | 30%+ |
| Explore-path cycle time vs commit-path | Baseline | 30% faster | 40% faster | 40% faster | 40% faster |

---

## Kill Criteria

> **Primary:** After 100 builds with evidence: if `sdp trace` never helped in a real incident — evidence isn't adding value.

> **Secondary:** After 6 months: if 0 external users — the product doesn't resonate with the market.

> **Tertiary:** After M3: if no compliance officer showed interest — the governance angle doesn't work.

> **Acceptance:** After 50 builds: if acceptance test pass rate isn't measurably higher than baseline vibe-coding — orchestration isn't adding value.

> **Efficiency:** By end of M2: if commit-path overhead is not justified by defect reduction + evidence value, narrow scope to high-stakes only.

---

## Risks

| Risk | Mitigation | Milestone |
|------|-----------|-----------|
| SDP produces non-working code with passing gates | Acceptance test gate | M1 |
| Evidence log adds friction | Automatic emission, no opt-in | M1 |
| Over-architecture for small projects | Progressive disclosure, adaptive depth | Research |
| Models improve, decomposition unnecessary | Provenance valuable regardless | All |
| Hash chain oversold as tamper-proof | Honest labeling: corruption detection only | M4 |
| Evidence never consulted during incidents | OTel bridge, deploy markers | M4 |
| Schemas diverge again | Single namespace, CI validation | M2 |
| Kiro/GitHub ships native evidence | Open format + tool-agnostic = standard wins | M3 |
| Log bloats repo | Budget per event, compression, archival policy | M1 |
| Scope collision detection is noisy | Signal, not block; tunable thresholds | M2 |
| No external users after 6 months | Kill or pivot | M3 |
| Discovery outputs don't translate into delivery work | Discovery → Delivery contract bridge with required acceptance fields | M2 |
| Tool/vendor lock-in pressure | Runtime-agnostic adapter interface + open evidence schema | M3 |

---

## What We Don't Build (confirmed)

- No more agents (19 is enough)
- No custom IDE (we're a protocol, not an IDE)
- No visual dashboard until Horizon
- No plugin marketplace until community exists
- No multi-stage review for simple checks
- No full decentralization of coordination

---

## Supersedes

This document replaces:
- `docs/vision/ROADMAP.md` v7.0 (P0-P3 framework)
- `docs/vision/FEATURES.md` v7.0 (feature map)
- `docs/roadmap/ROADMAP.md` v3.0 (quarterly roadmap)

---

*SDP Unified Roadmap v8.1 — February 15, 2026*
*From toy to sports car. Each milestone is self-contained. Evidence is the thread.*
