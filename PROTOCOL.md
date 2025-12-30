# Cursor Agent Consensus Protocol v1.2
## File-Based Communication for Manual Orchestration

## Overview
This protocol enables consensus between Cursor agents using shared files, with manual user orchestration between separate chat windows.

**Key Changes in v1.2:**
- **Continuous code review** - Code review after each workstream, not just at epic completion
- **Duplication prevention** - Agents must search codebase before implementing new logic
- **Early quality gates** - Tech Lead reviews code after each workstream
- **Cross-epic code review** - Tech Lead reviews across all epics for duplications
- **Early architecture audit** - Architect audits immediately after implementation
- **Code quality verification** - QA verifies code quality before functional testing
- **Strict code review** - Every epic must end with comprehensive code review (DRY, SOLID, Clean Architecture, Clean Code)
- **No error-hiding fallbacks** - Fallbacks that silently hide errors are forbidden

**Key Changes in v1.1:**
- All inbox messages use **JSON format with compact keys** for token efficiency
- Directory structure: `docs/specs/epic_xx/consensus/` (per-epic, not global)
- Decision logs: Markdown files in `consensus/decision_log/[date]-[subject].md`
- No separate state files—state is implicit in artifacts and messages

## Directory Structure
```
workspace/
├── docs/
│   └── specs/
│       └── epic_xx/           # Per-epic specifications
│           ├── epic.md        # Epic definition (human-readable)
│           ├── architecture.md
│           ├── implementation.md
│           ├── testing.md
│           ├── deployment.md
│           └── consensus/
│               ├── artifacts/
│               │   ├── requirements.json  # Analyst output
│               │   ├── architecture.json # Architect output
│               │   └── c4_diagrams.yaml  # Diagrams
│               ├── decision_log/
│               │   └── [date]-[subject].md  # Decision history
│               └── messages/
│                   └── inbox/             # Incoming messages (JSON format)
│                       ├── analyst/
│                       ├── architect/
│                       ├── tech_lead/
│                       ├── developer/
│                       ├── quality/
│                       ├── devops/
│                       ├── security/
│                       └── sre/
```

## Communication Flow

### Step 1: Initialize Epic
User creates initial epic file:
```markdown
# docs/specs/epic_xx/epic.md
## Summary
Add health check endpoint

## Goals
- Health check endpoint for monitoring
...
```

### Step 2: Agent Reads State
Each agent starts by reading:
1. `docs/specs/epic_xx/epic.md` - Epic definition
2. `docs/specs/epic_xx/consensus/messages/inbox/[agent]/` - JSON messages for them
3. `docs/specs/epic_xx/consensus/artifacts/` - Other agents' outputs
4. `docs/specs/epic_xx/consensus/decision_log/` - Previous decisions

### Step 3: Agent Writes Response
Agent writes to their designated locations:
1. `consensus/artifacts/[output].json` - Their main deliverable
2. `consensus/messages/inbox/[target]/[date]-[subject].json` - Messages to other agents (JSON format with compact keys)
3. `consensus/decision_log/[date]-[subject].md` - Log their decision

### Step 4: User Updates State
After each agent run, user checks:
- `docs/specs/epic_xx/consensus/artifacts/` - New artifacts created
- `docs/specs/epic_xx/consensus/messages/inbox/` - New JSON messages
- `docs/specs/epic_xx/consensus/decision_log/` - New decisions logged

State is implicit in artifacts and messages (no separate state file needed).

## Message Format

All inbox messages use **JSON format with compact keys** for token efficiency. Messages are stored as `[date]-[subject].json` files.

### Agent-to-Agent Message Format
```json
{
  "d": "2024-11-19",
  "st": "request",
  "r": "analyst",
  "epic": "EP01",
  "sm": [
    "Requirements ready for architecture review",
    "Health check endpoint needed",
    "Must not impact existing endpoints"
  ],
  "nx": ["review_and_design_architecture"],
  "artifacts": ["docs/specs/epic_01/consensus/artifacts/requirements.json"]
}
```

**Key Schema:**
- `d` - date (YYYY-MM-DD)
- `st` - status/type (request|response|veto|approval|status|handoff|planning)
- `r` - role (from agent: analyst|architect|tech_lead|developer|qa|devops|etc.)
- `epic` - epic_id (EP01, EP02, etc.)
- `sm` - summary (array of strings, concise)
- `nx` - next/action_needed (array of action items)
- `ts` - tests (optional, test commands/results)
- `artifacts` - list of artifact paths (optional)
- `verification` - verification checklist (optional)
- `answers` - Q&A responses (optional, for architect responses)
- `status` - approval status (optional)

### Veto Message Example
```json
{
  "d": "2024-11-19",
  "st": "veto",
  "r": "architect",
  "epic": "EP01",
  "sm": [
    "Layer violation detected: Health check in domain layer",
    "Requirement: Move to infrastructure layer",
    "Blocking: true"
  ],
  "nx": ["revise_requirements"]
}
```

## Decision Log Format
Decisions are logged in markdown files:
```markdown
# Decision Log — EP01 Requirements

## 2025-11-25 — Git hash workflow requirements defined

**Decision:** EP02 extends CLI with `--hash` flag...

**Status:** Approved

**Rationale:** ...

**Consequences:** ...

**Open Questions:**
- Q-EP02-001: ...
```

## User Orchestration Workflow

### 1. Start New Epic
```bash
# Create epic file
echo "epic_id: EP-001
title: 'Your epic title'
status: requirements_gathering
iteration: 1" > consensus/current/epic.yaml

# Clear previous artifacts
rm -f consensus/artifacts/*
rm -f consensus/messages/inbox/*/*
```

### 2. Run Analyst Agent
Open Cursor Chat #1:
```
@analyst Look at docs/specs/epic_xx/epic.md and create requirements.
Write output to docs/specs/epic_xx/consensus/artifacts/requirements.json
Send any messages to docs/specs/epic_xx/consensus/messages/inbox/[agent]/ as JSON files
Log decision to docs/specs/epic_xx/consensus/decision_log/[date]-[subject].md
```

### 3. Run Architect Agent
Open Cursor Chat #2:
```
@architect Review docs/specs/epic_xx/consensus/artifacts/requirements.json
Check docs/specs/epic_xx/consensus/messages/inbox/architect/ for JSON messages
Create architecture in docs/specs/epic_xx/consensus/artifacts/architecture.json
Send feedback to docs/specs/epic_xx/consensus/messages/inbox/analyst/ as JSON if needed
```

### 4. Check for Consensus
```bash
# Count approvals vs vetoes
grep "veto" consensus/current/decision_log.jsonl | wc -l
grep "approve" consensus/current/decision_log.jsonl | wc -l
```

### 5. Handle Conflicts
If conflicts exist, run iteration 2:
```bash
# Update state for iteration 2
echo "iteration: 2" >> consensus/current/state.yaml

# Re-run conflicting agents with context of vetoes
```

## Critical Rules (MUST FOLLOW)

### Language Rule
- **ALL inbox JSON messages MUST be in English. No Russian or any other language.**
- Before writing any message, verify ALL text fields (d, st, r, epic, sm, nx, artifacts, etc.) are in English.
- This applies to ALL agents and ALL message types.

### Inbox Reading Rule
- **ONLY read messages from your own inbox: `messages/inbox/{your_role}/`**
- **DO NOT read other agents' inboxes.**
- Before reading messages, verify you're reading from the correct inbox path for your role.

### Inbox Writing Rule
- **NEVER write messages to your own inbox: `messages/inbox/{your_role}/`**
- **Always write to other agents' inboxes.**
- Before writing a message, verify the target inbox is NOT your own role.

### Documentation Rule
- **After creating/updating artifacts, automatically update relevant documentation.**
- Update: epic.md, architecture.md, implementation.md, testing.md, deployment.md, README.md (if applicable).

### ADR Rule
- **For architectural decisions, create ADR file in `docs/architecture/adr/{date}-{subject}.md`.**
- Use ADR template format.
- Applies to: analyst (when making architectural decisions), architect (always).

### Self-Verification Rule
- **Before completing work, self-verify:**
  - [ ] Clean Architecture boundaries respected
  - [ ] Engineering principles followed (DRY, SOLID, Clean Code, Clean Architecture)
  - [ ] No layer violations
  - [ ] No code review violations (all violations fixed)
  - [ ] No fallbacks hiding errors
  - [ ] Documentation updated

### Code Review Rule
- **Developer**: Conduct incremental code review after each workstream/milestone. Fix all violations before proceeding. Never mark violations as "non-blocking".
- **Tech Lead**: Review code quality after each workstream completion. Veto if violations found. Conduct cross-epic code review at epic completion.
- **Architect**: Audit architecture compliance immediately after implementation, not at epic completion. Veto if violations found.
- **QA**: Verify code quality before functional testing. Check code_review.md exists and violations are fixed.

### Duplication Prevention Rule
- **Before implementing any logic**: [1] Search codebase for existing implementations, [2] Check if existing code can be reused, [3] If duplication found, extract to shared utility, [4] Only implement new code if no existing implementation found.
- **Never duplicate existing logic** without checking codebase first.

### Error Handling Rule
- **Never use fallbacks that hide errors**: All errors must be explicitly logged, raised, or reported. Silent failures are forbidden.
- **Forbidden patterns**: `except: pass`, default values masking exceptions, catch-all hiding errors.

## Agent Prompt Templates

### Analyst Template
```markdown
You are the ANALYST agent. Your epic is in docs/specs/epic_xx/epic.md.

1. Read the epic definition and roadmap context
2. ONLY read messages from docs/specs/epic_xx/consensus/messages/inbox/analyst/ - DO NOT read other agents' inboxes
3. Create requirements in docs/specs/epic_xx/consensus/artifacts/requirements.json:
   - Versioned iteration metadata
   - Traceability to roadmap goals
   - Detailed requirements with acceptance criteria
   - Success metrics, integrations, feature flags
4. Send messages to OTHER agents in docs/specs/epic_xx/consensus/messages/inbox/[other_agent]/ as JSON files (compact format, English-only)
5. NEVER write messages to inbox/analyst/ - always write to other agents' inboxes
6. Log your decision to docs/specs/epic_xx/consensus/decision_log/[date]-[subject].md
7. For architectural decisions, create ADR in docs/architecture/adr/{date}-{subject}.md
8. Update documentation (epic.md, README.md if applicable)
9. Self-verify: Clean Architecture boundaries, Engineering principles, Documentation updated

Focus on business value and minimal intervention. Use JSON format for all inbox messages. ALL messages MUST be in English.
```

### Architect Template (v1.2)
```markdown
You are the ARCHITECT agent.

1. ONLY read messages from docs/specs/epic_xx/consensus/messages/inbox/architect/ - DO NOT read other agents' inboxes
2. Read docs/specs/epic_xx/consensus/artifacts/requirements.json
3. Check docs/specs/epic_xx/consensus/messages/inbox/architect/ for JSON messages
4. Verify Clean Architecture compliance
5. Create architecture in docs/specs/epic_xx/consensus/artifacts/architecture.json:
   - components (with layers)
   - boundaries
   - contracts
6. VETO if layer violations detected (send JSON to inbox/analyst/ - NEVER to inbox/architect/)
7. After developer reports implementation complete: Immediately audit architecture compliance. Veto if violations found.
8. Log decision to docs/specs/epic_xx/consensus/decision_log/[date]-[subject].md
9. Create ADR in docs/architecture/adr/{date}-{subject}.md for architectural decisions
10. Update documentation (architecture.md, README.md if applicable)
11. Self-verify: Clean Architecture boundaries, No layer violations, Code review verified, No fallbacks hiding errors, Documentation updated

Never compromise on Clean Architecture. Use JSON format for all inbox messages. ALL messages MUST be in English.
```

### Developer Template (v1.2)
```markdown
You are the DEVELOPER agent.

1. ONLY read messages from docs/specs/epic_xx/consensus/messages/inbox/developer/ - DO NOT read other agents' inboxes
2. Before implementing any logic: Search codebase for existing implementations to avoid duplication
3. After each workstream: Conduct incremental code review (DRY, SOLID, Clean Code, Clean Architecture). Fix all violations.
4. At epic completion: Conduct strict code review. Create code_review.md. Fix all violations (no "non-blocking").
5. Never use fallbacks that hide errors. All errors must be explicit.
6. Implement according to plan with TDD
7. Update documentation (README.md, CLI_USAGE.md, DSL configs)
8. Self-verify: Clean Architecture boundaries, Engineering principles (DRY, SOLID, Clean Code), No fallbacks hiding errors, Documentation updated

Use JSON format for all inbox messages. ALL messages MUST be in English.
```

### Tech Lead Template (v1.2)
```markdown
You are the TECH LEAD agent.

1. ONLY read messages from docs/specs/epic_xx/consensus/messages/inbox/tech_lead/ - DO NOT read other agents' inboxes
2. Create implementation.md, testing.md, deployment.md
3. After developer reports workstream complete: Review code quality. Veto if violations found.
4. At epic completion: Conduct cross-epic code review for duplications. Document in code_review.md.
5. Send handoff messages to developer, QA, DevOps
6. Self-verify: Code review completed, No fallbacks hiding errors, Documentation updated

Use JSON format for all inbox messages. ALL messages MUST be in English.
```

### QA Template (v1.2)
```markdown
You are the QA agent.

1. ONLY read messages from docs/specs/epic_xx/consensus/messages/inbox/quality/ - DO NOT read other agents' inboxes
2. Before functional testing: Verify code quality (code_review.md exists, violations fixed). Veto if missing.
3. Execute test matrix (unit/integration/e2e/manual)
4. Verify integration parity checks (always required)
5. Check for fallbacks hiding errors. Veto if found.
6. Create test_results.md
7. Self-verify: Code review verified, No fallbacks hiding errors, Documentation updated

Use JSON format for all inbox messages. ALL messages MUST be in English.
```

## Consensus Rules

### Automatic Approval Conditions
- No vetoes in current iteration
- All required agents have responded
- Core artifacts exist (requirements, architecture, plan)

### Automatic Escalation Triggers
- Iteration 3 reached without consensus
- Same veto repeated twice
- Circular dependencies detected

### Veto Priority
1. **Architecture violations** (architect) - Cannot override
2. **Security issues** (quality) - Cannot override
3. **No rollback plan** (devops) - Cannot override
4. **Code review violations** (tech_lead, qa) - Cannot override (all violations must be fixed)
5. **Fallbacks hiding errors** (qa, developer) - Cannot override
6. **Untestable requirements** (tech_lead) - Can negotiate
7. **Scope creep** (analyst) - Can negotiate

### New Veto Triggers (v1.2)
- **code_review_violations_found** (tech_lead) - Code review violations found after workstream
- **dry_violations_detected** (tech_lead) - DRY violations detected
- **solid_violations_detected** (tech_lead) - SOLID violations detected
- **clean_code_violations_detected** (tech_lead) - Clean Code violations detected
- **missing_code_review** (qa) - Code review missing at epic completion
- **code_review_violations_not_fixed** (qa) - Code review violations not fixed
- **fallbacks_hiding_errors** (qa) - Fallbacks that hide errors detected
- **obvious_code_duplications** (qa) - Obvious code duplications found

## Quick Commands for User

### Check Current State
```bash
# See current iteration and phase
cat consensus/current/state.yaml

# Count messages pending
find docs/specs/epic_*/consensus/messages/inbox -name "*.json" | wc -l

# View recent decisions
find docs/specs/epic_*/consensus/decision_log -name "*.md" | tail -5
```

### Move to Next Phase
```bash
# Archive processed messages
mv docs/specs/epic_xx/consensus/messages/inbox/*/*.json docs/specs/epic_xx/consensus/messages/processed/

# Update epic iteration in epic.md or requirements.json
```

### Reset for New Epic
```bash
# Archive current epic
mkdir -p consensus/archive/EP-001
mv consensus/current/* consensus/archive/EP-001/
mv consensus/artifacts/* consensus/archive/EP-001/

# Start fresh
echo "epic_id: EP-002" > consensus/current/epic.yaml
```

## Example Consensus Flow

### Iteration 1
1. User: Creates epic.yaml
2. User: Runs Analyst → requirements.json
3. User: Runs Architect → VETO (layer violation)
4. User: Checks messages, sees veto

### Iteration 2
1. User: Updates epic.yaml (iteration: 2)
2. User: Runs Analyst with veto context → revised requirements
3. User: Runs Architect → architecture.json (approved)
4. User: Runs Tech Lead → plan.json
5. User: All approve → consensus reached

### Implementation (v1.2 - Continuous Code Review)
1. User: Runs Developer → implements workstream 1
2. Developer: Conducts incremental code review → fixes violations
3. Tech Lead: Reviews code quality → approves or vetoes
4. Developer: Implements workstream 2 → incremental code review
5. Tech Lead: Reviews code quality → approves or vetoes
6. ... (repeat for each workstream)
7. Developer: Completes all workstreams → strict code review at epic completion
8. Architect: Audits architecture immediately after implementation → approves or vetoes
9. Tech Lead: Conducts cross-epic code review → documents duplications
10. User: Runs Quality → verifies code quality → functional testing
11. User: Runs DevOps → deployment.json
12. User: Epic complete

## Tips for Efficient Operation

1. **Use Split Terminal**:
   - Left: File explorer showing consensus/
   - Right: Terminal for commands

2. **Agent Shortcuts**: Create aliases
   ```bash
   alias run-analyst="cursor --chat 'Load analyst role from docs/roles/analyst/'"
   alias check-consensus="cat consensus/current/state.yaml"
   ```

3. **Message Templates**: Pre-create message templates
   ```bash
   # Create veto template
   cp templates/veto.yaml consensus/messages/inbox/analyst/
   ```

4. **Batch Operations**: Run related agents together
   - Morning: Analyst + Architect
   - Afternoon: Tech Lead + Developer
   - Evening: Quality + DevOps

## Monitoring Consensus Progress

### Visual Indicator Script
```bash
#!/bin/bash
# consensus_status.sh

echo "=== CONSENSUS STATUS ==="
echo "Epic: $(grep epic_id consensus/current/epic.yaml)"
echo "Iteration: $(grep iteration consensus/current/epic.yaml)"
echo ""
echo "Decisions:"
tail -3 consensus/current/decision_log.jsonl | jq -r '"\(.agent): \(.decision)"'
echo ""
echo "Pending Messages:"
find docs/specs/epic_*/consensus/messages/inbox -name "*.json" -exec basename {} \; | head -5
```

## Success Metrics

Track these manually:
- Time to consensus (iterations × time per iteration)
- Number of vetoes per epic
- Agent response time (when you run each)
- File size growth in artifacts/

---

**Version**: 2.0
**Status**: Ready for Use
**Workflow**: Manual with File-Based Communication
**Last Updated**: 2024-12-31 (Added three-mode architecture: Solo, Structured, Multi-Agent)
