# Schema Registry

SDP publishes machine-readable schemas in `schema/`. This document maps each schema family to its purpose.

For implementation patterns and rollout guidance, see `docs/reference/integration-contracts.md`.

## Registry Index

- Canonical index: `schema/index.json`
- Contract: every index entry must point to valid JSON
- Verification: `TestSchemaRegistryLoads` in `sdp-plugin/internal/parser/workstream_test.go`

## Protocol Contracts

| Schema | Purpose |
|---|---|
| `schema/next-action.schema.json` | Legacy machine-readable next-action payload for orchestration flows |
| `schema/contracts/instructions.schema.json` | Machine-readable next-step instruction payload for CLI and agents |
| `schema/contracts/status-view.schema.json` | Machine-readable project state payload for `sdp status --json` |
