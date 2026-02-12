---
name: oneshot
description: Autonomous multi-agent execution with review-fix loop and PR creation
tools: Task, Read, Bash
version: 5.0.0
---

# @oneshot - Autonomous Feature Execution with Review-Fix Loop

Execute all workstreams, run review, auto-fix findings, create PR to dev.

## When to Use

- Feature has multiple workstreams (5-30 WS)
- Want autonomous execution with quality gates
- Need review-fix loop until approval
- Create PR for CI validation

## Invocation

```bash
@oneshot F050                       # Execute + review + fix + PR to dev
@oneshot F050 --max-reviews 5       # Allow 5 review iterations (default: 3)
@oneshot F050 --resume abc123       # Resume from checkpoint
@oneshot F050 --no-pr               # Skip PR creation (just review)
```

## How It Works

```
@oneshot F051
  │
  ├─► Phase 1: Execute Workstreams
  │     └─► @build 00-051-01, 00-051-02, ...
  │
  ├─► Phase 2: Review-Fix Loop (max 3 iterations)
  │     ├─► @review F051
  │     │     ├─► APPROVED → Exit loop
  │     │     └─► CHANGES_REQUESTED → Fix findings
  │     │           ├─► P0: Direct fix (security, blockers)
  │     │           ├─► P1: @bugfix for each finding
  │     │           └─► P2: Track only (don't block)
  │     └─► Repeat until APPROVED or max iterations
  │
  ├─► Phase 3: Verify No Blocking Findings
  │     └─► sdp guard finding list (must show 0 blocking)
  │
  └─► Phase 4: Create PR to dev
        ├─► git push origin feature/F051-xxx
        ├─► gh pr create --base dev --head feature/F051-xxx
        └─► CI validates automatically
```

## Orchestrator Agent Prompt

```python
Task(
    subagent_type="general-purpose",
    prompt=f"""
You are executing feature {feature_id} autonomously with review-fix loop.

**READ FIRST:** Read(".claude/agents/orchestrator.md")

**Workstreams to execute:**
{workstreams_list}

**Your workflow:**

## Phase 1: Execute Workstreams
1. Build dependency graph
2. Execute in topological order: @build {{ws_id}}
3. Update checkpoint after each WS
4. Commit after each WS

## Phase 2: Review-Fix Loop (max {max_reviews} iterations)

```
iteration = 1
while iteration <= {max_reviews}:
    # Run review
    result = @review {feature_id}

    if result.verdict == "APPROVED":
        print("✅ Review passed!")
        break

    # Review failed - fix findings
    print(f"⚠️ Review iteration {{iteration}}: CHANGES_REQUESTED")

    for finding in result.findings:
        if finding.priority == 0:  # P0 - critical
            # Fix immediately
            fix_security_issue(finding)
            git commit -m "fix: {{finding.title}}"

        elif finding.priority == 1:  # P1 - high
            # Create bugfix and execute
            @bugfix {{finding.beads_id}}

        else:  # P2+ - track only
            print("📋 Tracking: {{finding.title}}")

    iteration += 1

if iteration > {max_reviews}:
    escalate("Max review iterations reached")
```

## Phase 3: Verify Clean State
```bash
# Check for blocking findings
sdp guard finding list
# Must show: "0 blocking"

# If blocking findings exist, resolve them
for finding in $(sdp guard finding list --blocking); do
    resolve_finding "$finding"
done
```

## Phase 4: Create PR
```bash
# Ensure branch is pushed
git push origin {branch_name}

# Create PR to dev (NOT main)
gh pr create \\
    --base dev \\
    --head {branch_name} \\
    --title "feat({feature_id}): {feature_title}" \\
    --body "## Summary
{summary}

## Test plan
- [ ] All workstreams completed
- [ ] Review passed
- [ ] Tests pass locally

🤖 Generated with [Claude Code](https://claude.com/claude-code)"

echo "✅ PR created: $PR_URL"
echo "CI will validate automatically"
```

**CRITICAL RULES:**
- DO NOT merge to main - create PR to dev only
- Max {max_reviews} review iterations
- P0 findings must be fixed immediately
- P1 findings use @bugfix
- P2+ findings are tracked, not blocking
- Escalate if stuck after max iterations

**Progress format:**
```
[HH:MM] Phase 1: Executing {{ws_id}}...
[HH:MM] ✅ WS complete (Xm, Y% coverage)
[HH:MM] Phase 2: Review iteration {{n}}
[HH:MM] ⚠️ Found {{count}} findings, fixing...
[HH:MM] ✅ Review PASSED
[HH:MM] Phase 4: Creating PR...
[HH:MM] ✅ PR created: https://github.com/...
```

**Checkpoint file (.oneshot/{feature_id}-checkpoint.json):**
{{
  "feature": "{feature_id}",
  "phase": "review_fix",
  "iteration": 1,
  "completed_ws": ["..."],
  "review_result": "CHANGES_REQUESTED",
  "findings_fixed": ["sdp-xxx", ...],
  "pr_url": null
}}
"""
)
```

## Review-Fix Logic

### Finding Priority Handling

| Priority | Name | Action | Blocks PR? |
|----------|------|--------|------------|
| P0 | Critical | Fix immediately in code | YES |
| P1 | High | @bugfix for each | YES |
| P2 | Medium | Track in beads, don't fix | NO |
| P3 | Low | Track only | NO |

### Fix Commands by Review Area

```bash
# Security findings (P0)
# Fix in code directly, no new WS needed
# Example: Replace filepath.Join with security.SafeJoinPath

# TechLead findings (P1 - LOC violations)
# Split file into smaller modules
@build 99-{FEATURE}-01  # Refactor workstream

# SRE findings (P1 - observability)
# Add logging/context directly
@bugfix sdp-xxx

# Documentation findings (P1 - missing features)
# May need new workstream or @bugfix
```

### Guard Integration

```bash
# Before creating PR, verify no blocking findings
sdp guard finding list

# If blocking exists:
sdp guard finding resolve finding-xxx --by="Fixed in commit abc123"

# Clear resolved findings
sdp guard finding clear
```

## Output

**Success:**
```
✅ Feature F051 Complete

Phase 1: 9/9 workstreams executed
Phase 2: Review passed (2 iterations, 5 findings fixed)
Phase 3: 0 blocking findings
Phase 4: PR created

PR: https://github.com/owner/repo/pull/123
Branch: feature/F051-long-term-memory → dev
CI: Running...

Duration: 2h 30m
Avg Coverage: 85%
```

**Max Iterations Reached:**
```
⚠️ Feature F051 Needs Human Attention

Phase 1: 9/9 workstreams executed ✅
Phase 2: Review iterations exhausted (3/3)
Phase 3: 2 blocking findings remain

Blocking findings:
  - [Security] P0 Path traversal in drift detector
  - [TechLead] P1 store.go exceeds 200 LOC

Checkpoint: .oneshot/F051-checkpoint.json
Resume: @oneshot F051 --resume agent-xxx

Manual fix required for remaining findings.
```

## Checkpoint Format

```json
{
  "feature": "F051",
  "phase": "review_fix",
  "iteration": 2,
  "completed_ws": ["00-051-01", "00-051-02", ...],
  "review_results": [
    {"iteration": 1, "verdict": "CHANGES_REQUESTED", "findings": 5},
    {"iteration": 2, "verdict": "CHANGES_REQUESTED", "findings": 2}
  ],
  "findings_fixed": ["sdp-xxx", "sdp-yyy"],
  "pr_url": null,
  "started_at": "2026-02-12T10:00:00Z"
}
```

## CLI Flags

| Flag | Description |
|------|-------------|
| `--max-reviews N` | Max review iterations (default: 3) |
| `--no-pr` | Skip PR creation |
| `--resume ID` | Resume from checkpoint |
| `--dry-run` | Show plan without execution |

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Review stuck in loop | Check findings - may need human fix |
| CI failing on PR | Run `go test ./...` locally first |
| Can't create PR | Check branch exists and is pushed |
| Guard blocking PR | `sdp guard finding list` then fix |

## Example Session

```bash
User: @oneshot F051

Claude:
→ Launching orchestrator with review-fix loop...
→ Max iterations: 3

[Phase 1: Execute Workstreams]
→ [10:00] 00-051-01: Memory Store... ✅ (15m, 82%)
→ [10:15] 00-051-02: Search Engine... ✅ (20m, 85%)
→ ... (7 more workstreams)
→ [11:30] All workstreams complete

[Phase 2: Review-Fix Loop]
→ [11:30] Review iteration 1/3
→ ⚠️ CHANGES_REQUESTED: 5 findings
→   - P0: Path traversal (fixing...)
→   - P1: store.go 226 LOC (creating bugfix...)
→   - P1: Missing context (creating bugfix...)
→ [11:45] Fixed 3 findings, 2 tracked

→ [11:45] Review iteration 2/3
→ ✅ APPROVED

[Phase 3: Verify Clean]
→ sdp guard finding list: 0 blocking ✅

[Phase 4: Create PR]
→ git push origin feature/F051-long-term-memory
→ gh pr create --base dev --head feature/F051-long-term-memory
→ ✅ PR #123 created

✅ Feature F051 Complete
PR: https://github.com/owner/repo/pull/123
CI: Running validation...
```

---

**Version:** 5.0.0 (Review-Fix Loop + PR Creation)
**See Also:** `@review`, `@bugfix`, `@deploy`, `@build`
**Agent:** `.claude/agents/orchestrator.md`
