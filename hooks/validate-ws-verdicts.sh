#!/bin/sh
# Validate docs/ws-verdicts/*.json against schema/ws-verdict.schema.json.
# Used by post-build pipeline hook. Exits 1 if any verdict fails schema validation.
# CWD = project root (set by RunHooks).
set -e
if [ ! -f schema/ws-verdict.schema.json ]; then
  echo "ws-verdict-validate: schema not found" >&2
  exit 1
fi
if ! command -v go >/dev/null 2>&1; then
  echo "ws-verdict-validate: go not found, skipping" >&2
  exit 0
fi
go run ./cmd/sdp-ws-verdict-validate . 2>&1 || exit 1
