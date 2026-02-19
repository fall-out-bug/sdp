#!/bin/sh
# Post-merge hook: clear Go caches after merge operations
# Part of F063 follow-up - keep local build state aligned

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

if [ "${SDP_SKIP_GO_CACHE_CLEAN:-0}" = "1" ]; then
    exit 0
fi

if command -v go >/dev/null 2>&1; then
    if go clean -cache -testcache >/dev/null 2>&1; then
        echo "Go build/test caches cleared"
    else
        echo "Warning: failed to clear Go caches (go clean -cache -testcache)" >&2
    fi

    TMP_BASE=${TMPDIR:-/tmp}
    for dir in "$TMP_BASE"/go-build*; do
        [ -d "$dir" ] || continue
        rm -rf "$dir" >/dev/null 2>&1 || true
    done
fi

exit 0
