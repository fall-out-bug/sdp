#!/bin/bash
# Pre-push hook: regression tests before pushing
# Extracted to Python: src/sdp/hooks/pre_push.py

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

poetry run python -m sdp.hooks.pre_push
exit $?
