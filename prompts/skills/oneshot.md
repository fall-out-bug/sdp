# oneshot

> **CLI:** `sdp orchestrate <feature-id>`
> **Aliases:** `@oneshot`, `/oneshot`

Autonomous feature execution with review-fix loop and PR creation.

---

## Quick Start (CLI)

```bash
sdp orchestrate F067              # Execute all workstreams for F067
sdp orchestrate resume F067       # Resume from checkpoint
sdp orchestrate --retry 3 F067    # Allow 3 retries per workstream
```

---

## What It Does

Executes all workstreams of a feature autonomously in 4 phases:

```
Phase 1: Execute Workstreams
    └─ Load WS from docs/workstreams/backlog/00-XXX-*.md
    └─ Build dependency graph
    └─ Execute in topological order
    └─ Create checkpoint after each WS

Phase 2: Review-Fix Loop (max 3 iterations)
    └─ Run @review
    └─ If APPROVED → proceed
    └─ If CHANGES_REQUESTED → fix findings
    └─ Repeat until approved or max iterations

Phase 3: Verify Clean State
    └─ sdp guard finding list
    └─ Must show 0 blocking findings

Phase 4: Create PR
    └─ Push to feature branch
    └─ Create PR to dev (NOT main)
```

---

## AI Behavior (when invoked as @oneshot)

When user invokes `@oneshot F067`, you MUST execute these steps:

### Step 1: Initialize

Read the workstream files:
```bash
ls docs/workstreams/backlog/00-067-*.md
```

Load checkpoint if exists:
```bash
cat .sdp/checkpoints/F067-*.json
```

### Step 2: Execute CLI Command

The primary execution is via CLI:
```bash
sdp orchestrate F067
```

This command:
- Loads workstreams from `docs/workstreams/backlog/`
- Builds dependency graph
- Executes workstreams in topological order
- Creates checkpoints after each workstream

### Step 3: Handle Results

**On success:**
```
✅ Feature F067 completed successfully
```

**On failure:**
- Check `.sdp/checkpoints/F067-*.json` for state
- Use `sdp orchestrate resume F067` to continue
- Fix blocking issues before resuming

---

## Review-Fix Logic

### Finding Priority Handling

| Priority | Name | Action | Blocks PR? |
|----------|------|--------|------------|
| P0 | Critical | Fix immediately | YES |
| P1 | High | Create bugfix | YES |
| P2 | Medium | Track only | NO |
| P3 | Low | Track only | NO |

### Fix Commands

```bash
# P0: Fix directly
# Make code change, commit

# P1: Create bugfix
bd create --title="Fix: description" --priority=1

# P2+: Track only
bd create --title="Track: description" --priority=2
```

---

## CLI Reference

### orchestrate

```bash
sdp orchestrate <feature-id>           # Execute feature
sdp orchestrate resume <feature-id>    # Resume from checkpoint
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `--checkpoint-dir` | `.sdp/checkpoints` | Checkpoint directory |
| `--workstream-dir` | `docs/workstreams/backlog` | Workstream directory |
| `--retry` | 2 | Max retries per workstream |

### Checkpoint Format

Location: `.sdp/checkpoints/F067-<timestamp>.json`

```json
{
  "feature": "F067",
  "status": "in_progress",
  "completed_ws": ["00-067-01", "00-067-02"],
  "failed_ws": [],
  "execution_order": ["00-067-01", ..., "00-067-14"],
  "started_at": "2026-02-13T14:30:00Z",
  "last_updated": "2026-02-13T15:45:00Z"
}
```

---

## Output Examples

### Success

```
🚀 Orchestrating feature F067
   Workstream dir: docs/workstreams/backlog
   Checkpoint dir: .sdp/checkpoints
   Max retries: 2

→ [14:30] Executing 00-067-01...
→ [14:35] ✅ 00-067-01 complete
→ [14:35] Executing 00-067-02...
...
✅ Feature F067 completed successfully
```

### Resume

```
🔄 Resuming from checkpoint F067
   Workstream dir: docs/workstreams/backlog
   Checkpoint dir: .sdp/checkpoints

→ [15:00] Resuming at 00-067-05...
...
✅ Checkpoint F067 completed successfully
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Command not found | Run `go install ./sdp-plugin/cmd/sdp` |
| Checkpoint corrupted | Delete `.sdp/checkpoints/F067-*.json` and restart |
| Workstream blocked | Check dependencies in WS frontmatter |
| Review stuck | May need manual intervention |

---

## Implementation

**Go Package:** `internal/orchestrator`
**CLI Command:** `sdp-plugin/cmd/sdp/orchestrate.go`

Key components:
- `BeadsLoader` - Loads workstreams from files
- `CLIExecutor` - Executes workstreams
- `checkpoint.Manager` - Saves/restores state

---

## See Also

- `@build` - Execute single workstream
- `@review` - Quality review
- `@deploy` - Deployment workflow
- `.claude/patterns/tdd.md` - TDD pattern
- `.claude/patterns/quality-gates.md` - Quality gates

---

**Version:** 6.0.0 (LLM-agnostic)
