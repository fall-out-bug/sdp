#!/bin/sh
# Post-checkout hook: update session file when branch changes
# Part of F065 - Agent Git Safety Protocol

# This hook is called with three parameters:
# $1 = previous HEAD ref
# $2 = new HEAD ref
# $3 = flag indicating whether the checkout was a branch checkout (1) or file checkout (0)

BRANCH_CHECKOUT="$3"

# Only update session on branch checkout (not file checkout)
if [ "$BRANCH_CHECKOUT" != "1" ]; then
    exit 0
fi

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Update session if it exists
if [ -f ".sdp/session.json" ]; then
    # Check if jq is available
    if command -v jq &> /dev/null; then
        NEW_BRANCH=$(git branch --show-current)

        # Update expected_branch in session file
        # Note: This does NOT recalculate the hash - the session is now "synced"
        # If strict hash validation is needed, use `sdp session sync` instead
        TEMP_FILE=$(mktemp)
        jq --arg branch "$NEW_BRANCH" '.expected_branch = $branch' .sdp/session.json > "$TEMP_FILE" 2>/dev/null
        if [ $? -eq 0 ]; then
            mv "$TEMP_FILE" .sdp/session.json
            echo "Session updated: now on branch $NEW_BRANCH"
            echo "Note: Run 'sdp session sync' to recalculate hash if needed"
        else
            rm -f "$TEMP_FILE"
        fi
    fi
fi

# Clear Go build/test caches on branch switch to avoid stale artifacts
if [ "${SDP_SKIP_GO_CACHE_CLEAN:-0}" != "1" ] && command -v go >/dev/null 2>&1; then
    if go clean -cache -testcache >/dev/null 2>&1; then
        echo "Go build/test caches cleared"
    else
        echo "Warning: failed to clear Go caches (go clean -cache -testcache)" >&2
    fi

    (
        TMP_BASE=${TMPDIR:-/tmp}
        for dir in "$TMP_BASE"/go-build* /private/var/folders/*/*/T/go-build* /var/folders/*/*/T/go-build*; do
            [ -d "$dir" ] || continue
            rm -rf "$dir" >/dev/null 2>&1 || true
        done
    ) >/dev/null 2>&1 &
fi

exit 0
