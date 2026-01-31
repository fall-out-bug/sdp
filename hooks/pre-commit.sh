#!/bin/bash
# Pre-commit hook: quality checks on staged files
# Extracted to Python: src/sdp/hooks/pre_commit.py

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

poetry run python -m sdp.hooks.pre_commit
exit $?
