# 00-052-00: Backup & Worktree Setup

> **Beads ID:** sdp-wqv8
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 0 - Preparation
> **Size:** SMALL
> **Duration:** 1 day
> **Dependencies:** None (must be first)

## Goal

Create safety net and isolated workspace for multi-agent implementation.

## Acceptance Criteria

- **AC1:** Git tag `before-multi-agent-vision` created and pushed
- **AC2:** Git worktree created at `/Users/fall_out_bug/projects/vibe_coding/sdp-multi-agent`
- **AC3:** Worktree verified to be on `dev` branch
- **AC4:** Original repository remains untouched

## Files

**Create:**
- Git tag: `before-multi-agent-vision`

**Create:**
- Worktree: `/Users/fall_out_bug/projects/vibe_coding/sdp-multi-agent`

## Steps

### Step 1: Create Backup Tag

```bash
git tag before-multi-agent-vision
git push origin before-multi-agent-vision
```

**Verify:**
```bash
git tag -l before-multi-agent-vision
git log --oneline -1 before-multi-agent-vision
```

Expected: Tag shows current commit SHA

### Step 2: Create Worktree

```bash
cd /Users/fall_out_bug/projects/vibe_coding
git worktree add sdp-multi-agent sdp/dev
cd sdp-multi-agent
```

**Verify:**
```bash
pwd
git status
```

Expected: `/Users/fall_out_bug/projects/vibe_coding/sdp-multi-agent`, on `dev` branch, clean working tree

### Step 3: Document Worktree Location

Create `docs/workstreams/WORKTREE.md` in worktree:

```markdown
# Implementation Worktree

**Location:** `/Users/fall_out_bug/projects/vibe_coding/sdp-multi-agent`

**Branch:** `dev`

**Purpose:** Isolated workspace for multi-agent implementation

**Sync to main:**
```bash
cd /Users/fall_out_bug/projects/vibe_coding/sdp-multi-agent
git pull origin dev
git push origin dev
```

**Cleanup (after completion):**
```bash
cd /Users/fall_out_bug/projects/vibe_coding
git worktree remove sdp-multi-agent
```
```

## Quality Gates

- Tag exists in local and remote
- Worktree directory exists
- Worktree is clean git status
- Original repo unaffected

## Success Metrics

- Zero impact to main repository
- Safe rollback possible via git tag
- Isolated workspace ready for implementation
