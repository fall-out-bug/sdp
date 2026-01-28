# Beads Integration - Phase 2 Progress

> **Date:** 2026-01-28
> **Status:** Phase 2.1 Complete - @idea Skill Updated
> **Worktree:** `feature/beads-integration`

---

## What's Done

### Phase 1: Foundation ✅ (Complete)

- ✅ `src/sdp/beads/` module (~830 LOC)
- ✅ MockBeadsClient + CLIBeadsClient
- ✅ Bidirectional sync service
- ✅ 16 tests passing
- ✅ Documentation

### Phase 2.1: @idea Skill ✅ (Just Complete)

- ✅ Updated `.claude/skills/idea/SKILL.md` for Beads workflow
- ✅ Creates Beads task instead of markdown draft
- ✅ Returns hash-based ID (e.g., `bd-0001`)
- ✅ Stores interview answers in `sdp_metadata`
- ✅ Optional markdown export for git history
- ✅ Environment variable `BEADS_USE_MOCK` for dev/prod切换

**Key Changes:**
```python
# Old workflow
@idea "Add auth"  # → docs/drafts/idea-add-auth.md

# New Beads workflow  
@idea "Add auth"  # → bd-0001 (Beads task)
```

---

## What's Next

### Phase 2.2: @design Skill (Next)

**Goal:** Decompose Beads task into sub-tasks with dependencies

**Implementation:**
```python
def design(beads_id: str) -> list[str]:
    # Read parent task
    parent = client.get_task(beads_id)
    
    # Decompose into workstreams
    workstreams = decompose_feature(parent.title)
    
    # Create sub-tasks with parent_id
    ws_ids = []
    for ws in workstreams:
        task = client.create_task(BeadsTaskCreate(
            title=ws.title,
            parent_id=beads_id,
            dependencies=ws.dependencies  # BeadsDependency objects
        ))
        ws_ids.append(task.id)
    
    return ws_ids
```

**Example:**
```bash
@design bd-0001

# Output:
✅ Created 3 workstreams:
   bd-0001.1: Domain entities (ready)
   bd-0001.2: Repository layer (blocked by bd-0001.1)
   bd-0001.3: Service layer (blocked by bd-0001.2)
```

### Phase 2.3: @build Skill

**Goal:** Execute workstream, update Beads status

**Implementation:**
```python
def build(beads_id: str):
    # Update to IN_PROGRESS
    client.update_task_status(beads_id, BeadsStatus.IN_PROGRESS)
    
    try:
        # Execute TDD cycle
        execute_tdd_cycle(beads_id)
        
        # Mark as done → unblocks next tasks
        client.update_task_status(beads_id, BeadsStatus.CLOSED)
        
        # Show newly ready tasks
        ready = client.get_ready_tasks()
        print(f"✅ Now ready: {ready}")
    except Exception as e:
        client.update_task_status(beads_id, BeadsStatus.BLOCKED)
        raise
```

### Phase 2.4: @oneshot Skill

**Goal:** Multi-agent coordination using Beads ready detection

**Implementation:**
```python
def oneshot(feature_id: str, num_agents: int = 3):
    with ThreadPoolExecutor(max_workers=num_agents) as executor:
        while True:
            ready = client.get_ready_tasks()
            if not ready:
                break
            
            # Execute ready tasks in parallel
            futures = [executor.submit(build, tid) for tid in ready]
            for f in futures:
                f.result()
```

---

## Current State

### Files Modified (This Session)

```
.claude/skills/idea/SKILL.md  # Updated for Beads workflow
```

### Files Created (This Session)

```
docs/workstreams/backlog/BEADS-001-skills-integration.md  # Workstream spec
```

### Total Progress

```
Phase 1: ████████████████████ 100% (Complete)
Phase 2: ████░░░░░░░░░░░░░░░░  20% (@idea done)
  2.1 @idea:    ✅ Complete
  2.2 @design:  ⏳ Next (1-2 days)
  2.3 @build:   📋 Pending (3-4 days)
  2.4 @oneshot: 📋 Pending (2-3 days)
  2.5 @review:  📋 Pending (1-2 days)
  2.6 Migration:📋 Pending (2-3 days)
```

---

## Testing

### Manual Test (Next Step)

```bash
cd /Users/fall_out_bug/projects/vibe_coding/sdp-beads-integration

# Test @idea with mock
export BEADS_USE_MOCK=true
@idea "Test feature"

# Expected output:
# ✅ Created Beads task: bd-0001
#    Title: Test feature
#    Status: BeadsStatus.OPEN
#    Priority: BeadsPriority.MEDIUM

# Verify task was created
poetry run python -c "
from sdp.beads import create_beads_client
client = create_beads_client(use_mock=True)
task = client.get_task('bd-0001')
print(f'✅ Task found: {task.title if task else 'Not found'}')
"
```

---

## Quick Reference

### Workflow Comparison

**Old (Markdown):**
```bash
@idea "Add auth"        # → docs/drafts/idea-add-auth.md
@design idea-add-auth   # → docs/workstreams/backlog/WS-*.md
@build WS-001-01        # → Updates file location
```

**New (Beads):**
```bash
@idea "Add auth"        # → bd-0001
@design bd-0001         # → bd-0001.1, bd-0001.2, bd-0001.3
@build bd-0001.1        # → Updates status (OPEN → IN_PROGRESS → CLOSED)
                        # → Automatically unblocks bd-0001.2
```

### Benefits

| Aspect | Old | New |
|--------|-----|-----|
| **ID conflicts** | Possible (manual PP-FFF-SS) | Impossible (hash-based) |
| **Multi-agent** | Manual F012 | Built-in |
| **Ready detection** | Manual script | `bd ready` |
| **Dependencies** | YAML list | Native DAG |
| **Status tracking** | File location | Status field |

---

## Next Actions

1. **Test @idea skill** - Manual test with mock
2. **Start @design skill** - Create sub-tasks with dependencies
3. **Update @build skill** - Status updates + ready detection
4. **Multi-agent test** - @oneshot with 3 parallel agents
5. **Documentation** - Update CLAUDE.md, README

---

**Timeline Estimate:** 2-3 weeks to complete Phase 2 (all skills)

**Decision Point:** After Phase 2, test with real Beads (Go + bd CLI) vs F012 comparison.
