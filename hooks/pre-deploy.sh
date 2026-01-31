#!/bin/bash
# Pre-deploy hook: E2E tests before deployment
# Extracted to Python: src/sdp/hooks/pre_deploy.py
# Usage: ./pre-deploy.sh F{XX} [staging|prod]

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

poetry run python -m sdp.hooks.pre_deploy "$@"
exit $?
