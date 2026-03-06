# Schema Registry

SDP publishes machine-readable schemas in `schema/`. This document maps each schema family to its purpose.

## Registry Index

- Canonical index: `schema/index.json`
- Contract: every index entry must point to valid JSON
- Verification: `TestSchemaRegistryLoads` in `sdp-plugin/internal/parser/workstream_test.go`

## Protocol Contracts

| Schema | Purpose |
|---|---|
| `schema/contracts/orchestration-event.schema.json` | Event contract for orchestration telemetry (`task.started`, `quality.gate.failed`, etc.) |
| `schema/contracts/runtime-decision.schema.json` | Runtime governance decision contract with `allow` / `ask` / `deny` semantics |

## Findings Reports

| Schema | Purpose |
|---|---|
| `schema/findings/protocol-findings.schema.json` | Structured output for protocol findings from CI checks |
| `schema/findings/docs-findings.schema.json` | Structured output for documentation findings from CI checks |

## Agent Handoff Contracts

| Schema | Purpose |
|---|---|
| `schema/handoff-analyst.schema.json` | Analyst -> coder handoff payload |
| `schema/handoff-coder.schema.json` | Coder -> reviewer handoff payload |
| `schema/handoff-reviewer.schema.json` | Reviewer verdict payload (`approve` / `needs_changes` / `reject`) |

## Hook Template

- `.sdp/pipeline-hooks.yaml.example` provides the default shape for project-level pipeline hooks.
- Security model and fail behavior are documented in `docs/reference/pipeline-hooks-security.md`.
