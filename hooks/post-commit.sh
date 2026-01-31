#!/bin/bash
# Post-commit hook: auto-comment on GitHub issues
# Extract WS-XXX-YY from commit message and post comment

set -e

# Configuration
GITHUB_TOKEN="${GITHUB_TOKEN:-}"
GITHUB_REPO="${GITHUB_REPO:-}"
REPO_ROOT=$(git rev-parse --show-toplevel)
WS_DIR="${SDP_WORKSTREAM_DIR:-docs/workstreams}"
if [ ! -d "$REPO_ROOT/$WS_DIR" ]; then
    WS_DIR="workstreams"
fi
if [ ! -d "$REPO_ROOT/$WS_DIR" ]; then
    WS_DIR="tools/hw_checker/docs/workstreams"
fi

# Skip if GitHub not configured
if [ -z "$GITHUB_TOKEN" ] || [ -z "$GITHUB_REPO" ]; then
    # Silent skip - not an error
    exit 0
fi

# Get commit info
COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_MSG=$(git log -1 --pretty=%B)
COMMIT_SHORT=$(git rev-parse --short HEAD)
REPO_URL=$(git config --get remote.origin.url | sed 's/\.git$//')

# Extract WS-XXX-YY from commit message
WS_ID=$(echo "$COMMIT_MSG" | grep -oP 'WS-\d{3}-\d{2}' | head -1)

if [ -z "$WS_ID" ]; then
    # No WS ID in commit - skip
    exit 0
fi

# Find WS file
WS_FILE=$(find "$REPO_ROOT/$WS_DIR" -name "${WS_ID}*.md" -type f 2>/dev/null | head -1)

if [ ! -f "$WS_FILE" ]; then
    echo "⚠️  WS file not found: $WS_ID" >&2
    exit 0
fi

# Extract github_issue from frontmatter
ISSUE_NUMBER=$(grep -oP 'github_issue:\s*\K\d+' "$WS_FILE" || echo "")

if [ -z "$ISSUE_NUMBER" ] || [ "$ISSUE_NUMBER" = "null" ]; then
    # No GitHub issue linked - skip
    exit 0
fi

# Format comment body
COMMENT_BODY=$(cat <<EOF
🔨 **Commit:** [\`${COMMIT_SHORT}\`](${REPO_URL}/commit/${COMMIT_HASH})

\`\`\`
${COMMIT_MSG}
\`\`\`

---
*Auto-posted by SDP post-commit hook*
EOF
)

# Post comment to GitHub issue via API
curl -s -X POST \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "Accept: application/vnd.github.v3+json" \
    "https://api.github.com/repos/${GITHUB_REPO}/issues/${ISSUE_NUMBER}/comments" \
    -d "$(jq -n --arg body "$COMMENT_BODY" '{body: $body}')" \
    > /dev/null

echo "✅ Posted commit comment to issue #${ISSUE_NUMBER} (${WS_ID})"
