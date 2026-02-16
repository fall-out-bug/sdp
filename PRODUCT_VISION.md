# PRODUCT_VISION.md

> **Last Updated:** 2026-02-14
> **Version:** 3.2 (Multi-Agent, Public OSS Layered Platform)

## Mission

Transform SDP into a **multi-agent, multi-layer platform** for autonomous software delivery: protocol-first planning, guided CLI UX, and resilient orchestration.

### Architecture Evolution (v2.0 → v3.0)

**Original Vision (v2.0):** Language-agnostic Claude Plugin for protocol enforcement
**Pivot to Multi-Agent (v3.0):** Orchestration system for autonomous development

**Rationale for Pivot:**
- Claude Plugin marketplace not yet available (platform limitation)
- Multi-agent architecture enables autonomous execution (higher value)
- Strategic planning (@vision) and reality checking (@reality) prevent drift
- Parallel execution achieves 5x speedup vs sequential workflows

## Platform Topology (Public OSS Scope)

SDP is currently delivered as a monorepo, but product direction is explicit: split into multiple repositories once interfaces stabilize.

| Layer | Future Repository | Responsibility |
|-------|-------------------|----------------|
| Protocol | `sdp-protocol` | Core spec, schemas, event formats, frontmatter contracts, compatibility policy |
| CLI | `sdp-cli` | Developer-facing command UX, onboarding, status/help/next-step flow |
| Orchestrator | `sdp-orchestrator` | Planning/execution engine, dependency graph, dispatch, checkpoint/recovery |

Private extensions may exist outside this repository, but are intentionally out of scope for public OSS planning docs.

## Progressive Adoption Model (L0-L2, Public OSS)

SDP must be adoptable without "the whole spaceship." The operating model is progressive:

| Level | Product Scope | Primary Buyer/User | Distribution Channel | License |
|-------|---------------|--------------------|----------------------|---------|
| `L0` | Protocol assets only (prompts, guides, templates, schemas) | Individual developers and teams validating method | Claude plugin + prompt package | MIT |
| `L1` | Safety and evidence layer (hooks, guard, traces, provenance) | Teams needing control and auditability | Homebrew package | MIT |
| `L2` | Orchestrator and core SDP tooling | Teams scaling spec-driven execution | Homebrew package | MIT |

### Adoption Principle

- Every level must be independently valuable.
- `L0-L2` must be installable without non-public dependencies.
- Upgrade path is additive (`L0` -> `L1` -> `L2`), not all-or-nothing.
- Feature work must declare which level it belongs to.

### Monorepo-to-Multi-Repo Rules

- Keep layer boundaries explicit in workstreams and architecture docs
- Evolve contracts before extracting code
- Avoid hidden cross-layer dependencies that block future split
- Treat extraction readiness as a roadmap KPI, not an ad-hoc refactor

## Public OSS Commitment

- Public SDP layers in this repository (`L0-L2`) remain MIT.
- Public docs in this repository track public OSS scope only.
- Non-public extensions must not become a hard dependency for public layers.

## Users

1. **Development Teams Using AI Agents**
   - Want autonomous feature execution
   - Need quality assurance before merge
   - Require fault tolerance and checkpointing

2. **Strategic Planners**
   - Product managers defining roadmaps
   - Tech leads planning quarterly goals
   - Architects analyzing codebase health

3. **Solo Developers**
   - Want AI assistance for feature planning
   - Need reality checks on code quality
   - Require autonomous execution of repetitive tasks

## Success Metrics (v3.2 - Public OSS Platform)

### Delivered ✅
- [x] Multi-agent orchestration with 19 specialized agents
- [x] Parallel execution achieving 4.96x speedup
- [x] Fault tolerance (circuit breaker + checkpoint/resume)
- [x] Strategic planning via @vision (7 expert agents)
- [x] Codebase analysis via @reality (8 expert agents)
- [x] Two-stage quality review (implementer → spec reviewer → quality)
- [x] Agent synthesis for conflict resolution
- [x] Progressive disclosure for reduced question fatigue
- [x] 83.2% test coverage across 26 workstreams

### Deferred 🔄
- [ ] Claude Plugin marketplace distribution (awaiting platform support)
- [ ] Zero runtime dependency (Go binary required for orchestration)
- [ ] Prompts-only distribution (binary provides agent coordination)

### Maintained ✅
- [x] Language-agnostic validation (works on any project)
- [x] Backward compatibility with Python SDP (protocol unchanged)

## Strategic Tradeoffs (v3.1)

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| **Architecture** | Multi-agent orchestration | Enables autonomous execution, parallel speedup, fault tolerance |
| **Planning** | 4-level model (@vision/@reality/@feature/@oneshot) | Strategic planning → Reality check → Feature planning → Autonomous execution |
| **Execution** | Parallel dispatcher with DAG | 4.96x speedup via dependency-aware parallelization |
| **Quality** | Two-stage review + synthesis | Implementer → Spec reviewer → Quality reviewer with conflict resolution |
| **Fault Tolerance** | Circuit breaker + checkpoint | Crash-safe execution with automatic recovery |
| **Language Support** | Go binary with protocol enforcement | Language-agnostic rules via prompts, Go provides performance |
| **Distribution** | Binary (Go) + Claude skills | Binary for orchestration, skills for protocol definition |
| **Repository Strategy** | Monorepo now, split later | Faster iteration now; clearer ownership and scaling after contract maturity |
| **Adoption Model** | Progressive (`L0-L2`) | Low-friction entry with independent OSS adoption layers |

## Architecture Overview (v3.x Monorepo Runtime)

```
Strategic Level                 Analysis Level                 Feature Level                Execution Level
┌──────────────────┐           ┌──────────────────┐         ┌──────────────────┐         ┌──────────────────┐
│     @vision      │           │    @reality      │         │    @feature      │         │    @oneshot      │
│  (strategic)     │           │  (codebase anal) │         │   (planning)     │         │   (execution)    │
└──────────────────┘           └──────────────────┘         └──────────────────┘         └──────────────────┘
         │                              │                             │                             │
         ▼                              ▼                             ▼                             ▼
  7 Expert Agents             8 Expert Agents              @idea + @design            @build (all WS)
  (product analysis)           (codebase analysis)          (requirements + WS)         (implement)
         │                              │                             │                             │
         ▼                              ▼                             ▼                             ▼
  Product Artifacts            Reality Report               workstreams                 @review + @deploy
  (VISION, PRD, ROADMAP)      (health, gaps, debt)         (00-FFF-SS.md)              (quality + merge)
```

## Delivered Features (F052)

### Phase 1A: Strategic Planning (@vision)
- 7 expert agents (product, market, technical, UX, business, growth, risk)
- PRODUCT_VISION.md, PRD.md, ROADMAP.md generation
- Feature extraction from PRD with priority tagging

### Phase 1B: Codebase Analysis (@reality)
- 8 expert agents (architecture, quality, testing, security, performance, docs, debt, standards)
- Language/framework detection (Go, Python, JS, TS, Rust, Java)
- Health scoring and gap analysis
- Vision vs reality drift detection

### Phase 2: Quality Lock-in (Two-Stage Review)
- Implementer agent (TDD discipline: Red → Green → Refactor)
- Spec compliance reviewer (DO NOT TRUST pattern, verifies actual code)
- @build orchestration (3 stages with max 2 retries per stage)

### Phase 3: Speed Track (Parallel Execution)
- Parallel dispatcher with Kahn's algorithm
- Dependency graph from workstream files
- 4.96x speedup (5 WS: 55ms → 11ms, 10 WS: 109ms → 22ms)

### Phase 4: Synthesis Track (Agent Coordination)
- Agent proposal system with confidence scoring
- Synthesizer with priority-based rules (unanimous → domain expertise → quality gate → merge → escalate)
- Hierarchical supervisor for specialist agent coordination
- Conflict type detection (major, medium, minor)

### Phase 5: UX Track (Progressive Disclosure)
- 3-question cycles with trigger points
- TMI detection and --quiet mode
- Target: 12-27 questions per feature (down from unbounded)
- Verbosity tiers (--quiet, --verbose, --debug)

### Phase 6: Documentation Track
- Agent catalog (21 agents documented)
- Migration guide (v3.x to v4.0)
- Updated CLAUDE.md with 4-level planning model
