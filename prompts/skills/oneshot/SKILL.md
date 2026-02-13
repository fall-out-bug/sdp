---
name: oneshot
description: Autonomous multi-agent execution with review-fix loop and PR creation
cli: sdp orchestrate (file ops only - requires @build for actual work)
version: 7.1.0
---

# oneshot

> **CLI:** `sdp orchestrate <feature-id>` — handles file loading, dependency graph, checkpoints
> **LLM:** Required for actual workstream execution via `@build`

Autonomous feature execution with review-fix loop and PR creation.

---

## EXECUTE THIS NOW

When user invokes `@oneshot F067`:

### Step 1: Load Workstreams

```bash
ls docs/workstreams/backlog/00-067-*.md
```

Read each file for: WS ID, dependencies, AC, scope files.

### Step 2: Build Dependency Graph

Determine execution order:
- Empty `depends_on` → run first
- Has dependencies → run after deps complete
- Topological sort for valid order

### Step 3: Execute Workstreams (LLM required)

For each WS in dependency order:

1. **Invoke @build** with the workstream ID:
   ```
   @build 00-067-01
   ```

2. **Wait for completion** before starting dependent WS

3. **Checkpoint** after each successful WS

4. **Handle failures**: Retry up to 2 times, then escalate

**CRITICAL:** The CLI (`sdp orchestrate`) only handles file operations. Actual code changes require `@build` which spawns LLM subagents.

### Step 4: Review-Fix Loop (max 3 iterations)

After all workstreams complete:

1. Run `@review F067`
2. If APPROVED → proceed to Step 5
3. If CHANGES_REQUESTED:
   - P0: Fix immediately
   - P1: Create bugfix, then fix
   - P2+: Track only
4. Repeat

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

## CLI Reference (Optional Optimization)

The CLI can handle Steps 1-2 and checkpointing:

```bash
sdp orchestrate F067              # Load workstreams, build graph
sdp orchestrate resume F067       # Resume from checkpoint
```

But Step 3 (actual execution) still requires `@build` invocation by an LLM.

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

Continue from the first incomplete workstream.

---

## See Also

- `@build` - Execute single workstream (REQUIRED)
- `@review` - Quality review
- `.claude/patterns/tdd.md` - TDD pattern
