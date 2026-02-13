# oneshot

> **Protocol:** This skill describes the complete execution process
> **CLI (optional):** `sdp orchestrate <feature-id>`

Autonomous feature execution with review-fix loop and PR creation.

---

## EXECUTE THIS NOW

When user invokes `@oneshot F067`, follow this protocol:

### Step 1: Load Workstreams

```bash
ls docs/workstreams/backlog/00-067-*.md
```

Read each file to understand:
- Workstream ID
- Dependencies (`depends_on` in frontmatter)
- Acceptance Criteria
- Scope files

### Step 2: Build Dependency Graph

Determine execution order:
- WS with empty `depends_on` → can run first
- WS that depends on others → run after dependencies complete
- Use topological sort for valid order

### Step 3: Execute Workstreams

For each workstream in dependency order:

1. **Invoke @build** with the workstream ID
2. **Wait for completion** before starting dependent WS
3. **Checkpoint** after each successful WS
4. **Handle failures**: Retry up to 2 times, then escalate

### Step 4: Review-Fix Loop (max 3 iterations)

After all workstreams complete:

1. Run `@review F067`
2. If APPROVED → proceed to Step 5
3. If CHANGES_REQUESTED:
   - P0: Fix immediately
   - P1: Create bugfix, then fix
   - P2+: Track only
4. Repeat until APPROVED or max iterations

### Step 5: Verify Clean State

```bash
sdp guard finding list  # Must show 0 blocking
```

### Step 6: Create PR

```bash
git push origin feature/F067-xxx
gh pr create --base dev --head feature/F067-xxx
```

---

## CLI Optimization (Optional)

If CLI is available and working, it can handle Steps 1-3:

```bash
sdp orchestrate F067              # Execute all workstreams
sdp orchestrate resume F067       # Resume from checkpoint
```

If CLI fails or is unavailable, use the protocol above.

---

## Finding Priority

| Priority | Action | Blocks? |
|----------|--------|---------|
| P0 | Fix immediately | YES |
| P1 | Create bugfix with `bd create` | YES |
| P2+ | Track only | NO |

---

## Resume After Interruption

```bash
cat .sdp/checkpoints/F067-*.json
```

Continue from the first incomplete workstream.

---

## See Also

- `@build` - Execute single workstream
- `@review` - Quality review
- `.claude/patterns/tdd.md` - TDD pattern
