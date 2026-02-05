---
name: review
description: Multi-agent quality review (QA + Security + DevOps + SRE + TechLead + Documentation)
tools: Read, Bash, Grep, Task
version: 6.0.0
---

# @review - Multi-Agent Quality Review

Spawn specialist agents for comprehensive quality review.

## Invocation

```bash
@review F01       # Feature ID
@review sdp-xxx   # Beads task ID
```

## Workflow

### Step 1: List Workstreams

**Detect Beads:**
```bash
if bd --version &>/dev/null && [ -d .beads ]; then
  BEADS_ENABLED=true
  bd list --parent {feature-id}
else
  BEADS_ENABLED=false
  ls docs/workstreams/completed/{feature-id}-*.md
fi
```

### Step 2: Spawn Review Agents (PARALLEL)

```python
# Agent 1: QA
Task(
    subagent_type="general-purpose",
    prompt="""You are the QA expert.

Read .claude/agents/qa.md for your specification.

FEATURE: {feature_id}
WORKSTREAMS: {list of completed WS}

Your task:
1. Review test coverage (target: 80%+)
2. Check test quality (pytest/jest/etc)
3. Verify quality metrics (defect density, pass rate)
4. Validate quality gates (entry/exit criteria)

Output:
## QA Review
- Coverage: {percentage}
- Tests: {passing/total}
- Quality metrics: {table}
- Verdict: {PASS/FAIL}

BEADS_INTEGRATION:
If Beads enabled:
- Block workstreams that fail gates
- Update quality metrics in tasks
""",
    description="QA review"
)

# Agent 2: Security
Task(
    subagent_type="general-purpose",
    prompt="""You are the SECURITY expert.

Read .claude/agents/security.md for your specification.

FEATURE: {feature_id}

Your task:
1. Review security controls (auth, input validation, encryption)
2. Check for vulnerabilities (OWASP Top 10)
3. Verify compliance (GDPR/SOC2/etc if applicable)
4. Review secrets management

Output:
## Security Review
- Threats: {mitigated/partial/open}
- Vulnerabilities: {none/low/medium/high}
- Compliance: {status}
- Verdict: {PASS/FAIL}

BEADS_INTEGRATION:
If Beads enabled:
- Create security tasks for gaps
- Track compliance in Beads
""",
    description="Security review"
)

# Agent 3: DevOps
Task(
    subagent_type="general-purpose",
    prompt="""You are the DEVOPS expert.

Read .claude/agents/devops.md for your specification.

FEATURE: {feature_id}

Your task:
1. Review CI/CD pipeline (build, test, deploy)
2. Check infrastructure (Terraform/K8s)
3. Verify deployment strategy (rollback procedures)
4. Check environment management

Output:
## DevOps Review
- CI/CD: {status}
- Infrastructure: {review findings}
- Deployment: {safe/unsafe}
- Verdict: {PASS/FAIL}

BEADS_INTEGRATION:
If Beads enabled:
- Track deployment status in tasks
""",
    description="DevOps review"
)

# Agent 4: SRE
Task(
    subagent_type="general-purpose",
    prompt="""You are the SRE expert.

Read .claude/agents/sre.md for your specification.

FEATURE: {feature_id}

Your task:
1. Review SLOs/SLIs (defined? measurable?)
2. Check monitoring (metrics, logs, traces)
3. Verify incident response procedures
4. Review disaster recovery plan

Output:
## SRE Review
- SLOs: {defined/measurable}
- Monitoring: {coverage}
- Incidents: {procedures}
- Verdict: {PASS/FAIL}

BEADS_INTEGRATION:
If Beads enabled:
- Track SLO compliance in tasks
""",
    description="SRE review"
)

# Agent 5: Tech Lead
Task(
    subagent_type="general-purpose",
    prompt="""You are the TECH LEAD expert.

Read .claude/agents/tech-lead.md for your specification.

FEATURE: {feature_id}
WORKSTREAMS: {list}

Your task:
1. Review code quality (SOLID, clean code)
2. Check architecture decisions (ADRs)
3. Verify team coordination (blockers)
4. Review technical debt

Output:
## Technical Review
- Code quality: {assessment}
- Architecture: {review}
- Blockers: {none/identified}
- Verdict: {PASS/FAIL}

BEADS_INTEGRATION:
If Beads enabled:
- Unblock stuck tasks
- Update tasks with guidance
""",
    description="Technical lead review"
)

# Agent 6: Documentation & Drift
Task(
    subagent_type="general-purpose",
    prompt="""You are the DOCUMENTATION & DRIFT expert.

FEATURE: {feature_id}
WORKSTREAMS: {list of completed WS}

Your task: Check drift at THREE levels

## Level 1: Vision → Specifications
**Question:** Does what we planned match what we wanted?

Find and read:
- PRODUCT_VISION.md (if exists)
- docs/specs/{feature_id}.md (feature spec)
- docs/drafts/idea-{feature_id}.md (original requirements)

Analyze:
1. What were the original business requirements?
2. Did workstreams cover all requirements?
3. Any requirements missed in workstream decomposition?

**Check for gaps:**
- Required features not in any workstream
- User stories not implemented
- Acceptance criteria missing

## Level 2: Specifications → Implementation
**Question:** Does what we built match what we planned?

For each workstream:
```bash
sdp drift detect {ws_id}
```

Check:
- Do all scope_files exist? (NEW files allowed)
- Are all declared entities implemented? (functions, classes, types)
- Does file purpose match documentation?
- Any TODO/FIXME/HACK comments in production code?

## Level 3: Vision → Implementation
**Question:** Did we deliver what we promised?

Final cross-check:
1. Read all completed workstream specs
2. Read actual implementation code
3. Compare to original vision/requirements
4. Identify gaps:
   - Features in vision but not implemented
   - Features implemented but not in vision (scope creep)
   - Quality gaps (security, performance, UX)

Output:
## Documentation & Drift Review

### Level 1: Vision → Specifications
- Vision document: {found/missing}
- Requirements covered: {X%}
- Missing in workstreams: {count/list}
- **Verdict:** {PASS/FAIL}

### Level 2: Specifications → Implementation
- Workstreams checked: {N}
- Drift percentage: {X%} (target: 0%)
- Missing files: {count/list}
- Missing entities: {count/list}
- **Verdict:** {PASS/FAIL}

### Level 3: Vision → Implementation
- Original requirements delivered: {X%}
- Scope creep detected: {yes/no}
- Critical gaps: {count/list}
- **Verdict:** {PASS/FAIL}

### Overall Verdict
**{PASS/FAIL}**

Criteria:
- PASS: All 3 levels PASS, ≤5% drift, no critical gaps
- FAIL: Any level FAIL, >10% drift, or missing critical features

BEADS_INTEGRATION:
If Beads enabled:
- Block workstreams with high drift (>10%)
- Create tasks for missing implementation
- Track gaps in Beads issues
""",
    description="Documentation and drift review"
)
```

### Step 3: Synthesize Verdict

Wait for all 6 agents, then:

```markdown
## Feature Review: {feature_id}

### QA Review
{coverage, tests, metrics, verdict}

### Security Review
{threats, vulnerabilities, compliance, verdict}

### DevOps Review
{CI/CD, infrastructure, deployment, verdict}

### SRE Review
{SLOs, monitoring, incidents, verdict}

### Tech Lead Review
{code quality, architecture, verdict}

### Documentation & Drift Review
**Level 1 (Vision → Specs):** {verdict, coverage%}
**Level 2 (Specs → Code):** {verdict, drift%}
**Level 3 (Vision → Code):** {verdict, gaps}
{overall verdict}

## Overall Verdict

**APPROVED** if all 6 PASS
**CHANGES_REQUESTED** if any FAIL

No middle ground.

## Findings (if CHANGES_REQUESTED)

| Type | Description | Action | Owner |
|------|-------------|--------|-------|
| Bug | {...} | @issue | TBD |
| Tech debt | {...} | @issue | TBD |
| Missing | {...} | New WS | TBD |
```

### Step 4: Post-Review (MANDATORY if CHANGES_REQUESTED)

**Track findings:**
- Bugs → `@issue` → route to `/bugfix`
- Planned work → Add WS to **same feature**
- Tech debt → `@issue` for triage

**Rules:**
- Never create new feature for follow-up
- All findings tracked in Beads (if enabled)

## Output

**Success:**
```
✅ APPROVED
📊 QA: PASS (82% coverage)
🔒 Security: PASS (no vulnerabilities)
⚙️ DevOps: PASS (CI/CD validated)
⏱️ SRE: PASS (SLOs defined)
👨‍💻 TechLead: PASS (code quality good)
📚 Documentation: PASS
   - L1 (Vision → Specs): PASS (100% coverage)
   - L2 (Specs → Code): PASS (0% drift)
   - L3 (Vision → Code): PASS (all delivered)
📌 Beads: {updated if enabled}
```

**Failure:**
```
❌ CHANGES_REQUESTED
📊 QA: FAIL (coverage 65%)
🔒 Security: PASS
⚙️ DevOps: FAIL (no rollback)
⏱️ SRE: PASS
👨‍💻 TechLead: PASS
📚 Documentation: FAIL
   - L1 (Vision → Specs): FAIL (2 requirements missing)
   - L2 (Specs → Code): PASS (3% drift)
   - L3 (Vision → Code): FAIL (scope creep: +3 untracked features)

Findings tracked: {N issues}
```

## Parallel Execution Pattern

6 agents spawned simultaneously (via 6 Task calls) following `.claude/skills/think/SKILL.md` pattern.

## Version

**6.0.0** - Multi-agent review (QA + Security + DevOps + SRE + TechLead + Documentation & Drift)
- Agent 6 checks drift at 3 levels: Vision → Specs → Code
- Validates implementation matches original requirements
