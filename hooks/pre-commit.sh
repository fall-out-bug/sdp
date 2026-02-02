#!/bin/bash
# Pre-commit hook: quality checks on staged files
# SDP v0.7.1+

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

# Check if sdp is available via submodule
if [ -f ".sdp/src/sdp/hooks/pre_commit.py" ]; then
    # Use SDP submodule hooks
    export PYTHONPATH="${REPO_ROOT}/.sdp/src:${PYTHONPATH}"
    python3 .sdp/src/sdp/hooks/pre_commit.py
elif [ -f "src/sdp/hooks/pre_commit.py" ]; then
    # Use local sdp hooks (if sdp is the main repo)
    export PYTHONPATH="${REPO_ROOT}/src:${PYTHONPATH}"
    python3 src/sdp/hooks/pre_commit.py
elif command -v sdp &> /dev/null; then
    # Use installed sdp command
    sdp pre-commit
else
    echo "⚠️  SDP hooks not available, skipping validation"
fi

exit $?
