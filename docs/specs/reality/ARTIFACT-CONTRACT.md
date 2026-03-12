# Reality Artifact Contract

## Status

Draft.

## Goal

This document defines the open artifact contract shared by `reality` and `reality-pro`.

The contract must remain open even if premium logic stays private. Private value should live in analysis depth and orchestration quality, not in opaque output formats.

## Design Rules

1. Artifacts must be machine-readable first and human-readable second.
2. Claims must distinguish observation from inference.
3. Premium outputs may extend the contract, but not break base artifacts.
4. OSS and Pro must share IDs, enums, and evidence semantics.

## Suggested Repository Layout

Human-readable outputs:

- `docs/reality/summary.md`
- `docs/reality/architecture.md`
- `docs/reality/quality.md`
- `docs/reality/bootstrap.md`
- `docs/reality/c4-*.md` for Pro

Machine-readable outputs:

- `.sdp/reality/*.json`

Schema definitions, once formalized:

- `schema/reality/*.schema.json`

Product specs:

- `docs/specs/reality/OSS-SPEC.md`
- `docs/specs/reality/PRO-SPEC.md`
- `docs/specs/reality/ARTIFACT-CONTRACT.md`

## Core Types

### Claim

Every important conclusion is a `claim`.

Required fields:

- `claim_id`
- `title`
- `statement`
- `status`
- `confidence`
- `source_ids`
- `review_state`

Recommended fields:

- `affected_repos`
- `affected_paths`
- `affected_components`
- `counter_evidence`
- `open_questions`
- `tags`

### Claim Status Enum

- `observed`: directly supported by code, config, test, or executable runtime evidence
- `documented`: stated in docs but not yet confirmed in implementation
- `inferred`: reasoned from indirect evidence
- `conflicted`: evidence sources disagree
- `unknown`: insufficient support

### Review State Enum

- `unreviewed`
- `cross_checked`
- `challenged`
- `arbitrated`

### Source

Every source referenced by claims must be normalized.

Required fields:

- `source_id`
- `kind`
- `locator`
- `revision`

Recommended fields:

- `repo`
- `path`
- `uri`
- `line_range`
- `excerpt_hash`
- `captured_at`

### Source Kind Enum

- `code`
- `test`
- `config`
- `manifest`
- `doc`
- `issue`
- `pull_request`
- `commit`
- `runtime_trace`
- `external_contract`

## Artifact Families

### 1. Summary Family

Purpose: top-level run summary.

Files:

- `reality-summary.json`
- `summary.md`

Must include:

- run scope
- analyzed repos
- top findings
- readiness verdict
- artifact inventory

### 2. Feature Inventory Family

Purpose: reconstructed features and implementation status.

Files:

- `feature-inventory.json`

Per feature entry:

- synthetic or canonical feature id
- title
- summary
- status: `implemented|partial|candidate|dead`
- evidence claims
- confidence
- mapped components

### 3. Architecture Family

Purpose: reconstructed structure of the system.

Files:

- `architecture-map.json`
- `c4-system-context.json` (Pro)
- `c4-container.json` (Pro)
- `c4-component.json` (Pro)

Must support:

- nodes
- edges
- boundaries
- external integrations
- data stores
- hotspots

Pro C4 payload expectations:

- `c4-system-context.json`: system scope, internal/external systems, relationships
- `c4-container.json`: container inventory, responsibilities, inter-container relationships
- `c4-component.json`: component inventory, code-path mapping, relationships inside a container

### 4. Integration Family

Purpose: external and internal dependency mapping.

Files:

- `integration-map.json`

Must include:

- integration type
- producer and consumer
- contract type if known
- confidence
- failure and risk notes if known

### 5. Quality Family

Purpose: code-quality and test-quality assessment.

Files:

- `quality-report.json`

Must include:

- maintainability findings
- test posture
- hotspot ranking
- major code smells
- structural risks

### 6. Drift and Intent Family

Purpose: document contradictions and gaps.

Files:

- `drift-report.json`
- `intent-gap-report.json` (Pro)
- `conflicts-report.json` (Pro)

Must include:

- contradictory claims
- stale docs or misleading docs
- plan-vs-implementation mismatches
- unresolved questions

Pro payload expectations:

- `intent-gap-report.json`: expected state, observed state, gap type, severity, status
- `conflicts-report.json`: competing claim IDs, conflict severity, arbitration status, resolution notes if any

### 7. Readiness Family

Purpose: answer whether the system is safe for future agent work.

Files:

- `readiness-report.json`
- `agent-readiness-plan.json` (Pro)
- `bootstrap-backlog.json` (Pro)

Required readiness enum:

- `ready`
- `ready_with_constraints`
- `not_ready`

Required readiness dimensions:

- boundary clarity
- verification coverage
- hotspot concentration
- integration fragility
- documentation trust level

Pro payload expectations:

- `agent-readiness-plan.json`: current verdict, target verdict, phased readiness plan, allowed scope, blocked zones, exit criteria
- `bootstrap-backlog.json`: proposed workstreams, evidence-backed rationale, dependencies, exit criteria

### 8. Persistent Memory Family

Purpose: carry normalized memory across later reality-pro refreshes.

Files:

- `repo-memory.json` (Pro)

Must include:

- analyzed repo set
- module summaries
- feature-to-code mappings
- unresolved questions
- reusable hotspot context

## OSS vs Pro Contract Rules

OSS may emit:

- summary family
- feature inventory family
- architecture family without full C4 depth
- integration family
- quality family
- drift family without premium intent recovery
- readiness family baseline

Pro may extend with:

- persistent memory artifacts
- full C4 artifacts
- intent-gap family
- conflict arbitration outputs
- bootstrap backlog and agent-readiness plan

Pro must not rename or change base enums used by OSS.

## Validation Rules

1. Every finding of consequence must be linked to at least one source.
2. `observed` claims must reference executable evidence classes: code, config, manifest, test, or runtime trace.
3. `documented` claims may not be promoted to `observed` without stronger evidence.
4. `conflicted` claims must preserve both sides of the contradiction.
5. Final readiness verdict must cite the claims that justify it.

Validation command for OSS artifacts:

```bash
go run ./cmd/sdp-reality-validate .
```

## Schema Split

Schema files are split by artifact family:

- `schema/reality/claim.schema.json`
- `schema/reality/source.schema.json`
- `schema/reality/reality-summary.schema.json`
- `schema/reality/feature-inventory.schema.json`
- `schema/reality/architecture-map.schema.json`
- `schema/reality/integration-map.schema.json`
- `schema/reality/quality-report.schema.json`
- `schema/reality/drift-report.schema.json`
- `schema/reality/readiness-report.schema.json`
- `schema/reality/repo-memory.schema.json` (Pro)
- `schema/reality/conflicts-report.schema.json` (Pro)
- `schema/reality/intent-gap-report.schema.json` (Pro)
- `schema/reality/bootstrap-backlog.schema.json` (Pro)
- `schema/reality/agent-readiness-plan.schema.json` (Pro)
- `schema/reality/c4-*.schema.json` (Pro)
