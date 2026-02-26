#!/bin/sh
# Pre-commit hook: go build, optional ws-verdict validation (if docs/ws-verdicts/*.json changed).
# CWD = repo root. Exit 1 on any failure.
# When scripts/hooks/validate-ws-verdicts.sh exists (sdp_dev), runs ws-verdict validation.
set -e

# 1. go build ./...
go build ./... || { echo "pre-commit: go build failed" >&2; exit 1; }

# 2. If staged files touch docs/ws-verdicts/*.json and validate script exists — validate
if git diff --cached --name-only | grep -q '^docs/ws-verdicts/.*\.json$'; then
  if [ -f ./scripts/hooks/validate-ws-verdicts.sh ]; then
    sh ./scripts/hooks/validate-ws-verdicts.sh || { echo "pre-commit: ws-verdict validation failed" >&2; exit 1; }
  fi
fi

exit 0
