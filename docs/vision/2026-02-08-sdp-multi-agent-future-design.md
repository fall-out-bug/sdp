# SDP Multi-Agent Future: Protocol for the Post-Swarm Era

> **Status:** Research complete
> **Date:** 2026-02-08
> **Goal:** Design SDP's evolution from a single-provider coding framework to a universal multi-agent orchestration protocol that survives the next 5 evolutionary steps of AI coding.

---

## Table of Contents

1. [Overview](#overview)
2. [Protocol vs Implementation](#1-protocol-vs-implementation)
3. [Agent Contract Standard](#2-agent-contract-standard)
4. [Execution Runtime Abstraction](#3-execution-runtime-abstraction)
5. [State & Checkpoint Portability](#4-state--checkpoint-portability)
6. [Workstream as Universal Unit](#5-workstream-as-universal-unit)
7. [Trust & Verification in Heterogeneous Swarms](#6-trust--verification-in-heterogeneous-swarms)
8. [Economic Model](#7-economic-model)
9. [Competitive Landscape & Moat](#8-competitive-landscape--moat)
10. [Progressive Autonomy](#9-progressive-autonomy)
11. [SDK/Plugin Ecosystem](#10-sdkplugin-ecosystem)
12. [Implementation Roadmap](#implementation-roadmap)
13. [Success Metrics](#success-metrics)

---

## Overview

### The Thesis

SDP predicted the multi-agent coding wave before "Swarm" was a buzzword. The repo already has:
- 19 specialized agents with adversarial review chains
- Parallel dispatcher with DAG-based dependency resolution
- Synthesis engine resolving multi-agent conflicts
- Four-level planning model (@vision → @reality → @feature → @oneshot)

But SDP is currently **a framework tightly coupled to Claude Code**. To survive the next 5 steps of the agentic evolution, it must become **a protocol that any agent can implement**.

### The 5 Steps Ahead

| Step | Era | What Happens | SDP Must... |
|------|-----|-------------|-------------|
| **1** (Now) | Multi-IDE | Cursor, Claude Code, OpenCode coexist | Abstract provider layer |
| **2** (6mo) | API Agents | GLM, Kimi, QWEN agents via API | Support remote execution backends |
| **3** (12mo) | Agent Marketplaces | Teams pick best agent per task | Have agent contracts + economic routing |
| **4** (18mo) | Autonomous Swarms | Agents self-organize without humans | Progressive autonomy dials |
| **5** (24mo) | Protocol Wars | Competing standards for agent coordination | Be the OCI of agent orchestration |

### Key Decisions

| Aspect | Decision |
|--------|----------|
| Protocol vs Implementation | Strangler Fig: extract `protocol/` from existing code, generate provider dirs |
| Agent Contracts | Hybrid YAML contracts + JSON Schema validation with scored capability matching |
| Runtime Abstraction | `AgentBackend` interface with per-backend circuit breakers |
| State Portability | Event log with pluggable transport (File → Git → KV) |
| Workstream Granularity | Layered: Goal → Workstream (coordination) → Task (agent-internal) |
| Trust Model | Defense-in-depth: output gates + trust tiers + audit trail + reproducibility |
| Economic Model | Full economic router: budget constraints + quality-tiered routing + feedback loop |
| Strategic Positioning | Protocol Core + Framework Distribution (the Kubernetes model) |
| Autonomy | Per-scope autonomy dials in quality-gate.toml with trust score |
| Ecosystem | Plugin format with GitHub-based registry, `sdp plugin add` CLI |

---

## 1. Protocol vs Implementation

> **Expert:** Sam Newman (bounded contexts, microservices architecture)

### The Problem

SDP is coupled to Claude Code through `.claude/skills/`, `.claude/agents/`, and `.claude/settings.json`. The `sdp-plugin/prompts/` directory already contains provider-neutral prompts — but `.claude/` is treated as canonical, not generated.

### The Architecture

**Four-layer decomposition:**

| Layer | Contains | Changes When |
|-------|----------|-------------|
| **Protocol** (`protocol/`) | Schemas, canonical prompts, quality gates, hooks | SDP protocol evolves |
| **Engine** (`src/sdp/`, `internal/`) | Go binary (graph, dispatcher, guard, quality) | Core capabilities change |
| **Provider** (`providers/{name}/`) | Generation scripts, overrides, config templates | Provider API changes |
| **Generated** (`.claude/`, `.cursor/`, `.opencode/`) | Provider-specific artifacts | Regenerated, not hand-edited |

### Directory Structure

```
sdp/
├── protocol/                           # THE protocol layer (source of truth)
│   ├── schema/
│   │   ├── workstream.schema.json     # Already exists
│   │   ├── agent-contract.schema.json # NEW
│   │   └── quality-gate.schema.json   # NEW
│   ├── prompts/                        # Canonical prompt text
│   │   ├── skills/                     # Provider-neutral skill definitions
│   │   └── agents/                     # Provider-neutral agent prompts
│   ├── quality/
│   │   └── ci-gates.toml              # Already exists
│   └── hooks/
│       ├── pre-commit.sh
│       └── pre-push.sh
│
├── providers/                          # Per-provider adapter layer
│   ├── claude/
│   │   ├── generate.sh                # Generates .claude/ from protocol/
│   │   └── overrides/                  # Claude-specific (Task(), settings.json)
│   ├── cursor/
│   │   ├── generate.sh
│   │   └── overrides/                  # .cursorrules generation
│   └── opencode/
│       ├── generate.sh
│       └── overrides/                  # opencode.json generation
│
├── .claude/                            # GENERATED — DO NOT EDIT
├── .cursor/                            # GENERATED — DO NOT EDIT
└── .opencode/                          # GENERATED — DO NOT EDIT
```

### Migration Path (Strangler Fig)

1. **Phase 1**: Create `protocol/prompts/` by copying `sdp-plugin/prompts/`. Write `providers/claude/generate.sh`.
2. **Phase 2**: Make `.claude/`, `.cursor/`, `.opencode/` generated. CI validates: `make generate && git diff --exit-code`.
3. **Phase 3**: New providers (Kimi, QWEN, GLM) only need `providers/{name}/generate.sh`.
4. **Phase 4**: `sdp init --provider=cursor` generates the right directory on install.

### Key Insight

> "The prompt IS the product. In SDP, the protocol is prompt text — the natural language instructions. Provider differences are just packaging. The base prompt (the 'what to do') is 80%+ identical across providers."

---

## 2. Agent Contract Standard

> **Expert:** Theo Browne (type-safe API design)

### The Problem

Current agents are defined implicitly in markdown. The `synthesis.Agent` interface only declares `ID()`, `Available()`, `Consult()`. No capabilities, no quality guarantees, no cost profile. The `AgentSpawner` hardcodes valid types as a string enum.

### The Solution: YAML Contract + JSON Schema Validation

Agent contracts are authored in YAML, validated against a JSON Schema meta-schema. The orchestrator loads contracts, validates them, and builds a capability registry with **scored matching**.

### Contract Format

```yaml
# .agents/build-claude-sonnet.agent.yaml
sdp_agent_contract: v1
agent_id: build-claude-sonnet
version: 1.0.0
provider: anthropic
model: claude-sonnet-4-5
tier: T1

capabilities:
  - id: tdd
    level: expert          # expert | proficient | basic | none
    evidence: "95% success rate on 127 workstreams"
  - id: code_generation
    level: expert
    languages: [go, python, typescript]
  - id: test_generation
    level: expert
  - id: refactoring
    level: proficient
  - id: review
    level: none

input:
  required:
    - name: workstream_spec
      type: file_path
    - name: scope_files
      type: file_path[]

output:
  required:
    - name: verdict
      type: enum[PASS, FAIL]
    - name: changed_files
      type: file_path[]
    - name: test_results
      type: test_report

quality_guarantees:
  coverage_min: 80
  file_loc_max: 200
  complexity_max: 10
  success_rate: 0.90

cost:
  per_million_tokens: 3.00
  avg_tokens_per_task: 50000
  latency:
    p50: 120s
    p99: 600s

negotiation:
  degraded_modes:
    no_tool_use:
      drops: [tdd]
      retains: [code_generation, test_generation]
    small_context:
      max_scope_files: 3
      drops: [refactoring]
  alternatives:
    - agent_id: build-gpt4o
      priority: 2
    - agent_id: build-qwen-72b
      priority: 3
```

### Capability Matching

The orchestrator asks: "I need TDD in Go with >=80% coverage." It gets back a scored list:

```
1. build-claude-sonnet  score=0.95  degradations=[]          cost=$0.15
2. build-gpt4o          score=0.88  degradations=[]          cost=$0.25
3. build-qwen-72b       score=0.72  degradations=[refactor]  cost=$0.00
```

### Implementation

```go
// Extend existing synthesis.Agent interface
type AgentContract struct {
    AgentID       string
    Capabilities  []Capability
    Input         Schema
    Output        Schema
    Quality       QualityGuarantees
    Cost          CostProfile
    Negotiation   NegotiationSpec
}

type CapabilityScore struct {
    Agent        string
    Score        float64   // 0.0 - 1.0
    Degradations []string  // what's lost
    Cost         float64   // estimated cost
}

func (r *AgentRegistry) MatchAgent(required RequiredCapabilities) ([]CapabilityScore, error)
```

---

## 3. Execution Runtime Abstraction

> **Expert:** Kelsey Hightower (infrastructure, platform design)

### The Problem

`ExecuteFunc func(wsID string) error` — a single function signature that can only execute locally. No agent identity, no streaming, no transport abstraction, no per-agent failure handling.

### The Solution: AgentBackend Interface

```go
type AgentBackend interface {
    ID() string
    Kind() BackendKind          // Local, HTTP, gRPC, WebSocket
    Execute(ctx context.Context, req *ExecutionRequest) (*ExecutionResult, error)
    Stream(ctx context.Context, req *ExecutionRequest) (<-chan ProgressEvent, error)
    Capabilities() AgentCapabilities
    HealthCheck(ctx context.Context) error
}

type BackendKind int
const (
    BackendLocal BackendKind = iota
    BackendHTTP
    BackendGRPC
    BackendWebSocket
)
```

### Per-Backend Circuit Breakers

The current single global circuit breaker is catastrophically wrong for heterogeneous backends. If GLM rate-limits, it shouldn't block Claude.

```go
type FailureClass int
const (
    FailureTransient  FailureClass = iota  // Timeout, 503 — retry soon
    FailureRateLimit                        // 429 — exponential backoff
    FailureAuth                             // 401/403 — don't retry, escalate
    FailurePermanent                        // 400 — don't retry
)

type Dispatcher struct {
    graph    *DependencyGraph
    backends *BackendRegistry
    breakers map[string]*AgentCircuitBreaker  // per-backend
    store    CheckpointStore
}
```

### Architecture Diagram

```
┌────────────────────────────────────────────────────────────────┐
│                    SDP Dispatcher (Coordinator)                  │
│                                                                  │
│  ┌──────────┐   ┌──────────────┐   ┌────────────────────────┐  │
│  │ DAG      │──▶│ Scheduler    │──▶│ Backend Registry       │  │
│  │ (ready?) │   │ (route to?)  │   │ ┌────────────────────┐ │  │
│  └──────────┘   └──────────────┘   │ │ LocalBackend       │ │  │
│                                     │ │ HTTPBackend(Claude)│ │  │
│  ┌──────────┐   ┌──────────────┐   │ │ HTTPBackend(GLM)   │ │  │
│  │ Checkpoint│◀─▶│ Per-Backend  │   │ │ HTTPBackend(QWEN)  │ │  │
│  │ Store     │   │ Breakers     │   │ └────────────────────┘ │  │
│  └──────────┘   └──────────────┘   └────────────────────────┘  │
└────────────────────────────────────────────────────────────────┘
         │                │                     │
         ▼                ▼                     ▼
   Local Agent      Claude API            GLM/QWEN API
   (goroutine)      (remote)              (remote)
```

### Migration Path

| Version | Change |
|---------|--------|
| v0.9 | Define `AgentBackend` interface. `LocalBackend` wraps existing `ExecuteFunc`. Zero behavior change. |
| v0.9.1 | Add `HTTPBackend` for one remote agent. Per-backend circuit breakers. |
| v1.0 | Streaming progress. Unify `synthesis.Agent` and `AgentBackend`. |
| v1.1 | Reconciler loop for self-healing. Capability-based routing. |

---

## 4. State & Checkpoint Portability

> **Expert:** Martin Kleppmann (distributed systems, event sourcing)

### The Problem

Checkpoints are local JSON files with atomic-rename durability. When Agent B (remote GLM) completes a workstream, it has no way to update the local checkpoint. Two dispatchers running simultaneously corrupt state.

### The Solution: Event Log with Pluggable Transport

```go
type CheckpointStore interface {
    RecordEvent(event ExecutionEvent) error
    GetState(featureID string) (*CheckpointState, error)
    AcquireLock(featureID string, holder string, ttl time.Duration) (bool, error)
    ReleaseLock(featureID string, holder string) error
}

// Implementations:
type FileCheckpointStore struct { ... }  // Current behavior, single machine
type GitCheckpointStore struct { ... }   // git commit + push for multi-machine
type KVCheckpointStore struct { ... }    // Redis/etcd for hosted SDP (future)
```

### Event Model

```go
type ExecutionEvent struct {
    EventID      string          `json:"event_id"`
    Timestamp    time.Time       `json:"timestamp"`
    AgentID      string          `json:"agent_id"`
    WorkstreamID string          `json:"workstream_id"`
    EventType    EventType       `json:"event_type"`   // Started, Completed, Failed
    Payload      json.RawMessage `json:"payload"`
}
```

### Key Insight

> "Separate the data model from the transport. The current Checkpoint struct captures the right things. The problem isn't the data, it's the I/O."

---

## 5. Workstream as Universal Unit

> **Expert:** Martin Fowler (evolutionary architecture)

### The Problem

Static LOC limits (SMALL <500, MEDIUM 500-1500, LARGE >1500) become increasingly mismatched to agent capability. A 2027 agent handling 5000 LOC is forced through 3-4 artificial WS boundaries.

### The Solution: Layered Abstraction

```
Goal (what to achieve)      ← Human/AI sets this
  └── Workstream (coordination) ← SDP manages this
        └── Task (execution)    ← Agent decides internally
```

Workstreams become **contract boundaries** — defined by acceptance criteria and scope, not LOC limits. LOC limits become advisory output constraints checked by quality gates, not input constraints blocking `@design`.

```go
type ExecuteResult struct {
    WorkstreamID string
    Success      bool
    Error        error
    Duration     int64
    LOCProduced  int     // NEW: actual output size
    TaskCount    int     // NEW: internal decomposition count
    TokensUsed   int64   // NEW: cost tracking
}
```

### Key Insight

> "Make the unit of change match the unit of value. A workstream's value is its acceptance criteria, not its LOC count."

---

## 6. Trust & Verification in Heterogeneous Swarms

> **Expert:** Troy Hunt (security, defense in depth)

### The Problem

`Proposal.Confidence` is the agent grading its own homework. In a heterogeneous swarm with unknown providers, this is a vulnerability.

### The Solution: Defense-in-Depth Hybrid

| Layer | Mechanism | Always | High-Risk Only |
|-------|-----------|--------|----------------|
| **L1: Output Gates** | Tests, linting, SAST, coverage | Yes | Yes |
| **L2: Trust Tiers** | Permission scoping by provider | Yes | Yes |
| **L3: Audit Trail** | Signed, immutable action log | Yes | Yes |
| **L4: Reproducibility** | Second-agent verification | No | Yes |
| **L5: Human Review** | Escalation for critical paths | No | Yes |

### Trust Tiers

```go
type TrustLevel int
const (
    TrustUntrusted  TrustLevel = iota  // Can only propose, not execute
    TrustRestricted                     // Execute in sandbox, review required
    TrustStandard                       // Execute with quality gates
    TrustPrivileged                     // Execute security-sensitive code
)

type AgentPermissions struct {
    AllowedPaths    []string   // File glob patterns
    DeniedPaths     []string   // e.g., "**/.env", "**/credentials*"
    CanExecuteTests bool
    CanModifyDeps   bool
    CanAccessNetwork bool
}
```

### Adversarial Scenario Protection

A compromised agent injecting malicious code is caught at multiple layers:
1. **L1**: SAST catches known vulnerability patterns
2. **L2**: Agent can only touch files within its permission scope
3. **L3**: Audit trail identifies which agent introduced the change
4. **L4**: Reproducibility check with different provider catches divergent behavior
5. **L5**: Security-sensitive paths always require human review

---

## 7. Economic Model

> **Expert:** Nir Eyal (optimization, behavioral economics)

### The Problem

Zero cost awareness. The dispatcher treats all executions equally. Real-world pricing varies 100x across providers.

### The Solution: Full Economic Router

```
Workstream Spec → Budget Check → Quality Router → Cost Optimizer → Execute → Record Cost
                       ↑                                                          │
                       └──────────── Feedback Loop ◀──────────────────────────────┘
```

### Quality-Tiered Routing

```go
type QualityTier int
const (
    TierDraft    QualityTier = iota  // Tests, docs → cheap agent
    TierStandard                      // Normal code → mid-tier
    TierCritical                      // Security, core → premium
    TierReview                        // Verification → premium + reproducibility
)
```

### Budget Declaration

```yaml
# In feature spec or quality-gate.toml
budget:
  feature_total: 50.00      # USD max for entire feature
  per_workstream_max: 15.00  # USD max per workstream
  alert_threshold: 0.80      # Alert at 80% spend
  hard_limit: false           # Warn, don't stop

quality_overrides:
  "00-053-00": critical      # Contract synthesis = premium
  "00-053-01": standard      # Code analysis = mid-tier
  "00-053-07": draft         # Documentation = cheap
```

### Routing Algorithm

```go
func (router *EconomicRouter) Route(ws WorkstreamSpec, budget *Budget) (*RoutingDecision, error) {
    // 1. Estimate tokens
    // 2. Classify quality tier (auto from WS metadata)
    // 3. Filter providers by capability + budget
    // 4. Score: quality / normalized_cost
    // 5. Return Pareto-optimal choice
}
```

---

## 8. Competitive Landscape & Moat

> **Expert:** Sam Newman (strategic positioning)

### SDP vs Competitors

| Competitor | Their Strength | SDP's Unique Angle |
|------------|----------------|-------------------|
| **OpenAI Swarm** | Lightweight handoff | SDP has DAG dispatch, circuit breakers, checkpoints |
| **LangGraph** | General graph orchestration | SDP encodes domain knowledge (TDD, quality gates) into execution |
| **CrewAI** | Easy role-based setup | SDP's agents are adversarial (implementer vs reviewer) |
| **AutoGen** | Research flexibility | SDP orchestrates structured work units, not conversations |
| **Devin** | Fully autonomous product | SDP is a protocol; Devin is a product. SDP makes YOUR tool autonomous |
| **Aider** | Excellent pair programming | SDP provides orchestration Aider lacks. Complementary. |

### The One-Line Differentiator

> SDP is the only system that treats AI coding as a **verified manufacturing process** (spec → build → adversarial review → synthesis) rather than a conversation or a single agent execution.

### The Moat (3 Layers)

1. **Weakest: Technical implementation** — reproducible in weeks
2. **Medium: Protocol specification** — encodes methodology, harder to copy
3. **Strongest: Accumulated workflow knowledge** — 24 skills, 19 agents, adversarial review chains, progressive disclosure patterns. This is like Rails' moat: not MVC (anyone can do that), but the accumulated opinions about how things *should* work.

### What SDP Should Be

**Protocol Core + Framework Distribution** (the Kubernetes model):

1. **SDP Protocol Spec** (the standard): Workstream format, agent contract interface, checkpoint format. Published independently. Anyone can implement.
2. **SDP Framework** (the reference implementation): Go binary, skills, agents, dispatcher. Opinionated "Rails" that implements the protocol.

### Partnership Strategy

| Partner | Relationship |
|---------|-------------|
| Claude Code / Anthropic | Primary runtime host — propose SDP workstreams as a standard |
| Cursor | Secondary runtime host — formalize integration |
| Aider | Build executor — dispatch individual WS to Aider |
| LangGraph | Optional infrastructure layer under the hood |
| GitHub Actions | Quality gate execution in CI |

---

## 9. Progressive Autonomy

> **Expert:** Dan Abramov (progressive disclosure, developer experience)

### The Autonomy Spectrum

```
Level 0: ██████████ Human writes everything        — NOT SDP's concern
Level 1: ██████████ AI assists (autocomplete)       — NOT SDP's concern
Level 2: ████████░░ AI executes with approval        — @build (solid)
Level 3: ██████░░░░ AI executes, human reviews       — @oneshot (works)
Level 4: ████░░░░░░ AI executes AND reviews          — needs auto-remediation loop
Level 5: ██░░░░░░░░ AI handles everything            — needs production feedback loop
```

### Autonomy Dials in Config

```toml
# In quality-gate.toml
[autonomy]
default_level = 3

[autonomy.overrides]
"**/auth/**" = 2        # Security files: human approval
"**/tests/**" = 4       # Tests: fully autonomous
"**/docs/**" = 5        # Docs: no approval needed

[autonomy.escalation]
consecutive_failures = 3          # 3 fails → force human
coverage_drop_threshold = 5       # 5% coverage drop → human
security_finding = true           # Any security issue → human
```

### Trust Score

After every execution, track:
```json
{
  "successful_builds": 47,
  "failed_builds": 3,
  "calculated_trust": 0.92,
  "recommended_level": 4,
  "current_level": 3
}
```

When trust crosses thresholds (0.8→L3, 0.9→L4, 0.95→L5), SDP suggests bumping autonomy. Circuit breaker pattern applied to autonomy.

### DX at Each Level

| Level | Interaction Model |
|-------|-------------------|
| **2** | Terminal-centric: watch TodoWrite, approve stages |
| **3** | Fire-and-forget: check checkpoint files for status |
| **4** | PR-centric: interact via GitHub, not terminal |
| **5** | Dashboard-centric: metrics, not code review |

### Key Insight

> "The DX transition is: Terminal → Files → PRs → Dashboards. Each autonomy level shifts the human interaction surface further from code."

---

## 10. SDK/Plugin Ecosystem

> **Expert:** Kelsey Hightower (platform design)

### Plugin Format

```
sdp-terraform-deploy/
├── plugin.yaml              # Manifest
├── skills/
│   └── terraform-deploy/
│       └── SKILL.md
├── agents/
│   └── terraform-agent.md
├── quality-gates/
│   └── terraform.toml       # Extends quality gate config
└── README.md
```

### Plugin Manifest

```yaml
name: sdp-terraform-deploy
version: 1.0.0
sdp_compatibility:
  min: "0.9.0"
  max: "1.x"
provides:
  skills:
    - name: terraform-deploy
      path: skills/terraform-deploy/SKILL.md
  agents:
    - name: terraform-agent
      path: agents/terraform-agent.md
  quality_gates:
    - name: terraform
      path: quality-gates/terraform.toml
```

### Registry Strategy

1. **Now**: Git-based (plugins are repos, install via `sdp plugin add github.com/org/plugin`)
2. **50+ plugins**: GitHub-based registry repo (like Homebrew taps)
3. **1000+ users**: Hosted registry at `plugins.sdp.dev`

### Governance

| Layer | Control | Process |
|-------|---------|---------|
| Protocol Spec | SDP maintainers | RFC process, major version for breaking changes |
| Reference Implementation | Open source, community PRs | Standard review process |
| Plugin Ecosystem | Community | Registry acceptance via PR, "Verified" tier for audited plugins |

---

## Implementation Roadmap

### Phase 1: Foundation (v0.10) — "Make the Implicit Explicit"

- [ ] Extract `protocol/prompts/` from `sdp-plugin/prompts/`
- [ ] Write `providers/claude/generate.sh`
- [ ] Add `TokensUsed`, `LOCProduced`, `EstimatedCost` to `ExecuteResult`
- [ ] Add `AgentIdentity` to `Proposal` (backward-compatible)
- [ ] Define `AgentBackend` interface, create `LocalBackend` wrapping `ExecuteFunc`
- [ ] Define `CheckpointStore` interface, create `FileCheckpointStore`
- [ ] Extract per-backend circuit breakers from single global breaker
- [ ] Define `plugin.yaml` manifest format
- [ ] Publish SDP Protocol Spec v1.0-draft

### Phase 2: Multi-Provider (v0.11) — "Not Just Claude"

- [ ] Add `HTTPBackend` for first remote agent (Claude API or GLM)
- [ ] Add `GitCheckpointStore` for multi-machine state
- [ ] Implement agent contract YAML format + JSON Schema validation
- [ ] Implement `AgentRegistry` with basic capability matching
- [ ] Add `[autonomy]` section to quality-gate.toml
- [ ] Implement `sdp plugin add/remove` CLI
- [ ] Create GitHub-based plugin registry repo
- [ ] Ship 3 example plugins (terraform, docker, kubernetes)

### Phase 3: Economic Intelligence (v0.12) — "Cost-Aware Routing"

- [ ] Implement `BudgetTracker` integrating with circuit breaker
- [ ] Implement `ProviderCatalog` and `QualityTier` classification
- [ ] Build `EconomicRouter` between GetReady() and Execute()
- [ ] Add trust score tracking per project
- [ ] Implement `AgentPermissions` scoping in dispatcher
- [ ] Add `AuditEntry` logging for all agent actions
- [ ] Declarative synthesis rules (YAML, not just Go)
- [ ] GitHub Action for quality gates in CI

### Phase 4: Autonomous Swarms (v1.0) — "The Human Disappears"

- [ ] Auto-remediation loop (review → fix → re-review)
- [ ] Streaming progress from remote agents
- [ ] Unify `synthesis.Agent` and `AgentBackend` interfaces
- [ ] Reconciler loop for self-healing (Kubernetes operator pattern)
- [ ] Event log alongside snapshots for full audit trail
- [ ] Reproducibility checks for security-sensitive workstreams
- [ ] Per-file autonomy overrides via glob patterns

### Phase 5: Protocol Standard (v2.0) — "The OCI of Agent Orchestration"

- [ ] Extract `sdp-spec` as standalone repository
- [ ] Publish formal protocol specification
- [ ] Reference implementation conformance tests
- [ ] Verified plugins tier
- [ ] Hosted plugin registry
- [ ] Third-party SDP implementations (other languages/frameworks)

---

## Success Metrics

| Metric | Baseline (v0.9) | Target (v1.0) | Target (v2.0) |
|--------|-----------------|---------------|---------------|
| Supported providers | 3 (Claude, Cursor, OpenCode) | 6+ (add GLM, Kimi, QWEN) | 10+ |
| Agent contracts defined | 0 | 19 (all current agents) | 50+ (with plugins) |
| Remote backend support | None | 2+ remote backends | Any HTTP/gRPC agent |
| Cost tracking | None | Per-WS token reporting | Full economic routing |
| Plugin ecosystem | 0 plugins | 3 example + 10 community | 50+ community plugins |
| Autonomy levels supported | 2-3 | 2-5 (configurable) | 0-5 (full spectrum) |
| Trust model | None | Output gates + trust tiers | Full defense-in-depth |
| Protocol spec | Embedded in code | Draft spec published | Formal standard, RFCs |

---

## Appendix: The Convergence Point

All 10 aspects converge on the same first step: **enrich `ExecuteResult`**.

This single struct is the natural integration point:
- The `Dispatcher` already produces it
- The `CheckpointManager` already serializes it
- The `CircuitBreaker` already reacts to it

Adding identity (who), cost (how much), scope (how big), and trust (how reliable) to this struct is the minimal change that unlocks every roadmap above.

```go
type ExecuteResult struct {
    WorkstreamID string              // existing
    Success      bool                // existing
    Error        error               // existing
    Duration     int64               // existing
    AgentID      string              // NEW: who executed
    Provider     string              // NEW: which provider
    LOCProduced  int                 // NEW: output size
    TokensUsed   int64               // NEW: token count
    EstimatedCost float64            // NEW: USD cost
    TrustScore   float64             // NEW: confidence metric
    TaskCount    int                 // NEW: internal decomposition
}
```

Start here. Everything else follows.
