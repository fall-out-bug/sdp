---
name: design
description: Analyze Beads task and decompose into workstreams with dependencies. Creates sub-tasks with parent-child relationships.
tools: Read, Write, Edit, Bash, Glob, Grep, AskUserQuestion
version: 2.0.0-beads
---

# @design - Feature Decomposition (Beads Integration)

Analyze requirements and decompose Beads feature tasks into workstreams with sequential dependencies.

## When to Use

- After `@idea` creates a Beads task
- When a feature needs to be broken into workstreams
- Before starting implementation
- When architectural decisions need user input

## Beads vs Markdown Workflow

**This skill creates Beads sub-tasks** with hash-based IDs and dependencies.

For traditional markdown workflow, use `prompts/commands/design.md` instead.

## Invocation

```bash
@design bd-0001
```

**Environment Variables:**
- `BEADS_USE_MOCK=true` - Use mock Beads (default for dev)
- `BEADS_USE_MOCK=false` - Use real Beads CLI (requires Go + bd installed)

## Workflow

**IMPORTANT:** Use AskUserQuestion for architectural decisions before decomposition.

### Step 1: Initialize Beads Client

```python
from sdp.beads import create_beads_client, FeatureDecomposer
import os

use_mock = os.getenv("BEADS_USE_MOCK", "true").lower() == "true"
client = create_beads_client(use_mock=use_mock)
decomposer = FeatureDecomposer(client)
```

### Step 2: Read Parent Task

```python
# Get feature task from @idea
feature = client.get_task(beads_id)

if not feature:
    print(f"❌ Task not found: {beads_id}")
    return

print(f"📋 Designing: {feature.title}")
print(f"   Description: {feature.description[:100]}...")
```

### Step 3: Interactive Planning

**Use AskUserQuestion** for architectural decisions:

```markdown
AskUserQuestion({
  "questions": [{
    "question": "What is the complexity level of this feature?",
    "header": "Complexity",
    "options": [
      {"label": "Simple (1-2 workstreams)", "description": "Straightforward, minimal integration"},
      {"label": "Medium (3-5 workstreams)", "description": "Standard complexity, some integration points"},
      {"label": "Large (6+ workstreams)", "description": "Complex, multiple integrations, significant changes"}
    ],
    "multiSelect": false
  }, {
    "question": "Which layers need implementation?",
    "header": "Layers",
    "options": [
      {"label": "Domain", "description": "Business logic, entities, value objects"},
      {"label": "Repository", "description": "Data access, persistence layer"},
      {"label": "Service", "description": "Application services, use cases"},
      {"label": "API/Presentation", "description": "Endpoints, controllers, UI"}
    ],
    "multiSelect": true
  }]
})
```

**Continue interviewing** about:
- Database schema changes
- External API integrations
- Authentication/authorization needs
- Performance requirements
- Testing strategy

### Step 4: Determine Workstreams

Based on interview answers, determine workstreams:

**Simple feature:**
```python
# Default 3 workstreams
ws_ids = decomposer.decompose(beads_id)
```

**Medium feature:**
```python
# Custom workstreams
from sdp.beads import WorkstreamSpec

custom_workstreams = [
    WorkstreamSpec(title="Domain model", sequence=1, size="MEDIUM"),
    WorkstreamSpec(title="Database schema", sequence=2, size="MEDIUM", 
                    dependencies=["ws-001"]),
    WorkstreamSpec(title="Repository layer", sequence=3, size="MEDIUM",
                    dependencies=["ws-002"]),
    WorkstreamSpec(title="Service layer", sequence=4, size="MEDIUM",
                    dependencies=["ws-003"]),
    WorkstreamSpec(title="API endpoints", sequence=5, size="MEDIUM",
                    dependencies=["ws-004"]),
]

ws_ids = decomposer.decompose(beads_id, workstreams=custom_workstreams)
```

**Large feature:**
```python
# More granular workstreams
custom_workstreams = [
    WorkstreamSpec(title="Domain entities", sequence=1, size="MEDIUM"),
    WorkstreamSpec(title="Value objects", sequence=2, size="SMALL",
                    dependencies=["ws-001"]),
    WorkstreamSpec(title="Database migration", sequence=3, size="MEDIUM"),
    WorkstreamSpec(title="Repository interface", sequence=4, size="SMALL",
                    dependencies=["ws-001"]),
    WorkstreamSpec(title="Repository implementation", sequence=5, size="MEDIUM",
                    dependencies=["ws-003", "ws-004"]),
    WorkstreamSpec(title="Service interface", sequence=6, size="SMALL",
                    dependencies=["ws-001"]),
    WorkstreamSpec(title="Service implementation", sequence=7, size="MEDIUM",
                    dependencies=["ws-005", "ws-006"]),
    WorkstreamSpec(title="API controllers", sequence=8, size="MEDIUM",
                    dependencies=["ws-007"]),
]

ws_ids = decomposer.decompose(beads_id, workstreams=custom_workstreams)
```

### Step 5: Verify Dependencies

```python
# Check ready tasks
ready = client.get_ready_tasks()
print(f"✅ Ready to start: {ready}")

# Verify dependencies
for i, ws_id in enumerate(ws_ids):
    ws = client.get_task(ws_id)
    deps = ws.dependencies
    
    if deps:
        print(f"  {ws_id} ({ws.title}) blocked by {len(deps)} tasks")
    else:
        print(f"  {ws_id} ({ws.title}) ready to start")
```

### Step 6: Export to Markdown (Optional)

```python
# Export for human reference
markdown_path = f"docs/workstreams/beads-{beads_id}.md"

with open(markdown_path, "w") as f:
    f.write(f"# Workstreams for {feature.title}\n\n")
    f.write(f"> **Parent Task:** {beads_id}\n")
    f.write(f"> **Created:** {datetime.utcnow().isoformat()}\n\n")
    
    for i, ws_id in enumerate(ws_ids, 1):
        ws = client.get_task(ws_id)
        f.write(f"## {i}. {ws.title}\n")
        f.write(f"**ID:** {ws_id}\n")
        f.write(f"**Status:** {ws.status.value}\n")
        f.write(f"**Priority:** {ws.priority.value}\n")
        
        if ws.dependencies:
            dep_ids = [d.task_id for d in ws.dependencies]
            f.write(f"**Dependencies:** {', '.join(dep_ids)}\n")
        f.write("\n")
```

## Output

**Primary:** List of Beads workstream IDs (e.g., `[bd-0001.1, bd-0001.2, bd-0001.3]`)

**Secondary:** Optional markdown export to `docs/workstreams/beads-{parent_id}.md`

**Beads Sub-Tasks:**
- `id`: Hash-based IDs (auto-generated, e.g., `bd-0001.1`)
- `title`: Workstream title
- `parent_id`: Reference to parent feature task
- `status`: OPEN (default)
- `dependencies`: List of blocking dependencies
- `sdp_metadata`: Workstream sequence, size, etc.

## Next Steps

After decomposition:

1. **Check ready tasks:**
   ```bash
   bd ready
   
   # Output:
   # Ready tasks:
   # - bd-0001.1 (Domain entities)
   ```

2. **Start execution:**
   ```bash
   @build bd-0001.1
   ```

3. **Monitor progress:**
   ```bash
   bd status --watch
   
   # Automatically shows new tasks as they become ready
   ```

## Example Session

```bash
# Decompose feature
@design bd-0001

# ... (interviewing happens) ...

# Output:
✅ Created 3 workstreams:
   bd-0001.1: Domain entities
   bd-0001.2: Repository layer (blocked by bd-0001.1)
   bd-0001.3: Service layer (blocked by bd-0001.2)

# Check what's ready
bd ready

# Output:
Ready tasks:
- bd-0001 (parent)
- bd-0001.1 (Domain entities)

# Start execution
@build bd-0001.1

# After completion, bd-0001.2 automatically becomes ready!
bd ready

# Output:
Ready tasks:
- bd-0001
- bd-0001.2 (Repository layer)
```

## Key Principles

**Decomposition Strategy:**
1. **Sequential by default** - Each WS blocks the next (safe, simple)
2. **Size matters** - Keep workstreams SMALL/MEDIUM (< 500 LOC)
3. **Dependencies explicit** - Use BeadsDependency.BLOCKS for sequencing
4. **Parallel when possible** - Independent tasks can use same sequence number

**Beads Integration:**
1. **Hash-based IDs** - No conflicts, auto-generated
2. **Parent-child** - `parent_id` links workstreams to feature
3. **Native DAG** - Beads manages dependency graph
4. **Ready detection** - `bd ready` shows executable tasks

**Workstream Sizing:**
- **SMALL:** < 500 LOC, < 1500 tokens
- **MEDIUM:** 500-1500 LOC, 1500-5000 tokens
- **LARGE:** > 1500 LOC → Break into 2+ workstreams

## Migration from Markdown Workflow

**Old workflow:**
```bash
@design idea-add-auth
# → docs/workstreams/backlog/WS-001-01.md
# → docs/workstreams/backlog/WS-001-02.md
# ...
```

**New Beads workflow:**
```bash
@design bd-0001
# → bd-0001.1 (Domain)
# → bd-0001.2 (Repository)
# → bd-0001.3 (Service)
```

**Benefits:**
- No manual ID allocation (hash-based, automatic)
- Multi-agent ready (execute WS in parallel)
- Built-in ready detection (no manual scripts)
- Automatic unblocking (complete WS1 → WS2 becomes ready)

## Troubleshooting

**No workstreams created:**
```bash
# Check parent task exists
bd show bd-0001

# Check decomposer error logs
export DEBUG=1
```

**Dependencies incorrect:**
```bash
# View task details
bd show bd-0001.2

# Verify dependencies
bd dep list bd-0001.2
```

**Tasks not becoming ready:**
```bash
# Check blocking tasks
bd ready

# View dependency graph
bd graph bd-0001
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@design bd-0001` | Decompose into workstreams |
| `bd show {id}` | View task/workstream details |
| `bd ready` | List ready tasks |
| `bd dep list {id}` | List dependencies |
| `@build {ws_id}` | Execute workstream |

---

**Version:** 2.0.0-beads
**Status:** Beads Integration
**See Also:** `@idea`, `@build`, `@oneshot`
