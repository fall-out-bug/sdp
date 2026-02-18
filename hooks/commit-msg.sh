#!/bin/bash
# sdp/hooks/commit-msg.sh
# Git commit-msg hook for conventional commits validation
# and agent metadata trailers for provenance.
# Install: ln -sf ../../sdp/hooks/commit-msg.sh .git/hooks/commit-msg

COMMIT_MSG_FILE=$1
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

echo "🔍 Validating commit message..."

# Conventional commits pattern
# type(scope): description
# Types: feat, fix, docs, test, refactor, style, chore, perf, ci, build
PATTERN="^(feat|fix|docs|test|refactor|style|chore|perf|ci|build)(\([a-z0-9_-]+\))?: .{1,}"

# Also allow merge commits and revert commits
MERGE_PATTERN="^Merge "
REVERT_PATTERN="^Revert "

# Get first line of commit message (ignore Co-authored-by trailers)
FIRST_LINE=$(echo "$COMMIT_MSG" | grep -v "^Co-authored-by:" | head -1)

if echo "$FIRST_LINE" | grep -qE "$REVERT_PATTERN"; then
    echo "✓ Revert commit"
    VALID=1
elif echo "$FIRST_LINE" | grep -qE "$PATTERN"; then
    echo "✓ Valid conventional commit"
    VALID=1
elif echo "$FIRST_LINE" | grep -qE "$MERGE_PATTERN"; then
    echo "✓ Merge commit"
    VALID=1
else
    VALID=0
fi

if [ "$VALID" -ne 1 ]; then
    # Invalid commit message
    echo ""
    echo "❌ Invalid commit message format!"
    echo ""
    echo "Expected: type(scope): description"
    echo "Got:      $FIRST_LINE"
    echo ""
    echo "Valid types:"
    echo "  feat     - New feature"
    echo "  fix      - Bug fix"
    echo "  docs     - Documentation"
    echo "  test     - Tests"
    echo "  refactor - Refactoring"
    echo "  style    - Formatting"
    echo "  chore    - Maintenance"
    echo "  perf     - Performance"
    echo "  ci       - CI/CD"
    echo "  build    - Build system"
    echo ""
    echo "Examples:"
    echo "  feat(lms): WS-060-01 - implement domain layer"
    echo "  test(lms): WS-060-01 - add unit tests"
    echo "  docs(lms): WS-060-01 - execution report"
    echo "  fix(grading): resolve race condition in worker"
    echo ""
    echo "Scope should match feature slug (e.g., lms, grading, api)"
    exit 1
fi

detect_agent() {
    if [ -n "$SDP_AGENT_NAME" ]; then
        printf "%s" "$SDP_AGENT_NAME"
        return
    fi
    if [ -n "$SDP_AGENT" ]; then
        printf "%s" "$SDP_AGENT"
        return
    fi
    if [ "${OPENCODE:-}" = "1" ]; then
        printf "opencode"
        return
    fi
    printf "human"
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

detect_task() {
    if [ -n "$SDP_TASK_ID" ]; then
        printf "%s" "$SDP_TASK_ID"
        return
    fi
    from_msg=$(printf "%s\n" "$COMMIT_MSG" | grep -Eo 'WS-[0-9]{3}-[0-9]{2}|[0-9]{2}-[0-9]{3}-[0-9]{2}|F[0-9]{3,}' | head -1)
    if [ -n "$from_msg" ]; then
        printf "%s" "$from_msg"
        return
    fi
    printf "unknown"
}

append_trailer_if_missing() {
    key="$1"
    value="$2"
    if grep -qiE "^${key}:" "$COMMIT_MSG_FILE"; then
        return
    fi
    if [ "$TRAILER_BLOCK_STARTED" != "1" ]; then
        printf "\n" >> "$COMMIT_MSG_FILE"
        TRAILER_BLOCK_STARTED=1
    fi
    printf "%s: %s\n" "$key" "$value" >> "$COMMIT_MSG_FILE"
}

TRAILER_BLOCK_STARTED=0
append_trailer_if_missing "SDP-Agent" "$(detect_agent)"
append_trailer_if_missing "SDP-Model" "$(detect_model)"
append_trailer_if_missing "SDP-Task" "$(detect_task)"

echo "✓ Added provenance trailers (SDP-Agent, SDP-Model, SDP-Task)"
exit 0
