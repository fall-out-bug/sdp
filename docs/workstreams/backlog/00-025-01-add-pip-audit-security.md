---
ws_id: 00-025-01
feature: F025
status: backlog
size: SMALL
project_id: 00
github_issue: null
assignee: null
depends_on: []
---

## WS-00-025-01: Add pip-audit Security Scanning

### 🎯 Goal

**What must WORK after completing this WS:**
- pip-audit runs in CI/CD pipeline on every PR/push
- Dependencies with known vulnerabilities block merge
- Security scan results reported in PR comments
- Automated patching via Dependabot configured

**Acceptance Criteria:**
- [ ] AC1: pip-audit added to dev dependencies (pyproject.toml)
- [ ] AC2: pip-audit runs in GitHub Actions workflow (.github/workflows/sdp-quality-gate.yml)
- [ ] AC3: Vulnerability detection fails PR (hard blocking)
- [ ] AC4: PR comments include vulnerability details (CVE, severity, affected packages)
- [ ] AC5: Dependabot config created (.github/dependabot.yml)
- [ ] AC6: Documentation updated (dependency security policy)

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Current State**:
- ✅ Security gate for code-level issues (hardcoded secrets, eval usage)
- ✅ Poetry lock file (pins exact versions)
- ❌ **NO dependency vulnerability scanning**
- ❌ **NO automated dependency updates**
- ❌ **NO security patch notification**

**Risk**:
- Dependencies like PyGithub, python-dotenv, click, pyyaml could have unpatched CVEs
- Team must manually monitor security advisories (high error risk)
- No proactive patching (slow response to zero-day exploits)

**Solution**: Add pip-audit (official PyPA tool) + Dependabot for defense in depth.

### Dependencies

None

### Input Files

- `pyproject.toml` (add pip-audit dependency)
- `.github/workflows/sdp-quality-gate.yml` (existing CI/CD workflow)
- `docs/internals/development.md` (document security policy)

### Steps

1. **Add pip-audit to dev dependencies**

   ```toml
   # pyproject.toml
   [tool.poetry.group.dev.dependencies]
   # ... existing dependencies ...
   pip-audit = "^2.7"  # Official PyPA vulnerability scanning
   ```

   Install: `poetry install --with dev`

2. **Add pip-audit to GitHub Actions**

   ```yaml
   # .github/workflows/sdp-quality-gate.yml
   name: SDP Quality Gates

   on:
     pull_request:
       branches: [main, dev]
     push:
       branches: [main, dev]

   jobs:
     quality-gate:
       runs-on: ubuntu-latest
       timeout-minutes: 10

       steps:
         - uses: actions/checkout@v4
         - uses: actions/setup-python@v5
           with:
             python-version: '3.14'

         - name: Install dependencies
           run: |
             pip install poetry
             poetry install --with dev

         # ... existing quality gates ...

         - name: Run pip-audit security scan
           id: pip-audit
           run: |
             poetry run pip-audit --desc --format json --output audit-report.json
           continue-on-error: false  # Fail immediately on vulnerabilities

         - name: Check pip-audit results
           if: steps.pip-audit.outcome == 'failure'
           run: |
             echo "❌ Vulnerabilities detected:"
             cat audit-report.json
             echo ""
             echo "💡 Run: poetry update to patch vulnerabilities"
             exit 1

         - name: Comment PR with results
           if: github.event_name == 'pull_request'
           uses: actions/github-script@v7
           with:
             script: |
               const fs = require('fs');
               if (fs.existsSync('audit-report.json')) {
                 const report = JSON.parse(fs.readFileSync('audit-report.json', 'utf8'));
                 let vulnCount = report.length;
                 let comment = `## 🔒 Security Scan Results\n\n`;

                 if (vulnCount === 0) {
                   comment += `✅ No vulnerabilities found!\n`;
                 } else {
                   comment += `❌ Found ${vulnCount} vulnerabilities:\n\n`;
                   report.forEach(vuln => {
                     comment += `**${vuln.name} ${vuln.version}**\n`;
                     comment += `- Severity: ${vuln.severity}\n`;
                     comment += `- CVE: ${vuln.id.join(', ')}\n`;
                     comment += `- Fix: \`${vuln.fix_versions[0]}\`\n\n`;
                   });
                   comment += `💡 Run: poetry update to patch\n`;
                 }

                 github.rest.issues.createComment({
                   issue_number: context.issue.number,
                   owner: context.repo.owner,
                   repo: context.repo.repo,
                   body: comment
                 });
               }
   ```

3. **Add Dependabot configuration**

   ```yaml
   # .github/dependabot.yml (new file)
   version: 2
   updates:
     # Main dependencies
     - package-ecosystem: "pip"
       directory: "/"
       schedule:
         interval: "weekly"
         day: "monday"
         time: "09:00"
       open-pull-requests-limit: 3
       reviewers:
         - "fall-out-bug"
       labels:
         - "dependencies"
         - "security"
       commit-message:
         prefix: "chore(deps)"
         include: "scope"
       allow:
         - dependency-type: "direct"
         - dependency-type: "indirect"

     # GitHub Actions
     - package-ecosystem: "github-actions"
       directory: "/"
       schedule:
         interval: "weekly"
   ```

4. **Document dependency security policy**

   Create `docs/internals/development.md` (add section):
   ```markdown
   ## Dependency Management & Security

   ### Vulnerability Scanning
   - **Tool**: pip-audit (official PyPA tool)
   - **Runs**: CI/CD on every PR/push
   - **Failure threshold**: Any vulnerability blocks merge
   - **Exceptions**: Document in SECURTY.md with reason + workaround

   ### Dependency Updates
   - **Automated**: Dependabot creates PRs weekly (Mondays 9am)
   - **Patch versions** (X.Y.Z → X.Y.Z+1): Auto-merge if tests pass
   - **Minor versions** (X.Y → X.Y+1): Manual review required
   - **Major versions** (X → X+1): Create dedicated workstream

   ### Update Policy
   1. **Security patches**: Merge within 24 hours
   2. **Test compatibility**: Run `poetry lock --no-update` before committing
   3. **Update documentation**: Update CHANGELOG.md with dependency changes
   ```

5. **Create SECURITY.md template**

   ```markdown
   # Security Policy

   ## Reporting Vulnerabilities

   If you find a security vulnerability in SDP or its dependencies:

   1. DO NOT open a public issue
   2. Email security advisory to: [maintainer email]
   3. Include: CVE ID (if known), affected versions, reproduction steps

   We will respond within 48 hours with:
   - Confirmation of vulnerability
   - Severity assessment (CRITICAL/HIGH/MEDIUM/LOW)
   - Patch timeline (usually <7 days for CRITICAL/HIGH)

   ## Supported Versions

   | Version | Supported |
   |---------|-----------|
   | v0.5.x  | ✅ Yes      |
   | v0.4.x  | ⚠️ Security fixes only |
   | < v0.4  | ❌ No       |

   ## Dependency Security

   SDP uses pip-audit to scan for known vulnerabilities in dependencies.
   - Runs automatically on every PR via GitHub Actions
   - Blocks merge if vulnerabilities found
   - Automated patching via Dependabot
   ```

### Code

```bash
# Manual testing commands
# Install pip-audit
poetry add --group dev pip-audit

# Run vulnerability scan
poetry run pip-audit

# Check specific dependency
poetry run pip-audit --desc PyGithub

# Generate JSON report
poetry run pip-audit --format json --output audit-report.json

# Ignore specific vulnerability (document in SECURITY.md)
poetry run pip-audit --ignore-vuln <VULN_ID>
```

### Expected Outcome

**After completion:**
- CI/CD pipeline automatically scans dependencies for vulnerabilities
- PRs with vulnerable dependencies are blocked
- Dependabot creates automated patch PRs weekly
- Security policy documented in SECURITY.md
- Supply chain security achieved (Troy Hunt's defense in depth principle)

**Scope Estimate**
- Files: ~4
- Lines: ~200 (SMALL)
- Tokens: ~1000

### Completion Criteria

```bash
# Verify pip-audit installation
poetry show pip-audit

# Run security scan
poetry run pip-audit
# Should pass (no vulnerabilities) or show report

# Test CI/CD workflow
# Create test PR and verify pip-audit step runs

# Verify Dependabot config
yamllint .github/dependabot.yml

# Run tests
pytest tests/unit/test_security.py -v  # If security tests exist
```

### Constraints

- DO NOT ignore vulnerabilities without documentation in SECURITY.md
- DO NOT downgrade security patches (only upgrade)
- DO NOT remove pip-audit from dev dependencies
- DO NOT disable pip-audit in CI/CD (must always run)

---

## Execution Report

**Executed by:** ______
**Date:** ______
**Duration:** ______ minutes

### Goal Status
- [ ] AC1-AC6 — ✅

**Goal Achieved:** ______

### Files Changed
| File | Action | LOC |
|------|--------|-----|
| pyproject.toml | Modify | +2 |
| .github/workflows/sdp-quality-gate.yml | Modify | +50 |
| .github/dependabot.yml | Create | 25 |
| SECURITY.md | Create | 80 |
| docs/internals/development.md | Modify | +50 |

### Statistics
- **Files Changed:** 5
- **Lines Added:** ~207
- **Lines Removed:** ~0
- **Test Coverage:** N/A (infrastructure)
- **Tests Passed:** ______
- **Tests Failed:** ______

### Deviations from Plan
- ______

### Commit
______
