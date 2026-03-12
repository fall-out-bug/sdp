# SDP CLI Reference

Run `sdp --help` for the full command tree and `sdp <command> --help` for flags and subcommands.

## Core command groups

| Area | Commands | Purpose |
|------|----------|---------|
| Setup and project state | `sdp init`, `sdp doctor`, `sdp health`, `sdp status`, `sdp next`, `sdp demo`, `sdp hooks`, `sdp completion` | Bootstrap SDP, inspect current state, and manage local tooling |
| Planning and execution | `sdp parse`, `sdp plan`, `sdp build`, `sdp apply`, `sdp orchestrate`, `sdp verify`, `sdp tdd`, `sdp deploy` | Parse workstreams, create plans, execute work, verify completion, and record deployment approvals |
| Guard and context | `sdp guard ...`, `sdp session ...`, `sdp resolve`, `sdp git` | Enforce edit scope, validate context/branch state, resolve task identifiers, and keep session state in sync |
| Evidence and audit | `sdp log ...`, `sdp decisions ...`, `sdp checkpoint ...`, `sdp coordination ...`, `sdp design record`, `sdp idea record` | Inspect evidence, trace decision history, manage checkpoints, and record design/idea evidence |
| Quality and diagnostics | `sdp quality {coverage, complexity, size, types, all}`, `sdp drift detect`, `sdp diagnose`, `sdp watch`, `sdp reality {emit-oss, validate}`, `sdp collision check`, `sdp contract ...`, `sdp acceptance run` | Run quality gates, detect drift, emit or validate OSS reality artifacts, inspect failures, watch files, check collisions, validate contracts, and run smoke acceptance checks |
| Telemetry and metrics | `sdp telemetry {status, consent, enable, disable, analyze, export, upload}`, `sdp metrics {collect, classify, report}` | Manage local opt-in telemetry and derive benchmark/quality metrics |
| Workflow support | `sdp beads ...`, `sdp task create`, `sdp memory ...`, `sdp prd ...`, `sdp prototype`, `sdp skill ...` | Integrate with Beads, manage memory/search, work with PRDs, prototype features, and inspect skills |

## Frequently used commands

| Command | Purpose |
|---------|---------|
| `sdp init --auto` | Initialize prompts and SDP scaffolding without prompts |
| `sdp doctor` | Run health checks for hooks, config, telemetry, and repository setup |
| `sdp build <ws-id>` | Execute a single workstream with guard enforcement and tests |
| `sdp reality emit-oss --quick` | Emit the OSS reality baseline into `.sdp/reality/` and `docs/reality/` |
| `sdp reality validate` | Validate emitted OSS reality artifacts against the published schema contract |
| `sdp verify <ws-id>` | Validate workstream completion against evidence and checks |
| `sdp quality all` | Run all quality gates for the current project |
| `sdp telemetry status` | Show telemetry consent, event count, and storage path |
| `sdp telemetry export json` | Export local telemetry to `telemetry_export.json` |
| `sdp log show` | Show paginated evidence events with filters |
| `sdp decisions log` | Record a decision in the audit trail |

See [reference/reality-oss.md](reference/reality-oss.md) for the reality OSS contract and [reference/skills.md](reference/skills.md) for the skill catalog.
