# Git Safety Pattern

> Safe git operations that never destroy work

## Protected Actions

These actions are **FORBIDDEN** unless explicitly requested:

- `git push --force` (or `--force-with-lease` on shared branches)
- `git reset --hard` on shared branches
- `git clean -fd` without confirmation
- `git rebase -i` on pushed commits
- Deleting branches without merging

## Safe Defaults

```bash
# Preferred: rebase over merge for local sync
git pull --rebase

# Preferred: explicit push with branch
git push origin $(git branch --show-current)

# Always check before destructive operations
git status
git stash list
```

## Branch Convention

```
feature/F050-description    # New feature
bugfix/sdp-xxx-description  # Bug fix (P1/P2)
hotfix/sdp-xxx-description  # Emergency fix (P0)
docs/description            # Documentation only
```

## Commit Convention

```
feat(F050): add user authentication
fix(sdp-xxx): prevent panic on empty input
docs: update README with new commands
test(F051): add coverage for memory store
refactor: extract common validation logic
```

## PR Rules

- Target `dev` branch (NOT `main`) for features
- Target `main` only for releases
- Include test plan in PR body
- Wait for CI before merge

## Emergency Recovery

If something goes wrong:

```bash
# Undo last commit (keep changes)
git reset --soft HEAD~1

# Recover deleted branch
git reflog
git checkout -b recovered-branch <sha>

# Find lost commits
git fsck --lost-found
```

## See Also

- `@deploy` - Deployment workflow
- `.claude/skills/deploy.md` - Full skill definition
