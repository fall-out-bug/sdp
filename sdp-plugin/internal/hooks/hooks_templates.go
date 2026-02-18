package hooks

// Embedded hook templates for when hooks/ directory is not available

func getPreCommitTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Pre-commit
# Part of F065 - Agent Git Safety Protocol

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks"
fi

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

CURRENT_BRANCH=$(git branch --show-current)

# Session validation
if [ -f ".sdp/session.json" ] && command -v jq >/dev/null 2>&1; then
    EXPECTED_BRANCH=$(jq -r '.expected_branch' .sdp/session.json 2>/dev/null)
    if [ -n "$EXPECTED_BRANCH" ] && [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
        echo "ERROR: Branch mismatch! Expected: $EXPECTED_BRANCH, Current: $CURRENT_BRANCH"
        exit 1
    fi
fi

# Protected branch check
if [ -f ".sdp/session.json" ]; then
    case "$CURRENT_BRANCH" in
        main|dev)
            FEATURE_ID=$(jq -r '.feature_id' .sdp/session.json 2>/dev/null)
            if [ -n "$FEATURE_ID" ] && [ "$FEATURE_ID" != "null" ]; then
                echo "ERROR: Cannot commit to $CURRENT_BRANCH for feature $FEATURE_ID"
                exit 1
            fi
            ;;
    esac
fi

exit 0
`
}

func getPrePushTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Pre-push
# Part of F065 - Agent Git Safety Protocol

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks"
fi

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

CURRENT_BRANCH=$(git branch --show-current)

# Prevent pushing to protected branches
case "$CURRENT_BRANCH" in
    main|dev)
        echo "ERROR: Direct push to $CURRENT_BRANCH is not allowed!"
        echo "Create a feature branch and use PR workflow."
        exit 1
        ;;
esac

exit 0
`
}

func getPostCheckoutTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Post-checkout
# Part of F065 - Agent Git Safety Protocol

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks"
fi

# Only update session on branch checkout
if [ "$3" != "1" ]; then
    exit 0
fi

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

if [ -f ".sdp/session.json" ] && command -v jq >/dev/null 2>&1; then
    NEW_BRANCH=$(git branch --show-current)
    TEMP_FILE=$(mktemp)
    jq --arg branch "$NEW_BRANCH" '.expected_branch = $branch' .sdp/session.json > "$TEMP_FILE" 2>/dev/null
    if [ $? -eq 0 ]; then
        mv "$TEMP_FILE" .sdp/session.json
        echo "Session updated: now on branch $NEW_BRANCH"
    else
        rm -f "$TEMP_FILE"
    fi
fi

exit 0
`
}

func getPostMergeTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Post-merge
# Runs after a git merge completes

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks"
    exit 0
fi

# Check if .sdp/ exists
if [ -d ".sdp" ]; then
    echo "SDP: Post-merge checks..."
    # Add post-merge validation here
fi

exit 0
`
}

func getCommitMsgTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Commit message provenance trailers

set -e

COMMIT_MSG_FILE="$1"
[ -n "$COMMIT_MSG_FILE" ] || exit 0
[ -f "$COMMIT_MSG_FILE" ] || exit 0

append_if_missing() {
    key="$1"
    value="$2"
    if grep -qi "^${key}:" "$COMMIT_MSG_FILE"; then
        return
    fi
    if [ "${_SDP_TRAILER_STARTED:-0}" = "0" ]; then
        printf "\n" >> "$COMMIT_MSG_FILE"
        _SDP_TRAILER_STARTED=1
    fi
    printf "%s: %s\n" "$key" "$value" >> "$COMMIT_MSG_FILE"
}

agent="human"
if [ "${OPENCODE:-}" = "1" ]; then
    agent="opencode"
fi

model="unknown"
for key in SDP_MODEL_ID OPENCODE_MODEL ANTHROPIC_MODEL OPENAI_MODEL MODEL; do
    value=$(printenv "$key")
    if [ -n "$value" ]; then
        model="$value"
        break
    fi
done

task="unknown"
if [ -n "${SDP_TASK_ID:-}" ]; then
    task="$SDP_TASK_ID"
fi

append_if_missing "SDP-Agent" "$agent"
append_if_missing "SDP-Model" "$model"
append_if_missing "SDP-Task" "$task"

exit 0
`
}

func getPostCommitTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Post-commit provenance evidence

set -e

if ! command -v sdp >/dev/null 2>&1; then
    exit 0
fi

COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_SHORT=$(git rev-parse --short HEAD)
COMMIT_SUBJECT=$(git log -1 --pretty=%s)

agent="human"
if [ "${OPENCODE:-}" = "1" ]; then
    agent="opencode"
fi

model="unknown"
for key in SDP_MODEL_ID OPENCODE_MODEL ANTHROPIC_MODEL OPENAI_MODEL MODEL; do
    value=$(printenv "$key")
    if [ -n "$value" ]; then
        model="$value"
        break
    fi
done

sdp skill record \
    --skill commit \
    --type generation \
    --ws-id 00-000-00 \
    --data "commit_sha=$COMMIT_HASH" \
    --data "commit_short=$COMMIT_SHORT" \
    --data "commit_subject=$COMMIT_SUBJECT" \
    --data "agent=$agent" \
    --data "model=$model" \
    --data "source=embedded-post-commit" \
    >/dev/null 2>&1 || true

exit 0
`
}
