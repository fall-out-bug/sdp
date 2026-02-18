#!/bin/bash
# Post-commit hook: auto-comment on GitHub issues
# and emit commit provenance into SDP evidence log.

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

# Get commit info
COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_MSG=$(git log -1 --pretty=%B)
COMMIT_SUBJECT=$(git log -1 --pretty=%s)
COMMIT_SHORT=$(git rev-parse --short HEAD)
REPO_URL=$(git config --get remote.origin.url | sed 's/\.git$//')

extract_trailer() {
    key="$1"
    printf "%s\n" "$COMMIT_MSG" | sed -n "s/^${key}:[[:space:]]*//p" | tail -1
}

detect_model() {
    for key in SDP_MODEL_ID OPENCODE_MODEL ANTHROPIC_MODEL OPENAI_MODEL MODEL; do
        value=$(printenv "$key")
        if [ -n "$value" ]; then
            printf "%s" "$value"
            return
        fi
    done
    printf "unknown"
}

normalize_ws_id() {
    raw="$1"
    if echo "$raw" | grep -qE '^WS-[0-9]{3}-[0-9]{2}$'; then
        printf "00-%s" "${raw#WS-}"
        return
    fi
    if echo "$raw" | grep -qE '^[0-9]{2}-[0-9]{3}-[0-9]{2}$'; then
        printf "%s" "$raw"
        return
    fi
    printf "00-000-00"
}

# Extract task/workstream from trailers or commit message
TASK_ID=$(extract_trailer "SDP-Task")
if [ -z "$TASK_ID" ]; then
    TASK_ID=$(echo "$COMMIT_MSG" | grep -Eo 'WS-[0-9]{3}-[0-9]{2}|[0-9]{2}-[0-9]{3}-[0-9]{2}|F[0-9]{3,}' | head -1)
fi
if [ -z "$TASK_ID" ]; then
    TASK_ID="unknown"
fi

WS_ID=$(echo "$TASK_ID" | grep -Eo 'WS-[0-9]{3}-[0-9]{2}|[0-9]{2}-[0-9]{3}-[0-9]{2}' | head -1)
WS_ID_NORM=$(normalize_ws_id "$WS_ID")

AGENT_ID=$(extract_trailer "SDP-Agent")
if [ -z "$AGENT_ID" ]; then
    if [ "${OPENCODE:-}" = "1" ]; then
        AGENT_ID="opencode"
    else
        AGENT_ID="human"
    fi
fi

MODEL_ID=$(extract_trailer "SDP-Model")
if [ -z "$MODEL_ID" ]; then
    MODEL_ID=$(detect_model)
fi

# Emit commit provenance event (best effort, non-blocking)
if command -v sdp >/dev/null 2>&1; then
    sdp skill record \
        --skill commit \
        --type generation \
        --ws-id "$WS_ID_NORM" \
        --data "commit_sha=$COMMIT_HASH" \
        --data "commit_short=$COMMIT_SHORT" \
        --data "commit_subject=$COMMIT_SUBJECT" \
        --data "agent=$AGENT_ID" \
        --data "model=$MODEL_ID" \
        --data "task_id=$TASK_ID" \
        --data "source=post-commit-hook" \
        >/dev/null 2>&1 || true
fi

if [ -z "$WS_ID" ]; then
    # No WS ID in commit - skip
    # Still allow GitHub comment flow to skip silently.
    WS_ID=""
fi

# Skip GitHub comment if not configured
if [ -z "$GITHUB_TOKEN" ] || [ -z "$GITHUB_REPO" ] || [ -z "$WS_ID" ]; then
    exit 0
fi

# Find WS file
WS_FILE=$(find "$REPO_ROOT/$WS_DIR" -name "${WS_ID}*.md" -type f 2>/dev/null | head -1)

if [ ! -f "$WS_FILE" ]; then
    echo "⚠️  WS file not found: $WS_ID" >&2
    exit 0
fi

# Extract github_issue from frontmatter
ISSUE_NUMBER=$(sed -n 's/^github_issue:[[:space:]]*\([0-9][0-9]*\).*/\1/p' "$WS_FILE" | head -1)

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
