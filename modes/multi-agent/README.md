# Multi-Agent Mode

Full consensus workflow with specialized agents, file-based coordination, and explicit handoffs.

## When to Use

- Epic spanning weeks of work
- Multiple people/sessions working in parallel
- Need formal audit trail and documentation
- Compliance or enterprise requirements
- Learning exercise to understand team roles

## How It Works

```
Analyst    →  requirements.md
    ↓
Architect  →  architecture.md + ADRs  (may VETO)
    ↓
Tech Lead  →  implementation.md with workstreams
    ↓
Developer  →  code + tests (per workstream)
    ↓ (continuous review cycle)
QA         →  test-report.md
    ↓
DevOps     →  deployment.md
```

**Key features:**
- Each role has dedicated prompts and responsibilities
- Agents communicate through file-based messages
- Veto protocol prevents bad decisions from propagating
- Full audit trail in decision logs

## Agents and Responsibilities

| Agent | Focus | Output | Can Veto? |
|-------|-------|--------|-----------|
| **Analyst** | Requirements | requirements.md | No |
| **Architect** | System design | architecture.md, ADRs | Yes (architecture violations) |
| **Tech Lead** | Planning | implementation.md | Yes (code quality) |
| **Developer** | Implementation | Code, tests | No |
| **QA** | Verification | test-report.md | Yes (quality issues) |
| **DevOps** | Deployment | deployment.md | Yes (no rollback plan) |

Additional roles for specific needs:
- **Security** - Threat modeling, auth reviews (can veto)
- **SRE** - Observability, SLOs
- **Data/ML** - Data specs, ML artifacts

## Setup

### Directory Structure

```
project/
├── CLAUDE.md
├── docs/
│   └── specs/
│       └── epic-{name}/
│           ├── epic.md                    # Epic definition
│           └── consensus/
│               ├── artifacts/             # Agent outputs
│               │   ├── requirements.md
│               │   ├── architecture.md
│               │   └── ...
│               ├── messages/
│               │   └── inbox/{agent}/     # Inter-agent messages
│               └── decision_log/          # Audit trail
└── prompts/                               # Agent prompts
```

### Agent Prompts

Full prompts are in [/prompts/](/prompts/):
- `analyst_prompt.md` - Requirements analysis
- `architect_prompt.md` - System design
- `tech_lead_prompt.md` - Implementation planning
- `developer_prompt.md` - Code implementation
- `qa_prompt.md` - Quality verification
- `devops_prompt.md` - Deployment

Quick prompts (shorter, for routine tasks) in [/prompts/quick/](/prompts/quick/).

## Workflow

### Phase 1: Requirements

1. Create epic definition: `docs/specs/epic-{name}/epic.md`
2. Run Analyst with full prompt
3. Output: `consensus/artifacts/requirements.md`

```
# In Claude Code or Cursor:
Read prompts/analyst_prompt.md
Analyze docs/specs/epic-auth/epic.md
Create requirements following the protocol
```

### Phase 2: Architecture

1. Run Architect with full prompt
2. Reads requirements, creates design
3. **May VETO** if requirements violate Clean Architecture
4. Output: `consensus/artifacts/architecture.md` + ADRs

```
Read prompts/architect_prompt.md
Review requirements and create architecture
VETO if Clean Architecture would be violated
```

### Phase 3: Planning

1. Run Tech Lead with full prompt
2. Creates workstreams from architecture
3. Output: `consensus/artifacts/implementation.md`

```
Read prompts/tech_lead_prompt.md
Create implementation plan with workstreams
Each workstream should be 2-4 hours of work
```

### Phase 4: Implementation (Iterative)

For each workstream:
1. Developer implements with TDD
2. Developer conducts incremental code review
3. Tech Lead reviews code quality
4. **May VETO** if quality issues found
5. Proceed to next workstream

```
Read prompts/developer_prompt.md
Implement workstream 1 according to plan
Write tests first, then implementation
Conduct code review before reporting complete
```

### Phase 5: Verification

1. Run QA with full prompt
2. Verifies acceptance criteria
3. Runs test suite, checks coverage
4. **May VETO** if quality gates not met

```
Read prompts/qa_prompt.md
Verify all acceptance criteria from requirements
Run tests, check coverage
Report issues found
```

### Phase 6: Deployment

1. Run DevOps with full prompt
2. Creates deployment configuration
3. **May VETO** if no rollback plan

```
Read prompts/devops_prompt.md
Create deployment plan
Include rollback procedure
```

## Veto Protocol

Vetoes are **non-negotiable** quality gates:

| Veto Type | Agent | Trigger |
|-----------|-------|---------|
| Architecture violation | Architect | Clean Architecture broken |
| Security issue | Security | Vulnerabilities found |
| Missing rollback | DevOps | No recovery plan |
| Code quality | Tech Lead, QA | Review violations unfixed |

**When vetoed:**
1. Read veto message in your inbox
2. Fix the issue
3. Re-submit for review
4. Document in decision log

## File-Based Communication

### Message Format

JSON messages in `consensus/messages/inbox/{agent}/`:

```json
{
  "d": "2024-01-15",
  "r": "architect",
  "st": "approved",
  "epic": "epic-auth",
  "sm": "Architecture approved, proceed to planning",
  "nx": ["tech_lead"],
  "artifacts": ["consensus/artifacts/architecture.md"]
}
```

### Inbox Rules

- **READ**: Only your own inbox
- **WRITE**: Only to OTHER agents' inboxes

This prevents confusion and creates clear audit trail.

## When NOT to Use Multi-Agent

Multi-Agent is overkill for:
- Tasks < 1 week of work
- Single-person projects
- Rapid prototyping
- No compliance requirements

Use [Solo Mode](../solo/) or [Structured Mode](../structured/) instead.

## Protocol Reference

For complete protocol specification:
- [PROTOCOL.md](/PROTOCOL.md) - Full consensus protocol v1.2
- [RULES_COMMON.md](/RULES_COMMON.md) - Shared rules for all agents
- [USER_GUIDE.md](/USER_GUIDE.md) - Operational guide

## Tips

### Parallel Execution
Some agents can run in parallel:
- Analyst and initial research
- Developer (different workstreams, if independent)
- QA and DevOps preparation

### Context Switching
Clear context between agent roles:
- In Claude Code: Start new session or use `/clear`
- In Cursor: Use separate chat windows

### Model Selection
Match model capability to role:
- Analyst, Architect, Security: Most capable model
- Developer, QA, DevOps: Fast model is sufficient

### Keep Decision Logs
Every significant decision should be logged:
```markdown
# 2024-01-15 - Architecture Decision

**Decision**: Use PostgreSQL over MongoDB
**Agent**: Architect
**Reason**: Need ACID compliance for financial data
**ADR**: docs/adr/0003-use-postgresql.md
```
