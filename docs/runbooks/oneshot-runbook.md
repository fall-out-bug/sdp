# /oneshot Runbook

**Purpose:** Execute all workstreams of a feature autonomously with checkpoint/resume support.

**Prerequisites:**
- Feature designed with `/design`
- All workstreams in `tools/hw_checker/docs/workstreams/backlog/`
- Git repository initialized
- `sdp` package installed

## When to Use

- Feature has 3+ workstreams
- Want autonomous execution (hands-off)
- Need checkpoint/resume capability
- Complex feature with dependencies

## Workflow

### Phase 1: Feature Discovery

```bash
# 1. List all WS for feature
ls tools/hw_checker/docs/workstreams/backlog/WS-060-*.md

# 2. Verify dependencies are clear
# Check each WS for "Dependency:" section
```

### Phase 2: Execution (Autonomous)

**In Claude Code:**
```
/oneshot F60
```

**In Cursor:**
```
/oneshot F60
```

**In OpenCode:**
```
/oneshot F60
```

**Note:** If your local setup still has legacy command aliases, keep `/oneshot-simple` as a temporary fallback.

### Phase 3: What /oneshot Does

The orchestrator will:

1. **Create PR for approval** (GitHub integration)
2. **Wait for human approval**
3. **Read feature context**
4. **Build execution plan** (dependency order)
5. **Execute each WS** in loop:
   - Pre-build checks
   - Execute `/build WS-ID`
   - Post-build checks
   - Update checkpoint
   - Update progress JSON
   - Commit changes
6. **Handle errors**:
   - CRITICAL: Save checkpoint, notify, STOP
   - HIGH: Auto-fix, retry (max 2x)
   - MEDIUM: Mark needs_review, continue
7. **Run final review**:
   - Post-oneshot hooks
   - Execute `/codereview F{XX}`
   - Generate UAT guide

### Phase 4: Checkpoint/Resume

If execution fails, it can be resumed:

```bash
# Check checkpoint file
cat .sdp/checkpoint.json

# Resume execution (automatic when /oneshot is called again)
/oneshot F60
```

## Checkpoint Format

```json
{
  "feature_id": "F60",
  "current_ws": "WS-060-03",
  "completed_ws": ["WS-060-01", "WS-060-02"],
  "status": "failed",
  "error": "...",
  "timestamp": "2026-01-23T10:00:00Z"
}
```

## Progress Tracking

Real-time progress JSON:

```json
{
  "feature_id": "F60",
  "total_ws": 5,
  "completed_ws": 2,
  "current_ws": "WS-060-03",
  "progress_percent": 40,
  "estimated_remaining": "2h 30m",
  "status": "in_progress"
}
```

## Error Handling

| Severity | Action | Retry Policy |
|----------|--------|--------------|
| **CRITICAL** | Save checkpoint, notify, STOP | N/A |
| **HIGH** | Auto-fix, retry | Max 2 attempts |
| **MEDIUM** | Mark needs_review, continue | N/A |
| **LOW** | Log warning, continue | N/A |

## Verification After Execution

```bash
# 1. Check all WS completed
grep "status: completed" tools/hw_checker/docs/workstreams/backlog/WS-060-*.md

# 2. Verify all AC checked
grep "✅" tools/hw_checker/docs/workstreams/backlog/WS-060-*.md

# 3. Check execution reports
grep "Execution Report" tools/hw_checker/docs/workstreams/backlog/WS-060-*.md

# 4. Check commits
git log --oneline --grep="WS-060"

# 5. Run codereview
/codereview F60
```

## Troubleshooting

### /oneshot not found

**Problem:** Command not found

**Solution:**
- Claude Code: Check `.claude/skills/oneshot/SKILL.md` exists
- Cursor: Check `.cursor/commands/oneshot.md` exists
- OpenCode: Check `.opencode/commands/oneshot-simple.md` exists

### PR approval hangs

**Problem:** Waiting for GitHub PR approval

**Solution:**
1. Check `GITHUB_TOKEN` is set
2. Check `GITHUB_REPO` is set
3. Verify token has `repo` scope
4. Check PR was created: `gh pr list`

### Checkpoint not resuming

**Problem:** /oneshot starts from beginning instead of resuming

**Solution:**
1. Check checkpoint exists: `cat .sdp/checkpoint.json`
2. Verify feature ID matches
3. Check orchestrator agent can read checkpoint
4. Try manual resume: `/build WS-ID` for next WS

### Execution stuck on one WS

**Problem:** Looping on same WS

**Solution:**
1. Check error logs
2. Verify issue is not HIGH severity (auto-retry loop)
3. Manually fix issue: `/debug "WS stuck"`
4. Manually skip WS: Mark WS completed in checkpoint

## Integration with Other Commands

```
/idea "feature description"
    ↓
/design idea-slug
    ↓
/oneshot F{XX}
    ↓
/codereview F{XX}
    ↓
Human UAT
    ↓
/deploy F{XX}
```

## Best Practices

1. **Test manually first:** Execute first WS manually before using /oneshot
2. **Small features:** /oneshot works best with 3-10 workstreams
3. **Check dependencies:** Ensure all dependencies are clear before starting
4. **Monitor progress:** Watch progress JSON for real-time status
5. **Plan for errors:** Know how to resume if execution fails
6. **Review PR:** Always review and approve the GitHub PR before execution

## References

- Master prompt: `sdp/prompts/commands/oneshot.md`
- Orchestrator agent: `.claude/agents/orchestrator.md`
- Feature spec: `tools/hw_checker/docs/specs/feature_{XX}/feature.md`
