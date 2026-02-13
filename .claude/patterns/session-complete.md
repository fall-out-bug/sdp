# Session Complete Pattern

> Checklist before declaring work done

## Mandatory Steps

Work is NOT complete until ALL steps pass:

### 1. Code Verification

```bash
# All tests pass
go test ./...

# Coverage meets threshold
go test -cover ./...

# Lint passes
golangci-lint run
```

### 2. Git Status

```bash
# Check what changed
git status

# Stage changes
git add <files>

# Commit
git commit -m "type: description"
```

### 3. Beads Sync

```bash
# Close completed issues
bd close <issue-id>

# Sync with remote
bd sync
```

### 4. Push

```bash
# Pull first
git pull --rebase

# Push to remote
git push

# VERIFY
git status  # MUST show "up to date with origin"
```

### 5. Cleanup

```bash
# Clear stashes if used
git stash list
git stash drop  # if appropriate

# Prune merged branches (optional)
git fetch --prune
```

## Common Mistakes

| Mistake | Consequence | Fix |
|---------|-------------|-----|
| Forgot `bd close` | Issue stays open | Run `bd close` |
| Forgot `bd sync` | Remote out of sync | Run `bd sync` |
| Forgot `git push` | Work not saved | Run `git push` |
| Forgot `git status` check | Unknown state | Always verify |

## Handoff Template

When ending session, provide context for next:

```markdown
## Session Summary

**Completed:**
- [Task 1]
- [Task 2]

**In Progress:**
- [Task 3] - 50% complete, blocked by [X]

**Next Steps:**
1. [Action item 1]
2. [Action item 2]

**Files Changed:**
- path/to/file1.go
- path/to/file2.go
```

## See Also

- `.claude/patterns/git-safety.md` - Git safety rules
- `.claude/patterns/quality-gates.md` - Quality gates
