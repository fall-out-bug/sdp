# Reality Pro Spec

## Status

Draft.

## Purpose

`@reality-pro` is a consulting-grade, multi-source, multi-repo analysis workflow built on top of the open `reality` artifact contract. Its goal is to reconstruct how a system actually works, recover likely design intent, expose contradictions between plan and implementation, and prepare a system for high-trust agent execution.

It should behave less like a repo scanner and more like a repeatable technical consulting pipeline.

## Product Boundary

`reality-pro` is private logic over open contracts.

Open:

- artifact formats
- claim model
- readiness statuses
- base schema families

Private:

- orchestration logic
- scoring heuristics
- memory implementation
- connector implementations
- domain-specific analyzers
- arbitration and synthesis logic

## Goals

1. Analyze one repository or a coordinated set of repositories.
2. Ingest code plus optional external knowledge sources.
3. Auto-select the right specialist agents for the detected stack and domain.
4. Require cross-review of findings and review of synthesis.
5. Reconstruct features, intent, architecture, integrations, and risks.
6. Build C4 views and a code map for the analyzed system.
7. Produce agent-readiness and bootstrap plans for future SDP work.

## Non-Goals

- automatic remediation by default
- replacing subject-matter experts when evidence is missing
- asserting business intent without traceable evidence
- treating documentation as truth when code contradicts it

## Inputs

Required:

- one repository or explicit reposet
- git history
- build and test metadata
- local documentation

Optional premium inputs:

- Confluence/wiki/Notion pages
- Jira or other issue trackers
- PR history and review comments
- ADR collections
- runbooks and incident records
- schema registries and external contracts
- ownership and team metadata

## Operating Principles

### Multi-Source Truth

No source is privileged by default. `reality-pro` compares sources and explicitly records disagreement.

### Claim Typing

Every claim must be typed:

- `observed`
- `documented`
- `inferred`
- `conflicted`
- `unknown`

### Confidence Discipline

Every significant finding must include:

- confidence score
- source set
- opposing evidence, if any
- validation status

### Adversarial Review

No important conclusion should survive on a single-agent pass.

## Modes

| Mode | Purpose |
|---|---|
| `--repo <path>` | deep analysis of one repository |
| `--reposet <path1,path2,...>` | coordinated multi-repo analysis |
| `--with-docs` | include external documentation sources |
| `--reconstruct-intent` | emphasize intent recovery and plan-vs-implementation gaps |
| `--bootstrap-sdp` | generate bootstrap backlog and readiness plan |
| `--domain <name>` | bias agent selection for domain-heavy systems |

Current private lab baseline:

- `go run ./cmd/sdp-reality-pro-ingest --repo .`
- `go run ./cmd/sdp-reality-pro-ingest --reposet .,/abs/path/to/other/repo`
- `go run ./cmd/sdp-reality-pro-review --project-root .`
- `go run ./cmd/sdp-reality-pro-report --project-root .`

## Phase Model

### Phase 1: Ingestion

Ingest:

- repositories
- local docs
- optional external docs
- optional issue/PR/ADR history

Normalize all evidence into a common source model.

### Phase 2: Persistent Index and Memory Build

Build and maintain persistent repo memory:

- module summaries
- glossary and domain vocabulary
- feature-to-code mappings
- integration endpoints
- previous validated claims
- unresolved questions
- hotspots and risk zones

This memory must support incremental refresh after later repository changes.

Current private baseline writes:

- `.sdp/reality/repo-memory.json`
- `docs/reality/multi-repo-map.md`

### Phase 3: Specialist Selection

Auto-select specialists from stack and domain signals.

Selection dimensions:

- language and framework
- storage and database types
- platform and infrastructure stack
- security surface
- data flow and integration style
- product/domain signals

Example specialist pool:

- architecture analyst
- runtime analyst
- database analyst
- API analyst
- frontend analyst
- documentation analyst
- test and quality analyst
- security analyst
- domain analyst
- intent analyst

### Phase 4: Parallel Specialist Passes

Run specialists in parallel against the same normalized source graph.

Each specialist must produce:

- scoped findings
- evidence references
- confidence
- open questions
- contradictions with other evidence

### Phase 5: Cross-Review and Dissent

For important findings, require:

1. primary finding
2. opposing review or skeptic pass
3. arbitration decision

Unresolved disagreement becomes a first-class output, not hidden synthesis noise.

Current private lab baseline emits:

- `.sdp/reality/conflicts-report.json`
- `.sdp/reality/intent-gap-report.json`

The current implementation uses deterministic specialist heuristics over `repo-memory.json`. It is not yet the full consulting-grade agent mesh.

### Phase 6: System Reconstruction

Recover:

- feature inventory
- stated vs implemented intent
- system context
- containers
- components
- code-area map
- integration and data-flow map
- constraints and hotspots

Current private lab baseline emits:

- `.sdp/reality/c4-system-context.json`
- `.sdp/reality/c4-container.json`
- `.sdp/reality/c4-component.json`
- `docs/reality/c4-system-context.md`
- `docs/reality/c4-containers.md`
- `docs/reality/c4-components.md`

### Phase 7: Readiness and Bootstrap Synthesis

Produce:

- readiness verdict
- first safe workstreams
- high-risk zones to avoid or isolate
- test-first recommendations
- documentation priorities
- bootstrap backlog for SDP delivery

Current private lab baseline emits:

- `.sdp/reality/bootstrap-backlog.json`
- `.sdp/reality/agent-readiness-plan.json`
- `docs/reality/intent-gap.md`
- `docs/reality/multi-repo-map.md`

These outputs are deterministic syntheses over repo memory plus reviewed findings. They are executable now, but they are not yet the full multi-source consulting pipeline described in this spec.

### Phase 8: Synthesis Review

The final synthesis itself must be reviewed.

The synthesis reviewer checks:

- unsupported claims
- hidden contradictions
- overconfident language
- missing risk qualifiers
- missing next steps

## Required Outputs

All OSS outputs, plus:

- `docs/reality/c4-system-context.md`
- `docs/reality/c4-containers.md`
- `docs/reality/c4-components.md`
- `docs/reality/intent-gap.md`
- `docs/reality/multi-repo-map.md`
- `.sdp/reality/repo-memory.json`
- `.sdp/reality/conflicts-report.json`
- `.sdp/reality/c4-system-context.json`
- `.sdp/reality/c4-container.json`
- `.sdp/reality/c4-component.json`
- `.sdp/reality/intent-gap-report.json`
- `.sdp/reality/bootstrap-backlog.json`
- `.sdp/reality/agent-readiness-plan.json`

## C4 Scope

`reality-pro` should support these layers:

- system context
- container
- component
- code area map

For multi-repo runs, it must also emit a repo landscape view:

- repo roles
- ownership zones
- dependency direction
- contract and schema sharing
- version skew risks

## Readiness Model

`reality-pro` should judge not only code quality, but readiness for autonomous development.

Dimensions:

- boundary clarity
- scope extraction quality
- verification surfaces
- operational safety
- integration stability
- documentation trustworthiness
- memory completeness
- unresolved architectural ambiguity

Verdicts:

- `ready`
- `ready_with_constraints`
- `not_ready`

## Consulting-Grade Expectations

To justify premium value, `reality-pro` must provide:

- evidence-backed claims
- explicit unknowns
- prioritized risks
- preservation advice for fragile areas
- safe-first modernization guidance
- recommended first workstreams

## Exit Criteria

A `reality-pro` run is complete when it has:

1. reconstructed the analyzed system at feature and architecture level
2. produced C4 views and repo landscape views
3. identified the top contradictions between intent, docs, and implementation
4. mapped main integrations and data boundaries
5. produced a reviewed synthesis with confidence and dissent tracking
6. emitted a bootstrap-ready SDP backlog and agent-readiness plan
