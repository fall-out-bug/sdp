---
ws_id: 00-032-14
feature: F032
status: completed
completed: "2026-01-30"
size: SMALL
project_id: 00
github_issue: null
assignee: Claude
depends_on: ["00-032-12", "00-032-13"]
---

## 00-032-14: Branch Protection Config

### 🎯 Goal

**What must WORK after completing this WS:**
- Документация по настройке branch protection
- Script для автоматической настройки через gh CLI
- main и dev branches защищены

**Acceptance Criteria:**
- [x] AC1: `docs/runbooks/branch-protection-setup.md` created
- [x] AC2: `scripts/setup-branch-protection.sh` created
- [x] AC3: Script requires critical workflow for merge
- [x] AC4: Instructions for manual setup via GitHub UI

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: CI workflow создан, но branch protection не настроен.

**Solution**: Документация + automation script.

### Dependencies

- **00-032-12**: Critical Gate Workflow
- **00-032-13**: Warning Gate Workflow

### Steps

1. **Create runbook**

   ```markdown
   # docs/runbooks/branch-protection-setup.md
   # Branch Protection Setup
   
   ## Overview
   
   This guide explains how to configure branch protection for SDP repositories.
   
   ## Automatic Setup
   
   ```bash
   # Requires: gh CLI authenticated with admin access
   ./scripts/setup-branch-protection.sh
   ```
   
   ## Manual Setup (GitHub UI)
   
   ### Step 1: Navigate to Settings
   
   1. Go to repository Settings
   2. Click "Branches" in sidebar
   3. Click "Add rule" or edit existing
   
   ### Step 2: Configure main branch
   
   **Branch name pattern:** `main`
   
   **Settings:**
   - ✅ Require a pull request before merging
     - ✅ Require approvals: 1
   - ✅ Require status checks to pass before merging
     - ✅ Require branches to be up to date
     - **Required checks:**
       - `Critical Checks (Blocking)` ← from ci-critical.yml
   - ✅ Do not allow bypassing the above settings
   
   ### Step 3: Configure dev branch
   
   **Branch name pattern:** `dev`
   
   **Settings:**
   - ✅ Require status checks to pass before merging
     - **Required checks:**
       - `Critical Checks (Blocking)`
   
   ## Verify Setup
   
   ```bash
   # Check protection status
   gh api repos/{owner}/{repo}/branches/main/protection
   
   # Test by creating failing PR
   git checkout -b test-protection
   echo "test" > test.txt
   git add test.txt
   git commit -m "test: verify protection"
   git push -u origin test-protection
   gh pr create --title "Test protection" --body "Should be blocked"
   ```
   
   ## Troubleshooting
   
   ### Check not appearing
   
   Ensure workflow name matches exactly:
   - Workflow: `name: Critical Quality Gates`
   - Job: `name: Critical Checks (Blocking)`
   
   Required check uses job name, not workflow name.
   
   ### Admin bypass
   
   By default, admins can bypass. To prevent:
   - ✅ Do not allow bypassing the above settings
   ```

2. **Create setup script**

   ```bash
   #!/bin/bash
   # scripts/setup-branch-protection.sh
   # Configure branch protection for SDP repository
   
   set -e
   
   REPO="${GITHUB_REPOSITORY:-$(gh repo view --json nameWithOwner -q .nameWithOwner)}"
   
   echo "Setting up branch protection for: $REPO"
   
   # Main branch protection
   echo "Configuring main branch..."
   gh api -X PUT "repos/$REPO/branches/main/protection" \
     -f required_status_checks='{"strict":true,"contexts":["Critical Checks (Blocking)"]}' \
     -f enforce_admins=true \
     -f required_pull_request_reviews='{"required_approving_review_count":1}' \
     -f restrictions=null \
     -f allow_force_pushes=false \
     -f allow_deletions=false
   
   echo "✅ main branch protected"
   
   # Dev branch protection
   echo "Configuring dev branch..."
   gh api -X PUT "repos/$REPO/branches/dev/protection" \
     -f required_status_checks='{"strict":true,"contexts":["Critical Checks (Blocking)"]}' \
     -f enforce_admins=false \
     -f required_pull_request_reviews=null \
     -f restrictions=null \
     -f allow_force_pushes=false \
     -f allow_deletions=false
   
   echo "✅ dev branch protected"
   
   echo ""
   echo "Branch protection configured successfully!"
   echo "Required status check: 'Critical Checks (Blocking)'"
   ```

### Output Files

- `docs/runbooks/branch-protection-setup.md`
- `scripts/setup-branch-protection.sh`

### Completion Criteria

```bash
# Runbook exists
test -f docs/runbooks/branch-protection-setup.md

# Script exists and is executable
test -x scripts/setup-branch-protection.sh

# Script syntax valid
bash -n scripts/setup-branch-protection.sh
```

---

## Execution Report

**Executed by:** Claude (AI Agent)  
**Date:** 2026-01-30

### Goal Status
- [x] AC1-AC4 — ✅

**Goal Achieved:** YES

### Implementation Details

Created two files:

1. **Runbook** (`docs/runbooks/branch-protection-setup.md`)
   - Overview of branch protection
   - Automatic setup instructions (using script)
   - Manual setup via GitHub UI (step-by-step)
   - Verification commands
   - Troubleshooting section
   - Clear explanation of required checks

2. **Automation Script** (`scripts/setup-branch-protection.sh`)
   - Configures main branch protection:
     - Requires PR with 1 approval
     - Requires "Critical Checks (Blocking)" to pass
     - Enforce admins = true
     - No force pushes/deletions
   - Configures dev branch protection:
     - Requires "Critical Checks (Blocking)" to pass
     - No PR approval requirement
     - No force pushes/deletions
   - Uses gh CLI API
   - Auto-detects repository

### Verification

```bash
$ chmod +x scripts/setup-branch-protection.sh
✅ Script executable

$ bash -n scripts/setup-branch-protection.sh
✅ Script syntax valid

$ test -f docs/runbooks/branch-protection-setup.md
✅ Runbook exists
```

### Notes

- Script requires `gh` CLI with admin permissions
- Can be run manually or in CI/CD
- Branch protection requires "Critical Checks (Blocking)" job from ci-critical.yml
- Warning workflow is NOT required and doesn't block
- Manual setup instructions provided for teams without gh CLI access
