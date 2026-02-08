# F057: CLI plan/apply/log

> Beads: sdp-5lca | Priority: P1

---

## Problem

Currently SDP is a Claude Code plugin with 19 skills. No standalone CLI UX. Users can't:
- Plan features from the terminal
- Execute builds without Claude Code
- Browse evidence without reading JSONL manually

The Terraform-inspired `plan/apply/log` UX makes SDP accessible outside Claude Code.

## Solution

Three top-level commands:

### `sdp plan <description>`
```bash
sdp plan "Add OAuth2 login"
# → Decomposition output (workstreams, dependencies, estimates)
# → Emits plan event to evidence log
# → Creates workstream files in docs/workstreams/backlog/

sdp plan "Add OAuth2 login" --interactive  # Drive mode: asks questions
sdp plan "Add OAuth2 login" --auto-apply   # Ship mode: plan + execute
```

Internally calls `@idea` + `@design` flow.

### `sdp apply`
```bash
sdp apply                    # Execute all ready workstreams
sdp apply --ws 00-054-01     # Execute specific workstream
sdp apply --retry 3          # Retry failed unit up to 3 times
sdp apply --dry-run          # Show what would execute

# Streaming progress:
# [00-054-01] ████████░░ 80% — running tests
# [00-054-02] waiting (blocked by 00-054-01)
```

Internally calls `@build` for each workstream.

### `sdp log`
```bash
sdp log                      # Recent 20 events
sdp log trace abc1234        # Chain backwards from commit
sdp log trace --ws 00-054-03 # Chain for workstream
sdp log show                 # Interactive browser
sdp log show --type=generation --model=claude-sonnet-4
sdp log export --format=csv  # For analysis
```

Already partially built in F054 (WS-07). This feature adds `show`, filtering, export.

## Constraints

- All commands emit `--output=json` for scripting
- All commands work offline (evidence is local)
- `sdp plan` without model API = decomposition only (no AI generation)
- `sdp apply` without model API = error with clear message
- Streaming progress for long-running operations

## Users

- Developers who prefer terminal over IDE
- CI/CD pipelines
- Enterprise teams standardizing on CLI workflows

## Success Metrics

- Full feature lifecycle: `sdp plan` → `sdp apply` → `sdp log trace` works E2E
- JSON output parseable by jq for all commands
- `sdp apply` progress visible in real-time

## Dependencies

- F054 (evidence log, acceptance gate, `sdp log trace`)
- F056 (full instrumentation — for complete evidence chain in `sdp apply`)

## Notes

- `sdp plan` is the new entry point — replaces `@feature` for non-IDE users
- `sdp apply` is `@oneshot` but standalone
- This is the path to CI/CD integration (F058) — actions call CLI commands
