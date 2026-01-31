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

## Notes

- The `Critical Checks (Blocking)` job from `ci-critical.yml` is the required status check
- This workflow runs on every PR to main/dev branches
- PRs cannot merge until this check passes
- The warning workflow (`ci-warnings.yml`) is NOT required and doesn't block
