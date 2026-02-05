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

Your task:
1. Run drift detection on all completed workstreams
2. Verify implementation matches specifications
3. Check for missing files or entities
4. Validate completeness of documentation

For each workstream:
```bash
sdp drift detect {ws_id}
```

Check:
- Do all scope_files exist? (NEW files allowed)
- Are all declared entities implemented? (functions, classes, types)
- Does file purpose match documentation?
- Any TODO/FIXME/HACK comments in production code?

Output:
## Documentation & Drift Review
- Workstreams checked: {N}
- Drift percentage: {X%} (target: 0%)
- Missing files: {count/list}
- Missing entities: {count/list}
- Documentation completeness: {assessment}
- Verdict: {PASS/FAIL}

Drift criteria:
- PASS: ≤5% drift, all critical files present, no blocking gaps
- FAIL: >10% drift, missing critical files, or incomplete documentation

BEADS_INTEGRATION:
If Beads enabled:
- Block workstreams with high drift (>10%)
- Create tasks for missing implementation
- Update drift status in Beads
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
{drift percentage, missing files, documentation completeness, verdict}

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
📚 Documentation: PASS (0% drift, complete)
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
📚 Documentation: FAIL (15% drift, missing files)

Findings tracked: {N issues}
```

## Parallel Execution Pattern

6 agents spawned simultaneously (via 6 Task calls) following `.claude/skills/think/SKILL.md` pattern.

## Version

**6.0.0** - Multi-agent review (QA + Security + DevOps + SRE + TechLead + Documentation)
