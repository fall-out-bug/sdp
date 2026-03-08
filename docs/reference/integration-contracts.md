# Integration Contracts Guide

Practical guide for teams integrating SDP's control-tower contracts into CLI tooling, adapters, and automation.

## Canonical Artifacts

| Schema | Purpose |
|---|---|
| `schema/contracts/instructions.schema.json` | Machine-readable next-step payload shared by `sdp next --json` and `sdp status --json` |
| `schema/contracts/status-view.schema.json` | Machine-readable project state payload returned by `sdp status --json` |
| `schema/next-action.schema.json` | Legacy next-action payload used by orchestration flows |

Registry source of truth: `schema/index.json`.

## Integration Patterns

### 1) Shared Next-Step Guidance

Use `instructions` when you need one deterministic recommendation with machine-readable routing fields.

- `action_id` gives consumers a stable action handle
- `required_context` and `optional_context` let adapters preload the right inputs
- `policy_expectations` and `evidence_expectations` expose what gates the recommended action is expected to touch

### 2) Shared Project State

Use `status-view` when you need the whole operator-facing state in one payload.

- `workstreams.ready`, `workstreams.blocked`, and `workstreams.in_progress` are already sorted deterministically
- `next_action` and embedded `next_step` must stay aligned so humans and agents see the same guidance
- `environment` and `active_session` make the recommendation explainable without scraping text output

### 3) Compatibility Surface

`next-action` remains in the registry because orchestration flows still use it.

- prefer `instructions` and `status-view` for new CLI-facing integrations
- keep `next-action` consumers stable until they are explicitly migrated

## Quickstart Snippet

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

## Validation Hooks

- Registry integrity: `go test ./internal/parser -run TestSchemaRegistryLoads` (in `sdp-plugin`)
- Contract payloads: `go test ./internal/nextstep -run TestStatusAndInstructionSchemasValidateContracts`

## Migration Notes

- For new consumers, integrate `instructions` and `status-view` first; treat `next-action` as a compatibility layer.
- Validate JSON against schema before making automation decisions from `action_id`, `category`, or `next_action`.
