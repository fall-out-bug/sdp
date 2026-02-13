---
name: oneshot
description: Autonomous multi-agent execution with review-fix loop and PR creation
cli: sdp orchestrate
version: 7.0.0
---

# oneshot

> **CLI (primary):** `sdp orchestrate <feature-id>`
> **Protocol (fallback):** Manual execution if CLI unavailable

Autonomous feature execution with review-fix loop and PR creation.

---

## EXECUTE THIS NOW

When user invokes `@oneshot F067`:

### Try CLI First (if available)

```bash
sdp orchestrate F067
```

If CLI succeeds, skip to **Step 4: Review-Fix Loop**.

If CLI fails or is unavailable, continue with manual execution below.

---

### Manual Execution (if CLI unavailable)

#### Step 1: Load Workstreams

```bash
ls docs/workstreams/backlog/00-067-*.md
```

Read each file for: WS ID, dependencies, AC, scope files.

#### Step 2: Build Dependency Graph

- Empty `depends_on` → run first
- Has dependencies → run after deps complete
- Topological sort for valid order

#### Step 3: Execute Workstreams

For each WS in order:
1. Invoke `@build {ws_id}`
2. Wait for completion
3. Checkpoint
4. Retry failures (max 2)

---

### Step 4: Review-Fix Loop (max 3 iterations)

After all workstreams complete:

1. Run `@review F067`
2. If APPROVED → proceed to Step 5
3. If CHANGES_REQUESTED:
   - P0: Fix immediately
   - P1: Create bugfix, then fix
   - P2+: Track only
4. Repeat

---

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

## CLI Reference

```bash
sdp orchestrate F067              # Execute all workstreams
sdp orchestrate resume F067       # Resume from checkpoint
sdp orchestrate --retry 3 F067    # Allow 3 retries per WS
```

---

## Finding Priority

| Priority | Action | Blocks? |
|----------|--------|---------|
| P0 | Fix immediately | YES |
| P1 | Create bugfix | YES |
| P2+ | Track only | NO |

---

## Resume After Interruption

```bash
cat .sdp/checkpoints/F067-*.json
```

---

## See Also

- `@build` - Execute single workstream
- `@review` - Quality review
- `.claude/patterns/tdd.md` - TDD pattern
