#!/bin/sh
# Post-merge hook: clear Go caches and run doc-sync for architectural changes
# Part of F063 follow-up - keep local build state aligned
# Part of sdplab-665 - auto-fix documentation on merge

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Skip if requested
if [ "${SDP_SKIP_POST_MERGE:-0}" = "1" ]; then
    exit 0
fi

# Clear Go caches if requested
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

# Run doc-sync fix for architectural changes
# This automatically fixes documentation inconsistencies when merging changes
if command -v sdp-doc-sync >/dev/null 2>&1; then
    if [ "${SDP_SKIP_DOC_SYNC:-0}" != "1" ]; then
        echo "Running sdp-doc-sync fix for architectural changes..."
        if sdp-doc-sync --mode fix 2>&1 | grep -q "nothing to fix"; then
            echo "Documentation is consistent"
        else
            echo "Documentation inconsistencies fixed automatically"
        fi
    fi
fi

exit 0
