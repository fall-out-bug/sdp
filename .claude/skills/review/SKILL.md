---
name: review
description: Multi-agent quality review (QA + Security + DevOps + SRE + TechLead + Documentation + Contract Validation)
tools: Read, Bash, Grep, Task
version: 7.0.0
---

# @review - Multi-Agent Quality Review

Spawn specialist agents for comprehensive quality review including contract validation.

## Invocation

```bash
@review F01       # Feature ID
@review sdp-xxx   # Beads task ID
```

## Verbosity Tiers

```bash
@review F01 --quiet     # Exit status only: ✅
@review F01             # Summary: ✅ APPROVED (QA:82%, Security:PASS, 5 agents)
@review F01 --verbose   # Step-by-step progress
@review F01 --debug     # Internal state + API calls
```

**Examples:**

```bash
# Quiet mode
@review F01 --quiet
# Output: ✅

# Default mode
@review F01
# Output: ✅ APPROVED (QA:82%, Security:PASS, 5 agents)

# Verbose mode
@review F01 --verbose
# Output:
# → Spawning 6 review agents...
# → QA review: PASS (82% coverage, 145/145 tests)
# → Security review: PASS (no vulnerabilities)
# → DevOps review: PASS (CI/CD validated)
# → SRE review: PASS (SLOs defined)
# → TechLead review: PASS (code quality good)
# → Documentation review: PASS (0% drift)
# → Contract validation: PASS (0 mismatches)
# ✅ APPROVED

# Debug mode
@review F01 --debug
# Output:
# [DEBUG] Feature: F01
# [DEBUG] Workstreams: 5
# [DEBUG] Beads enabled: true
# [DEBUG] Spawning agents via Task tool...
# [DEBUG] Agent 1: QA (subagent_type=general-purpose)
# [DEBUG] Agent 2: Security (subagent_type=general-purpose)
# [DEBUG] Agent 3: DevOps (subagent_type=general-purpose)
# [DEBUG] Agent 4: SRE (subagent_type=general-purpose)
# [DEBUG] Agent 5: TechLead (subagent_type=general-purpose)
# [DEBUG] Agent 6: Documentation (subagent_type=general-purpose)
# → Spawning 6 review agents...
# [QA agent output...]
# → QA review: PASS (82% coverage, 145/145 tests)
# [Security agent output...]
# → Security review: PASS (no vulnerabilities)
# [DevOps agent output...]
# → DevOps review: PASS (CI/CD validated)
# [SRE agent output...]
# → SRE review: PASS (SLOs defined)
# [TechLead agent output...]
# → TechLead review: PASS (code quality good)
# [Documentation agent output...]
# → Documentation review: PASS (0% drift)
# [DEBUG] Running contract validation...
# → Contract validation: PASS (0 mismatches)
# ✅ APPROVED
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
If Beads enabled (bd --version &>/dev/null):
1. Create issues for quality gaps:
   ```bash
   bd create --title="QA: Coverage {package} below 80%" --type=task --priority=2
   bd create --title="QA: Add tests for {package}" --type=task --priority=1
   ```

2. Block workstreams failing gates:
   ```bash
   bd update beads-{ws-id} --status=blocked --notes="QA: Coverage {X}% < 80%"
   ```

3. Update quality metrics:
   ```bash
   bd update beads-{feature-id} --notes="QA: {coverage}%, {passing}/{total} tests"
   ```
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
If Beads enabled (bd --version &>/dev/null):
1. Create issues for security findings:
   ```bash
   bd create --title="SEC: {vulnerability}" --type=bug --priority=0
   bd create --title="SEC: Hardening {component}" --type=task --priority=1
   ```

2. Track compliance:
   ```bash
   bd update beads-{feature-id} --notes="Security: {OWASP} compliance, {vuln} vulns"
   ```

3. Link security issues to feature:
   ```bash
   bd dep add beads-{sec-issue} beads-{feature-id}
   ```
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
If Beads enabled (bd --version &>/dev/null):
1. Track deployment status:
   ```bash
   bd update beads-{feature-id} --notes="DevOps: CI/CD {status}, deployment {safe/unsafe}"
   ```

2. Create deployment tasks:
   ```bash
   bd create --title="DevOps: Add missing CI workflow" --type=task --priority=2
   bd create --title="DevOps: Fix Goreleaser config" --type=bug --priority=0
   ```
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
If Beads enabled (bd --version &>/dev/null):
1. Track SLO compliance:
   ```bash
   bd update beads-{feature-id} --notes="SRE: Monitoring {X}%, SLOs {defined/not}"
   ```

2. Create incident tasks:
   ```bash
   bd create --title="SRE: Add incident procedures" --type=task --priority=1
   bd create --title="SRE: Implement alerting" --type=task --priority=2
   ```
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
If Beads enabled (bd --version &>/dev/null):
1. Unblock stuck tasks:
   ```bash
   for task in $(bd blocked | grep "blocked by: {feature-id}"); do
       bd update $task --status=ready --notes="Unblocked by tech lead review"
   done
   ```

2. Create technical debt tasks:
   ```bash
   bd create --title="TechDebt: Refactor {component}" --type=task --priority=3
   bd create --title="TechDebt: Add ADR for {decision}" --type=task --priority=2
   ```

3. Update tasks with guidance:
   ```bash
   bd update beads-{ws-id} --notes="TechLead: {guidance on architecture}"
   ```
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
If Beads enabled (bd --version works):
1. Create issues for gaps:
   ```bash
   bd create --title="Drift: {description}" --type=bug --priority=2
   bd create --title="Missing: {file/function}" --type=task --priority=1
   ```

2. Block workstreams with >10% drift:
   ```bash
   for ws in $(sdp drift detect {ws_id} | grep "Drift: >10%"); do
       bd dep add beads-{ws} beads-parent
   done
   ```

3. Update drift status in issues:
   ```bash
   bd update beads-{id} --notes="Drift check: {X}% drift, {N} gaps"
   ```

4. Link findings to parent feature:
   ```bash
   bd dep add beads-{finding-id} beads-{feature-id}
   ```
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

**Track findings in Beads (if enabled):**

```bash
# Bugs
bd create --title="BUG: {description}" --type=bug --priority=1
bd dep add beads-{bug-id} beads-{feature-id}

# Planned work (new workstreams)
bd create --title="WS: {workstream name}" --type=task --priority=2
bd dep add beads-{ws-id} beads-{feature-id}

# Tech debt
bd create --title="TechDebt: {description}" --type=task --priority=3
bd dep add beads-{debt-id} beads-{feature-id}

# Sync to remote
bd sync
```

**Rules:**
- Never create new feature for follow-up
- All findings MUST be tracked in Beads (if enabled)
- Use dependencies: findings depend on feature parent

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
📝 Contract: PASS (0 mismatches)
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

### Step 5: Contract Validation (NEW - CRITICAL QUALITY GATE)

**Purpose:** Ensure implementation matches agreed contract before approval.

**Run contract validation:**
```bash
# Check if contract exists
if [ -f ".contracts/{feature_id}.yaml" ]; then
  echo "✓ Contract found"

  # Run contract validation
  sdp contract validate \
    --contracts .contracts/{feature_id}.yaml \
    --contracts .contracts/{component}-backend.yaml \
    --output .contracts/validation-report.md

  # Check validation result
  error_count=$(grep -c "^|.*ERROR" .contracts/validation-report.md || echo "0")
  warning_count=$(grep -c "^|.*WARNING" .contracts/validation-report.md || echo "0")

  if [ "$error_count" -gt 0 ]; then
    echo "✗ Contract validation failed: $error_count errors found"
    echo "Review: .contracts/validation-report.md"
    return 1  # Block review
  fi

  if [ "$warning_count" -gt 0 ]; then
    echo "⚠ Contract warnings: $warning_count warnings (check report)"
  fi

  echo "✓ Contract validation passed"
else
  echo "⚠ No contract found - skipping validation"
fi
```

**Review contract compliance:**
```bash
# Verify contract lock exists (if contract was synthesized)
if [ -f ".contracts/{feature_id}.yaml" ]; then
  if [ ! -f ".contracts/{feature_id}.yaml.lock" ]; then
    echo "⚠ Contract exists but not locked - consider locking"
  else
    echo "✓ Contract locked"
  fi
fi
```

**Checklist:**
- [ ] Contract validation passed (if contract exists)
- [ ] No endpoint mismatches
- [ ] No schema incompatibilities
- [ ] All integration points verified
- [ ] Contract locked (if synthesized)

**Output:**
```
## Contract Validation
- Contract: {feature_id}.yaml
- Validation: {PASS/SKIP}
- Errors: {count}
- Warnings: {count}
- Verdict: {PASS/FAIL}
```

**CRITICAL:** Contract validation with ERRORs blocks review approval.

## Parallel Execution Pattern

6 agents spawned simultaneously (via 6 Task calls) following `.claude/skills/think/SKILL.md` pattern.

## Version

**7.0.0** - Contract Validation Integration
- Added Step 5: Contract Validation (CRITICAL QUALITY GATE)
- Validates implementation matches agreed contract
- Checks for endpoint mismatches and schema incompatibilities
- Blocks review if contract validation fails
- Integrates with `sdp contract validate` command

**6.0.0** - Multi-agent review (QA + Security + DevOps + SRE + TechLead + Documentation & Drift)
- Agent 6 checks drift at 3 levels: Vision → Specs → Code
- Validates implementation matches original requirements
