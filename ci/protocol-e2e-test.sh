#!/usr/bin/env bash
# Protocol E2E test - runs inside Docker container
# Collects all errors and reports at end (no stop-on-first)

set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

ERRORS=()

err() {
  ERRORS+=("$1")
}

# Phase 1: SELF-CONSISTENCY
echo "=== Phase 1: Self-Consistency ==="
MAPPING_COUNT=$(wc -l < .beads-sdp-mapping.jsonl 2>/dev/null || echo 0)
WS_COUNT=$(ls docs/workstreams/backlog/*.md 2>/dev/null | wc -l)
if [ "$MAPPING_COUNT" != "$WS_COUNT" ]; then
  err "beads-mapping-count: mapping=$MAPPING_COUNT, ws-files=$WS_COUNT (MISMATCH)"
fi

# Phase 2: CLI VERIFICATION
echo "=== Phase 2: CLI Verification ==="
# sdp-evidence has no --help (exits 2); verify it runs and prints usage (ignore exit for pipefail)
if ! (sdp-evidence 2>&1 || true) | grep -q "Usage"; then
  err "cli-sdp-evidence: binary failed"
fi
for bin in sdp-guard sdp-orchestrate sdp-ci-loop sdp-eval; do
  if ! $bin --help &>/dev/null; then
    err "cli-$bin: --help failed"
  fi
done

# sdp CLI commands from CLAUDE.md (subset - key commands)
for cmd in "doctor" "status" "init" "parse" "guard activate" "guard check" "guard status" "guard deactivate" \
           "session show" "session clear" "log show" "log trace" "log export" "log stats" \
           "memory index" "memory search" "memory stats" "drift detect" \
           "metrics report" "metrics classify" "telemetry status" "telemetry analyze" \
           "skill list" "skill show" "skill validate"; do
  if ! sdp $cmd --help &>/dev/null 2>&1; then
    err "phantom-cli: sdp $cmd -> exit non-zero"
  fi
done

# Beads (bd --version must succeed; ready/sync may exit 0 or 1)
if ! bd --version &>/dev/null; then
  err "beads: bd --version failed"
  echo "beads-debug: bd --version: $(bd --version 2>&1 || true)"
fi
bd_ready_exit=0; bd ready &>/dev/null || bd_ready_exit=$?
if [ "$bd_ready_exit" -ne 0 ] && [ "$bd_ready_exit" -ne 1 ]; then
  err "beads: bd ready failed (exit $bd_ready_exit)"
  echo "beads-debug: bd ready: $(bd ready 2>&1 || true)"
fi
bd_sync_exit=0; bd sync &>/dev/null || bd_sync_exit=$?
if [ "$bd_sync_exit" -ne 0 ] && [ "$bd_sync_exit" -ne 1 ]; then
  err "beads: bd sync failed (exit $bd_sync_exit)"
  echo "beads-debug: bd sync: $(bd sync 2>&1 || true)"
fi

# Phase 2b: GLOBAL INSTALL + INIT (no local prompts - must fetch or use cache)
echo "=== Phase 2b: Global Install + Init (fresh project) ==="
FRESH_DIR=$(mktemp -d)
if ! (cd "$FRESH_DIR" && git init -q && sdp init --auto 2>/tmp/sdp-init-fresh.log); then
  err "sdp-init-fresh: sdp init --auto failed in fresh project"
  echo "sdp-init-fresh-debug: $(cat /tmp/sdp-init-fresh.log 2>/dev/null | tail -30)"
fi
if [ ! -d "$FRESH_DIR/.claude/skills" ] || [ -z "$(ls -A $FRESH_DIR/.claude/skills 2>/dev/null)" ]; then
  err "sdp-init-fresh: .claude/skills not created (prompts copy failed)"
  echo "sdp-init-fresh-debug: $(cat /tmp/sdp-init-fresh.log 2>/dev/null | tail -30)"
fi
rm -rf "$FRESH_DIR"

# Phase 3: PROTOCOL COMMANDS (happy + negative)
echo "=== Phase 3: Protocol Commands ==="

# sdp-evidence validate (happy)
if ! sdp-evidence validate --require-pr-url=false ci/protocol-e2e-fixtures/valid-evidence.json &>/dev/null; then
  err "sdp-evidence-validate: valid fixture should pass"
fi

# sdp-evidence validate (negative)
if sdp-evidence validate --require-pr-url=false ci/protocol-e2e-fixtures/invalid-evidence.json &>/dev/null; then
  err "sdp-evidence-validate: invalid fixture should fail"
fi

# sdp-evidence inspect
if ! sdp-evidence inspect ci/protocol-e2e-fixtures/valid-evidence.json | grep -q "intent"; then
  err "sdp-evidence-inspect: should show intent section"
fi

# sdp-orchestrate --next-action (F016 exists; also creates checkpoint + runs when none exist)
if ! sdp-orchestrate --feature F016 --next-action 2>/dev/null | grep -qE '"action"|"phase"'; then
  err "sdp-orchestrate: --next-action should output JSON"
fi

# Create feature branch (advance pre-build hook expects it)
git checkout -b feature/F016-e2e 2>/dev/null || git checkout feature/F016-e2e 2>/dev/null || true

# sdp-orchestrate --hydrate
if ! sdp-orchestrate --feature F016 --hydrate --ws 00-016-01 &>/dev/null; then
  err "sdp-orchestrate: --hydrate should succeed"
fi
if [ ! -f .sdp/context-packet.json ]; then
  err "sdp-orchestrate: context-packet.json not created"
fi

# sdp-orchestrate --advance (init→build): creates checkpoint + runs for Phase 4
if ! sdp-orchestrate --feature F016 --advance &>/dev/null; then
  err "sdp-orchestrate: --advance should succeed (init→build)"
fi

# sdp-orchestrate --feature FXXX (negative)
if sdp-orchestrate --feature FXXX --next-action &>/dev/null; then
  err "sdp-orchestrate: non-existent feature should fail"
fi

# sdp-guard: verify binary runs (exit 0=pass or 1=violations both valid)
guard_exit=0
sdp-guard --ws 00-023-01 2>/dev/null || guard_exit=$?
if [ "${guard_exit}" -ne 0 ] && [ "${guard_exit}" -ne 1 ]; then
  err "sdp-guard: unexpected exit ${guard_exit} (expected 0 or 1)"
fi

# Phase 4: TRACING VERIFICATION
echo "=== Phase 4: Tracing ==="
if [ ! -f .sdp/checkpoints/F016.json ]; then
  err "tracing: .sdp/checkpoints/F016.json not created"
fi
if [ ! -d .sdp/runs ] || [ -z "$(ls -A .sdp/runs 2>/dev/null)" ]; then
  err "tracing: .sdp/runs/ should have run files"
fi

# Provenance contract tests (per plan: docs/ARTIFACT_PROVENANCE_HASH_CHAIN_CONTRACT.md)
# Skip if internal/artifact does not exist (package may be added in future)
if [ -d internal/artifact ]; then
  if ! go test ./internal/artifact/... -count=1 &>/dev/null; then
    err "provenance: go test ./internal/artifact/... failed"
  fi
fi

# Phase 5: LLM INTEGRATION (required - validates protocol end-to-end)
echo "=== Phase 5: LLM Integration (opencode) ==="
if [ -z "${GLM_API_KEY:-}" ]; then
  err "llm: GLM_API_KEY required for Phase 5. Add GLM_API_KEY to repo secrets for full E2E (sdp-orchestrate --runtime opencode)."
else
  # Copy E2E fixtures
  mkdir -p docs/workstreams/backlog
  cp ci/protocol-e2e-fixtures/docs/workstreams/backlog/00-999-01.md docs/workstreams/backlog/
  cp ci/protocol-e2e-fixtures/docs/workstreams/backlog/00-999-02.md docs/workstreams/backlog/
  cat ci/protocol-e2e-fixtures/beads-sdp-mapping-e2e.jsonl >> .beads-sdp-mapping.jsonl 2>/dev/null || true

  # Create branch for E2E
  git checkout -b feature/F999-e2e-test 2>/dev/null || git checkout feature/F999-e2e-test 2>/dev/null || true
  git add docs/workstreams/backlog/00-999-*.md .beads-sdp-mapping.jsonl 2>/dev/null || true
  git commit -m "E2E: add F999 fixtures" 2>/dev/null || true

  # Run orchestrate with timeout (15 min; LLM can be slow in CI — API latency, retries)
  if timeout 900 sdp-orchestrate --feature F999 --runtime opencode &>/tmp/e2e-llm.log; then
    if [ ! -f .sdp/checkpoints/F999.json ]; then
      err "llm: checkpoint F999.json not created"
    fi
    if [ ! -f internal/e2e/hello.go ]; then
      err "llm: internal/e2e/hello.go not created by LLM"
    fi
    if [ ! -f internal/e2e/hello_test.go ]; then
      err "llm: internal/e2e/hello_test.go not created by LLM"
    fi
    if ! go test ./internal/e2e/... -count=1 &>/dev/null; then
      err "llm: go test ./internal/e2e/... failed"
    fi
  else
    err "llm: sdp-orchestrate --runtime opencode failed (see /tmp/e2e-llm.log)"
  fi
fi

# Report
echo ""
if [ ${#ERRORS[@]} -gt 0 ]; then
  echo "PROTOCOL E2E FAILED (${#ERRORS[@]} errors)"
  for e in "${ERRORS[@]}"; do
    echo "[ERR] $e"
  done
  echo ""
  echo "=== Debug (for CI investigation) ==="
  echo "beads: which bd=$(which bd 2>/dev/null || echo 'not found'), bd --version=$(bd --version 2>&1 || true)"
  echo "Phase 1: mapping lines=$(wc -l < .beads-sdp-mapping.jsonl 2>/dev/null || echo 0), ws files=$(ls docs/workstreams/backlog/*.md 2>/dev/null | wc -l)"
  echo "Phase 4: .sdp/checkpoints/F016.json exists=$([ -f .sdp/checkpoints/F016.json ] && echo yes || echo no)"
  echo "Phase 4: .sdp/runs: $(ls -la .sdp/runs 2>/dev/null || echo 'dir missing')"
  if [ -f /tmp/e2e-llm.log ]; then
    echo ""
    echo "Phase 5: /tmp/e2e-llm.log (last 100 lines):"
    echo "---"
    tail -100 /tmp/e2e-llm.log
    echo "---"
  fi
  exit 1
fi
echo "Protocol E2E: all phases passed"
exit 0
