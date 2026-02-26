# Git Safety Guidelines for Agents

> **Purpose:** Prevent agents from committing to wrong worktrees/branches.
> **Last Updated:** 2026-02-12
> **Related:** Agent Git Safety Protocol

---

## CRITICAL RULES

1. **NEVER run git commands without verifying context**
2. **ALWAYS check pwd before git operations**
3. **ALWAYS verify branch before commits**
4. **NEVER push to branches you weren't assigned to**
5. **Features MUST be implemented in feature branches (not dev/main)**

---

## Pre-Flight Checklist

Before ANY git operation, run:

```bash
pwd                         # Verify worktree path
git branch --show-current   # Verify branch name
sdp guard context check     # Validate session
```

---

## Required Pattern

```bash
# CORRECT - verify first, then commit
sdp guard context check && git commit -m "message"

# WRONG - no verification
git commit -m "message"

# WRONG - assumes CWD is correct (it resets after tool calls!)
cd /path && git commit
```

---

## Error Recovery

If context check fails:

1. **DO NOT proceed** with git command
2. Run: `sdp guard context find $EXPECTED_FEATURE`
3. Run: `sdp guard context go $EXPECTED_FEATURE`
4. Retry git command

---

## Forbidden Commands

| Command | Reason |
|---------|--------|
| `git push --force` | **NEVER** - destroys history |
| `git checkout main` | Only with explicit user request |
| `git checkout dev` | Only for non-feature work |
| `git merge` across features | Can cause cross-contamination |
| `git rebase` | Only with user approval |
| Direct commits to `main` or `dev` | Features go in feature branches |

---

## Feature Branch Rule

ALL feature implementation MUST happen in feature branches:

```bash
# Create feature branch
git checkout -b feature/<feature-id>

# Never commit to dev or main directly
# Guard will reject commits to protected branches
```

### Branch Naming Convention

| Branch Type | Pattern | Example |
|-------------|---------|---------|
| Feature | `feature/<feature-id>` | `feature/auth-login` |
| Bugfix | `bugfix/issue-id` | `bugfix/sdp-1234` |
| Hotfix | `hotfix/issue-id` | `hotfix/sdp-1234` |

### Protected Branches

| Branch | Allowed Operations |
|--------|-------------------|
| `main` | Merge only (via PR) |
| `dev` | Merge only (via PR) |

---

## CWD Reset Issue

**IMPORTANT:** After Bash tool calls, the shell state (including CWD) may reset.

### Problem

```bash
# This does NOT work reliably!
cd /path/to/worktree && git status
# CWD may reset after the command!
```

### Solution

Always use absolute paths or verify CWD before git operations:

```bash
# Verify context first
sdp guard context check

# Or use absolute paths
git -C /absolute/path/to/worktree status
```

---

## Session File Format

Sessions are stored in `.sdp/session.json`:

```json
{
  "version": "1.0",
  "worktree_path": "/Users/user/projects/sdp-<feature-id>",
  "feature_id": "<feature-id>",
  "expected_branch": "feature/<feature-id>",
  "expected_remote": "origin/feature/<feature-id>",
  "created_at": "2026-02-12T10:00:00Z",
  "created_by": "sdp worktree create",
  "hash": "sha256:abc123..."
}
```

### Session Commands

```bash
# Initialize session for feature
sdp session init --feature=<feature-id>

# Sync session with actual git state
sdp session sync

# Show session details
sdp session show

# Repair corrupted session
sdp session repair --force
```

---

## Quick Reference

| Scenario | Action |
|----------|--------|
| Before any git command | `sdp guard context check` |
| Context check fails | `sdp guard context go F###` |
| Find worktree for feature | `sdp guard context find F###` |
| Create feature branch | `git checkout -b feature/F###` |
| Session file corrupted | `sdp session repair --force` |

---

## See Also

- [docs/plans/2026-02-12-agent-git-safety-protocol.md](../docs/plans/2026-02-12-agent-git-safety-protocol.md)
- [@build skill](skills/build/SKILL.md)
- [@deploy skill](skills/deploy/SKILL.md)
