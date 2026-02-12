# Incident: Agent Branch Confusion (2026-02-12)

## Summary

Agents committed to wrong branches, mixing features together and corrupting branch tracking.

## Symptoms

1. **Branch tracking corruption:** `feature/F063` was tracking `origin/dev` instead of `origin/feature/F063`
2. **Cross-feature commits:** F064 commits appeared in F063 branch history
3. **Direct pushes to dev:** F063 commits pushed directly to `origin/dev` without PR
4. **Bugfix branch confusion:** `bugfix/sdp-67l6` was merged into feature branches instead of being standalone

## Evidence

```bash
# feature/F063 tracking wrong remote
* feature/F063 76e88f9 [origin/dev: ahead 1]  # Should track origin/feature/F063

# F064 commits in F063 branch
origin/feature/F063 contains:
- d255df9 fix(resolver): path traversal (F064 commit!)
- 14e5d08 refactor(resolver,task) (F064 commit!)

# F063 commits pushed directly to dev (no PR)
origin/dev contains F063 commits without going through PR #28
```

## Root Cause Analysis

### Likely Causes

1. **Working directory reset after tool calls:**
   ```
   Shell cwd is reset to /Users/fall_out_bug/projects/vibe_coding/sdp
   ```
   After Bash commands, shell returns to main repo instead of worktree.

2. **Agent doesn't verify current directory before git operations:**
   - Agent runs `git commit` without checking `pwd`
   - Agent runs `git push` without verifying branch

3. **No branch isolation enforcement:**
   - No guard rails preventing commits to wrong branch
   - No verification that worktree is on correct branch

## Prevention Measures

### Immediate (Critical)

1. **Always verify directory before git operations:**
   ```bash
   # BEFORE any git operation
   pwd  # Must be in correct worktree
   git branch --show-current  # Must be correct branch
   ```

2. **Use absolute paths for all git commands in worktrees:**
   ```bash
   cd /Users/fall_out_bug/projects/vibe_coding/sdp-F063 && git status
   # Instead of assuming we're already there
   ```

3. **Add pre-commit hook to verify branch:**
   ```bash
   # .git/hooks/pre-commit
   CURRENT_BRANCH=$(git branch --show-current)
   EXPECTED_BRANCH=$(cat .sdp/active-branch 2>/dev/null)
   if [ -n "$EXPECTED_BRANCH" ] && [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
     echo "ERROR: Wrong branch! Expected $EXPECTED_BRANCH, on $CURRENT_BRANCH"
     exit 1
   fi
   ```

### Long-term (Recommended)

1. **Add branch guard to orchestrator:**
   - Orchestrator should verify worktree branch before spawning agents
   - Agents should inherit branch context

2. **Git operations wrapper:**
   ```go
   func SafeGit(op string, worktreePath string, expectedBranch string) error {
       // 1. Verify worktreePath exists
       // 2. Verify current branch == expectedBranch
       // 3. Execute git command
       // 4. Verify still on expected branch
   }
   ```

3. **Worktree isolation:**
   - Each worktree should have `.sdp/active-feature` file
   - Hooks verify this file matches current branch

## Action Items

- [ ] Add branch verification to all agent prompts
- [ ] Create git safety wrapper in sdp CLI
- [ ] Add pre-commit hook for branch verification
- [ ] Document "always verify pwd" rule in agent guidelines

## Related Issues

- beads: sdp-45q9, sdp-5ho2, sdp-67l6 (created during confused agent session)
