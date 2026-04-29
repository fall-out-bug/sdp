#!/bin/sh
# Pre-commit hook: go build, ws-verdict schema validation (if docs/ws-verdicts/*.json changed).
# CWD = repo root. Exit 1 on any failure.
set -e

# 1. go build ./...
go build -tags "sqlite_fts5" ./... || { echo "pre-commit: go build failed" >&2; exit 1; }

# 2. If staged files touch docs/ws-verdicts/*.json — validate
if git diff --cached --name-only | grep -q '^docs/ws-verdicts/.*\.json$'; then
  sh ./scripts/hooks/validate-ws-verdicts.sh || { echo "pre-commit: ws-verdict validation failed" >&2; exit 1; }
fi

exit 0
