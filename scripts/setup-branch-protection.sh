#!/bin/bash
# scripts/setup-branch-protection.sh
# Configure branch protection for SDP repository

set -e

REPO="${GITHUB_REPOSITORY:-$(gh repo view --json nameWithOwner -q .nameWithOwner 2>/dev/null)}"

if [ -z "$REPO" ]; then
  echo "❌ Error: Could not determine repository"
  echo "Set GITHUB_REPOSITORY or run from a git repository with gh CLI"
  exit 1
fi

echo "Setting up branch protection for: $REPO"

# Main branch protection
echo ""
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
echo ""
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
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Branch protection configured successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Required status check: 'Critical Checks (Blocking)'"
echo ""
echo "To verify:"
echo "  gh api repos/$REPO/branches/main/protection"
echo ""
echo "To test:"
echo "  Create a PR with failing tests and verify it cannot merge"
