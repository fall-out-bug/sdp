# Integration Contracts Guide

Practical guide for teams integrating SDP protocol artifacts into CI, adapters, and review workflows.

## What This Covers

- Runtime contracts for orchestration and policy decisions
- Findings contracts for CI output and local improvement loops
- Handoff contracts for analyst/coder/reviewer payloads
- Evidence provenance fields (`prompt_hash`, `context_sources`) for reproducibility

## Canonical Artifacts

| Family | Schema(s) | Primary Use |
|---|---|---|
| Runtime contracts | `schema/contracts/orchestration-event.schema.json`, `schema/contracts/runtime-decision.schema.json` | Event stream and allow/ask/deny decisions across adapters |
| Findings reports | `schema/findings/protocol-findings.schema.json`, `schema/findings/docs-findings.schema.json` | Machine-readable CI findings for sync/automation |
| Findings examples | `schema/findings/examples/protocol-findings-example.json`, `schema/findings/examples/docs-findings-example.json` | Golden payloads for consumers and fixtures |
| Handoff contracts | `schema/handoff-analyst.schema.json`, `schema/handoff-coder.schema.json`, `schema/handoff-reviewer.schema.json` | Typed cross-agent exchange during implementation/review |
| Evidence envelope | `schema/evidence-envelope.schema.json` | Strict run evidence including prompt provenance |

Registry source of truth: `schema/index.json`.

## Integration Patterns

### 1) Runtime Event Ingestion

Use `orchestration-event` to normalize execution telemetry from adapters and orchestrators.

- Emit event type + metadata at every critical transition (`task.started`, `quality.gate.failed`, etc.)
- Validate payloads against `schema/contracts/orchestration-event.schema.json` before publishing
- Keep event names stable; add new names in backward-compatible manner

Use `runtime-decision` when policy or guard logic returns a decision.

- Decision surface is explicit: `allow`, `ask`, `deny`
- Record reason and context so downstream explainability remains deterministic
- Treat unknown decision values as schema violations

### 2) CI Findings -> Local Improvement Loop

CI producers (`sdp-protocol-check`, `sdp-doc-sync`) should emit one findings JSON per check run.

- Protocol findings: `schema/findings/protocol-findings.schema.json`
- Docs findings: `schema/findings/docs-findings.schema.json`
- Deduplicate on `finding_key` at consumer side
- Include remediation hints to allow automated patch generation

Use these examples as compatibility fixtures:

- `schema/findings/examples/protocol-findings-example.json`
- `schema/findings/examples/docs-findings-example.json`

### 3) Typed Handoffs Across Agent Roles

Use dedicated handoff schemas instead of free-form JSON blobs.

- Analyst output -> `schema/handoff-analyst.schema.json`
- Coder output -> `schema/handoff-coder.schema.json`
- Reviewer output -> `schema/handoff-reviewer.schema.json`

Benefits:

- deterministic parser behavior
- simpler contract tests
- less coupling between prompt wording and integration code

### 4) Evidence Provenance for Reproducibility

`schema/evidence-envelope.schema.json` includes:

- `provenance.prompt_hash`: hash of the rendered prompt
- `provenance.context_sources`: typed list of context inputs with digest

Use these fields to:

- verify what inputs shaped model output without storing raw prompt text
- correlate behavior changes with prompt/context drift
- support compliance and incident postmortems

## Quickstart Snippets

### Validate a findings payload against schema (Python)

```bash
python3 - <<'PY'
import json
from jsonschema import validate

schema = json.load(open("schema/findings/protocol-findings.schema.json", "r", encoding="utf-8"))
doc = json.load(open("schema/findings/examples/protocol-findings-example.json", "r", encoding="utf-8"))

validate(instance=doc, schema=schema)
print("protocol findings payload is valid")
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

Validate against `schema/contracts/runtime-decision.schema.json` before publish.

### Minimal provenance extension inside evidence envelope

```json
{
  "provenance": {
    "prompt_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "context_sources": [
      {
        "type": "workstream_spec",
        "path": "docs/workstreams/backlog/00-077-01.md",
        "hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
      }
    ]
  }
}
```

Keep these fields as hashes/metadata only; do not store raw prompts in evidence.

## Producer/Consumer Checklist

Producer side:

- Emit only fields defined by target schema
- Include stable IDs (`findings_id`, `finding_key`, run identifiers)
- Fail CI step if payload no longer validates

Consumer side:

- Pin to known schema family + version path
- Reject unknown enum values for policy-critical fields
- Log schema validation failures with payload source metadata

## Validation Hooks

- Registry integrity: `go test ./internal/parser -run TestSchemaRegistryLoads` (in `sdp-plugin`)
- Findings examples stay valid via tests in `internal/evidenceenv/findings_examples_test.go`
- Evidence envelope parity is guarded by `internal/evidenceenv/schema_test.go`

## Migration Notes

- If you currently parse untyped handoff or findings JSON, migrate parsers to schema-first validation before business logic.
- If you already store evidence envelopes, ensure your parser accepts `prompt_hash` and `context_sources` in provenance.
- Keep custom extensions outside canonical objects, or namespace them explicitly to avoid future collisions.
