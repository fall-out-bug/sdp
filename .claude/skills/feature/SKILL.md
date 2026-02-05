---
name: feature
description: Feature development with multi-agent discovery
tools: Read, Write, Edit, Bash, AskUserQuestion, Task
version: 5.0.0
---

# @feature - Multi-Agent Feature Development

Spawn expert agents for discovery → analysis → planning.

## When to Use

- New feature development
- Product ideation
- Requirements gathering
- Feature planning

## Workflow

### Step 1: Quick Interview (3-5 questions)

AskUserQuestion:
- What problem do we solve?
- Who are the users?
- Success metrics?
- Timeline/urgency?

### Step 2: Spawn Discovery Agents (PARALLEL)

**Detect Beads:**
```bash
if bd --version &>/dev/null && [ -d .beads ]; then
  BEADS_ENABLED=true
else
  BEADS_ENABLED=false
fi
```

**Spawn 4 agents in parallel (single message with multiple Task calls):**

```python
# Agent 1: Business Analyst
Task(
    subagent_type="general-purpose",
    prompt="""You are the BUSINESS ANALYST expert.

Read .claude/agents/business-analyst.md for your specification.

FEATURE: {feature_description}

Your task:
1. Discover user needs
2. Write user stories (Given/When/Then)
3. Define success metrics (KPIs)
4. Identify stakeholders

Output format (as per business-analyst.md):
## Business Requirements
{stakeholders, problem, user stories, KPIs}

BEADS_INTEGRATION:
If Beads enabled ($BEADS_ENABLED=true):
- Create feature task: `bd create "Feature: {title}" --parent {project}`
- Link user stories to task
- Otherwise: Skip Beads operations
""",
    description="Business analysis"
)

# Agent 2: Product Manager
Task(
    subagent_type="general-purpose",
    prompt="""You are the PRODUCT MANAGER expert.

Read .claude/agents/product-manager.md for your specification.

FEATURE: {feature_description}

Your task:
1. Define product vision
2. Prioritize features (RICE framework)
3. Create roadmap (quarterly)
4. Define success metrics

Output format:
## Product Requirements
{vision, prioritization, roadmap, KPIs}

BEADS_INTEGRATION:
If Beads enabled:
- Update feature task with roadmap
- Create child tasks for high-priority features
- Otherwise: Skip
""",
    description="Product management"
)

# Agent 3: Systems Analyst
Task(
    subagent_type="general-purpose",
    prompt="""You are the SYSTEMS ANALYST expert.

Read .claude/agents/systems-analyst.md for your specification.

FEATURE: {feature_description}

Your task:
1. Define functional requirements
2. Specify APIs/interfaces
3. Design data models
4. Document use cases

Output format:
## Functional Specification
{FRs, APIs, data models, use cases}

BEADS_INTEGRATION:
If Beads enabled:
- Note: Workstreams will be created by Technical Decomposition agent
- Otherwise: Skip
""",
    description="Systems analysis"
)

# Agent 4: Technical Decomposition
Task(
    subagent_type="general-purpose",
    prompt="""You are the TECHNICAL DECOMPOSITION expert.

Read .claude/agents/technical-decomposition.md for your specification.

FEATURE: {feature_description}
INPUT FROM:
- Business Analyst (user stories)
- Product Manager (priorities)
- Systems Analyst (functional specs)

Your task:
1. Break into workstreams
2. Define dependencies
3. Estimate effort (T-shirt sizing)
4. Identify critical path

Output format:
## Workstream Breakdown
{workstreams with AC, dependencies, estimates}

BEADS_INTEGRATION:
If Beads enabled:
- Create task per workstream: `bd create "WS-XXX: {title}" --parent {feature}`
- Set dependencies: `bd update {ws} --blocks {other_ws}`
- Update .beads-sdp-mapping.jsonl with ws_id → beads_id
- Otherwise: Skip Beads operations
""",
    description="Technical decomposition"
)
```

### Step 3: Synthesize Results

Wait for all 4 agents to complete, then:

```markdown
## Feature Specification: {feature_name}

### Vision (from Product Manager)
{vision statement}

### User Stories (from Business Analyst)
{prioritized stories with acceptance criteria}

### Functional Requirements (from Systems Analyst)
{FRs, APIs, data models}

### Workstreams (from Technical Decomposition)
| WS | Title | Size | Dependencies | Priority |
|----|-------|------|--------------|----------|
| WS-001 | {title} | M | None | P0 |

### Success Metrics (from BA + PM)
{KPIs with targets}

### Next Steps
- Review workstreams
- Execute: @design {feature_id} or @oneshot {feature_id}
```

### Step 4: Save Outputs

```bash
# Save feature spec
mkdir -p docs/drafts
cat > docs/drafts/{feature_id}.md << 'EOF'
{synthesized specification}
EOF

# If Beads enabled: Update parent feature task
if [ "$BEADS_ENABLED" = true ]; then
  bd show {feature_id}  # Verify parent exists
fi
```

## Output

**Success:**
```
✅ Feature specification created
📄 docs/drafts/FXXX.md
📊 Workstreams: N defined (P0: X, P1: Y)
📌 Beads: {N tasks created if enabled}
```

## Beads Detection

All agents check `$BEADS_ENABLED` before Beads operations:
- If enabled: Create/update/link tasks
- If disabled: Markdown-only workflow

## Parallel Execution Pattern

Spawning 4 agents simultaneously (via 4 Task calls in one response) follows the `.claude/skills/think/SKILL.md` pattern for parallel expert analysis.

## Version

**5.0.0** - Multi-agent discovery (BA + PM + SA + TD)
