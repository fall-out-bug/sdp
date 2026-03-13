# SDP Plugin Library Capabilities

This reference maps the main Go packages in `sdp-plugin` to the user-visible functionality they implement.

## CLI surface

| Package | Capability |
|---------|------------|
| `cmd/sdp` | CLI entrypoint and command tree for planning, execution, evidence, diagnostics, telemetry, metrics, guard, sessions, Beads integration, and utility flows |

## Planning and execution

| Package | Capability |
|---------|------------|
| `internal/parser` | Parse workstream files and extract structured metadata |
| `internal/planner` | Generate planning artifacts and design-oriented evidence |
| `internal/nextstep` | Recommend next actions from repo and workflow state |
| `internal/orchestrator` | Execute multi-workstream feature flows, resume from checkpoints, and manage role/router state |
| `internal/checkpoint` | Persist and restore long-running execution checkpoints |
| `internal/session` | Track per-worktree CLI session state |
| `internal/worktree` | Create and inspect feature worktrees and related session metadata |

## Evidence and audit

| Package | Capability |
|---------|------------|
| `internal/evidence` | Build evidence events, write/read the hash-chained event log, export events, filter/search/traverse evidence, and page browser views |
| `internal/decision` | Decision audit logging and reporting helpers |
| `internal/coordination` | Coordination event storage, stats, and integrity verification |

## Quality, diagnostics, and safety

| Package | Capability |
|---------|------------|
| `internal/quality` | Coverage, complexity, file-size, and type checks for Python, Go, and Java projects |
| `internal/watcher` | Debounced file watching, include/exclude filtering, and live quality violation tracking |
| `internal/doctor` | Environment, config, hook, dependency, and deep diagnostic checks |
| `internal/diagnostics` | Diagnostic report generation for failures and environment analysis |
| `internal/drift` | Detect documentation-to-code drift and workstream inconsistencies |
| `internal/collision` | Detect overlapping scopes and shared-boundary risks across workstreams/features |
| `internal/guard` | Scope enforcement, context recovery, branch validation, and review findings management |
| `internal/security` | Safe command execution and path/process safety helpers |

## Telemetry and metrics

| Package | Capability |
|---------|------------|
| `internal/telemetry` | Local opt-in command telemetry, consent persistence, analysis, export, and upload packaging |
| `internal/metrics` | Derive benchmark and operational metrics from evidence and telemetry data |

## Supporting subsystems

| Package | Capability |
|---------|------------|
| `internal/sdpinit` | Initialize prompts/config layout, interactive bootstrap, and project type detection |
| `internal/memory` | Long-term memory indexing, search, and stats |
| `internal/task` | Task creation helpers and workstream-linked task operations |
| `internal/beads` | Beads integration for ready/create/show/update/close/sync flows with repo snapshot persistence |
| `internal/verify` | Verify workstream completion against protocol expectations |
| `internal/errors` | Structured error taxonomy and formatted SDP errors |
