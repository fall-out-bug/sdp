# ADR-0002: Adopt File-Native, Schema-Validated Protocol as the Canonical Core (Reliability-First)

## Status

Proposed

## Context

This repository targets workflows from small student projects to large multi-week epics. The current multi-agent protocol (v1.2) is file-based and educational, but it has reliability risks that become costly as projects scale:

- **Invalid / drifting structure**: messages and artifacts are described in prose or loosely, so agents sometimes produce malformed or inconsistent outputs.
- **Implicit state**: determining “what phase are we in?” requires reading many files, increasing context cost and errors.
- **Execution gap**: plans can be too high-level to execute deterministically; resumability is weak.
- **Inconsistency across docs/examples/prompts**: differing artifact formats and paths reduce beginner success and automation reliability.
- **Tool lock-in risk**: IDE-native “special features” can improve DX, but must not become mandatory for protocol correctness.

Two v2 drafts exist with different philosophies:

- **Cursor-native**: TypeScript interfaces + IDE “consoles” (strong DX, weaker portability/runtime validation).
- **File-native**: JSON Schema + runtime validation + portable prompts (strong reliability/portability).

Project direction preference: **Automation / LLM reliability-first**.

## Decision

We will make **File-Native + JSON Schema validation** the canonical core of the protocol (“Schema is Law”), and treat IDE-native experiences (Cursor consoles, TypeScript typings) as **optional adapters**.

Specifically:

1. **Canonical contract**: JSON Schema defines all message and artifact structures (runtime-validatable).
2. **Single source of truth for state**: each epic has `consensus/status.json` (validated) as the traffic-light state machine.
3. **Canonical artifacts are machine-readable**: JSON is the primary artifact format used by agents and validators; human-friendly Markdown summaries are optional.
4. **Validation is a required gate**: agents (and humans) must validate before advancing phases.
5. **Adapters are optional**: Cursor/TypeScript convenience layers may exist, but they MUST NOT be required for correctness.

## Scope

This ADR defines the canonical v2 protocol core:

- directory layout (per-epic)
- schemas for status/messages/artifacts
- validation and phase gates
- artifact format policy (JSON-first)
- optional adapter strategy

This ADR does NOT redesign Solo/Structured modes; it only defines the canonical core for Multi-Agent (and optionally Structured) when reliability is needed.

## Non-Goals

- Mandating any specific AI tool (Cursor/Claude Code/Codex/etc.).
- Requiring compilation/tooling (TypeScript build) as a prerequisite for participation.
- Building a full orchestration engine; the protocol remains file-based and human-in-the-loop.

## Detailed Design

### 1) Canonical Directory Structure (Per Epic)

```
docs/specs/{epic}/
├── epic.md                      # Human-readable problem statement (language flexible)
├── architecture.md              # Human-readable summary (optional but recommended)
├── implementation.md            # Human-readable plan summary (optional but recommended)
├── testing.md                   # Human-readable test strategy (optional but recommended)
├── deployment.md                # Human-readable deployment plan (optional but recommended)
└── consensus/
    ├── status.json              # REQUIRED, validated against status.schema.json
    ├── schema/                  # REQUIRED (or symlink to shared schemas)
    │   ├── status.schema.json
    │   ├── message.schema.json
    │   ├── requirements.schema.json
    │   ├── architecture.schema.json
    │   ├── plan.schema.json
    │   ├── implementation.schema.json
    │   └── test_results.schema.json
    ├── artifacts/               # REQUIRED outputs (JSON-first)
    │   ├── requirements.json
    │   ├── architecture.json
    │   ├── plan.json            # workstreams/tasks; replaces prose-only planning as canonical executable plan
    │   ├── implementation.json  # per-task progress + evidence pointers
    │   └── test_results.json
    ├── messages/
    │   ├── inbox/{role}/        # REQUIRED, JSON messages validated against message.schema.json
    │   └── processed/{role}/    # OPTIONAL but recommended lifecycle hygiene
    └── decision_log/            # OPTIONAL but recommended; Markdown audit notes
```

### 2) Canonical “JSON-first” Artifact Policy

- **Canonical**: `requirements.json`, `architecture.json`, `plan.json`, `implementation.json`, `test_results.json`.
- **Optional human summaries**: `requirements.md`, `architecture.md`, etc. may exist, but are not required for automation.
- **Rule**: automation and phase gates depend only on JSON artifacts + status.json.

### 3) Centralized State Machine (`status.json`)

`status.json` is the single, validated state indicator. It includes:

- `epic_id`, `phase`, `mode`, `iteration`
- `approvals`, `vetoes`, `blockers`
- `workstreams` (or references to `plan.json`)
- `updated_at`, `updated_by`

State transitions are explicit, validated, and phase-owned:

- only the current phase owner role can advance phase (except moving to `blocked` with reason).

### 4) Messages (Inbox) are Validated and Minimal

All inter-agent messages:

- are JSON with compact keys (token efficiency)
- validated against `message.schema.json`
- include artifact references when handing off
- do NOT require reading other roles’ inboxes (existing rule preserved)

### 5) Validation Gate (Required)

A platform-agnostic validator (script or documented command) must:

- validate JSON syntax
- validate JSON Schema conformance
- verify required artifacts exist for the current phase
- optionally enforce “English-only for machine fields” (see Open Questions)

Agents MUST:

- read `status.json` before acting
- validate outputs before writing/advancing
- update `status.json` after completing deliverables

### 6) Conflict / Concurrency Strategy

Because this is file-based:

- we adopt a simple rule: **single-writer by phase**
- optional improvement: a lightweight lock file (`status.lock`) with TTL; if present and fresh, others stop and report

### 7) Optional Adapters (DX Layer)

Adapters MAY provide:

- TypeScript type generation from JSON Schema
- Cursor “agent consoles” as curated context anchors
- UI helpers (kanban views, clickable actions)

Adapters MUST:

- be additive only
- never change canonical schema meanings
- never be required to validate or participate

## Alternatives Considered

### 1) TypeScript interfaces as the primary contract

Rejected as canonical because it is compile-time oriented and not universally runtime-validatable without extra tooling.

Kept as an optional generated artifact from JSON Schema.

### 2) No centralized status file (implicit state)

Rejected due to context overhead and ambiguity at scale.

### 3) External orchestration engines (LangGraph/CrewAI/etc.)

Rejected because they reduce human-in-the-loop control and add infrastructure complexity.

## Consequences

### Positive

- **Higher reliability**: schema validation catches malformed outputs early.
- **Lower context cost**: `status.json` provides immediate phase visibility.
- **Resumability**: `plan.json` + `implementation.json` provide deterministic task tracking.
- **Tool-agnostic**: any system that can read/write files can participate.
- **Beginner success (reliability-first)**: fewer “it broke because format drifted” failures.

### Negative

- **More upfront structure**: schemas/status require setup and discipline.
- **Schema evolution**: needs versioning strategy.
- **Strictness can block progress**: must provide clear errors and quick fixes.

## Rollout Plan (Proposed)

- Phase A: introduce schemas + validator + status.json as “v2 experimental” alongside v1.2.
- Phase B: migrate multi-agent examples to JSON-first canonical artifacts.
- Phase C: mark v1.2 implicit-state workflow as “legacy”; keep for education/quick usage only.
- Phase D: optional Cursor adapter docs and generated TypeScript typings.

## Migration Notes

A migration tool (optional) can:

- infer `status.json` from existing artifacts/messages
- convert markdown artifacts to JSON canonical equivalents (where possible)
- create placeholder `plan.json` from `implementation.md` workstreams

## Open Questions

1. **English-only scope**: should we require English only in machine artifacts/messages, while allowing `epic.md` to be any language?
2. **Schema versioning**: embed `schema_version` in status/artifacts vs directory versioning?
3. **Task granularity**: what is the recommended maximum size of a “task” for `plan.json`?
4. **Audit trail**: do we require `status_history.jsonl` to record all state transitions?
5. **Role naming**: standardize on `qa` vs `quality` as canonical role key (affects schema/enforcement).

## References

- Builds on v1.2 file-based protocol and quality gates.
- Aligns with the “File-Native v2.0” direction; Cursor-native features become optional adapters.


