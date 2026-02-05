---
name: design
description: Multi-agent system design (Arch + Security + SRE)
tools: Read, Write, Bash, AskUserQuestion, Task
version: 4.0.0
---

# @design - Multi-Agent System Design

Spawn expert agents for architecture + security + reliability design.

## When to Use

- After @feature (requirements complete)
- Before implementation
- Architecture decisions needed
- Security/reliability requirements

## Workflow

### Step 1: Read Feature Spec

**Priority:**
1. Markdown: `Read("docs/drafts/{feature}.md")`
2. Beads (optional): `bd show {feature-id}`

**Detect Beads:**
```bash
if bd --version &>/dev/null && [ -d .beads ]; then
  BEADS_ENABLED=true
else
  BEADS_ENABLED=false
fi
```

### Step 2: Spawn Design Agents (PARALLEL)

```python
# Agent 1: System Architect
Task(
    subagent_type="general-purpose",
    prompt="""You are the SYSTEM ARCHITECT expert.

Read .claude/agents/system-architect.md for your specification.

FEATURE: {feature_spec}

Your task:
1. Design system architecture (pattern, components)
2. Select technology stack
3. Define quality attributes (performance, scalability)
4. Document ADRs (Architecture Decision Records)

Output format:
## System Architecture
{pattern, components, tech stack, quality attributes, ADRs}

BEADS_INTEGRATION:
If Beads enabled:
- Review architecture in feature Beads task
- Link ADRs to workstreams
- Otherwise: Skip
""",
    description="System architecture"
)

# Agent 2: Security
Task(
    subagent_type="general-purpose",
    prompt="""You are the SECURITY expert.

Read .claude/agents/security.md for your specification.

FEATURE: {feature_spec}

Your task:
1. Identify threats (threat modeling)
2. Design authentication/authorization
3. Define data protection (encryption)
4. Ensure compliance (GDPR/SOC2/etc)

Output format:
## Security Assessment
{threat model, auth design, data protection, compliance}

BEADS_INTEGRATION:
If Beads enabled:
- Create security tasks for gaps
- Track compliance in Beads
- Otherwise: Skip
""",
    description="Security design"
)

# Agent 3: SRE
Task(
    subagent_type="general-purpose",
    prompt="""You are the SRE expert.

Read .claude/agents/sre.md for your specification.

FEATURE: {feature_spec}

Your task:
1. Define SLOs/SLIs
2. Design monitoring strategy (metrics, logs, traces)
3. Plan incident response
4. Define disaster recovery

Output format:
## Reliability Strategy
{SLOs, monitoring, incidents, DR}

BEADS_INTEGRATION:
If Beads enabled:
- Create reliability tasks
- Track SLO compliance in Beads
- Otherwise: Skip
""",
    description="Reliability engineering"
)
```

### Step 3: Synthesize Design

Wait for all 3 agents, then combine:

```markdown
## System Design: {feature_name}

### Architecture (from System Architect)
{pattern, components, tech stack, ADRs}

### Security (from Security)
{threats, auth, data protection, compliance}

### Reliability (from SRE)
{SLOs, monitoring, incident response}

### Tradeoffs Analysis
| Aspect | Decision | Rationale |
|--------|----------|-----------|
| Performance | {choice} | {why} |
| Security | {choice} | {why} |
| Scalability | {choice} | {why} |

### Open Questions
{What remains to be decided}

### Next Steps
- Review design
- Execute: @oneshot {feature_id}
```

### Step 4: Update Workstreams

```bash
# Detailed workstreams (with scope_files)
sdp ws create {feature} --detailed

# If Beads enabled: Update with design artifacts
if [ "$BEADS_ENABLED" = true ]; then
  # Design agents already updated Beads
  bd list --parent {feature}  # Verify
fi
```

## Output

**Success:**
```
✅ System design complete
🏗️ Architecture: {pattern}
🔒 Security: {threats mitigated}
⏱️ Reliability: {SLOs defined}
📄 docs/designs/{feature}.md
📌 Beads: {N design tasks if enabled}
```

## Parallel Execution Pattern

3 agents spawned simultaneously (via 3 Task calls) following `.claude/skills/think/SKILL.md` pattern.

## Version

**4.0.0** - Multi-agent design (Arch + Security + SRE)
