# Integration Contracts Guide

> Historical design note. Do not use this file as the source of truth for current CLI behavior or public onboarding. Prefer `sdp <command> --help`, `docs/CLI_REFERENCE.md`, and `docs/PROTOCOL.md`.

Practical guide for teams integrating SDP protocol artifacts into CI, adapters, status surfaces, and review workflows.

## What This Covers

- next-step and status contracts for CLI and automation
- runtime contracts for orchestration and policy decisions
- findings contracts for CI output and local improvement loops
- handoff contracts for analyst/coder/reviewer payloads
- evidence provenance fields for reproducibility and auditability

## Canonical Artifacts

| Family | Schema(s) | Primary Use |
|---|---|---|
| Next-step guidance | `schema/contracts/instructions.schema.json`, `schema/contracts/status-view.schema.json` | Shared payloads for `sdp next --json` and `sdp status --json` |
| Runtime contracts | `schema/contracts/orchestration-event.schema.json`, `schema/contracts/runtime-decision.schema.json` | Event stream and allow/ask/deny decisions across adapters |
| Findings reports | `schema/findings/protocol-findings.schema.json`, `schema/findings/docs-findings.schema.json` | Machine-readable CI findings for sync/automation |
| Findings examples | `schema/findings/examples/protocol-findings-example.json`, `schema/findings/examples/docs-findings-example.json` | Golden payloads for consumers and fixtures |
| Handoff contracts | `schema/handoff-analyst.schema.json`, `schema/handoff-coder.schema.json`, `schema/handoff-reviewer.schema.json` | Typed cross-agent exchange during implementation/review |
| Evidence envelope | `schema/evidence-envelope.schema.json` | Strict run evidence including provenance fields |
| Compatibility payload | `schema/next-action.schema.json` | Legacy next-action payload still consumed by orchestration flows |

Registry source of truth: `schema/index.json`.

## Integration Patterns

### 1) Shared Next-Step Guidance

Use `instructions` when you need one deterministic recommendation with machine-readable routing fields.

- `action_id` gives consumers a stable action handle
- `required_context` and `optional_context` let adapters preload the right inputs
- `policy_expectations` and `evidence_expectations` expose the gates the action is expected to touch

### 2) Shared Project State

Use `status-view` when you need the whole operator-facing state in one payload.

- `workstreams.ready`, `workstreams.blocked`, and `workstreams.in_progress` are already sorted deterministically
- `next_action` and embedded `next_step` should stay aligned so humans and agents see the same guidance
- `environment` and `active_session` make the recommendation explainable without scraping text output

### 3) Runtime Events and Decisions

Use `orchestration-event` to normalize execution telemetry from adapters and orchestrators.

- emit event type + metadata at each critical transition (`task.started`, `quality.gate.failed`, etc.)
- validate payloads against `schema/contracts/orchestration-event.schema.json` before publishing
- keep event names stable; add new names in backward-compatible manner

Use `runtime-decision` when policy or guard logic returns a decision.

- decision surface is explicit: `allow`, `ask`, `deny`
- record reason and context so downstream explainability stays deterministic
- treat unknown decision values as schema violations

### 4) CI Findings -> Local Improvement Loop

CI producers (`sdp-protocol-check`, `sdp-doc-sync`) should emit one findings JSON per check run.

- protocol findings: `schema/findings/protocol-findings.schema.json`
- docs findings: `schema/findings/docs-findings.schema.json`
- deduplicate on `finding_key` at consumer side
- include remediation hints to allow automated patch generation

### 5) Typed Handoffs Across Agent Roles

Use dedicated handoff schemas instead of free-form JSON blobs.

- analyst output -> `schema/handoff-analyst.schema.json`
- coder output -> `schema/handoff-coder.schema.json`
- reviewer output -> `schema/handoff-reviewer.schema.json`

### 6) Evidence Provenance for Reproducibility

`schema/evidence-envelope.schema.json` includes:

- `provenance.prompt_hash`
- `provenance.context_sources`

Use these fields to verify what inputs shaped model output without storing raw prompt text and to support compliance/postmortem workflows.

## Quickstart Snippets

### Validate `status-view` against schema (Python)

```bash
python3 - <<'PY'
import json
from jsonschema import validate

schema = json.load(open("schema/contracts/status-view.schema.json", "r", encoding="utf-8"))
doc = json.load(open("status-view.json", "r", encoding="utf-8"))

validate(instance=doc, schema=schema)
print("status-view payload is valid")
PY
```

### Minimal runtime decision payload

```json
{
  "spec_version": "v1.0",
  "decision_id": "2dfd1087-7b77-4df4-9ec5-6ea6a6d6f4b5",
  "timestamp": "2026-03-06T12:00:00Z",
  "decision_type": "quality.gate",
  "decision": "allow",
  "reason": {
    "code": "QUALITY_GATES_PASSED",
    "message": "all required quality gates passed"
  },
  "context": {
    "request": {
      "action": "merge",
      "resource": "pull_request"
    },
    "workstream_id": "00-077-01",
    "feature_id": "F077",
    "session_id": "run-20260306-120000"
  }
}
```

## Validation Hooks

- registry integrity: `go test ./internal/parser -run TestSchemaRegistryLoads` (in `sdp-plugin`)
- next-step contracts: `go test ./internal/nextstep -run TestStatusAndInstructionSchemasValidateContracts`
- findings examples: tests in `internal/evidenceenv`
- evidence envelope parity: `internal/evidenceenv/schema_test.go`

## Migration Notes

- for new consumers, integrate `instructions` and `status-view` first; treat `next-action` as a compatibility layer
- validate JSON against schema before making automation decisions from `action_id`, `category`, or `next_action`
- if you currently parse untyped handoff or findings JSON, migrate parsers to schema-first validation before business logic
