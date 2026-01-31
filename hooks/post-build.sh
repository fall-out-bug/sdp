#!/bin/bash
# Post-build hook: quality checks after workstream execution
# Extracted to Python: src/sdp/hooks/post_build.py
# Usage: ./post-build.sh WS-ID [module_path]

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

poetry run python -m sdp.hooks.post_build "$@"
exit $?
