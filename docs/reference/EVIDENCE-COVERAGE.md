# Evidence Coverage Matrix (F056)

Skill → event types emitted. Use for pipeline verification and `sdp log show --type=X`.

---

## Evidence Tracking Policy (WS-067-07: AC4)

### Tracked in Git (Source of Truth)

| Path | Purpose | Why Tracked |
|------|---------|-------------|
| `.sdp/log/events.jsonl` | Evidence event log | Audit trail for AI decisions |
| `.sdp/config.yml` | Project configuration | Required for reproducibility |
| `.sdp/guard-rules.yml` | Quality gate rules | Required for CI/local consistency |
| `.beads/issues.jsonl` | Task tracking | Session persistence |

### Not Tracked (Generated/Runtime)

| Pattern | Why Not Tracked |
|---------|-----------------|
| `*.out` | Build/test coverage artifacts |
| `.sdp/memory.db` | SQLite index (rebuildable) |
| `.sdp/checkpoints/` | Runtime state (resumable) |
| `coverage.html` | Generated report |
| `bin/`, `dist/` | Compiled binaries |

### Merge Strategy (.gitattributes)

```gitattributes
# Union merge for concurrent appends
.sdp/log/events.jsonl merge=union

# Beads custom merge
.beads/issues.jsonl merge=beads
```

---

## Pipeline chain (idea → deploy)

| Phase   | Skill     | Event type(s)     | CLI / trigger                          |
|---------|-----------|-------------------|----------------------------------------|
| Idea    | @idea     | plan              | `sdp idea record`, `sdp parse` (per WS) |
| Design  | @design   | plan              | `sdp design record`, `sdp parse`       |
| Build   | @build    | generation        | TDD runner (F054)                       |
| Review  | @review   | verification      | `sdp verify` (per gate)                 |
| Deploy  | @deploy   | approval          | `sdp deploy --target main`              |

## Skill × event types

| Skill      | plan | generation | verification | approval | Notes                          |
|------------|------|------------|--------------|----------|--------------------------------|
| @vision    | ✓    |            |              |          | `sdp skill record --skill vision --type plan` |
| @reality   |      |            | ✓            |          | `sdp skill record --skill reality --type verification` |
| @idea      | ✓    |            |              |          | `sdp idea record`             |
| @design    | ✓    |            |              |          | `sdp design record`, `sdp parse` |
| @build     |      | ✓          |              |          | Evidence layer (F054)          |
| @review    |      |            | ✓            |          | `sdp verify` (per gate)       |
| @deploy    |      |            |              | ✓        | `sdp deploy`                  |
| @oneshot   | ✓    |            |              | ✓        | `sdp orchestrate` (plan start, approval end) |
| @prototype |      | ✓          |              |          | `sdp prototype` (after WS gen) |
| @hotfix    |      | ✓          |              | ✓        | `sdp skill record` (2 calls)  |
| @bugfix    |      | ✓          | ✓            |          | `sdp skill record` (gen + verification) |
| @issue     | ✓    |            |              |          | `sdp skill record --skill issue --type plan` |
| @debug     |      |            | ✓            |          | `sdp skill record --skill debug --type verification` |

## Commands

- **Show recent events:** `sdp log show` (last 20)
- **Filter by type:** `sdp log show --type=plan` (or `generation`, `verification`, `approval`, `decision`, `lesson`)
- **Trace by commit/WS:** `sdp log trace [commit-sha]` or `sdp log trace --ws 00-056-01`
- **Chain integrity:** `sdp log trace --verify`

## Schema

Event types: `plan`, `generation`, `verification`, `approval`, `decision`, `lesson`.  
Schema: [schema/evidence.schema.json](../../schema/evidence.schema.json).
